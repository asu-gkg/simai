package common

// SchedulingPolicy 表示调度策略
type SchedulingPolicy int

const (
	SchedulingPolicyLIFO SchedulingPolicy = iota
	SchedulingPolicyFIFO
	SchedulingPolicyHighest
	SchedulingPolicyNone
)

// IntraDimensionScheduling 表示维度内调度
type IntraDimensionScheduling int

const (
	IntraDimensionSchedulingFIFO IntraDimensionScheduling = iota
	IntraDimensionSchedulingRG
	IntraDimensionSchedulingSmallestFirst
	IntraDimensionSchedulingLessRemainingPhaseFirst
)

// InterDimensionScheduling 表示维度间调度
type InterDimensionScheduling int

const (
	InterDimensionSchedulingAscending InterDimensionScheduling = iota
	InterDimensionSchedulingOnlineGreedy
	InterDimensionSchedulingRoundRobin
	InterDimensionSchedulingOfflineGreedy
	InterDimensionSchedulingOfflineGreedyFlex
)
