package main

import (
	"flag"
	"log"

	"example.com/event-booker/db"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	db.InitDB()
	r := setupEngine()
	r.Run(":8080")
}

func setupEngine() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.Recovery())
	r.Use(middlewares.ErrorHandler())

	routes.RegisterRoutes(r)
	r.Use(middlewares.NotFoundHandler())

	return r
}

func loadEnv() {
	env := flag.String("env", "dev", "Environment (dev|prod)")
	flag.Parse()

	var envFile string
	if *env == "prod" {
		envFile = ".env.prod"
		gin.SetMode(gin.ReleaseMode)
	} else {
		envFile = ".env.dev"
		gin.SetMode(gin.DebugMode)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}
