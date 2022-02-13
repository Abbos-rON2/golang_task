package models

type Response struct {
	Action string      `json:"action"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}

type UserData struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type GuestData struct {
	Ip        string `json:"ip"`
	LastVisit string `json:"last_visit"`
}
