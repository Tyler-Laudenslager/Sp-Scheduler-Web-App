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
	First string `json:"First"`
	Last  string `json:"Last"`
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
	Name  Name   `json:"Name"`
	Title string `json:"Title"`
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
	Name              Name       `json:"Name"`
	Role              Role       `json:"Role"`
	AssignedPatients  []*SpUser  `json:"AssignedPatients"`
	SessionsManaged   []*Session `json:"SessionsManaged"`
	SessionsUnmanaged []*Session `json:"SessionsUnmanaged"`
	Password          string     `json:"Password"`
	Email             string     `json:"Email"`
}

func (spManager SpManager) Create(name Name, role Role, email string) *SpManager {
	return &SpManager{
		Name:              name,
		Role:              role,
		AssignedPatients:  []*SpUser{},
		SessionsManaged:   []*Session{},
		SessionsUnmanaged: []*Session{},
		Password:          "",
		Email:             email,
	}
}

type Session struct {
	Date                string        `json:"Date"`
	Time                string        `json:"Time"`
	Duration            string        `json:"Duration"`
	Location            string        `json:"Location"`
	Description         string        `json:"Description"`
	Information         *SessionInfo  `json:"Information"`
	Instructors         []*Instructor `json:"Instructors"`
	PatientsNeeded      int           `json:"PatientsNeeded"`
	PatientsAssigned    []*SpUser     `json:"PatientsAssigned"`
	PatientsAvailable   []*SpUser     `json:"PatientsAvailable"`
	PatientsUnavailable []*SpUser     `json:"PatientsUnavailable"`
	PatientsNoResponse  []*SpUser     `json:"PatientsNoResponse"`
}

func (s Session) Create(date string, time string, duration string, location string, description string) *Session {
	return &Session{
		Information: &SessionInfo{Date: date,
			Time:        time,
			Duration:    duration,
			Location:    location,
			Description: description},
		Instructors:         []*Instructor{},
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

type SessionInfo struct {
	Date        string `json:"Date"`
	Time        string `json:"Time"`
	Duration    string `json:"Duration"`
	Location    string `json:"Location"`
	Description string `json:"Description"`
}

func (s Session) Info() *SessionInfo {
	if s.Information.Date == "" &&
		s.Information.Time == "" &&
		s.Information.Duration == "" &&
		s.Information.Location == "" {

		s.Information.Date = s.Date
		s.Information.Time = s.Time
		s.Information.Duration = s.Duration
		s.Information.Location = s.Location
		s.Information.Description = s.Description
		return s.Information
	} else {
		return s.Information
	}
}
