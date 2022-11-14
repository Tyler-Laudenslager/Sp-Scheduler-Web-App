package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type Role int

const (
	SP Role = iota + 1
	Manager
	SuperUser
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

func (n Name) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Name) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &n)
}

type Instructor struct {
	Name  Name   `json:"Name"`
	Title string `json:"Title"`
}

func (i Instructor) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Instructor) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &i)
}

func (i Instructor) Create(full_name string, title string) *Instructor {
	return &Instructor{
		Name:  *Name{}.Create(full_name),
		Title: title,
	}
}

type SpUser struct {
	Id                    uint           `json:"Id"`
	Name                  Name           `json:"Name"`
	Username              string         `json:"Username"`
	Role                  Role           `json:"Role"`
	TotalSessionsAssigned uint32         `json:"TotalSessionsAssigned"`
	SessionsPool          []*SessionInfo `json:"SessionsPool"`
	SessionsAvailable     []*SessionInfo `json:"SessionsAvailable"`
	SessionsUnavailable   []*SessionInfo `json:"SessionsUnavailable"`
	SessionsAssigned      []*SessionInfo `json:"SessionsAssigned"`
	Password              string         `json:"Password"`
	Email                 string         `json:"Email"`
}

func (sp SpUser) Value() (driver.Value, error) {
	return json.Marshal(sp)
}

func (sp *SpUser) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &sp)
}

func (spUser SpUser) Create(name Name, username string, role Role, email string) *SpUser {
	return &SpUser{
		Name:                  name,
		Username:              username,
		Role:                  role,
		TotalSessionsAssigned: 0,
		SessionsPool:          make([]*SessionInfo, 0),
		SessionsAvailable:     make([]*SessionInfo, 0),
		SessionsUnavailable:   make([]*SessionInfo, 0),
		SessionsAssigned:      make([]*SessionInfo, 0),
		Password:              "",
		Email:                 email,
	}
}

type SpManager struct {
	Id                uint       `json:"Id"`
	Name              Name       `json:"Name"`
	Username          string     `json:"Username"`
	Role              Role       `json:"Role"`
	AssignedPatients  []*SpUser  `json:"AssignedPatients"`
	SessionsManaged   []*Session `json:"SessionsManaged"`
	SessionsUnmanaged []*Session `json:"SessionsUnmanaged"`
	Password          string     `json:"Password"`
	Email             string     `json:"Email"`
}

func (sp SpManager) Value() (driver.Value, error) {
	return json.Marshal(sp)
}

func (sp *SpManager) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &sp)
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
	Id                  uint          `json:"Id"`
	Information         *SessionInfo  `json:"Information"`
	Instructors         []*Instructor `json:"Instructors"`
	PatientsNeeded      int           `json:"PatientsNeeded"`
	PatientsAssigned    []*SpUser     `json:"PatientsAssigned"`
	PatientsAvailable   []*SpUser     `json:"PatientsAvailable"`
	PatientsUnavailable []*SpUser     `json:"PatientsUnavailable"`
	PatientsNoResponse  []*SpUser     `json:"PatientsNoResponse"`
}

func (s Session) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Session) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}

func (s Session) Create(title string, date string, starttime string, endtime string, location string, description string) *Session {
	return &Session{
		Information: &SessionInfo{
			Title:       title,
			Date:        date,
			StartTime:   starttime,
			EndTime:     endtime,
			Location:    location,
			Description: description,
			ShowSession: true},
		Instructors:         []*Instructor{},
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

type SessionInfo struct {
	Title       string `json:"Title"`
	Date        string `json:"Date"`
	StartTime   string `json:"StartTime"`
	EndTime     string `json:"EndTime"`
	Location    string `json:"Location"`
	Description string `json:"Description"`
	ShowSession bool   `json:"ShowSession"`
}

func (si SessionInfo) Value() (driver.Value, error) {
	return json.Marshal(si)
}

func (si *SessionInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &si)
}
