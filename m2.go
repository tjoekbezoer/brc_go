package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

func m2(file io.Reader, dst io.Writer) {
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
			if ch == ';' || (isTemp && ch == '.') {
				isTemp = true
			} else if ch == 10 {
				station := string(stationName)
				val, _ := strconv.ParseFloat(string(temp), 64)

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
	}

	// Sort station names
	result := make([]string, 0, len(stations))
	for n := range stations {
		result = append(result, n)
	}
	sort.Sort(sort.StringSlice(result))

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
