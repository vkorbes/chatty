package types

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User contains the user fields as per specification.
type User struct {
	ID        bson.ObjectId `json:"id"        bson:"_id,omitempty"` // The unique indentifier of the object. Read only.
	Budget    int           `json:"budget"    bson:"budget"`        // The remaining budget to send messages. Read only.
	Name      string        `json:"name"      bson:"name"`          // The human readable name of the user.
	Username  string        `json:"username"  bson:"username"`      // The unique name of the user. '^[a-z][a-z_\.\-0-9]*$'.
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`     // The UTC date and time user has been created. Read only.
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`     // The UTC date and time user has been updated. Read only.
}

// MarshalJSON is a hack to hijack JSON encoding for this type and format the createdAt and updatedAt fields as per specification.
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	utc, _ := time.LoadLocation("UTC")
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt" bson:"createdAt"`
		UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
	}{
		Alias:     (*Alias)(u),
		CreatedAt: u.CreatedAt.In(utc).Format("2006-01-02T15:04:05.999Z0700"),
		UpdatedAt: u.UpdatedAt.In(utc).Format("2006-01-02T15:04:05.999Z0700"),
	})
}
