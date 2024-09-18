package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

func splitFile(file *os.File, numParts int) []part {
	fi, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	// Split the file into equal parts. In every part, we
	// will search for the last occurring newline in the last
	// 100 bytes of the part. The remaining trailing bytes will
	// be the starting bytes of the next part.
	var (
		fileSize  int64 = fi.Size()
		partSize  int64 = fileSize / int64(numParts)
		offset    int64 = 0
		searchLen int64 = 100
	)

	parts := make([]part, 0, numParts)
	buf := make([]byte, searchLen)

	for i := range numParts {
		if i == numParts-1 {
			if offset < fileSize {
				parts = append(parts, part{offset, fileSize - offset})
			}
			break
		}

		searchStart := offset + partSize - searchLen
		if searchStart < 0 {
			log.Panic("Huh?")
		}

		// Read bytes from offset to searchlen (seek + readfull)
		file.Seek(searchStart, io.SeekStart)
		if _, err := io.ReadFull(file, buf); err != nil {
			log.Panic(err)
		}

		newline := int64(bytes.LastIndexByte(buf, '\n'))
		if newline < 0 {
			log.Panic("Huh?")
		}

		// Next part starts after newline. Size of the current part
		// should include the newline but obviously not the char after.
		newOffset := searchStart + newline + 1
		parts = append(parts, part{offset, newOffset - offset - 1})
		offset = newOffset
	}

	return parts
}
