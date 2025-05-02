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

var fileTracker atomic.Int32
var (
    fileYear = ""
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

func UnpackSchemas() ([]string, error) {
    res, err := http.Get("https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-teos-schemas")
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    doc, err := html.Parse(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    os.Mkdir("./data/990_xsd", 0777)

    var links []string
    var walk func(*html.Node)
    walk = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    if strings.Contains(attr.Val, ".zip") {
                        fetchSchema(attr.Val)
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

    os.Mkdir(`./data/990_zips/`, 0777) 

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

func splitYear(uri string, strategy string) string {
    switch strategy {
    case "schema":
        res := strings.Split(uri, "-")
        year := strings.Split(res[len(res) - 1], "v")
        
        return year[0]
    case "zips":
        res := strings.Split(uri, "/")
        
        return res[7]
    }

    return "" 
}


func fetchSchema(uri string) {
    res, err := http.Get(uri)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    year := splitYear(uri, "schema")
    if year != fileYear {
        fileTracker.Store(0)
        fileYear = year
    }

    out, err := os.Create(fmt.Sprintf(`./data/990_xsd/%s_%d.zip`, year, fileTracker.Load()))
    if err != nil {
        fmt.Println(err)
    }
    defer out.Close()

    _, err = io.Copy(out, res.Body)
    fmt.Println(err)
}

func fetchZip(uri string) string {
    res, err := http.Get(uri)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    year := splitYear(uri, "zips")
    if year != fileYear {
        fileTracker.Store(0)
        fileYear = year
    }

    tracker := fmt.Sprintf(`./data/990_zips/%s_%d.zip`, year, fileTracker.Load())
    out, err := os.Create(tracker)
    if err != nil {
        fmt.Println(err)
    }
    defer out.Close()

    _, err = io.Copy(out, res.Body)
    fmt.Println(err)
    
    fileTracker.Add(1)
    return tracker 
}
