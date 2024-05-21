package main

import (
	"bufio"
	"fmt"
	"math"
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

	var antCount int         // Karınca sayısını tutacak değişken
	startID, endID := -1, -1 //ilk önce baslangic ve bitis satırlarına bi id atıyoruz
	graph := Graph{}         //struct yapısını değişkene atıyoruz

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if antCount == 0 { // Karınca sayısını oku
			antCount, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("Karınca sayısı okunurken hata oluştu:", err)
				return
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

	// İlişkili düğümleri kullanarak her bir kenarın uzunluğunu hesaplayın
	fmt.Println("\nKenar Uzunluklari:")
	for _, edge := range graph.Edges {
		startNode := findNodeByID(graph.Nodes, edge.Start)         //dügümler kümesi içinden başlangıç olanı alıyor
		endNode := findNodeByID(graph.Nodes, edge.End)             //düğümler kümesi içinden bitiş olanı alıyor
		uzunluk := hesapla(startNode, endNode)                     //uzunluğunu bulmak için bu iki bilgiyi hesapla fonksiyonuna yolluyor
		fmt.Printf("%d - %d: %f\n", edge.Start, edge.End, uzunluk) //ve bunu bastırma işlemi yapıyor ama hangi kenarlar arasındaki olduğunu da yazıyor
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

// İki düğüm arasındaki mesafeyi hesapla
func hesapla(node1, node2 Node) float64 {
	dx := float64(node2.X - node1.X)
	dy := float64(node2.Y - node1.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
