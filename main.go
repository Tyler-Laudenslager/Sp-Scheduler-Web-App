package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("dinosaursarecool")
	store = sessions.NewCookieStore(key)
)

func makeSP(name, username, password string) {
	var err error
	hashedPassword, _ := HashPassword(password)
	spuser := SpUser{
		Name:     *Name{}.Create(name),
		Username: username,
		Role:     SP,
		Password: hashedPassword,
		Email:    username + "@duck.com",
	}

	err = spuser.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
}

func makeManager(name, username, password string) {
	var err error
	hashedPassword, _ := HashPassword(password)
	spmanager3 := SpManager{
		Name:     *Name{}.Create(name),
		Username: username,
		Role:     Manager,
		Password: hashedPassword,
	}

	spmanager3.AssignedPatients, err = GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error in Init Get All Sp User Records: ", err)
	}
	err = spmanager3.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Record -> ", err)
		return
	}
}

func resetAllSpPasswords() {
	allSpUsers, _ := GetAllSpUserRecords(db)
	for _, su := range allSpUsers {
		hashedPassword, _ := HashPassword("Stlukes800!")
		su.Password = hashedPassword
		su.UpdateRecord(db)
	}
}

func makeSession(name, date, starttime, endtime, location, description string) {

	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}

	session := Session{}.Create(name, date, starttime, endtime, location, description)
	session.Information.CreatedDate = time.Now().In(loc).Format("01/02/2006")
	session.Information.ExpiredDate = time.Now().In(loc).Format("01/02/2006")
	session.PatientsNeeded = 3
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Making Session Record 1: ", err)
	}
}

func saveDatabase() {
	SpUsersBox := make(SpUsersBox, 0)
	SpManagersBox := make(SpManagersBox, 0)
	SessionsBox := make(SpSessionsBox, 0)

	allSpUsers, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error getting all sp user records in json: ", err)
	}
	allManagers, err := GetAllSpManagerRecords(db)
	if err != nil {
		fmt.Println("Error getting all sp manager records in json: ", err)
	}
	allSessions, err := GetAllSessionRecords(db)
	if err != nil {
		fmt.Println("Error getting all session records in json: ", err)
	}

	for _, spuser := range allSpUsers {
		SpUsersBox = append(SpUsersBox, spuser)
	}

	for _, spmanager := range allManagers {
		SpManagersBox = append(SpManagersBox, spmanager)
	}

	for _, session := range allSessions {
		SessionsBox = append(SessionsBox, session)
	}

	HospitalCalendar := HospitalCalendar{
		Users:    SpUsersBox,
		Managers: SpManagersBox,
		Sessions: SessionsBox,
	}

	output, err := json.MarshalIndent(HospitalCalendar, "", "\t\t")
	if err != nil {
		fmt.Println("Error Marshal Indent in Main")
	}
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	timenow := time.Now().In(loc)
	date := timenow.Format("01_02_2006")
	fileName := "backup_session_data_" + date + ".json"

	os.WriteFile(fileName, output, 0600)
}

// This function will run on initilization of program
func init() {
	//resetAllSpPasswords()

	// Session Creation
	makeSession("Sacred Heart Check-UP", "06/12/2023", "11:00am", "12:00pm", "Sacred Heart", "Check-Up")
	makeSession("Anderson Follow UP", "06/20/2023", "1:00pm", "2:00pm", "Warren", "Follow-Up")
	makeSession("Practice Session", "06/30/2023", "1:30pm", "3:00pm", "Sacred Heart", "Follow-Up")
	makeSession("Practice Session", "06/29/2023", "7:30am", "8:30am", "Sacred Heart", "Follow-Up")
	// SP Creation
	makeSP("Charles Darwin", "cdarwin", "letmein2")
	makeSP("Robert Pike", "rpike", "letmein")

	// Manager Creation
	// saveDatabase()
	makeManager("Tyler Lauden", "tlaud", "letmeinman")

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
	http.HandleFunc("/confirmAllSPs", confirmAllSPs)
	http.HandleFunc("/selectedToAssigned", selectedToAssigned)
	http.HandleFunc("/makeSessionsAvailable", makeSessionsAvailable)
	http.HandleFunc("/addcomment", addcomment)
	http.HandleFunc("/assignsp", assignsp)
	http.HandleFunc("/signupavailable", signupavailable)
	http.HandleFunc("/signupnotavailable", signupnotavailable)
	http.HandleFunc("/sessionbackup", sessionbackup)
	http.HandleFunc("/changeemail", changeemail)
	http.HandleFunc("/changepassword", changepassword)
	http.HandleFunc("/createSPRecord", createSPRecord)
	http.HandleFunc("/deleteSPRecord", deleteSPRecord)
	http.HandleFunc("/toggleshowsession", toggleshowsession)
	http.HandleFunc("/togglehourglass", togglehourglass)
	http.HandleFunc("/togglechecksquare", togglechecksquare)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/", login)
	server.ListenAndServe()
	// server.ListenAndServeTLS("/etc/letsencrypt/live/sluhnspcalendar.com/fullchain.pem",
	//                          "/etc/letsencrypt/live/sluhnspcalendar.com/privkey.pem")
}
