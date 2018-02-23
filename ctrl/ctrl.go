package ctrl

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/ellenkorbes/chatty/db"
	"github.com/ellenkorbes/chatty/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Controller
type Controller struct {
	DB *mgo.Session
}

// NewController
func NewController(db *mgo.Session) *Controller {
	return &Controller{
		DB: db,
	}
}

// ListAllUsers
func (c *Controller) ListAllUsers(response http.ResponseWriter, request *http.Request) {
	c.ListAll(response, request, &[]types.User{})
}

// ListAllMessages
func (c *Controller) ListAllMessages(response http.ResponseWriter, request *http.Request) {
	c.ListAll(response, request, &[]types.Message{})
}

func (c *Controller) ListAll(response http.ResponseWriter, request *http.Request, items interface{}) {
	err := db.GetAll(c.DB, items)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.AddUser"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(items)
}

// NewUser
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
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.AddUser"])
		return
	}
	if !unique {
		Error(response, request, http.StatusConflict, ErrorMessage["TakenUsername"])
		return
	}
	err = db.Add(c.DB, &newUser)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.AddUser"])
		return
	}
	check, err := db.GetUser(c.DB, newUser.Username)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.GetUser"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&check)
}

// GetUserByUsername
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
			Error(response, request, http.StatusInternalServerError, ErrorMessage["db.GetUser"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// GetUserByID
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
			Error(response, request, http.StatusInternalServerError, ErrorMessage["dbGetUserByID"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// NewMessage
func (c *Controller) NewMessage(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var newMessage types.Message
	err := decoder.Decode(&newMessage)
	if err != nil {
		Error(response, request, http.StatusBadRequest, ErrorMessage["BadJSON"])
		return
	}
	// newMessage.Body maxLength: 280
	if newMessage.To == "" || newMessage.From == "" || newMessage.Body == "" {
		errors := ""
		switch {
		case newMessage.To == "":
			errors += "The message sender is empty. "
		case newMessage.From == "":
			errors += "The message recipient is empty. "
		case newMessage.Body == "":
			errors += "The message has no content. "
		}
		Error(response, request, http.StatusBadRequest, errors)
		http.Error(response, errors, http.StatusBadRequest)
		log.Println(errors, err)
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
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.AddMessage"])
		return
	}
	check := types.Message{}
	err = db.Get(c.DB, newMessage.ID, &check)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.GetMessage"])
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

// GetMessages
func (c *Controller) GetMessages(response http.ResponseWriter, request *http.Request) {
	user := request.URL.Query().Get("to")
	messages, err := db.GetMessagesByUser(c.DB, user)
	if err != nil {
		Error(response, request, http.StatusInternalServerError, ErrorMessage["db.GetMessagesByUser"])
		return
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&messages)
}

// GetMessage
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
			Error(response, request, http.StatusInternalServerError, ErrorMessage["db.GetMessagesByID"])
			return
		}
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&query)
}

// MessageRouter
func (c *Controller) MessageRouter(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		c.NewMessage(response, request)
	}
	if request.Method == "GET" {
		c.GetMessages(response, request)
	}
}
