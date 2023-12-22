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

type problem struct {
	question string
	answer   string
}

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

	//create a new Reader from the file then read all lines from the file
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	//parse the lines read from the CSV into a slice of problem structs
	problems := parseLines(records)

	//if the flag is true shuffle the list of questions
	if *shuffle {
		for i := range problems {
			ix := rand.Intn(i + 1)
			problems[i], problems[ix] = problems[ix], problems[i]
		}
	}

	//initialize channel for user answers and counter for correct answers
	answerCh := make(chan string)
	correct := 0

	//initialize timer to run during the course of the quiz
	clock := time.NewTimer(time.Duration(*limit) * time.Second)

questionLoop:

	//for each quiz question, print it and try to get a answer from user or what for timer to run out
	for i, entry := range problems {

		//print the question to the user, triming leading and trailing spaces
		fmt.Printf("Problem #%d: %s = ", i+1, entry.question)

		//in a goroutine read input from user and send through answer channel
		go func() {
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				log.Fatal(err)
			}

			answerCh <- strings.ToLower(strings.TrimSpace(answer))
		}()

		//handle timer or user input, whichever comes in first
		select {

		//CASE: timer ended first
		//EFFECT: break out of loop asking driving the quiz
		case <-clock.C:
			break questionLoop

		//CASE: user answered before timer ended
		//EFFECT: trim leading/trailing spaces from user input and increment counter if it matches correct answer
		case answer := <-answerCh:
			if answer == entry.answer {
				correct++
			}
		}
	}

	//inform user of their score before terminating
	fmt.Printf("\nYou scored %d out of %d total.\n", correct, len(records))
}

func parseLines(records [][]string) []problem {
	//initialize a slice to store the new structs
	problems := make([]problem, len(records))

	//for each line read in from the CSV, turn trim the spaces from its question and answer and create a problem struct of it
	for i, record := range records {
		problems[i] = problem{question: strings.TrimSpace(record[0]), answer: strings.TrimSpace(record[1])}
	}

	//return the slice of problems to the user
	return problems
}
