package stations

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"iter"
	"slices"
)

const BucketSize = 1 << 14

type Station struct {
	Name          []byte
	Min, Max, Sum float64
	Num           int
}

type Map struct {
	len   int
	items []*Station
}

func NewMap() *Map {
	// Create a slice that will hold all the possible stations
	// the input file can contain.
	stations := &Map{
		items: make([]*Station, BucketSize),
	}
	return stations
}

func (s *Map) Set(name []byte, minTemp, maxTemp, sumTemp float64, num int) error {
	// Arbitrary limit to the amount of stations this map can contain,
	// since the challenge states that there is a maximum of 10.000
	// different station names.
	if s.len > 11000 {
		return errors.New("too many stations in hash map")
	}

	hash := fnv.New64()
	hash.Write(name)
	bucket := hash.Sum64() & (BucketSize - 1)

	limit := 100
	for limit > 0 {
		station := s.items[bucket]
		if station == nil || bytes.Equal(station.Name, name) {
			break
		}
		limit--
		bucket++
	}

	if limit == 0 {
		return fmt.Errorf("collision limit reached for %s", name)
	}

	if s.items[bucket] == nil {
		s.items[bucket] = &Station{
			Name: bytes.Clone(name),
			Min:  minTemp,
			Max:  maxTemp,
			Sum:  sumTemp,
			Num:  num,
		}
		s.len++
	} else {
		station := s.items[bucket]
		station.Min = min(station.Min, minTemp)
		station.Max = max(station.Max, maxTemp)
		station.Sum += sumTemp
		station.Num += num
	}
	return nil
}

func (s *Map) Get(name []byte) (*Station, error) {
	hash := fnv.New64()
	hash.Write(name)
	bucket := hash.Sum64() & (BucketSize - 1)

	limit := 100
	// fmt.Printf("%s = %v\n", name, s.items[bucket])
	for station := s.items[bucket]; station != nil && limit > 0; {
		if bytes.Equal(station.Name, name) {
			return station, nil
		}
		limit--
		bucket++
	}

	return nil, fmt.Errorf("Station not found: %s", name)
}

func (s *Map) Sorted() iter.Seq[*Station] {
	stations := make([]*Station, 0, s.len)
	for _, station := range s.items {
		if station == nil {
			continue
		}
		stations = append(stations, station)
	}
	slices.SortFunc(stations, func(a, b *Station) int {
		return bytes.Compare(a.Name, b.Name)
	})

	return func(yield func(*Station) bool) {
		for _, station := range stations {
			if !yield(station) {
				return
			}
		}
	}
}
