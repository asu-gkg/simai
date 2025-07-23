package system

import (
	"simai/common"
	"simai/event"
	"simai/logical_topology"
	"simai/memory"
	"simai/network_backend/network_api"
	"simai/stream"
	"simai/system/offline_schedule"
	"simai/system/schedule"
	"simai/workload"
	"sync"
	"time"
)

type Sys struct {
	NetworkAPI network_api.NetworkAPI
	MemoryAPI  memory.MemAPI

	// 基本配置
	ID           int
	NPUOffset    int
	NumGPUs      int
	NGPUsPerNode int
	GPUType      common.GPUType
	AllGPUs      []int
	NVSwitches   []int
	NVSwitchID   int

	// 系统状态
	Initialized       bool
	BoostMode         bool
	RendezvousEnabled bool

	// 时间管理
	StartSimTime time.Time
	EndSimTime   time.Time
	Offset       common.Tick

	// 事件系统
	EventQueue map[common.Tick][]*event.Event
	EventMutex sync.RWMutex

	// 流管理
	ReadyList           []stream.BaseStream
	RunningList         []stream.BaseStream
	ActiveStreams       map[int]stream.BaseStream
	StreamPriorities    map[int]int
	StreamCounter       int
	StreamsInjected     int64
	StreamsFinished     int64
	TotalRunningStreams int

	// 调度器
	SchedulerUnit *schedule.SchedulingUnit
	VLevels       *QueueLevels

	// 集体通信
	LogicalTopologies      map[string]logical_topology.LogicalTopology
	AllReduceImpl          []*common.CollectiveImplementation
	ReduceScatterImpl      []*common.CollectiveImplementation
	AllGatherImpl          []*common.CollectiveImplementation
	AllToAllImpl           []*common.CollectiveImplementation
	CollectiveOptimization common.CollectiveOptimization

	// 工作负载
	Workloads []*workload.Workload
	MemBus    *memory.MemBus

	// 配置参数
	PhysicalDims []int
	QueuesPerDim []int

	// 运行时
	AllQueues         int
	TotalNodes        int
	MaxRunning        int
	ConcurrentStreams int
	ActiveFirstPhase  int

	// 调度策略
	SchedulingPolicy            common.SchedulingPolicy
	IntraDimScheduling          common.IntraDimensionScheduling
	InterDimScheduling          common.InterDimensionScheduling
	RoundRobinInterDimScheduler int

	// 延迟和带宽
	ProcessingLatency   int
	CommunicationDelay  int
	LocalReductionDelay int

	// 缩放因子
	CommScale      float64
	ComputeScale   float64
	InjectionScale float64

	// 全局状态
	AllGenerators     []*Sys
	FinishedWorkloads int

	// 离线调度
	OfflineGreedy           offline_schedule.OfflineGreedy
	LastScheduledCollective common.Tick

	// 待处理发送
	PendingSends        map[common.Pair]chan *SimSendCaller
	IsTherePendingSends map[common.Pair]bool
}
