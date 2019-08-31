//Copyright 2019. Jeongwon Her. All right reserved.

package KMP

import (
	"time"
)

//public
var DEBUG bool = false
var FindPtrs []int

//private
var numCPUs int
var readNum int
var inited = false
var seq []byte
var reads [][]byte

func Init(CPUs int, sequence []byte, rds [][]byte) {
	numCPUs = CPUs
	readNum = len(rds)
	seq = sequence
	reads = rds
	FindPtrs = make([]int, len(rds))
	inited = true
}

func KMP(text []byte, pattern []byte) int {
	//println("in KMP")

	//making sp
	var sp []int = make([]int, len(pattern))
	for i := 0; i < len(pattern); i++ {
		sp[i] = -1
	}
	i := 0
	j := 1
	for j < len(pattern) {
		if pattern[i] == pattern[j] {
			sp[j] = i
			i++
			j++
		} else if i == 0 {
			j++
		} else {
			i = sp[i-1] + 1
		}
	}

	//return pointer
	s := 0
	isExist := false
	for i := 0; s < len(text); {
		if text[s] == pattern[i] {
			s++
			i++
			if i == len(pattern) {
				isExist = true
				break
			}
		} else if i == 0 {
			s++
		} else {
			i = sp[i-1] + 1
		}
	}

	if isExist {
		return s - len(pattern)
	} else {
		return -1
	}
}

func kMPGR(text []byte, patterns [][]byte, index int, end int, c chan int) {
	for i := index; i < end; i++ {
		ptr := KMP(text, patterns[i])
		FindPtrs[i] = ptr
		if DEBUG {
			print(ptr)
			print(" ")
		}
	}
	c <- 0
}

func KMPPRL() {
	println("in KMPPRL")
	if !inited {
		println("please init first")
		return
	}

	var c chan int = make(chan int)
	start := time.Now()
	for i := 0; i < numCPUs; i++ {
		index := readNum / numCPUs * i
		end := readNum / numCPUs * (i + 1)
		if i == numCPUs-1 {
			end = readNum
		}
		go kMPGR(seq, reads, index, end, c)
	}
	for i := 0; i < numCPUs; i++ {
		<-c
	}
	end := time.Now()
	if DEBUG {
		println()
	}

	print("Works done ")
	println(end.Sub(start).String())
}
