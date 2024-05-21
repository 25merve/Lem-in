package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Düğüm yapısını tanımlar
type Dugum struct {
	ID   int
	Isim string
	X    int
	Y    int
}

// Kenar yapısını tanımlar
type Kenar struct {
	Baslangic int
	Bitis     int
}

// Grafik yapısını tanımlar
type Grafik struct {
	Dugumler         []Dugum
	Kenarlari        []Kenar
	BaslangicDugumID int
	BitisDugumID     int
	KomusulukListesi map[int][]int
}

// Düğümleri yazdırır
func dugumleriYazdir(dugumler []Dugum) {
	fmt.Println("\nodalar:")
	for _, dugum := range dugumler {
		fmt.Printf("%d: %s (%d, %d)\n", dugum.ID, dugum.Isim, dugum.X, dugum.Y)
	}
}

// Kenarları yazdırır
func kenarlariYazdir(kenarlar []Kenar) {
	fmt.Println("\ntuneller:")
	for _, kenar := range kenarlar {
		fmt.Printf("%d - %d\n", kenar.Baslangic, kenar.Bitis)
	}
}

// A* algoritması ile en kısa yolu bulur
func (g *Grafik) AStar(baslangicDugumID int, bitisDugumID int) []int {
	type Node struct {
		id       int
		priority float64
		index    int
		path     []int
	}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	startNode := &Node{
		id:       baslangicDugumID,
		priority: 0,
		path:     []int{baslangicDugumID},
	}
	heap.Push(&pq, startNode)

	ziyaretEdildi := make(map[int]bool)
	gScore := make(map[int]float64)
	gScore[baslangicDugumID] = 0

	heuristic := func(a, b Dugum) float64 {
		return math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y)))
	}

	for pq.Len() > 0 {
		currentNode := heap.Pop(&pq).(*Node)
		current := currentNode.id

		if current == bitisDugumID {
			return currentNode.path
		}

		if ziyaretEdildi[current] {
			continue
		}
		ziyaretEdildi[current] = true

		for _, komsu := range g.KomusulukListesi[current] {
			if ziyaretEdildi[komsu] {
				continue
			}

			tentativeGScore := gScore[current] + heuristic(g.Dugumler[current], g.Dugumler[komsu])

			if score, exists := gScore[komsu]; !exists || tentativeGScore < score {
				gScore[komsu] = tentativeGScore
				fScore := tentativeGScore + heuristic(g.Dugumler[komsu], g.Dugumler[bitisDugumID])
				newPath := append([]int{}, currentNode.path...)
				newPath = append(newPath, komsu)
				heap.Push(&pq, &Node{
					id:       komsu,
					priority: fScore,
					path:     newPath,
				})
			}
		}
	}

	return nil
}

// İsme göre düğüm ID'si bulur
func isimeGoreDugumIDBul(dugumler []Dugum, isim string) int {
	for _, dugum := range dugumler {
		if dugum.Isim == isim {
			return dugum.ID
		}
	}
	return -1
}

// PriorityQueue implementasyonu
type PriorityQueue []*Node

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

	item := x.(*Node)
	item.index = len(*pq)
	*pq = append(*pq, item)
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

