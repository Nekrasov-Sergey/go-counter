# Тестовое задание "go-counter"

## Для запуска прописать команды из директории проекта
* go build go-counter.go
* echo -e 'https://golang.org\n/etc/passwd\nhttps://golang.org\nhttps://golang.org' | ./go-counter

## Задание
Процессу на stdin приходят строки, содержащие URL или названия файлов. Каждый такой URL нужно запросить,
каждый файл нужно прочитать, и посчитать кол-во вхождений строки "Go" в ответе.
В конце работы приложение выводит на экран общее кол-во найденных строк "Go" во всех источниках данных, например:
$ echo -e 'https://golang.org\n/etc/passwd\nhttps://golang.org\nhttps://golang.org' | ./go-counter
Count for https://golang.org: 9
Count for /etc/passwd: 0
Count for https://golang.org: 9
Count for https://golang.org: 9
Total: 27

Каждый источник данных должен начать обрабатываться сразу после вычитывания и параллельно с вычитыванием следующего.
Источники должны обрабатываться параллельно,но не более k=5 одновременно.
Обработчики данных не должны порождать лишних горутин,
т.е. если k=1000 а обрабатываемых источников нет, не должно создаваться 1000 горутин.
Нужно обойтись без глобальных переменных и использовать только стандартные библиотеки.
Представьте, что этой тулзой в будущем будут пользоваться и поддерживать ваши коллеги,
поэтому писать стоит так, чтоб вам самим такой код было не больно использовать и дополнять.