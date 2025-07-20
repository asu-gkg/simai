package common

// BusType 表示总线类型
type BusType int

const (
	BusTypeBoth BusType = iota
	BusTypeShared
	BusTypeMem
)
