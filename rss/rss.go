package rss

import (
	"encoding/xml"
	"net/http"
	"time"
)

// RSSFeed is the root of the RSS feed XML document
type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Language string `xml:"language"`
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

func urlToRSSFeed(url string) (RSSFeed, error) {
	httpClient := http.Client {
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err!= nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	var rssFeed RSSFeed
	err = xml.NewDecoder(resp.Body).Decode(&rssFeed)
	if err!= nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}


type AtomFeed struct {
    XMLName xml.Name `xml:"feed"`
    Title   string   `xml:"title"`
    Link    []struct {
        Rel  string `xml:"rel,attr"`
        Href string `xml:"href,attr"`
    } `xml:"link"`
    Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
    Title   string `xml:"title"`
    Link    struct {
        Href string `xml:"href,attr"`
    } `xml:"link"`
    ID        string `xml:"id"`
    Updated   string `xml:"updated"`
    Published string `xml:"published"`
    Summary   string `xml:"summary"`
    Author    struct {
        Name string `xml:"name"`
        URI  string `xml:"uri"`
    } `xml:"author"`
    Content struct {
        Data string `xml:",innerxml"`
    } `xml:"content"`
}

func urlToAtomFeed(url string) (AtomFeed, error) {
    httpClient := http.Client{
        Timeout: 10 * time.Second,
    }
    resp, err := httpClient.Get(url)
    if err != nil {
        return AtomFeed{}, err
    }
    defer resp.Body.Close()

    var atomFeed AtomFeed
    err = xml.NewDecoder(resp.Body).Decode(&atomFeed)
    if err != nil {
        return AtomFeed{}, err
    }

    return atomFeed, nil
}