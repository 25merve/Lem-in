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
	ID   int    // Düğüm ID'si
	Name string // Düğüm Adı
	X    int    // X koordinatı
	Y    int    // Y koordinatı
}

type Edge struct {
	Start int // Kenarın başlangıç düğümü ID'si
	End   int // Kenarın bitiş düğümü ID'si
}

type Graph struct {
	Nodes       []Node        // Grafikteki düğümlerin listesi
	Edges       []Edge        // Grafikteki kenarların listesi
	StartNodeID int           // Başlangıç düğümünün ID'si
	EndNodeID   int           // Bitiş düğümünün ID'si
	AdjList     map[int][]int // Grafiği temsil eden bitişiklik listesi
}

// Tüm düğümleri yazdıran fonksiyon
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: %s (%d, %d)\n", node.ID, node.Name, node.X, node.Y)
	}
}

// Tüm kenarları yazdıran fonksiyon
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// BFS kullanarak startNodeID'den endNodeID'ye tüm yolları bulan fonksiyon
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
				newPath := append([]int{}, path...) //Bu boş dilime, path dilimindeki tüm elemanları ekler.
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return paths
}

// Bir dilimde 
func contains(slice []int, item int) bool {
	for _, v := range slice {//_ burada indeks değerini alır, ancak kullanılmadığı için _ olarak belirtilir.
		if v == item {
			return true
		}
	}
	return false
}

// Düğüm adından düğüm ID'sini bulan fonksiyon
func findNodeIDByName(nodes []Node, name string) int {
	for _, node := range nodes {
		if node.Name == name {
			return node.ID
		}
	}
	return -1
}

// Eğer ana yol tıkanırsa bir karınca için alternatif bir yol bulan fonksiyon
func findAlternativePath(graph Graph, currentPos int, occupied map[int]bool) []int {
	allPaths := graph.BFSAllPaths(currentPos, graph.EndNodeID)
	for _, path := range allPaths {
		valid := true
		for _, node := range path {
			if occupied[node] {
				valid = false
				break
			}
		}
		if valid {
			return path
		}
	}
	return nil
}

// Çakışma olmadığından emin olarak yolları filtreleyen yeni fonksiyon
func FilterPaths(paths [][]int) [][]int {
	var maxPaths [][]int
	var currentPaths [][]int
	usedNodes := make(map[int]bool)

	var backtrack func(int)
	backtrack = func(start int) {
		if len(currentPaths) > len(maxPaths) {
			maxPaths = make([][]int, len(currentPaths))
			copy(maxPaths, currentPaths)
		}

		for i := start; i < len(paths); i++ {
			path := paths[i]
			keepPath := true

			for _, node := range path[1 : len(path)-1] {
				if usedNodes[node] {
					keepPath = false
					break
				}
			}

			if keepPath {
				currentPaths = append(currentPaths, path)
				for _, node := range path[1 : len(path)-1] {
					usedNodes[node] = true
				}

				backtrack(i + 1)

				// Geriye git
				currentPaths = currentPaths[:len(currentPaths)-1]
				for _, node := range path[1 : len(path)-1] {
					delete(usedNodes, node)
				}
			}
		}
	}

	backtrack(0)
	return maxPaths
}

