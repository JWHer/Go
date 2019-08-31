// Copyright 2019. Jeongwon Her. All right reserved.
package main

import (
	"BWT"
	"KMP"
	"NASeq"
	"Naive"
	"runtime"
	"time"
)

func main() {
	println("In main")
	//get cpu info and set max process
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS((numCPUs))

	//One Billion challenge
	genLen := 1000
	readLen := 50
	readNum := genLen / readLen * 2
	NASeq.Init(genLen, readNum, readLen, numCPUs, false)

	NASeq.SeqGen()
	seq := NASeq.GetSeq()

	NASeq.ReadGen()
	reads := NASeq.GetReads()

	NASeq.Save("gen.txt")

	KMP.Init(numCPUs, seq, reads)
	KMP.DEBUG = false
	KMP.KMPPRL()

	//no parallel
	println("IN KMP")
	start := time.Now()
	for i := 0; i < len(reads); i++ {
		KMP.KMP(seq, reads[i])
	}
	end := time.Now()
	print("Works done ")
	println(end.Sub(start).String())

	Naive.Init(seq, reads, numCPUs)
	Naive.NaivePRL(0)

	BWT.Init(seq, reads, numCPUs)
	BWT.DEBUG = false
	BWT.BWTTable()
	BWT.BWTPRL(0)

	//for debugging match pointers
	if false {
		for i := 0; i < len(reads); i++ {
			print(Naive.FindPtrs[i])
			print(",")
			print(BWT.FindPtrs[i])
			print(" ")
		}
		println()
	}
}
