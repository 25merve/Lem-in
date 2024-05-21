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
	ID   int
	Name string
	X    int
	Y    int
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
		fmt.Printf("%s: (%d, %d)\n", node.Name, node.X, node.Y)
	}
}

// Function to print edges
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

func validateInput(lines []string) error {
	// Perform input validation according to the specified rules
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: No input file specified.")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: File opening error:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		AntCounts: make(map[int]int),
		AdjList:   make(map[int][]Edge),
	}
	antCount := 0
	startID, endID := -1, -1
	lines := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := validateInput(lines); err != nil {
		fmt.Println(err)
		return
	}

	lineIndex := 0
	for lineIndex < len(lines) {
		line := lines[lineIndex]
		if antCount == 0 {
			antCount, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("ERROR: invalid data format, invalid number of Ants")
				return
			}
			fmt.Println("number_of_ants:", antCount)
			lineIndex++
			continue
		}

		if strings.HasPrefix(line, "##start") {
			lineIndex++
			fields := strings.Fields(lines[lineIndex])
			startID = lineIndex
			startName := fields[0]
			startx, _ := strconv.Atoi(fields[1])
			starty, _ := strconv.Atoi(fields[2])
			graph.StartNodeID = startID
			graph.Nodes = append(graph.Nodes, Node{ID: startID, Name: startName, X: startx, Y: starty})
			fmt.Println("start_room:", graph.StartNodeID)
		} else if strings.HasPrefix(line, "##end") {
			lineIndex++
			fields := strings.Fields(lines[lineIndex])
			endID = lineIndex
			endName := fields[0]
			endx, _ := strconv.Atoi(fields[1])
			endy, _ := strconv.Atoi(fields[2])
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, Name: endName, X: endx, Y: endy})
			fmt.Println("end_room:", graph.EndNodeID)
		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				id := lineIndex
				name := fields[0]
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				graph.Nodes = append(graph.Nodes, Node{ID: id, Name: name, X: x, Y: y})
			} else if len(fields) == 1 && strings.Contains(line, "-") {
				edgeParts := strings.Split(fields[0], "-")
				start, end := -1, -1
				for _, node := range graph.Nodes {
					if node.Name == edgeParts[0] {
						start = node.ID
					}
					if node.Name == edgeParts[1] {
						end = node.ID
					}
				}
				if start == -1 || end == -1 {
					fmt.Println("ERROR: invalid data format, unknown room name in link")
					return
				}
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end, Cost: 1}) // Default edge cost is 1
				graph.AdjList[start] = append(graph.AdjList[start], Edge{Start: start, End: end, Cost: 1})
				graph.AdjList[end] = append(graph.AdjList[end], Edge{Start: end, End: start, Cost: 1}) // Bidirectional link
			}
		}
		lineIndex++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ERROR: Reading error:", err)
		return
	}

	if startID == -1 || endID == -1 {
		fmt.Println("ERROR: invalid data format, no start or end room found.")
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
				moves = append(moves, fmt.Sprintf("L%d-%s", i+1, graph.Nodes[nextPosition].Name))
				moved = true
			}
		}

		if !moved {
			break
		}

		steps = append(steps, moves)
		for k := range occupied {
			delete(occupied, k)
		}
	}

	// Print moves
	for _, move := range steps {
		fmt.Println(strings.Join(move, " "))
	}
}
