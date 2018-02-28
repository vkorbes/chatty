package ctrl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// "github.com/ellenkorbes/chatty/db"
	db "github.com/ellenkorbes/chatty/nodb"
)

func TestListAllUsers(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	ts := httptest.NewServer(http.HandlerFunc(ctrl.ListAllUsers))
	defer ts.Close()
	response, err := http.Get(ts.URL)
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeUsers)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}

}

// To test:

// // Lists all users. Not on spec; added to make development easier.
// mux.HandleFunc("/listusers", ctrl.ListAllUsers)

// // Lists all messages. Not on spec; added to make development easier.
// mux.HandleFunc("/listmsg", ctrl.ListAllMessages)

// // New user.
// mux.HandleFunc("/users", ctrl.NewUser)

// // Get user by id.
// mux.HandleFunc("/users/", ctrl.GetUserByID)

// // POST: New message. GET: Get messages for user.
// mux.HandleFunc("/messages", ctrl.MessageRouter)

// // Get message by id.
// mux.HandleFunc("/message/", ctrl.GetMessage)
