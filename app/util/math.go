package util

import types "github.com/your-org/go-start-monolithic-kit/shared/types"

func GetMinInMap[K comparable, T types.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap {
		res = min(res, value)
	}
	return res
}

func GetMaxInMap[K comparable, T types.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap {
		res = max(res, value)
	}
	return res
}
