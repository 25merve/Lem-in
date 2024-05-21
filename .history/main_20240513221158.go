package main 

import
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
        if len(fields) == 1 {
				// Eğer satır sadece bir sayıdan oluşuyorsa, bu karınca sayısıdır.
				antCount, _ = strconv.Atoi(fields[0])
        }if err := scanner.Err(); err != nil {
		fmt.Println("Okuma hatasi:", err)
		return
	}