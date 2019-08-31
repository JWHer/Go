//Copyright 2019. Jeongwon Her. All right reserved.

package BWT

import (
	"bytes"
	"time"
)

type index struct {
	bwt []byte
	rtt []int //virtual rotation table
	ind []int //sa table
}

//public
var DEBUG bool = false
var FindPtrs []int

//private
var seq []byte
var seqLen int
var reads [][]byte
var readNum int
var readLen int
var indexs index
var inited bool = false
var numCPUs int

func Init(sequence []byte, rds [][]byte, CPUs int) {
	seq = append(sequence, []byte("$")...)
	seqLen = len(seq)
	reads = rds
	readNum = len(rds)
	readLen = len(rds[0])
	FindPtrs = make([]int, readNum)
	numCPUs = CPUs

	//make virtual rotate
	indexs.rtt = make([]int, seqLen)
	for i := 0; i < seqLen; i++ {
		indexs.rtt[i] = i
	}
	inited = true
}

func BWTTable() {
	println("in BWTTable")
	if !inited {
		println("please init first")
		return
	}

	//make BWT tables
	var c chan int = make(chan int)
	start := time.Now()

	//sort
	cpus := numCPUs
	for {
		indexSize := readNum / cpus
		for i := 0; i < cpus; i++ {
			s := indexSize * i
			e := indexSize * (i + 1)
			if i == cpus-1 {
				e = seqLen
			}
			go func(c chan int) {
				mergeSort(indexs.rtt[s:e])
				c <- 0
			}(c)
		}
		for i := 0; i < cpus; i++ {
			<-c
		}
		cpus = cpus / 2

		if cpus == 0 {
			break
		}
	}
	indexs.bwt = make([]byte, seqLen)
	indexs.ind = make([]int, seqLen)
	for i := 0; i < seqLen; i++ {
		indexs.bwt[i] = seq[seqLen-indexs.rtt[i]-1]
		indexs.ind[i] = (seqLen - indexs.rtt[i]) % seqLen
	}
	end := time.Now()

	if DEBUG {
		println(string(indexs.bwt))
	}

	print("Works done ")
	println(end.Sub(start).String())
}

func merge(ary []int) {
	aryLen := len(ary)
	mid := aryLen / 2
	i := 0
	j := mid
	k := 0
	var sorted []int = make([]int, aryLen)

	for {
		if i >= mid || j >= aryLen {
			break
		}

		if bytes.Compare(seq[(seqLen-ary[i])%seqLen:], seq[(seqLen-ary[j])%seqLen:]) <= 0 {
			sorted[k] = ary[i]
			i++
		} else {
			sorted[k] = ary[j]
			j++
		}
		k++
	}

	if i >= mid {
		for l := j; l < aryLen; l++ {
			sorted[k] = ary[l]
			k++
		}
	} else {
		for l := i; l < mid; l++ {
			sorted[k] = ary[l]
			k++
		}
	}

	for l := 0; l < aryLen; l++ {
		ary[l] = sorted[l]
	}
}

func mergeSort(ary []int) {
	aryLen := len(ary)
	if aryLen <= 1 {
		return
	}
	mid := aryLen / 2
	mergeSort(ary[:mid])
	mergeSort(ary[mid:])
	//merge
	merge(ary)
}

func BWTPRL(mis int) {
	if mis < 0 {
		mis = 0
	}
	println("in BWTPRL")

	indSize := readNum / numCPUs
	var c chan int = make(chan int)

	start := time.Now()
	for i := 0; i < numCPUs; i++ {
		indStart := indSize * i
		indEnd := indSize * (i + 1)
		if i == numCPUs-1 {
			indEnd = readNum
		}

		if mis == 0 {
			go func(c chan int) {
				for j := indStart; j < indEnd; j++ {
					FindPtrs[j] = bWTGR(reads[j])
				}
				c <- 0
			}(c)
		} else {
			go func(c chan int) {
				for j := indStart; j < indEnd; j++ {
					FindPtrs[j] = bWTGR2(reads[j], mis)
				}
				c <- 0
			}(c)
		}
	}
	for i := 0; i < numCPUs; i++ {
		<-c
	}
	end := time.Now()

	print("Works done ")
	println(end.Sub(start).String())
}

func bWTGR(pat []byte) int {
	m := len(pat)
	n := seqLen
	l := 0
	r := n - 1
	for {
		if !(l <= r) {
			break
		}

		mid := l + (r-l)/2

		var res int
		if len(seq[indexs.ind[mid]:]) < m {
			res = bytes.Compare(pat, seq[indexs.ind[mid]:])
		} else {
			res = bytes.Compare(pat, seq[indexs.ind[mid]:indexs.ind[mid]+m])
		}

		if res == 0 {
			return indexs.ind[mid]
		} else if res < 0 {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}

func bWTGR2(pat []byte, mis int) int {
	m := len(pat)
	n := seqLen
	l := 0
	r := n - 1
	for {
		if !(l <= r) {
			break
		}

		mid := l + (r-l)/2

		var res int
		if len(seq[indexs.ind[mid]:]) < m {
			res = misComp(pat, seq[indexs.ind[mid]:], mis)
		} else {
			res = misComp(pat, seq[indexs.ind[mid]:indexs.ind[mid]+m], mis)
		}

		if res == 0 {
			return indexs.ind[mid]
		} else if res < 0 {
			r = mid - 1
		} else {
			l = mid + 1
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
