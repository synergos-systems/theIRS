package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocomply/xsd2go/pkg/xsd2go"
)

func UnzipSchemas() {
    entries, err := os.ReadDir("./data/990_xsd/")
    if err != nil {
        fmt.Println(err)
    }

    for _, val := range entries {
        dst := "./data/990_xsd/output"
        fmt.Println(val)
        archive, err := zip.OpenReader("./data/990_xsd/" + val.Name())
        if err != nil {
            panic(err)
        }
        defer archive.Close()

        for _, f := range archive.File {
            filePath := filepath.Join(dst, f.Name)
            fmt.Println("unzipping file ", filePath)

            if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
                fmt.Println("invalid file path")
                return
            }
            if f.FileInfo().IsDir() {
                fmt.Println("creating directory...")
                os.MkdirAll(filePath, os.ModePerm)
                continue
            }

            if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
                panic(err)
            }

            dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                panic(err)
            }

            fileInArchive, err := f.Open()
            if err != nil {
                panic(err)
            }

            if _, err := io.Copy(dstFile, fileInArchive); err != nil {
                panic(err)
            }
            
            dstFile.Close()
            fileInArchive.Close()
        }
    }
}

func GlobWalk(rootDir, pattern string) ([]string, error) {
    var matches []string

    err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, walkErr error) error {
        if walkErr != nil {
            // abort the walk on error
            return walkErr
        }

        if d.IsDir() {
            // skip directories
            return nil
        }

        ok, err := filepath.Match(pattern, filepath.Base(path))
        if err != nil {
            // bad pattern
            return err
        }
        if ok {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}

func SchemaGenerator(uri string) {
    err := xsd2go.Convert(
        uri,
        "main",
        "./data/990_xsd/output/generated_templates",
        nil,
    )
    if err != nil {
        fmt.Println(err)
    }
}
