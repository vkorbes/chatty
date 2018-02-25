package db

import (
	"errors"
	"time"

	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Init opens a connection to the database. Hostname, post, user, and password must be supplied by package secrets.
func Init(arg string) *mgo.Session {
	session, err := mgo.Dial(arg)
	if err != nil {
		panic(err)
	}
	return session
}

// Add adds an entry to the database. The interface{} argument must be a pointer.
func Add(db *mgo.Session, entry interface{}) error {
	return db.DB("chatty").C(collectionByType(entry)).Insert(entry)
}

// Get gets an entry from the database. The interface{} argument must be a pointer.
func Get(db *mgo.Session, id bson.ObjectId, saveTo interface{}) error {
	return db.DB("chatty").C(collectionByType(saveTo)).FindId(id).One(saveTo)
}

// GetAll gets all items in a collection. The interface{} argument must be a pointer.
func GetAll(db *mgo.Session, saveTo interface{}) error {
	return db.DB("chatty").C(collectionByType(saveTo)).Find(bson.M{}).All(saveTo)
}

// GetUser gets the full User object for a username.
func GetUser(db *mgo.Session, user string) (types.User, error) {
	data := types.User{}
	err := db.DB("chatty").C("users").Find(bson.M{"username": user}).One(&data)
	if err != nil {
		return types.User{}, err
	}
	return data, nil
}

// DecreaseBudget decreases a user's budget by 1.
func DecreaseBudget(db *mgo.Session, sender types.User) error {
	userCheck := types.User{}
	budget := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"budget": -1}, "$set": bson.M{"updatedAt": time.Now()}},
		ReturnNew: true,
	}
	_, err := db.DB("chatty").C("users").Find(bson.M{"username": sender.Username}).Apply(budget, &userCheck)
	if err != nil {
		return err
	}
	if sender.Budget-1 != userCheck.Budget {
		return errors.New("budget discrepancy")
	}
	return nil
}

// GetMessagesByUser gets all messages addressed to a specific user.
func GetMessagesByUser(db *mgo.Session, user string) (types.Messages, error) {
	sm := []types.Message{}
	err := db.DB("chatty").C("messages").Find(bson.M{"to": user}).All(&sm)
	if err != nil {
		return types.Messages{}, err
	}
	return types.Messages{sm}, nil
}

// IsUnique checks whether a username is already present in the database.
func IsUnique(db *mgo.Session, user types.User) (bool, error) {
	c := db.DB("chatty").C("users")
	count, err := c.Find(bson.M{"username": user.Username}).Limit(1).Count()
	if err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
	return true, nil
}

// collectionByType returns the fitting collection name based on the type of the object supplied.
func collectionByType(x interface{}) string {
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
