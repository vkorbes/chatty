package ctrl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// "github.com/ellenkorbes/chatty/db"
	db "github.com/ellenkorbes/chatty/nodb"
	"github.com/ellenkorbes/chatty/types"
)

// TestListAllUsers tests the functioning of the ListAllUsers controller method.
func TestListAllUsers(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.ListAllUsers))
	defer ts.Close()
	// And a fake GET request.
	response, err := http.Get(ts.URL)
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting.
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeUsers)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}
}

// TestListAllMessages tests the functioning of the ListAllMessages controller method.
func TestListAllMessages(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.ListAllMessages))
	defer ts.Close()
	// And a fake GET request.
	response, err := http.Get(ts.URL)
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting.
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeMessages2)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}
}

// TestNewUser tests the functioning of the NewUser controller method.
func TestNewUser(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.NewUser))
	defer ts.Close()
	// And a fake POST request.
	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(db.FakeUser))
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting. Here we're unmarshalling since we only care about the Name and Username fields.
	actual, expected := types.User{}, types.User{}
	json.Unmarshal(read, &actual)
	json.Unmarshal(db.FakeUser, &expected)
	if actual.Name != expected.Name || actual.Username != expected.Username {
		t.Error(fmt.Sprintf("Actual: %s - %s\tExpected: %s - %s", actual.Name, actual.Username, expected.Name, expected.Username))
	}
}

// TestGetUserByID tests the functioning of the GetUserByID controller method.
func TestGetUserByID(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.GetUserByID))
	defer ts.Close()
	// And a fake GET request.
	response, err := http.Get(ts.URL + "/5a8d75057d9b53706595116a")
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting.
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeUser)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}
}

// TestGetMessage tests the functioning of the GetMessage controller method.
func TestGetMessage(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.GetMessage))
	defer ts.Close()
	// And a fake GET request.
	response, err := http.Get(ts.URL + "/5a93000c7d9b532f98e8bba2")
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting.
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeMessage)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}
}

// TestMessageRouter tests the functioning of the MessageRouter controller method. This calls either NewMessage or GetMessages depending on whether we receive a POST or a GET request, so we need to test both cases.
func TestMessageRouter(t *testing.T) {
	d := db.NewSession("")
	defer d.Session.Close()
	ctrl := NewController(d)
	// Creating a fake HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(ctrl.MessageRouter))
	defer ts.Close()
	// And a fake GET request.
	response, err := http.Get(ts.URL + "?to=orange")
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	read, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting.
	actual := strings.TrimSpace(string(read))
	expected := string(db.FakeMessages1)
	if actual != expected {
		t.Error(fmt.Sprintf("Actual:\n%sExpected:\n%s", actual, expected))
	}
	// Now for the POST request.
	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(db.FakeMessage))
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	responsePost, err := client.Do(request)
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	readPost, err := ioutil.ReadAll(responsePost.Body)
	responsePost.Body.Close()
	if err != nil {
		t.Error(fmt.Sprintln("Unknown error:", err))
	}
	// And checking whether the response matches what we're expecting. Here we're unmarshalling since we only care about the To, From, and Body fields.
	actualPost, expectedPost := types.Message{}, types.Message{}
	json.Unmarshal(readPost, &actualPost)
	json.Unmarshal(db.FakeMessage, &expectedPost)
	if actualPost.To != expectedPost.To || actualPost.From != expectedPost.From || actualPost.Body != expectedPost.Body {
		t.Error(fmt.Sprintf("Actual: %s - %s - %s\tExpected: %s - %s - %s", actualPost.To, actualPost.From, actualPost.Body, expectedPost.To, expectedPost.From, expectedPost.Body))
	}
}
