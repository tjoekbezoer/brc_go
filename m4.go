package main

// Introduces concurrency to process the input file in multiple
// goroutines simultaneously.

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Station4 struct {
	min, max, sum float64
	num           int
}

type part struct {
	offset, size int64
}

func m4(file *os.File, dst io.Writer) {
	fi, _ := file.Stat()
	fileName := fi.Name()

	stations := map[string]*Station4{}
	output := make(chan map[string]*Station4)

	cpus := runtime.NumCPU()
	parts := splitFile(file, cpus)
	for _, part := range parts {
		go processPart(fileName, part, output)
	}

	// Receive the input of all goroutines, and add them to
	// the main station map.
	for range parts {
		partResult := <-output
		for name, station := range partResult {
			if s, ok := stations[name]; ok {
				s.min = min(s.min, station.min)
				s.max = max(s.max, station.max)
				s.sum += station.sum
				s.num += station.num
			} else {
				stations[name] = &Station4{
					station.min, station.max, station.sum, station.num,
				}
			}
		}
	}
	close(output)

	// Sort station names
	result := make([]string, 0, len(stations))
	for n := range stations {
		result = append(result, n)
	}
	sort.Sort(sort.StringSlice(result))

	dst.Write([]byte("{"))
	for i, name := range result {
		s := stations[name]
		minTemp := float64(s.min)
		meanTemp := (math.Round(float64(s.sum)/float64(s.num)*10) + -0) / 10
		maxTemp := float64(s.max)
		var prefix []byte

		if i > 0 {
			prefix = []byte(", ")
		}
		dst.Write([]byte(fmt.Sprintf("%s%v=%.1f/%.1f/%.1f", prefix, name, minTemp, meanTemp, maxTemp)))
	}
	dst.Write([]byte("}"))
}

func processPart(fileName string, p part, output chan map[string]*Station4) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	if _, err := file.Seek(p.offset, io.SeekStart); err != nil {
		log.Panic(err)
	}

	partReader := io.LimitReader(file, p.size)
	scanner := bufio.NewScanner(partReader)
	stations := map[string]*Station4{}
	for scanner.Scan() {
		line := scanner.Text()
		stationName, tempStr, found := strings.Cut(line, ";")
		if !found {
			log.Panic("Huh?")
		}

		// TODO: Improve performance by handrolling conversion?
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			log.Panic(err)
		}

		if s, ok := stations[stationName]; ok {
			s.min = min(s.min, temp)
			s.max = max(s.max, temp)
			s.sum += temp
			s.num++
		} else {
			stations[stationName] = &Station4{temp, temp, temp, 1}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	output <- stations
}
