package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type User struct {
	Username string
	Password string
	Role     string
}

type Course struct {
	ID          int
	Title       string
	Description string
}

type Task struct {
	ID          int
	CourseID    int
	Title       string
	Description string
}

type Quiz struct {
	ID        int
	CourseID  int
	Title     string
	Questions []string
}

type Forum struct {
	ID          int
	CourseID    int
	Topic       string
	Description string
	Posts       []string
}

type StudentAnswer struct {
	ID          int    // ID unik jawaban
	TaskID      int    // ID tugas terkait
	StudentName string // Nama mahasiswa
	Answer      string // Jawaban mahasiswa
	Graded      bool   // Status apakah sudah dinilai
	Grade       string // Nilai yang diberikan dosen
}

var studentAnswers []StudentAnswer
var nextAnswerID = 1

type LMSData struct {
	Courses []Course `json:"courses"`
	Tasks   []Task   `json:"tasks"`
	Quizzes []Quiz   `json:"quizzes"`
	Forums  []Forum  `json:"forums"`
}

var users = []User{
	{"Angel", "An123", "dosen"},
	{"Udin", "Ud123", "dosen"},
	{"Firtri", "Fi123", "mahasiswa"},
	{"Agus", "Ag123", "mahasiswa"},
}

var lmsData LMSData
var nextCourseID = 1
var nextTaskID = 1
var nextQuizID = 1
var nextForumID = 1

const dataFile = "lms_data.json"

