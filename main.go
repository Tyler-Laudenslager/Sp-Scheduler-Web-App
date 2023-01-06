package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("dinosaursarecool")
	store = sessions.NewCookieStore(key)
)

func init() {
	hashedPassword, err := HashPassword("letmein")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	hashedPassword2, err := HashPassword("letmeinman")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	hashedPassword3, err := HashPassword("letmein2")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	hashedPassword5, err := HashPassword("letmeinman2")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	session := Session{}.Create("Sacred Heart Check-UP", "01/25/2023", "11:00AM", "12:00PM", "Sacred Heart", "Check-Up")
	session.Information.CreatedDate = time.Now().Format("01/02/2006")
	session.Information.ExpiredDate = time.Now().AddDate(0, 0, 7).Format("01/02/2006")
	session.PatientsNeeded = 6
	session.Instructors = append(session.Instructors, Instructor{}.Create("Joe Thompson", "Director"))
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 1: ", err)
	}
	fmt.Println("Created Session -> ", session.Information)
	session2 := Session{}.Create("Anderson Follow UP", "02/25/2023", "12:00PM", "2:00PM", "Warren", "Follow-Up")
	session2.Information.CreatedDate = time.Now().Format("01/02/2006")
	session2.Information.ExpiredDate = time.Now().AddDate(0, 0, 7).Format("01/02/2006")
	session2.PatientsNeeded = 4
	err = session2.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 2: ", err)
	}
	fmt.Println("Created Session -> ", session2.Information)
	session3 := Session{}.Create("Allentown Skills Workshop", "03/25/2023", "1:00PM", "4:00PM", "Allentown", "Skills Workshop")
	session3.Information.CreatedDate = time.Now().Format("01/02/2006")
	session3.Information.ExpiredDate = time.Now().AddDate(0, 0, 7).Format("01/02/2006")
	session3.PatientsNeeded = 2
	err = session3.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 3: ", err)
	}
	fmt.Println("Created Session -> ", session3.Information)
	session4 := Session{}.Create("Anderson ED Skills", "04/25/2023", "2:00PM", "6:00PM", "Anderson", "ED Skills Assessment")
	session4.Information.CreatedDate = time.Now().Format("01/02/2006")
	session4.Information.ExpiredDate = time.Now().AddDate(0, 0, 7).Format("01/02/2006")
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
		Password: hashedPassword,
		Email:    "rpike@duck.com",
	}

	spuser2 := SpUser{
		Name:     *Name{}.Create("Charles Darwin"),
		Username: "cdarwin",
		Role:     SP,
		Password: hashedPassword3,
		Email:    "cdarwin@duck.com",
	}

	spmanager := SpManager{
		Name:     *Name{}.Create("Emily Garey"),
		Username: "gareye",
		Role:     Manager,
		Password: hashedPassword2,
		Email:    "egarey@duck.com",
	}

	spmanager2 := SpManager{
		Name:     *Name{}.Create("Megan Augustine"),
		Username: "augustm",
		Role:     Manager,
		Password: hashedPassword5,
		Email:    "maugustine@duck.com",
	}

	err = spuser.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spuser.Name)

	err = spuser2.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spuser2.Name)

	spmanager.AssignedPatients, err = GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error in Init Get All Sp User Records: ", err)
	}
	err = spmanager.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spmanager.Name)

	spmanager2.AssignedPatients, err = GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error in Init Get All Sp User Records: ", err)
	}
	err = spmanager2.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spmanager2.Name)
}

func main() {
	//TLS Server
	// server := http.Server {
	//	Addr: ":443",
	//}
	server := http.Server{
		Addr: "127.0.0.1:6600",
	}

	http.Handle("/images/logos/", http.StripPrefix("/images/logos/", http.FileServer(http.Dir("./images/logos"))))
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/json", sendjson)
	http.HandleFunc("/createsession", createsession)
	http.HandleFunc("/updatesession", updatesession)
	http.HandleFunc("/deletesession", deletesession)
	http.HandleFunc("/assignsp", assignsp)
	http.HandleFunc("/signupavailable", signupavailable)
	http.HandleFunc("/signupnotavailable", signupnotavailable)
	http.HandleFunc("/changeemail", changeemail)
	http.HandleFunc("/changepassword", changepassword)
	http.HandleFunc("/createSPRecord", createSPRecord)
	http.HandleFunc("/deleteSPRecord", deleteSPRecord)
	http.HandleFunc("/toggleshowsession", toggleshowsession)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/", login)
	server.ListenAndServe()
}
