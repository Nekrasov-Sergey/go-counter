package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// Reader Реализующие этот интерфейс типы должны содержать метод ReadData, который выполняет
// чтение данных и возвращает количество слов "Go" и ошибку.
type Reader interface {
	ReadData() (count uint64, err error)
}

// URLReader Содержит путь до сайта и текущее суммарное кол-во слов "Go"
type URLReader struct {
	Path  string
	Total *uint64
}

// FileReader Содержит путь до файла и текущее суммарное кол-во слов "Go"
type FileReader struct {
	Path  string
	Total *uint64
}

// ReadData Читает тело ответа сайта и считает кол-во слов "Go"
func (r *URLReader) ReadData() (count uint64, err error) {
	response, err := http.Get(r.Path)
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
	atomic.AddUint64(r.Total, count)

	return count, nil
}

// ReadData Читает файл и считает кол-во слов "Go"
func (r *FileReader) ReadData() (count uint64, err error) {
	absolutePath, err := filepath.Abs(r.Path)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при получении абсолютного пути: %w\n", err)
	}

	file, err := os.Open(absolutePath)
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
	atomic.AddUint64(r.Total, count)

	return count, nil
}