func main() {
	loadData()
	defer saveData()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("====================================")
		fmt.Println("     Aplikasi Moodle Sederhana!")
		fmt.Println("====================================")

		fmt.Println("1. Login sebagai Dosen")
		fmt.Println("2. Login sebagai Mahamahasiswa")
		fmt.Println("3. Keluar")
		fmt.Println("====================================")
		fmt.Print("Pilih opsi login: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		var role string
		if choice == "1" {
			role = "dosen"
		} else if choice == "2" {
			role = "mahasiswa"
		} else if choice == "3" {
			fmt.Println("Keluar dari program.")
			return
		} else {
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
			continue
		}

		// Login
		fmt.Print("Masukkan username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Masukkan password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		// Autentikasi
		user := authenticate(username, password, role)
		if user == nil {
			fmt.Println("\n====================================")
			fmt.Println("Login gagal. Username atau password salah.")
			fmt.Println("====================================")
			continue
		}

		fmt.Printf("\n====================================\n")
		fmt.Printf("Login berhasil sebagai %s (%s)\n", user.Username, user.Role)
		fmt.Println("====================================")

		// Menu berdasarkan role
		if user.Role == "dosen" {
			dosenMenu(reader)
		} else if user.Role == "mahasiswa" {
			mahasiswaMenu(reader, user)
		}
	}
}

func authenticate(username, password, role string) *User {
	for _, user := range users {
		if user.Username == username && user.Password == password && user.Role == role {
			return &user
		}
	}
	return nil
}

func dosenMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n====================================")
		fmt.Println("              Menu dosen")
		fmt.Println("====================================")
		fmt.Println("1. Tambah Course")
		fmt.Println("2. Lihat Data")
		fmt.Println("3. Tambah Tugas")
		fmt.Println("4. Tambah Quiz")
		fmt.Println("5. Tambah Forum Diskusi")
		fmt.Println("6. Hapus Data")
		fmt.Println("7. Lihat Jawaban Mahasiswa")
		fmt.Println("8. Beri Nilai")
		fmt.Println("9. Logout")
		fmt.Println("====================================")
		fmt.Print("Pilih menu: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			tambahCourse(reader)
		case "2":
			lihatSemuaData()
		case "3":
			tambahTugas(reader)
		case "4":
			tambahQuiz(reader)
		case "5":
			tambahForum(reader)
		case "6":
			hapusData(reader)
		case "7":
			lihatJawabanMahasiswa()
		case "8":
			beriNilai(reader)
		case "9":
			fmt.Println("Logout dari menu dosen.")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func searchMenu(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Menu Pencarian")
	fmt.Println("====================================")
	fmt.Println("1. Cari Course")
	fmt.Println("2. Cari Tugas")
	fmt.Println("3. Cari Quiz")
	fmt.Println("4. Cari Forum")
	fmt.Println("5. Kembali")
	fmt.Println("====================================")
	fmt.Print("Pilih menu: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		searchCourseByTitle(reader)
	case "2":
		searchTaskByTitle(reader)
	case "3":
		searchQuizByTitle(reader)
	case "4":
		searchForumByTitle(reader)
	case "5":
		return
	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

func searchCourseByTitle(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Cari Course")
	fmt.Println("====================================")
	fmt.Print("Masukkan kata kunci judul course: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	found := false
	for _, course := range lmsData.Courses {
		if strings.Contains(strings.ToLower(course.Title), strings.ToLower(keyword)) {
			fmt.Printf("ID: %d, Title: %s, Description: %s\n", course.ID, course.Title, course.Description)
			found = true
		}
	}

	if !found {
		fmt.Println("Course tidak ditemukan.")
	}
}

func searchQuizByTitle(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Cari Quiz")
	fmt.Println("====================================")
	fmt.Print("Masukkan kata kunci judul quiz: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	found := false
	for _, quiz := range lmsData.Quizzes {
		if strings.Contains(strings.ToLower(quiz.Title), strings.ToLower(keyword)) {
			fmt.Printf("ID: %d, Course ID: %d, Title: %s, Questions: %v\n", quiz.ID, quiz.CourseID, quiz.Title, quiz.Questions)
			found = true
		}
	}

	if !found {
		fmt.Println("Quiz tidak ditemukan.")
	}
}
func searchTaskByTitle(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Cari Tugas")
	fmt.Println("====================================")
	fmt.Print("Masukkan kata kunci judul tugas: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	found := false
	for _, task := range lmsData.Tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(keyword)) {
			fmt.Printf("ID: %d, Course ID: %d, Title: %s, Description: %s\n", task.ID, task.CourseID, task.Title, task.Description)
			found = true
		}
	}

	if !found {
		fmt.Println("Tugas tidak ditemukan.")
	}
}

func searchForumByTitle(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Cari Forum")
	fmt.Println("====================================")
	fmt.Print("Masukkan kata kunci topik forum: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	found := false
	for _, forum := range lmsData.Forums {
		if strings.Contains(strings.ToLower(forum.Topic), strings.ToLower(keyword)) {
			fmt.Printf("ID: %d, Course ID: %d, Topic: %s, Description: %s\n", forum.ID, forum.CourseID, forum.Topic, forum.Description)
			found = true
		}
	}

	if !found {
		fmt.Println("Forum tidak ditemukan.")
	}
}

func hapusData(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("              Hapus Data")
	fmt.Println("====================================")
	fmt.Println("1. Hapus Course")
	fmt.Println("2. Hapus Tugas")
	fmt.Println("3. Hapus Quiz")
	fmt.Println("4. Hapus Forum")
	fmt.Print("Pilih data yang ingin dihapus: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Print("Masukkan ID course yang ingin dihapus: ")
		var courseID int
		fmt.Scan(&courseID)
		for i, course := range lmsData.Courses {
			if course.ID == courseID {
				lmsData.Courses = append(lmsData.Courses[:i], lmsData.Courses[i+1:]...)
				fmt.Println("Course berhasil dihapus.")
				fmt.Println("Loading......")
				time.Sleep(3 * time.Second)
				return
			}
		}
		fmt.Println("Course tidak ditemukan.")
	case "2":
		fmt.Print("Masukkan ID tugas yang ingin dihapus: ")
		var taskID int
		fmt.Scan(&taskID)
		for i, task := range lmsData.Tasks {
			if task.ID == taskID {
				lmsData.Tasks = append(lmsData.Tasks[:i], lmsData.Tasks[i+1:]...)
				fmt.Println("Tugas berhasil dihapus.")
				fmt.Println("Loading......")
				time.Sleep(3 * time.Second)
				return
			}
		}
		fmt.Println("Tugas tidak ditemukan.")
	case "3":
		fmt.Print("Masukkan ID quiz yang ingin dihapus: ")
		var quizID int
		fmt.Scan(&quizID)
		for i, quiz := range lmsData.Quizzes {
			if quiz.ID == quizID {
				lmsData.Quizzes = append(lmsData.Quizzes[:i], lmsData.Quizzes[i+1:]...)
				fmt.Println("Quiz berhasil dihapus.")
				fmt.Println("Loading......")
				time.Sleep(3 * time.Second)
				return
			}
		}
		fmt.Println("Quiz tidak ditemukan.")
	case "4":
		fmt.Print("Masukkan ID forum yang ingin dihapus: ")
		var forumID int
		fmt.Scan(&forumID)
		for i, forum := range lmsData.Forums {
			if forum.ID == forumID {
				lmsData.Forums = append(lmsData.Forums[:i], lmsData.Forums[i+1:]...)
				fmt.Println("Forum berhasil dihapus.")
				fmt.Println("Loading......")
				time.Sleep(3 * time.Second)
				return
			}
		}
		fmt.Println("Forum tidak ditemukan.")
	default:
		fmt.Println("Pilihan tidak valid.")
	}
	fmt.Println("Loading......")
	time.Sleep(3 * time.Second)
}

func lihatJawabanMahasiswa() {
	fmt.Println("\n====================================")
	fmt.Println("        Jawaban Mahasiswa")
	fmt.Println("====================================")
	if len(studentAnswers) == 0 {
		fmt.Println("Belum ada jawaban yang dikirimkan.")
		return
	}

	for _, answer := range studentAnswers {
		fmt.Printf("ID Jawaban: %d\nID Tugas/Quiz: %d\nNama Mahasiswa: %s\nJawaban:\n%s\nGraded: %v\nGrade: %s\n------------------------------------\n",
			answer.ID, answer.TaskID, answer.StudentName, answer.Answer, answer.Graded, answer.Grade)
	}
}

func beriNilai(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Beri Nilai")
	fmt.Println("====================================")
	fmt.Print("Masukkan ID jawaban yang akan dinilai: ")

	var answerID int
	if _, err := fmt.Scan(&answerID); err != nil {
		fmt.Println("Input tidak valid. Harap masukkan angka.")
		return
	}

	// Membersihkan buffer setelah fmt.Scan
	reader.ReadString('\n')

	for i, answer := range studentAnswers {
		if answer.ID == answerID {
			fmt.Printf("ID Tugas/Quiz: %d\nNama Mahasiswa: %s\nJawaban:\n%s\n", answer.TaskID, answer.StudentName, answer.Answer)
			fmt.Print("Masukkan nilai: ")
			grade, _ := reader.ReadString('\n')
			grade = strings.TrimSpace(grade)

			studentAnswers[i].Graded = true
			studentAnswers[i].Grade = grade

			fmt.Println("Nilai berhasil disimpan!")
			return
		}
	}
	fmt.Println("Jawaban tidak ditemukan.")
}

func mahasiswaMenu(reader *bufio.Reader, user *User) {
	for {
		fmt.Println("\n====================================")
		fmt.Println("              Menu mahasiswa")
		fmt.Println("====================================")
		fmt.Println("1. Lihat Tugas")
		fmt.Println("2. Kerjakan Tugas")
		fmt.Println("3. Kerjakan Quiz")
		fmt.Println("4. Ikut Forum Diskusi")
		fmt.Println("5. Logout")
		fmt.Println("====================================")
		fmt.Print("Pilih menu: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			lihatSemuaData()
		case "2":
			kerjakanTugas(reader, user)
		case "3":
			kerjakanQuiz(reader, user)
		case "4":
			ikutiForum(reader)
		case "5":
			fmt.Println("Logout dari menu mahasiswa.")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

func kerjakanTugas(reader *bufio.Reader, student *User) {
	fmt.Println("\n====================================")
	fmt.Println("           Kerjakan Tugas")
	fmt.Println("====================================")
	lihatSemuaData()
	fmt.Print("Masukkan ID tugas yang akan dikerjakan: ")
	var taskID int
	fmt.Scan(&taskID)

	// Buang karakter newline yang tersisa di buffer
	reader.ReadString('\n')

	for _, task := range lmsData.Tasks {
		if task.ID == taskID {
			fmt.Printf("Mengakses tugas: %s\nDescription: %s\n", task.Title, task.Description)
			fmt.Print("Masukkan jawaban Anda: ")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)

			// Simpan jawaban
			studentAnswers = append(studentAnswers, StudentAnswer{
				ID:          nextAnswerID,
				TaskID:      taskID,
				StudentName: student.Username,
				Answer:      answer,
				Graded:      false,
				Grade:       "",
			})
			nextAnswerID++

			fmt.Println("Jawaban Anda berhasil disimpan!")
			return
		}
	}
	fmt.Println("Tugas tidak ditemukan.")
}

func kerjakanQuiz(reader *bufio.Reader, student *User) {
	fmt.Println("\n====================================")
	fmt.Println("           Kerjakan Quiz")
	fmt.Println("====================================")
	lihatSemuaData()
	fmt.Print("Masukkan ID quiz yang akan dikerjakan: ")
	var quizID int
	fmt.Scan(&quizID)

	for _, quiz := range lmsData.Quizzes {
		if quiz.ID == quizID {
			fmt.Printf("Mengakses quiz: %s\n", quiz.Title)
			answers := ""
			for i, question := range quiz.Questions {
				fmt.Printf("%d. %s\n", i+1, question)
				fmt.Print("Jawaban: ")
				answer, _ := reader.ReadString('\n')
				answers += fmt.Sprintf("Q%d: %s\n", i+1, strings.TrimSpace(answer))
			}

			// Simpan jawaban quiz
			studentAnswers = append(studentAnswers, StudentAnswer{
				ID:          nextAnswerID,
				TaskID:      quizID,
				StudentName: student.Username,
				Answer:      answers,
				Graded:      false,
				Grade:       "",
			})
			nextAnswerID++

			fmt.Println("Jawaban quiz Anda berhasil disimpan!")
			return
		}
	}
	fmt.Println("Quiz tidak ditemukan.")
}

func ikutiForum(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Ikuti Forum Diskusi")
	fmt.Println("====================================")
	lihatSemuaData()
	fmt.Print("Masukkan ID forum yang akan diikuti: ")
	var forumID int
	fmt.Scan(&forumID)

	// Membersihkan buffer untuk menghindari sisa karakter
	reader.ReadString('\n')

	for i, forum := range lmsData.Forums {
		if forum.ID == forumID {
			fmt.Printf("Mengakses forum: %s\n", forum.Topic)
			fmt.Println("Post sebelumnya:")
			for _, post := range forum.Posts {
				fmt.Printf("- %s\n", post)
			}

			// Tambahkan post baru
			fmt.Print("Masukkan post baru: ")
			newPost, _ := reader.ReadString('\n')
			newPost = strings.TrimSpace(newPost)
			lmsData.Forums[i].Posts = append(lmsData.Forums[i].Posts, newPost)

			fmt.Println("Post berhasil ditambahkan ke forum.")

			// Simpan perubahan ke file
			saveData()

			fmt.Println("Loading......")
			time.Sleep(3 * time.Second)
			return
		}
	}

	fmt.Println("Forum tidak ditemukan.")
	fmt.Println("Loading......")
	time.Sleep(3 * time.Second)
}

func lihatSemuaData() {
	fmt.Println("\n====================================")
	fmt.Println("           Semua Data LMS")
	fmt.Println("====================================")
	fmt.Println("Courses:")
	fmt.Println("------------------------------------")
	if len(lmsData.Courses) == 0 {
		fmt.Println("Belum ada course yang tersedia.")
	} else {
		for _, course := range lmsData.Courses {
			fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\n------------------------------------\n", course.ID, course.Title, course.Description)
		}
	}

	fmt.Println("\nTasks:")
	fmt.Println("------------------------------------")
	if len(lmsData.Tasks) == 0 {
		fmt.Println("Belum ada tugas yang tersedia.")
	} else {
		for _, task := range lmsData.Tasks {
			fmt.Printf("ID: %d\nCourse ID: %d\nTitle: %s\nDescription: %s\n------------------------------------\n", task.ID, task.CourseID, task.Title, task.Description)
		}
	}

	fmt.Println("\nQuizzes:")
	fmt.Println("------------------------------------")
	if len(lmsData.Quizzes) == 0 {
		fmt.Println("Belum ada quiz yang tersedia.")
	} else {
		for _, quiz := range lmsData.Quizzes {
			fmt.Printf("ID: %d\nCourse ID: %d\nTitle: %s\nQuestions: %v\n------------------------------------\n", quiz.ID, quiz.CourseID, quiz.Title, quiz.Questions)
		}
	}

	fmt.Println("\nForums:")
	fmt.Println("------------------------------------")
	if len(lmsData.Forums) == 0 {
		fmt.Println("Belum ada forum yang tersedia.")
	} else {
		for _, forum := range lmsData.Forums {
			fmt.Printf("ID: %d\nCourse ID: %d\nTopic: %s\nDeskripsi: %v\n------------------------------------\n", forum.ID, forum.CourseID, forum.Topic, forum.Description)
		}
	}
	fmt.Println("====================================")
}

func tambahCourse(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Tambah Course")
	fmt.Println("====================================")
	fmt.Print("Masukkan judul course: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Masukkan deskripsi course: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	course := Course{
		ID:          nextCourseID,
		Title:       title,
		Description: description,
	}

	lmsData.Courses = append(lmsData.Courses, course)
	nextCourseID++
	fmt.Println("\nCourse berhasil ditambahkan.")
	fmt.Println("====================================")
}

func tambahTugas(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Tambah Tugas")
	fmt.Println("====================================")
	fmt.Print("Masukkan ID course untuk tugas ini: ")
	var courseID int
	_, err := fmt.Scanf("%d\n", &courseID)
	if err != nil {
		fmt.Println("Input tidak valid. Silakan coba lagi.")
		return
	}

	fmt.Print("Masukkan judul tugas: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Masukkan deskripsi tugas: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	task := Task{
		ID:          nextTaskID,
		CourseID:    courseID,
		Title:       title,
		Description: description,
	}

	lmsData.Tasks = append(lmsData.Tasks, task)
	nextTaskID++
	fmt.Println("\nTugas berhasil ditambahkan.")
	fmt.Println("====================================")
}

func tambahQuiz(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Tambah Quiz")
	fmt.Println("====================================")
	fmt.Print("Masukkan ID course untuk quiz ini: ")
	var courseID int
	fmt.Scan(&courseID)

	fmt.Print("Masukkan judul quiz: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	questions := []string{}
	for {
		fmt.Print("Masukkan pertanyaan (atau ketik 'selesai' untuk selesai): ")
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)
		if question == "selesai" {
			break
		}
		questions = append(questions, question)
	}

	quiz := Quiz{
		ID:        nextQuizID,
		CourseID:  courseID,
		Title:     title,
		Questions: questions,
	}

	lmsData.Quizzes = append(lmsData.Quizzes, quiz)
	nextQuizID++
	fmt.Println("\nQuiz berhasil ditambahkan.")
	fmt.Println("====================================")
}

func tambahForum(reader *bufio.Reader) {
	fmt.Println("\n====================================")
	fmt.Println("           Tambah Forum Diskusi")
	fmt.Println("====================================")

	// Meminta ID course
	fmt.Print("Masukkan ID course untuk forum ini: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)
	var courseID int
	fmt.Sscanf(idInput, "%d", &courseID)

	// Meminta topik forum
	fmt.Print("Masukkan topik forum: ")
	topic, _ := reader.ReadString('\n')
	topic = strings.TrimSpace(topic)

	// Meminta deskripsi forum
	fmt.Print("Masukkan deskripsi forum: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Membuat forum baru
	forum := Forum{
		ID:          nextForumID,
		CourseID:    courseID,
		Topic:       topic,
		Description: description,
		Posts:       []string{},
	}

	// Menyimpan forum ke data
	lmsData.Forums = append(lmsData.Forums, forum)
	nextForumID++
	fmt.Println("\nForum diskusi berhasil ditambahkan.")
	fmt.Println("====================================")
}

func loadData() {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // Jika file tidak ada, abaikan
		}
		fmt.Println("Gagal membuka file data:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&lmsData)
	if err != nil {
		fmt.Println("Gagal decode data dari file:", err)
	}

	if len(lmsData.Courses) > 0 {
		nextCourseID = lmsData.Courses[len(lmsData.Courses)-1].ID + 1
	}
	if len(lmsData.Tasks) > 0 {
		nextTaskID = lmsData.Tasks[len(lmsData.Tasks)-1].ID + 1
	}
	if len(lmsData.Quizzes) > 0 {
		nextQuizID = lmsData.Quizzes[len(lmsData.Quizzes)-1].ID + 1
	}
	if len(lmsData.Forums) > 0 {
		nextForumID = lmsData.Forums[len(lmsData.Forums)-1].ID + 1
	}
}

func saveData() {
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Gagal menyimpan data:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&lmsData)
	if err != nil {
		fmt.Println("Gagal encode data ke file:", err)
	}
}
