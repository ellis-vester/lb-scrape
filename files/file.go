package files

import (
	"bufio"
	"log"
	"os"
)

func GetFilmListUrls(path string) (urls []string, err error) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to open file", err)
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	listUrls := []string{}

	for scanner.Scan() {
		listUrls = append(listUrls, scanner.Text())
	}

	return listUrls, nil
}
