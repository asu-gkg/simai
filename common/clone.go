package common

type Clone interface {
	GetType() CollectiveImplementationType
	Clone() CollectiveImplementation
}
