package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"golang.org/x/net/html"
)

type Crawler struct {}

const (
    upTo20 = `https://apps.irs.gov/pub/epostcard/990/xml/`
    post20 = `https://apps.irs.gov/pub/epostcard/990/xml/`
    currentStart = 2019
    currentYear = 2025
)

var zipTracker atomic.Int32
var (
    zipYear = ""
    schemaList = []string{
        "Form 990N Schema 2024v1.0 ZIP",
        "Form 990T Schema 2024v5.0 ZIP",
        "Form 990X Schema 2024v5.0 ZIP",
        "2023 redacted schema package (990N, 990x, 990T)",
        "Form 990N Schema 2023v1.0 ZIP",
        "Form 990T Schema includes Shared 2023v6.0 ZIP",
        "Form 990T Schema includes Shared 2023v7.0 ZIP",
        "Form 990X Schema 2023v4.0 ZIP",
        "Form 990X Schema 2023v5.0 ZIP",
        "Form 990X Schema 2023v5.1 ZIP",
        "2022 redacted schema package (990N, 990x, 990T)",
        "Form 990N Schema 2022 ZIP",
        "Form 990N Schema 2022v1.0 ZIP",
        "Form 990T Schema includes Shared 2022 ZIP",
        "Form 990T Schema includes Shared 2022v6.0 ZIP",
        "Form 990T Schema includes Shared 2022v7.0 ZIP",
        "Form 990X Schema 2022 ZIP",
        "Form 990X Schema 2022v4.0 ZIP",
        "Form 990X Schema 2022v4.1 ZIP",
        "Form 990X Schema 2022v5.0 ZIP",
        "2021 redacted schema package (990N, 990x, 990T)",
        "Form 990N Schema 2021V1.0 ZIP",
        "Form 990N Schema 2021V1.1 ZIP",
        "Form 990N Schema 2021V1.2 ZIP",
        "Form 990T Schema includes Shared 2021v4.0 ZIP",
        "Form 990T Schema includes Shared 2021v4.1 ZIP",
        "Form 990T Schema includes Shared 2021v4.2 ZIP",
        "Form 990T Schema includes Shared 2021v4.3 ZIP",
        "Form 990T Schema includes Shared 2021v4.4 ZIP",
        "Form 990X Schema 2021v4.0 ZIP",
        "Form 990X Schema 2021v4.1 ZIP",
        "Form 990X Schema 2021v4.2 ZIP",
        "Form 990X Schema 2021v4.3 ZIP",
        "2020 redacted schema package (990N, 990x, 990T)",
        "Form 990N Schema 2020V3.0 ZIP",
        "Form 990T Schema includes Shared 2020v1.0 ZIP",
        "Form 990T Schema includes Shared 2020v1.1 ZIP",
        "Form 990T Schema includes Shared 2020v1.2 ZIP",
        "Form 990T Schema includes Shared 2020v1.3 ZIP",
        "Form 990X Schema 2022v4.0 ZIP",
        "Form 990X Schema 2022v4.1 ZIP",
        "Form 990X Schema 2022v4.2 ZIP",
        "2019 redacted schema package (990N, 990x, 990T)",
        "Form 990N Schema 2019v1.0 ZIP",
        "Form 990X Schema 2019v5.0 ZIP",
        "Form 990X Schema 2019v5.1 ZIP",
        "2018 redacted schema package (990N, 990x, 990T)",
        "efile990N 2018v1.0 ZIP",
        "Form 990X Schema 2018v3.0 ZIP",
        "Form 990X Schema 2018v3.1 ZIP",
        "Form 990X Schema 2018v3.2 ZIP",
        "Form 990X Schema 2018v3.3 ZIP",
    }
)



func ScrapeURLs() {
    var template string
    for year := currentStart; year <= currentYear; year++ {
        counter := 12
        for counter > 0 {
            if year < 2021 {
                template = upTo20 + fmt.Sprintf(`%d/download990xml_%d_%d.zip`, year, year, counter)
            } else {
                template = upTo20 + fmt.Sprintf(`%d/%d_TEOS_XML_%02dA.zip`, year, year, counter)
            }
            fmt.Println(template)
            res, err := http.Get(template)
            if err != nil {
                fmt.Println(err)
            }
            defer res.Body.Close()

            out, err := os.Create(fmt.Sprintf(`%d_%d.zip`, year, counter))
            if err != nil {
                fmt.Println(err)
            }
            defer out.Close()

            _, err = io.Copy(out, res.Body)
            fmt.Println(err)
            counter--
        }
    }
}

func GenerateSchemas() ([]string, error) {
    res, err := http.Get("https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-teos-schemas")
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    doc, err := html.Parse(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    var links []string
    var walk func(*html.Node)
    walk = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    if strings.Contains(attr.Val, ".zip") {
                        links = append(links, attr.Val)
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            walk(c)
        }
    }
    walk(doc)

    return links, nil
}

func UnpackZips() ([]string, error) {
    res, err := http.Get("https://www.irs.gov/charities-non-profits/form-990-series-downloads")
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    doc, err := html.Parse(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    var links []string
    var walk func(*html.Node)
    walk = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    if strings.Contains(attr.Val, ".zip") {
                        links = append(links, attr.Val)
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            walk(c)
        }
    }
    walk(doc)
    
    var zipData []string
    for _, uri := range links {
        zipData = append(zipData, fetchZip(uri))
    }
    return links, nil
}

func splitYear(link string) string {
    res := strings.Split(link, "/")
    
    return res[7]
}

func fetchZip(link string) string {
    res, err := http.Get(link)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()
    
    _ = os.Mkdir(`./data/990_zips/`, 0777) 

    year := splitYear(link)
    if year != zipYear {
        zipTracker.Store(0)
        zipYear = year
    }

    tracker := fmt.Sprintf(`./data/990_zips/%s_%d.zip`, year, zipTracker.Load())
    out, err := os.Create(tracker)
    if err != nil {
        fmt.Println(err)
    }
    defer out.Close()

    _, err = io.Copy(out, res.Body)
    fmt.Println(err)
    
    zipTracker.Add(1)
    return tracker 
}
