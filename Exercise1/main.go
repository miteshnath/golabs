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
)

var (
	filePath string
)

func init() {
	flag.StringVar(&filePath, "filePath", "questions.csv", "path/to/cs_file")
	flag.Parse()
}

func main() {
	// This is quiz exercise of gopherisces.
	// No randomization of questions , No timeout
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

	var totalQuest = len(csvData)
	questions := make(map[int]string, totalQuest)
	answers := make(map[int]string, totalQuest)
	//responses := make(map[int]string, totalQuest)

	for i, data := range csvData {
		questions[i] = data[0]
		answers[i] = data[1]
	}

	totalCorrect := 0

	for key, value := range questions {
		num := askQuestion(os.Stdout, os.Stdin, value, answers[key])
		if num == 1 {
			totalCorrect = totalCorrect + 1
		}
	}
	fmt.Printf("totalCorrect: %d", totalCorrect)
}

func askQuestion(w io.Writer, r io.Reader, question string, answer string) int {
	reader := bufio.NewReader(r)
	fmt.Fprintln(w, "Question: "+question)
	fmt.Fprintln(w, "Answer: ")
	resp, err := reader.ReadString('\n')
	if err == io.EOF {
		if err == io.EOF {
			return 0
		}
		log.Fatalln(err)
	}
	return checkCorrect(answer, strings.TrimRight(resp, "\n"))
}

func checkCorrect(answer, resp string) int {
	//fmt.Println("answer: ", answer)
	fmt.Println("resp: ", resp)
	if answer == resp {
		return 1
	}
	return 0
}

