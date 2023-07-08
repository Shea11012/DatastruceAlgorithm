package sort

type Cmp func(*PQ) func(i,j int) bool

func max(q *PQ) func(i,j int) bool {
	return func(i, j int) bool {
		return q.data[i] < q.data[j]
	}
}

func less(q *PQ) func(i,j int) bool {
	return func(i, j int) bool {
		return q.data[i] > q.data[j]
	}
}

// PQ 优先队列
type PQ struct {
	size int
	data []int

	// 根据传入的函数确定是最大堆还是最小堆
	cmp func(i,j int) bool
}

func NewPQ(size int,fn Cmp) *PQ {
	p := &PQ {
		data: make([]int,size),
	}

	p.cmp = fn(p)

	return p
}

func (p *PQ) Insert(vlaue int) {
	if p.size >= len(p.data)-1 {
		p.resize()
	}
	// 将新元素加到最后，然后让它上浮到正确的位置
	p.size++
	p.data[p.size] = vlaue
	p.up(p.size)
}

func (p *PQ) Delete() int {
	if len(p.data) <= 0 {
		return -1
	}

	value := p.data[1]
	// 把最大的元素换到最后，删除
	swap(p.data,1,p.size)
	p.data[p.size] = 0
	p.size--
	// 让堆顶元素下沉到正确的位置
	p.down(1)
	return value
}

func (p *PQ) GetTop() int {
	return p.data[1]
}

// parentIndex 计算父节点索引
func (q PQ) parentIndex(index int) int {
	return index / 2
}

// leftChildIndex 左子节点索引
func (q PQ) leftChildIndex(index int) int {
	return index * 2
}

// rightChildIndex 右子节点索引
func (q PQ) rightChildIndex(index int) int {
	return index * 2 + 1
}

// up 上浮
func (q *PQ) up(index int) {
	// 如果浮到堆顶，就不能再上浮
	for index > 1 && q.cmp(index,q.parentIndex(index)) {
		// 如果第i个元素比上层大，将i换上去
		swap(q.data,q.parentIndex(index),index)
		index = q.parentIndex(index)
	}
}

// down 下沉，下沉时需要比较两个子节点
func (q *PQ) down(index int) {
	// 判断是否到了堆底
	for q.leftChildIndex(index) <= q.size {
		// 先假设左子节点较大
		older := q.leftChildIndex(index)
		// 判断右子节点是否存在，并且两个子节点比较大小
		if q.rightChildIndex(index) <= q.size && q.cmp(older,q.rightChildIndex(index)) {
			older = q.rightChildIndex(index)
		}

		// 父节点与两个子节点比较
		if q.cmp(older,index) {
			break
		}
		
		swap(q.data,index,older)
		index = older
	}
}

func (p *PQ) resize() {
	newSize := len(p.data)  * 2
	newData := make([]int,newSize)
	copy(newData,p.data)
	p.data = newData
}