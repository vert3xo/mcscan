package types

type ServerPayload struct {
	Ip   string
	Port int
}

type Player struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

type PingResponse struct {
	Host          string        `json:"host"`
	Port          int           `json:"port"`
	Version       string        `json:"version"`
	MaxPlayers    int           `json:"maxPlayers"`
	OnlinePlayers int           `json:"onlinePlayers"`
	PlayersList   []Player		`json:"playersList"`
	Description   interface{}   `json:"description"`
	Favicon       string        `json:"favicon"`
}