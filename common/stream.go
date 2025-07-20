package common

// StreamState 表示流状态
type StreamState int

const (
	StreamStateCreated StreamState = iota
	StreamStateTransferring
	StreamStateReady
	StreamStateExecuting
	StreamStateZombie
	StreamStateDead
)
