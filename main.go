package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func init() {
	hashedPassword, err := HashPassword("letmein")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	hashedPassword2, err := HashPassword("rxpt221!@#")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Sacred Heart", "Check-Up")
	session.PatientsNeeded = 6
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 1: ", err)
	}
	fmt.Println("Created Session -> ", session.Information)
	session2 := Session{}.Create("12/25/2022", "12:00AM", "2H", "Anderson", "Follow-Up")
	session2.PatientsNeeded = 4
	err = session2.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 2: ", err)
	}
	fmt.Println("Created Session -> ", session2.Information)
	session3 := Session{}.Create("1/25/2023", "1:00PM", "3H", "Allentown", "Skills Workshop")
	session3.PatientsNeeded = 2
	err = session3.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 3: ", err)
	}
	fmt.Println("Created Session -> ", session3.Information)
	session4 := Session{}.Create("2/25/2024", "2:00PM", "4H", "Anderson", "ED Skills Assessment")
	session4.PatientsNeeded = 3
	err = session4.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 4: ", err)
	}
	fmt.Println("Created Session -> ", session4.Information)
	spuser := SpUser{
		Name:     *Name{}.Create("Robert Pike"),
		Username: "rpike",
		Role:     SP,
		Sex:      Male,
		Password: hashedPassword,
		Email:    "rpike@duck.com",
	}

	spmanager := SpManager{
		Name:     *Name{}.Create("Emily Garey"),
		Username: "egarey",
		Role:     Manager,
		Password: hashedPassword2,
		Email:    "egarey@duck.com",
	}

	err = spuser.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spuser.Name)

	err = spmanager.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spmanager.Name)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:6600",
	}

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/json", sendjson)
	http.HandleFunc("/signupavailable", signupavailable)
	http.HandleFunc("/signupnotavailable", signupnotavailable)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/", login)
	server.ListenAndServe()
}
