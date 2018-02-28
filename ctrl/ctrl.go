package ctrl

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/ellenkorbes/chatty/types"
	"gopkg.in/mgo.v2/bson"
)

// DBInterface allows us to inject as our DB package anything that fulfills this interface.
type DBInterface interface {
	Add(interface{}) error
	Get(bson.ObjectId, interface{}) error
	GetAll(interface{}) error
	GetUser(string) (types.User, error)
	DecreaseBudget(types.User) error
	GetMessagesByUser(string) (types.Messages, error)
	IsUnique(types.User) (bool, error)
}

// Controller is... pretty simple, just look at it.
type Controller struct {
	DB DBInterface
}

// NewController returns a new Controller.
func NewController(db DBInterface) *Controller {
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

// ListAll lists all items of type *[]types.User or *[]types.Message.
func (c *Controller) ListAll(response http.ResponseWriter, request *http.Request, items interface{}) {
	err := c.DB.GetAll(items)
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
	// Only lowercase alphanumericals, dots, and hyphens.
	r, _ := regexp.Compile(`^[a-z][a-z_\.\-0-9]*$`)
	if !r.MatchString(newUser.Username) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadUsername"])
		return
	}
	if newUser.Name == "" {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BlankUsername"])
		return
	}
	// Creating the new object.
	newUser.ID = bson.NewObjectId()
	newUser.Budget = 10
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	unique, err := c.DB.IsUnique(newUser)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewUser:"+ErrorMessage["db.IsUnique"])
		return
	}
	if !unique {
		Error(response, request, http.StatusConflict, ErrorMessage["TakenUsername"])
		return
	}
	// And off it goes.
	err = c.DB.Add(&newUser)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewUser:"+ErrorMessage["db.Add"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&newUser)
}

// GetUserByUsername returns a full User object based on the username.
func (c *Controller) GetUserByUsername(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		Error(response, request, http.StatusMethodNotAllowed, ErrorMessage["PleaseGET"])
		return
	}
	// Gets the bit of the URL after the last "/"
	user := path.Base(request.URL.Path)
	query, err := c.DB.GetUser(user)
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
	// Gets the bit of the URL after the last "/"
	id := path.Base(request.URL.Path)
	if !bson.IsObjectIdHex(id) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadObjectID"])
		return
	}
	query := types.User{}
	// Get user.
	err := c.DB.Get(bson.ObjectIdHex(id), &query)
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
	// From, To, and Body fields can't be empty. Body can't be larger than 280 characters. We're lumping all of these checks together *before* making a DB call. Because we're cheap.
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
	// Hey database, is the sender real or just an imaginary friend? Habout the recipient?
	sender, err := c.DB.GetUser(newMessage.From)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["SenderNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, ErrorMessage["UnexpectedSender"])
			return
		}
	} else if sender.Budget < 1 {
		// No cheapskates here!
		Error(response, request, http.StatusForbidden, ErrorMessage["BudgetExceeded"])
		return
	}
	_, err = c.DB.GetUser(newMessage.To)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["RecipientNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, ErrorMessage["UnexpectedRecipient"])
			return
		}
	}
	// Filling in the rest of the field.
	newMessage.ID = bson.NewObjectId()
	newMessage.SentAt = time.Now()
	// And boom! New message!
	c.DB.Add(&newMessage)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, "c.NewMessage:"+ErrorMessage["db.Add"])
		return
	}
	// We charge the sender only after checking that there were no errors in adding the message.
	err = c.DB.DecreaseBudget(sender)
	if err != nil {
		if err.Error() == "budget discrepancy" {
			log.Println("Budget discrepancy.", err)
		} else {
			log.Println("Unknown error in dbDecreaseBudget call.", err)
		}
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&newMessage)
}

// GetMessages gets all messages addressed to a specific user.
func (c *Controller) GetMessages(response http.ResponseWriter, request *http.Request) {
	// Hey, look, a param!
	user := request.URL.Query().Get("to")
	_, err := c.DB.GetUser(user)
	if err != nil {
		if err.Error() == "not found" {
			Error(response, request, http.StatusNotFound, ErrorMessage["UserNotFound"])
			return
		} else {
			Error(response, request, http.StatusInternalServerError, "c.GetMessages:"+ErrorMessage["UnexpectedRecipient"])
			return
		}
	}
	messages, err := c.DB.GetMessagesByUser(user)
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
	// Do I need to keep writing these comments? I'll just assume you got the hang of it by now.
	id := path.Base(request.URL.Path)
	if !bson.IsObjectIdHex(id) {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadObjectID"])
		return
	}
	query := types.Message{}
	err := c.DB.Get(bson.ObjectIdHex(id), &query)
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
