package main

import (
	"fmt"
	"strings"
)

func main() {
	builder := strings.Builder{}

	// 添加字符串片段到构造器
	builder.WriteString("Hello, ")
	builder.WriteString("World!")

	// 获取构建好的字符串
	result := builder.String()

	fmt.Println(result)
}
