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
	First string `json:'First'`
	Last  string `json:'Last'`
}

func (n Name) Create(full_name string) *Name {
	fullname_split := strings.Split(full_name, " ")
	first, last := fullname_split[0], fullname_split[1]
	return &Name{
		First: first,
		Last:  last,
	}
}

type Instructor struct {
	Name  Name   `json:'Name'`
	Title string `json:'Title'`
}

func (i Instructor) Create(full_name string, title string) *Instructor {
	return &Instructor{
		Name:  *Name{}.Create(full_name),
		Title: title,
	}
}

type SpUser struct {
	Name                Name           `json:"Name"`
	Role                Role           `json:"Role"`
	Sex                 Sex            `json:"Sex"`
	SessionsAvailable   []*SessionInfo `json:"SessionsAvailable"`
	SessionsUnavailable []*SessionInfo `json:"SessionsUnavailable"`
	SessionsAssigned    []*SessionInfo `json:"SessionsAssigned"`
	Password            string         `json:"Password"`
	Email               string         `json:"Email"`
}

func (spUser SpUser) Create(name Name, role Role, sex Sex, email string) *SpUser {
	return &SpUser{
		Name:                name,
		Role:                role,
		Sex:                 sex,
		SessionsAvailable:   []*SessionInfo{},
		SessionsUnavailable: []*SessionInfo{},
		SessionsAssigned:    []*SessionInfo{},
		Password:            "",
		Email:               email,
	}
}

type SpManager struct {
	Name              Name      `json:'Name'`
	Role              Role      `json:'Role'`
	AssignedPatients  []*SpUser `json:'AssignedPatients'`
	SessionsManaged   []*SpUser `json:'SessionsManaged'`
	SessionsUnmanaged []*SpUser `json:'SessionsUnmanaged'`
	Password          string    `json:'Password`
	Email             string    `json:'Email`
}

func (spManager SpManager) Create(name Name, role Role, email string) *SpManager {
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
	Date                string        `json:'Date'`
	Time                string        `json:'Time'`
	Duration            string        `json:'Duration'`
	Location            string        `json:'Location'`
	Description         string        `json:'Description'`
	Instructors         []*Instructor `json:'Instructors'`
	PatientsNeeded      int           `json:'PatientsNeeded'`
	PatientsAssigned    []*SpUser     `json:'PatientsAssigned'`
	PatientsAvailable   []*SpUser     `json:'PatientsAvailable'`
	PatientsUnavailable []*SpUser     `json:'PatientsUnavailable'`
	PatientsNoResponse  []*SpUser     `json:'PatientsNoResponse'`
}

func (s Session) Create(date string, time string, duration string, location string) *Session {
	return &Session{
		Date:                date,
		Time:                time,
		Duration:            duration,
		Location:            location,
		Description:         "",
		Instructors:         []*Instructor{},
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

type SessionInfo struct {
	Date     string
	Time     string
	Duration string
	Location string
}

func (s Session) Info() *SessionInfo {
	return &SessionInfo{
		Date:     s.Date,
		Time:     s.Time,
		Duration: s.Duration,
		Location: s.Location,
	}
}

// type Admin struct {
// }
