// Go version of the C++ UserParam and NetWorkParam system
package config

import (
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type GPUType string

const (
	A100 GPUType = "A100"
	A800 GPUType = "A800"
	H100 GPUType = "H100"
	H800 GPUType = "H800"
	H20  GPUType = "H20"
	NONE GPUType = "NONE"
)

type NetWorkParam struct {
	NodeNum        uint32
	SwitchNum      uint32
	LinkNum        uint32
	TraceNum       uint32
	NvswitchNum    uint32
	GpusPerServer  uint32
	NicsPerServer  uint32
	NvlinkBw       float64
	BwPerNic       float64
	NicType        string
	Visual         bool
	DpOverlapRatio float64
	TpOverlapRatio float64
	EpOverlapRatio float64
	PpOverlapRatio float64
	GpuType        GPUType
	NVswitchs      []int
	AllGpus        [][]int
}

type UserParam struct {
	Thread       int
	Gpus         []int
	Workload     string
	Res          string
	ResFolder    string
	CommScale    int
	Mode         string
	NetWorkParam NetWorkParam
}

func ParseArgs() *UserParam {
	param := &UserParam{
		Thread:    1,
		CommScale: 1,
		Mode:      "MOCKNCCL",
	}

	// Define flags
	var gpusStr string
	var gpuType string
	flag.StringVar(&param.Workload, "workload", "", "Workload path")
	flag.StringVar(&gpusStr, "gpus", "", "Comma-separated GPU counts")
	flag.StringVar(&param.Res, "result", "None", "Result name")
	flag.StringVar(&param.ResFolder, "result_folder", "None", "Result folder")
	flag.IntVar(&param.NetWorkParam.GpusPerServer, "gpus_per_server", 8, "GPUs per server")
	flag.Float64Var(&param.NetWorkParam.NvlinkBw, "nvlink", -1.0, "Nvlink bandwidth")
	flag.Float64Var(&param.NetWorkParam.BwPerNic, "nic_busbw", -1.0, "NIC bandwidth")
	flag.IntVar((*int)(&param.NetWorkParam.NicsPerServer), "nic_per_server", 1, "NICs per server")
	flag.StringVar(&param.NetWorkParam.NicType, "nic_type", "cx7", "NIC type")
	flag.StringVar(&gpuType, "gpu_type", "NONE", "GPU type")
	flag.BoolVar(&param.NetWorkParam.Visual, "visual", false, "Enable visualization")
	flag.Float64Var(&param.NetWorkParam.DpOverlapRatio, "dp_overlap", 0.0, "DP overlap ratio")
	flag.Float64Var(&param.NetWorkParam.TpOverlapRatio, "tp_overlap", 0.0, "TP overlap ratio")
	flag.Float64Var(&param.NetWorkParam.EpOverlapRatio, "ep_overlap", 0.0, "EP overlap ratio")
	flag.Float64Var(&param.NetWorkParam.PpOverlapRatio, "pp_overlap", 1.0, "PP overlap ratio")

	flag.Parse()

	// Parse GPU list
	for _, g := range strings.Split(gpusStr, ",") {
		if val, err := strconv.Atoi(g); err == nil {
			param.Gpus = append(param.Gpus, val)
		}
	}

	switch strings.ToUpper(gpuType) {
	case "A100":
		param.NetWorkParam.GpuType = A100
	case "A800":
		param.NetWorkParam.GpuType = A800
	case "H100":
		param.NetWorkParam.GpuType = H100
	case "H800":
		param.NetWorkParam.GpuType = H800
	case "H20":
		param.NetWorkParam.GpuType = H20
	default:
		param.NetWorkParam.GpuType = NONE
	}

	// Post-process
	if len(param.Gpus) > 0 {
		param.NetWorkParam.NvswitchNum = uint32(param.Gpus[0]) / param.NetWorkParam.GpusPerServer
		param.NetWorkParam.SwitchNum = 120 + param.NetWorkParam.GpusPerServer
		param.NetWorkParam.NodeNum = param.NetWorkParam.NvswitchNum + param.NetWorkParam.SwitchNum + uint32(param.Gpus[0])
	}

	if param.Res == "None" && param.Workload != "" {
		modelInfo := filepath.Base(param.Workload)
		pattern := regexp.MustCompile(`(?i)(world_size|tp|pp|ep|gbs|mbs|seq)(\d+)`)
		matches := pattern.FindAllStringSubmatch(modelInfo, -1)

		var worldSize, tp, pp, ep, gbs, mbs int
		for _, match := range matches {
			val, _ := strconv.Atoi(match[2])
			switch match[1] {
			case "world_size":
				worldSize = val
			case "tp":
				tp = val
			case "pp":
				pp = val
			case "ep":
				ep = val
			case "gbs":
				gbs = val
			case "mbs":
				mbs = val
			}
		}
		dp := worldSize / (tp * pp)
		ga := float64(gbs) / float64(dp*mbs)
		param.Res = fmt.Sprintf("%s-tp%d-pp%d-dp%d-ga%d-ep%d-NVL%d-%.1fG-DP%.1f",
			strings.TrimSuffix(modelInfo, filepath.Ext(modelInfo)), tp, pp, dp, int(ga), ep,
			param.NetWorkParam.GpusPerServer, param.NetWorkParam.BwPerNic*8, param.NetWorkParam.DpOverlapRatio)
	}

	if param.ResFolder != "None" {
		param.Res = filepath.Join(param.ResFolder, param.Res)
	}

	return param
}
