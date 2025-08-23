package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"

	"github.com/loffy1337/go_todo/internal/config"
	httprouter "github.com/loffy1337/go_todo/internal/transport/http"
	"github.com/loffy1337/go_todo/pkg/logger"
)

func main() {
	// Загрузка конфигурации проекта
	config, err := config.LoadConfig()
	if err != nil {
		logger.Error.Fatalf("config load: %v", err)
	}

	fmt.Println(config.MySQL.Port)

	// Подключение к базе данных
	db, err := sql.Open("mysql", config.GetMySQLDSN())
	if err != nil {
		logger.Error.Fatalf("database open: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Error.Fatalf("database ping: %v", err)
	}

	// Получение http-роутера
	var router *chi.Mux = httprouter.NewRouter()

	// Создание и объединение с дефолтным http-обработчиком
	var server *http.Server = &http.Server{
		Addr:              ":" + config.AppPort,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Запуск HTTP-сервера
	logger.Info.Printf("%s listening on :%s", config.AppName, config.AppPort)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error.Fatalf("server: %v", err)
	}
}
