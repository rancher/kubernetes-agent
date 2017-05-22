package kubernetesevents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

type DeltaFIFO struct {
	l sync.RWMutex
	c sync.Cond

	items map[string]model.WatchEvent
	queue []string

	handler  SyncHandler
	doneChan chan error
}

func NewDeltaFIFO(handler SyncHandler, doneChan chan error) *DeltaFIFO {
	dF := &DeltaFIFO{
		handler:  handler,
		doneChan: doneChan,
		items:    map[string]model.WatchEvent{},
		queue:    []string{},
	}

	dF.c.L = &dF.l
	return dF
}

func (d *DeltaFIFO) startProcessing() {
	for {
		event := d.Pop()
		resource, err := d.handler.Decode(event)
		if err != nil {
			continue
		}
		switch event.Type {
		case "MODIFIED":
			fallthrough
		case "ADDED":
			err := d.handler.Add(resource)
			if err != nil {
				log.Errorf("Error Processing event %v", err)
				// Push this event back to the end of the queue?
				continue
			}
		case "DELETED":
			err := d.handler.Delete(resource)
			if err != nil {
				log.Errorf("Error Processing event %v", err)
				// Push this event back to the end of the queue?
				continue
			}
		}
	}
}

//thread safe add
func (d *DeltaFIFO) Add(event model.WatchEvent) error {
	d.l.Lock()
	defer d.l.Unlock()
	key, err := d.handler.GetKey(event)
	if err != nil {
		return err
	}
	if _, ok := d.items[key]; !ok {
		d.queue = append(d.queue, key)
	}
	d.items[key] = event
	d.c.Broadcast()
	return nil
}

//blocks until a value is available
func (d *DeltaFIFO) Pop() model.WatchEvent {
	d.l.Lock()
	defer d.l.Unlock()
	if len(d.queue) == 0 {
		d.c.Wait()
	}
	key := d.queue[0]
	d.queue = d.queue[1:]
	val, ok := d.items[key]
	delete(d.items, key)
	if !ok {
		return model.WatchEvent{}
	}
	return val
}

func (d *DeltaFIFO) Process() {
	listURL := d.handler.GetListURL()
	watchURL := d.handler.GetWatchURL()

	go d.startProcessing()

	go func(done chan error) {
		wait := 1
		dialer := &websocket.Dialer{
			TLSClientConfig: kubernetesclient.GetTLSClientConfig(),
		}
		headers := http.Header{}
		headers.Add("Origin", "http://kubernetes-agent")
		headers.Add("Authorization", kubernetesclient.GetAuthorizationHeader())

		var err error
		var ws *websocket.Conn

		for {
			ws, _, err = dialer.Dial(watchURL, headers)
			if err != nil {
				if wait > 16 {
					d.doneChan <- fmt.Errorf("Error connecting to %s. Giving up. Err: %v", watchURL, err)
					return
				}
				wait = wait * 2
				time.Sleep(time.Second * time.Duration(wait))
				continue
			}
			break
		}
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				d.doneChan <- fmt.Errorf("Error reading ws message %v", err)
				return
			}
			var event model.WatchEvent
			err = json.Unmarshal(msg, &event)
			if err != nil {
				d.doneChan <- fmt.Errorf("Error unmarshalling event %v", err)
				return
			}
			d.Add(event)
		}
	}(d.doneChan)

	listClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: kubernetesclient.GetTLSClientConfig(),
		},
	}
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		log.Errorf("Error creating list request %v", err)
		d.doneChan <- err
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", kubernetesclient.GetAuthorizationHeader())

	resp, err := listClient.Do(req)
	if err != nil {
		log.Errorf("Error requesting for list of objects %v", err)
		d.doneChan <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Errorf("Couldn't retrieve list of objects, status code: %d", resp.StatusCode)
		d.doneChan <- fmt.Errorf("Couldn't retrieve list of objects, status code: %d", resp.StatusCode)
		return
	}
	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error decoding response %v", err)
		d.doneChan <- fmt.Errorf("Error decoding response %v", err)
		return
	}

	eventObj := map[string]interface{}{}
	err = json.Unmarshal(byteContent, &eventObj)
	if err != nil {
		log.Errorf("Error unmarshalling json %v", err)
		d.doneChan <- fmt.Errorf("Error unmarshalling json %v", err)
		return
	}

	items, ok := eventObj["items"].([]interface{})
	if ok {
		for _, obj := range items {
			var event model.WatchEvent
			event.Object = obj.(map[string]interface{})
			event.Type = "ADDED"
			d.Add(event)
		}
	}

}
