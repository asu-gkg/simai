package collective

import (
	"simai/common"
	"simai/logical_topology"
	"sync"
)

type Ring struct {
	commType        common.CommType
	id              int
	layerNum        int
	ringTopology    *logical_topology.RingTopology
	dataSize        uint64
	direction       logical_topology.RingDirectionType
	injectionPolicy common.InjectionPolicy
	boostMode       bool

	dimension logical_topology.RingDimensionType

	mu sync.RWMutex
}

func NewRing(commType common.CommType,
	id, layerNum int,
	ringTopology *logical_topology.RingTopology,
	dataSize uint64,
	direction logical_topology.RingDirectionType,
	injectionPolicy common.InjectionPolicy,
	boostMode bool) *Ring {
	return &Ring{
		commType:        commType,
		id:              id,
		layerNum:        layerNum,
		ringTopology:    ringTopology,
		dataSize:        dataSize,
		direction:       direction,
		injectionPolicy: injectionPolicy,
		boostMode:       boostMode,
	}
}
