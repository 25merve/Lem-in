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
		for k := range occupied {
			delete(occupied, k)
		}
	}

	// Print moves
	for _, move := range steps {package main

		import (
			"bufio"
			"fmt"
			"os"
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
		}
		
		type Graph struct {
			Nodes       []Node
			Edges       []Edge
			StartNodeID int
			EndNodeID   int
			AdjList     map[int][]int
		}
		
		func printNodes(nodes []Node) {
			fmt.Println("\nthe_rooms:")
			for _, node := range nodes {
				fmt.Printf("%d: %s (%d, %d)\n", node.ID, node.Name, node.X, node.Y)
			}
		}
		
		func printEdges(edges []Edge) {
			fmt.Println("\nthe_links:")
			for _, edge := range edges {
				fmt.Printf("%d - %d\n", edge.Start, edge.End)
			}
		}
		
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
		
		func contains(slice []int, item int) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		}
		
		func findNodeIDByName(nodes []Node, name string) int {
			for _, node := range nodes {
				if node.Name == name {
					return node.ID
				}
			}
			return -1
		}
		
		func main() {
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
		
			scanner.Scan()
			antCountLine := scanner.Text()
			antCount, err := strconv.Atoi(antCountLine)
			if err != nil {
				fmt.Println("Karınca sayısı okunamadı:", err)
				return
			}
		
			startID, endID := -1, -1
		
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "##start") {
					scanner.Scan()
					fields := strings.Fields(scanner.Text())
					name := fields[0]
					x, _ := strconv.Atoi(fields[1])
					y, _ := strconv.Atoi(fields[2])
					startID = len(graph.Nodes)
					graph.StartNodeID = startID
					graph.Nodes = append(graph.Nodes, Node{ID: startID, Name: name, X: x, Y: y})
					fmt.Println("start_room:", name)
		
				} else if strings.HasPrefix(line, "##end") {
					scanner.Scan()
					fields := strings.Fields(scanner.Text())
					name := fields[0]
					x, _ := strconv.Atoi(fields[1])
					y, _ := strconv.Atoi(fields[2])
					endID = len(graph.Nodes)
					graph.EndNodeID = endID
					graph.Nodes = append(graph.Nodes, Node{ID: endID, Name: name, X: x, Y: y})
					fmt.Println("end_room:", name)
		
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
				fmt.Println("Okuma hatası:", err)
				return
			}
		
			if startID == -1 || endID == -1 {
				fmt.Println("Hata: Başlangıç veya bitiş düğümü belirtilmemiş.")
				return
			}
		
			printNodes(graph.Nodes)
			printEdges(graph.Edges)
		
			allPaths := graph.BFSAllPaths(graph.StartNodeID, graph.EndNodeID)
			if len(allPaths) == 0 {
				fmt.Println("Hata: Başlangıç ve bitiş arasında geçerli bir yol yok.")
				return
			}
			fmt.Println("BFS ile bulunan tüm yollar:", allPaths)
		
			antPaths := make([][]int, antCount)
			for i := 0; i < antCount; i++ {
				antPaths[i] = allPaths[i%len(allPaths)]
			}
		
			antPositions := make([]int, antCount)
			for i := 0; i < antCount; i++ {
				antPositions[i] = graph.StartNodeID
			}
		
			step := 1
			for {
				moves := []string{}
				allAtEnd := true
				occupied := make(map[int]bool)
				for i := 0; i < antCount; i++ {
					if antPositions[i] != graph.EndNodeID {
						allAtEnd = false
						path := antPaths[i]
						for j := 0; j < len(path)-1; j++ {
							if path[j] == antPositions[i] && !occupied[path[j+1]] {
								occupied[path[j+1]] = true
								antPositions[i] = path[j+1]
								moves = append(moves, fmt.Sprintf("L%d-%s", i+1, graph.Nodes[antPositions[i]].Name))
								break
							}
						}
					}
				}
				if allAtEnd {
					break
				}
				if len(moves) > 0 {
					fmt.Printf("Adım %d: %s\n", step, strings.Join(moves, " "))
				}
				step++
			}
		}
		fmt.Println(strings.Join(move, " "))
	}
}
