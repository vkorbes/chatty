package ctrl

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/ellenkorbes/chatty/db"
	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Controller is... pretty simple, just look at it.
type Controller struct {
	DB *mgo.Session
}

// NewController returns a new Controller.
func NewController(db *mgo.Session) *Controller {
	return &Controller{
		DB: db,
	}
}

// ListAllUsers lists all registered users.
func (c *Controller) ListAllUsers(response http.ResponseWriter, request *http.Request) {
	c.ListAll(response, request, &[]types.User{})
}

// ListAllMessages lists all messages.
func (c *Controller) ListAllMessages(response http.ResponseWriter, request *http.Request) {
	c.ListAll(response, request, &[]types.Message{})
}

// ListAll lists all items of the interface{} type. Valid types are *[]types.User and *[]types.Message.
func (c *Controller) ListAll(response http.ResponseWriter, request *http.Request, items interface{}) {
	err := db.GetAll(c.DB, items)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.ListAll: "+ErrorMessage["db.GetAll"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(items)
}

// NewUser creates a new user and returns the resulting object.
func (c *Controller) NewUser(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		Error(response, request, http.StatusMethodNotAllowed, ErrorMessage["PleasePOST"])
		return
	}
	decoder := json.NewDecoder(request.Body)
	var newUser types.User
	err := decoder.Decode(&newUser)
	if err != nil {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadJSON"])
		return
	}
	r, _ := regexp.Compile(`^[a-z][a-z_\.\-0-9]*$`)
	if !r.MatchString(newUser.Username) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadUsername"])
		return
	}
	if newUser.Name == "" {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BlankUsername"])
		return
	}
	newUser.ID = bson.NewObjectId()
	newUser.Budget = 10
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	unique, err := db.IsUnique(c.DB, newUser)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewUser:"+ErrorMessage["db.IsUnique"])
		return
	}
	if !unique {
		Error(response, request, http.StatusConflict, ErrorMessage["TakenUsername"])
		return
	}
	err = db.Add(c.DB, &newUser)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewUser:"+ErrorMessage["db.Add"])
		return
	}
	check, err := db.GetUser(c.DB, newUser.Username)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewUser:"+ErrorMessage["db.GetUser"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&check)
}

// GetUserByUsername returns a full User object based on the username.
func (c *Controller) GetUserByUsername(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		Error(response, request, http.StatusMethodNotAllowed, ErrorMessage["PleaseGET"])
		return
	}
	user := path.Base(request.URL.Path)
	query, err := db.GetUser(c.DB, user)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["UserNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, "c.GetUserByUsername:"+ErrorMessage["db.GetUser"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// GetUserByID returns a full User object based on the ID.
func (c *Controller) GetUserByID(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		Error(response, request, http.StatusMethodNotAllowed, ErrorMessage["PleaseGET"])
		return
	}
	id := path.Base(request.URL.Path)
	if !bson.IsObjectIdHex(id) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadObjectID"])
		return
	}
	query := types.User{}
	err := db.Get(c.DB, bson.ObjectIdHex(id), &query)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["UserNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, "c.GetUserByID:"+ErrorMessage["db.Get"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// NewMessage creates a new message and returns the resulting object.
func (c *Controller) NewMessage(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var newMessage types.Message
	err := decoder.Decode(&newMessage)
	if err != nil {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadJSON"])
		return
	}
	if newMessage.To == "" || newMessage.From == "" || newMessage.Body == "" || len(newMessage.Body) > 280 {
		errors := ""
		if newMessage.To == "" {
			errors += ErrorMessage["EmptyTo"]
		}
		if newMessage.From == "" {
			errors += ErrorMessage["EmptyFrom"]
		}
		switch {
		case newMessage.Body == "":
			errors += ErrorMessage["EmptyBody"]
		case len(newMessage.Body) > 280:
			errors += ErrorMessage["LengthExceeded"]
		}
		Error(response, request, http.StatusBadRequest, strings.TrimSpace(errors))
		return
	}
	sender, err := db.GetUser(c.DB, newMessage.From)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["SenderNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, ErrorMessage["UnexpectedSender"])
			return
		}
	} else if sender.Budget < 1 {
		Error(response, request, http.StatusForbidden, ErrorMessage["BudgetExceeded"])
		return
	}
	_, err = db.GetUser(c.DB, newMessage.To)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["RecipientNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, ErrorMessage["UnexpectedRecipient"])
			return
		}
	}
	newMessage.ID = bson.NewObjectId()
	newMessage.SentAt = time.Now()
	db.Add(c.DB, &newMessage)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewMessage:"+ErrorMessage["db.Add"])
		return
	}
	check := types.Message{}
	err = db.Get(c.DB, newMessage.ID, &check)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewMessage:"+ErrorMessage["db.Get"])
		return
	}
	err = db.DecreaseBudget(c.DB, sender)
	if err != nil {
		if err.Error() == "budget discrepancy" {
			log.Println("Budget discrepancy.", err)
		} else {
			log.Println("Unknown error in dbDecreaseBudget call.", err)
		}
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&check)
}

// GetMessages gets all messages addressed to a specific user.
func (c *Controller) GetMessages(response http.ResponseWriter, request *http.Request) {
	user := request.URL.Query().Get("to")
	_, err := db.GetUser(c.DB, user)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["UserNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, "c.GetMessages:"+ErrorMessage["UnexpectedRecipient"])
			return
		}
	}
	messages, err := db.GetMessagesByUser(c.DB, user)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.GetMessages:"+ErrorMessage["db.GetMessagesByUser"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&messages)
}

// GetMessage returns a full Message object based on the ID.
func (c *Controller) GetMessage(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		Error(response, request, http.StatusMethodNotAllowed, ErrorMessage["PleaseGET"])
		return
	}
	id := path.Base(request.URL.Path)
	if !bson.IsObjectIdHex(id) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadObjectID"])
		return
	}
	query := types.Message{}
	err := db.Get(c.DB, bson.ObjectIdHex(id), &query)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["MessageNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, "c.GetMessage:"+ErrorMessage["db.Get"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// MessageRouter routes requests to /messages to either NewMessage or GetMessages based on the request method.
func (c *Controller) MessageRouter(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		c.NewMessage(response, request)
	}
	if request.Method == "GET" {
		c.GetMessages(response, request)
	}
}
