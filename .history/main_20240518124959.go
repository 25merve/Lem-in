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

	// Find paths with BFS
	paths := graph.BFS(graph.StartNodeID, graph.EndNodeID)
	fmt.Println("Paths found with BFS:", paths)

	if len(paths) == 0 {
		fmt.Println("No path found.")
		return
	}

	// Sort paths by length to ensure we get the shortest paths first
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Function to print moves
	printMoves := func(steps [][]string) {
		for _, move := range steps {
			fmt.Println(strings.Join(move, " "))
		}
	}

	steps := [][]string{}
	occupied := make(map[int]bool)
	antPositions := make([]int, antCount)
	for i := range antPositions {
		antPositions[i] = graph.StartNodeID
	}

	for {
		moves := []string{}
		moved := false

		for i := 0; i < antCount; i++ {
			if antPositions[i] == graph.EndNodeID {
				continue
			}

			nextPositionIndex := -1
			for j, pos := range paths[i%len(paths)] {
				if pos == antPositions[i] {
					nextPositionIndex = j + 1
					break
				}
			}

			if nextPositionIndex >= len(paths[i%len(paths)]) {
				continue
			}

			nextPosition := paths[i%len(paths)][nextPositionIndex]
			if !occupied[nextPosition] {
				occupied[antPositions[i]] = false
				antPositions[i] = nextPosition
				occupied[nextPosition] = true
				moves = append(moves, fmt.Sprintf("L%d-%d", i+1, nextPosition))
				moved = true
			} else {
				// If the next position is occupied, find an alternative path for the ant
				for _, path := range paths {
					if !contains(path, nextPosition) {
						for j, pos := range path {
							if pos == antPositions[i] {
								nextPositionIndex = j + 1
								break
							}
						}
						if nextPositionIndex < len(path) {
							alternativePosition := path[nextPositionIndex]
							if !occupied[alternativePosition] {
								occupied[antPositions[i]] = false
								antPositions[i] = alternativePosition
								occupied[alternativePosition] = true
								moves = append(moves, fmt.Sprintf("L%d-%d", i+1, alternativePosition))
								moved = true
								break
							}
						}
					}
				}
			}
		}

		if !moved {
			break
		}

		steps = append(steps, moves)
	}

	printMoves(steps)
}

// Helper function to check if a slice contains an element
func contains(slice []int, elem int) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}
