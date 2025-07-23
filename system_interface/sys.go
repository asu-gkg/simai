package system_interface

import "simai/common"

type Sys interface {
	BoostedTick() common.Tick
}
