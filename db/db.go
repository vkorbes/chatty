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

// dbAddUser
func AddUser(db *mgo.Session, user types.User) error {
	c := db.DB("chatty").C("users")
	count, err := c.Find(bson.M{"username": user.Username}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("409")
	}
	return c.Insert(user)
	// TODO: look into mgo.IsDup(err) func
}

// dbGetUser
func GetUser(db *mgo.Session, user string) (types.User, error) {
	data := types.User{}
	err := db.DB("chatty").C("users").Find(bson.M{"username": user}).One(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// dbAddMessage
func AddMessage(db *mgo.Session, message types.Message) error {
	c := db.DB("chatty").C("messages")
	return c.Insert(message)
}

// dbGetMessage
func GetMessage(db *mgo.Session, id bson.ObjectId) (types.Message, error) {
	data := types.Message{}
	err := db.DB("chatty").C("messages").FindId(id).One(&data)
	if err != nil {
		return data, err
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

// dbItemsInCollection
func ItemsInCollection(db *mgo.Session, collection string) (interface{}, error) {
	c := db.DB("chatty").C(collection)
	switch {
	case collection == "users":
		user := types.User{}
		find := c.Find(bson.M{})
		items := find.Iter()
		response := []types.User{}
		for items.Next(&user) {
			response = append(response, user)
		}
		return response, nil
	case collection == "messages":
		message := types.Message{}
		find := c.Find(bson.M{})
		items := find.Iter()
		response := []types.Message{}
		for items.Next(&message) {
			response = append(response, message)
		}
		return response, nil
	}
	return nil, errors.New("Valid collections are: users, messages.")
}

// dbGetUserByID
func GetUserByID(db *mgo.Session, id bson.ObjectId) (types.User, error) {
	data := types.User{}
	err := db.DB("chatty").C("users").FindId(id).One(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// dbGetMessagesByUser
func GetMessagesByUser(db *mgo.Session, user string) (types.Messages, error) {
	c := db.DB("chatty").C("messages")
	sm := []types.Message{}
	err := c.Find(bson.M{"to": user}).All(&sm)
	if err != nil {
		return types.Messages{}, err
	}
	return types.Messages{sm}, nil
}

// dbGetMessagesByID
func GetMessagesByID(db *mgo.Session, id bson.ObjectId) (types.Message, error) {
	m := types.Message{}
	err := db.DB("chatty").C("messages").FindId(id).One(&m)
	if err != nil {
		return types.Message{}, err
	}
	return m, nil
}
