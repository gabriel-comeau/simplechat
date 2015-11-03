package main

import (
	"strconv"
	"sync"
)

type clientHolder struct {
	clients     map[int64]*client
	nickNames   map[string]int64
	*sync.Mutex //embeded
}

func NewClientHolder() *clientHolder {
	return &clientHolder{
		clients:   make(map[int64]*client),
		nickNames: make(map[string]int64),
		Mutex:     new(sync.Mutex),
	}
}

func (ch *clientHolder) GetClientById(id int64) *client {
	ch.Lock()
	defer ch.Unlock()
	client, ok := ch.clients[id]
	if ok {
		return client
	}
	return nil
}

func (ch *clientHolder) GetClientByNick(nick string) *client {
	ch.Lock()
	defer ch.Unlock()
	clientId, ok := ch.nickNames[nick]
	if ok {
		return ch.clients[clientId]
	}

	// didn't find it by name?  check if we can turn that nick into a number and if it's
	// a real id number
	clientId, err := strconv.ParseInt(nick, 10, 64)
	if err != nil {
		return nil
	}

	client, ok := ch.clients[clientId]
	if ok {
		return client
	}

	return nil
}

func (ch *clientHolder) AddClient(c *client) {
	ch.Lock()
	defer ch.Unlock()
	ch.clients[c.id] = c
	if c.username != "" {
		ch.nickNames[c.username] = c.id
	}
}

func (ch *clientHolder) RemoveClient(id int64) {
	ch.Lock()
	defer ch.Unlock()
	client, ok := ch.clients[id]
	if !ok {
		return
	}

	if client.username != "" {
		delete(ch.nickNames, client.username)
	}

	delete(ch.clients, id)
}

func (ch *clientHolder) GetClients() []*client {
	ch.Lock()
	defer ch.Unlock()
	clients := make([]*client, 0)
	for _, c := range ch.clients {
		clients = append(clients, c)
	}
	return clients
}

func (ch *clientHolder) IsNickInUse(nick string) bool {
	ch.Lock()
	defer ch.Unlock()
	_, ok := ch.nickNames[nick]

	if ok {
		return true
	}

	// OK, the actual nick string isn't in use, but what if it's a number and someone's trying
	// to mess up the default-nick-is-the-id thing?
	_, err := strconv.ParseInt(nick, 10, 64)
	if err == nil {
		// valid numeric string - just reject it out of hand
		return true
	}

	return false
}
