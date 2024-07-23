package main

import (
	// "bufio"
	"bytes"
	"fmt"
	"io"
	// "strings"

	// _ "net/http/pprof"
	"runtime/pprof"
	"sort"

	"log"
	"os"
	"strconv"
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

	type Station struct {
		min, max, sum int64
		num           int
	}
	stations := map[string]*Station{}

	buf := make([]byte, 1024*1024*10)
	skipNum := 0
	for {
		n, err := file.Read(buf[skipNum:])
		if err != nil && err != io.EOF {
			log.Panic(err)
		}
		if n+skipNum == 0 {
			break
		}

		chunk := buf[:skipNum+n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline == -1 {
			log.Panic("Huh?")
		}

		remainder := chunk[newline+1:]
		chunk = chunk[:newline+1]

		isTemp := false
		var temp []byte
		var stationName []byte
		for _, ch := range chunk {
			if ch == ';' || ch == '.' {
				isTemp = true
			} else if ch == 10 {
				station := string(stationName)
				val, _ := strconv.ParseInt(string(temp), 10, 64)

				if s, ok := stations[station]; ok {
					s.min = min(s.min, val)
					s.max = max(s.max, val)
					s.sum += val
					s.num++
				} else {
					stations[station] = &Station{val, val, val, 1}
				}

				isTemp = false
				temp = []byte{}
				stationName = []byte{}
			} else if isTemp {
				temp = append(temp, ch)
			} else {
				stationName = append(stationName, ch)
			}
		}

		skipNum = copy(buf, remainder)
		// fmt.Println(n, newline, skipNum)
	}

	// type Station struct {
	// 	min, max, sum float64
	// 	num           int
	// }
	// stations := map[string]*Station{}
	//
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	station, temp, found := strings.Cut(line, ";")
	// 	if !found {
	// 		log.Panic("Huh?")
	// 	}
	//
	// 	val, err := strconv.ParseFloat(temp, 64)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	//
	// 	if s, ok := stations[station]; ok {
	// 		s.min = min(s.min, val)
	// 		s.max = max(s.max, val)
	// 		s.sum += val
	// 		s.num++
	// 	} else {
	// 		stations[station] = &Station{val, val, val, 1}
	// 	}
	// }
	//
	// if err := scanner.Err(); err != nil {
	// 	log.Panic(err)
	// }

	// Sort station names
	result := make([]string, 0, len(stations))
	for n := range stations {
		result = append(result, n)
	}
	sort.Sort(sort.StringSlice(result))

	// Write out station data based on the sorted name slice
	dst := io.Discard
	// dst, err := os.Create("result.txt")
	// defer func() { dst.Close() }()
	if err != nil {
		log.Panic(err)
	}

	dst.Write([]byte("{ "))
	for i, name := range result {
		station := stations[name]
		mean := float64(station.sum) / float64(station.num) / 10
		var prefix, suffix []byte

		if i < len(result)-1 {
			suffix = []byte(", ")
		}
		dst.Write([]byte(fmt.Sprintf("%s%v=%.1f/%.1f/%.1f%s", prefix, name, float64(station.min)/10, mean, float64(station.max)/10, suffix)))
	}
	dst.Write([]byte(" }"))

	fmt.Println(result[:10])
}
