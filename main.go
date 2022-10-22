package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/", home)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func home(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h1>Home!</h1>")
}
