package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Station3 struct {
	min, max, sum float64
	num           int
}
type Line struct {
	name string
	temp float64
}

func m3(file io.Reader, dst io.Writer) {
	stations := map[string]*Station3{}
	input := make(chan string, 100)
	output := make(chan Line, 100)
	cpus := runtime.NumCPU()

	wg := &sync.WaitGroup{}
	wg.Add(cpus)

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			input <- scanner.Text()
		}
		close(input)

		if err := scanner.Err(); err != nil {
			log.Panic(err)
		}

		wg.Wait()
		close(output)
	}()

	for i := 0; i < cpus; i++ {
		go parse(input, output, wg)
	}

	for line := range output {
		if s, ok := stations[line.name]; ok {
			s.min = min(s.min, line.temp)
			s.max = max(s.max, line.temp)
			s.sum += line.temp
			s.num++
		} else {
			stations[line.name] = &Station3{line.temp, line.temp, line.temp, 1}
		}
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

func parse(input <-chan string, output chan<- Line, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range input {
		stationName, temp, found := strings.Cut(line, ";")
		if !found {
			log.Panic("Huh?")
		}

		val, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			log.Panic(err)
		}

		output <- Line{stationName, val}
	}
}
