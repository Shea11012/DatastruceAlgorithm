package bloom_filter

import (
	"github.com/vcaesar/murmur"
)

type BloomFilter struct {
	// m: bitmap 长度
	// k：哈希函数个数
	// n：输入元素个数
	m, k, n uint32
	// 长度为 m/32+1
	bitmap []int
}

func NewBloomFilter(m, k uint32) *BloomFilter {
	return &BloomFilter{
		m:      m,
		k:      k,
		bitmap: make([]int, m/32+1),
	}
}

func (b *BloomFilter) Add(val string) {
	b.n++
	for _, offset := range b.getEncrypt(val) {
		// offset >> 5 等于 offset / 32，相当于将offset转为对应的数组索引位置
		index := offset >> 5
		// 获得位图索引内的具体偏移量
		bitOffset := offset & 31
		b.bitmap[index] |= 1 << bitOffset
	}
}

func (b *BloomFilter) Exist(val string) bool {
	for _, offset := range b.getEncrypt(val) {
		index := offset >> 5
		bitOffset := offset & 31
		if b.bitmap[index]&(1<<bitOffset) == 0 {
			return false
		}
	}

	return true
}

func (b *BloomFilter) getEncrypt(src string) []uint32 {
	encrypteds := make([]uint32, 0, b.k)
	for i := uint32(0); i < b.k; i++ {
		encrypted := murmur.Sum32(src, i)
		// 通过对m取模，将哈希值限制在位图索引的范围内
		encrypteds = append(encrypteds, encrypted%b.m)
	}
	return encrypteds
}
