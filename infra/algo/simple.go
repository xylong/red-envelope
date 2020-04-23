package algo

import (
	"math/rand"
	"time"
)

const min = int64(1)

// SimpleRand 简单随机算法
// 红包的数量，红包金额
// 金额单位：分
func SimpleRand(count, amount int64) int64 {
	if count == 1 {
		return amount
	}
	// 计算最大可调度金额
	max := amount - min*count
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(max) + min
	return x
}
