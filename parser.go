package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ParseXMLs() {
    pathway := "./data/990_zips/"

    reader, err := os.ReadDir(pathway)
    if err != nil {
        fmt.Println(err)
    }
    
    re := regexp.MustCompile(`.zip`)
    for _, zipper := range reader {
        if !re.Match([]byte(zipper.Name())) {
            zReader, err := os.ReadDir(pathway + zipper.Name())
            if err != nil {
                fmt.Println(err)
            }
            convertCollection(pathway + zipper.Name(), zReader)     
        }
    }
}

func convertCollection(root string, files []os.DirEntry) {
    for _, file := range files {
        fmt.Println(root + file.Name())
        f, err := os.Open(root + "/" + file.Name())
        if err != nil {
            panic(err)
        }
        
        
        decoder := xml.NewDecoder(f)
        for {
            tok, err := decoder.Token()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatal(err)
            }

            switch elem := tok.(type) {
            case xml.StartElement:
                fmt.Println("Start element:", elem.Name.Local)
            case xml.EndElement:
                fmt.Println("End element:", elem.Name.Local)
            case xml.CharData:
                data := strings.TrimSpace(string(elem))
                if len(data) > 0 {
                    fmt.Println("Text:", data)
                }
            }
        }
    }
}

func UnzipXMLs() {
    pathway := "./data/990_zips/"

    reader, err := os.ReadDir(pathway)
    if err != nil {
        fmt.Println(err)
    }

    for _, zipper := range reader {
        template := pathway + zipper.Name()
        unzipXMLs(template, template[:len(template) - 4])
    }
}

func unzipXMLs(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0777)
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, f.Name)
        if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}

