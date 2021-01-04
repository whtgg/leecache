package main

import (
	"fmt"
	"sync"
)

var set = make(map[int]bool,0)
var set1 sync.Map
func printOnce(num int) {
	if _,exist := set[num];exist {
		println(num)
	}
	set[num] = true
}

func printOnce1(num int) {
	if val,ok := set1.Load(num);ok {
		fmt.Println(val)
	}
	fmt.Println(num)
	set1.Store(num,true)
}

func main() {
	bb := 1 << 9  //总结:2^9*1
	cc := 2 << 9  //总结:2^9*2
	dd := 1 << 7  //总结:2^7*1
	fmt.Println(bb)
	fmt.Println(cc)
	fmt.Println(dd)
	ee := 1 >> 2 		//总结:1/2^2 < 1 = 0
	ff := 2 >> 1 		//总结:2/2^1 = 1
	gg := 128 >> 4		//总结:128/2^4  = 8
	fmt.Println(ee)
	fmt.Println(ff)
	fmt.Println(gg)
}
