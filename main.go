package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func main() {
	defer elapsed("main")()
	arg := os.Args[1]
	numbers, err := readFile(arg)
	if err != nil {
		log.Fatalf("Error while read from file: %s", err)
	}
	sort.Ints(numbers)

	reader := bufio.NewReader(os.Stdin)
	c0 := make(chan int)
	c1 := make(chan int)
	go readStdIn(reader, c0)
	go indexOf(numbers, c0, c1)
	printResult(c1)
}

func readStdIn(reader *bufio.Reader, downstream chan int) {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while read from stdin: %s", err)
		}
		number, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			log.Fatalf("Error while read from stdin: %s", err)
		}
		downstream <- number
	}
	close(downstream)
}

func indexOf(numbers []int, upstream, downstream chan int) {
	for key := range upstream {
		lo := 0
		hi := len(numbers) - 1
		notFound := true
		for lo <= hi {
			mid := lo + (hi-lo)/2
			if key < numbers[mid] {
				hi = mid - 1
			} else if key > numbers[mid] {
				lo = mid + 1
			} else {
				notFound = false
				break
			}
		}
		if notFound {
			downstream <- key
		}
	}
	close(downstream)
}

func printResult(upstream chan int) {
	count := 0
	for key := range upstream {
		count++
		fmt.Println(key)
	}
	fmt.Println(count)
}

func readFile(path string) ([]int, error) {
	defer elapsed("readFile")()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	return numbers, scanner.Err()
}
