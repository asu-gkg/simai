package common

type GPUType int

const (
	GPUTypeA100 GPUType = iota
	GPUTypeA800
	GPUTypeH100
	GPUTypeH800
	GPUTypeNONE
	GPUTypeH20
)
