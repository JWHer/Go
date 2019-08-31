//Copyright 2019. Jeongwon Her. All right reserved.

//Nucleic Acid Sequence
package NASeq

import (
	"bytes"
	"math/rand"
	"os"
	"strconv" //string convert
	"time"
)

//enum
type NA int

const (
	Adenine = 0 + iota
	Cytosine
	Guanine
	Thymine
)

/*class variables*/
//public
var DEBUG bool = false

//private
var seq bytes.Buffer
var reads [][]byte
var coverage []int
var genLen int
var readNum int
var readLen int
var inited bool = false
var numCPUs int
var modify bool = false

func Init(M int, N int, L int, CPUs int, mod bool) {
	genLen = M
	readNum = N
	readLen = L
	numCPUs = CPUs
	modify = mod
	inited = true
	//println("init success")
}

//Generate Sequence
func SeqGen() {
	println("In SeqGen")
	if !inited {
		panic("using function without init")
	}
	coverage = make([]int, genLen)

	rand.Seed(time.Now().UTC().UnixNano())
	//t := rand.Int31n(100000000)
	var c chan string = make(chan string)

	//println("Generating Sequence...")
	start := time.Now()
	for i := 0; i < numCPUs; i++ {
		index := (genLen / numCPUs) * i
		var loopSize int = genLen / numCPUs
		//if final routine
		if i == numCPUs-1 {
			loopSize = genLen - index
		}

		//start goroutine
		go seqGenGR(loopSize, c)
	}

	//join
	for i := 0; i < numCPUs; i++ {
		seq.WriteString(<-c)
	}
	end := time.Now()

	println("Works done " + end.Sub(start).String())
	//println(string(seq.Bytes()[:10]))
}

//seqGen goroutin
func seqGenGR(loopSize int, c chan string) {
	var ret bytes.Buffer
	for i := 0; i < loopSize; i++ {
		switch rand.Intn(4) {
		case Adenine:
			ret.WriteString("A")
		case Cytosine:
			ret.WriteString("C")
		case Guanine:
			ret.WriteString("G")
		case Thymine:
			ret.WriteString("T")
		default:
			return
		}
	}

	c <- ret.String()
}

//Generate Reads
func ReadGen() {
	println("in ReadGen")
	if !inited {
		panic("using function without init")
	}

	if seq.Len() == 0 {
		println("please SeqGen first")
		return
	}

	var c chan int = make(chan int)
	reads = make([][]byte, readNum)
	var size = genLen / numCPUs

	start := time.Now()
	for i := 0; i < numCPUs; i++ {
		var index int = readNum / numCPUs * i
		var end int
		var start int = genLen / numCPUs * i
		if i == numCPUs-1 {
			end = readNum
			size = genLen - start - readLen
		} else {
			end = readNum / numCPUs * (i + 1)
		}

		go readGenGR(index, end, start, size, c)
	}
	for i := 0; i < numCPUs; i++ {
		<-c
	}
	end := time.Now()

	println("Works done " + end.Sub(start).String())

}

func readGenGR(index int, end int, start int, size int, c chan int) {
	if size < 0 {
		size = 1
	}
	for i := index; i < end; i++ {
		cPtr := rand.Int31n(int32(size)) + int32(start)
		ePtr := cPtr + int32(readLen)
		if DEBUG {
			print(cPtr)
			print(" to ")
			println(ePtr)
		}

		reads[i] = make([]byte, readLen)
		copy(reads[i], seq.Bytes()[cPtr:ePtr])

		if modify {
			//1%
			if rand.Intn(100)>readLen/100{
				switch reads[i][rand.Intn(readLen)]{
				case 'A':
					reads[i][rand.Intn(readLen)]='T'
				case 'C':
					reads[i][rand.Intn(readLen)]='G'
				case 'G':
					reads[i][rand.Intn(readLen)]='C'
				case 'T':
					reads[i][rand.Intn(readLen)]='A'
				default:
					reads[i][rand.Intn(readLen)]='U'
				}
			}
		}

		if DEBUG {
			println(string(reads[i]))
		}

		for j := cPtr; j < ePtr; j++ {
			coverage[j]++
		}
	}
	c <- 0
}

//Save sequence
func Save(filename string) {
	println("Saving...")

	//open file
	file, err := os.Create(filename)
	if err != nil {
		println("fail to create file")
		panic(err)
	}
	//close file at the end of main
	defer file.Close()

	_, err = file.Write(seq.Bytes())
	if err != nil {
		println("fail to write file")
		panic(err)
	}
	file.WriteString("\r\n")
	for i := 0; i < seq.Len(); i++ {
		file.WriteString(strconv.Itoa(coverage[i]))
	}
	file.WriteString("\r\n")
	for i := 0; i < readNum; i++ {
		file.Write(reads[i])
		file.WriteString("\r\n")
	}

	println("saved")
}

func GetSeq() []byte {
	return seq.Bytes()
}

func GetReads() [][]byte {
	return reads
}
