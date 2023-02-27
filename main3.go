package main

import "fmt"

func main() {
	// 遍历, 决定处理第几行
	for i := 1; i < 10; i++ {
		// 遍历, 决定这一行有多少列
		for j := 1; j <= i; j++ {
			fmt.Printf("%d*%d=%d ", j, i, i*j)
		}
		// 手动生成回车
		fmt.Println()
	}
}
