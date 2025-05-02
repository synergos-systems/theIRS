package main

import "fmt"

func main() {
    //ConstructURLs()
   // links, err := GenerateSchemas()
   // if err != nil {
   //     fmt.Println(err)
   // }

    zips, err := UnpackZips()
    if err != nil {
        fmt.Println(err)
    }
    

    fmt.Println(zips)
}

