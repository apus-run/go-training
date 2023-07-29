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

	j := 0
	for i, v := range src {
		// fmt.Printf("value = %d, value-addr = %x, slice-addr = %x\n", v, &v, &src[i])
		if i != index {
			src[j] = v
			j++
		}
	}
	src = src[:j]
	return src
}
