package param_parse

import (
	"encoding/json"
	"fmt"
	"os"
)

type TopologyConfig struct {
	TopologyType      string           `json:"topology_type"`
	Description       string           `json:"description"`
	NetworkParameters ClusterTopoParam `json:"network_parameters"`
}

func (c TopologyConfig) Validate() error {
	return nil
}

type ClusterTopoParam struct {
	NumGPUs             int     `json:"num_gpus"`
	GPUType             string  `json:"gpu_type"`
	AllGPUs             [][]int `json:"all_gpus"`
	NumServers          int     `json:"num_servers"`
	NVSwitchNum         int     `json:"nv_switch_num"`
	NumSwitches         int     `json:"num_switches"`
	NumToRs             int     `json:"num_tors"`
	NumSpines           int     `json:"num_spines"`
	NumAggregators      int     `json:"num_aggregators"`
	GPUPerServer        int     `json:"gpu_per_server"`
	NICsPerServer       int     `json:"nics_per_server"`
	ServersPerToR       int     `json:"servers_per_tor"`
	ToRsPerSpine        int     `json:"tors_per_spine"`
	SpinesPerAggregator int     `json:"spines_per_aggregator"`
	BWPerNIC            float64 `json:"bw_per_nic"`
	NICType             string  `json:"nic_type"`
	LinkNum             int     `json:"link_num"`
	NVLinkBW            float64 `json:"nv_link_bw"`
	NICLatency          float64 `json:"nic_latency"`
	NVLinkLatency       float64 `json:"nv_link_latency"`
}

func (p ClusterTopoParam) Validate() error {
	return nil
}

func NewClusterTopoParam(filename string) (*ClusterTopoParam, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var config TopologyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	return &config.NetworkParameters, nil
}
