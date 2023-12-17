package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

// URL Структура для работы с сайтом
type URL struct {
}

// Check Проверяет сайт на существование
func (r *URL) Check(path string) bool {
	if !strings.HasPrefix(path, "http") {
		return false
	}

	return true
}

// Read Читает тело ответа сайта и считает кол-во слов "Go"
func (r *URL) Read(path string, total *uint64) (count uint64, err error) {
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
