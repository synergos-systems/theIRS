package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func XMLParser() {
    unzipFile()
}

func unzipFile() {
    pathway := "./data/990_zips/"
    reader, err := os.ReadDir(pathway)
    if err != nil {
        fmt.Println(err)
    }


    for _, zipper := range reader {
        zipperName := pathway + zipper.Name() 


        zReader, err := zip.OpenReader(zipperName)
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println(zipperName[:len(zipperName) - 4])
        err = os.MkdirAll(zipperName[:len(zipperName) - 4], os.FileMode(os.O_CREATE | os.O_TRUNC | os.O_RDWR))
        if err != nil {
            fmt.Println(err)
        }

        for _, file := range zReader.File {
            rc, err := file.Open()
            if err != nil {
                fmt.Println(err)
            }

            newFile, err := os.Open(zipperName[:len(zipperName) - 4])
            if err != nil {
                fmt.Println(err)
            }

            io.Copy(newFile, rc)
        }

        zReader.Close()
    }
    
}


