package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
				<-maxRequests
				wg.Done()
			}()

			if len(path) == 0 {
				return
			}

			url := &URL{}
			file := &File{}
			objects := []ReaderChecker{url, file}

			var readerChecker ReaderChecker

			for _, object := range objects {
				if object.Check(path) {
					readerChecker = object
					break
				}
			}
			if readerChecker == nil {
				log.Printf("Неверный путь %s", path)
				return
			}

			count, err := readerChecker.Read(path, &total)
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
