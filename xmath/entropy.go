package xmath

import "math"

const epsilon = 1e-6

// CalcEntropy 计算信息熵
// 熵表示任何一种能量在空间中分布的均匀程度，能量分布得越均匀，熵就越大
// 适用用统计数字分布情况
func CalcEntropy(m map[interface{}]int) float64 {
	if len(m) == 0 || len(m) == 1 {
		return 1
	}

	var entropy float64
	var total int
	for _, v := range m {
		total += v
	}

	for _, v := range m {
		proba := float64(v) / float64(total)
		if proba < epsilon {
			proba = epsilon
		}
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}
