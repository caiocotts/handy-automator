package model

type Workflow struct {
	Id      string
	Name    string
	Devices []Device
	UserId  string
}
