package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var UrlFilePath = "ListOfAsciiSiteUrls.txt"
var DownloadFileSuffix = "html_"
var WritingFileDir = "/opt/htmls/"

func writeToPath(htmlData io.ReadCloser, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Println("Failed creating path:", path)
	} else {
		defer f.Close()
		_, _ = io.Copy(f, htmlData)
	}
}

func urlBuffer(urls chan []string) {
	urlsFile, err := os.Open(UrlFilePath)
	if err != nil {
		log.Fatalf("Failed reading file: %s", err)
	} else {
		defer urlsFile.Close()
		scanner := bufio.NewScanner(urlsFile)
		index := 0
		for scanner.Scan() {
			index++
			uniqueFileSuffix := DownloadFileSuffix + strconv.Itoa(index)
			urls <- []string{uniqueFileSuffix, scanner.Text()}
		}
		close(urls)
	}
}

func downloadHtml(urls chan []string) {
	defer wg.Done()
	for urlData := range urls {
		fileSuffix := urlData[0]
		url := urlData[1]
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Failed downloading:", url)
		} else {
			defer resp.Body.Close()
			writingPath := WritingFileDir + fileSuffix
			fmt.Println("Writing ", url, "to ", writingPath)
			writeToPath(resp.Body, writingPath)
		}
	}
}

func main() {
	var numberOfWorkers int
	urls := make(chan []string)
	fmt.Println("How many workers would you like to use?")
	_, err := fmt.Scanf("%d", &numberOfWorkers)
	if err != nil {
		log.Fatalf("failed reading from stdin, %s", err)
	}

	startTime := time.Now()
	for worker := 0; worker < numberOfWorkers; worker++ {
		wg.Add(1)
		go downloadHtml(urls)
	}

	go urlBuffer(urls)
	wg.Wait()
	elapsed := time.Since(startTime)
	log.Printf("reading urls took %s", elapsed)
}
