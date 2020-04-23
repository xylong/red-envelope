package algo

import (
	"math/rand"
	"time"
)

// DoubleRandom 二次随机算法
func DoubleRandom(count, amount int64) int64 {
	if count == 1 {
		return amount
	}
	max := amount - min*count
	rand.Seed(time.Now().UnixNano())
	// 1次随机，计算出一个种子金额作为基数
	seed := rand.Int63n(count*2) + 1
	n := max/seed + 1
	// 2次随机，计算红包金额序列元素
	x := rand.Int63n(n)
	return x + min
}

// DoubleAverage 2倍均值算法
func DoubleAverage(count, amount int64) int64 {
	if count == 1 {
		return amount
	}
	max := amount - min*count
	// 计算最大可用平均值
	avg := max / count
	avg2 := avg*2 + min
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(avg2) + min
	return x
}
