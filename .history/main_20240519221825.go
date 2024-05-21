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
	Cost  int // Edge cost
}

type Graph struct {
	Nodes       []Node
	Edges       []Edge
	StartNodeID int
	EndNodeID   int
	AntCounts   map[int]int
	AdjList     map[int][]Edge // Adjacency list
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

// BFS algorithm
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

// Function to print nodes
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: (%d, %d)\n", node.ID, node.X, node.Y)
	}
}

// Function to print edges
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// PriorityQueue items
type Item struct {
	nodeID   int
	priority int
	index    int
}

// PriorityQueue structure
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
		fmt.Println("ERROR: unable to open file")
		return
	}
	defer file.Close()

	graph := Graph{
		AntCounts: make(map[int]int),
		AdjList:   make(map[int][]Edge),
	}
	startID, endID := -1, -1

	scanner := bufio.NewScanner(file)
	antCount := -1

	// Helper function to check if a node ID already exists
	nodeExists := func(id int) bool {
		for _, node := range graph.Nodes {
			if node.ID == id {
				return true
			}
		}
		return false
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		// Check number of ants
		if antCount == -1 {
			antCount, err = strconv.Atoi(fields[0])
			if err != nil || antCount <= 0 {
				fmt.Println("ERROR: invalid data format, invalid number of ants")
				return
			}
			continue
		}

		if strings.HasPrefix(line, "##start") {
			if scanner.Scan() {
				fields = strings.Fields(scanner.Text())
				if len(fields) != 3 {
					fmt.Println("ERROR: invalid data format, invalid start room")
					return
				}
				startID, err = strconv.Atoi(fields[0])
				startX, err1 := strconv.Atoi(fields[1])
				startY, err2 := strconv.Atoi(fields[2])
				if err != nil || err1 != nil || err2 != nil || nodeExists(startID) {
					fmt.Println("ERROR: invalid data format, invalid start room")
					return
				}
				graph.StartNodeID = startID
				graph.Nodes = append(graph.Nodes, Node{ID: startID, X: startX, Y: startY})
			} else {
				fmt.Println("ERROR: invalid data format, missing start room")
				return
			}
			continue
		}

		if strings.HasPrefix(line, "##end") {
			if scanner.Scan() {
				fields = strings.Fields(scanner.Text())
				if len(fields) != 3 {
					fmt.Println("ERROR: invalid data format, invalid end room")
					return
				}
				endID, err = strconv.Atoi(fields[0])
				endX, err1 := strconv.Atoi(fields[1])
				endY, err2 := strconv.Atoi(fields[2])
				if err != nil || err1 != nil || err2 != nil || nodeExists(endID) {
					fmt.Println("ERROR: invalid data format, invalid end room")
					return
				}
				graph.EndNodeID = endID
				graph.Nodes = append(graph.Nodes, Node{ID: endID, X: endX, Y: endY})
			} else {
				fmt.Println("ERROR: invalid data format, missing end room")
				return
			}
			continue
		}

		if len(fields) == 3 {
			id, err1 := strconv.Atoi(fields[0])
			x, err2 := strconv.Atoi(fields[1])
			y, err3 := strconv.Atoi(fields[2])
			if err1 != nil || err2 != nil || err3 != nil || nodeExists(id) {
				fmt.Println("ERROR: invalid data format, invalid room")
				return
			}
			graph.Nodes = append(graph.Nodes, Node{ID: id, X: x, Y: y})
		} else if len(fields) == 1 && strings.Contains(fields[0], "-") {
			edgeParts := strings.Split(fields[0], "-")
			if len(edgeParts) != 2 {
				fmt.Println("ERROR: invalid data format, invalid edge")
				return
			}
			start, err1 := strconv.Atoi(edgeParts[0])
			end, err2 := strconv.Atoi(edgeParts[1])
			if err1 != nil || err2 != nil || !nodeExists(start) || !nodeExists(end) || start == end {
				fmt.Println("ERROR: invalid data format, invalid edge")
				return
			}
			graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Cost: 1})
			graph.AdjList[start] = append(graph.AdjList[start], Edge{Start: start, End: end, Cost: 1})
			graph.AdjList[end] = append(graph.AdjList[end], Edge{Start: end, End: start, Cost: 1})
		} else {
			fmt.Println("ERROR: invalid data format, unknown line format")
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ERROR: reading file", err)
		return
	}

	if startID == -1 || endID == -1 {
		fmt.Println("ERROR: invalid data format, start or end node not specified")
		return
	}

	printNodes(graph.Nodes)
	printEdges(graph.Edges)

	// Find paths from start to end using BFS
	paths := graph.BFS(graph.StartNodeID, graph.EndNodeID)
	if len(paths) == 0 {
		fmt.Println("ERROR: no path found from start to end")
		return
	}

	fmt.Println("Paths found with BFS:", paths)

	// Assign paths to ants
	antPaths := make([][]int, antCount)
	for i := 0; i < antCount; i++ {
		antPaths[i] = paths[i%len(paths)]
	}

	steps := make([][]string, 0)
	occupied := make(map[int]bool)

	// Initialize positions of ants
	antPositions := make([]int, antCount)
	for i := range antPositions {
		antPositions[i] = graph.StartNodeID
	}

	// Move ants
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

	// Print moves
	for _, move := range steps {
		fmt.Println(strings.Join(move, " "))
	}
}

func (g *Graph) CheckAntCount(nodeID int) int {
	antCount, exist := g.AntCounts[nodeID]
	if !exist {
		return 0
	}
	return antCount
}

func (g *Graph) UpdateAntCounts(newAntCounts map[int]int) {
	for id, count := range newAntCounts {
		g.AntCounts[id] = count
	}
}