func main() {
	// Dosyayı açar
	dosya, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer dosya.Close()

	grafik := Grafik{
		KomusulukListesi: make(map[int][]int),
	}

	tarama := bufio.NewScanner(dosya)

	// İlk satırdaki karınca sayısını okur
	tarama.Scan()
	karincaSayisiSatiri := tarama.Text()
	karincaSayisi, err := strconv.Atoi(karincaSayisiSatiri)
	if err != nil {
		fmt.Println("Karınca sayısı okunamadı:", err)
		return
	}

	baslangicID, bitisID := -1, -1

	// Dosya satırlarını okur
	for tarama.Scan() {
		satir := tarama.Text()
		if strings.HasPrefix(satir, "##start") {
			tarama.Scan()
			alanlar := strings.Fields(tarama.Text())
			isim := alanlar[0]
			x, _ := strconv.Atoi(alanlar[1])
			y, _ := strconv.Atoi(alanlar[2])
			baslangicID = len(grafik.Dugumler)
			grafik.BaslangicDugumID = baslangicID
			grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: baslangicID, Isim: isim, X: x, Y: y})
			fmt.Println("baslangic_oda:", grafik.BaslangicDugumID)

		} else if strings.HasPrefix(satir, "##end") {
			tarama.Scan()
			alanlar := strings.Fields(tarama.Text())
			isim := alanlar[0]
			x, _ := strconv.Atoi(alanlar[1])
			y, _ := strconv.Atoi(alanlar[2])
			bitisID = len(grafik.Dugumler)
			grafik.BitisDugumID = bitisID
			grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: bitisID, Isim: isim, X: x, Y: y})
			fmt.Println("bitis_oda:", grafik.BitisDugumID)

		} else {
			alanlar := strings.Fields(satir)
			if len(alanlar) == 3 {
				isim := alanlar[0]
				x, _ := strconv.Atoi(alanlar[1])
				y, _ := strconv.Atoi(alanlar[2])
				id := len(grafik.Dugumler)
				grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: id, Isim: isim, X: x, Y: y})
			} else if len(alanlar) == 1 && strings.Contains(satir, "-") {
				kenarParcalari := strings.Split(alanlar[0], "-")
				baslangicIsmi := kenarParcalari[0]
				bitisIsmi := kenarParcalari[1]
				baslangicID := isimeGoreDugumIDBul(grafik.Dugumler, baslangicIsmi)
				bitisID := isimeGoreDugumIDBul(grafik.Dugumler, bitisIsmi)
				grafik.Kenarlari = append(grafik.Kenarlari, Kenar{Baslangic: baslangicID, Bitis: bitisID})
				grafik.KomusulukListesi[baslangicID] = append(grafik.KomusulukListesi[baslangicID], bitisID)
				grafik.KomusulukListesi[bitisID] = append(grafik.KomusulukListesi[bitisID], baslangicID)
			}
		}
	}

	if err := tarama.Err(); err != nil {
		fmt.Println("Okuma hatası:", err)
		return
	}

	if baslangicID == -1 || bitisID == -1 {
		fmt.Println("Hata: Başlangıç veya bitiş düğümü belirtilmemiş.")
		return
	}

	dugumleriYazdir(grafik.Dugumler)
	kenarlariYazdir(grafik.Kenarlari)

	enKisaYol := grafik.AStar(grafik.BaslangicDugumID, grafik.BitisDugumID)
	if enKisaYol == nil {
		fmt.Println("Hata: Başlangıç ve bitiş arasında geçerli bir yol yok.")
		return
	}
	fmt.Println("En kısa yol:", enKisaYol)

	karincaYollari := make([][]int, karincaSayisi)
	for i := 0; i < karincaSayisi; i++ {
		karincaYollari[i] = enKisaYol
	}

	karincaKonumlari := make([]int, karincaSayisi)
	for i := 0; i < karincaSayisi; i++ {
		karincaKonumlari[i] = grafik.BaslangicDugumID
	}

	adim := 1
	for {
		hareketler := []string{}
		hepsiBitisteMi := true
		yerlesik := make(map[int]bool)
		for i := 0; i < karincaSayisi; i++ {
			if karincaKonumlari[i] != grafik.BitisDugumID {
				hepsiBitisteMi = false
				yol := karincaYollari[i]
				for j := 0; j < len(yol)-1; j++ {
					if yol[j] == karincaKonumlari[i] && !yerlesik[yol[j+1]] {
						yerlesik[yol[j+1]] = true
						karincaKonumlari[i] = yol[j+1]
						hareketler = append(hareketler, fmt.Sprintf("L%d-%s", i+1, grafik.Dugumler[karincaKonumlari[i]].Isim))
						break
					}
				}
			}
		}
		if hepsiBitisteMi {
			break
		}
		fmt.Printf("Adım %d: %s\n", adim, strings.Join(hareketler, " "))
		adim++
	}
	Node:=
}