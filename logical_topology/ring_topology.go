package logical_topology

import (
	"fmt"
	"simai/common"
)

type RingDirectionType int

const (
	RingDirectionClockwiseType RingDirectionType = iota
	RingDirectionAnticlockwiseType
)

func (d RingDirectionType) String() string {
	switch d {
	case RingDirectionClockwiseType:
		return "Clockwise"
	case RingDirectionAnticlockwiseType:
		return "Anticlockwise"
	default:
		return "Unknown"
	}
}

// RingDimensionType Dimension环形拓扑的维度
type RingDimensionType int

const (
	RingDimensionLocalType RingDimensionType = iota
	RingDimensionVerticalType
	RingDimensionHorizontalType
	RingDimensionNAType
)

func (d RingDimensionType) String() string {
	switch d {
	case RingDimensionLocalType:
		return "Local"
	case RingDimensionVerticalType:
		return "Vertical"
	case RingDimensionHorizontalType:
		return "Horizontal"
	case RingDimensionNAType:
		return "NA"
	default:
		return "Unknown"
	}
}

type RingTopology struct {
	complexity       Complexity
	name             string
	id               int
	nextNodeID       int
	previousNodeID   int
	offset           int
	totalNodesInRing int
	indexInRing      int
	dimension        RingDimensionType
	idToIndex        map[int]int
}

func NewRingTopology(dimension RingDimensionType, id, totalNodesInRing, indexInRing, offset int) *RingTopology {
	ring := &RingTopology{
		name:             fmt.Sprintf("Ring_%s_%d", dimension.String(), id),
		id:               id,
		totalNodesInRing: totalNodesInRing,
		indexInRing:      indexInRing,
		offset:           offset,
		dimension:        dimension,
		idToIndex:        make(map[int]int),
	}

	ring.findNeighbors()
	return ring
}

func (r *RingTopology) GetNumOfNodesInDimension(dimension int) int {
	return r.GetNodesInRing()
}

// GetNodesInRing 获取环形中的节点数量
func (r *RingTopology) GetNodesInRing() int {
	return r.totalNodesInRing
}

func (r *RingTopology) findNeighbors() {
	r.nextNodeID = r.getReceiverNode(r.id, RingDirectionClockwiseType)
	r.previousNodeID = r.getReceiverNode(r.id, RingDirectionAnticlockwiseType)

	// 构建ID到索引的映射
	for i := 0; i < r.totalNodesInRing; i++ {
		nodeID := r.id + i*r.offset
		r.idToIndex[nodeID] = i
	}
}

func (r *RingTopology) getReceiverNode(nodeID int, direction RingDirectionType) int {
	if direction == RingDirectionClockwiseType {
		return (nodeID + r.offset) % (r.totalNodesInRing * r.offset)
	} else {
		return (nodeID - r.offset + r.totalNodesInRing*r.offset) % (r.totalNodesInRing * r.offset)
	}
}

func (r *RingTopology) getSenderNode(nodeID int, direction RingDirectionType) int {
	if direction == RingDirectionClockwiseType {
		return (nodeID - r.offset + r.totalNodesInRing*r.offset) % (r.totalNodesInRing * r.offset)
	} else {
		return (nodeID + r.offset) % (r.totalNodesInRing * r.offset)
	}
}

func (r *RingTopology) GetDimension() RingDimensionType {
	return r.dimension
}

func (r *RingTopology) GetOffset() int {
	return r.offset
}

func (r *RingTopology) GetID() int {
	return r.id
}

func (r *RingTopology) GetName() string {
	return r.name
}

func (r *RingTopology) GetNumOfDimensions() int {
	return 1
}

func (r *RingTopology) GetLogicalTopologyAtDimension(dimension int, comType common.CommType) LogicalTopology {
	return r
}

func (r *RingTopology) GetComplexity() Complexity {
	return r.complexity
}

func (r *RingTopology) GetLogicalTopologyType() LogicalTopologyType {
	return LogicalTopologyRing
}
