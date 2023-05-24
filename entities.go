package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
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

type Comment struct {
	Author      string `json:"Creator"`
	DateCreated string `json:"DateCreated"`
	TimeCreated string `json:"TimeCreated"`
	Content     string `json:"Content"`
}

func (c Comment) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Comment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
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
	SessionsSorted        []*SessionInfo `json:"SessionsSorted"`
	SessionsAvailable     []*SessionInfo `json:"SessionsAvailable"`
	SessionsUnavailable   []*SessionInfo `json:"SessionsUnavailable"`
	SessionsSelected      []*SessionInfo `json:"SessionsSelected"`
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
		SessionsSorted:        make([]*SessionInfo, 0),
		SessionsAvailable:     make([]*SessionInfo, 0),
		SessionsUnavailable:   make([]*SessionInfo, 0),
		SessionsSelected:      make([]*SessionInfo, 0),
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
	PatientsSelected    []*SpUser     `json:"PatientsSelected"`
	PatientsAvailable   []*SpUser     `json:"PatientsAvailable"`
	PatientsUnavailable []*SpUser     `json:"PatientsUnavailable"`
	PatientsNoResponse  []*SpUser     `json:"PatientsNoResponse"`
	Comments            []*Comment    `json:"Comments"`
}

type SessionContainer []*Session

func (a SessionContainer) Len() int { return len(a) }
func (a SessionContainer) Less(i, j int) bool {
	iDate := a[i].Information.Date
	jDate := a[j].Information.Date

	iStartTime := strings.ReplaceAll(a[i].Information.StartTime, " ", "")

	// AM or PM
	iStartEnding := strings.ToUpper(iStartTime[len(iStartTime)-2:])

	jStartTime := strings.ReplaceAll(a[j].Information.StartTime, " ", "")
	// AM or PM
	jStartEnding := strings.ToUpper(jStartTime[len(jStartTime)-2:])

	iTimeOfDay := iStartTime[:len(iStartTime)-2]
	jTimeOfDay := jStartTime[:len(jStartTime)-2]

	iParsed, _ := time.Parse("01/02/2006", iDate)
	jParsed, _ := time.Parse("01/02/2006", jDate)

	if iDate == jDate {

		if iStartEnding == "PM" && jStartEnding == "AM" {
			return iParsed.Before(jParsed)
		} else if iStartEnding == "AM" && jStartEnding == "PM" {
			return !iParsed.Before(jParsed)
		} else if iStartEnding == jStartEnding {
			iHour, _ := strconv.Atoi(iTimeOfDay[:strings.Index(iTimeOfDay, ":")])
			jHour, _ := strconv.Atoi(jTimeOfDay[:strings.Index(jTimeOfDay, ":")])
			iMinutes, _ := strconv.Atoi(iTimeOfDay[strings.Index(iTimeOfDay, ":")+1:])
			jMinutes, _ := strconv.Atoi(jTimeOfDay[strings.Index(jTimeOfDay, ":")+1:])

			iMin := float64(iMinutes) / 60.0
			jMin := float64(jMinutes) / 60.0

			iH := float64(iHour) + iMin
			jH := float64(jHour) + jMin
			if iHour != 12.0 {
				iH += 12.0
			}
			if jHour != 12.0 {
				jH += 12.0
			}

			if iH < jH {
				return !iParsed.Before(jParsed)
			} else {
				return iParsed.Before(jParsed)
			}

		}
	}
	return iParsed.Before(jParsed)

}
func (a SessionContainer) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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
			Title:             title,
			Date:              date,
			StartTime:         starttime,
			EndTime:           endtime,
			Location:          location,
			Description:       description,
			Status:            "noresponse",
			CreatedDate:       "",
			ExpiredDate:       "",
			CheckMarkAssigned: false,
			CheckXCanceled:    false,
			ShowSession:       false,
			Comments:          map[string][]*Comment{}},
		Instructors:         []*Instructor{},
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsSelected:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

type SessionInfo struct {
	Title             string                `json:"Title"`
	Date              string                `json:"Date"`
	StartTime         string                `json:"StartTime"`
	EndTime           string                `json:"EndTime"`
	Location          string                `json:"Location"`
	Description       string                `json:"Description"`
	Status            string                `json:"Status"`
	CreatedDate       string                `json:"CreatedDate"`
	ExpiredDate       string                `json:"ExpiredDate"`
	ShowSession       bool                  `json:"ShowSession"`
	CheckMarkAssigned bool                  `json:"CheckMarkAssigned"`
	CheckXCanceled    bool                  `json:"CheckXCanceled"`
	Comments          map[string][]*Comment `json:"Comments"`
}
type SessionInfoContainer []*SessionInfo

func (a SessionInfoContainer) Len() int { return len(a) }
func (a SessionInfoContainer) Less(i, j int) bool {
	iDate := a[i].Date
	jDate := a[j].Date

	iStartTime := strings.ReplaceAll(a[i].StartTime, " ", "")

	// AM or PM
	iStartEnding := strings.ToUpper(iStartTime[len(iStartTime)-2:])

	jStartTime := strings.ReplaceAll(a[j].StartTime, " ", "")
	// AM or PM
	jStartEnding := strings.ToUpper(jStartTime[len(jStartTime)-2:])

	iTimeOfDay := iStartTime[:len(iStartTime)-2]
	jTimeOfDay := jStartTime[:len(jStartTime)-2]

	iParsed, _ := time.Parse("01/02/2006", iDate)
	jParsed, _ := time.Parse("01/02/2006", jDate)

	if iDate == jDate {

		if iStartEnding == "PM" && jStartEnding == "AM" {
			return iParsed.Before(jParsed)
		} else if iStartEnding == "AM" && jStartEnding == "PM" {
			return !iParsed.Before(jParsed)
		} else if iStartEnding == jStartEnding {
			iHour, _ := strconv.Atoi(iTimeOfDay[:strings.Index(iTimeOfDay, ":")])
			jHour, _ := strconv.Atoi(jTimeOfDay[:strings.Index(jTimeOfDay, ":")])
			iMinutes, _ := strconv.Atoi(iTimeOfDay[strings.Index(iTimeOfDay, ":")+1:])
			jMinutes, _ := strconv.Atoi(jTimeOfDay[strings.Index(jTimeOfDay, ":")+1:])

			iMin := float64(iMinutes) / 60.0
			jMin := float64(jMinutes) / 60.0

			iH := float64(iHour) + iMin
			jH := float64(jHour) + jMin
			if iH != 12.0 {
				iH += 12.0
			}
			if jH != 12.0 {
				jH += 12.0
			}
			if iH < jH {
				return !iParsed.Before(jParsed)
			} else {
				return iParsed.Before(jParsed)
			}

		}
	}
	return iParsed.Before(jParsed)

}
func (a SessionInfoContainer) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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
