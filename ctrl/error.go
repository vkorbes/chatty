package ctrl

import (
	"encoding/json"
	"net/http"

	"github.com/ellenkorbes/chatty/types"
)

// Error returns a Problem JSON object according to RFC 7807.
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

// ErrorMessage is a central location to store all error messages in the system.
var ErrorMessage map[interface{}]string = map[interface{}]string{
	// These go on Problem.Title:
	201: "Created",
	400: "Bad Request",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	409: "Conflict",
	500: "Internal Server Error",
	// These go on Problem.Detail:
	"PleasePOST":           "Please use a POST request for this endpoint.",
	"PleaseGET":            "Please use a GET request for this endpoint.",
	"db.GetUser":           "Unknown error in db.GetUser call.",
	"db.GetAll":            "Unknown error in db.GetAll call.",
	"db.GetMessagesByUser": "Unknown error in db.GetMessagesByUser call.",
	"db.IsUnique":          "Unknown error in db.IsUnique call.",
	"db.Add":               "Unknown error in db.Add call.",
	"db.Get":               "Unknown error in db.Get call.",
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
	"EmptyTo":              "The message sender is empty. ",
	"EmptyFrom":            "The message recipient is empty. ",
	"EmptyBody":            "The message has no content.",
	"LengthExceeded":       "Message maximum length exceeded: it can contain no more than 280 characters.",
	"BlankMessage":         "",
}
