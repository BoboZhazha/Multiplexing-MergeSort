package main

import (
	"Multiplexing-MergeSort/pipeline"
	"bufio"
	"fmt"
	"os"
)

func main() {
	p := createPipeline("small.in", 512, 4)
	writeFile(p, "small.out")
	printFile("small.out")
}

func printFile(filenane string) {
	file, err := os.Open(filenane)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := pipeline.ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
	}
}

func writeFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	pipeline.WriterSink(writer, p)
}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int {

	// 计算一个每一块的大小
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)

		sortResults = append(sortResults, pipeline.InMemSort(source))

	}

	return pipeline.MergeN(sortResults...)

}
