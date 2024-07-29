package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
)

type Station2 struct {
	min, max, sum float64
	num           int
}

func m2(file io.Reader, dst io.Writer) {
	stations := map[string]*Station2{}

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
					stations[station] = &Station2{val, val, val, 1}
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

	dst.Write([]byte("{"))
	for i, name := range result {
		s := stations[name]
		minTemp := float64(s.min) / 10
		// math.Round and `+ -0` to tackle a situation where the mean is a very
		// small fraction, that will result in a -0.0 in the Sprintf below.
		meanTemp := (math.Round(float64(s.sum)/float64(s.num)) + -0) / 10
		maxTemp := float64(s.max) / 10

		if i > 0 {
			dst.Write([]byte(", "))
		}
		dst.Write([]byte(fmt.Sprintf("%v=%.1f/%.1f/%.1f", name, minTemp, meanTemp, maxTemp)))
	}
	dst.Write([]byte("}"))
}
