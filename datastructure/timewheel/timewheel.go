package timewheel

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type task struct {
	// 轮次
	cycle int
	// 定时任务在插槽中的位置
	pos int
	key string
	job func()
}

type TimeWheel struct {
	// 时间轮运行时间间隔
	interval time.Duration
	// 定时器
	ticker *time.Ticker

	// 插槽(循环数组)，每个插槽由一个链表连接
	slots []*list.List
	// 当前运行到的插槽位置
	curSlot int
	// key与task的映射，方便在list中删除任务节点
	keyToTask map[string]*list.Element

	// 确保停止操作只执行一次
	sync.Once
	// stopC 停止时间轮
	stopC chan struct{}
	// addTaskC 新增任务
	addTaskC chan *task
	// removeTaskC 删除任务
	removeTaskC chan string
}

var (
	defaultSlotNum  = 10
	defaultInterval = time.Second
)

func NewTimeWheel(slotNum int, interval time.Duration) *TimeWheel {
	if slotNum <= 0 {
		slotNum = defaultSlotNum
	}

	if interval <= 0 {
		interval = defaultInterval
	}

	t := &TimeWheel{
		interval:    interval,
		stopC:       make(chan struct{}),
		keyToTask:   map[string]*list.Element{},
		slots:       make([]*list.List, slotNum),
		addTaskC:    make(chan *task),
		removeTaskC: make(chan string),
	}

	for i := 0; i < slotNum; i++ {
		t.slots[i] = list.New()
	}

	return t
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *TimeWheel) Stop() {
	tw.Do(func() {
		tw.ticker.Stop()
		close(tw.stopC)
	})
}

func (tw *TimeWheel) AddTask(key string, delay time.Duration, job func()) {
	pos, cycle := tw.getPosAndCircle(delay)

	tw.addTaskC <- &task{
		pos:   pos,
		cycle: cycle,
		job:   job,
		key:   key,
	}
}

func (tw *TimeWheel) getPosAndCircle(d time.Duration) (int, int) {
	delay := int(d.Seconds())
	interval := int(tw.interval.Seconds())
	// 定时任务需要执行的轮次
	cycle := delay / interval / len(tw.slots)
	// 定时任务在循环数组的位置
	pos := (tw.curSlot + delay/interval) % len(tw.slots)

	return pos, cycle
}

func (tw *TimeWheel) RemoveTask(key string) {
	tw.removeTaskC <- key
}

func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.stopC:
			return
		case task := <-tw.addTaskC:
			tw.addTask(task)
		case key := <-tw.removeTaskC:
			tw.removeTask(key)
		case <-tw.ticker.C:
			tw.tickHandler()
		}
	}
}

func (tw *TimeWheel) addTask(t *task) {
	l := tw.slots[t.pos]
	if _, ok := tw.keyToTask[t.key]; ok {
		tw.removeTask(t.key)
	}

	eTask := l.PushBack(t)
	tw.keyToTask[t.key] = eTask
}

func (tw *TimeWheel) removeTask(key string) {
	eTask, ok := tw.keyToTask[key]
	if !ok {
		return
	}

	// 从映射里删除
	delete(tw.keyToTask, key)
	t, _ := eTask.Value.(*task)
	// 从插槽中删除
	tw.slots[t.pos].Remove(eTask)
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.curSlot]
	tw.curSlot = (tw.curSlot + 1) % len(tw.slots)
	tw.execute(l)
}

func (tw *TimeWheel) execute(l *list.List) {
	for e := l.Front(); e != nil; {
		t, _ := e.Value.(*task)
		// 任务轮次还没到，则只执行轮次扣减
		if t.cycle > 0 {
			t.cycle--
			e = e.Next()
			continue
		}

		// 给以满足条件的定时任务，开启goroutine负责执行任务
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()

			job := t.job
			job()
		}()

		// 任务已执行，将任务从对应插槽和映射中删除
		next := e.Next()
		l.Remove(e)
		delete(tw.keyToTask, t.key)
		e = next
	}
}
