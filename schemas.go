package main

import (
	"fmt"
	"os"
)


func SchemaGenerator(uri string) {
    entries, err := os.ReadDir("./data/990_xsd/")
    if err != nil {
        fmt.Println(err)
    }

    for _, val := range entries {
        fmt.Println(val)
    }
}
