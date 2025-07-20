package workload

type ParallelismPolicyType int

const (
	ParallelismPolicyMicroBenchmark ParallelismPolicyType = iota
	ParallelismPolicyData

	// todo: add more
)
