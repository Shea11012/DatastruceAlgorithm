# 时间轮

- 环形数据结构
- 每个刻度对应一个时间范围
- 创建定时任务时，根据距当前的相对时长，计算出需要后移的刻度值
- 遇到环形数组结尾时，重新从起点开始计算，并把执行轮次+1
- 一个刻度可能存在多个定时任务，每个刻度需要挂载一个定时任务链表

## 多级时间轮