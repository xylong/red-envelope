package algo

import (
	"math/rand"
	"time"
)

// AfterShuffle 后洗牌算法
func AfterShuffle(count, amount int64) []int64 {
	list := make([]int64, 0)
	max := amount - min*count
	remain := max
	for i := int64(0); i < count; i++ {
		x := SimpleRand(count-i, remain)
		remain -= x
		list = append(list, x)
	}
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list
}

// BeforeShuffle 先洗牌算法
func BeforeShuffle(count, amount int64) int64 {
	if count == 1 {
		return amount
	}
	list := make([]int64, 0)
	max := amount - min*count
	size := count / 2
	if size < 3 {
		size = 3
	}
	for i := int64(0); i < count; i++ {
		x := max / (i + 1)
		list = append(list, x)
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Int63n(int64(len(list)))
	x := rand.Int63n(list[index])
	return x + min
}
