package kubernetesevents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-model/model"
)

var (
	pathTemplates = []string{"/api/v1/%s", "/apis/extensions/v1beta1/%s"}
	waits         = []int{0, 1, 2, 4, 8, 16, 0}
)

type Handler interface {
	Handle(event model.WatchEvent) error
	GetKindHandled() string
}

func SyncAndWatchEventStream(handlers []SyncHandler) error {
	doneChan := make(chan error)
	for _, handler := range handlers {
		fifo := NewDeltaFIFO(handler, doneChan)
		go fifo.Process()
	}
	return <-doneChan
}

func ConnectToEventStream(handlers []Handler, conf config.Config) error {
	log.Infof("Starting kubernetes event listener configuration: %+v", conf)
	baseUrl := conf.KubernetesURL
	baseUrl = strings.Replace(baseUrl, "http", "ws", 1)

	doneChan := make(chan error)

	for _, handler := range handlers {
		dialer := &websocket.Dialer{}
		headers := http.Header{}
		headers.Add("Origin", "http://kubernetes-agent")

	outer:
		for idx, wait := range waits {
			for _, template := range pathTemplates {
				url := buildURL(baseUrl, handler.GetKindHandled(), template)
				log.WithFields(log.Fields{"url": url}).Info("Connecting to event stream.")

				ws, _, err := dialer.Dial(url, headers)
				if err != nil {
					if idx < len(waits)-1 {
						if idx > 0 {
							log.Warnf("Error connecting to %s. Try %v of %v. Will wait %v seconds and try again. Error: %#v", url, idx, len(waits), wait, err)
						}
						time.Sleep(time.Second * time.Duration(wait))
						continue
					} else {
						log.Error("Failed to connet to %s. Giving up. Error: %#v", url, err)
						return err
					}
				} else {
					go readMessages(ws, url, doneChan, handler)
					break outer
				}
			}
		}
	}

	err := <-doneChan
	return err
}

func readMessages(ws *websocket.Conn, url string, rc chan<- error, handler Handler) (e error) {
	defer func() {
		rc <- e
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error reading from websocket for [%v]: %s", url, err)
		}

		var event model.WatchEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("Error parsing event: %s", err)
		}
		log.Infof("Received event: [%s]", msg)

		err = handler.Handle(event)
		if err != nil {
			log.Errorf("Error handling event: %#v", err)
		}
	}
}

func buildURL(baseUrl, resource, template string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		// Fatal logging. Will cause program exit
		log.WithFields(log.Fields{"error": err, "baseUrl": baseUrl}).Fatal("Couldn't parse URL.")
	}
	path := fmt.Sprintf(template, resource)
	u.Path = path
	q := u.Query()
	q.Set("watch", "true")
	u.RawQuery = q.Encode()
	return u.String()
}
