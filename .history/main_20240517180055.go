package main

import (
	"bufio"
	"container/heap"
	"fmt"
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

func (g *Graph) Dijkstra(startNodeID int) map[int]int {
	dist := make(map[int]int)
	for _, node := range g.Nodes {
		dist[node.ID] = int(^uint(0) >> 1) // Max int
	}
	dist[startNodeID] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{nodeID: startNodeID, priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		current := item.nodeID

		for _, edge := range g.AdjList[current] {
			next := edge.End
			newDist := dist[current] + edge.Cost
			if newDist < dist[next] {
				dist[next] = newDist
				heap.Push(pq, &Item{nodeID: next, priority: newDist})
			}
		}
	}
	return dist
}

// BFS Algoritması
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

// Düğümleri yazdırmak için fonksiyon
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: (%d, %d)\n", node.ID, node.X, node.Y)
	}
}

// Kenarları yazdırmak için fonksiyon
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// PriorityQueue item'ları
type Item struct {
	nodeID   int
	priority int
	index    int
}

// PriorityQueue yapısı
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
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

	scanner := bufio.NewScanner(file)
	section := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			section = strings.TrimSpace(line[1:])
			continue
		}
		if section == "the_ants" {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				nodeID, _ := strconv.Atoi(parts[0])
				count, _ := strconv.Atoi(parts[1])
				graph.AntCounts[nodeID] = count
			}
		} else if section == "the_rooms" {
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				id, _ := strconv.Atoi(parts[0])
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graph.Nodes = append(graph.Nodes, Node{ID: id, X: x, Y: y})
			}
		} else if section == "the_links" {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				start, _ := strconv.Atoi(parts[0])
				end, _ := strconv.Atoi(parts[1])
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Cost: 1})
				graph.AdjList[start] = append(graph.AdjList[start], Edge{Start: start, End: end, Cost: 1})
				graph.AdjList[end] = append(graph.AdjList[end], Edge{Start: end, End: start, Cost: 1})
			}
		} else if section == "start" {
			graph.StartNodeID, _ = strconv.Atoi(line)
		} else if section == "end" {
			graph.EndNodeID, _ = strconv.Atoi(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okuma hatası:", err)
		return
	}

	// Print nodes and edges
	printNodes(graph.Nodes)
	printEdges(graph.Edges)

	// Find all paths using BFS
	paths := graph.BFS(graph.StartNodeID, graph.EndNodeID)

	// Print paths
	fmt.Println("\nPaths found:")
	for _, path := range paths {
		fmt.Println(path)
	}

	// Ant movement simulation
	moveAnts(graph, paths)
}

// Function to move ants along found paths
func moveAnts(graph Graph, paths [][]int) {
	antID := 1
	antPositions := make(map[int]int) // Map of antID to current position index in the path
	antPaths := make(map[int][]int)   // Map of antID to path

	// Assign paths to ants
	for antID <= graph.AntCounts[graph.StartNodeID] {
		for _, path := range paths {
			if antID > graph.AntCounts[graph.StartNodeID] {
				break
			}
			antPaths[antID] = path
			antPositions[antID] = 0
			antID++
		}
	}

	// Move ants step by step
	step := 1
	for {
		fmt.Printf("\nStep %d:\n", step)
		moved := false
		for antID, position := range antPositions {
			if position < len(antPaths[antID])-1 {
				moved = true
				antPositions[antID]++
				fmt.Printf("Ant %d moves to node %d\n", antID, antPaths[antID][antPositions[antID]])
			}
		}
		if !moved {
			break
		}
		step++
	}
}
