// Finley McIlwaine
// COSC 5825 - Inverted index text search
// Dec. 8, 2019

package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "regexp"
    "strings"
//    "encoding/json"
)

func main() {
    // Get the path of the directory to walk over
    searchPath := os.Args[1]
    ixFiles := []string {}

    // Walk over the directory, getting the list of files to index
    err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("Something went wrong accessing path '%q': %v\n", path, err)
            return err
        }
        if !info.IsDir() {
            // Only index .txt or .md files
            if ext := filepath.Ext(path); ext == ".md" || ext == ".txt" {
                ixFiles = append(ixFiles,path)
            }
        }
        return nil
    })
    if err != nil {
        fmt.Printf("Something went wrong walking path '%q': %v\n", searchPath, err)
        return
    }

    // Read the files, putting the stuff in indexes
    wIndex := WordIndex{}
    wIndex.Init()
    wpIndex := WordPairIndex{}
    wpIndex.Init()
    pairs := map[string]map[string]bool{}
    notWord, err := regexp.Compile("[^a-zA-Z0-9]+")
    if err != nil {
        fmt.Printf("Something went wrong generating the alphanumeric regular expression: %v\n",err)
        return
    }
    for i, fp := range ixFiles {
        fBytes, err := ioutil.ReadFile(fp)
        if err != nil {
            fmt.Printf("Something went wrong reading file '%q': %v\n", fp, err)
            return
        }
        fText := string(fBytes)
        fText  = notWord.ReplaceAllString(fText," ")
        words := strings.Fields(fText)
        for j, w := range words {
            wIndex.addWord(w, i, j)
            if _, ok := pairs[w]; !ok {
                pairs[w] = map[string]bool{}
            }
            if j>0 {
                if _, ok := pairs[words[j-1]][w]; ok {
                    wpIndex.addWordPair(words[j-1],w,i,j)
                    oldW := wIndex.Index[words[j-1]]
                    wIndex.Index[words[j-1]] = word{ oldW.Ft,oldW.Fdt,true }
                } else {
                    pairs[words[j-1]][w] = true
                }
            }
        }
    }

    // Write the indexes to a file
    // indexJson, err := json.MarshalIndent(wIndex, "", "    ")
    if err!=nil {
        fmt.Printf("Something went wrong marshaling the word index: %v\n", err)
        return
    }
    fmt.Printf("The word index json:\n%+v",wIndex)
}
