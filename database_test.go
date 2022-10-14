package main

import (
	"fmt"
	"testing"
)

func testSpUser() (err error) {
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
func testSpManager() (err error) {
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

func testSession() (err error) {
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

func TestDB(t *testing.T) {

	err := testSpUser()
	if err != nil {
		fmt.Println("Error test db user in main: ", err)
	}

	err = testSpManager()
	if err != nil {
		fmt.Println("Error test db user in main: ", err)
	}

	err = testSession()
	if err != nil {
		fmt.Println("Error test db session in main: ", err)
	}

}
