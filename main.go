package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://postgres:Zxcvbnm123@localhost:5432/todolist?sslmode=disable"
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL: ", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("База недоступна: ", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Todo API работает",
		})
	})

	r.GET("/todos", func(c *gin.Context) {
		rows, err := db.Query(
			context.Background(),
			`SELECT id, title, completed, created_at
			 FROM todos
			 ORDER BY id`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Не удалось получить список задач",
			})
			return
		}
		defer rows.Close()

		todos := make([]Todo, 0)

		for rows.Next() {
			var todo Todo

			err := rows.Scan(
				&todo.ID,
				&todo.Title,
				&todo.Completed,
				&todo.CreatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка чтения данных",
				})
				return
			}

			todos = append(todos, todo)
		}

		c.JSON(http.StatusOK, todos)
	})

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}