package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	const maxAsyncRequests = 5

	var (
		total uint64
		wg    sync.WaitGroup
	)

	maxRequests := make(chan struct{}, maxAsyncRequests)
	defer close(maxRequests)

	scanner := bufio.NewScanner(os.Stdin)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка при считывании строк: %v\n", err)
			return
		}
	}()

	for scanner.Scan() {
		wg.Add(1)
		maxRequests <- struct{}{}

		path := scanner.Text()

		go func(path string) {
			defer func() {
				wg.Done()
				<-maxRequests
			}()

			if len(path) == 0 {
				return
			}

			var (
				count uint64
				err   error
			)

			if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "https") {
				count, err = ReadBodyFromURL(path, &total)
			} else {
				count, err = ReadFile(path[1:], &total)
			}

			if err != nil {
				log.Print(err)
				return
			}

			fmt.Printf("Count for %s: %d\n", path, count)
		}(path)
	}

	wg.Wait()
	fmt.Printf("Total: %d\n", total)
}

func ReadBodyFromURL(path string, total *uint64) (count uint64, err error) {
	response, err := http.Get(path)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при выполнении запроса: %w\n", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("Ошибка при закрытии тела ответа: %v\n", err)
			return
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при чтении тела ответа: %w\n", err)
	}

	count = uint64(strings.Count(string(body), "Go"))
	atomic.AddUint64(total, count)

	return count, nil
}

func ReadFile(path string, total *uint64) (count uint64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при открытии файла: %w\n", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Ошибка при закрытии файла: %v\n", err)
			return
		}
	}()

	body, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при чтении файла: %w\n", err)
	}

	count = uint64(strings.Count(string(body), "Go"))
	atomic.AddUint64(total, count)

	return count, nil
}
