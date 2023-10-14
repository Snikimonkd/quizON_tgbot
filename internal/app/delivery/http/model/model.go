package model

type Register struct {
	TgContact   string  `json:"tg_contact"`
	TeamID      *string `json:"team_id,omitempty"`
	TeamName    string  `json:"team_name"`
	CaptainName string  `json:"captain_name"`
	Phone       string  `json:"phone"`
	GroupName   string  `json:"group_name"`
	Amount      string  `json:"amount"`
}

type Registration struct {
	Number      int64   `json:"number"`
	TgContact   string  `json:"tg_contact"`
	TeamID      *string `json:"team_id,omitempty"`
	TeamName    string  `json:"team_name"`
	CaptainName string  `json:"captain_name"`
	Phone       string  `json:"phone"`
	GroupName   string  `json:"group_name"`
	Amount      string  `json:"amount"`
}

type RegisterAvailable struct {
	Ok bool `json:"ok"`
}
