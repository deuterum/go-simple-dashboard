package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, getStatsJson())
	})

	http.HandleFunc("/procs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, getProcs())
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, getHostinfo())
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Сервер запущен: http://localhost:8080/dashboard")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Ошибка запуска сервера: ", err)
	}
}
