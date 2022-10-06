package main

import (
	"strings"
)

type Role int
type Sex int

const (
	SP Role = iota + 1
	Manager
	SuperUser
)

const (
	Male Sex = iota + 1
	Female
)

type Name struct {
	First string
	Last  string
}

func (n Name) Create(full_name string) *Name {
	fullname_split := strings.Split(full_name, " ")
	first, last := fullname_split[0], fullname_split[1]
	return &Name{
		First: first,
		Last:  last,
	}
}

type SpUser struct {
	Name                Name
	Role                Role
	Sex                 Sex
	SessionsAvailable   []*Session
	SessionsUnavailable []*Session
	SessionsAssigned    []*Session
	Password            string
	Email               string
}

func (spUser SpUser) Create(name Name,
	role Role,
	sex Sex,
	email string) *SpUser {
	return &SpUser{
		Name:                name,
		Role:                role,
		Sex:                 sex,
		SessionsAvailable:   []*Session{},
		SessionsUnavailable: []*Session{},
		SessionsAssigned:    []*Session{},
		Password:            "",
		Email:               email,
	}
}

type SpManager struct {
	Name              Name
	Role              Role
	AssignedPatients  []*SpUser
	SessionsManaged   []*SpUser
	SessionsUnmanaged []*SpUser
	Password          string
	Email             string
}

func (spManager SpManager) Create(name Name,
	role Role,
	email string) *SpManager {
	return &SpManager{
		Name:              name,
		Role:              role,
		AssignedPatients:  []*SpUser{},
		SessionsManaged:   []*SpUser{},
		SessionsUnmanaged: []*SpUser{},
		Password:          "",
		Email:             email,
	}
}

type Session struct {
	Date                string
	Time                string
	Duration            string
	Location            string
	Description         string
	Instructors         string
	PatientsNeeded      int
	PatientsAssigned    []*SpUser
	PatientsAvailable   []*SpUser
	PatientsUnavailable []*SpUser
	PatientsNoResponse  []*SpUser
}

func (s Session) Create(date string,
	time string,
	duration string,
	location string) *Session {
	return &Session{
		Date:                date,
		Time:                time,
		Duration:            duration,
		Location:            location,
		Description:         "",
		Instructors:         "",
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

// type Admin struct {
// }
