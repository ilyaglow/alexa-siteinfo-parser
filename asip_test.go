package asip

import (
	"io"
	"os"
	"reflect"
	"testing"
)

const (
	successTestDocLoc = "testdata/body.html"
	nodataTestDocLoc  = "testdata/nodata.html"
)

var successTestSite = &Site{
	Title:        "Сбербанк России",
	Description:  "Сведения об истории создания, руководстве, филиалах и подразделениях. Перечень услуг. Тарифы.",
	MainCountry:  "Russia",
	GlobalRank:   506,
	LocalRank:    17,
	LinkingTotal: 8491,
	Visitors: []Visitor{
		Visitor{
			Country:   "Russia",
			Percent:   "83.8%",
			LocalRank: 17,
		},
		Visitor{
			Country:   "Netherlands",
			Percent:   "2.0%",
			LocalRank: 182,
		},
		Visitor{
			Country:   "Germany",
			Percent:   "1.7%",
			LocalRank: 1366,
		},
		Visitor{
			Country:   "United Kingdom",
			Percent:   "1.4%",
			LocalRank: 1234,
		},
		Visitor{
			Country:   "United States",
			Percent:   "1.3%",
			LocalRank: 7997,
		},
	},
	Keywords: []Keyword{
		Keyword{
			Word:    "сбербанк онлайн",
			Percent: "49.69%",
		},
		Keyword{
			Word:    "сбербанк",
			Percent: "7.87%",
		},
		Keyword{
			Word:    "сбербанк бизнес онлайн",
			Percent: "7.74%",
		},
		Keyword{
			Word:    "sberbank online",
			Percent: "3.63%",
		},
		Keyword{
			Word:    "sberbank",
			Percent: "2.65%",
		},
	},
	Upstreams: []Upstream{
		Upstream{
			Site:    "yandex.ru",
			Percent: "21.4%",
		},
		Upstream{
			Site:    "google.com",
			Percent: "10.1%",
		},
		Upstream{
			Site:    "vk.com",
			Percent: "5.6%",
		},
		Upstream{
			Site:    "mail.ru",
			Percent: "4.3%",
		},
		Upstream{
			Site:    "youtube.com",
			Percent: "2.3%",
		},
	},
	LinksFrom: []Link{
		Link{
			Site: "yandex.ru",
			Page: "http://money.yandex.ru/doc.xml?id=242350",
		},
		Link{
			Site: "mail.ru",
			Page: "http://card.krugdoveriya.mail.ru/articles.html?id=19376",
		},
		Link{
			Site: "fc2.com",
			Page: "http://10rank.blog.fc2.com/blog-entry-264.html",
		},
		Link{
			Site: "mit.edu",
			Page: "http://misti.mit.edu/hosts-partners/featured-hosts",
		},
		Link{
			Site: "wixsite.com",
			Page: "http://belov-72.wixsite.com/ocenka72",
		},
	},
	Related: []string{
		"sbrf.ru",
		"sravni.ru",
		"gosuslugi.ru",
		"banki.ru",
		"avito.ru",
	},
	Categories: []string{
		"World",
		"Russian",
		"Страны и регионы",
		"Европа",
		"Россия",
		"Бизнес и экономика",
		"Финансовые услуги",
		"Банки",
	},
	Subdomains: []Subdomain{
		Subdomain{
			Domain:  "online.sberbank.ru",
			Percent: "69.69%",
		},
		Subdomain{
			Domain:  "sberbank.ru",
			Percent: "28.30%",
		},
		Subdomain{
			Domain:  "securepayments.sberbank.ru",
			Percent: "6.72%",
		},
		Subdomain{
			Domain:  "sbi.sberbank.ru",
			Percent: "4.53%",
		},
		Subdomain{
			Domain:  "info.sberbank.ru",
			Percent: "0.58%",
		},
	},
}

func testDoc(filename string) (body io.ReadCloser, err error) {
	return os.Open(filename)
}

func TestSiteInfo(t *testing.T) {
	body, err := testDoc(successTestDocLoc)
	if err != nil {
		t.Fatal(err)
	}

	si, err := parse(body)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(si, successTestSite) {
		t.Fatalf("want %v, got %v", successTestSite, si)
	}
}

func TestNoData(t *testing.T) {
	body, err := testDoc(nodataTestDocLoc)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parse(body)
	if err != ErrNoEnoughData {
		t.Fatal("want error, but got no error")
	}
}
