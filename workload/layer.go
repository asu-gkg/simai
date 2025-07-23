package workload

import (
	"simai/common"
	"simai/dataset"
	"simai/mock_nccl"
	"simai/system_interface"
	"sync"
)

type LoopState int

const (
	LoopStateForwardPass LoopState = iota
	LoopStateWeightGradient
	LoopStateInputGradient
	LoopStateWaitForSimFinish
	LoopStateForwardInBackPass
)

func (l LoopState) String() string {
	switch l {
	case LoopStateForwardPass:
		return "Forward_Pass"
	case LoopStateWeightGradient:
		return "Weight_Gradient"
	case LoopStateInputGradient:
		return "Input_Gradient"
	case LoopStateWaitForSimFinish:
		return "Wait_For_Sim_Finish"
	case LoopStateForwardInBackPass:
		return "Forward_In_BackPass"
	default:
		return "Unknown"
	}
}

type Layer struct {
	id        string
	layerNum  int
	generator *system_interface.Sys
	workload  *Workload
	name      string

	// 前向传播相关
	fwdPassComputeTime            int64
	fwdPassCommType               common.CommType         // 由并行模式指定
	fwdPassGroupType              mock_nccl.NcclGroupType // 由并行模式指定
	fwdPassCommSize               uint64
	fwdUpdateTime                 int64
	fwdPassCommInvolvedDimensions []bool

	// 输入梯度相关，常伴随 DP AllReduce
	inputGradComputeTime            int64
	inputGradCommType               common.CommType         // 由并行模式指定
	inputGradGroupType              mock_nccl.NcclGroupType // 由并行模式指定
	inputGradCommSize               uint64
	inputGradUpdateTime             int64
	inputGradCommInvolvedDimensions []bool

	// 权重梯度相关，常伴随 TP/PP 的 Send/Recv
	weightGradComputeTime            int64
	weightGradCommType               common.CommType         // 由并行模式指定
	weightGradGroupType              mock_nccl.NcclGroupType // 由并行模式指定
	weightGradCommSize               uint64
	weightGradUpdateTime             int64
	weightGradCommInvolvedDimensions []bool

	needsFwdInBackwardInitiation bool
	isCheckpoint                 bool                  // 由上层指定
	specificParallelism          ParallelismPolicyType // 由上层指定

	// 统计信息
	lookupTableSize   int
	collectiveCounter int

	// 数据集管理
	fwdPassDatasets    map[int]*dataset.Dataset
	inputGradDatasets  map[int]*dataset.Dataset
	weightGradDatasets map[int]*dataset.Dataset

	// 等待时间记录
	startedWaitingForFwdPass    []common.Tick
	startedWaitingForInputGrad  []common.Tick
	startedWaitingForWeightGrad []common.Tick

	// 总计算和通信时间
	totalForwardPassCompute common.Tick
	totalInputGradCompute   common.Tick
	totalWeightGradCompute  common.Tick
	totalWeightGradComm     common.Tick
	totalInputGradComm      common.Tick
	totalFwdComm            common.Tick

	// 最后完成时间
	lastFwdFinished common.Tick
	lastWgFinished  common.Tick
	lastIgFinished  common.Tick

	// 等待通信时间
	totalWaitingForWgComm  common.Tick
	totalWaitingForIgComm  common.Tick
	totalWaitingForFwdComm common.Tick

	fwdBarrier common.CollectiveBarrier
	wgBarrier  common.CollectiveBarrier
	igBarrier  common.CollectiveBarrier

	allGPUs  []int // 该层使用的GPU
	tenantId string

	mu sync.RWMutex
}

type LayerDataReport struct {
}
