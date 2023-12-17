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

// ReaderChecker Реализующие этот интерфейс типы должны содержать методы Check и Read. Check проверяет путь
// на существование, а Read выполняет чтение данных и возвращает количество слов "Go" и ошибку.
type ReaderChecker interface {
	Check() bool
	Read() (count uint64, err error)
}

// URL Содержит путь до сайта и текущее суммарное кол-во слов "Go"
type URL struct {
	Path  string
	Total *uint64
}

// File Содержит путь до файла и текущее суммарное кол-во слов "Go"
type File struct {
	Path  string
	Total *uint64
}

// Check Проверяет сайт на существование
func (r *URL) Check() bool {
	if !strings.HasPrefix(r.Path, "http") {
		return false
	}

	head, err := http.Head(r.Path)
	if err != nil {
		return false
	}

	if head.StatusCode != http.StatusOK {
		return false
	}

	return true
}

// Read Читает тело ответа сайта и считает кол-во слов "Go"
func (r *URL) Read() (count uint64, err error) {
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

// Check Проверяет файл на существование
func (r *File) Check() bool {
	absolutePath, err := filepath.Abs(r.Path)
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

	r.Path = absolutePath

	return true
}

// Read Читает файл и считает кол-во слов "Go"
func (r *File) Read() (count uint64, err error) {
	file, err := os.Open(r.Path)
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
