package main

// ReaderChecker Реализующие этот интерфейс типы должны содержать методы Check и Read. Check проверяет путь
// на существование, а Read выполняет чтение данных и возвращает количество слов "Go" и ошибку.
type ReaderChecker interface {
	Check() bool
	Read() (count uint64, err error)
}
