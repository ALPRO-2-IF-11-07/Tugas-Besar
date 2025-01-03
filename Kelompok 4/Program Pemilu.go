package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type Calon struct {
	ID     int
	Nama   string
	Partai string
	Suara  int
}

type Pemilih struct {
	ID   int
	Nama string
}

var (
	calons  []Calon
	pemilih []Pemilih
)

func saveResultsToFile() {
	file, err := os.Create("hasil_pemilihan.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	file.WriteString("HASIL PEMILIHAN UMUM\n")
	file.WriteString("===================\n\n")

	sort.Slice(calons, func(i, j int) bool {
		return calons[i].Suara > calons[j].Suara
	})

	for _, calon := range calons {
		result := fmt.Sprintf("ID: %d\nNama: %s\nPartai: %s\nJumlah Suara: %d\n\n",
			calon.ID, calon.Nama, calon.Partai, calon.Suara)
		file.WriteString(result)
	}

	fmt.Println("Hasil pemilihan telah disimpan dalam file 'hasil_pemilihan.txt'")
}

func startVoting(duration time.Duration, done chan bool) {
	fmt.Println("Pemilihan dimulai! Anda memiliki waktu:", duration.Minutes(), "menit.")
	time.Sleep(duration)
	fmt.Println("Waktu pemilihan telah habis.")
	done <- true
}

func displayCalons() {
	fmt.Println("Daftar Calon:")
	for _, calon := range calons {
		fmt.Printf("ID: %d, Nama: %s, Partai: %s, Suara: %d\n", calon.ID, calon.Nama, calon.Partai, calon.Suara)
	}
}

func addCalon(nama, partai string) {
	calon := Calon{ID: len(calons) + 1, Nama: nama, Partai: partai, Suara: 0}
	calons = append(calons, calon)
	fmt.Println("Calon berhasil ditambahkan:", calon)
}

func editCalon(id int, nama, partai string) {
	for i, calon := range calons {
		if calon.ID == id {
			calons[i].Nama = nama
			calons[i].Partai = partai
			fmt.Println("Calon berhasil diubah:", calons[i])
			return
		}
	}
	fmt.Println("Calon tidak ditemukan.")
}

func deleteCalon(id int) {
	for i, calon := range calons {
		if calon.ID == id {
			calons = append(calons[:i], calons[i+1:]...)
			fmt.Println("Calon berhasil dihapus.")
			return
		}
	}
	fmt.Println("Calon tidak ditemukan.")
}

func vote(pemilihID, calonID int) {
	for i, calon := range calons {
		if calon.ID == calonID {
			calons[i].Suara++
			fmt.Printf("Pemilih %d telah memberikan suara untuk calon %s\n", pemilihID, calon.Nama)
			return
		}
	}
	fmt.Println("Calon tidak ditemukan.")
}

func displaySortedCalonsByVotes() {
	sort.Slice(calons, func(i, j int) bool {
		return calons[i].Suara > calons[j].Suara
	})
	displayCalons()
}

func displaySortedCalonsByName() {
	sort.Slice(calons, func(i, j int) bool {
		return calons[i].Nama < calons[j].Nama
	})
	displayCalons()
}

func searchCalonByPartai(partai string) {
	fmt.Println("Calon dari Partai:", partai)
	for _, calon := range calons {
		if calon.Partai == partai {
			fmt.Printf("ID: %d, Nama: %s, Suara: %d\n", calon.ID, calon.Nama, calon.Suara)
		}
	}
}

func searchCalonByName(nama string) {
	fmt.Println("Calon dengan Nama:", nama)
	for _, calon := range calons {
		if calon.Nama == nama {
			fmt.Printf("ID: %d, Partai: %s, Suara: %d\n", calon.ID, calon.Partai, calon.Suara)
		}
	}
}

func checkThreshold(threshold int) {
	fmt.Println("Calon yang memenuhi threshold:", threshold)
	for _, calon := range calons {
		if calon.Suara >= threshold {
			fmt.Printf("ID: %d, Nama: %s, Suara: %d\n", calon.ID, calon.Nama, calon.Suara)
		}
	}
}

func main() {
	for {
		var choice int
		fmt.Println("\nSelamat datang di Aplikasi Pemilihan Umum")
		fmt.Println("1. Login sebagai Pemilih")
		fmt.Println("2. Login sebagai Petugas KPU")
		fmt.Println("3. Simpan Hasil dan Keluar")
		fmt.Print("Pilih opsi: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var pemilihID int
			fmt.Print("Masukkan ID Pemilih: ")
			fmt.Scan(&pemilihID)

			done := make(chan bool)
			go startVoting(30*time.Second, done)

			for {
				select {
				case <-done:
					fmt.Println("Waktu pemilihan telah habis. Kembali ke menu utama.")
					goto MainMenu
				default:
					displayCalons()
					var calonID int
					fmt.Print("Masukkan ID Calon yang ingin dipilih (atau 0 untuk kembali): ")
					fmt.Scan(&calonID)

					if calonID == 0 {
						fmt.Println("Kembali ke menu utama.")
						goto MainMenu
					}

					vote(pemilihID, calonID)
				}
			}

		case 2:
			for {
				var adminChoice int
				fmt.Println("\nMenu Petugas KPU:")
				fmt.Println("1. Tambah Calon")
				fmt.Println("2. Edit Calon")
				fmt.Println("3. Hapus Calon")
				fmt.Println("4. Tampilkan Calon")
				fmt.Println("5. Tampilkan Calon Terurut Berdasarkan Suara")
				fmt.Println("6. Tampilkan Calon Terurut Berdasarkan Nama")
				fmt.Println("7. Pencarian Calon Berdasarkan Partai")
				fmt.Println("8. Pencarian Calon Berdasarkan Nama")
				fmt.Println("9. Cek Threshold Calon")
				fmt.Println("10. Kembali ke Menu Utama")
				fmt.Print("Pilih opsi: ")
				fmt.Scan(&adminChoice)

				switch adminChoice {
				case 1:
					var nama, partai string
					fmt.Print("Masukkan Nama Calon: ")
					fmt.Scan(&nama)
					fmt.Print("Masukkan Partai: ")
					fmt.Scan(&partai)
					addCalon(nama, partai)

				case 2:
					var id int
					var nama, partai string
					fmt.Print("Masukkan ID Calon yang ingin diubah: ")
					fmt.Scan(&id)
					fmt.Print("Masukkan Nama Baru: ")
					fmt.Scan(&nama)
					fmt.Print("Masukkan Partai Baru: ")
					fmt.Scan(&partai)
					editCalon(id, nama, partai)

				case 3:
					var id int
					fmt.Print("Masukkan ID Calon yang ingin dihapus: ")
					fmt.Scan(&id)
					deleteCalon(id)

				case 4:
					displayCalons()

				case 5:
					displaySortedCalonsByVotes()

				case 6:
					displaySortedCalonsByName()

				case 7:
					var partai string
					fmt.Print("Masukkan Nama Partai untuk pencarian: ")
					fmt.Scan(&partai)
					searchCalonByPartai(partai)

				case 8:
					var nama string
					fmt.Print("Masukkan Nama Calon untuk pencarian: ")
					fmt.Scan(&nama)
					searchCalonByName(nama)

				case 9:
					var threshold int
					fmt.Print("Masukkan nilai threshold: ")
					fmt.Scan(&threshold)
					checkThreshold(threshold)

				case 10:
					goto MainMenu

				default:
					fmt.Println("Opsi tidak valid.")
				}
			}

		case 3:
			saveResultsToFile()
			fmt.Println("Terima kasih telah menggunakan Aplikasi Pemilihan Umum. Sampai jumpa!")
			return

		default:
			fmt.Println("Opsi tidak valid.")
		}
	MainMenu:
	}
}
