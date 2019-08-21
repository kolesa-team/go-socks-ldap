package ldap

import "github.com/kolesa-team/go-socks-ldap/validator"

// entry struct
type Entry struct {
	uid        string
	password   string
	authClient validator.Client
}

func NewEntry(auth validator.Client, uid string, password string) *Entry {
	return &Entry{uid, password, auth}
}

// validate password
func (e *Entry) Validate(password string) bool {
	return e.authClient.Validate(e.password, password)
}

// get entry uid
func (e *Entry) GetUID() string {
	return e.uid
}
