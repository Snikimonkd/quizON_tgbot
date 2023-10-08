package model

type Register struct {
	UserID      int64  `json:"user_id"`
	TgContact   string `json:"tg_contact"`
	TeamID      *int64 `json:"team_id,omitempty"`
	TeamName    string `json:"team_name"`
	CaptainName string `json:"captain_name"`
	Phone       string `json:"phone"`
	GroupName   string `json:"group_name"`
	Amount      string `json:"amount"`
}
