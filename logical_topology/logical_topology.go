package logical_topology

import "simai/common"

type Complexity int

const (
	ComplexityBasic Complexity = iota
	ComplexityComplex
)

func (c Complexity) String() string {
	switch c {
	case ComplexityBasic:
		return "Basic"
	case ComplexityComplex:
		return "Complex"
	default:
		return "Unknown"
	}
}

type LogicalTopologyType int

const (
	LogicalTopologyRing LogicalTopologyType = iota
	LogicalTopologyBinaryTree
)

func (topo LogicalTopologyType) String() string {
	switch topo {
	case LogicalTopologyRing:
		return "Ring"
	case LogicalTopologyBinaryTree:
		return "BinaryTree"
	default:
		return "Unknown"
	}
}

type LogicalTopology interface {
	// GetNumOfDimensions 返回拓扑的维度数量
	GetNumOfDimensions() int

	// GetNumOfNodesInDimension 返回指定维度中的节点数量
	GetNumOfNodesInDimension(dimension int) int

	// GetLogicalTopologyAtDimension 返回指定维度的逻辑拓扑
	GetLogicalTopologyAtDimension(dimension int, comType common.CommType) LogicalTopology

	// GetComplexity 返回拓扑的复杂度
	GetComplexity() Complexity

	// GetLogicalTopologyType 返回逻辑拓扑类型
	GetLogicalTopologyType() LogicalTopologyType
}
