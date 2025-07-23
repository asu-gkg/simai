package param_parse

import (
	"flag"
	"path/filepath"
	"simai/common"
	"strings"
)

type UserParam struct {
	GPUs         []int
	TotalGPUs    int
	WorkloadsDir string  // 可支持多个workloads文件
	ResultPrefix string  // 结果文件名前缀
	ResultDir    string  // 结果输出目录
	CommScale    float64 // 通信放大比例，用于“虚拟”放大通信负载，方便仿真不同通信压力场景。例如将通信大小扩大 10 倍，用于测试带宽瓶颈。
	Mode         ModeType
	NetWorkParam NetWorkParamV1
}

func NewUserParam() *UserParam {
	return &UserParam{}
}

func (p *UserParam) Parse(args []string) error {
	var gpuType string
	var mode string
	var cluster_topology string
	flag.StringVar(&p.WorkloadsDir, "workloads", "", "Workloads directory path")
	flag.IntVar(&p.TotalGPUs, "total_gpus", 0, "Total GPUs")
	flag.StringVar(&p.ResultPrefix, "result_prefix", "results", "Result name")
	flag.StringVar(&p.ResultDir, "result_dir", "None", "Result folder")
	flag.UintVar(&p.NetWorkParam.GpusPerServer, "gpus_per_server", 8, "GPUs per server")
	flag.Float64Var(&p.NetWorkParam.NvLinkBw, "nvlink", -1.0, "Nvlink bandwidth")
	flag.Float64Var(&p.NetWorkParam.BwPerNic, "nic_busbw", -1.0, "NIC bandwidth")
	flag.UintVar(&p.NetWorkParam.NicsPerServer, "nic_per_server", 1, "NICs per server")
	flag.StringVar(&p.NetWorkParam.NicType, "nic_type", "cx7", "NIC type")
	flag.StringVar(&gpuType, "gpu_type", "NONE", "GPU type")
	flag.BoolVar(&p.NetWorkParam.Visual, "visual", false, "Enable visualization")
	flag.Float64Var(&p.NetWorkParam.DpOverlapRatio, "dp_overlap", 0.0, "DP overlap ratio")
	flag.Float64Var(&p.NetWorkParam.TpOverlapRatio, "tp_overlap", 0.0, "TP overlap ratio")
	flag.Float64Var(&p.NetWorkParam.EpOverlapRatio, "ep_overlap", 0.0, "EP overlap ratio")
	flag.Float64Var(&p.NetWorkParam.PpOverlapRatio, "pp_overlap", 1.0, "PP overlap ratio")
	flag.StringVar(&cluster_topology, "cluster_topology", "", "Cluster Topology File")
	// 默认情况下设置的switch数量
	flag.UintVar(&p.NetWorkParam.SwitchNum, "switch_num", 0, "Switch number")
	flag.StringVar(&mode, "mode", "ANALYTICAL", "mode type")

	flag.Parse()

	p.NetWorkParam.GPUType = common.ToGPUType(gpuType)

	switch strings.ToUpper(mode) {
	case "ANALYTICAL":
		p.Mode = ModeTypeANALYTICAL
	case "NS3":
		p.Mode = ModeTypeNS3
	case "HTSIM":
		p.Mode = ModeTypeHTSIM
	case "TECCL":
		p.Mode = ModeTypeTECCL
	default:
		p.Mode = ModeTypeNONE
	}

	if len(p.GPUs) > 0 {
		p.NetWorkParam.NvSwitchNum = uint32(uint(p.TotalGPUs) / p.NetWorkParam.GpusPerServer)
		if p.NetWorkParam.SwitchNum == 0 {
			p.NetWorkParam.SwitchNum = 120 + p.NetWorkParam.GpusPerServer
		}
		p.NetWorkParam.NodeNum = p.NetWorkParam.NvSwitchNum + uint32(p.NetWorkParam.SwitchNum) + uint32(p.TotalGPUs)
	}

	if p.ResultDir == "None" && p.WorkloadsDir != "" {
		base := filepath.Base(p.WorkloadsDir)
		p.ResultDir = filepath.Join("results", base+"-"+p.ResultPrefix)
	}

	if cluster_topology != "" {
		topo, err := NewClusterTopoParam(cluster_topology)
		if err != nil {
			panic(err)
		}
		p.NetWorkParam.GPUType = common.ToGPUType(topo.GPUType)
		p.NetWorkParam.SwitchNum = uint(topo.NumSwitches)
		p.NetWorkParam.NvLinkBw = topo.NVLinkBW
		p.NetWorkParam.NvSwitchNum = uint32(topo.NVSwitchNum)
		p.NetWorkParam.GpusPerServer = uint(topo.GPUPerServer)
		p.NetWorkParam.AllGpus = topo.AllGPUs
		p.NetWorkParam.BwPerNic = topo.BWPerNIC

		p.TotalGPUs = 0
		for _, gpus := range p.NetWorkParam.AllGpus {
			for _, gpu := range gpus {
				p.GPUs = append(p.GPUs, gpu)
				p.TotalGPUs += 1
			}
		}
	}

	return nil
}
