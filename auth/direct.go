package auth

import "errors"

// Direct is a "fake" Auth implementation, used internally when parsing
// direct connect URLs.
type Direct struct {
	username string
}

func (d *Direct) Login() (bool, error) {
	return true, nil
}

func (d *Direct) ServerList() ([]Server, error) {
	return nil, errors.New("kyubu: Direct connect auth has no server list.")
}

func (d *Direct) Username() string {
	return d.username
}

func NewDirect(username, password string) Auth {
	return &Direct{username}
}
