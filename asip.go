// Package asip is a Alexa Website Info page parser.
package asip

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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
	asiLocation    = "https://www.alexa.com/siteinfo/%s"
)

// ErrNoEnoughData is returned when a domain is not in top 1M.
const ErrNoEnoughData = "no enough data"

// Conf is a asip configuration.
type Conf struct {
	client *http.Client
}

// NewWithClient bootstraps configuration with a customized client.
func NewWithClient(c *http.Client) *Conf {
	return &Conf{c}
}

// Site is Website Traffic Statistics from alexa.com.
type Site struct {
	Title        string
	Description  string
	MainCountry  string
	GlobalRank   uint
	LocalRank    uint
	LinkingTotal uint
	Visitors     []Visitor
	Keywords     []Keyword
	Upstreams    []Upstream
	Related      []string
	Subdomains   []Subdomain
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

// Subdomain represent subdomains where visitors go from the site.
type Subdomain struct {
	Domain  string
	Percent string
}

type findable interface {
	Find(string) *goquery.Selection
	Text() string
}

func parse(body io.Reader) (*Site, error) {
	d, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	if noEnoughData(d) {
		return nil, errors.New(ErrNoEnoughData)
	}

	var s Site
	gr, err := globalRank(d)
	if err != nil {
		return nil, err
	}
	s.GlobalRank = uint(gr)

	lr, err := localRank(d)
	if err != nil {
		return &s, err
	}
	s.LocalRank = uint(lr)

	country, err := country(d)
	if err != nil {
		return &s, err
	}
	s.MainCountry = country

	lt, err := linkingTotal(d)
	if err != nil {
		return &s, err
	}
	s.LinkingTotal = uint(lt)

	tt, err := title(d)
	if err != nil {
		return &s, err
	}
	s.Title = tt

	dsc, err := description(d)
	if err != nil {
		return &s, err
	}
	s.Description = dsc

	vst, err := visitors(d)
	if err != nil {
		return &s, err
	}
	s.Visitors = vst

	kws, err := keywords(d)
	if err != nil {
		return &s, err
	}
	s.Keywords = kws

	ups, err := upstreams(d)
	if err != nil {
		return &s, err
	}
	s.Upstreams = ups

	ls, err := linksFrom(d)
	if err != nil {
		return &s, err
	}
	s.LinksFrom = ls

	rs, err := related(d)
	if err != nil {
		return &s, err
	}
	s.Related = rs

	cts, err := categories(d)
	if err != nil {
		return &s, err
	}
	s.Categories = cts

	ss, err := subdomains(d)
	if err != nil {
		return &s, err
	}
	s.Subdomains = ss

	return &s, nil
}

type getFunc func(string) (*http.Response, error)

func siteInfo(domain string, f getFunc) (*Site, error) {
	resp, err := f(domain)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, no data for %s?", resp.StatusCode, domain)
	}

	return parse(resp.Body)
}

// SiteInfo parses webpage of Alexa Website Info.
func SiteInfo(domain string) (*Site, error) {
	return siteInfo(fmt.Sprintf(asiLocation, domain), http.Get)
}

// SiteInfo parses webpage of Alexa Website Info with customised parameters.
func (c *Conf) SiteInfo(domain string) (*Site, error) {
	return siteInfo(fmt.Sprintf(asiLocation, domain), c.client.Get)
}

func getUint(d findable, selector string, kind string) (uint64, error) {
	var s string
	if selector != "" {
		s = strings.TrimSpace(d.Find(selector).Text())
	} else {
		s = strings.TrimSpace(d.Text())
	}

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

func getString(d findable, selector string, kind string) (string, error) {
	var s string
	if selector != "" {
		s = strings.TrimSpace(d.Find(selector).Text())
	} else {
		s = strings.TrimSpace(d.Text())
	}
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

func noEnoughData(d *goquery.Document) bool {
	return d.Find(seNoData).Length() > 0
}

func visitors(d *goquery.Document) ([]Visitor, error) {
	tbody := d.Find(seVisitors)
	if tbody.Length() == 0 {
		return nil, errors.New("no visitors found")
	}

	var (
		v                []Visitor
		country, percent string
		countryRank      uint64
	)
	tbody.Find("tr").Each(func(i int, tr *goquery.Selection) {
		country = strings.TrimSpace(tr.Find("td a").Text())
		percent = strings.TrimSpace(tr.Find("td span").First().Text())
		countryRank, _ = getUint(
			tr.Find("td span").Last(),
			"",
			fmt.Sprintf("%d country rank", i),
		)

		v = append(v, Visitor{
			Country:   country,
			Percent:   percent,
			LocalRank: uint(countryRank),
		})
	})

	return v, nil
}

func keywords(d *goquery.Document) ([]Keyword, error) {
	tbody := d.Find(seKeywords)
	if tbody.Length() == 0 {
		return nil, errors.New("no keywords found")
	}

	var (
		ks              []Keyword
		key, percentage string
	)
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		key = strings.TrimSpace(tr.Find("td:first-child span:last-child").Text())
		percentage = strings.TrimSpace(tr.Find("td:last-child span").Text())
		ks = append(ks, Keyword{
			Word:    key,
			Percent: percentage,
		})
	})

	return ks, nil
}

func upstreams(d *goquery.Document) ([]Upstream, error) {
	tbody := d.Find(seUpstreams)
	if tbody.Length() == 0 {
		return nil, errors.New("no upstream servers found")
	}

	var (
		us            []Upstream
		site, percent string
	)
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		site = strings.TrimSpace(tr.Find("td a").Text())
		percent = strings.TrimSpace(tr.Find("td:last-child span").Text())
		us = append(us, Upstream{
			Site:    site,
			Percent: percent,
		})
	})

	return us, nil
}

func linksFrom(d *goquery.Document) ([]Link, error) {
	tbody := d.Find(seLinks)
	if tbody.Length() == 0 {
		return nil, errors.New("no linking sites found")
	}

	var (
		ls         []Link
		site, page string
	)
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		site = strings.TrimSpace(tr.Find("span.word-wrap a").Text())
		page, _ = tr.Find("a.word-wrap").Attr("href")
		ls = append(ls, Link{
			Site: site,
			Page: page,
		})
	})

	return ls, nil
}

func related(d *goquery.Document) ([]string, error) {
	tbody := d.Find(seRelated)
	if tbody.Length() == 0 {
		return nil, errors.New("no related sites found")
	}

	var rs []string
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		rs = append(rs, tr.Find("a").Text())
	})

	return rs, nil
}

func categories(d *goquery.Document) ([]string, error) {
	tbody := d.Find(seCategories)
	if tbody.Length() == 0 {
		return nil, errors.New("no categories found")
	}

	var cts []string
	tbody.Find("a").Each(func(_ int, a *goquery.Selection) {
		cts = append(cts, a.Text())
	})

	return cts, nil
}

func subdomains(d *goquery.Document) ([]Subdomain, error) {
	tbody := d.Find(seSubdomains)
	if tbody.Length() == 0 {
		return nil, errors.New("no subdomains found")
	}

	var (
		ss              []Subdomain
		domain, percent string
	)
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		domain = tr.Find("td:first-child span").Text()
		percent = tr.Find("td:last-child span").Text()
		ss = append(ss, Subdomain{
			Domain:  domain,
			Percent: percent,
		})
	})

	return ss, nil
}
