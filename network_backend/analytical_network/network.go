package analytical_network

import (
	"simai/memory"
	"simai/network_backend/network_api"
)

type AnalyticalNetwork struct {
	NpuOffset int
}

func NewAnalyticalNetwork(localRank int) *AnalyticalNetwork {
	return &AnalyticalNetwork{
		NpuOffset: 0,
	}
}

// SimCommSize 输出通信器的大小
func (an *AnalyticalNetwork) SimCommSize(comm network_api.SimComm) (int, error) {
	return 0, nil
}

func (an *AnalyticalNetwork) SimFinish() error {
	return nil
}

func (an *AnalyticalNetwork) SimTimeResolution() float64 {
	return 0.0
}

func (an *AnalyticalNetwork) SimInit(mem memory.MemAPI) error {
	return nil
}
func (an *AnalyticalNetwork) SimSchedule(delta network_api.TimeSpec, callback func(interface{}), arg interface{}) {
	panic("implement me")
}

func (an *AnalyticalNetwork) SimSend(buffer []byte, count uint64, dataType network_api.RequestType, dst int, tag int,
	request *network_api.SimRequest, callback func(interface{}), arg interface{}) error {
	return nil
}

func (an *AnalyticalNetwork) SimReceive(buffer []byte, count uint64, dataType network_api.RequestType, src int, tag int,
	request *network_api.SimRequest, callback func(interface{}), arg interface{}) error {
	return nil
}

func (an *AnalyticalNetwork) SimGetTime() network_api.TimeSpec {
	panic("implement me")
}
