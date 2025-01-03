package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type User struct {
	Username string
	Password string
	Status   string
	Friends  []string
}

var users []User
var currentUser *User

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n--- Menu Utama ---")
		fmt.Println("1. Registrasi Pengguna")
		fmt.Println("2. Login")
		fmt.Println("3. Keluar")
		fmt.Print("Pilih menu: ")
		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "1":
			registerUser(scanner)
		case "2":
			login(scanner)
		case "3":
			fmt.Println("Terima kasih!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func registerUser(scanner *bufio.Scanner) {
	fmt.Print("Masukkan username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Print("Masukkan password: ")
	scanner.Scan()
	password := scanner.Text()

	for _, user := range users {
		if user.Username == username {
			fmt.Println("Username sudah terdaftar.")
			return
		}
	}

	users = append(users, User{Username: username, Password: password})
	fmt.Println("Registrasi berhasil.")
}

func login(scanner *bufio.Scanner) {
	fmt.Print("Masukkan username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Print("Masukkan password: ")
	scanner.Scan()
	password := scanner.Text()

	for i, user := range users {
		if user.Username == username && user.Password == password {
			currentUser = &users[i]
			fmt.Println("Login berhasil.")
			userMenu(scanner)
			return
		}
	}
	fmt.Println("Username atau password salah.")
}

func userMenu(scanner *bufio.Scanner) {
	for {
		fmt.Println("\n--- Menu Pengguna ---")
		fmt.Println("1. Lihat Status dan Beri Komentar")
		fmt.Println("2. Tambah/Hapus Teman")
		fmt.Println("3. Edit Profil")
		fmt.Println("4. Lihat Teman")
		fmt.Println("5. Cari Pengguna")
		fmt.Println("6. Logout")
		fmt.Print("Pilih menu: ")
		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "1":
			viewAndCommentStatus(scanner)
		case "2":
			manageFriends(scanner)
		case "3":
			editProfile(scanner)
		case "4":
			viewFriendsSorted()
		case "5":
			searchUser(scanner)
		case "6":
			fmt.Println("Logout berhasil.")
			currentUser = nil
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func viewAndCommentStatus(scanner *bufio.Scanner) {
	for _, user := range users {
		fmt.Printf("\n%s: %s\n", user.Username, user.Status)
	}
	fmt.Print("Ingin beri komentar pada pengguna (masukkan username atau kosongkan untuk kembali)? ")
	scanner.Scan()
	username := scanner.Text()
	if username == "" {
		return
	}
	for i, user := range users {
		if user.Username == username {
			fmt.Print("Masukkan komentar: ")
			scanner.Scan()
			comment := scanner.Text()
			users[i].Status += "\nKomentar: " + comment
			fmt.Println("Komentar berhasil ditambahkan.")
			return
		}
	}
	fmt.Println("Pengguna tidak ditemukan.")
}

func manageFriends(scanner *bufio.Scanner) {
	fmt.Println("1. Tambah Teman")
	fmt.Println("2. Hapus Teman")
	fmt.Print("Pilih aksi: ")
	scanner.Scan()
	choice := scanner.Text()

	fmt.Print("Masukkan username teman: ")
	scanner.Scan()
	friend := scanner.Text()

	switch choice {
	case "1":
		if !isFriend(friend) {
			currentUser.Friends = append(currentUser.Friends, friend)
			fmt.Println("Teman berhasil ditambahkan.")
		} else {
			fmt.Println("Pengguna sudah menjadi teman.")
		}
	case "2":
		for i, f := range currentUser.Friends {
			if f == friend {
				currentUser.Friends = append(currentUser.Friends[:i], currentUser.Friends[i+1:]...)
				fmt.Println("Teman berhasil dihapus.")
				return
			}
		}
		fmt.Println("Pengguna bukan teman Anda.")
	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

func editProfile(scanner *bufio.Scanner) {
	fmt.Print("Masukkan status baru: ")
	scanner.Scan()
	currentUser.Status = scanner.Text()
	fmt.Println("Profil berhasil diperbarui.")
}

func viewFriendsSorted() {
	sort.Strings(currentUser.Friends)
	fmt.Println("Teman Anda (terurut):")
	for _, friend := range currentUser.Friends {
		fmt.Println(friend)
	}
}

func searchUser(scanner *bufio.Scanner) {
	fmt.Print("Masukkan nama pengguna yang dicari: ")
	scanner.Scan()
	keyword := scanner.Text()

	fmt.Println("Hasil pencarian:")
	for _, user := range users {
		if strings.Contains(user.Username, keyword) {
			fmt.Println(user.Username)
		}
	}
}

func isFriend(username string) bool {
	for _, friend := range currentUser.Friends {
		if friend == username {
			return true
		}
	}
	return false
}