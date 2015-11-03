package main

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	clients *clientHolder

	whisperCmdRegex *regexp.Regexp = regexp.MustCompile(`^\/w(|hisper) .+ .+$`)
	nickCmdRegex    *regexp.Regexp = regexp.MustCompile(`^\/nick.+$`)
	logoutCmdRegex  *regexp.Regexp = regexp.MustCompile(`^\/logout.+$`)
)

func init() {
	clients = NewClientHolder()
}

const PORT = 1337

func main() {

	// init a TCP server
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil || err != nil {
		panic("couldn't start listening: " + err.Error())
	}

	log.Print("Server Listening")

	for {
		tcpConn, err := server.Accept()

		if err != nil {
			log.Print("Error during server accept: " + err.Error())
		}

		if tcpConn != nil {
			user := createNewClient(tcpConn)
			clients.AddClient(user)

			go handleClient(user)
		}
	}
}

func handleClient(c *client) {
	b := bufio.NewReader(c.conn)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			break
		}

		if line == "" || line == "\n" {
			continue
		}

		dispatch(strings.Trim(line, " \n\t"), c)
	}

	// hit EOF - this client's gone
	log.Printf("Client: %v went EOF and is being removed", c.id)
	clients.RemoveClient(c.id)
}

func dispatch(line string, c *client) {
	// Check if the message is a special command.  If it is validate it and then act.  If not, broadcast.

	if nickCmdRegex.MatchString(line) {
		// Split the command apart to extract the desired nickname
		lineParts := strings.Split(line, " ")
		if len(lineParts) < 2 {
			sendError("Invalid nickname command", c)
		}

		newNick := lineParts[1]
		if clients.IsNickInUse(newNick) {
			sendError("Nick in use or invalid", c)
		}

		c.username = newNick
		clients.AddClient(c)
		sendServerMessage("OK\n", c)
		log.Printf("Client id %v nickname set to %v", c.id, newNick)

	} else if whisperCmdRegex.MatchString(line) {

		lineParts := strings.Split(line, " ")
		targetNick := lineParts[1]

		// Get the rest of the message
		msg := ""
		for i, part := range lineParts[2:] {
			msg += part
			if i+2 < len(lineParts)-1 {
				msg += " "
			}
		}

		// check if we've got a user with that nickname
		target := clients.GetClientByNick(targetNick)
		if target != nil {
			sendWhisper(msg, c, target)
		} else {
			sendError("No user with that nickname online", c)
		}

	} else if logoutCmdRegex.MatchString(line) {
		logout(c)
	} else {
		broadCast(line, c)
	}
}

func sendError(msg string, target *client) {
	formatted := formatErrorMessage(msg)
	target.conn.Write([]byte(formatted))
}

func sendServerMessage(msg string, target *client) {
	target.conn.Write([]byte(msg))
}

func sendWhisper(msg string, sender *client, target *client) {
	formatted := formatWhisperMessage(msg, sender)
	target.conn.Write([]byte(formatted))
}

func broadCast(msg string, sender *client) {
	formatted := formatBroadCastMessage(msg, sender)
	for _, c := range clients.GetClients() {
		c.conn.Write([]byte(formatted))
	}
}

func logout(c *client) {
	sendServerMessage("BYE\n", c)
	clients.RemoveClient(c.id)
	c.conn.Close()
}
