package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct { //düğümler için struct yapısı
	ID int
	X  int
	Y  int
}

type Edge struct { //kenarlar için struct yapısı
	Start int
	End   int
}

type Graph struct { //düğümler ve kenarlar için oluşturulmuş bir graf struct yapısı
	Nodes []Node
	Edges []Edge
}

func main() {
	file, err := os.Open("text.txt") //dosya açma işlemi
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}

	defer file.Close()

	startID, endID := -1, -1 //ilk önce başlangıç ve bitiş satırlarına bir id atıyoruz
	graph := Graph{}         //struct yapısını değişkene atıyoruz

	// Dosyadan okunan başlangıç ve bitiş düğümleri id'lerini atayın
	if strings.HasPrefix(line, "##start") {
		scanner.Scan()
		fields := strings.Fields(scanner.Text())
		startID, _ = strconv.Atoi(fields[0])
	} else if strings.HasPrefix(line, "##end") {
		scanner.Scan()
		fields := strings.Fields(scanner.Text())
		endID, _ = strconv.Atoi(fields[0])
	}
	
	var antCount int
	for scanner.Scan() {
		line := scanner.Text()
		if antCount == 0 { // Karınca sayısı daha önce belirlenmediyse
			antCount, err = strconv.Atoi(line) // Dosyadan karınca sayısını oku
			if err != nil {
				fmt.Println("Karınca sayısı okunamadı:", err)
				return
			}
			fmt.Println("Karınca sayısı :", antCount)
		}

		fields := strings.Fields(line)
		if len(fields) == 1 {
			// Eğer satır sadece bir sayıdan oluşuyorsa, bu karınca sayısıdır.
			antCount, _ = strconv.Atoi(fields[0])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatası:", err)
		return
	}
}
