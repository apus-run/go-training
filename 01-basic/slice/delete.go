package slice

// Delete 删除 index 位置的元素
func Delete[T any](src []T, index int) ([]T, T) {
	length := len(src)
	if index < 0 || index >= length {
		var zero T
		return nil, zero
	}

	t := src[index]
	res := src[:index]
	res = append(res, src[index+1:]...)

	return res, t
}

// DeleteV1 删除 index 位置的元素
func DeleteV1(src []int64, index int) []int64 {
	length := len(src)
	if index < 0 || index >= length {
		return nil
	}

	res := src[:index]
	res = append(res, src[index+1:]...)

	return res
}

// DeleteV2 删除 index 位置的元素
func DeleteV2(src []int64, index int) []int64 {
	length := len(src)
	if index < 0 || index >= length {
		return nil
	}
	
	// 从 index+1 开始，将后面的元素向前移动一位
	for i := index + 1; i < length; i++ {
		src[i-1] = src[i]
	}

	return src[:length-1]
}
