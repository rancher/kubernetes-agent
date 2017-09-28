package events

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rancher/go-rancher/v3"
)

var slashRegex = regexp.MustCompile("[/]{2,}")

// EventHandler Defines the function "interface" that handlers must conform to.
type EventHandler func(*Event, *client.RancherClient) error

type EventRouter struct {
	apiClient     *client.RancherClient
	subscribeURL  string
	eventHandlers map[string]EventHandler
	workerCount   int
	eventStream   *websocket.Conn
	PingConfig    PingConfig
}

func NewEventRouter(apiClient *client.RancherClient, workerCount int, eventHandlers map[string]EventHandler) (*EventRouter, error) {
	subscribeURL := ""

	if apiClient.RancherBaseClient != nil {
		schema, ok := apiClient.GetTypes()["subscribe"]
		if !ok {
			return nil, errors.New("Client is not able to subscribe to events")
		}

		subscribeURL = schema.Links["collection"]
		if strings.HasPrefix(subscribeURL, "http") {
			subscribeURL = strings.Replace(subscribeURL, "http", "ws", 1)
		}
	}

	return &EventRouter{
		apiClient:     apiClient,
		subscribeURL:  subscribeURL,
		eventHandlers: eventHandlers,
		workerCount:   workerCount,
		PingConfig:    DefaultPingConfig,
	}, nil
}

func (router *EventRouter) StartHandler(name string, ready chan<- bool) error {
	wp := SkippingWorkerPool(router.workerCount, resourceIDLocker)
	return router.run(wp, ready, ";handler="+name)
}

func (router *EventRouter) Start(ready chan<- bool) error {
	wp := SkippingWorkerPool(router.workerCount, resourceIDLocker)
	return router.run(wp, ready, "")
}

func (router *EventRouter) RunWithWorkerPool(wp WorkerPool) error {
	return router.run(wp, nil, "")
}

func (router *EventRouter) run(wp WorkerPool, ready chan<- bool, eventSuffix string) (err error) {

	log.WithFields(log.Fields{
		"workerCount": router.workerCount,
	}).Info("Initializing event router")

	handlers := map[string]EventHandler{}

	if pingHandler, ok := router.eventHandlers["ping"]; ok {
		// Ping doesnt need registered in the POST and ping events don't have the handler suffix.
		//If we start handling other non-suffix events, we might consider improving this.
		handlers["ping"] = pingHandler
	} else {
		handlers["ping"] = DropEvent
	}

	subscribeParams := url.Values{}
	for event, handler := range router.eventHandlers {
		fullEventKey := event + eventSuffix
		subscribeParams.Add("eventNames", fullEventKey)
		handlers[fullEventKey] = handler
	}

	accessKey := ""
	secretKey := ""
	if router.apiClient.RancherBaseClient != nil {
		accessKey = router.apiClient.GetOpts().AccessKey
		secretKey = router.apiClient.GetOpts().SecretKey
	}

	eventStream, err := router.subscribeToEvents(router.subscribeURL, accessKey, secretKey, subscribeParams)
	if err != nil {
		return err
	}
	log.Info("Connection established")
	router.eventStream = eventStream
	defer router.Stop()

	if ready != nil {
		ready <- true
	}

	ph := newPongHandler(router)
	defer ph.stop()
	router.eventStream.SetPongHandler(ph.handle)
	go router.sendWebsocketPings()

	for {
		_, message, err := router.eventStream.ReadMessage()
		if err != nil {
			// Error here means the connection is closed. It's normal, so just return.
			return nil
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		event := &Event{}
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.WithFields(log.Fields{
				"message": string(message),
			}).Warnf("Error parsing message: %s", err)
			continue
		}
		wp.HandleWork(event, handlers, router.apiClient)
	}
}

func (router *EventRouter) Stop() {
	router.eventStream.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
	router.eventStream.Close()
}

func (router *EventRouter) subscribeToEvents(subscribeURL string, accessKey string, secretKey string, data url.Values) (*websocket.Conn, error) {
	// gorilla websocket will blow up if the path starts with //
	parsed, err := url.Parse(subscribeURL)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(parsed.Path, "//") {
		parsed.Path = slashRegex.ReplaceAllString(parsed.Path, "/")
		subscribeURL = parsed.String()
	}

	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 30,
	}
	headers := http.Header{}
	headers.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(accessKey+":"+secretKey)))
	subscribeURL = subscribeURL + "?" + data.Encode()
	ws, resp, err := dialer.Dial(subscribeURL, headers)

	if err != nil {
		log.WithFields(log.Fields{
			"subscribeUrl": subscribeURL,
		}).Errorf("Error subscribing to events: %s", err)
		if resp != nil {
			log.WithFields(log.Fields{
				"status":          resp.Status,
				"statusCode":      resp.StatusCode,
				"responseHeaders": resp.Header,
			}).Error("Got error response")
			if resp.Body != nil {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				log.Errorf("Error response: %s", body)
			}
		}
		if ws != nil {
			ws.Close()
		}
		return nil, err
	}
	return ws, nil
}

func (router *EventRouter) GetWebSocketConn() *websocket.Conn {
	return router.eventStream
}
