package utils

// 元素是否在数组中存在
func InArray(needle interface{}, haystack []interface{}) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// 元素是否在数组中存在
func InArrayString(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// 查询元素是否在数组中，存在则从数组中删除
func RemoveIfExist(arr []string, elem string) []string {
	index := -1
	for i, val := range arr {
		if val == elem {
			index = i
			break
		}
	}
	if index == -1 {
		// 如果元素不存在，则返回原数组
		return arr
	}
	// 从数组中删除元素
	arr = append(arr[:index], arr[index+1:]...)
	return arr
}

// IntersectionString 获取多个数组中相同的元素组成的新数组
func IntersectionString(arrays [][]string) []string {
	if len(arrays) == 0 {
		return nil
	}
	if len(arrays) == 1 {
		return arrays[0]
	}

	// 将第一个数组作为基准，遍历其中的元素
	// 如果元素在后面的数组中都有出现，则加入结果数组
	var result []string
	for _, elem := range arrays[0] {
		if contains(arrays[1:], elem) {
			result = append(result, elem)
		}
	}
	return result
}

// 判断数组中是否包含某个元素
func contains(arrays [][]string, elem string) bool {
	count := 0
	for _, arr := range arrays {
		for _, e := range arr {
			if elem == e {
				count++
				break
			}
		}
	}
	if count == len(arrays) {
		return true
	}
	return false
}

// UnsetArray 删除数组指定元素
func UnsetArray[T string | int](arr []T, index int) []T {
	j := 0
	for k, v := range arr {
		if k != index {
			arr[j] = v
			j++
		}
	}
	return arr[:j]
}

func Intersection(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	res := make([]string, 0)
	for _, v := range arr1 {
		m[v] = true
	}
	for _, v := range arr2 {
		if m[v] {
			res = append(res, v)
		}
	}
	return res
}

// ArrayUnique 数组去重
func ArrayUnique[T comparable](arr []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, item := range arr {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// HasDuplicate 是否包含包含重复元素
func HasDuplicate[T comparable](arr []T) bool {
	seen := make(map[T]bool)
	for _, num := range arr {
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}
