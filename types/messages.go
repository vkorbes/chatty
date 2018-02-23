package types

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Messages is a slice of Message.
type Messages struct {
	Entries []Message `json:"messages" bson:"messages"`
}

// Message contains the message fields as per specification.
type Message struct {
	ID     bson.ObjectId `json:"id"      bson:"_id,omitempty"` // The unique indentifier of the object. Read only.
	From   string        `json:"from"    bson:"from"`          // The sender user id.
	To     string        `json:"to"      bson:"to"`            // The recipient user id.
	Body   string        `json:"body"    bson:"body"`          // The message body content. Length: 1–280.
	SentAt time.Time     `json:"sentAt"  bson:"sentAt"`        // The UTC date and time message was sent. Read only.
}

// MarshalJSON is a hack to hijack JSON encoding for this type and format the sentAt field as per specification.
func (u *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	utc, _ := time.LoadLocation("UTC")
	return json.Marshal(&struct {
		*Alias
		SentAt string `json:"sentAt"  bson:"sentAt"`
	}{
		Alias:  (*Alias)(u),
		SentAt: u.SentAt.In(utc).Format("2006-01-02T15:04:05.999Z0700"),
	})
}
