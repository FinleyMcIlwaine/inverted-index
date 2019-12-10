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
    "encoding/json"
    "bufio"
    "math"
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
        // Only index .txt or .md files
        if ext := filepath.Ext(path); ext == ".md" || ext == ".txt" {
            ixFiles = append(ixFiles,path)
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
        words := strings.Fields(strings.ToLower(fText))
        for j, w := range words {
            wIndex.addWord(w, i, j)
            if _, ok := pairs[w]; !ok {
                pairs[w] = map[string]bool{}
            }
            if j>0 {
                if _, ok := pairs[words[j-1]][w]; ok {
                    wpIndex.addWordPair(words[j-1],w,i,j)
                    oldW := wIndex.Index[words[j-1]]
                    wIndex.Index[words[j-1]] = &word{ oldW.Ft,oldW.Fdt,true }
                } else {
                    pairs[words[j-1]][w] = true
                }
            }
        }
    }

    // Write the indexes to a file
    os.Remove("indexes.log");
    f, err := os.OpenFile("indexes.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err!=nil {
        fmt.Printf("Something went wrong opening the log file: %v\n", err)
    }
    indexJson, err := json.MarshalIndent(wIndex, "", "  ")
    if err!=nil {
        fmt.Printf("Something went wrong marshaling the word index: %v\n", err)
        return
    }
    if _, err := f.Write(append([]byte("Word bag index:\n\n"),indexJson...)); err != nil {
        f.Close();
        fmt.Printf("Something went wrong writing the word index to the log: %v\n", err)
        return
    }
    indexJson, err = json.MarshalIndent(wpIndex, "", "  ")
    if err!=nil {
        fmt.Printf("Something went wrong marshaling the pair index: %v\n", err)
        return
    }
    if _, err := f.Write(
        append([]byte("\n\nWord pair index:\n\n"),
        append(indexJson,[]byte("\n")...)...)); err != nil {
            f.Close();
            fmt.Printf("Something went wrong writing the pair index to the log: %v\n", err)
            return
    }

    // Now prompt for a user search and do the actual search
    reader := bufio.NewReader(os.Stdin);
    for {
        fmt.Print("Enter a search query (:q to quit): ")
        text, _ := reader.ReadString('\n')
        if text==":q\n" {
            break
        } else {
            text = notWord.ReplaceAllString(text," ")
            words := strings.Fields(strings.ToLower(text));

            // Accumulators for documents
            aDocs := make([]float64, len(ixFiles))
            wDocs := make([]float64, len(ixFiles))
            sDocs := make([]float64, len(ixFiles))
            // For each word and pair
            for i, w := range words {
                // Calculate wqt for each word
                wqt := wIndex.Wqt(w,len(ixFiles))
                // For each term doc pair
                if wqt > 0 {
                    for _, fdt := range wIndex.Index[w].Fdt {
                        wdt := 1 + math.Log(float64(fdt.Frequency))
                        wDocs[fdt.Document] += wdt*wdt
                        aDocs[fdt.Document] += wqt*wdt
                    }
                }
                if i>0 {
                    pwqt := wpIndex.Wqt(words[i-1], w, len(ixFiles))
                    if pwqt > 0 {
                        for _, fdt := range wpIndex.Index[words[i-1]][w].Fdt {
                            pwdt := 1 + math.Log(float64(fdt.Frequency))
                            wDocs[fdt.Document] += pwdt*pwdt
                            aDocs[fdt.Document] += 4.0*pwqt*pwdt
                        }
                    }
                }
            }
            for i, d := range aDocs {
                if d > 0 {
                    sDocs[i] = d / wDocs[i]
                }
            }
            fCopy := make([]string, len(ixFiles))
            copy(fCopy,ixFiles)
            for i := len(sDocs)-1; i>=0; i-- {
                for j := i; j<len(sDocs)-1; j++ {
                    if sDocs[j+1] > sDocs[j] {
                        sDocs[j+1],sDocs[j] = sDocs[j],sDocs[j+1]
                        fCopy[j+1],fCopy[j] = fCopy[j],fCopy[j+1]
                    }
                }
            }
            fmt.Print("\nResults:\n")
            for i, s := range sDocs {
                fmt.Printf("\n    %d: %s, s = %.5f",i+1,fCopy[i],s)
            }
            fmt.Print("\n\n")
        }
    }
}
