package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	seGlobalRank   = "span.globleRank span div strong"
	seLocalRank    = "span.countryRank span div strong"
	seCountry      = "span.countryRank span h4 a"
	seVisitors     = "table#demographics_div_country_table tbody"
	seKeywords     = "table#keywords_top_keywords_table tbody"
	seUpstreams    = "table#keywords_upstream_site_table tbody"
	seLinks        = "table#linksin_table tbody"
	seLinkingTotal = "section#linksin-panel-content div span div span.font-4.box1-r"
	seRelated      = "table#audience_overlap_table tbody"
	seCategories   = "table#category_link_table tbody"
	seSubdomains   = "table#subdomain_table tbody"
	seTitle        = "div.row-fluid.siteinfo-site-summary span div p"
	seDescription  = "section#contact-panel-content div.row-fluid span.span8 p.color-s3"
	seNoData       = "section#no-enough-data"
)

// Site is Website Traffic Statistics from alexa.com.
type Site struct {
	Title        string
	Description  string
	Country      string
	GlobalRank   uint
	LocalRank    uint
	LinkingTotal uint
	Visitors     []Visitor
	Keywords     []Keyword
	Upstreams    []Upstream
	Related      []string
	Subdomains   []string
	Categories   []string
	LinksFrom    []Link
}

// Link is a site and page that links to the website.
type Link struct {
	Site string
	Page string
}

// Visitor represents a variety of visitors from a single country.
type Visitor struct {
	Country   string
	Percent   string
	LocalRank uint
}

// Keyword is a one of the top keywords from search engines.
type Keyword struct {
	Word    string
	Percent string
}

// Upstream sites people visited immediately before this site.
type Upstream struct {
	Site    string
	Percent string
}

func main() {
	f, err := os.Open("testdata/body.html")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	d, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	gr, err := globalRank(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(gr)

	lr, err := localRank(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(lr)

	country, err := country(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(country)

	lt, err := linkingTotal(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(lt)

	tt, err := title(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tt)

	dsc, err := description(d)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dsc)

}

func getUint(d *goquery.Document, selector string, kind string) (uint64, error) {
	s := strings.TrimSpace(d.Find(selector).Text())
	if s == "" {
		return 0, fmt.Errorf("no %s found", kind)
	}

	s = strings.ReplaceAll(s, ",", "") // remove commas from string like 1,111,111

	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func getString(d *goquery.Document, selector string, kind string) (string, error) {
	s := strings.TrimSpace(d.Find(selector).Text())
	if s == "" {
		return "", fmt.Errorf("no %s found", kind)
	}
	return s, nil
}

func globalRank(d *goquery.Document) (uint64, error) {
	return getUint(d, seGlobalRank, "global rank")
}

func localRank(d *goquery.Document) (uint64, error) {
	return getUint(d, seLocalRank, "local rank")
}

func country(d *goquery.Document) (string, error) {
	return getString(d, seCountry, "country")
}

func linkingTotal(d *goquery.Document) (uint64, error) {
	return getUint(d, seLinkingTotal, "linking total")
}

func title(d *goquery.Document) (string, error) {
	return getString(d, seTitle, "site title")
}

func description(d *goquery.Document) (string, error) {
	return getString(d, seDescription, "site description")
}
