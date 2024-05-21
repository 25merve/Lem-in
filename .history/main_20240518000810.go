package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	ID int
	X  int
	Y  int
}

type Edge struct {
	Start int
	End   int
	Cost  int // Kenar maliyeti
}

type Graph struct {
	Nodes       []Node
	Edges       []Edge
	StartNodeID int
	EndNodeID   int
	AntCounts   map[int]int
	AdjList     map[int][]Edge // Komşuluk listesi
}

func (g *Graph) BFS(startNodeID int, endNodeID int) [][]int {
	paths := [][]int{}
	queue := [][]int{{startNodeID}}
	visited := map[int]bool{}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node == endNodeID {
			paths = append(paths, path)
			continue
		}

		if visited[node] {
			continue
		}

		visited[node] = true

		for _, edge := range g.AdjList[node] {
			if !visited[edge.End] {
				newPath := make([]int, len(path))
				copy(newPath, path)
				newPath = append(newPath, edge.End)
				queue = append(queue, newPath)
			}
		}
	}
	return paths
}

// Düğümleri yazdırma işlevi
func printNodes(nodes []Node) {
	fmt.Println("\nodalar:")
	for _, node := range nodes {
		fmt.Printf("%d: (%d, %d)\n", node.ID, node.X, node.Y)
	}
}

// Kenarları yazdırma işlevi
func printEdges(edges []Edge) {
	fmt.Println("\nkenarlar:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		AntCounts: make(map[int]int),
		AdjList:   make(map[int][]Edge),
	}
	antCount := 0
	startID, endID := -1, -1

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if antCount == 0 {
			antCount, err = strconv.Atoi(line)
			if err != nil {
				return
			}
			fmt.Println("karınca_sayısı:", antCount)
			continue
		}

		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			startID, _ = strconv.Atoi(fields[0])
			startx, _ := strconv.Atoi(fields[1])
			starty, _ := strconv.Atoi(fields[2])
			graph.StartNodeID = startID
			graph.Nodes = append(graph.Nodes, Node{ID: startID, X: startx, Y: starty})
			fmt.Println("başlangıç_odası:", graph.StartNodeID)
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			endID, _ = strconv.Atoi(fields[0])
			endx, _ := strconv.Atoi(fields[1])
			endy, _ := strconv.Atoi(fields[2])
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, X: endx, Y: endy})
			fmt.Println("bitiş_odası:", graph.EndNodeID)
		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				id, _ := strconv.Atoi(fields[0])
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				graph.Nodes = append(graph.Nodes, Node{ID: id, X: x, Y: y})
			} else if len(fields) == 1 && strings.Contains(line, "-") {
				edgeParts := strings.Split(fields[0], "-")
				start, _ := strconv.Atoi(edgeParts[0])
				end, _ := strconv.Atoi(edgeParts[1])
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Cost: 1}) // Kenar maliyeti varsayılan olarak 1
				graph.AdjList[start] = append(graph.AdjList[start], Edge{Start: start, End: end, Cost: 1})
				graph.AdjList[end] = append(graph.AdjList[end], Edge{Start: end, End: start, Cost: 1}) // Çift yönlü bağlantı
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatası:", err)
		return
	}

	if startID == -1 || endID == -1 {
		fmt.Println("Hata: Başlangıç veya bitiş düğümü belirtilmedi.")
		return
		}
		printNodes(graph.Nodes)
		printEdges(graph.Edges)
	}// BFS ile yolları bul
paths := graph.BFS(graph.StartNodeID, graph.EndNodeID)
fmt.Println("BFS ile bulunan yollar:", paths)

if len(paths) == 0 {
	fmt.Println("Yol bulunamadı.")
	return
}

// Yolları uzunluğa göre sırala, en kısa yolu almak için
sort.Slice(paths, func(i, j int) bool {
	return len(paths[i]) < len(paths[j])
})

// Karıncalara yolları ata
antPaths := make([][]int, antCount)
for i := 0; i < antCount; i++ {
	antPaths[i] = paths[i%len(paths)]
}

steps := make([][]string, 0)
occupied := make(map[int]bool)

// Karınca pozisyonlarını başlat
antPositions := make([]int, antCount)
for i := range antPositions {
	antPositions[i] = graph.StartNodeID
}

// Karıncaları taşı
for {
	moves := []string{}
	moved := false
	for i := 0; i < antCount; i++ {
		if antPositions[i] == graph.EndNodeID {
			continue
		}
		nextPositionIndex := -1
		for j, pos := range antPaths[i] {
			if pos == antPositions[i] {
				nextPositionIndex = j + 1
				break
			}
		}
		if nextPositionIndex >= len(antPaths[i]) {
			continue
		}
		nextPosition := antPaths[i][nextPositionIndex]
		if !occupied[nextPosition] {
			occupied[antPositions[i]] = false
			antPositions[i] = nextPosition
			occupied[nextPosition] = true
			moves = append(moves, fmt.Sprintf("L%d-%d", i+1, nextPosition))
			moved = true
		}
	}
	if !moved {
		break
	}
	steps = append(steps, moves)
}

// Adımları yazdır
for _, move := range steps {
	fmt.Println(strings.Join(move, " "))
}
