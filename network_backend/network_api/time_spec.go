package network_api

type TimeType int

const (
	TimeTypeSE TimeType = iota
	TimeTypeMS
	TimeTypeUS
	TimeTypeNS
	TimeTypeFS
)

type TimeSpec struct {
	TimeRes TimeType
	TimeVal float64
}
