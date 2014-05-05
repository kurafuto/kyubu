package client

import (
	"errors"
	"github.com/sysr-q/kyubu/auth"
	"regexp"
	"strconv"
)

var directRegexp = regexp.MustCompile(`^mc://(?P<host>[a-z0-9.-]+):(?P<port>\d+)/(?P<user>[a-zA-Z0-9_]{1,16})/(?P<hash>[a-fA-F0-9]{32})$`)

func Direct(url string) (s Settings, err error) {
	direct := directRegexp.FindStringSubmatch(url)
	if direct == nil {
		err = errors.New("kyubu: Failed to parse direct connect URL")
		return
	}
	port, err := strconv.Atoi(direct[2])
	if err != nil {
		return
	}
	s = Settings{
		Server:  auth.Server{Address: direct[1], Port: port, MpPass: direct[4]},
		Auth:    auth.NewDirect(direct[3], ""),
		Trickle: 25,
		Debug:   false,
	}
	return
}
