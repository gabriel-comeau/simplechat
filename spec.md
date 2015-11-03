# Chat Protocol Specification

All messages sent to the server should be a single line of text ending with a \n.
Special commands are started with a / like with IRC, in the format:

/<command> <args>

## The commands

**/nick <nickname>** Changes the user's nickname.  Returns an error if that nickname is already in use.

**/whisper <target> <message>** Sends a message only to another user

**/logout** Does the same thing as EOF - logs user out of server

## Flow

A user connects to the chat server.  The server assigns them a unique ID number, which will be used as their nick name
until they specify one with /nick.

They can send messages to all connected users by just sending the strings and ending them with a \n

They can send a private message to only one user with /whisper