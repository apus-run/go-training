package slice

// RemoveV1 删除一个 value 值, 如果有多个, 只删除第一个
func RemoveV1(src []int64, value int64) []int64 {
	length := len(src)
	if length == 0 {
		return nil
	}

	for i, v := range src {
		// fmt.Printf("value = %d, value-addr = %x, slice-addr = %x\n", v, &v, &src[i])
		if v == value {
			src = append(src[:i], src[i+1:]...)
			return src
		}
	}

	return src
}
