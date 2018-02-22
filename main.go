package main

import (
	"log"
	"net/http"

	"github.com/ellenkorbes/chatty/ctrl"
	"github.com/ellenkorbes/chatty/db"
)

func main() {

	d := db.Init()
	defer d.Close()
	ctrl := ctrl.NewController(d)
	mux := http.NewServeMux()
	mux.HandleFunc("/listusers", ctrl.ListAllUsers)  // List all users.
	mux.HandleFunc("/listmsg", ctrl.ListAllMessages) // List all messages.
	mux.HandleFunc("/users", ctrl.NewUser)           // New user.
	mux.HandleFunc("/users/", ctrl.GetUserByID)      // Get user by id.
	mux.HandleFunc("/messages", ctrl.MessageRouter)  // POST: New message. GET: Get messages for user.
	mux.HandleFunc("/message/", ctrl.GetMessage)     // Get message by id.

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}

}

/*

Operations:

Create a user.
POST /users
{
	"name": "Peter Gibbons",
	"username": "peter.gibbons"
  }
Responses:
201 - The user object representation. application/json.
400 - The user object is bad formatted, missing attributes or has invalid values. application/problem+json.
409 - The username is already taken by another user. application/problem+json.
default - Unexpected Error. application/problem+json.

Get a user by id.
GET /users/{id}
Responses:
200 - The user object representation. application/json.
404 - The user was not found. application/problem+json.
default - Unexpected Error. application/problem+json.

Send message from one user to another.
POST /messages
{
	"from": "string",
	"to": "string",
	"body": "string"
}
Responses:
201 - The message object representation. application/json.

List the messages a user has received.
GET /messages?to={user}
Responses:
200 - The message listing representation. application/json.
400 - The request is missing required attributes. application/problem+json.
404 - The user was not found. application/problem+json.
default - Unexpected Error. application/problem+json.

Get a message by id.
GET /message/{id}
Responses:
200 - The message object representation. application/json.
404 - The message was not found. application/problem+json.
default - Unexpected Error. application/problem+json.

*/
