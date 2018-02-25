package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ellenkorbes/chatty/ctrl"
	"github.com/ellenkorbes/chatty/db"
	"github.com/ellenkorbes/chatty/secrets"
)

func main() {

	argPort := flag.String("p", "8000", "The port the server will listen on")
	argMongo := flag.String("m", secrets.Mongo(), "The MongoDB address URL in the format: mongodb://user:password@yourdatabase.com:12345/dbname")
	flag.Parse()

	d := db.Init(*argMongo)
	defer d.Close()
	ctrl := ctrl.NewController(d)
	mux := http.NewServeMux()

	// Lists all users. Not on spec; added to make development easier.
	mux.HandleFunc("/listusers", ctrl.ListAllUsers)

	// Lists all messages. Not on spec; added to make development easier.
	mux.HandleFunc("/listmsg", ctrl.ListAllMessages)

	// New user.
	mux.HandleFunc("/users", ctrl.NewUser)

	// Get user by id.
	mux.HandleFunc("/users/", ctrl.GetUserByID)

	// POST: New message. GET: Get messages for user.
	mux.HandleFunc("/messages", ctrl.MessageRouter)

	// Get message by id.
	mux.HandleFunc("/message/", ctrl.GetMessage)

	if err := http.ListenAndServe(":"+*argPort, mux); err != nil {
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
