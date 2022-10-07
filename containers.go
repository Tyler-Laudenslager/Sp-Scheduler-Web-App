package main

type HospitalCalendar struct {
	Users    SpUsers    `json:"SP Users"`
	Managers SpManagers `json:"SP Managers"`
	Sessions SpSessions `json:"SP Sessions"`
}
type SpUsers []*SpUser
type SpManagers []*SpManager
type SpSessions []*Session
