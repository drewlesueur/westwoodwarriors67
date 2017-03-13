package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)
import "io/ioutil"
import "log"

type Page struct {
	Letters       []string
	LettersString string
	Number        int
	URL           string
	LocalFile     string
	SubNumber     int
}

func (p Page) FileName() string {
	if p.SubNumber == 0 {
		return p.LettersString + ".html"
	}

	return fmt.Sprintf("%s-%d.html", p.LettersString, p.SubNumber + 1)
}

func (p Page) ImageName() string {
	if p.SubNumber == 0 {
		return p.LettersString + ".jpg"
	}

	return fmt.Sprintf("%s-%d.jpg", p.LettersString, p.SubNumber + 1)
}

func (p Page) Title() string {
	if p.SubNumber == 0 {
		return p.LettersString
	}

	return fmt.Sprintf("%s (%d)", p.LettersString, p.SubNumber + 1)
}

func main() {
	b, err := ioutil.ReadFile("./yearbook.txt")
	if err != nil {
		log.Fatal(err)
	}

	t := string(b)
	lines := strings.Split(t, "\n")

	lines = lines[2:]

	pages := []Page{}

	setI := 0

	var p Page
	for _, line := range lines {
		switch setI {
		case 0:
			number, err := strconv.Atoi(line)
			if err != nil {
				log.Fatalf("could not convert page number: %v", err)
			}

			p = Page{Number: number}
		case 1:
			fmt.Println("letters are", line)
			letters := parseLetters(line)
			fmt.Println(letters)
			p.Letters = letters
			p.LettersString = line

			// add the SubNumber (like if there are 2 pages of B's)
			// NOTE: same letter have to be in order for now
			if len(pages) > 0 {
				lastPage := pages[len(pages)-1]
				if len(lastPage.Letters) == 1 && len(p.Letters) == 1 && lastPage.Letters[0] == p.Letters[0] {
					p.SubNumber = lastPage.SubNumber + 1
				}
			}
		case 2:
			p.URL = line

			// you don't need to run this every time
			//_, err := exec.Command("wget", "-O", p.ImageName(), p.URL).Output()
			//if err != nil {
			//	log.Fatalf("error wgetting the file: %v", err)
			//}
		case 3:
			setI = -1
			pages = append(pages, p)
		}
		setI += 1
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Number <= pages[j].Number
	})

	for _, p := range pages {
		fmt.Printf("%+v\n", p)
	}

	footer := `<div class="index">
<a href="../">Home</a>
`
	for _, p := range pages {
		if (p.SubNumber > 0) {

		}
		footer += `<a href="` + p.FileName() + `">` + p.Title() + `</a> `
	}

	footer += `</div>`

	for _, p := range pages {
		f, err := os.Create(p.FileName())

		if err != nil {
			log.Fatal("error creating file")
		}

		_, err = f.WriteString("<!doctype html>")
		if err != nil {
			log.Fatal("error writing file")
		}

		_, err = f.WriteString(footer)
		if err != nil {
			log.Fatal("error writing footer")
		}

		_, err = f.WriteString("<h1>" + p.Title() + "</h1>")
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(`<div><img src="`+p.ImageName()+`" /></div>`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(footer)
		if err != nil {
			log.Fatal("error writing footer")
		}

		f.Close()
	}

	_, err = exec.Command("cp", pages[0].LettersString+".html", "index.html").Output()
	if err != nil {
		log.Fatalf("error making index file: %v", err)
	}

}

//parseLetters takes a string like "A-F" and converts it to a slice of strings {"A", "B", "C", "D", "E", "F"}
func parseLetters(line string) []string {
	parts := strings.Split(line, "-")
	if len(parts) < 2 {
		return []string{}
	}

	first := []byte(parts[0])
	if len(first) == 0 {
		return []string{}
	}

	second := []byte(parts[1])
	if len(second) == 0 {
		return []string{}
	}

	ret := []string{}
	start := int(first[0])
	end := int(second[0])

	for i := start; i <= end; i++ {
		ret = append(ret, string(i))
	}
	return ret
}
