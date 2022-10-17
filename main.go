package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		dashboard(w, r)
		return
	}
	t, _ := template.ParseFiles("html-boilerplate.html", "login-content.html")
	t.ExecuteTemplate(w, "html-boilerplate", "")
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)

	if err != nil {
		fmt.Println("Error Get SP user record in authenticate: ", err)
		return
	}
	/* fmt.Fprintln(w, "User Authenticated with Cookie!")
	fmt.Fprintln(w, "Welcome to the Dashboard ", spuser.Name.First, spuser.Name.Last) */
	t, _ := template.ParseFiles("html-boilerplate.html", "dashboard-content.html", "session-content.html")
	session_records, err := GetAllSessionInfoRecords(db)
	if err != nil {
		fmt.Println("Error Get All Session Records: ", err)
	}

	spuser.SessionsAvailable = session_records

	t.ExecuteTemplate(w, "html-boilerplate", spuser)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	username := r.PostFormValue("userid")
	password := r.PostFormValue("password")
	spuser, err := GetSpUserRecord(username, db)
	if err != nil {
		login(w, r)
	}
	if !CheckPasswordHash(password, spuser.Password) {
		login(w, r)
	} else {
		session.Values["authenticated"] = true
		session.Values["username"] = spuser.Username
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		login(w, r)
	}
	session.Save(r, w)
	dashboard(w, r)
}

func sendjson(w http.ResponseWriter, r *http.Request) {
	bob_marcs := SpUser{}.Create(*Name{}.Create("Bob Marcs"), SP, Male, "bob@marcs.com")
	susan_miller := SpUser{}.Create(*Name{}.Create("Susan Miller"), SP, Female, "susan@miller.com")

	andy_thomas := SpManager{}.Create(*Name{}.Create("Andy Thomas"), Manager, "andy@thomas.com")

	session1 := Session{}.Create("11/15/2022", "11:00AM", "1H", "Anderson", "Check-Up")

	andy_thomas.AssignedPatients = append(andy_thomas.AssignedPatients, bob_marcs, susan_miller)
	andy_thomas.SessionsManaged = append(andy_thomas.SessionsManaged, session1)

	bob_marcs.SessionsAvailable = append(bob_marcs.SessionsAvailable, session1.Information)
	bob_marcs.SessionsAssigned = append(bob_marcs.SessionsAssigned, session1.Information)

	susan_miller.SessionsAvailable = append(susan_miller.SessionsAvailable, session1.Information)
	susan_miller.SessionsAssigned = append(susan_miller.SessionsAssigned, session1.Information)

	session1.PatientsAvailable = append(session1.PatientsAvailable, bob_marcs)
	session1.PatientsAvailable = append(session1.PatientsAvailable, susan_miller)
	session1.PatientsAssigned = append(session1.PatientsAssigned, bob_marcs)
	session1.PatientsAssigned = append(session1.PatientsAssigned, susan_miller)

	SpUsersBox := append(make(SpUsersBox, 0, 2), bob_marcs, susan_miller)
	SpManagersBox := append(make(SpManagersBox, 0, 1), andy_thomas)
	SessionsBox := append(make(SpSessionsBox, 0, 1), session1)

	HospitalCalendar := HospitalCalendar{
		Users:    SpUsersBox,
		Managers: SpManagersBox,
		Sessions: SessionsBox,
	}

	output, err := json.MarshalIndent(HospitalCalendar, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:6600",
	}
	hashedPassword, err := HashPassword("letmein")
	if err != nil {
		fmt.Println("Error Hashing Password")
		return
	}
	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-Up")
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 1: ", err)
	}
	fmt.Println("Created Session -> ", session.Information)
	session2 := Session{}.Create("12/25/2022", "12:00AM", "2H", "Anderson", "Follow-Up")
	err = session2.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 2: ", err)
	}
	fmt.Println("Created Session -> ", session2.Information)
	session3 := Session{}.Create("1/25/2023", "1:00AM", "3H", "Anderson", "Invasion")
	err = session3.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 3: ", err)
	}
	fmt.Println("Created Session -> ", session3.Information)
	session4 := Session{}.Create("2/25/2024", "2:00AM", "4H", "Anderson", "Holy-Grail")
	err = session4.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 4: ", err)
	}
	fmt.Println("Created Session -> ", session4.Information)
	spuser := SpUser{
		Name:     *Name{}.Create("Robert Pikert"),
		Username: "rpike",
		Role:     SP,
		Sex:      Male,
		Password: hashedPassword,
		Email:    "rpike@duck.com",
	}

	err = spuser.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
	fmt.Println("Created Record in Database -> ", spuser.Name)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/json", sendjson)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/", login)
	server.ListenAndServe()
}
