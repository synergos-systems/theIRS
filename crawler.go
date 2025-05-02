package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

type Version struct {
    Schedule string
    Major int
    Minor int
    Sep string
}

var ledger map[string]Version

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

func UnpackSchemas() (map[string]Version, error) {
    ledger = make(map[string]Version)
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
    fmt.Printf("%+v", ledger)
    return ledger, nil
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
        fmt.Println(uri)
        res := strings.Split(uri, "/")
        schedule := res[5]
        
        var sep string
        var major, minor int
        var test, register []string
        if strings.Contains(schedule, "v") {
            test = strings.Split(schedule, "v")
            register = strings.Split(test[0], "-")
        } else {
            register = strings.Split(schedule, "-")
        }
         
        if len(test) > 0 {
            var erra error
            major, erra = strconv.Atoi(string(test[1][0]))
            if erra != nil {
                fmt.Println(erra)
            }

            minor, erra = strconv.Atoi(string(test[1][2]))
            if erra != nil {
                fmt.Println(erra)
            }
            sep = string(test[1][1])
        }

        year := register[len(register) - 1]
        key := year + ":" + string(register[0][3])
        if found, ok := ledger[key]; ok {
            if found.Major < major {
                ledger[key] = Version{
                    Schedule: schedule,
                    Major: major,
                    Minor: minor,
                    Sep: sep,
                }
            } else if found.Major == major {
                if found.Minor < minor {
                    ledger[key] = Version{
                        Schedule: schedule,
                        Major: major,
                        Minor: minor,
                        Sep: sep,
                    } 
                }
            }
        } else {
            ledger[key] = Version{
                Schedule: schedule,
                Major: major,
                Minor: minor,
                Sep: sep,
            } 
        }

        return res[5]
    case "zips":
        res := strings.Split(uri, "/")
        
        return res[7]
    }

    return "" 
}


func fetchSchema(uri string) {
    fmt.Println(uri)
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
    fmt.Println(year)
    out, err := os.Create(fmt.Sprintf(`./data/990_xsd/%s`, year))
    if err != nil {
        fmt.Println(err)
    }
    defer out.Close()

    _, err = io.Copy(out, res.Body)
    if err != nil {
        log.Println(err)
    }
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
