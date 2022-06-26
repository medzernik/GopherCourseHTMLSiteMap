package main

import (
	link "HTMLLinkParser"
	"flag"
	"net/http"
	"regexp"
	"strings"
)

var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page
    <span> some span  </span>
  </a>
  <a href="/page-two">A link to a second page</a>
</body>
</html>
`

var url *string
var urlStripped string

func main() {
	url = flag.String("URL", "http://calhoun.io", "Put your URL here to access the site")
	flag.Parse()

	links := getAWebPageBody(*url)

	urlStripped = stripInputLink(*url)

	linksToVisit := checkLinksDomainValid(links)
	evaluateLinks(&linksToVisit)

	//fmt.Printf("Links to visit in a domain: %+v", linksToVisit)
}

func checkLinksDomainValid(links []link.Link) []link.Link {
	var newSubListLinks []link.Link
	for _, j := range links {
		if strings.Contains(j.Href, "mailto:") {
			checkEmail := regexp.MustCompile("mailto:*")
			urlStripped = strings.ReplaceAll(j.Href, checkEmail.String(), "")

		} else if strings.Contains(j.Href, urlStripped) || strings.Contains(j.Href, *url) {
			newSubListLinks = append(newSubListLinks, j)
		}
	}
	return newSubListLinks
}

//gets the body of a function
func getAWebPageBody(url string) []link.Link {
	result, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer result.Body.Close()

	body, err := link.Parse(result.Body)
	if err != nil {
		panic(err)
	}
	return body

}

var alreadyVisitedLinks = make(map[link.Link]bool)

//checks and orders the links (recursive function)
func evaluateLinks(links *[]link.Link) {
	var linksToVisitFurther []link.Link
	for _, j := range *links {
		alreadyVisitedLinks[j] = true
	}
	evaluateLinks(&linksToVisitFurther)

}

//strip the url to a short form (for domain purposes)
func stripInputLink(input string) string {
	var urlStripped string

	if strings.Contains(input, "https://") {
		urlStripped = strings.ReplaceAll(input, "https://", "")
	} else if strings.Contains(input, "http://") {
		urlStripped = strings.ReplaceAll(input, "http://", "")
	}

	return urlStripped
}
