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

type Node struct {
	ID   int    // Node ID
	Name string // Node Name
	X    int    // X coordinate
	Y    int    // Y coordinate
}

type Edge struct {
	Start int // Starting node ID of the edge
	End   int // Ending node ID of the edge
}

type Graph struct {
	Nodes       []Node        // List of nodes in the graph
	Edges       []Edge        // List of edges in the graph
	StartNodeID int           // ID of the start node
	EndNodeID   int           // ID of the end node
	AdjList     map[int][]int // Adjacency list representing the graph
}

// Function to print all nodes
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: %s (%d, %d)\n", node.ID, node.Name, node.X, node.Y)
	}
}

// Function to print all edges
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// Function to find all paths from startNodeID to endNodeID using BFS
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

// Helper function to check if a slice contains a particular item
func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Function to find the node ID by its name
func findNodeIDByName(nodes []Node, name string) int {
	for _, node := range nodes {
		if node.Name == name {
			return node.ID
		}
	}
	return -1
}

func main() {
	startTime := time.Now() // Başlangıç zamanını al
	endTime := time.Now().Add(2 * time.Minute)
	if len(os.Args) < 2 {
		fmt.Println("Dosya adı belirtilmedi.")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		AdjList: make(map[int][]int),
	}

	scanner := bufio.NewScanner(file)

	// Read the number of ants
	scanner.Scan()
	antCountLine := scanner.Text()
	antCount, err := strconv.Atoi(antCountLine)
	if err != nil {
		fmt.Println("Karınca sayısı okunamadı:", err)
		return
	}
	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format:karınca sayısı 0 veya daha kücük olamaz.")
		return
	} else {
		// Print the number of ants
		fmt.Printf("Number of ants: %d\n", antCount)
	}

	// Read the graph data
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			name := fields[0]
			x, _ := strconv.Atoi(fields[1])
			y, _ := strconv.Atoi(fields[2])
			startID := len(graph.Nodes)
			graph.StartNodeID = startID
			graph.Nodes = append(graph.Nodes, Node{ID: startID, Name: name, X: x, Y: y})
			fmt.Println("start_room:", name) // Print start room

		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			name := fields[0]
			x, _ := strconv.Atoi(fields[1])
			y, _ := strconv.Atoi(fields[2])
			endID := len(graph.Nodes)
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, Name: name, X: x, Y: y})
			fmt.Println("end_room:", name) // Print end room

		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				name := fields[0]
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				id := len(graph.Nodes)
				graph.Nodes = append(graph.Nodes, Node{ID: id, Name: name, X: x, Y: y})
			} else if len(fields) == 1 && strings.Contains(line, "-") {
				edgeParts := strings.Split(fields[0], "-")
				startName := edgeParts[0]
				endName := edgeParts[1]
				startID := findNodeIDByName(graph.Nodes, startName)
				endID := findNodeIDByName(graph.Nodes, endName)
				graph.Edges = append(graph.Edges, Edge{Start: startID, End: endID})
				graph.AdjList[startID] = append(graph.AdjList[startID], endID)
				graph.AdjList[endID] = append(graph.AdjList[endID], startID)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatası:", err) // Reading error
		return
	}

	if graph.StartNodeID == -1 || graph.EndNodeID == -1 {
		fmt.Println("Hata: Başlangıç veya bitiş düğümü belirtilmemiş.") // Error: Start or end node not specified
		return
	}

	printNodes(graph.Nodes) // Print all nodes
	printEdges(graph.Edges) // Print all edges

	allPaths := graph.BFSAllPaths(graph.StartNodeID, graph.EndNodeID)

	// Sort paths by length (from shortest to longest)
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Filter paths to ensure no overlapping
	filteredPaths := filterPaths(allPaths, antCount)

	if len(filteredPaths) == 0 || len(filteredPaths[0]) == 0 {
		fmt.Println("ERROR: invalid data format:Hata: Geçerli bir yol yok.") // Error: No valid path
		return
	} else {
		fmt.Println("Filtrelenmiş yollar:") // Filtered paths

		for _, path := range filteredPaths {
			for i, node := range path {
				if i > 0 {
					fmt.Print(" -> ")
				}
				fmt.Print(graph.Nodes[node].Name)
			}
			fmt.Println()
		}

		// Assign paths to ants
		antPaths := assignPathsToAnts(antCount, filteredPaths)

		// Initialize ant positions
		antPositions := make([]int, antCount)
		for i := 0; i < antCount; i++ {
			antPositions[i] = graph.StartNodeID
		}

		step := 1
		for {
			moves := []string{}
			allAtEnd := true
			occupied := make(map[int]bool)

			// Collect the planned moves for this step
			for i := 0; i < antCount; i++ {
				if antPositions[i] != graph.EndNodeID {
					allAtEnd = false
					path := antPaths[i]
					for j := 0; j < len(path)-1; j++ {
						if path[j] == antPositions[i] && (!occupied[path[j+1]] || path[j+1] == graph.EndNodeID) {
							occupied[path[j+1]] = true
							antPositions[i] = path[j+1]
							moves = append(moves, fmt.Sprintf("L%d-%s", i+1, graph.Nodes[antPositions[i]].Name))
							break
						}
					}
				}
			}

			// Print the moves for this step
			if len(moves) > 0 {
				fmt.Printf("Adım %d: %s\n", step, strings.Join(moves, " "))
			}

			// If all ants are at the end, break the loop
			if allAtEnd {
				break
			}

			step++

			// Check if the last step was redundant
			allAtEnd = true
			for i := 0; i < antCount; i++ {
				if antPositions[i] != graph.EndNodeID {
					allAtEnd = false
					break
				}
			}

		}
		// startTime ve endTime arasındaki farkı hesaplayın
		geçensure := endTime.Sub(startTime)

		// Dakika cinsine dönüştürün
		geçensureDakika := geçensure.Minutes()

		fmt.Printf("Toplam süre: %.1f dakika\n", geçensureDakika)
	}
}

// Function to assign paths to ants ensuring no overlap and choosing the shortest path for the first ant
func assignPathsToAnts(antCount int, filteredPaths [][]int) [][]int {
	antPaths := make([][]int, antCount)
	remainingPaths := filteredPaths

	// Assign paths to all ants
	for i := 0; i < antCount; i++ {
		if len(remainingPaths) > 0 {
			antPaths[i] = remainingPaths[0]
			remainingPaths = remainingPaths[1:]
		} else {
			// If no more unique paths, reuse available paths
			antPaths[i] = filteredPaths[i%len(filteredPaths)]
		}
	}

	return antPaths
}

// Function to filter paths to ensure they don't overlap more than allowed by ant count
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

// Function to check if two paths have overlapping nodes (except start and end)
func hasOverlap(path1 []int, path2 []int) bool {
	for i := 1; i < len(path1)-1; i++ {
		for j := 1; j < len(path2)-1; j++ {
			if path1[i] == path2[j] {
				return true
			}
		}
	}
	return false
}
}