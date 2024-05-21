package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct { //dügümler için struct yapısı
	ID int
	X  int
	Y  int
}

type Edge struct { //kenarlar için struct yapısı
	Start int
	End   int
}

type Graph struct { //düğümler ve knarlar için oluşturulmuş bir graf struct yapısı
	Nodes []Node
	Edges []Edge
}

func main() {
	file, err := os.Open("text.txt") //dosya açma işlemi
	if err != nil {
		fmt.Println("Dosya açma hatasi:", err)
		return
	}
	defer file.Close()

	startID, endID := -1, -1 //ilk önce baslangic ve bitis satırlarına bi id atıyoruz
	graph := Graph{}         //struct yapısını değişkene atıyoruz

	scanner := bufio.NewScanner(file)
	var antCount int
	for scanner.Scan() {
		line := scanner.Text()
		if antCount == 0 { // Karınca sayısı daha önce belirlenmediyse
			antCount, err = strconv.Atoi(line) // Dosyadan karınca sayısını oku
			if err != nil {
				return
			}
			fmt.Println("karınca sayısı :", antCount)
		}
		if strings.HasPrefix(line, "##start") { //içinde start var mı diye kontrol ediyor
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			startID, _ = strconv.Atoi(fields[0]) //ilk önce atadığımız id değiştiği için onun atamasını yapıyoruz
		} else if strings.HasPrefix(line, "##end") { //satır bitis satırı mı diye kontrol ediyor
			scanner.Scan()
			fields := strings.Fields(scanner.Text())
			endID, _ = strconv.Atoi(fields[0]) //ilk önce atadığımız id değiştiği için onun atamasını yapıyoruz
		} else {
			fields := strings.Fields(line)
			if len(fields) == 3 { //eger satır 3 argumandan oluşuyorsa o zaman  bunun bir dügüm koordinatı belirtme olduğunu anlayıp onaa göre atama yapıyoruz
				id, _ := strconv.Atoi(fields[0])
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				graph.Nodes = append(graph.Nodes, Node{ID: id, X: x, Y: y}) //grafdaki dügüm struct yapısına ekleme yapıyoruz çünkü artık id,x,y yapılarımız belirlendi
			} else if len(fields) == 1 {
				// Eğer satır sadece bir sayıdan oluşuyorsa, bu karınca sayısıdır.
				antCount, _ = strconv.Atoi(fields[0])
			} else if len(fields) == 2 {
				edgeParts := strings.Split(fields[0], "-") //- ile ayrılmış olanları alıyoruz
				start, _ := strconv.Atoi(edgeParts[0])     //bunlardan ilki bizim için başlangıç ikincisi bitiş olaark alınıyor çünkü aradaki uzunluğu bulmak için kullanıcaz
				end, _ := strconv.Atoi(edgeParts[1])
				graph.Edges = append(graph.Edges, Edge{Start: start, End: end}) //kenarlar struct  yapısına ekleme yapıyoruz çünkü artık başlangıç ve bitiş noktaları belirlendi
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatasi:", err)
		return
	}

	// Başlangıç ve bitiş düğümlerini kontrol et
	if startID == -1 || endID == -1 {
		fmt.Println("Hata: Başlangiç veya bitiş düğümü belirtilmemiş.")
		return
	}

	// Düğümleri ve kenarları yazdır
	fmt.Println("Düğümler:")
	for _, node := range graph.Nodes {
		fmt.Printf("%d: (%d, %d)\n", node.ID, node.X, node.Y)

	}

	fmt.Println("\nKenarlar:")
	for _, edge := range graph.Edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}

}

// ID'ye göre düğümü bul
func findNodeByID(nodes []Node, id int) Node {
	for _, node := range nodes {
		if node.ID == id { //eger dügüm id bilgisi aranan id bilgisine eşitse id bu deyip döndürme işlemi yapıyor
			return node
		}
	}
	return Node{} // Bulunamadıysa boş düğüm döndür
}
func shortestPathBFS(graph Graph, startID, targetID int) []int {
	queue := make([]int, 0)       // BFS için bir kuyruk oluştur
	visited := make(map[int]bool) // Ziyaret edilen düğümleri tutmak için bir harita
	parents := make(map[int]int)  // Düğümlerin ebeveynlerini tutmak için bir harita
	found := false                // Hedef düğüm bulundu mu?

	queue = append(queue, startID) // Başlangıç düğümünü kuyruğa ekle
	visited[startID] = true        // Başlangıç düğümünü ziyaret et

	for len(queue) > 0 {
		currentNode := queue[0] // Kuyruğun başındaki düğümü al
		queue = queue[1:]       // Kuyruğun başındaki düğümü kaldır

		if currentNode == targetID { // Hedef düğümü bulundu mu?
			found = true
			break
		}

		// Şimdi, currentDüğümün komşularını bul ve kuyruğa ekle
		for _, edge := range graph.Edges {
			if edge.Start == currentNode {
				neighborID := edge.End
				if !visited[neighborID] {
					queue = append(queue, neighborID)
					visited[neighborID] = true
					parents[neighborID] = currentNode
				}
			}
		}
	}

	if !found {
		fmt.Println("Hedef düğüm bulunamadı!")
		return nil
	}

	// En kısa yolu oluşturmak için hedef düğümden başlayarak ebeveynleri takip et
	shortestPath := make([]int, 0)
	currentNode := targetID
	for currentNode != startID {
		shortestPath = append([]int{currentNode}, shortestPath...)
		currentNode = parents[currentNode]
	}
	shortestPath = append([]int{startID}, shortestPath...)

	return shortestPath
}
