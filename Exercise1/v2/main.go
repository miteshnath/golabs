package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	filePath string
	wg       sync.WaitGroup
)

func init() {
	flag.StringVar(&filePath, "filePath", "questions.csv", "path/to/cs_file")
	flag.Parse()
}

func main() {
	csvPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalln("Unable to parse path" + csvPath)
	}
	fmt.Println(csvPath)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	timeOver := time.After(10 * time.Second)
	totalQuestions := len(csvData)
	totalCorrect := 0
	questions := make(map[int]string, totalQuestions)
	correctAnswers := make(map[int]string, totalQuestions)
	userResponses := make(map[int]string, totalQuestions)

	responseChan := make(chan string)

	for i, data := range csvData {
		questions[i] = data[0]
		correctAnswers[i] = data[1]
	}

	wg.Add(1)
	go func() {
	iterate:
		for i := 0; i < totalQuestions; i++ {
			go askQuestion(os.Stdout, os.Stdin, questions[i], responseChan)
			select {
			case <-timeOver:
				fmt.Fprintln(os.Stderr, "Time Over!!")
				break iterate
			case ans, ok := <-responseChan:
				if ok {
					userResponses[i] = ans
				} else {
					break iterate
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	for i := 0; i < len(userResponses); i++ {
		if verfityAnswer(correctAnswers[i], userResponses[i]) {
			totalCorrect++
		}
	}

	fmt.Printf("totalCorrect: %d", totalCorrect)
}

func askQuestion(w io.Writer, r io.Reader, question string, responseChan chan string) {
	reader := bufio.NewReader(r)
	fmt.Fprintln(w, "Question: "+question)
	fmt.Fprintln(w, "Answer: ")
	resp, err := reader.ReadString('\n')
	if err == io.EOF {
		close(responseChan)
		if err == io.EOF {
			return
		}
		log.Fatalln(err)
	}
	responseChan <- strings.TrimRight(resp, "\n")
}

func verfityAnswer(correctAnswer, userResponse string) bool {
	if correctAnswer == userResponse {
		return true
	}
	return false
}

