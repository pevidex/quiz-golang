package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
Simple exercise to learn GO.
Disclaimers:
- No tests were added since this is only a way to get in touch with go language
- Neither performance or memory concerns were taken into consideration even though we are dealing with external files
- Some error handling scenarios were ignored
*/

type MathQuestion struct {
	question string
	solution string
}

func buildQuizObj(csvLines [][]string) []MathQuestion {
	var questions []MathQuestion
	for i := 0; i < len(csvLines); i++ {
		questions = append(questions, MathQuestion{question: csvLines[i][0], solution: csvLines[i][1]})
	}
	return questions
}

func readCSVFile(csvFilePath string) ([][]string, error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	// memory concerns here
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	return csvLines, nil
}

func play(questions []MathQuestion, timeInt int) {
	fmt.Print("Press enter to start...")
	fmt.Scanln()
	c1 := make(chan string, 1)
	nrCorrect := 0
	go func() {
		for i := 0; i < len(questions); i++ {
			fmt.Print(questions[i].question, " ")
			var answer string
			fmt.Scanf("%s\n", &answer)
			if answer == questions[i].solution {
				nrCorrect++
			}
		}
		c1 <- "end"
	}()
	select {
	case res := <-c1:
		_ = res
		fmt.Print("Number of questions: ", len(questions), ", Correct answers: ", nrCorrect)
	case <-time.After(time.Duration(timeInt) * time.Second):
		fmt.Print("Timeout! Number of questions: ", len(questions), ", Correct answers: ", nrCorrect)
		return
	}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file with the format 'question, answer'")
	time := flag.Int("time", 30, "set a timer for the quiz")
	flag.Parse()

	fmt.Println("Quiz game by me!")
	csvFilePath := filepath.Join("..", "quizzes", *csvFilename)

	csvLines, err := readCSVFile(csvFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	questions := buildQuizObj(csvLines)
	play(questions, int(*time))
}
