package main

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
	

	file, err := os.Open("")
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
