package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

// File Структура для работы с файлом
type File struct {
}

// Check Проверяет файл на существование
func (r *File) Check(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileInfo.Size() == 0 || fileInfo.IsDir() {
		return false
	}

	return true
}

// Read Читает файл и считает кол-во слов "Go"
func (r *File) Read(path string, total *uint64) (count uint64, err error) {
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
