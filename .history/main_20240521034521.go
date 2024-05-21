package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Düğüm yapısı, bir grafikteki bir düğümü temsil eder.
type Node struct {
	ID   int    // Düğüm ID'si
	Name string // Düğüm Adı
	X    int    // X koordinatı
	Y    int    // Y koordinatı
}

// Kenar yapısı, grafikte iki düğüm arasındaki bir kenarı temsil eder.
type Edge struct {
	Start int // Kenarın başlangıç düğümü ID'si
	End   int // Kenarın bitiş düğümü ID'si
}

// Grafik yapısı, düğümleri, kenarları ve bir grafikteki başlangıç ve bitiş düğümlerini temsil eder.
type Graph struct {
	Nodes       []Node        // Grafikteki tüm düğümlerin listesi
	Edges       []Edge        // Grafikteki tüm kenarların listesi
	StartNodeID int           // Başlangıç düğümünün ID'si
	EndNodeID   int           // Bitiş düğümünün ID'si
	AdjList     map[int][]int // Grafik için komşuluk listesi
}

// Tüm düğümleri yazdırmak için fonksiyon
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: %s (%d, %d)\n", node.ID, node.Name, node.X, node.Y)
	}
}

// Tüm kenarları yazdırmak için fonksiyon
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// BFS kullanarak startNodeID'den endNodeID'ye tüm yolları bulmak için fonksiyon
func (g *Graph) BFSAllPaths(startNodeID int, endNodeID int) [][]int {
	paths := [][]int{}
	queue := [][]int{{startNodeID}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node == endNodeID {
			paths = append(paths, path)
			continue
		}

		for _, neighbor := range g.AdjList[node] {
			if !contains(path, neighbor) {
				newPath := append([]int{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return paths
}

// Belirli bir öğeyi içeren bir dilim kontrol etmek için yardımcı fonksiyon
func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Adıyla düğüm ID'sini bulmak için fonksiyon
func findNodeIDByName(nodes []Node, name string) int {
	for _, node := range nodes {
		if node.Name == name {
			return node.ID
		}
	}
	return -1
}

// Karıncalara yolları atamak için fonksiyon, her karıncanın üzerine binme olmadan en kısa yolu seçer
func assignPathsToAnts(antCount int, filteredPaths [][]int) [][]int {
	antPaths := make([][]int, antCount)
	remainingPaths := filteredPaths

	// Tüm karıncalara yolları ata
	for i := 0; i < antCount; i++ {
		if len(remainingPaths) > 0 {
			antPaths[i] = remainingPaths[0]
			remainingPaths = remainingPaths[1:]
		} else {
			// Eğer daha fazla benzersiz yol yoksa, mevcut yolları tekrar kullan
			antPaths[i] = filteredPaths[i%len(filteredPaths)]
		}
	}

	return antPaths
}

// Karıncaların paylaşmasına izin verilmeyen, yani üzerine binmeyen yolları sağlamak için fonksiyon
func filterPaths(allPaths [][]int, antCount int) [][]int {
	filteredPaths := [][]int{}
	for _, path := range allPaths {
		valid := true
		for _, p := range filteredPaths {
			if hasOverlap(path, p) {
				valid = false
				break
			}
		}
		if valid {
			filteredPaths = append(filteredPaths, path)
			if len(filteredPaths) == antCount {
				break
			}
		}
	}
	return filteredPaths
}

// İki yolun (başlangıç ve bitiş haricinde) üzerine bindiğini kontrol etmek için fonksiyon
func hasOverlap(path1 []int, path2 []int) bool {
	for i := 1; i < len(path1)-1; i++ {
		for j := 1; j < len(path2)-1; j++ {
			if path1[i] == path2[j] {
				return true
			}
		}
	}
}