package main

import (
	"bufio"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	pf, err := os.Create("calc1.prof")
	if err != nil {
		log.Panic(err)
	}
	// runtime.SetBlockProfileRate(1)
	pprof.StartCPUProfile(pf)
	defer pprof.StopCPUProfile()

	file, err := os.Open("measurements2.txt")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// Write out station data based on the sorted name slice
	// dst := io.Discard

	// dst, err := os.Create("result.txt")
	// defer func() { dst.Close() }()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// dst := &strings.Builder{}
	dst := bufio.NewWriter(os.Stdout)

	m3(file, dst)

	// fmt.Println(dst.String())
	dst.Flush()
}
