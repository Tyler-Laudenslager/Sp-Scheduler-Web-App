package main

type HospitalCalendar struct {
	Users    SpUsersBox    `json:"SP Users"`
	Managers SpManagersBox `json:"SP Managers"`
	Sessions SpSessionsBox `json:"SP Sessions"`
}
type SpUsersBox []*SpUser
type SpManagersBox []*SpManager
type SpSessionsBox []*Session