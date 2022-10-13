package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (sp *SpUser) Display() {
	fmt.Println("Name: ", sp.Name)
	fmt.Println("Username: ", sp.Username)
	fmt.Println("Password: ", sp.Password)
	fmt.Println("Email: ", sp.Email)
	fmt.Println("Role: ", sp.Role)
	fmt.Println("Sex: ", sp.Sex)
	fmt.Println("Sessions Assigned: ")
	for i := 0; i < len(sp.SessionsAssigned); i++ {
		fmt.Println(*sp.SessionsAssigned[i])
	}
	fmt.Println("Sessions Available: ")
	for i := 0; i < len(sp.SessionsAvailable); i++ {
		fmt.Println(*sp.SessionsAvailable[i])
	}
	fmt.Println("Sessions Unavailable: ")
	for i := 0; i < len(sp.SessionsUnavailable); i++ {
		fmt.Println(*sp.SessionsUnavailable[i])
	}
}
func (sp *SpUser) MakeRecord(db *sql.DB) (err error) {
	statement := "insert into spusers (name, username, role, email, sex, sessionsassigned, sessionsavailable, sessionsunavailable, password) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	sessionsAssignedByte, err := json.Marshal(&sp.SessionsAssigned)
	sessionsAvailableByte, err := json.Marshal(&sp.SessionsAvailable)
	sessionsUnavailableByte, err := json.Marshal(&sp.SessionsUnavailable)

	err = stmt.QueryRow(sp.Name, sp.Username, sp.Role, sp.Email, sp.Sex, sessionsAssignedByte,
		sessionsAvailableByte, sessionsUnavailableByte, sp.Password).Scan(&sp.Id)

	if err != nil {
		return
	}
	return
}
func GetSpUserRecord(username string, db *sql.DB) (sp SpUser, err error) {
	sp = SpUser{}

	var sessionsAvailableByte []byte
	var sessionsUnavailableByte []byte
	var sessionsAssignedByte []byte

	err = db.QueryRow("select Id, name, username, role, email, sex, sessionsassigned, "+
		"sessionsavailable, sessionsunavailable, password "+
		"from spusers where username = $1 ", username).Scan(&sp.Id, &sp.Name, &sp.Username, &sp.Role, &sp.Email,
		&sp.Sex, &sessionsAssignedByte,
		&sessionsAvailableByte, &sessionsUnavailableByte,
		&sp.Password)

	if err != nil {
		return
	}

	err = json.Unmarshal(sessionsAssignedByte, &sp.SessionsAssigned)
	if err != nil {
		return
	}
	err = json.Unmarshal(sessionsAvailableByte, &sp.SessionsAvailable)
	if err != nil {
		return
	}
	err = json.Unmarshal(sessionsUnavailableByte, &sp.SessionsUnavailable)
	if err != nil {
		return
	}
	return
}
func (sp *SpUser) UpdateRecord(db *sql.DB) (err error) {

	sessionsAssignedByte, err := json.Marshal(&sp.SessionsAssigned)
	if err != nil {
		return
	}
	sessionsAvailableByte, err := json.Marshal(&sp.SessionsAvailable)
	if err != nil {
		return
	}
	sessionsUnavailableByte, err := json.Marshal(&sp.SessionsUnavailable)
	if err != nil {
		return
	}

	_, err = db.Exec("update spusers set sessionsavailable = $2, "+
		"sessionsunavailable = $3, sessionsassigned = $4, "+
		"password = $5, email = $6 where username = $1 ",
		sp.Username, sessionsAvailableByte, sessionsUnavailableByte,
		sessionsAssignedByte, sp.Password, sp.Email)

	return
}
func (sp *SpUser) DeleteRecord(db *sql.DB) (err error) {
	_, err = db.Exec("delete from spusers where username = $1", sp.Username)
	return
}

