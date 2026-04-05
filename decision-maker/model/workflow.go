package model

type Workflow struct {
	Id        string
	Name      string
	Devices   []Device
	UserId    string
	GestureId *string
	State     string
}

type DeviceTriggerStatus struct {
	DeviceId string
	Ok       bool
	Error    string
}
