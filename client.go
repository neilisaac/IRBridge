package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/tarm/serial"
)

func runClient(cmd *cobra.Command, args []string) error {
	logger := log15.New()
	addr := args[0]
	device, err := cmd.Flags().GetString("device")
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(addr, http.Header{})
	if err != nil {
		return err
	}

	commands := make(chan command)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				logger.Error("websocket ping failed", "error", err)
				close(commands)
				return
			}
		}
	}()

	go func() {
		for {
			c := command{}
			if err := conn.ReadJSON(&c); err != nil {
				logger.Error("failed to read json", "error", err)
				close(commands)
				break
			}
			commands <- c
		}
	}()

	dev, err := serial.OpenPort(&serial.Config{Name: device, Baud: 9600})
	if err != nil {
		// return err
	}
	defer dev.Close()

	for c := range commands {
		if c.Repeat == 0 {
			c.Repeat = 1
		}
		logger.Info("sending command", "code", c.Code, "times", c.Repeat)
		for i := 0; i < c.Repeat; i++ {
			if i > 0 {
				time.Sleep(10 * time.Millisecond)
			}
			if _, err := fmt.Fprintf(dev, "%s\n", c.Code); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("websocket closed")
}
