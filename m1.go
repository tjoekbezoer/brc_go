package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Station struct {
	min, max, sum float64
	num           int
}

func m1(file io.Reader, dst io.Writer) {
	stations := map[string]*Station{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stationName, temp, found := strings.Cut(line, ";")
		if !found {
			log.Panic("Huh?")
		}

		val, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			log.Panic(err)
		}

		if s, ok := stations[stationName]; ok {
			s.min = min(s.min, val)
			s.max = max(s.max, val)
			s.sum += val
			s.num++
		} else {
			stations[stationName] = &Station{val, val, val, 1}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
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
