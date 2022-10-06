package main

import (
	"reflect"
	"testing"
)

func TestSpUserCreation(t *testing.T) {
	name := Name{
		First: "Bob",
		Last:  "Miller",
	}
	SpUserObj := SpUser{}.Create(name, SP, Male, "generic@email.com")
	if SpUserObj.Name.First != "Bob" {
		t.Error("Expecting 'Bob' : Received", SpUserObj.Name.First)
	}
	if SpUserObj.Name.Last != "Miller" {
		t.Error("Expecting 'Miller' : Received", SpUserObj.Name.Last)
	}
	if SpUserObj.Role != SP {
		t.Error("Expecting 'SP' : Received", SpUserObj.Role)
	}
	if SpUserObj.Sex != Male {
		t.Error("Expecting 'Male' : Received", SpUserObj.Sex)
	}
	if SpUserObj.Password != "" {
		t.Error("Expecting 'letmein' : Received", SpUserObj.Password)
	}
	if SpUserObj.Email != "generic@email.com" {
		t.Error("Expecting 'generic@email.com' : Received", SpUserObj.Email)
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
}
