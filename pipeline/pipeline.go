package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

/**
这个函数把所有的数据扔进channel里
*/
// a ... 叫可变长参数, 返回一个int类型的channel, 读作chan of int
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			// 在goroutine里把数据送进去
			out <- v
		}
		// 一般来说channel不需要close, 但是并行计算的时候要告诉别人我结束了
		close(out)
	}()
	//channel
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// Read into memory
		a := []int{}
		for v := range in {
			// a是slice, 是一个不可变对象,所以append以后要再接收一下
			a = append(a, v)
		}
		// Sort
		sort.Ints(a)
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// 核心方法
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2

		// 只要有一个channel有数据, 那么就可以送出数据
		for ok1 || ok2 {
			// 什么时候送出channel1的数据呢? 答案是 chan2没有数据的时候 或者 chan1有数据的时候并且chan1的数据比chan2小的时候,等于取chan1的数据(2路归并)
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
