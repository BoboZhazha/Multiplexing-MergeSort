package main

import (
	"Multiplexing-MergeSort/pipeline"
	"bufio"
	"fmt"
	"os"
)

func main() {

	const filename = "small.in"
	const count = 64
	// file就是一个io.Writer
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	p := pipeline.RandomSource(count)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	p = pipeline.ReaderSource(bufio.NewReader(f), -1)
	for v := range p {
		fmt.Println(v)
	}

}

func mergeDemo() {
	//这个p是一个channel, 是可以不断的从里面获得数据的
	p := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(3, 2, 6, 7, 4)),
		pipeline.InMemSort(pipeline.ArraySource(8, 123, 16, 77, 94)),
	)
	for i := range p {
		fmt.Println(i)
	}
}
