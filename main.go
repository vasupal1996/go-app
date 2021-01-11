package main

import (
	"go-app/server"
	"net"
	"os"
	"os/signal"
)

func main() {
	s := server.NewServer()
	s.StartServer()
	defer s.StopServer()

	err := sendSystemdNotification()

	if err != nil {
		s.Log.Error().Err(err).Msg("Systemd notification error")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Log.Debug().Msg("Stopped server")
	s.StopServer()
	os.Exit(0)
}

func sendSystemdNotification() error {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket != "" {
		state := "READY=1"
		socketAddr := &net.UnixAddr{
			Name: notifySocket,
			Net:  "unixgram",
		}
		conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
		if err != nil {
			return err
		}
		defer conn.Close()
		_, err = conn.Write([]byte(state))
		return err
	}
	return nil
}
