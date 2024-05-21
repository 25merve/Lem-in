package main

import (
	"bufio"
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

// BFS algoritması ile en kısa yolu bulur
func (g *Grafik) BFS(baslangicDugumID int, bitisDugumID int) []int {
	kuyruk := [][]int{{baslangicDugumID}}
	ziyaretEdildi := make(map[int]bool)
	ziyaretEdildi[baslangicDugumID] = true

	for len(kuyruk) > 0 {
		yol := kuyruk[0]
		kuyruk = kuyruk[1:]
		dugum := yol[len(yol)-1]

		if dugum == bitisDugumID {
			return yol
		}

		for _, komsu := range g.KomusulukListesi[dugum] {
			if !ziyaretEdildi[komsu] {
				yeniYol := append([]int{}, yol...)
				yeniYol = append(yeniYol, komsu)
				kuyruk = append(kuyruk, yeniYol)
				ziyaretEdildi[komsu] = true
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

	enKisaYol := grafik.BFS(grafik.BaslangicDugumID, grafik.BitisDugumID)
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
}

// Edmonds-Karp algoritması ile maksimum akışı ve en kısa yolu bulur
func (g *Grafik) EdmondsKarp(baslangicDugumID, bitisDugumID int, kapasiteler map[Kenar]int, filtreliDugumler map[int]bool) (int, []int) {
	// Başlangıçta maksimum akış miktarı sıfır olarak ayarlanır
	maxAkis := 0
	// En kısa yol ve artışlar için önceki düğümler
	oncekiDugumler := make(map[int]int)

	// BFS algoritması kullanılarak artan yolu bulunur
	for {
		// BFS kullanarak artan yolu bul
		yol, artis := g.BulArtanYol(baslangicDugumID, bitisDugumID, kapasiteler, filtreliDugumler)
		if yol == nil {
			// Artan yol bulunamadı, döngüden çık
			break
		}

		// Artış miktarını maksimum akışa ekle
		maxAkis += artis

		// Yolu izleyerek kapasite tablosunu güncelle
		g.kapasiteTablosunuGuncelle(yol, artis, kapasiteler)

		// Önceki düğümleri güncelle
		for i := 1; i < len(yol); i++ {
			oncekiDugumler[yol[i]] = yol[i-1]
		}
	}

	// En kısa yol bulunur
	enKisaYol := g.BulEnKisaYol(baslangicDugumID, bitisDugumID, oncekiDugumler)
	if enKisaYol == nil {
		return 0, nil
	}

	return maxAkis, enKisaYol
}

// BulArtanYol, başlangıç ve bitiş düğümleri arasında artan bir yol ve artış miktarı bulur
func (g *Grafik) BulArtanYol(baslangicDugumID, bitisDugumID int, kapasiteler map[Kenar]int, filtreliDugumler map[int]bool) ([]int, int) {
	// Önceki düğümleri saklamak için bir harita oluştur
	oncekiDugumler := make(map[int]int)

	// BFS kuyruğu oluştur
	kuyruk := []int{baslangicDugumID}
	oncekiDugumler[baslangicDugumID] = -1
	artis := math.MaxInt64

	for len(kuyruk) > 0 {
		// Kuyruğun başındaki düğümü al
		dugum := kuyruk[0]
		kuyruk = kuyruk[1:]

		// Bitiş düğümüne ulaşıldıysa, artan yolu oluştur ve artış miktarını döndür
		if dugum == bitisDugumID {
			yol := []int{}
			for dugum != -1 {
				yol = append([]int{dugum}, yol...)
				dugum = oncekiDugumler[dugum]
			}
			return yol, artis
		}

		// Düğümün komşularını kontrol et
		for _, komsu := range g.KomusulukListesi[dugum] {
			// Kapasite kontrolü yap
			kenar := Kenar{Baslangic: dugum, Bitis: komsu}
			geciciArtis := kapasiteler[kenar]
			if geciciArtis > 0 && !filtreliDugumler[komsu] {
				// Artış miktarını güncelle
				if geciciArtis < artis {
					artis = geciciArtis
				}
				// Önceki düğümü kaydet
				oncekiDugumler[komsu] = dugum
				// Kuyruğa ekle
				kuyruk = append(kuyruk, komsu)
			}
		}
	}

	// Bitiş düğümüne ulaşılamadı, artan yol yok
	return nil, 0
}

// BulEnKisaYol, başlangıç ve bitiş düğümleri arasında en kısa yolu bulur
func (g *Grafik) BulEnKisaYol(baslangicDugumID, bitisDugumID int, oncekiDugumler map[int]int) []int {
	// En kısa yolun uzunluğunu bul
	uzunluk := 0
	dugum := bitisDugumID
	for dugum != -1 {
		uzunluk++
		dugum = oncekiDugumler[dugum]
	}

	// En kısa yol dizisini oluştur
	yol := make([]int, uzunluk)
	dugum = bitisDugumID
	for i := uzunluk - 1; i >= 0; i-- {
		yol[i] = dugum
		dugum = oncekiDugumler[dugum]
	}

	return yol
}

// Kapasite tablosunu günceller
func (g *Grafik) kapasiteTablosunuGuncelle(yol []int, akis int, kapasiteler map[Kenar]int) {
	// Yolu izleyerek kapasite tablosunu güncelle
	for i := 0; i < len(yol)-1; i++ {
		kenar := Kenar{Baslangic: yol[i], Bitis: yol[i+1]}
		// İlgili kenarın kapasitesinden akış miktarı çıkarılır
		kapasiteler[kenar] -= akis
		// Ters yöndeki kenarın kapasitesine akış miktarı eklenir (geri dönüş)
		tersKenar := Kenar{Baslangic: yol[i+1], Bitis: yol[i]}
		kapasiteler[tersKenar] += akis
	}
}
