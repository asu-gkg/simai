package common

// PacketRouting 表示数据包路由
type PacketRouting int

const (
	PacketRoutingHardware PacketRouting = iota
	PacketRoutingSoftware
)
