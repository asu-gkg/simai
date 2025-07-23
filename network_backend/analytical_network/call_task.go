package analytical_network

type CallTask struct {
	time   uint64            // 执行时间
	funPtr func(interface{}) // 函数指针
	funArg interface{}       // 函数参数
	index  int               // 在堆中的索引，用于heap包
}
type CallTaskHeap []*CallTask

// Len 返回堆的长度
func (h CallTaskHeap) Len() int { return len(h) }

// Less 定义堆的排序规则（最小堆，按时间排序）
func (h CallTaskHeap) Less(i, j int) bool {
	return h[i].time < h[j].time
}

// Swap 交换两个元素
func (h CallTaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Push 向堆中添加元素
func (h *CallTaskHeap) Push(x interface{}) {
	n := len(*h)
	task := x.(*CallTask)
	task.index = n
	*h = append(*h, task)
}

// Pop 从堆中移除并返回最小元素
func (h *CallTaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	task.index = -1 // 标记为已移除
	*h = old[0 : n-1]
	return task
}
