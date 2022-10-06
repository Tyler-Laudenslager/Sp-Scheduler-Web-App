package main

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

type name struct {
	first string
	last  string
}

type SPUser struct {
	Name                name
	Role                Role
	Sex                 Sex
	SessionsAvailable   []*Session
	SessionsUnavailable []*Session
	SessionsAssigned    []*Session
	Password            string
	Email               string
}

type SPManager struct {
	Name              name
	Role              Role
	AssignedPatients  []*SPUser
	SessionsManaged   []*SPUser
	SessionsUnmanaged []*SPUser
	Password          string
	Email             string
}

type Session struct {
	Date                string
	Time                string
	Duration            string
	Location            string
	Description         string
	Instructors         string
	PatientsNeeded      int
	PatientsAssigned    []*SPUser
	PatientsAvailable   []*SPUser
	PatientsUnavailable []*SPUser
	PatientsNoResponse  []*SPUser
}

type Admin struct {
}