func (sp *SpManager) Display() {
	fmt.Println("Name: ", sp.Name)
	fmt.Println("Username: ", sp.Username)
	fmt.Println("Password: ", sp.Password)
	fmt.Println("Email: ", sp.Email)
	fmt.Println("Role: ", sp.Role)
	fmt.Println("Assigned Patients: ")
	for i := 0; i < len(sp.AssignedPatients); i++ {
		fmt.Println(*sp.AssignedPatients[i])
	}
	fmt.Println("Sessions Available: ")
	for i := 0; i < len(sp.SessionsManaged); i++ {
		fmt.Println(*sp.SessionsManaged[i])
	}
	fmt.Println("Sessions Unavailable: ")
	for i := 0; i < len(sp.SessionsUnmanaged); i++ {
		fmt.Println(*sp.SessionsUnmanaged[i])
	}
}
func (sp *SpManager) MakeRecord(db *sql.DB) (err error) {
	statement := "insert into spmanagers (name, username, role, password, email, assignedpatients, sessionsmanaged, sessionsunmanaged) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	assignedPatientsByte, err := json.Marshal(&sp.AssignedPatients)
	sessionsManagedByte, err := json.Marshal(&sp.SessionsManaged)
	sessionsUnmanagedByte, err := json.Marshal(&sp.SessionsUnmanaged)

	err = stmt.QueryRow(sp.Name, sp.Username, sp.Role, sp.Password, sp.Email, assignedPatientsByte,
		sessionsManagedByte, sessionsUnmanagedByte).Scan(&sp.Id)

	if err != nil {
		return
	}
	return
}
func GetSpManagerRecord(username string, db *sql.DB) (sp SpManager, err error) {
	sp = SpManager{}

	var sessionsManagedByte []byte
	var sessionsUnmanagedByte []byte
	var assignedPatientsByte []byte

	err = db.QueryRow("select Id, name, username, role, email, password, assignedpatients, "+
		"sessionsmanaged, sessionsunmanaged "+
		"from spmanagers where username = $1 ", username).Scan(&sp.Id, &sp.Name, &sp.Username, &sp.Role, &sp.Email,
		&sp.Password, &assignedPatientsByte,
		&sessionsManagedByte, &sessionsUnmanagedByte)

	if err != nil {
		return
	}

	err = json.Unmarshal(assignedPatientsByte, &sp.AssignedPatients)
	if err != nil {
		return
	}
	err = json.Unmarshal(sessionsManagedByte, &sp.SessionsManaged)
	if err != nil {
		return
	}
	err = json.Unmarshal(sessionsUnmanagedByte, &sp.SessionsUnmanaged)
	if err != nil {
		return
	}
	return
}
func (sp *SpManager) UpdateRecord(db *sql.DB) (err error) {
	assignedPatientsByte, err := json.Marshal(&sp.AssignedPatients)
	if err != nil {
		return
	}
	sessionsManagedByte, err := json.Marshal(&sp.SessionsManaged)
	if err != nil {
		return
	}
	sessionsUnmanagedByte, err := json.Marshal(&sp.SessionsUnmanaged)
	if err != nil {
		return
	}

	_, err = db.Exec("update spmanagers set sessionsmanaged = $2, "+
		"sessionsunmanaged = $3, assignedpatients = $4, "+
		"password = $5, email = $6 where username = $1 ",
		sp.Username, sessionsManagedByte, sessionsUnmanagedByte,
		assignedPatientsByte, sp.Password, sp.Email)

	return
}
func (sp *SpManager) DeleteRecord(db *sql.DB) (err error) {
	_, err = db.Exec("delete from spmanagers where username = $1", sp.Username)
	return
}

