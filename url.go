package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

// URL Содержит путь до сайта и текущее суммарное кол-во слов "Go"
type URL struct {
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
