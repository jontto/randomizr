package main

import (
	"bufio"
	crypto_rand "crypto/rand"
	"encoding/binary"
	tm "github.com/buger/goterm"
	"math/rand"
	"os"
	"strconv"
)

type participant struct {
	name  string
	count int
}

const DEFAULT_FILENAME = "entries.txt"

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	return fileLines, err
}

func main() {

	// read file
	entriesFromFile, ferr := readFile(DEFAULT_FILENAME)
	if ferr != nil {
		panic(ferr)
	}

	// map file lines to participants
	var participants []participant
	for _, entry := range entriesFromFile {
		part := participant{
			name: entry,
		}
		participants = append(participants, part)
	}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("Failed to create truly random seed for the random number picker!")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	rand.Shuffle(len(participants), func(i, j int) { participants[i], participants[j] = participants[j], participants[i] })

	previousN := -1
	round := 0

	tm.Clear()
	tm.MoveCursor(1,1)
	tm.Flush()

	for {
		if previousN < 0 || participants[previousN].count <= 3 {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		tm.Clear()
		tm.MoveCursor(1, 1)

		round++
		tm.Println("Round " + strconv.Itoa(round))
		n := rand.Int() % len(participants)
		if previousN >= 0 && participants[previousN] == participants[n] {
			tm.Println("Ignoring second consecutive for same person.")
		} else {
			participants[n].count++
			tm.Println()
		}
		previousN = n
		tm.Println("Random list order")
		for _, part := range participants {
			if part.count >= 3 {
				tm.Printf(tm.Color(part.name+"\t"+strconv.Itoa(part.count), tm.RED))
			} else {
				tm.Printf(part.name+"\t%d", part.count)
			}
			if part == participants[previousN] {
				tm.Print(" <---")
			}
			tm.Println()
			tm.Flush()
		}
	}
}