func main() {
	startTime := time.Now() // Başlangıç zamanını al

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
		AdjList:     make(map[int][]int),
		StartNodeID: -1,
		EndNodeID:   -1,
	}

	scanner := bufio.NewScanner(file)

	// Karınca sayısını oku
	scanner.Scan()
	antCountLine := scanner.Text()
	antCount, err := strconv.Atoi(antCountLine)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}
	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Graf verilerini oku
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			if len(fields) < 3 {
				fmt.Println("ERROR: invalid data format")
				return
			}
			name := fields[0]
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("ERROR: invalid data format")
				return
			}
			y, err := strconv.Atoi(fields[2])
			if err != nil {
				fmt.Println("ERROR: invalid data format")
				return
			}
			startID := len(graph.Nodes)
			graph.StartNodeID = startID
			graph.Nodes = append(graph.Nodes, Node{ID: startID, Name: name, X: x, Y: y})

		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			if len(fields) < 3 {
				fmt.Println("ERROR: invalid data format")
				return
			}
			name := fields[0]
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("ERROR: invalid data format")
				return
			}
			y, err := strconv.Atoi(fields[2])
			if err != nil {
				fmt.Println("ERROR: invalid data format")
				return
			}
			endID := len(graph.Nodes)
			graph.EndNodeID = endID
			graph.Nodes = append(graph.Nodes, Node{ID: endID, Name: name, X: x, Y: y})

		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				name := fields[0]
				x, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("ERROR: invalid data format")
					return
				}
				y, err := strconv.Atoi(fields[2])
				if err != nil {
					fmt.Println("ERROR: invalid data format")
					return
				}
				id := len(graph.Nodes)
				graph.Nodes = append(graph.Nodes, Node{ID: id, Name: name, X: x, Y: y})
			} else if len(fields) == 1 && strings.Contains(line, "-") {
				edgeParts := strings.Split(fields[0], "-")
				if len(edgeParts) != 2 {
					fmt.Println("ERROR: invalid data format")
					return
				}
				startName := edgeParts[0]
				endName := edgeParts[1]
				startID := findNodeIDByName(graph.Nodes, startName)
				endID := findNodeIDByName(graph.Nodes, endName)
				if startID == -1 || endID == -1 {
					fmt.Println("ERROR: invalid data format")
					return
				}
				graph.Edges = append(graph.Edges, Edge{Start: startID, End: endID})
				graph.AdjList[startID] = append(graph.AdjList[startID], endID)
				graph.AdjList[endID] = append(graph.AdjList[endID], startID)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ERROR: invalid data format") // Okuma hatası
		return
	}

	if graph.StartNodeID == -1 || graph.EndNodeID == -1 {
		fmt.Println("ERROR: invalid data format") // Hata: Başlangıç veya bitiş düğümü belirtilmedi
		return
	}

	allPaths := graph.BFSAllPaths(graph.StartNodeID, graph.EndNodeID)
	if len(allPaths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	filteredPaths := FilterPaths(allPaths)

	if len(filteredPaths) == 0 || len(filteredPaths[0]) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Giriş verilerini yazdır
	fmt.Printf("Number of ants: %d\n", antCount)
	fmt.Printf("start_room: %d\n", graph.StartNodeID)
	fmt.Printf("end_room: %d\n", graph.EndNodeID)
	printNodes(graph.Nodes)
	printEdges(graph.Edges)

	antPaths := assignPathsToAnts(antCount, filteredPaths)
	antPositions := make([]int, antCount)
	antAtEnd := make([]bool, antCount)

	// Tüm karıncalar için pozisyonları ve yolları başlat
	for i := 0; i < antCount; i++ {
		if i == antCount-1 {
			antPaths[i] = filteredPaths[0] // Son karınca en kısa yolu takip eder
		} else {
			antPaths[i] = filteredPaths[(i)%len(filteredPaths)] // Diğer karıncalar sırayla takip eder
		}
		antPositions[i] = graph.StartNodeID // Bütün karıncalar başlangıç pozisyonundadır
		antAtEnd[i] = false                 // Hiçbir karınca bitişe ulaşmamıştır
	}

	step := 1
	for {
		moves := []string{}
		allAtEnd := true
		occupied := make(map[int]bool)

		// Bu adım için planlanan hareketleri topla
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

				// Tıkanmışsa alternatif yol kontrol et
				if antPositions[i] != graph.EndNodeID && occupied[antPositions[i]] {
					altPath := findAlternativePath(graph, antPositions[i], occupied)
					if altPath != nil {
						antPaths[i] = altPath
					}
				}
			} else {
				antAtEnd[i] = true
			}
		}

		// Bu adım için hareketleri yazdır
		if len(moves) > 0 {
			fmt.Printf("Adım %d: %s\n", step, strings.Join(moves, " "))
		}

		// Tüm karıncalar bitişte ise döngüyü kır
		if allAtEnd {
			break
		}

		step++

		// Son adımın gereksiz olup olmadığını kontrol et
		allAtEnd = true
		for i := 0; i < antCount; i++ {
			if !antAtEnd[i] {
				allAtEnd = false
				break
			}
		}
	}

	// Geçen zamanı hesapla
	elapsedTime := time.Since(startTime)

	// Toplam geçen süreyi kesirli kısmı ile birlikte yazdır
	fmt.Printf("Toplam süre: %.9f saniye\n", elapsedTime.Seconds())
}

// Karıncalara yolları atayan ve çakışmayı önleyen fonksiyon
func assignPathsToAnts(antCount int, filteredPaths [][]int) [][]int {
	antPaths := make([][]int, antCount)
	remainingPaths := filteredPaths

	// Tüm karıncalara yolları ata
	for i := 0; i < antCount; i++ {
		if len(remainingPaths) > 0 {
			antPaths[i] = remainingPaths[0]
			remainingPaths = remainingPaths[1:]
		} else {
			// Daha fazla benzersiz yol yoksa, mevcut yolları yeniden kullan
			antPaths[i] = filteredPaths[i%len(filteredPaths)]
		}
	}

	return antPaths
}
