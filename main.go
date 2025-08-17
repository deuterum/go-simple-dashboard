package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func handleFuncs() {
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
}

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalln("Ошибка загрузки конфига: ", err)
	}

	handleFuncs()

	log.Printf("Сервер запущен: http://%s:%d/dashboard", cfg.Host, cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	if err != nil {
		log.Fatalln("Ошибка запуска сервера: ", err)
	}
}
