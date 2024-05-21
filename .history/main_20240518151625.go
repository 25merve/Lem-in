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

func (g *Graph) DFS(startNodeID int, endNodeID int, visited map[int]bool, path []int, allPaths *[][]int) {
	if startNodeID == endNodeID {
		newPath := make([]int, len(path))
		copy(newPath, path)
		*allPaths = append(*allPaths, newPath)
		return
	}

	visited[startNodeID] = true

	for _, edge := range g.AdjList[startNodeID] {
		if !visited[edge.End] {
			path = append(path, edge.End)
			g.DFS(edge.End, endNodeID, visited, path, allPaths)
			path = path[:len(path)-1]
		}
	}

	visited[startNodeID] = false
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

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("File opening error:", err)
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
			fmt.Println("number_of_ants:", antCount)
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
			fmt.Println("start_room:", graph.StartNodeID)
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			endID, _ = strconv.Atoi(fields[0])
			endx, _ := strconv.Atoi(fields[1])
			endy, _ := strconv.Atoi(fields[2])
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, X: endx, Y: endy})
			fmt.Println("end_room:", graph.EndNodeID)
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
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Cost: 1}) // Default edge cost is 1
				graph.AdjList[start] = append(graph.AdjList[start], Edge{Start: start, End: end, Cost: 1})
				graph.AdjList[end] = append(graph.AdjList[end], Edge{Start: end, End: start, Cost: 1}) // Bidirectional link
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Reading error:", err)
		return
	}

	if startID == -1 || endID == -1 {
		fmt.Println("Error: Start or end node not specified.")
		return
	}
	printNodes(graph.Nodes)
	printEdges(graph.Edges)

	// Find paths with DFS
	allPaths := [][]int{}
	visited := make(map[int]bool)
	graph.DFS(graph.StartNodeID, graph.EndNodeID, visited, []int{graph.StartNodeID}, &allPaths)
	fmt.Println("Paths found with DFS:", allPaths)

	if len(allPaths) == 0 {
		fmt.Println("No path found.")
		return
	}

	// Sort paths by length to ensure we get the shortest paths first
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Assign paths to ants in a round-robin manner
	antPaths := make([][]int, antCount)
	for i := 0; i < antCount; i++ {
		antPaths[i] = allPaths[i%len(allPaths)]
	}

	steps := make([][]string, 0)
	occupied := make(map[int]bool)
	antPositions := make([]int, antCount)
	for i := range antPositions {
		antPositions[i] = graph.StartNodeID
	}

	// Move ants
	for {
		moves := []string{} // Karınca hareketlerini tutmak için boş bir dizi oluşturun
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

		// Karınca hareketlerini diziye ekle
		steps = append(steps, moves)
		// Mevcut konumları sıfırla
		for k := range occupied {
			delete(occupied, k)
		}
	}

	// Print moves
	for _, move := range steps {
		fmt.Println(strings.Join(move, " "))
	}

}
