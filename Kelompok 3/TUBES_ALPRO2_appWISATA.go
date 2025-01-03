package main

import (
	"fmt"
	"sort"
	"strings"
)

type TempatWisata struct {
	ID        int
	Nama      string
	Kategori  string
	Fasilitas []string
	Wahana    []string
	Biaya     float64
	Jarak     float64
}

var tempatWisata []TempatWisata
var idCounter int

func main() {
	for {
		fmt.Println("\nMenu Utama:")
		fmt.Println("1. Admin")
		fmt.Println("2. Pengguna")
		fmt.Println("3. Keluar")
		fmt.Print("Pilih mode: ")

		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			menuAdmin()
		case 2:
			menuPengguna()
		case 3:
			fmt.Println("Keluar dari aplikasi.")
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func menuAdmin() {
	for {
		fmt.Println("\nMenu Admin:")
		fmt.Println("1. Tambah Tempat Wisata")
		fmt.Println("2. Ubah Tempat Wisata")
		fmt.Println("3. Hapus Tempat Wisata")
		fmt.Println("4. Kembali ke Menu Utama")
		fmt.Print("Pilih menu: ")

		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			tambahTempatWisata()
		case 2:
			ubahTempatWisata()
		case 3:
			hapusTempatWisata()
		case 4:
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func menuPengguna() {
	for {
		fmt.Println("\nMenu Pengguna:")
		fmt.Println("1. Lihat Daftar Tempat Wisata")
		fmt.Println("2. Cari Tempat Wisata")
		fmt.Println("3. Kembali ke Menu Utama")
		fmt.Print("Pilih menu: ")

		var pilihan int
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			lihatDaftarTerurut()
		case 2:
			cariTempatWisataKategori()
		case 3:
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

// --- Fungsi Admin ---
func tambahTempatWisata() {
	idCounter++
	var nama, kategori, fasilitasInput, wahanaInput string
	var biaya, jarak float64

	fmt.Print("Nama Tempat Wisata: ")
	fmt.Scan(&nama)
	fmt.Print("Kategori: ")
	fmt.Scan(&kategori)
	fmt.Print("Biaya: ")
	fmt.Scan(&biaya)
	fmt.Print("Jarak(km): ")
	fmt.Scan(&jarak)
	fmt.Print("Fasilitas (pisahkan dengan koma): ")
	fmt.Scan(&fasilitasInput)
	fmt.Print("Wahana (pisahkan dengan koma): ")
	fmt.Scan(&wahanaInput)

	fasilitas := strings.Split(fasilitasInput, ",")
	wahana := strings.Split(wahanaInput, ",")

	tempatWisata = append(tempatWisata, TempatWisata{
		ID:        idCounter,
		Nama:      nama,
		Kategori:  kategori,
		Fasilitas: fasilitas,
		Wahana:    wahana,
		Biaya:     biaya,
		Jarak:     jarak,
	})

	fmt.Println("Tempat wisata berhasil ditambahkan!")
}

func ubahTempatWisata() {
	var id int
	fmt.Print("Masukkan ID tempat wisata yang ingin diubah: ")
	fmt.Scan(&id)

	for i, tempat := range tempatWisata {
		if tempat.ID == id {
			var nama, kategori, fasilitasInput, wahanaInput string
			var biaya, jarak float64

			fmt.Print("Nama Tempat Wisata baru: ")
			fmt.Scan(&nama)
			fmt.Print("Kategori baru: ")
			fmt.Scan(&kategori)
			fmt.Print("Biaya baru: ")
			fmt.Scan(&biaya)
			fmt.Print("Jarak(km) baru: ")
			fmt.Scan(&jarak)
			fmt.Print("Fasilitas baru (pisahkan dengan koma): ")
			fmt.Scan(&fasilitasInput)
			fmt.Print("Wahana baru (pisahkan dengan koma): ")
			fmt.Scan(&wahanaInput)

			fasilitas := strings.Split(fasilitasInput, ",")
			wahana := strings.Split(wahanaInput, ",")

			tempatWisata[i] = TempatWisata{
				ID:        id,
				Nama:      nama,
				Kategori:  kategori,
				Fasilitas: fasilitas,
				Wahana:    wahana,
				Biaya:     biaya,
				Jarak:     jarak,
			}

			fmt.Println("Tempat wisata berhasil diubah!")
			return
		}
	}
	fmt.Println("Tempat wisata tidak ditemukan!")
}

func hapusTempatWisata() {
	var id int
	fmt.Print("Masukkan ID tempat wisata yang ingin dihapus: ")
	fmt.Scan(&id)

	for i, tempat := range tempatWisata {
		if tempat.ID == id {
			tempatWisata = append(tempatWisata[:i], tempatWisata[i+1:]...)
			fmt.Println("Tempat wisata berhasil dihapus!")
			return
		}
	}
	fmt.Println("Tempat wisata tidak ditemukan!")
}

// --- Fungsi Pengguna ---
func lihatDaftarTerurut() {
	fmt.Println("\nPilih pengurutan:")
	fmt.Println("1. Berdasarkan Jarak")
	fmt.Println("2. Berdasarkan Biaya")
	fmt.Print("Pilihan: ")

	var pilihan int
	fmt.Scan(&pilihan)

	switch pilihan {
	case 1:
		sort.Slice(tempatWisata, func(i, j int) bool {
			return tempatWisata[i].Jarak < tempatWisata[j].Jarak
		})
	case 2:
		sort.Slice(tempatWisata, func(i, j int) bool {
			return tempatWisata[i].Biaya < tempatWisata[j].Biaya
		})
	default:
		fmt.Println("Pilihan tidak valid!")
		return
	}

	fmt.Println("Daftar Tempat Wisata:")
	for _, tempat := range tempatWisata {
		fmt.Printf("ID: %d, Nama: %s, Kategori: %s, Biaya: %.2f, Jarak: %.2f, Fasilitas: %v, Wahana: %v\n",
			tempat.ID, tempat.Nama, tempat.Kategori, tempat.Biaya, tempat.Jarak, tempat.Fasilitas, tempat.Wahana)
	}
}

func cariTempatWisataKategori() {
	var kataKunci, kategori string
	fmt.Print("Masukkan kata kunci: ")
	fmt.Scan(&kataKunci)
	fmt.Print("Masukkan kategori: ")
	fmt.Scan(&kategori)

	fmt.Println("Hasil Pencarian:")
	for _, tempat := range tempatWisata {
		if strings.Contains(strings.ToLower(tempat.Nama), strings.ToLower(kataKunci)) &&
			strings.EqualFold(tempat.Kategori, kategori) {
			fmt.Printf("ID: %d, Nama: %s, Kategori: %s, Biaya: %.2f, Jarak: %.2f, Fasilitas: %v, Wahana: %v\n",
				tempat.ID, tempat.Nama, tempat.Kategori, tempat.Biaya, tempat.Jarak, tempat.Fasilitas, tempat.Wahana)
		}
	}
}
