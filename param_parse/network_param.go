package param_parse

import "simai/common"

// NetWorkParamV1
// todo: 假如要支持ali hpn的双plane架构、或者上行链路和下行链路带宽不同的情况，则需要拓展到NetWorkParamV2
type NetWorkParamV1 struct {
	NodeNum        uint32
	SwitchNum      uint
	LinkNum        uint32
	TraceNum       uint32
	NvSwitchNum    uint32
	GpusPerServer  uint
	NicsPerServer  uint
	NvLinkBw       float64
	BwPerNic       float64
	NicType        string
	Visual         bool
	DpOverlapRatio float64
	TpOverlapRatio float64
	EpOverlapRatio float64
	PpOverlapRatio float64
	GpuType        common.GPUType
	NvSwitchs      []int
	AllGpus        [][]int
}
