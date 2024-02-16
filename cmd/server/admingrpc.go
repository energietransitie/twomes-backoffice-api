package main

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/energietransitie/twomes-backoffice-api/handlers"
	"github.com/sirupsen/logrus"
)

// Setup Admin GRPC Handler.
func setupAdminRPCHandler(adminHandler *handlers.AdminHandler) {

	rpc.Register(adminHandler)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp4", "127.0.0.1:8081")
	if err != nil {
		logrus.Fatal(err)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
