package common

import "strings"

type GPUType int

const (
	GPUTypeA100 GPUType = iota
	GPUTypeA800
	GPUTypeH100
	GPUTypeH800
	GPUTypeNONE
	GPUTypeH20
)

func ToGPUType(s string) GPUType {
	switch strings.ToUpper(s) {
	case "A100":
		return GPUTypeA100
	case "A800":
		return GPUTypeA800
	case "H100":
		return GPUTypeH100
	case "H800":
		return GPUTypeH800
	case "H20":
		return GPUTypeH20
	default:
		return GPUTypeNONE
	}
}
