package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inconshreveable/log15"
	"github.com/juju/ratelimit"
	"github.com/julienschmidt/httprouter"
)

type command struct {
	Code   string `json:"code"`
	Repeat int    `json:"repeat,omitempty"`
}

type publishChannel chan command
type subscriptions map[publishChannel]struct{}

type dispatch struct {
	Bucket        *ratelimit.Bucket
	Subscriptions map[string]subscriptions
	Logger        log15.Logger
	mutex         sync.Mutex
}

func (d *dispatch) subscribe(id string, c publishChannel) {
	d.Logger.Info("subscribe", "id", id)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.Subscriptions[id]; !exists {
		d.Subscriptions[id] = subscriptions{}
	}

	d.Subscriptions[id][c] = struct{}{}
}

func (d *dispatch) unsubscribe(id string, c publishChannel) {
	d.Logger.Info("unsubscribe", "id", id)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.Subscriptions[id], c)

	if len(d.Subscriptions[id]) == 0 {
		delete(d.Subscriptions, id)
	}
}

func (d *dispatch) publish(id string, data command) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	count := 0
	for pub := range d.Subscriptions[id] {
		pub <- data
		count++
	}

	d.Logger.Info("publish", "id", id, "data", data, "subscribers", count)
}

func (d *dispatch) Send(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if d.Bucket.TakeAvailable(1) != 1 {
		http.Error(w, "rate limit reached", http.StatusTooManyRequests)
		return
	}

	d.publish(params.ByName("id"), command{Code: params.ByName("code")})
}

func (d *dispatch) Trigger(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if d.Bucket.TakeAvailable(1) != 1 {
		http.Error(w, "rate limit reached", http.StatusTooManyRequests)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	codes := map[string]string{
		"on":                "dac-on",
		"off":               "dac-off",
		"select tv":         "dac-opt1",
		"select chromecast": "dac-opt2",
		"select usb":        "dac-usb",
	}

	code := string(data)
	if c, ok := codes[code]; ok {
		code = c
	}

	d.publish(params.ByName("id"), command{Code: code})
}

func (d *dispatch) Subscribe(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if d.Bucket.TakeAvailable(1) != 1 {
		http.Error(w, "rate limit reached", http.StatusTooManyRequests)
		return
	}

	conn, err := websocket.Upgrade(w, r, http.Header{}, 1024, 1024)
	if err != nil {
		d.Logger.Error("failed to upgrade websocket", "error", err)
		http.Error(w, "failed to upgrade to websocket", http.StatusBadRequest)
		return
	}

	c := make(chan command)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				d.Logger.Warn("ping failed", "error", err)
				close(c)
				return
			}
		}
	}()

	d.subscribe(params.ByName("id"), c)
	defer d.unsubscribe(params.ByName("id"), c)

	for data := range c {
		if err := conn.WriteJSON(data); err != nil {
			d.Logger.Error("failed to write json", "error", err)
			return
		}
	}
}
