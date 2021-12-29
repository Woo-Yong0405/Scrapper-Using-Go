package main

import (
	"os"
	"strings"

	"github.com/Woo-Yong0405/google-scraper/scraper"
	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	term := strings.ToLower(scraper.CleanString(c.FormValue("term")))
	var fileName string = term + "_jobs.csv"
	defer os.Remove(fileName)
	scraper.Scrape(term)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1324"))
}
