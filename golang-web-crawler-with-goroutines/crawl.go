// based on: http://stackoverflow.com/questions/42381426/change-the-sample-by-using-goroutine
package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var alreadyCrawledList []string
var pending []string
var brokenLinks []string

const localHostWithPort = "localhost:8080"

func IsLinkInPendingQueue(link string) bool {
	for _, x := range pending {
		if x == link {
			return true
		}
	}
	return false
}

func IsLinkAlreadyCrawled(link string) bool {
	for _, x := range alreadyCrawledList {
		if x == link {
			return true
		}
	}
	return false
}

func AddLinkInAlreadyCrawledList(link string) {
	alreadyCrawledList = append(alreadyCrawledList, link)
}

func AddLinkInPendingQueue(link string) {
	pending = append(pending, link)
}

func AddLinkInBrokenLinksQueue(link string) {
	brokenLinks = append(brokenLinks, link)
}

func TestBrokenLinks(links []string) {
	start := time.Now()
	for _, link := range links {
		AddLinkInPendingQueue(link)
	}

	for count := 0; len(pending) > 0; count++ {
		x := pending[0]
		pending = pending[1:]

		if err := crawlPage(x); err != nil {
			fmt.Errorf(err.Error())
		}
	}
	duration := time.Since(start)
	fmt.Println("________________")
	count := 0
	for _, l := range alreadyCrawledList {
		count++
		fmt.Println(count, "OK. | ", l)
	}

	count = 0
	for _, l := range brokenLinks {
		count++
		fmt.Println(count, "Broken. | ", l)
	}
	fmt.Println("Time taken:", duration)
}

func crawlPage(uri string) error {
	if IsLinkAlreadyCrawled(uri) {
		fmt.Println("Already visited: Ignoring uri | ", uri)
		return nil
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		fmt.Println("Got error: ", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		AddLinkInBrokenLinksQueue(uri)
		return errors.New(fmt.Sprintf("Got %v instead of 200", resp.StatusCode))
	}

	defer resp.Body.Close()

	links := ParseLinks(resp.Body)
	links = ConvertLinksToLocalHost(links)

	for _, link := range links {
		if !InOurDomain(link) {
			continue
		}

		absolute := FixURL(link, uri)
		if !IsLinkAlreadyCrawled(absolute) && !IsLinkInPendingQueue(absolute) && absolute != uri { // Don't enqueue a page twice!
			AddLinkInPendingQueue(absolute)
		}
	}
	AddLinkInAlreadyCrawledList(uri)
	return nil
}

func InOurDomain(link string) bool {
	uri, err := url.Parse(link)
	if err != nil {
		return false
	}

	if uri.Scheme == "http" || uri.Scheme == "https" {
		if uri.Host == localHostWithPort {
			return true
		}
		return false
	}
	return true
}

func ConvertLinksToLocalHost(links []string) []string {
	var convertedLinks []string
	for _, link := range links {
		convertedLinks = append(convertedLinks, strings.Replace(link, "leantricks.com", localHostWithPort, 1))
	}
	return convertedLinks
}

func FixURL(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

func ParseLinks(httpBody io.Reader) []string {
	var links []string
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}

		token := page.Token()
		switch tokenType {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			switch token.DataAtom.String() {
			case "a":
				fallthrough
			case "link":
				fallthrough
			case "script":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

func main() {
	httpLocalHostWithPort := "http://" + localHostWithPort
	TestBrokenLinks([]string{
		httpLocalHostWithPort + "/about",
		httpLocalHostWithPort + "/broken",
		httpLocalHostWithPort + "/home",
		httpLocalHostWithPort + "/about", // revisiting
		httpLocalHostWithPort + "/broken_again"})
}