package common

// InjectionPolicy 表示注入策略
// 注入策略是指packet的注入配置
type InjectionPolicy int

const (
	InjectionPolicyInfinite InjectionPolicy = iota
	InjectionPolicyAggressive
	InjectionPolicySemiAggressive
	InjectionPolicyExtraAggressive
	InjectionPolicyNormal
)
