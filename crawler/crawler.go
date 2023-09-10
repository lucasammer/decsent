package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/mvdan/xurls"
)


func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func isAllowed(path string, disallowList []string, allowList []string) (bool) {
	for _, pattern := range disallowList {
		matchingTo := strings.Replace(path, "*", ".*", -1)
		isMatch, err := regexp.MatchString(pattern, matchingTo)
		if err != nil{
			return false
		}
		if isMatch{
			return false
		}
	}
	for _, pattern := range allowList {
		matchingTo := strings.Replace(path, "*", ".*", -1)
		isMatch, err := regexp.MatchString(pattern, matchingTo)
		if err != nil{
			return false
		}
		if isMatch{
			return true
		}
	}
	return true
}

func main() {
	fmt.Println("Crawler intialising...")
	lines, err := readLines("scrapelist.txt")
	if err != nil{
		log.Fatalln(err)
	}
	fmt.Println("Crawler starting...")

	for i := 0; i < len(lines); i++ {
		fmt.Println("Crawling to site " + lines[i])
		resp, err := http.Get(lines[i] + "/robots.txt")

		var hasRobotsFile = true;
		if err != nil{
			hasRobotsFile = false;
		}else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			sb := string(body)
			if strings.HasPrefix(sb, "<!DOCTYPE html>"){
				hasRobotsFile = false;
			}
		}

		var allowed []string
		var disallowed []string

		var currSite string = lines[i];

		if hasRobotsFile {
			resp, err := http.Get(currSite + "/robots.txt")
			if err != nil {
				log.Fatalln(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			sb := string(body)
			foundRules := strings.Split(strings.ReplaceAll(sb, "\r\n", "\n"), "\n")
			for i := 0; i < len(foundRules); i++ {
				if strings.HasPrefix(foundRules[i], "Allow: "){
					allowed = append(allowed, strings.TrimPrefix(foundRules[i], "Allow: "))
					fmt.Println("Found allowed location " + strings.TrimPrefix(foundRules[i], "Allow: ") + " for " + currSite)
				}else if strings.HasPrefix(foundRules[i], "Disallow: "){
					disallowed = append(disallowed, strings.TrimPrefix(foundRules[i], "Disallow: "))
					fmt.Println("Found disallowed location " + strings.TrimPrefix(foundRules[i], "Disallow: ") + " for " + currSite)
				}
			}
		}else{
			fmt.Println("No robots.txt file found on " + currSite)
		}

		resp, err = http.Get(currSite)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		sb := string(body)
		if err != nil {
			log.Fatalln(err)
		}
		xurls.Relaxed.FindAllString(sb, -1)
	}
}