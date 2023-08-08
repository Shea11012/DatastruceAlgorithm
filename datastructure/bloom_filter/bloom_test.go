package bloom_filter

import "testing"

func TestBloomFilter_Add(t *testing.T) {
	b := NewBloomFilter(33, 5)
	b.Add("https://www.douyu.com/topic/dota2zfflj?rid=74960&dyshid=10c472f1-db14c456ce463ee7106d64ed00051601")
}
