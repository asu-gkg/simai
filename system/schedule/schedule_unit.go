package schedule

import (
	"simai/common"
	"simai/system_interface"
	"simai/usage"
)

type SchedulingUnit struct {
	Sys                system_interface.Sys
	ReadyListThreshold int
	QueueThreshold     int
	MaxRunningStreams  int
	RunningStreams     map[int]int
	StreamPointer      map[int]int

	LatencyPerDimension           []common.Tick
	TotalChunksPerDimension       []int64
	TotalActiveChunksPerDimension []int64
	QueueIDToDimension            map[int]int
	UsageTracker                  []*usage.Tracker
}
