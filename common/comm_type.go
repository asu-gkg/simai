package common

// ComType 表示通信类型
type CommType int

const (
	CommTypeNone CommType = iota
	CommTypeReduceScatter
	CommTypeAllGather
	CommTypeAllReduce
	CommTypeAllToAll
	CommTypeAllReduceAllToAll
	CommTypeAllReduceNVLS
)

// CollectiveOptimization 表示集合优化类型
type CollectiveOptimization int

const (
	CollectiveOptimizationBaseline CollectiveOptimization = iota
	CollectiveOptimizationLocalBWAware
)

// CollectiveImplementationType 表示集合实现类型
type CollectiveImplementationType int

const (
	CollectiveImplementationTypeRing CollectiveImplementationType = iota
	CollectiveImplementationTypeOneRing
	CollectiveImplementationTypeDirect
	CollectiveImplementationTypeOneDirect
	CollectiveImplementationTypeAllToAll
	CollectiveImplementationTypeDoubleBinaryTreeLocalAllToAll
	CollectiveImplementationTypeLocalRingNodeA2AGlobalDBT
	CollectiveImplementationTypeHierarchicalRing
	CollectiveImplementationTypeDoubleBinaryTree
	CollectiveImplementationTypeHalvingDoubling
	CollectiveImplementationTypeOneHalvingDoubling
	CollectiveImplementationTypeNcclFlowModel
	CollectiveImplementationTypeNcclTreeFlowModel
)

// CollectiveBarrier 表示集合屏障类型
type CollectiveBarrier int

const (
	CollectiveBarrierBlocking CollectiveBarrier = iota
	CollectiveBarrierNonBlocking
)

type BaseCollectiveImplementation struct {
	implType CollectiveImplementationType
}

func NewBaseCollectiveImplementation(implType CollectiveImplementationType) *BaseCollectiveImplementation {
	return &BaseCollectiveImplementation{
		implType: implType,
	}
}

// CollectiveImplementation 表示集合实现
type CollectiveImplementation interface {
	// GetType 返回实现类型
	GetType() CollectiveImplementationType
}
