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

const pathTemplate string = "/api/v1/%s"

var waits []int = []int{1, 2, 4, 4, 4, 0}

type Handler interface {
	Handle(event model.WatchEvent) error
	GetKindHandled() string
}

func ConnectToEventStream(handlers []Handler, conf config.Config) error {
	log.Infof("Starting kubernetes event listener configuration: %+v", conf)
	baseUrl := conf.KubernetesURL
	baseUrl = strings.Replace(baseUrl, "http", "ws", 1)

	doneChan := make(chan error)

	for _, handler := range handlers {
		url := buildURL(baseUrl, handler.GetKindHandled())
		log.WithFields(log.Fields{"url": url}).Info("Connecting to event stream.")

		dialer := &websocket.Dialer{}
		headers := http.Header{}
		headers.Add("Origin", "http://kubernetes-agent")

		for idx, wait := range waits {
			ws, _, err := dialer.Dial(url, headers)
			if err != nil {
				if idx < len(waits)-1 {
					log.Warn("Error connecting to %s. Try %v of %v. Will wait %v and try again. Error: %#v", url, idx, len(waits), wait, err)
					time.Sleep(time.Second * time.Duration(wait))
					continue
				} else {
					log.Error("Failed to connet to %s. Giving up. Error: %#v", url, err)
					return err
				}
			} else {
				go readMessages(ws, doneChan, handler)
			}
		}
	}

	err := <-doneChan
	return err
}

func readMessages(ws *websocket.Conn, rc chan<- error, handler Handler) (e error) {
	defer func() {
		rc <- e
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error reading from websocket: %s", err)
		}

		var event model.WatchEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("Error parsing event: %s", err)
		}
		log.Infof("Received event: [%s]", msg)

		err = handler.Handle(event)
		if err != nil {
			return fmt.Errorf("Error handling event: %s", err)
		}

	}
}

func buildURL(baseUrl string, resource string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		// Fatal logging. Will cause program exit
		log.WithFields(log.Fields{"error": err, "baseUrl": baseUrl}).Fatal("Couldn't parse URL.")
	}
	path := fmt.Sprintf(pathTemplate, resource)
	u.Path = path
	q := u.Query()
	q.Set("watch", "true")
	u.RawQuery = q.Encode()
	return u.String()
}
