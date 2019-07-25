package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readInput(answer chan string) {
	var choice string
	for {
		// Read entire line into buffer to avoid reading character by character
		_, err := fmt.Scanln(&choice)
		if err != nil {
			log.Println(err)
			continue
		}
		answer <- choice
	}
}

func askQuestions(problems [][]string, duration time.Duration) int {

	answer := make(chan string, 1)
	go readInput(answer)

	timeout := time.NewTimer(duration)
	goodAnswers := 0

	for _, problem := range problems {
		fmt.Printf("Question: %s : ", problem[0])

		select {
		case choice := <-answer:
			if strings.TrimSpace(choice) == problem[1] {
				goodAnswers++
			}
		case <-timeout.C:
			return goodAnswers
		}
	}

	return goodAnswers
}

func main() {

	problemsPath := flag.String("problems", "problems.csv", "Path to problems CSV file")
	quizDurationSecond := flag.Int("timeout", 30, "Maximum duration of the session")
	suffle := flag.Bool("shuffle", false, "Randomize quiz order")
	flag.Parse()

	quizDuration := time.Duration(*quizDurationSecond) * time.Second

	problemsFile, err := os.Open(*problemsPath)
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(problemsFile)

	problems, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if *suffle {
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed rng
		tProblems := make([][]string, len(problems))

		copy(tProblems, problems)
		for i, j := range r.Perm(len(problems)) {
			problems[i] = tProblems[j]
		}
	}

	fmt.Println("You will have", quizDuration, "seconds to answer as many questions as you can")
	fmt.Println("Press enter key to start")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	goodAnswers := askQuestions(problems, quizDuration)

	totalQuestions := len(problems)
	fmt.Println("\nYour score:", goodAnswers, "/", totalQuestions)
}
