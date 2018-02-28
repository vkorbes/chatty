package db

import (
	"errors"
	"time"

	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBObject struct {
	Session *mgo.Session
}

// Init opens a connection to the database. Hostname, post, user, and password must be supplied by package secrets.
func NewSession(arg string) DBObject {
	session, err := mgo.Dial(arg)
	if err != nil {
		panic(err)
	}
	return DBObject{session}
}

// Add adds an entry to the database. The interface{} argument must be a pointer.
func (db DBObject) Add(entry interface{}) error {
	return db.Session.DB("chatty").C(CollectionByType(entry)).Insert(entry)
}

// Get gets an entry from the database. The interface{} argument must be a pointer.
func (db DBObject) Get(id bson.ObjectId, saveTo interface{}) error {
	return db.Session.DB("chatty").C(CollectionByType(saveTo)).FindId(id).One(saveTo)
}

// GetAll gets all items in a collection. The interface{} argument must be a pointer.
func (db DBObject) GetAll(saveTo interface{}) error {
	return db.Session.DB("chatty").C(CollectionByType(saveTo)).Find(bson.M{}).All(saveTo)
}

// GetUser gets the full User object for a username.
func (db DBObject) GetUser(user string) (types.User, error) {
	data := types.User{}
	err := db.Session.DB("chatty").C("users").Find(bson.M{"username": user}).One(&data)
	if err != nil {
		return types.User{}, err
	}
	return data, nil
}

// DecreaseBudget decreases a user's budget by 1.
func (db DBObject) DecreaseBudget(sender types.User) error {
	userCheck := types.User{}
	budget := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"budget": -1}, "$set": bson.M{"updatedAt": time.Now()}},
		ReturnNew: true,
	}
	_, err := db.Session.DB("chatty").C("users").Find(bson.M{"username": sender.Username}).Apply(budget, &userCheck)
	if err != nil {
		return err
	}
	if sender.Budget-1 != userCheck.Budget {
		return errors.New("budget discrepancy")
	}
	return nil
}

// GetMessagesByUser gets all messages addressed to a specific user.
func (db DBObject) GetMessagesByUser(user string) (types.Messages, error) {
	sm := []types.Message{}
	err := db.Session.DB("chatty").C("messages").Find(bson.M{"to": user}).All(&sm)
	if err != nil {
		return types.Messages{}, err
	}
	return types.Messages{sm}, nil
}

// IsUnique checks whether a username is already present in the database.
func (db DBObject) IsUnique(user types.User) (bool, error) {
	c := db.Session.DB("chatty").C("users")
	count, err := c.Find(bson.M{"username": user.Username}).Limit(1).Count()
	if err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
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
