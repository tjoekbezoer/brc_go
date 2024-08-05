package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
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

	// dst := io.Discard
	dst := &strings.Builder{}
	// dst := bufio.NewWriter(os.Stdout)

	m4(file, dst)

	fmt.Println(dst.String())
	// dst.Flush()
}
