package analytical_network

import (
	"fmt"
	"simai/common"
	"simai/memory"
	"simai/param_parse"
	"simai/system"
	"sync"
)

var (
	callList *CallTaskHeap
	tick     common.Tick
	mu       sync.Mutex
)

func init() {
	callList = &CallTaskHeap{}
}

func Simulate(param *param_parse.UserParam) {
	// todo
	analyticalNetwork := NewAnalyticalNetwork(0)
	sys := system.NewSys(analyticalNetwork,
		memory.MemAPI{},
		0, 0, 1,
		nil, nil,
		"",
		param.WorkloadsDir,
		param.CommScale, 1, 1,
		1, 0,
		"", "Analytical_test", true, false,
		param.NetWorkParam.GPUType, param.GPUs, param.NetWorkParam.NvSwitchs,
		int(param.NetWorkParam.SwitchNum))

	sys.NVSwitchID = param.NetWorkParam.NvSwitchs[0]
	fmt.Printf("sys.NVSwitchID: %+v\n", sys.NVSwitchID)
	sys.NumGPUs = param.TotalGPUs
	fmt.Printf("sys.NumGPUs: %+v\n", sys.NumGPUs)
	fmt.Printf("sys.workloads: %+v\n", sys.Workloads)
	for _, workload := range sys.Workloads {
		go workload.Fire()
	}

	fmt.Printf("begin run Analytical\n")
	Run()
	Stop()
}

func Stop() {

}
