//Copyright 2019. Jeongwon Her. All right reserved.
package Naive

import "time"

//public
var FindPtrs []int

//private
var inited bool = false
var numCPUs int
var seq []byte
var reads [][]byte

func Init(seqeunce []byte, rds [][]byte, cpus int) {
	seq = seqeunce
	reads = rds
	FindPtrs = make([]int, len(reads))
	numCPUs = cpus
	inited = true
}

func NaivePRL(mis int) {
	println("in NaivePRL")
	if !inited {
		println("please init first")
		return
	}

	var c chan int = make(chan int)
	indSize := len(reads) / numCPUs

	start := time.Now()
	for i := 0; i < numCPUs; i++ {
		s := indSize * i//start index
		e := indSize * (i + 1)//end index
		if i == numCPUs-1 {
			e = len(reads)
		}

		go func(c chan int) {
			for j := s; j < e; j++ {
				FindPtrs[j] = naive(seq, reads[j], mis)
			}
			c <- 0
		}(c)
	}
	for i := 0; i < numCPUs; i++ {
		<-c
	}
	end := time.Now()

	print("Works done ")
	println(end.Sub(start).String())
}

func naive(txt []byte, pat []byte, mis int) int {
	txtLen := len(txt)
	patLen := len(pat)

	if patLen > txtLen {
		return -1
	}
	if mis < 0 {
		mis = 0
	}

	for i := 0; i < txtLen-patLen; i++ {
		ret := misComp(txt[i:i+patLen], pat, mis)
		if ret == 0 {
			return i
		}
	}

	return -1
}

func misComp(pat1 []byte, pat2 []byte, mis int) int {

	var maxLen int
	if len(pat1) > len(pat2) {
		maxLen = len(pat2)
	} else {
		maxLen = len(pat1)
	}

	first := true
	var rem int
	for i := 0; i < maxLen; i++ {

		if pat1[i] > pat2[i] {
			if first {
				first = false
				rem = 1
			}
			if mis > 0 {
				mis--
				continue
			}
			return rem
		} else if pat1[i] < pat2[i] {
			if first {
				first = false
				rem = -1
			}
			if mis > 0 {
				mis--
				continue
			}
			return rem
		}
	}

	if len(pat1) > len(pat2) {
		return 1
	} else if len(pat1) < len(pat2) {
		return -1
	}

	return 0
}
