package network_api

type SimComm struct {
	CommName string
}

type NetworkAPI interface {
	SimGetTime() TimeSpec
}
