package model

import "net"

type Device struct {
	Id          string
	Hostname    string
	LastKnownIp net.IP
	Name        *string
	Type        *string
}
