package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	ID int
	X  int
	Y  int
}

type Edge struct {
	Start   int
	End     int
	Weight  int // Kenarın ağırlığı
	Reverse bool // Kenarın tersi, geçici olarak ağırlığı tersine çevirmek için kullanılacak
}

type Graph struct {
	Nodes       []Node
	Edges       []Edge
	StartNodeID int
	EndNodeID   int
	AntCounts   map[int]int
}

func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: (%d, %d)\n", node.ID, node.X, node.Y)
	}
}

func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d (Weight: %d, Reverse: %t)\n", edge.Start, edge.End, edge.Weight, edge.Reverse)
	}
}

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Dosya açma hatasi:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		AntCounts: make(map[int]int),
	}

	nodeID := 1
	antCount := graph.CheckAntCount(nodeID)
	fmt.Println("Node", nodeID, "deki karınca sayısı:", antCount)
	startID, endID := -1, -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if antCount == 0 {
			antCount, err = strconv.Atoi(line)
			if err != nil {
				return
			}
			fmt.Println("number_of_ants:", antCount)
		}
		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			startID, _ = strconv.Atoi(fields[0])
			startX, _ := strconv.Atoi(fields[1])
			startY, _ := strconv.Atoi(fields[2])
			graph.StartNodeID = startID
			graph.Nodes = append(graph.Nodes, Node{ID: startID, X: startX, Y: startY})
			fmt.Println("start_room:", graph.StartNodeID)
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			endID, _ = strconv.Atoi(fields[0])
			endX, _ := strconv.Atoi(fields[1])
			endY, _ := strconv.Atoi(fields[2])
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, X: endX, Y: endY})
			fmt.Println("end_room:", graph.EndNodeID)
		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				id, _ := strconv.Atoi(fields[0])
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				graph.Nodes = append(graph.Nodes, Node{ID: id, X: x, Y: y})
			} else if len(fields) == 1 && !strings.Contains(line, "-") {
				antCount, _ = strconv.Atoi(fields[0])
			} else if len(fields) == 1 && strings.Contains(line, "-") {
				edgeParts := strings.Split(fields[0], "-")
				start, _ := strconv.Atoi(edgeParts[0])
				end, _ := strconv.Atoi(edgeParts[1])
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Weight: 1, Reverse: false})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatasi:", err)
		return
	}

	if startID == -1 || endID == -1 {
		fmt.Println("Hata: Başlangıç veya bitiş düğümü belirtilmemiş.")
		return
	}
	fmt.Println("son durum:")
	printNodes(graph.Nodes)
	printEdges(graph.Edges)

	// Suurballe algoritmasını uygula
	shortestPaths := graph.Suurballe(graph.StartNodeID)

	// En kısa yolları yazdır
	fmt.Println("\nEn kısa yollar:")
	for end, path := range shortestPaths {
		fmt.Printf("Başlangıç düğümünden %d düğümüne en kısa yol: %v\n", end, path)
	}
}

func (g *Graph) CheckAntCount(nodeID int) int {
	antCount, exist := g.AntCounts[nodeID]
	if !exist {
		return 0
	}
	return antCount
}

func (g *Graph) Suurballe(start int) map[int][]int {
	dist := make(map[int]int)
	prev := make(map[int]int)

	// Başlangıç düğümünden diğer tüm düğümlere olan mesafeleri sonsuza ve önceki düğümleri nil olarak başlat
	for _, edge := range g.Edges {
		dist[edge.Start] = math.MaxInt64
		dist[edge.End] = math.MaxInt64
		prev[edge.Start] = -1
		prev[edge.End] = -1
	}

	// Başlangıç düğümüne olan mesafeyi kendine eşitle
	dist[start] = 0

	// İlk Dijkstra geçişi
	for i := 0; i < len(g.Edges)-1; i++ {
		for _, edge := range g.Edges {
			if dist[edge.Start]+edge.Weight < dist[edge.End] {
				dist[edge.End] = dist[edge.Start] + edge.Weight
				prev[edge.End] = edge.Start
			}
		}
	}

	// İkinci Dijkstra geçişi
	for i := 0; i < len(g.Edges)-1; i++ {
		for _, edge := range g.Edges {
			// Geçici olarak kenarın ağırlığını tersine çevir
			edge.Weight *= -1

			if dist[edge.Start]+edge.Weight < dist[edge.End] {
				dist[edge.End] = dist[edge.Start] + edge.Weight
				prev[edge.End] = edge.Start
			}

			// Kenarın ağırlığını geri çevir
			edge.Weight *= -1
		}
	}

	// En kısa yolları hesapla
	shortestPaths := make(map[int][]int)
	for end, _
