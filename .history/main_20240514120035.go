package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	ID   string
	X, Y int
}

type Edge struct {
	Start, End string
}

type Graph struct {
	Nodes       map[string]Node
	Edges       []Edge
	StartNodeID string
	EndNodeID   string
}

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Dosya açma hatasi:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		Nodes: make(map[string]Node),
	}

	scanner := bufio.NewScanner(file)
	var antCount int
	for scanner.Scan() {
		line := scanner.Text()

		// Karınca sayısı
		if antCount == 0 {
			fmt.Println("Karınca sayısı:", line)
			antCount++
			continue
		}

		// Düğüm veya kenar işlemleri
		if strings.HasPrefix(line, "##start") {
			node := parseNode(strings.Fields(scanner.Text()))
			graph.Nodes[node.ID] = node
			graph.StartNodeID = node.ID
		} else if strings.HasPrefix(line, "##end") {
			node := parseNode(strings.Fields(scanner.Text()))
			graph.Nodes[node.ID] = node
			graph.EndNodeID = node.ID
		} else if strings.Contains(line, "-") {
			edgeParts := strings.Split(line, "-")
			start := strings.TrimSpace(edgeParts[0])
			end := strings.TrimSpace(edgeParts[1])
			graph.Edges = append(graph.Edges, Edge{Start: start, End: end})
		} else if len(strings.TrimSpace(line)) > 0 {
			node := parseNode(strings.Fields(line))
			graph.Nodes[node.ID] = node
		}
	}

	// Hata kontrolü
	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatasi:", err)
		return
	}

	// Düğümleri ve kenarları yazdır
	fmt.Println("Düğümler:")
	for _, node := range graph.Nodes {
		fmt.Printf("%s: (%d, %d)\n", node.ID, node.X, node.Y)
	}

	fmt.Println("\nKenarlar:")
	for _, edge := range graph.Edges {
		fmt.Printf("%s - %s\n", edge.Start, edge.End)
	}

	// Tüm yolları bul
	allPathsBFS := findAllPathsBFS(graph, graph.EndNodeID, graph.StartNodeID)

	// Tüm yolları sondan başa doğru yazdır
	fmt.Println("\nİlk odadan son odaya giden tüm yollar:")
	for i := len(allPathsBFS) - 1; i >= 0; i-- {
		fmt.Println(allPathsBFS[i])
	}
}

func parseNode(fields []string) Node {
	id := fields[0]
	x, _ := strconv.Atoi(fields[1])
	y, _ := strconv.Atoi(fields[2])
	return Node{ID: id, X: x, Y: y}
}

func findAllPathsBFS(graph Graph, startID, endID string) [][]string {
	// Tüm yolları tutacak slice
	var allPaths [][]string

	// BFS için kuyruk oluştur
	queue := [][]string{{startID}}

	// BFS döngüsü
	for len(queue) > 0 {
		// Kuyruğun başındaki yolu ve son düğümü al
		currentPath := queue[0]
		queue = queue[1:]
		currentNode := currentPath[len(currentPath)-1]

		// Bitiş düğümünü bulduysak, bu yolu ekle
		if currentNode == endID {
			allPaths = append(allPaths, currentPath)
			// Yeni yollar bulmak için devam et
			continue
		}

		// Mevcut düğümün komşularını al ve kuyruğa ekle
		for _, edge := range graph.Edges {
			if edge.Start == currentNode {
				// Yeni yolu oluştur
				newPath := append([]string{}, currentPath...)
				newPath = append(newPath, edge.End)
				// Kuyruğa ekle
				queue = append(queue, newPath)
			}
		}
	}

	return allPaths

}
