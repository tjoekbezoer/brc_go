package main

import (
	"fmt"
	"strings"

	// _ "net/http/pprof"
	"runtime/pprof"

	"log"
	"os"
)

func main() {
	pf, err := os.Create("calc1.prof")
	if err != nil {
		log.Panic(err)
	}
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

	dst := &strings.Builder{}

	m3(file, dst)

	fmt.Println(dst.String())
}
