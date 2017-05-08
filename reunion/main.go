package main

import "fmt"
import "io/ioutil"
import "log"
import "bytes"
import "os"
import "math"

func main() {
	
	fmt.Println("")
	fileInfos, err := ioutil.ReadDir("../reunion_pics")
	if err != nil {
		log.Fatal(err)
	}

	page := 1
	filesOnPage := []string{}

	perPage := 5
	header := makeHeader(len(fileInfos), perPage)
	var i int
	var fileInfo os.FileInfo
	for i, fileInfo = range fileInfos {	
		filesOnPage = append(filesOnPage, fileInfo.Name())
		if i % perPage == 0 {
			writePage(page, filesOnPage, header)
			page += 1
			filesOnPage = nil
		}
	}

	// make a page of the leftovers
	if i % 5 != 0 {
		writePage(page, filesOnPage, header)
	}

}

func makeHeader(fileLen, perPage int) ([]byte) {
	var b bytes.Buffer
	pagesFloat := float64(fileLen) / float64(perPage)
	pages := int(math.Ceil(pagesFloat))

	for i := 0; i < pages; i++ {
		startNumber := pages * perPage + 1
		endNumber := pages * (perPage + 1)
		b.WriteString(fmt.Sprintf(`<a href="%d.html">%d-%d</a>`, startNumber, endNumber))	
	}

	return b.Bytes()
}

func writePage(page int, files []string, header []byte) {	
	f, err := os.Create(fmt.Sprintf("%d.html", page))
	if err != nil {
		log.Fatal(err)	
	}
	defer f.Close()
	f.Write(header)
	f.WriteString("\n")
	for _, file := range files {
		f.WriteString(`<img src="` + file + `" />` )	
	}
	
}
