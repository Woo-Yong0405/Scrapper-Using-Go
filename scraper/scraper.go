package scraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJobs struct {
	URL      string
	Title    string
	Company  string
	Location string
}

func Scrape(term string) {
	c := make(chan []extractedJobs)
	var jobs []extractedJobs
	var baseURL string = "https://stackoverflow.com/jobs?q=" + term
	totalPages := getPages(baseURL)
	for i := 0; i < totalPages; i++ {
		fmt.Println("Getting jobs from page ", i+1)
		go getPage(i+1, baseURL, c)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	writePage(jobs, term)
}

func getPage(page int, url string, c chan<- []extractedJobs) {
	var jobs []extractedJobs
	pageURL := url + "&pg=" + strconv.Itoa(page)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkPage(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchJobs := doc.Find(".-job>.d-flex>.fl1")
	searchJobs.Each(func(i int, s *goquery.Selection) {
		jobs = append(jobs, extractJobs(s))
	})
	c <- jobs
}

func extractJobs(card *goquery.Selection) extractedJobs {
	titleA, _ := card.Find(".mb4>a").Attr("title")
	title := CleanString(titleA)
	URLA, _ := card.Find(".mb4>a").Attr("href")
	URL := CleanString(URLA)
	Company := CleanString(card.Find("h3>span:first-of-type").Text())
	Location := CleanString(card.Find("h3>.fc-black-500").Text())
	return extractedJobs{
		URL:      URL,
		Title:    title,
		Company:  Company,
		Location: Location,
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkPage(res)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".s-pagination :nth-child(6)").Each(func(i int, s *goquery.Selection) {
		ddd := CleanString(s.Text())
		finalInt, err := strconv.Atoi(ddd)
		checkErr(err)
		pages = finalInt
	})
	defer res.Body.Close()
	return pages
}

func writePage(jobs []extractedJobs, term string) {
	file, err := os.Create(term + "_jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	headers := []string{"URL", "Title", "Company", "Location"}
	wErr := w.Write(headers)
	checkErr(wErr)
	for _, jobs := range jobs {
		jobSlice := []string{"https://stackoverflow.com" + jobs.URL, jobs.Title, jobs.Company, jobs.Location}
		jErr := w.Write(jobSlice)
		checkErr(jErr)
	}
	defer w.Flush()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

func checkPage(res *http.Response) {
	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode)
		log.Fatalln()
	}
}
