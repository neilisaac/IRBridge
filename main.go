package main

import (
	"github.com/spf13/cobra"
)

func main() {
	server := &cobra.Command{
		Use:   "server",
		Short: "run server",
		RunE:  runServer,
		Args:  cobra.ExactArgs(0),
	}
	server.Flags().Int("port", 8080, "port to listen on")

	client := &cobra.Command{
		Use:   "client ws://server:port/subscribe/client_uuid",
		Short: "run client",
		RunE:  runClient,
		Args:  cobra.ExactArgs(1),
	}
	client.Flags().String("device", "/dev/ttyACM0", "serial device")

	cmd := &cobra.Command{Use: "IRBridge"}
	cmd.AddCommand(server)
	cmd.AddCommand(client)
	cmd.Execute()
}