func (s *Session) Display() {
	fmt.Println("Date: ", s.Information.Date)
	fmt.Println("Time: ", s.Information.Time)
	fmt.Println("Duration: ", s.Information.Duration)
	fmt.Println("Location: ", s.Information.Location)
	fmt.Println("Description: ", s.Information.Description)
	fmt.Println("Instructors: ")
	for i := 0; i < len(s.Instructors); i++ {
		fmt.Println(*s.Instructors[i])
	}
	fmt.Println("PatientsNeeded: ", s.PatientsNeeded)

	fmt.Println("Patients Assigned: ")
	for i := 0; i < len(s.PatientsAssigned); i++ {
		fmt.Println(*s.PatientsAssigned[i])
	}
	fmt.Println("Patients Available: ")
	for i := 0; i < len(s.PatientsAvailable); i++ {
		fmt.Println(*s.PatientsAvailable[i])
	}
	fmt.Println("Patients Unavailable: ")
	for i := 0; i < len(s.PatientsUnavailable); i++ {
		fmt.Println(*s.PatientsUnavailable[i])
	}

	fmt.Println("Patients No Response: ")
	for i := 0; i < len(s.PatientsNoResponse); i++ {
		fmt.Println(*s.PatientsNoResponse[i])
	}
}
func (s *Session) MakeRecord(db *sql.DB) (err error) {
	statement := "insert into sessions (date, time, duration, location, description, instructors, patientsneeded, patientsassigned, patientsavailable, patientsunavailable, patientsnoresponse) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	instructorsByte, err := json.Marshal(&s.Instructors)
	patientsAssignedByte, err := json.Marshal(&s.PatientsAssigned)
	patientsAvailableByte, err := json.Marshal(&s.PatientsAvailable)
	patientsUnavailableByte, err := json.Marshal(&s.PatientsUnavailable)
	patientsNoResponse, err := json.Marshal(&s.PatientsNoResponse)

	err = stmt.QueryRow(s.Information.Date, s.Information.Time, s.Information.Duration,
		s.Information.Location, s.Information.Description,
		instructorsByte, s.PatientsNeeded,
		patientsAssignedByte, patientsAvailableByte,
		patientsUnavailableByte,
		patientsNoResponse).Scan(&s.Id)

	if err != nil {
		return
	}
	return
}
func GetSessionRecord(id uint, db *sql.DB) (s Session, err error) {
	s = Session{
		Information: &SessionInfo{}}

	var instructorsByte []byte
	var patientsAssignedByte []byte
	var patientsAvailableByte []byte
	var patientsUnavailableByte []byte
	var patientsNoResponseByte []byte

	err = db.QueryRow("select id, date, time, duration, location, description, "+
		"instructors, patientsneeded, patientsassigned, patientsavailable, patientsunavailable, patientsnoresponse "+
		"from sessions where id = $1 ", id).Scan(&s.Id, &s.Information.Date, &s.Information.Time, &s.Information.Duration, &s.Information.Location,
		&s.Information.Description, &instructorsByte, &s.PatientsNeeded,
		&patientsAssignedByte, &patientsAvailableByte, &patientsUnavailableByte, &patientsNoResponseByte)

	if err != nil {
		return
	}

	err = json.Unmarshal(instructorsByte, &s.Instructors)
	if err != nil {
		return
	}
	err = json.Unmarshal(patientsAssignedByte, &s.PatientsAssigned)
	if err != nil {
		return
	}
	err = json.Unmarshal(patientsAvailableByte, &s.PatientsAvailable)
	if err != nil {
		return
	}
	err = json.Unmarshal(patientsUnavailableByte, &s.PatientsUnavailable)
	if err != nil {
		return
	}
	err = json.Unmarshal(patientsNoResponseByte, &s.PatientsNoResponse)
	if err != nil {
		return
	}
	return
}
func (s *Session) UpdateRecord(db *sql.DB) (err error) {
	patientsAssignedByte, err := json.Marshal(&s.PatientsAssigned)
	if err != nil {
		return
	}
	patientsAvailableByte, err := json.Marshal(&s.PatientsAvailable)
	if err != nil {
		return
	}
	patientsUnavailableByte, err := json.Marshal(&s.PatientsUnavailable)
	if err != nil {
		return
	}

	patientsNoResponseByte, err := json.Marshal(&s.PatientsNoResponse)
	if err != nil {
		return
	}
	instructorsByte, err := json.Marshal(&s.Instructors)
	if err != nil {
		return
	}
	_, err = db.Exec("update sessions set date = $2, "+
		"time = $3, duration = $4, "+
		"location = $5, description = $6, instructors = $7, patientsneeded = $8, patientsassigned = $9, patientsavailable = $10, patientsunavailable = $11, patientsnoresponse = $12 where id = $1",
		s.Id, s.Information.Date, s.Information.Time, s.Information.Duration,
		s.Information.Location, s.Information.Description,
		instructorsByte, s.PatientsNeeded,
		patientsAssignedByte, patientsAvailableByte,
		patientsUnavailableByte, patientsNoResponseByte)

	return
}
func (s *Session) DeleteRecord(db *sql.DB) (err error) {
	_, err = db.Exec("delete from sessions where id = $1", s.Id)
	return
}

