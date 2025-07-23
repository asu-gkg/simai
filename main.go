package main

import (
	"fmt"
	"os"
	"simai/network_backend/analytical_network"
	"simai/param_parse"
)

func main() {
	fmt.Println("SimAi starts!")
	param := param_parse.NewUserParam()
	if err := param.Parse(os.Args[1:]); err != nil {
		panic("parse err")
	}

	fmt.Printf("gpu_nums: %+v\n", param.TotalGPUs)
	fmt.Printf("workloads: %+v\n", param.WorkloadsDir)
	// 为每个GPU分配对应的NVSwitch
	node2nvswitch := make(map[int]int)
	for i := 0; i < param.TotalGPUs; i++ {
		nvSwitchId := param.TotalGPUs + i/int(param.NetWorkParam.GpusPerServer)
		node2nvswitch[i] = nvSwitchId
	}
	for i := param.TotalGPUs; i < param.TotalGPUs+int(param.NetWorkParam.NvSwitchNum); i++ {
		node2nvswitch[i] = i
		param.NetWorkParam.NvSwitchs = append(param.NetWorkParam.NvSwitchs, i)
	}

	if param.Mode == param_parse.ModeTypeANALYTICAL {
		analytical_network.Simulate(param)
	}
}
