package ctrl

import (
	"encoding/json"
	"net/http"

	"github.com/ellenkorbes/chatty/types"
)

// Error
func Error(response http.ResponseWriter, request *http.Request, status int, detail string) {
	response.Header().Set("Content-Type", "application/problem+json")
	response.WriteHeader(status)
	json.NewEncoder(response).Encode(&types.Problem{
		Type:     "",
		Title:    ErrorMessage[status],
		Status:   status,
		Detail:   detail,
		Details:  []string{},
		Instance: request.URL.Path,
	})
}

var ErrorMessage map[interface{}]string = map[interface{}]string{
	// These go on Problem.Title:
	201: "",
	400: "Bad Request",
	403: "",
	404: "",
	405: "Method Not Allowed",
	409: "Conflict",
	500: "Internal Server Error",
	// These go on Problem.Detail:
	"PleasePOST":           "Please use a POST request for this endpoint.",
	"PleaseGET":            "Please use a GET request for this endpoint.",
	"db.AddUser":           "Unknown error in db.AddUser call.",
	"db.GetUser":           "Unknown error in db.GetUser call.",
	"dbGetUserByID":        "Unknown error in dbGetUserByID call.",
	"db.AddMessage":        "Unknown error in db.AddMessage call.",
	"db.GetMessage":        "Unknown error in db.GetMessage call.",
	"db.GetMessagesByUser": "Unknown error in db.GetMessagesByUser call.",
	"db.GetMessagesByID":   "Unknown error in db.GetMessagesByID call.",
	"BadJSON":              "Error parsing JSON object.",
	"BadUsername":          "The username should only contain lowercase alphanumerical characters, dashes, and underscores.",
	"BlankUsername":        "The username value cannot be blank.",
	"TakenUsername":        "This username has already been taken by another user.",
	"UserNotFound":         "Username not found.",
	"MessageNotFound":      "Message not found.",
	"BadObjectID":          "The supplied object ID is invalid.",
	"SenderNotFound":       "Sender username not found.",
	"UnexpectedSender":     "Unknown error verifying sender.",
	"BudgetExceeded":       "The sender username has no budget left.",
	"RecipientNotFound":    "Recipient username not found.",
	"UnexpectedRecipient":  "Unknown error verifying recipient.",
	"BlankMessage":         "",
}

// Error(response, request, http., ErrorMessages[""])
