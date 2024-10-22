package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*

1. Buatkan sebuah func untuk countdown Func ini di jalankan di go countdown hasilnya di tampilkan di main func

2. buatkan 3 func ada parameter ada channel dan message nanti dijalankan di goroutine yang berbeda buatkan juga masing masing punya channelnya sendiri-sendiri tampilakan pesan dari 3 func ini yang paling cepat responnya

3. Buatkan 1 func, func itu memfilter nilai yang ganjil untuk di tampilkan kemudian jalankan func sebanyak 10x di masing2 gorutine memiliki parameter *sync.WaitGroup dan range angka yang jalankan

4. Buatkan 1 object struct di dalamnya ada properti Qty name, buatkan sebuah func yang fungsinya untuk mengurangi qty dari struct tersebut func ini memiliki param jlmQty yang akan di kurangi selanjutnya jalankan func itu dalam goroutine sebanyak 10x lalu terakhir cetak qty dari struct tersebut

5. Buatkan pengecekan suhu

buatkan aplikasi yang memantau 3 sensor (Misal: Suhu, Kelembapan dan tekanan) menggunakan goroutine. setiap sensor akan mengirim data secara berkala melalui buffered channel dan program harus mengambil data tersebut secara interval tertentu menggunakan time.Ticker

Persyaratan

- implementasi 3 Goroutine yang mensimulasikan pengiriman data sensor dengan buffered channel
- gunakan time.ticker untuk mengambil data dari setiap sensor setiap 2 detik
- terapkan timeout(5 detik) menggunakan time.after dan jika sensor tidak merespon dalam waktu tersebut, cetak pesan "Sensor Timeout"
- Pastikan semua sensor tertutup dengan benar dan goroutine selesai

*/

//  Jawaban No 2

func first( ch chan string)  {
	time.Sleep(3 * time.Second)
	ch <- "data terlambat"
}

func second( ch chan string)  {
	time.Sleep(2 * time.Second)
	ch <- "data sedang"

}

func third( ch chan string)  {
	time.Sleep(1 * time.Second)
	ch <- "Data tercepat"

}

// Jawaban no 3

func filterDataGanjil(num int, wg *sync.WaitGroup)  {
	defer wg.Done()
	for i := 0; i < num; i++ {
		if i%2 != 0 {
			fmt.Println("ini adalah nilai Ganjil", i)
		}
	}
}

// Jawaban No 4

type Product struct {
	Name string
	Qty int
	mu sync.Mutex
}

func (p *Product) dec(jmlQty int) {

	p.mu.Lock()
	p.Qty -= jmlQty
	p.mu.Unlock()
}



func main() {

	fmt. Println("--- Jawaban No 1 ---")

	ch := make(chan int, 3)

	go func() {
		
		for i := 3; i > 0; i-- {
			ch <- i
			fmt.Println("Hitungan Mundur Ke", i)
		}
		
		close(ch)
	}()
	time.Sleep(3 * time.Second)

	for value := range ch {
		fmt.Println("menerima", value)
	}

	fmt.Println("Semua data telah diterima.")

	fmt.Println("---Jawaban No 2---")

	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	go first(ch1)
	go second(ch2)
	go third(ch3)

	select {
	case data := <-ch1:
		fmt.Println("Menerima:", data)
	case data := <-ch2:
		fmt.Println("Menerima:", data)
	case data := <-ch3:
		fmt.Println("Menerima:", data)
	}

	ch4 := make(chan string)
	ch5 := make(chan string)
	ch6 := make(chan string)

	// Goroutine untuk mengirim data ke channel 1 setelah 2 detik
	go func() {
		time.Sleep(3 * time.Second)
		ch4 <- "Data dari channel 1"
	}()
	// Goroutine untuk mengirim data ke channel 2 setelah 1 detik
	go func() {
		time.Sleep(2 * time.Second)
		ch5 <- "Data dari channel 2"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch6 <- "Data dari channel 3"
	}()

	// Menggunakan select untuk menangani kedua channel
	// default statement
	for i := 0; i < 1; i++ {
		select {
		case msg1 := <-ch4:
			fmt.Println("Menerima dari channel 1:", msg1)
		case msg2 := <-ch5:
			fmt.Println("Menerima dari channel 2:", msg2)
		case msg3 := <-ch6:
			fmt.Println("Menerima dari channel 3:", msg3)
		default:
			fmt.Println("Tidak ada data yang siap untuk diterima.")
			time.Sleep(500 * time.Millisecond) // Simulasi operasi lain
		}
	}

	// Jawaban no 3

	fmt. Println("--- Jawaban No 3 ---")

	numCPU := runtime.NumCPU()

	fmt.Printf("Jumlah CPU %d\n", numCPU)

	runtime.GOMAXPROCS(2)

	for i := 0; i < 5; i++ {
		go func ()  {
			time.Sleep(1 * time.Second)
			fmt.Println("Selesai menjalankan Goroutine")
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Selesai Menjalankan Program")

	var wg sync.WaitGroup
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go filterDataGanjil(10, &wg)	
	}

	wg.Wait()
	fmt.Println("Semua Pekerjaan selesai")

	fmt. Println("--- Jawaban No 4 ---")

	var wa sync.WaitGroup

	product := Product{
		Qty: 100,
		Name: "Sambalado",
	}

	jmlQty := 5

	for i := 0; i < 10; i++ {
		wa.Add(1)
		
		go func ()  {
			defer wa.Done()
			go product.dec(jmlQty)
			fmt.Printf("Nama Produk %s sisa %d\n", product.Name, product.Qty)
		}()
	}

	wa.Wait()

	fmt.Println("Qty Akhir", product.Qty)

}