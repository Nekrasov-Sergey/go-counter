package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// File Содержит путь до файла и текущее суммарное кол-во слов "Go"
type File struct {
	Path  string
	Total *uint64
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
