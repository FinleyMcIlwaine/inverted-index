// Finley McIlwaine
// COSC 5825 - Inverted index text search
// Dec. 8, 2019

package main

type WordIndex struct {
    Index   map[string]word
}

type WordPairIndex struct {
    Index   map[string]map[string]wordPair
}

type fdtData struct {
    Frequency   int
    Document    int
    Positions   []int
}

type word struct {
    Ft      int
    Fdt     []fdtData
    Paired  bool
}

type wordPair struct {
    Ft      int
    Fdt     []fdtData
}

func (wi *WordIndex) Init() {
    wi.Index = make(map[string]word)
}

func (wpi *WordPairIndex) Init() {
    wpi.Index = make(map[string]map[string]wordPair)
}

func (w *word) setToPaired() {
    w.Paired=true
}

func (wi *WordIndex) addWord(t string, d int, p int) {
    if w, ok := wi.Index[t]; ok {
        w.Ft++
        ind := -1
        for i, fdt := range w.Fdt {
            if d == fdt.Document {
                ind = i
                break
            }
        }
        if ind==-1 {
            w.Fdt = append(w.Fdt,fdtData{1,d,[]int{p}})
        } else {
            w.Fdt[ind].Frequency++
            w.Fdt[ind].Positions = append(w.Fdt[ind].Positions, p)
        }
    } else {
        wi.Index[t] = word{ 1, []fdtData{fdtData{1,d,[]int{p}}}, false }
    }
}

func (wpi *WordPairIndex) addWordPair(w1 string, w2 string, d int, p int) {
    if wp, ok := wpi.Index[w1][w2]; ok {
        wp.Ft++
        ind := -1
        for i, fdt := range wp.Fdt {
            if d == fdt.Document {
                ind = i
                break
            }
        }
        if ind==-1 {
            wp.Fdt = append(wp.Fdt, fdtData{1,d,[]int{p}})
        } else {
            wp.Fdt[ind].Frequency++
            wp.Fdt[ind].Positions = append(wp.Fdt[ind].Positions, p)
        }
    } // TODO Handle the else case here!!!!! How to initialize inner map?
}