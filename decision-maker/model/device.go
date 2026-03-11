package model

import "net"

type Device struct {
	Id   string
	Ip   net.IP
	Name *string
	Type *string
}
