package main

import (
	"reflect"
	"testing"
)

func TestNameCreation(t *testing.T) {
	name := Name{}.Create("Alex Conway")
	if name.First != "Alex" {
		t.Error("Expecting 'Alex' : Received", name.First)
	}
	if name.Last != "Conway" {
		t.Error("Expecting 'Conway' : Received", name.Last)
	}
}

func TestInstructorCreation(t *testing.T) {
	InstructorObj := Instructor{}.Create("Joe Kingas", "Director of Marketing")
	if InstructorObj.Name.First != "Joe" {
		t.Error("Expected 'Joe' : Received", InstructorObj.Name.First)
	}
	if InstructorObj.Name.Last != "Kingas" {
		t.Error("Expected 'Kingas' : Received", InstructorObj.Name.Last)
	}
	if InstructorObj.Title != "Director of Marketing" {
		t.Error("Expected 'Director of Marketing' : Received", InstructorObj.Title)
	}
}

func TestSpUserCreation(t *testing.T) {
	name := Name{}.Create("Bob Miller")
	SpUserObj := SpUser{}.Create(*name, "bmiller", SP, "bob@example.com")
	if SpUserObj.Name.First != "Bob" {
		t.Error("Expecting 'Bob' : Received", SpUserObj.Name.First)
	}
	if SpUserObj.Name.Last != "Miller" {
		t.Error("Expecting 'Miller' : Received", SpUserObj.Name.Last)
	}
	if SpUserObj.Username != "bmiller" {
		t.Error("Expecting 'bmiller' : Received", SpUserObj.Username)
	}
	if SpUserObj.Role != SP {
		t.Error("Expecting 'SP' : Received", SpUserObj.Role)
	}
	if SpUserObj.TotalSessionsAssigned != 0 {
		t.Error("Expecting '0' : Received", SpUserObj.TotalSessionsAssigned)
	}
	if SpUserObj.Password != "" {
		t.Error("Expecting 'letmein' : Received", SpUserObj.Password)
	}
	if SpUserObj.Email != "bob@example.com" {
		t.Error("Expecting 'bob@example.com' : Received", SpUserObj.Email)
	}
	if reflect.TypeOf(SpUserObj.SessionsAssigned).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", SpUserObj.SessionsAssigned)
	}
	if reflect.TypeOf(SpUserObj.SessionsAvailable).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", SpUserObj.SessionsAvailable)
	}
	if reflect.TypeOf(SpUserObj.SessionsUnavailable).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", SpUserObj.SessionsUnavailable)
	}
	if reflect.TypeOf(SpUserObj.SessionsPool).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", SpUserObj.SessionsUnavailable)
	}
}

func TestSpManagerCreation(t *testing.T) {
	name := Name{}.Create("Tony Miller")
	SpManagerObj := SpManager{}.Create(*name, Manager, "tony@example.com")
	if SpManagerObj.Name.First != "Tony" {
		t.Error("Expecting 'Tony' : Received", SpManagerObj.Name.First)
	}
	if SpManagerObj.Name.Last != "Miller" {
		t.Error("Expecting 'Miller' : Received", SpManagerObj.Name.Last)
	}
	if SpManagerObj.Role != Manager {
		t.Error("Expecting 'Manager' : Received", SpManagerObj.Role)
	}
	if reflect.TypeOf(SpManagerObj.AssignedPatients).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", reflect.TypeOf(SpManagerObj.AssignedPatients))
	}
	if reflect.TypeOf(SpManagerObj.SessionsManaged).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", reflect.TypeOf(SpManagerObj.SessionsManaged))
	}
	if reflect.TypeOf(SpManagerObj.SessionsUnmanaged).Kind() != reflect.Slice {
		t.Error("Expecting 'slice' : Received", reflect.TypeOf(SpManagerObj.SessionsUnmanaged))
	}
	if SpManagerObj.Password != "" {
		t.Error("Expecting '' : Received", SpManagerObj.Password)
	}
	if SpManagerObj.Email != "tony@example.com" {
		t.Error("Expecting 'tony@example.com' : Received", SpManagerObj.Email)
	}
}

