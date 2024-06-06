package sysmanager

type SysActorMessage int

//go:generate stringer -type=SysManagerMessage
const (
	Nothing SysActorMessage = iota
	Stop
	SpeedUp
	SpeedDown
)
