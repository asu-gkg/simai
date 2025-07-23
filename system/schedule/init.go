package schedule

import (
	"simai/common"
	"simai/system_interface"
	"simai/usage"
)

func NewSchedulingUnit(sys system_interface.Sys, queues []int, maxRunningStreams, readyListThreshold, queueThreshold int) *SchedulingUnit {
	su := &SchedulingUnit{
		Sys:                           sys,
		ReadyListThreshold:            readyListThreshold,
		QueueThreshold:                queueThreshold,
		MaxRunningStreams:             maxRunningStreams,
		RunningStreams:                make(map[int]int),
		StreamPointer:                 make(map[int]int),
		LatencyPerDimension:           make([]common.Tick, len(queues)),
		TotalChunksPerDimension:       make([]int64, len(queues)),
		TotalActiveChunksPerDimension: make([]int64, len(queues)),
		QueueIDToDimension:            make(map[int]int),
		UsageTracker:                  make([]*usage.Tracker, len(queues)),
	}

	base := 0
	for dimension, q := range queues {
		for i := 0; i < q; i++ {
			su.RunningStreams[base] = 0
			su.StreamPointer[base] = 0
			su.QueueIDToDimension[base] = dimension
			base++
		}
		su.UsageTracker[dimension] = &usage.Tracker{Usage: 0}
	}

	return su
}