func TestSessionCreation(t *testing.T) {
	SessionObj := Session{}.Create("Anderson Clinical Nurse Session", "11/23/2022", "11:00AM", "12:00PM", "Anderson", "Check-Up")
	if SessionObj.Information.Title != "Anderson Clinical Nurse Session" {
		t.Error("Expected 'Anderson Clinical Nurse Session' : Received", SessionObj.Information.Title)
	}
	if SessionObj.Information.Date != "11/23/2022" {
		t.Error("Expected '11/23/2022' : Received", SessionObj.Information.Date)
	}
	if SessionObj.Information.StartTime != "11:00AM" {
		t.Error("Expected '11:00AM' : Received", SessionObj.Information.StartTime)
	}
	if SessionObj.Information.EndTime != "12:00PM" {
		t.Error("Expected '12:00PM' : Received", SessionObj.Information.EndTime)
	}
	if SessionObj.Information.Location != "Anderson" {
		t.Error("Expected 'Anderson' : Received", SessionObj.Information.Location)
	}
	if SessionObj.Information.Description != "Check-Up" {
		t.Error("Expected 'Check-Up' : Received", SessionObj.Information.Description)
	}
	if SessionObj.PatientsNeeded != 0 {
		t.Error("Expected '0' : Received", SessionObj.PatientsNeeded)
	}
	if reflect.TypeOf(SessionObj.Instructors).Kind() != reflect.Slice {
		t.Error("Expected 'slice' : Received", reflect.TypeOf(SessionObj.Instructors).Kind())
	}
	if reflect.TypeOf(SessionObj.PatientsAssigned).Kind() != reflect.Slice {
		t.Error("Expected 'slice' : Received", reflect.TypeOf(SessionObj.PatientsAssigned).Kind())
	}
	if reflect.TypeOf(SessionObj.PatientsAvailable).Kind() != reflect.Slice {
		t.Error("Expected 'slice' : Received", reflect.TypeOf(SessionObj.PatientsAvailable).Kind())
	}
	if reflect.TypeOf(SessionObj.PatientsNoResponse).Kind() != reflect.Slice {
		t.Error("Expected 'slice' : Received", reflect.TypeOf(SessionObj.PatientsNoResponse).Kind())
	}
	if reflect.TypeOf(SessionObj.PatientsUnavailable).Kind() != reflect.Slice {
		t.Error("Expected 'slice' : Received", reflect.TypeOf(SessionObj.PatientsUnavailable).Kind())
	}

}

func TestSessionInfoCreation(t *testing.T) {
	SessionObj := Session{}.Create("Anderson Clinical Nurse Session", "11/23/2022", "11:00AM", "12:00PM", "Anderson", "Check-Up")
	SessionInfo := SessionObj.Information
	if SessionInfo.Title != "Anderson Clinical Nurse Session" {
		t.Error("Expected 'Anderson Clinical Nurse Session' : Received", SessionObj.Information.Title)
	}
	if SessionInfo.Date != "11/23/2022" {
		t.Error("Expected '11/23/2022' : Received", SessionInfo.Date)
	}
	if SessionInfo.StartTime != "11:00AM" {
		t.Error("Expected '11:00AM' : Received", SessionObj.Information.StartTime)
	}
	if SessionInfo.EndTime != "12:00PM" {
		t.Error("Expected '12:00PM' : Received", SessionObj.Information.EndTime)
	}
	if SessionInfo.Location != "Anderson" {
		t.Error("Expected 'Anderson' : Received", SessionInfo.Location)
	}
	if SessionInfo.Description != "Check-Up" {
		t.Error("Expected 'Check-Up' : Received", SessionInfo.Description)
	}
}
