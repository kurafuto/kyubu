// Package auth implements a few Minecraft-related authentication services.
package auth

type Server struct {
	Name   string `json:"name"`
	Hash   string `json:"hash"`
	MpPass string `json:"mppass"`
	Uptime int    `json:"uptime"`

	Address string `json:"ip"`
	Port    int    `json:"port"`

	MaxPlayers int `json:"maxplayers"`
	Players    int `json:"players"`

	Direct string
}

type Auth interface {
	Login() (bool, error)
	ServerList() ([]Server, error)
	Username() string
	//Key() string
}

type AuthFunc func(string, string) Auth