func testSpUser(db *sql.DB) (err error) {
	spuser := SpUser{
		Name:     *Name{}.Create("Robert Pike"),
		Username: "rpike",
		Role:     SP,
		Email:    "robert@pike.com",
		Password: "letmein",
		Sex:      Male,
	}

	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-UP")
	session2 := Session{}.Create("10/14/2022", "10:30AM", "1H", "Park Ave", "Follow-UP")

	spuser.SessionsAssigned = append(spuser.SessionsAssigned, session.Information, session2.Information)

	// The database driver will call the Value() method and and marshall the
	// attrs struct to JSON before the INSERT.
	err = spuser.MakeRecord(db)

	if err != nil {
		fmt.Println("Error in Create Record:", err)
		return
	}

	fmt.Println("Created Record -> ", spuser.Name)
	fmt.Println()

	spuser, err = GetSpUserRecord("rpike", db)

	if err != nil {
		fmt.Println("Error in main:", err)
		return
	}

	fmt.Println("Retrieved Record -> ")
	spuser.Display()
	spuser.Email = "rpike@mousemail.com"
	err = spuser.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error in Update Main", err)
		return
	}
	fmt.Println()

	fmt.Println("Updated Record -> ")
	spuser.Display()

	fmt.Println()
	err = spuser.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error in Deletion of Record in Main: ", err)
		return
	}
	fmt.Println("Deleted Record -> ", spuser.Username)

	return
}
func testSpManager(db *sql.DB) (err error) {
	spmanager := SpManager{
		Name:     *Name{}.Create("Lisa Thomas"),
		Username: "lthomas",
		Role:     Manager,
		Password: "gxpt442!$%",
		Email:    "lisa@thomas.com",
	}

	err = spmanager.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Make Record Sp Manager in Main", err)
	}

	fmt.Println("Created Record -> ", spmanager.Name)
	fmt.Println()

	spmanager, err = GetSpManagerRecord("lthomas", db)

	if err != nil {
		fmt.Println("Error in retrieving SpManager record from DB", err)
	}

	fmt.Println("Retrieved Record ->")
	spmanager.Display()
	fmt.Println()

	spuser := SpUser{
		Name:     *Name{}.Create("Robert Pike"),
		Username: "rpike",
		Role:     SP,
		Email:    "robert@pike.com",
		Password: "letmein",
		Sex:      Male,
	}

	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-UP")
	session2 := Session{}.Create("10/14/2022", "10:30AM", "1H", "Park Ave", "Follow-UP")

	spuser.SessionsAssigned = append(spuser.SessionsAssigned, session.Information, session2.Information)

	spmanager.Email = "lthomas@duck.com"
	spmanager.AssignedPatients = append(spmanager.AssignedPatients, &spuser)
	err = spmanager.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error with Updating Sp Manager", err)
		return
	}

	fmt.Println("Updated Record -> ")
	spmanager.Display()
	fmt.Println()

	err = spmanager.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error Deleting Record Sp Manager in Main", err)
	}
	fmt.Println("Deleted Record -> ", spmanager.Username)
	return
}

func testSession(db *sql.DB) (err error) {
	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-Up")
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Make session: ", err)
	}
	fmt.Println()
	fmt.Println("Created Record ->", *session.Information)
	fmt.Println()

	session2, err := GetSessionRecord(1, db)

	if err != nil {
		fmt.Println("Error in main: ", err)
	}

	fmt.Println("Record Retrieved -> ")
	session2.Display()
	fmt.Println()

	session2.Information.Date = "12/15/2022"
	session2.Information.Description = "Follow Up"
	session2.Information.Location = "Park Ave"
	err = session2.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error Updating Record in Main", err)
		return
	}
	fmt.Println("Updated Record ->: ")
	session2.Display()
	fmt.Println()

	err = session2.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error Deleting Record in Main", err)
		return
	}
	fmt.Println("Deleted Record ->", *session2.Information)
	return
}

func TestDB() {
	db, err := sql.Open("postgres", "user=postgres dbname=sp_calendar password=rxpt221!@# sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = testSpUser(db)
	if err != nil {
		fmt.Println("Error test db user in main: ", err)
	}

	err = testSpManager(db)
	if err != nil {
		fmt.Println("Error test db user in main: ", err)
	}

	err = testSession(db)
	if err != nil {
		fmt.Println("Error test db session in main: ", err)
	}

}
