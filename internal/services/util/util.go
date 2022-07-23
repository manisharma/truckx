package util

func Contains(list []int, key int) (bool, int) {
	for idx, item := range list {
		if item == key {
			return true, idx
		}
	}
	return false, -1
}
