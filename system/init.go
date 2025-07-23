package system

import (
	"simai/common"
	"simai/event"
	"simai/logical_topology"
	"simai/memory"
	"simai/network_backend/network_api"
	"simai/stream"
	"simai/system/schedule"
	"time"
)

func NewSys(networkAPI network_api.NetworkAPI,
	memoryAPI memory.MemAPI,
	id int,
	npuOffset int,
	numPasses int,
	physicalDims []int,
	queuesPerDim []int,
	systemName string,
	workloadDir string,
	commScale float64,
	computeScale float64,
	injectionScale float64,
	totalStatRows int,
	statRow int,
	path string,
	runName string,
	separateLog bool,
	rendezvousEnabled bool,
	gpuType common.GPUType,
	allGPUs []int,
	nvSwitches []int,
	nGPUsPerNode int) *Sys {
	sys := &Sys{
		ID:                  id,
		NPUOffset:           npuOffset,
		NumGPUs:             len(allGPUs),
		NGPUsPerNode:        nGPUsPerNode,
		GPUType:             gpuType,
		AllGPUs:             allGPUs,
		NVSwitches:          nvSwitches,
		NetworkAPI:          networkAPI,
		MemoryAPI:           memoryAPI,
		Initialized:         false,
		BoostMode:           false,
		RendezvousEnabled:   rendezvousEnabled,
		StartSimTime:        time.Now(),
		Offset:              0,
		EventQueue:          make(map[common.Tick][]*event.Event),
		ReadyList:           make([]stream.BaseStream, 0),
		RunningList:         make([]stream.BaseStream, 0),
		ActiveStreams:       make(map[int]stream.BaseStream),
		StreamPriorities:    make(map[int]int),
		StreamCounter:       0,
		LogicalTopologies:   make(map[string]logical_topology.LogicalTopology),
		AllReduceImpl:       make([]*common.CollectiveImplementation, 0),
		ReduceScatterImpl:   make([]*common.CollectiveImplementation, 0),
		AllGatherImpl:       make([]*common.CollectiveImplementation, 0),
		AllToAllImpl:        make([]*common.CollectiveImplementation, 0),
		PhysicalDims:        physicalDims,
		QueuesPerDim:        queuesPerDim,
		MaxRunning:          100000000,
		ConcurrentStreams:   1,
		ActiveFirstPhase:    100000000,
		SchedulingPolicy:    common.SchedulingPolicyFIFO,
		IntraDimScheduling:  common.IntraDimensionSchedulingFIFO,
		InterDimScheduling:  common.InterDimensionSchedulingAscending,
		ProcessingLatency:   10,
		CommunicationDelay:  10,
		LocalReductionDelay: 1,
		CommScale:           commScale,
		ComputeScale:        computeScale,
		InjectionScale:      injectionScale,
		AllGenerators:       make([]*Sys, 0),
		FinishedWorkloads:   0,
		PendingSends:        make(map[common.Pair]chan *SimSendCaller),
		IsTherePendingSends: make(map[common.Pair]bool),
	}

	sys.SchedulerUnit = schedule.NewSchedulingUnit(sys, queuesPerDim,
		sys.MaxRunning, sys.ActiveFirstPhase, sys.ConcurrentStreams)
	if workloadDir != "" {
		panic("impl me")
	}
	return sys
}
