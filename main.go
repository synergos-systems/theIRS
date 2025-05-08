package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
)

func confirmation(s string, tries int) bool {
    r := bufio.NewReader(os.Stdin)

    for ; tries > 0; tries-- {
        fmt.Printf("%s Proceed? [y/n]: ", s)

        res, err := r.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        // Empty input (i.e. "\n")
        if len(res) < 2 {
            continue
        }

        return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
    }

    return false
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Nah need a command")
        return
    } else if len(os.Args) > 2 {
        fmt.Println("too many")
        return 
    }



    switch os.Args[1] {
    case "zips":
        proceed := confirmation(`
        This will delete any and all zip files in the ./data/990_zips directory.
        Do not proceed with this command if you have already used it.

        `, 3)
        if proceed {
            zips, err := UnpackZips()
            if err != nil {
                fmt.Println(err)
            }
            fmt.Println(zips)
        } else {
            fmt.Println("Aborting")
        }
        break

    case "schemas":
        versions, err := UnpackSchemas()
        if err != nil {
            fmt.Println(err)
        }
        links := generateLinks(versions) 
        fmt.Println(links)
        UnzipSchemas()
        files, err := GlobWalk("./data/990_xsd/output", "*.xsd")
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println(files) 

        cmd := exec.Command("bash", "chmod x+a ./models.sh; ./models.sh")
        if err := cmd.Run(); err != nil {
            fmt.Println("pipeline failed to run", err)
        } else {
            log.Println("Completed pipeline collapse")
        }

        break

    case "unzip":
       XMLParser() 
        

    default:
        fmt.Println("the argument provided doesn't exist")
    }
}

func generateLinks(versions map[string]Version) []string {
    var links []string  

    for _, version := range versions {
        links = append(links, fmt.Sprintf(`https://www.irs.gov/pub/irs-tege/%s`, version.Schedule))
    }

    return links
}

