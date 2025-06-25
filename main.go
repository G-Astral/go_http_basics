package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request)  {
	var msg string

	switch r.URL.Path {
	case "/":
		msg = "Привет, мир!"
	case "/hello":
		msg = "Привет, пользователь!"
	default:
		http.NotFound(w, r)
		return
	}

	if msg != "" {
		fmt.Fprintln(w, msg)
	}
}

func main()  {
	// http.HandleFunc("/", handler) — регистрирует функцию-обработчик
	// handler для пути /. То есть: если кто-то откроет http://localhost:8080/ —
	// Go вызовет твою функцию handler.
	http.HandleFunc("/", handler)
	http.HandleFunc("/hello", handler)

	fmt.Println("Сервер запущен на http://localhost:8080")

	// http.ListenAndServe(":8080", nil) — запускает сервер.
	// Он слушает порт 8080 и обрабатывает все HTTP-запросы
	http.ListenAndServe(":8080", nil)
}