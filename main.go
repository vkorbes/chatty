package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ellenkorbes/chatty/ctrl"
	"github.com/ellenkorbes/chatty/db"
	"github.com/ellenkorbes/chatty/secret"
	// db "github.com/ellenkorbes/chatty/nodb"
)

func main() {

	// Fun with flags!
	argPort := flag.String("p", "8000", "The port the server will listen on")
	argMongo := flag.String("m", secret.Secret(), "The MongoDB address URL in the format: mongodb://user:password@yourdatabase.com:12345/dbname")
	flag.Parse()

	// New database session, new controller, new http server.
	d := db.NewSession(*argMongo)
	defer d.Session.Close()
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

	// Off we go!
	if err := http.ListenAndServe(":"+*argPort, mux); err != nil {
		log.Fatal(err)
	}

}
