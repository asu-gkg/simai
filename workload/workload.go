package workload

import (
	"simai/common"
	"simai/system"
	"sync"
)

type Workload struct {
	layers    []*Layer
	size      int
	generator *system.Sys // 一个Sys Simulator可以有多个Workload
	runType   string

	// 状态管理
	counter             int64
	index               int
	currentState        LoopState
	delayLoaded         bool
	separateLog         bool
	checkpointInitiated bool
	collectiveIssued    bool
	initialized         bool

	// 训练参数
	totalPass              int // 模拟的训练轮数
	dlrmLastBottomLayer    int // DLRM 模式下的 MLP bottom 层序号，用于区分 dense/embedding，默认为-1
	passCounter            int // 当前模拟执行的 pass 数，内部计数器
	pendingCollectives     int // 当前未完成的 collective 数量，由系统维护，非手动设置
	modelParallelNPUGroup  int // TP Size
	expertParallelNPUGroup int // EP Size
	pipelineModelParallel  int // PP Size
	ga                     int // GA group size，主要用于 Group AllReduce 之类模拟，也即dp_size

	vpp        int    // 虚拟流水线并行宽度
	ppCommSize uint32 // pipeline 通信大小（用于控制 send/recv 的通信量），默认为0（不设定）

	// 多租户
	allGPUs   []int       // 该租户拥有的GPU
	startTick common.Tick // 任务开始时间
	tenantId  string

	// 策略和状态
	parallelismPolicy ParallelismPolicyType
	waitingForComm    int64

	// 统计和日志
	detailed                     *CSVWriter
	endToEnd                     *CSVWriter
	dimensionUtilization         *CSVWriter
	path                         string
	runName                      string
	statRow                      int
	totalRows                    int
	registeredForFinishedStreams bool

	mu sync.RWMutex
}
