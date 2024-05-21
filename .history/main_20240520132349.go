package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Düğüm yapısını tanımlar
type Dugum struct {
	ID   int
	Isim string
}

// Kenar yapısını tanımlar
type Kenar struct {
	Baslangic int
	Bitis     int
	Kapasite  int
}

// Grafik yapısını tanımlar
type Grafik struct {
	Dugumler         []Dugum
	Kenarlari        []Kenar
	BaslangicDugumID int
	BitisDugumID     int
	KomusulukListesi map[int][]Kenar
}

// Düğümleri yazdırır
func dugumleriYazdir(dugumler []Dugum) {
	fmt.Println("\nodalar:")
	for _, dugum := range dugumler {
		fmt.Printf("%d: %s\n", dugum.ID, dugum.Isim)
	}
}

// Kenarları yazdırır
func kenarlariYazdir(kenarlar []Kenar) {
	fmt.Println("\ntuneller:")
	for _, kenar := range kenarlar {
		fmt.Printf("%d - %d (%d)\n", kenar.Baslangic, kenar.Bitis, kenar.Kapasite)
	}
}

// BFS ile artırılabilir bir yol bulur
func (g *Grafik) bfs(akim [][]int, parent []int) bool {
	visited := make([]bool, len(g.Dugumler))
	queue := []int{g.BaslangicDugumID}
	visited[g.BaslangicDugumID] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for _, kenar := range g.KomusulukListesi[u] {
			if !visited[kenar.Bitis] && akim[u][kenar.Bitis] < kenar.Kapasite {
				queue = append(queue, kenar.Bitis)
				parent[kenar.Bitis] = u
				visited[kenar.Bitis] = true
				if kenar.Bitis == g.BitisDugumID {
					return true
				}
			}
		}
	}

	return false
}

// Edmonds-Karp algoritmasını uygular
func (g *Grafik) edmondsKarp() int {
	n := len(g.Dugumler)
	akim := make([][]int, n)
	for i := range akim {
		akim[i] = make([]int, n)
	}
	parent := make([]int, n)
	maxAkis := 0

	for g.bfs(akim, parent) {
		yolKapasitesi := int(^uint(0) >> 1)
		for v := g.BitisDugumID; v != g.BaslangicDugumID; v = parent[v] {
			u := parent[v]
			for _, kenar := range g.KomusulukListesi[u] {
				if kenar.Bitis == v {
					yolKapasitesi = min(yolKapasitesi, kenar.Kapasite-akim[u][v])
				}
			}
		}

		for v := g.BitisDugumID; v != g.BaslangicDugumID; v = parent[v] {
			u := parent[v]
			akim[u][v] += yolKapasitesi
			akim[v][u] -= yolKapasitesi
		}

		maxAkis += yolKapasitesi
	}

	return maxAkis
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func main() {
	// Dosyayı açar
	dosya, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer dosya.Close()

	grafik := Grafik{
		KomusulukListesi: make(map[int][]Kenar),
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
			baslangicID = len(grafik.Dugumler)
			grafik.BaslangicDugumID = baslangicID
			grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: baslangicID, Isim: isim})
			fmt.Println("baslangic_oda:", grafik.BaslangicDugumID)

		} else if strings.HasPrefix(satir, "##end") {
			tarama.Scan()
			alanlar := strings.Fields(tarama.Text())
			isim := alanlar[0]
			bitisID = len(grafik.Dugumler)
			grafik.BitisDugumID = bitisID
			grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: bitisID, Isim: isim})
			fmt.Println("bitis_oda:", grafik.BitisDugumID)

		} else {
			alanlar := strings.Fields(satir)
			if len(alanlar) == 2 && strings.Contains(alanlar[0], "-") {
				kenarParcalari := strings.Split(alanlar[0], "-")
				baslangicIsmi := kenarParcalari[0]
				bitisIsmi := kenarParcalari[1]
				baslangicID := isimeGoreDugumIDBul(grafik.Dugumler, baslangicIsmi)
				bitisID := isimeGoreDugumIDBul(grafik.Dugumler, bitisIsmi)
				grafik.Kenarlari = append(grafik.Kenarlari, Kenar{Baslangic: baslangicID, Bitis: bitisID, Kapasite: 1})
				grafik.KomusulukListesi[baslangicID] = append(grafik.KomusulukListesi[baslangicID], Kenar{Baslangic: baslangicID, Bitis: bitisID, Kapasite: 1})
				grafik.KomusulukListesi[bitisID] = append(grafik.KomusulukListesi[bitisID], Kenar{Baslangic: bitisID, Bitis: baslangicID, Kapasite: 0})
			} else if len(alanlar) == 3 {
				isim := alanlar[0]
				id := len(grafik.Dugumler)
				grafik.Dugumler = append(grafik.Dugumler, Dugum{ID: id, Isim: isim})
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

	maxAkis := grafik.edmondsKarp()
	fmt.Printf("Maksimum akış: %d\n", maxAkis)

	if maxAkis < karincaSayisi {
		fmt.Printf("Başlangıç ve bitiş arasında %d karıncayı geçirmek için yeterli kapasite yok.\n", karincaSayisi)
		return
	}

	// Edmonds-Karp algoritması ile minimum maksimum akışı bul
akis := grafik.EdmondsKarpMinimumMaksimumAkis()
if akis == 0 {
    fmt.Println("Hata: Başlangıç ve bitiş arasında geçerli bir yol yok.")
    return
}

// Karıncaların bitiş düğümüne olan en kısa yollarını bul
enKisaYollar := grafik.enKisaYollar()
enKisaYollariYazdir(enKisaYollar, grafik.Dugumler)

