package system

import (
	"github.com/spf13/cast"
	"simai/common"
)

func (sys *Sys) BoostedTick() common.Tick {
	var ts *Sys
	for _, gen := range sys.AllGenerators {
		if gen != nil {
			ts = gen
			break
		}
	}

	if ts == nil {
		return 0
	}

	timeSpec := ts.NetworkAPI.SimGetTime()
	tick := timeSpec.TimeVal / common.CLOCK_PERIOD
	return common.Tick(cast.ToUint64(tick)) + sys.Offset
}
