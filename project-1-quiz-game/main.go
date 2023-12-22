package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	filename := flag.String("filename", "problems.csv", "Name of the csv file to use as question bank.")
	limit := flag.Int("limit", 30, "Number of seconds to allow for participants to answer questions.")
	shuffle := flag.Bool("shuffle", true, "Whether to shuffle the questions each run.")
	flag.Parse()

	//open up the file
	file, err := os.Open(*filename)
	if err != nil {
		//Fatal is equivalent to calling Print following by os.Exit(1)
		log.Fatal(err)
	}

	//defer schedules a function call to be run immediately before the function executing the defer returns
	defer file.Close()

	//create a new Reader from the file then do one priming read from it
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	if *shuffle {
		for i := range records {
			ix := rand.Intn(i + 1)
			records[i], records[ix] = records[ix], records[i]
		}
	}

	//initialize channel for user answers and int to count correct answers
	answerCh := make(chan string)
	correct := 0

	clock := time.NewTimer(time.Duration(*limit) * time.Second)
questionLoop:
	for i, record := range records {
		//increment number of questions asked and ask the question
		fmt.Printf("Problem #%d: %s = ", i+1, strings.TrimSpace(record[0]))

		go func() {
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				log.Fatal(err)
			}

			answerCh <- strings.ToLower(strings.TrimSpace(answer))
		}()

		select {
		case <-clock.C:
			break questionLoop
		case answer := <-answerCh:
			if strings.ToLower(strings.TrimSpace(answer)) == record[1] {
				correct++
			}
		}
	}

	//Base case of answering all the questions before time runs out
	fmt.Printf("\nYou scored %d out of %d total.\n", correct, len(records))
}
