package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		path := scanner.Text()

		maxRequests <- struct{}{}
		wg.Add(1)
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

			var reader Reader

			switch {
			case CheckURL(path):
				reader = &URLReader{
					Path:  path,
					Total: &total,
				}
			case CheckFile(path[1:]):
				reader = &FileReader{
					Path:  path[1:],
					Total: &total,
				}
			default:
				log.Printf("Неверный путь %s", path)
				return
			}

			count, err = reader.ReadData()
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

// CheckURL Проверяет сайт на существование
func CheckURL(path string) bool {
	if !strings.HasPrefix(path, "http") {
		return false
	}

	head, err := http.Head(path)
	if err != nil {
		return false
	}

	if head.StatusCode != http.StatusOK {
		return false
	}

	return true
}

// CheckFile Проверяет файл на существование
func CheckFile(path string) bool {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return false
	}

	if fileInfo.Size() == 0 || fileInfo.IsDir() {
		return false
	}

	return true
}
