package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var rdb *redis.Client

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://robin.r:@localhost:5432/people?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer rdb.Close()
	_, err = rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatal(err)
	}

	// Create the Gin router and enable CORS
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Define the routes
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := db.Exec("INSERT INTO public.user (name) VALUES ($1)", user.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Message": "Successfully added to the DB"})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		cacheKey := "user-" + id
		val, err := rdb.Get(c, cacheKey).Result()
		if err == nil {
			err = json.Unmarshal([]byte(val), &user)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, gin.H{"data": user})
			return
		} else if err != redis.Nil {
			log.Fatal(err)
		}
		row := db.QueryRow("SELECT id, name FROM public.user WHERE id = $1", id)
		if err := row.Scan(&user.ID, &user.Name); err != nil {
			if err == sql.ErrNoRows {
				c.Status(http.StatusNotFound)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		data, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		err = rdb.Set(c, cacheKey, data, 100*time.Second).Err()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, user)
	})

	router.GET("/users", func(c *gin.Context) {
		var users []User
		cacheKey := "users"
		val, err := rdb.Get(c, cacheKey).Result()
		if err == nil {
			err = json.Unmarshal([]byte(val), &users)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, gin.H{"data": users})
			return
		} else if err != redis.Nil {
			log.Fatal(err)
		}
		rows, err := db.Query("SELECT id, name FROM public.user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data, err := json.Marshal(users)
		if err != nil {
			log.Fatal(err)
		}
		err = rdb.Set(c, cacheKey, data, 10*time.Second).Err()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, users)
	})

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// const webport = "80"

// type Config struct{}

// func main() {
// 	app := Config{}
// 	log.Printf("Starting broker server on port %s\n", webport)
// 	//define server
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webport),
// 		Handler: app.routes(),
// 	}
// 	//start server
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }