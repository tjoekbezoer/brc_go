package main

// Implement a custom hash map to process the station data

import (
	"brc/stations"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
)

func m5(file *os.File, dst io.Writer) {
	fi, _ := file.Stat()
	fileName := fi.Name()

	hash := stations.NewMap()
	output := make(chan *stations.Map, 100)

	cpus := runtime.NumCPU()
	parts := splitFile(file, cpus)
	for _, part := range parts {
		go processPart5(fileName, part, output)
	}

	// Receive the input of all goroutines, and add them to
	// the main station map.
	for range parts {
		partResult := <-output
		for s := range partResult.Sorted() {
			err := hash.Set([]byte(s.Name), s.Min, s.Max, s.Sum, s.Num)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	close(output)

	dst.Write([]byte("{"))
	i := 0
	for s := range hash.Sorted() {
		minTemp := float64(s.Min)
		meanTemp := (math.Round(float64(s.Sum)/float64(s.Num)*10) + -0) / 10
		maxTemp := float64(s.Max)
		var prefix []byte

		if i > 0 {
			prefix = []byte(", ")
		}
		dst.Write([]byte(fmt.Sprintf("%s%s=%.1f/%.1f/%.1f", prefix, s.Name, minTemp, meanTemp, maxTemp)))
		i++
	}
	dst.Write([]byte("}"))
}

func processPart5(fileName string, p part, output chan *stations.Map) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	if _, err := file.Seek(p.offset, io.SeekStart); err != nil {
		log.Panic(err)
	}

	partReader := io.LimitReader(file, p.size)
	scanner := bufio.NewScanner(partReader)
	hash := stations.NewMap()
	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, tempStr, found := bytes.Cut(line, []byte(";"))
		if !found {
			log.Panic("Huh?")
		}

		// TODO: Improve performance by handrolling conversion?
		temp, err := strconv.ParseFloat(string(tempStr), 64)
		if err != nil {
			log.Panic(err)
		}

		// Repeating temp 3 times seems unnecessary, but this function
		// is also used up top to union all parts together. In that case
		// the values for all parameters will be different.
		hash.Set(stationName, temp, temp, temp, 1)
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	output <- hash
}
