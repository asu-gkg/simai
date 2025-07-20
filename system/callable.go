package system

import "simai/common"

type Callable interface {
	Call(eventType common.EventType, callData CallData)
}
