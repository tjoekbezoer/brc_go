package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

func m1(file io.Reader, dst io.Writer) {
	stations := map[string]*Station{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		station, temp, found := strings.Cut(line, ";")
		if !found {
			log.Panic("Huh?")
		}

		val, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			log.Panic(err)
		}

		if s, ok := stations[station]; ok {
			s.min = min(s.min, val)
			s.max = max(s.max, val)
			s.sum += val
			s.num++
		} else {
			stations[station] = &Station{val, val, val, 1}
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
