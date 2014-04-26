package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const ccApiBase = "http://www.classicube.net/api/"

type resp struct {
	Username      string   `json:"username"`
	Authenticated bool     `json:"authenticated"`
	Token         string   `json:"token"`
	ErrorCount    int      `json:"errorcount"`
	Errors        []string `json:"errors"`
}

type classiCube struct {
	username, password string
	client             *http.Client
}

func (c *classiCube) Username() string {
	return c.username
}

func (c *classiCube) apiUrl(endpoint string) string {
	return ccApiBase + endpoint + "/"
}

func (c *classiCube) Login() (bool, error) {
	// GET to fetch token
	fetch, err := c.client.Get(c.apiUrl("login"))
	if err != nil {
		return false, err
	}
	defer fetch.Body.Close()

	var fetchResp resp
	fetchBody, err := ioutil.ReadAll(fetch.Body)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(fetchBody, &fetchResp); err != nil {
		return false, err
	}

	// POST to /actually/ login.
	auth, err := c.client.PostForm(c.apiUrl("login"),
		url.Values{
			"username": {c.username},
			"password": {c.password},
			"token":    {fetchResp.Token},
		})
	if err != nil {
		return false, err
	}
	defer auth.Body.Close()

	var authResp resp
	authBody, err := ioutil.ReadAll(auth.Body)
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(authBody, &authResp); err != nil {
		return false, err
	}

	if authResp.Authenticated {
		c.username = authResp.Username
	}
	return authResp.Authenticated, nil
}

func (c *classiCube) ServerList() ([]Server, error) {
	serverList, err := c.client.Get(c.apiUrl("serverlist"))
	if err != nil {
		return nil, err
	}
	defer serverList.Body.Close()

	servers := []Server{}
	serverBody, err := ioutil.ReadAll(serverList.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(serverBody, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func NewClassiCube(username, password string) Auth {
	cj, _ := cookiejar.New(nil)
	return &classiCube{
		username: username,
		password: password,
		client: &http.Client{
			Jar: cj,
		},
	}
}
