package scraper

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	lb "github.com/ellis-vester/lb-scrape/letterboxd"
	"github.com/gocolly/colly"
)

func ScrapeFilmListHtml(url string) (string, error) {

	collector := colly.NewCollector()

	var err error
	var html string

	collector.OnHTML("ul.poster-list", func(e *colly.HTMLElement) {
		html, err = e.DOM.Html()
	})

	err = collector.Visit(url)
	if err != nil {
		return html, err
	}

	return html, err
}

func ParseFilmList(content string) ([]lb.FilmListEntry, error) {

	listEntries := []lb.FilmListEntry{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return listEntries, errors.New("error creating film list reader")
	}

	success := true
	errorMessage := ""

	doc.Find("li.poster-container").Each(func(i int, selection *goquery.Selection) {

		listEntry := lb.FilmListEntry{}

		rating, exists := selection.Attr("data-owner-rating")
		if !exists || rating == "" {
			success = false
			errorMessage = "error parsing data-owner-rating from film list"
			return
		}

		ratingInt, err := strconv.ParseInt(rating, 10, 64)
		if err != nil {
			success = false
			errorMessage = "error parsing rating from film list"
			return
		}

		listEntry.Rating = int8(ratingInt)

		link, exists := selection.
			Find("div.film-poster").
			Attr("data-target-link")
		if !exists || link == "" {
			success = false
			errorMessage = "error parsing data-target-link from film list" + link
			return
		}

		listEntry.Link = link
		listEntries = append(listEntries, listEntry)
	})

	if !success {
		return nil, errors.New(errorMessage)
	}

	return listEntries, nil
}

func ParseUsername(url string) string {
	strings.Split(url, "/")
	return strings.Split(url, "/")[3]
}

func ScrapeFilmHtml(url string) (string, error) {

	collector := colly.NewCollector()

	var err error
	var html string

	collector.OnHTML("section.film-header-group", func(e *colly.HTMLElement) {
		html, err = e.DOM.Html()
	})

	err = collector.Visit(url)
	if err != nil {
		return html, err
	}

	return html, err
}

func ParseFilm(content string) (lb.Film, error) {

	film := lb.Film{}

	success := true
	errorMessage := ""

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return film, errors.New("error creating film reader")
	}

	titleSel := doc.Find("h1.filmtitle").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("span").Text()
		if title == "" {
			success = false
			errorMessage = "error parsing title from film"
			return
		}
		film.Title = title
	})
	if titleSel.Length() == 0 {
		success = false
		errorMessage = "error parsing title from film"
	}

	yearSel := doc.Find("div.releaseyear").Each(func(i int, selection *goquery.Selection) {
		yearText := selection.Find("a").Text()
		if yearText == "" {
			success = false
			errorMessage = "error parsing year from film"
			return
		}

		year, err := strconv.ParseInt(yearText, 10, 64)
		if err != nil {
			success = false
			errorMessage = "error parsing year from film"
			return
		}
		film.Year = int(year)
	})
	if yearSel.Length() == 0 {
		success = false
		errorMessage = "error parsing year from film"
	}

	directorSel := doc.Find("a.contributor").Each(func(i int, selection *goquery.Selection) {
		director := selection.Find("span").Text()
		if director == "" {
			success = false
			errorMessage = "error parsing director from film"
			return
		}
		film.Director = director
	})
	if directorSel.Length() == 0 {
		success = false
		errorMessage = "error parsing director from film"
	}

	if !success {
		return film, errors.New(errorMessage)
	}

	return film, nil
}

func SumFilmInclusions(lists [][]lb.FilmListEntry) []lb.FilmListEntry {

	var films = map[string]*lb.FilmListEntry{}

	for _, list := range lists {
		for _, listItem := range list {

			var inclusions int

			_, exists := films[listItem.Link]
			if exists {
				inclusions = films[listItem.Link].Inclusions + 1
				films[listItem.Link] = &listItem
				films[listItem.Link].Inclusions = inclusions
			} else {
				films[listItem.Link] = &listItem
				films[listItem.Link].Inclusions = 1
			}
		}
	}

	filmListEntries := []lb.FilmListEntry{}

	for _, value := range films {
		filmListEntries = append(filmListEntries, *value)
	}

	sort.Slice(filmListEntries, func(i, j int) bool {
		return filmListEntries[i].Inclusions > filmListEntries[j].Inclusions
	})

	return filmListEntries
}

func SumDirectorInclusions(list []lb.Film) []lb.Director {
	var directorsMap = map[string]int{}

	for _, listItem := range list {
		directorsMap[listItem.Director]++
	}

	var directors = []lb.Director{}

	for key, value := range directorsMap {
		directors = append(directors, lb.Director{
			Name:       key,
			Inclusions: value,
		})
	}

	sort.Slice(directors, func(i, j int) bool {
		return directors[i].Inclusions > directors[j].Inclusions
	})

	return directors
}
