package db

import (
	"errors"
	"time"

	"github.com/ellenkorbes/chatty/secrets"
	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// dbInit
func Init() *mgo.Session {
	session, err := mgo.Dial(secrets.Mongo())
	if err != nil {
		panic(err)
	}
	return session
}

// Add - interface must be pointer!
func Add(db *mgo.Session, entry interface{}) error {
	return db.DB("chatty").C(collectionByType(entry)).Insert(entry)
}

// Get - argument saveTo must be a pointer!
func Get(db *mgo.Session, id bson.ObjectId, saveTo interface{}) error {
	return db.DB("chatty").C(collectionByType(saveTo)).FindId(id).One(saveTo)
}

func GetAll(db *mgo.Session, saveTo interface{}) error {
	return db.DB("chatty").C(collectionByType(saveTo)).Find(bson.M{}).All(saveTo)
}

// dbGetUser
func GetUser(db *mgo.Session, user string) (types.User, error) {
	data := types.User{}
	err := db.DB("chatty").C("users").Find(bson.M{"username": user}).One(&data)
	if err != nil {
		return types.User{}, err
	}
	return data, nil
}

// dbDecreaseBudget
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

// dbGetMessagesByUser
func GetMessagesByUser(db *mgo.Session, user string) (types.Messages, error) {
	sm := []types.Message{}
	err := db.DB("chatty").C("messages").Find(bson.M{"to": user}).All(&sm)
	if err != nil {
		return types.Messages{}, err
	}
	return types.Messages{sm}, nil
}

// IsUnique
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

// collectionByType
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
