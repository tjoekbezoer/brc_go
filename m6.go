package main

// Implement a custom hash map to process the station data

import (
	"brc/stations"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
)

func m6(file *os.File, dst io.Writer) {
	fi, _ := file.Stat()
	fileName := fi.Name()

	hash := stations.NewMap()
	output := make(chan *stations.Map, 100)

	cpus := runtime.NumCPU()
	parts := splitFile(file, cpus)
	for _, part := range parts {
		go processPart6(fileName, part, output)
	}

	// Receive the input of all goroutines, and add them to
	// the main station map.
	for range parts {
		partResult := <-output
		for s := range partResult.Sorted() {
			err := hash.Set([]byte(s.Name), s.Min, s.Max, s.Sum, s.Num)
			if err != nil {
				panic(err)
			}
		}
	}
	close(output)

	dst.Write([]byte("{"))
	i := 0
	for s := range hash.Sorted() {
		minTemp := float64(s.Min) / 10
		meanTemp := (math.Round(float64(s.Sum)/float64(s.Num)) + -0) / 10
		maxTemp := float64(s.Max) / 10
		var prefix []byte

		if i > 0 {
			prefix = []byte(", ")
		}
		dst.Write([]byte(fmt.Sprintf("%s%s=%.1f/%.1f/%.1f", prefix, s.Name, minTemp, meanTemp, maxTemp)))
		i++
	}
	dst.Write([]byte("}"))
}

func processPart6(fileName string, p part, output chan *stations.Map) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	if _, err := file.Seek(p.offset, io.SeekStart); err != nil {
		panic(err)
	}

	partReader := io.LimitReader(file, p.size)
	scanner := bufio.NewScanner(partReader)
	hash := stations.NewMap()
	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, tempStr, found := bytes.Cut(line, []byte(";"))
		if !found {
			panic("Huh?")
		}

		var (
			isNegative = false
			temp       int
		)

		for _, ch := range tempStr {
			if ch == '.' {
				continue
			} else if ch == '-' {
				isNegative = true
			} else {
				temp = temp*10 + int(ch-48)
			}
		}
		if isNegative {
			temp *= -1
		}

		// Repeating temp 3 times seems unnecessary, but this function
		// is also used up top to union all parts together. In that case
		// the values for all parameters will be different.
		err = hash.Set(stationName, temp, temp, temp, 1)
		if err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	output <- hash
}
