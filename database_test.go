package main

import (
	"database/sql"
	"fmt"
	"testing"
)

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

	spuser, err = GetSpUserRecord("rpike", db)

	if err != nil {
		fmt.Println("Error in Retrieving Record SP:", err)
		return
	}

	spuser.Email = "rpike@mousemail.com"
	err = spuser.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error in Update Record SP: ", err)
		return
	}

	err = spuser.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error in Deletion of Record SP: ", err)
		return
	}

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

	spmanager, err = GetSpManagerRecord("lthomas", db)

	if err != nil {
		fmt.Println("Error in retrieving SpManager record from DB", err)
	}

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

	err = spmanager.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error Deleting Record Sp Manager in Main", err)
	}
	return
}

func testSession(db *sql.DB) (err error) {
	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-Up")
	err = session.MakeRecord(db)
	if err != nil {
		fmt.Println("Error Make Session: ", err)
	}

	session2, err := GetSessionRecord(session.Information, db)

	if err != nil {
		fmt.Println("Error in Get Session: ", err)
	}

	session2.Information.Date = "12/15/2022"
	session2.Information.Description = "Follow Up"
	session2.Information.Location = "Park Ave"
	err = session2.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error Updating Record Session: ", err)
		return
	}

	err = session2.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error Deleting Record Session: ", err)
		return
	}
	return
}

func testGetAllSessions(db *sql.DB) (err error) {
	session := Session{}.Create("11/25/2022", "11:00AM", "1H", "Anderson", "Check-Up")
	session2 := Session{}.Create("12/25/2022", "12:00AM", "2H", "Anderson", "Follow-Up")
	session3 := Session{}.Create("1/25/2023", "1:00AM", "3H", "Anderson", "Invasion")
	session4 := Session{}.Create("2/25/2024", "2:00AM", "4H", "Anderson", "Holy-Grail")

	session.MakeRecord(db)
	session2.MakeRecord(db)
	session3.MakeRecord(db)
	session4.MakeRecord(db)

	session_records, err := GetAllSessionRecords(db)

	if err != nil {
		fmt.Println("Error Get All Session Records: ", err)
		return
	}

	for i := 0; i < 4; i++ {
		session_records[i].Display()
		fmt.Println()
	}
	session.DeleteRecord(db)
	session2.DeleteRecord(db)
	session3.DeleteRecord(db)
	session4.DeleteRecord(db)
	return
}

func TestDB(t *testing.T) {

	db, err := sql.Open("postgres", "user=postgres dbname=sp_calendar password=rxpt221!@# sslmode=disable")
	if err != nil {
		fmt.Println("Error after opening db: ", err)
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

	err = testGetAllSessions(db)
	if err != nil {
		fmt.Println("Error test get all session in main: ", err)
	}

}
