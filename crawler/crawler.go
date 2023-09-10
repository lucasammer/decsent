package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mvdan/xurls"
	"golang.org/x/exp/slices"
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
		matchingTo := strings.Replace(pattern, "*", ".*", -1)
		isMatch, err := regexp.MatchString(matchingTo, path)
		if err != nil{
			return false
		}
		if isMatch{
			return false
		}
	}
	for _, pattern := range allowList {
		matchingTo := strings.Replace(pattern, "*", ".*", -1)
		isMatch, err := regexp.MatchString(matchingTo, path)
		if err != nil{
			return false
		}
		if isMatch{
			return true
		}
	}
	return true
}

var visited []string;

func contains(list []string, query string) (listHasQuery bool){
	for _, item := range list {
		if item == query{
			return true
		}
	}
	return false
}

func externalURL(urllink string) {
	if contains(visited, urllink){
		fmt.Println("Already visited adress")
		return
	}
	parsedURL, err := url.Parse(urllink)
	if err != nil {
		fmt.Println("Failed to parse")
		return;
	}
	visited = append(visited, urllink)
	fmt.Println("preparing external url " + urllink)
	parsedURL.Path = ""
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	resp, err := http.Get(parsedURL.String() + "/robots.txt")
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
	if hasRobotsFile {
		resp, err := http.Get(parsedURL.String() + "/robots.txt")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		sb := string(body)
		foundRules := strings.Split(strings.ReplaceAll(sb, "\r\n", "\n"), "\n")
		// sitemap := "none"
		for i := 0; i < len(foundRules); i++ {
			if strings.HasPrefix(foundRules[i], "Allow: "){
				allowed = append(allowed, strings.TrimPrefix(foundRules[i], "Allow: "))
				fmt.Println("Found allowed location " + strings.TrimPrefix(foundRules[i], "Allow: ") + " for " + parsedURL.String())
			}else if strings.HasPrefix(foundRules[i], "Disallow: "){
				disallowed = append(disallowed, strings.TrimPrefix(foundRules[i], "Disallow: "))
				fmt.Println("Found disallowed location " + strings.TrimPrefix(foundRules[i], "Disallow: ") + " for " + parsedURL.String())
			}
			// else if strings.HasPrefix(foundRules[i], "Sitemap: "){
			// 	sitemap = strings.TrimPrefix(foundRules[i], "Sitemap: ")
			// 	followSitemap(sitemap, urllink)
			// 	return
			// }
		}
	}else{
		fmt.Println("No robots.txt file found on " + parsedURL.String())
	}
	parsedURL, err = url.Parse(urllink)
	if err != nil {
		log.Fatalln(err)
	}
	if isAllowed(parsedURL.Path, disallowed, allowed){
		fmt.Println("visiting " + urllink + " as result of external url")
		forOneUrl(urllink)
	}else{
		fmt.Println("NOT visiting " + urllink + " as result of external url")
	}
}

func forOneUrl(urllink string) {
	if contains(visited, urllink){
		fmt.Println("Already visited adress")
		return
	}
	visited = append(visited, urllink)
	fmt.Println("Crawling to site " + urllink)
	resp, err := http.Get(urllink + "/robots.txt")

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

	var currSite string = urllink

	if hasRobotsFile {
		resp, err := http.Get(currSite + "/robots.txt")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		sb := string(body)
		foundRules := strings.Split(strings.ReplaceAll(sb, "\r\n", "\n"), "\n")
		// sitemap := "none"
		for i := 0; i < len(foundRules); i++ {
			if strings.HasPrefix(foundRules[i], "Allow: "){
				allowed = append(allowed, strings.TrimPrefix(foundRules[i], "Allow: "))
				fmt.Println("Found allowed location " + strings.TrimPrefix(foundRules[i], "Allow: ") + " for " + currSite)
			}else if strings.HasPrefix(foundRules[i], "Disallow: "){
				disallowed = append(disallowed, strings.TrimPrefix(foundRules[i], "Disallow: "))
				fmt.Println("Found disallowed location " + strings.TrimPrefix(foundRules[i], "Disallow: ") + " for " + currSite)
			}
			// else if strings.HasPrefix(foundRules[i], "Sitemap: "){
			// 	sitemap = strings.TrimPrefix(foundRules[i], "Sitemap: ")
			// 	followSitemap(sitemap, urllink)
			// 	return
			// }
		}
	}else{
		fmt.Println("No robots.txt file found on " + currSite)
	}

	resp, err = http.Get(currSite)
	if err != nil {
		fmt.Println("Failed to GET " + currSite)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	sb := string(body)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader((sb)))
    if err != nil {
        log.Fatal(err)
    }
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "description" {
			description, _ := s.Attr("content")
			fmt.Printf("Description field: %s\n", description)
		}
		if name, _ := s.Attr("name"); name == "title" {
			title, _ := s.Attr("content")
			fmt.Printf("Title field: %s\n", title)
		}
	})

	linksInPage := xurls.Relaxed.FindAllString(sb, -1)
	splitquote := strings.Split(sb, "\"")
	for i := 0; i < len(splitquote); i++ {
		if strings.HasSuffix(splitquote[i], "href="){
			parsedUrl, err := url.Parse(urllink);
			if err != nil {
				log.Fatalln(err)
			}
			parsedUrl.Path = ""
			parsedUrl.RawQuery = ""
			parsedUrl.Fragment = ""
			if splitquote[i+1][0] == '/' || splitquote[i+1][0] == '#' {
				linksInPage = append(linksInPage, parsedUrl.String() + splitquote[i+1])
			}else if strings.HasPrefix(splitquote[i+1], "http"){
				linksInPage = append(linksInPage, splitquote[i + 1])
			}else{
				linksInPage = append(linksInPage, parsedUrl.String() + "/" + splitquote[i+1])
			}
		}
	}
	fmt.Println("all links in page: " + strings.Join(linksInPage, " & "))
	thisURL, err := url.Parse(currSite)
	if err != nil {
		log.Fatalln(err)
	}
	thisHost := thisURL.Hostname()
	var localToVisit []string
	for _, link := range linksInPage {
		parsedUrl, err := url.Parse(link);
		if err != nil {
			fmt.Println(err)
			continue
		}
		hostName := parsedUrl.Hostname()
		if hostName != thisHost {
			fmt.Println(link + " is an external url.");
			externalURL(link)
		}else{
			fmt.Println(link + " allowed: " + strconv.FormatBool(isAllowed(parsedUrl.Path, disallowed, allowed)))
			if isAllowed(parsedUrl.Path, disallowed, allowed){
				localToVisit = append(localToVisit, parsedUrl.Path)
			}
		}
	}
	fmt.Println("to visit: ", strings.Join(localToVisit, ", "))

	for _, path := range localToVisit {
		if path != thisURL.Path{
			parsedUrl, err := url.Parse(urllink);
			if err != nil {
				log.Fatalln(err)
			}
			parsedUrl.Path = ""
			parsedUrl.RawQuery = ""
			parsedUrl.Fragment = ""

			if slices.Contains(visited, parsedUrl.String() + path){
				fmt.Println("already been here!")
				continue
			}
			forOneUrl(parsedUrl.String() + path)
		}
	}
}

func main() {
	fmt.Println("Crawler intialising...")
	lines, err := readLines("scrapelist.txt")
	if err != nil{
		log.Fatalln(err)
	}
	fmt.Println("Crawler starting...")

	for i := 0; i < len(lines); i++ {
		forOneUrl(lines[i])
	}

	fmt.Println("All visited sites: " + strings.Join(visited, ", "))
}
