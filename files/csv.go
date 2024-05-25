package files

import (
	"encoding/csv"
	"os"
	"strconv"

	lb "github.com/ellis-vester/lb-scrape/letterboxd"
)

func WriteFilmsToCsv(films []lb.Film, path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()
	}()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{
		"Title",
		"Director",
		"Year",
		"Rating",
		"Inclusions",
		"Link"})
	if err != nil {
		return err
	}

	for _, film := range films {
		err = writer.Write([]string{
			film.Title,
			film.Director,
			strconv.Itoa(film.Year),
			strconv.Itoa(int(film.Rating)),
			strconv.Itoa(film.Inclusions),
			film.Link})
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return err
}

func WriteDirectorsToCsv(directors []lb.Director, path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()
	}()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{"Director", "Inclusions"})
	if err != nil {
		return err
	}

	for _, director := range directors {
		err = writer.Write([]string{director.Name, strconv.Itoa(director.Inclusions)})
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return err
}
