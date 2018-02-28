package db

import (
	"encoding/json"

	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DBObject carries the MongoDB session and serves to inject all the code below.
type DBObject struct {
	Session *mgo.Session
}

// NewSession opens a connection to the database.
func NewSession(arg string) DBObject {
	return DBObject{&mgo.Session{}}
}

// FakeUser is a mock user, to be used for testing.
var FakeUser = []byte(`{"id":"5a8d75057d9b53706595116a","budget":7,"name":"Orange","username":"orange","createdAt":"2018-02-21T13:32:53.509Z","updatedAt":"2018-02-25T18:27:25.239Z"}`)

// FakeUsers is a mock list of users, to be used for testing.
var FakeUsers = []byte(`[{"id":"5a8d75057d9b53706595116a","budget":7,"name":"Orange","username":"orange","createdAt":"2018-02-21T13:32:53.509Z","updatedAt":"2018-02-25T18:27:25.239Z"},{"id":"5a8d750d7d9b53706595116b","budget":7,"name":"Banana","username":"banana","createdAt":"2018-02-21T13:33:01.239Z","updatedAt":"2018-02-25T16:50:55.969Z"}]`)

// FakeMessage is a mock message, to be used for testing.
var FakeMessage = []byte(`{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}`)

// FakeMessages1 is a mock list of messages, to be used for testing.
var FakeMessages1 = []byte(`{"messages":[{"id":"5a8d766c7d9b537448d19b2f","from":"banana","to":"orange","body":"Message.","sentAt":"2018-02-21T13:38:52.358Z"},{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}]}`)

// FakeMessages2 is another mock list of messages, to be used for testing.
var FakeMessages2 = []byte(`[{"id":"5a8d766c7d9b537448d19b2f","from":"banana","to":"orange","body":"Message.","sentAt":"2018-02-21T13:38:52.358Z"},{"id":"5a93000c7d9b532f98e8bba2","from":"orange","to":"banana","body":"This is a test message.","sentAt":"2018-02-25T18:27:24.885Z"}]`)

// Add returns nil to simulate a successful DB addition.
func (db DBObject) Add(entry interface{}) error {
	return nil
}

// Get gets a fake user or message, depending on the type of the saveTo argument.
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

// GetAll gets list of "all" fake users or messages, depending on the type of the saveTo argument.
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

// GetUser returns a fake user object.
func (db DBObject) GetUser(user string) (types.User, error) {
	x := types.User{}
	json.Unmarshal(FakeUser, &x)
	return x, nil
}

// DecreaseBudget returns nil to simulate a successful DecreaseBudget operation.
func (db DBObject) DecreaseBudget(sender types.User) error {
	return nil
}

// GetMessagesByUser returns a fake list of messages.
func (db DBObject) GetMessagesByUser(user string) (types.Messages, error) {
	x := types.Messages{}
	json.Unmarshal(FakeMessages1, &x)
	return x, nil
}

// IsUnique returns fake value indicating there are no duplicates in the database.
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
