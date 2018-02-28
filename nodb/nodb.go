package db

import (
	"encoding/json"

	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBObject struct {
	Session *mgo.Session
}

// Init opens a connection to the database. Hostname, post, user, and password must be supplied by package secrets.
func NewSession(arg string) DBObject {
	return DBObject{&mgo.Session{}}
}

var FakeUser = []byte(`{"id":"5a8d75057d9b53706595116a","budget":7,"name":"Orange","username":"orange","createdAt":"2018-02-21T13:32:53.509Z","updatedAt":"2018-02-25T18:27:25.239Z"}`)

var FakeUsers = []byte(`[{"id":"5a8d75057d9b53706595116a","budget":7,"name":"Orange","username":"orange","createdAt":"2018-02-21T13:32:53.509Z","updatedAt":"2018-02-25T18:27:25.239Z"},{"id":"5a8d750d7d9b53706595116b","budget":7,"name":"Banana","username":"banana","createdAt":"2018-02-21T13:33:01.239Z","updatedAt":"2018-02-25T16:50:55.969Z"}]`)

var FakeMessage = []byte(`{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}`)

var FakeMessages1 = []byte(`{"messages":[{"id":"5a8d766c7d9b537448d19b2f","from":"banana","to":"orange","body":"Message.","sentAt":"2018-02-21T13:38:52.358Z"},{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}]}`)

var FakeMessages2 = []byte(`[{"id":"5a8d766c7d9b537448d19b2f","from":"banana","to":"orange","body":"Message.","sentAt":"2018-02-21T13:38:52.358Z"},{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}]`)

// Add adds an entry to the database. The interface{} argument must be a pointer.
func (db DBObject) Add(entry interface{}) error {
	return nil
}

// Get gets an entry from the database. The interface{} argument must be a pointer.
func (db DBObject) Get(id bson.ObjectId, saveTo interface{}) error {
	switch saveTo.(type) {
	case *types.User:
		json.Unmarshal(FakeUser, saveTo)
		return nil
	case *types.Message:
		json.Unmarshal(FakeMessage, saveTo)
		return nil
	}
	return nil
}

// GetAll gets all items in a collection. The interface{} argument must be a pointer.
func (db DBObject) GetAll(saveTo interface{}) error {
	switch saveTo.(type) {
	case *[]types.User:
		json.Unmarshal(FakeUsers, saveTo)
		return nil
	case *[]types.Message:
		json.Unmarshal(FakeMessages2, saveTo)
		return nil
	}
	return nil
}

// GetUser gets the full User object for a username.
func (db DBObject) GetUser(user string) (types.User, error) {
	x := types.User{}
	json.Unmarshal(FakeUser, &x)
	return x, nil
}

// DecreaseBudget decreases a user's budget by 1.
func (db DBObject) DecreaseBudget(sender types.User) error {
	return nil
}

// GetMessagesByUser gets all messages addressed to a specific user.
func (db DBObject) GetMessagesByUser(user string) (types.Messages, error) {
	x := types.Messages{}
	json.Unmarshal(FakeMessages1, &x)
	return x, nil
}

// IsUnique checks whether a username is already present in the database.
func (db DBObject) IsUnique(user types.User) (bool, error) {
	return true, nil
}

// CollectionByType returns the fitting collection name based on the type of the object supplied.
func CollectionByType(x interface{}) string {
	switch x.(type) {
	case *types.User:
		return "users"
	case *[]types.User:
		return "users"
	case *types.Message:
		return "messages"
	case *[]types.Message:
		return "messages"
	}
	return ""
}
