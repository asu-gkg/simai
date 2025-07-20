package common

// DirectCollectiveImplementation 表示直接集合实现
type DirectCollectiveImplementation struct {
	*BaseCollectiveImplementation
	DirectCollectiveWindow int
}

// NewDirectCollectiveImplementation 创建新的直接集合实现
func NewDirectCollectiveImplementation(implType CollectiveImplementationType, window int) *DirectCollectiveImplementation {
	return &DirectCollectiveImplementation{
		BaseCollectiveImplementation: NewBaseCollectiveImplementation(implType),
		DirectCollectiveWindow:       window,
	}
}

func (d *DirectCollectiveImplementation) Clone() *DirectCollectiveImplementation {
	return &DirectCollectiveImplementation{
		BaseCollectiveImplementation: NewBaseCollectiveImplementation(d.implType),
		DirectCollectiveWindow:       d.DirectCollectiveWindow,
	}
}
