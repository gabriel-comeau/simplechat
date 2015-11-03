package main

import "net"

type client struct {
	id       int64
	username string
	conn     net.Conn
}

// Creates a new client without setting a nickname since one won't be
// set until the client sends the command
func createNewClient(tcpConn net.Conn) *client {
	return &client{
		id:       getNextId(),
		conn:     tcpConn,
		username: "",
	}
}
