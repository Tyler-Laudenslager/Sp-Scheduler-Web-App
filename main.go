package main

import (
	"encoding/json"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	bob_marcs := SpUser{}.Create(*Name{}.Create("Bob Marcs"), SP, Male, "bob@marcs.com")
	susan_miller := SpUser{}.Create(*Name{}.Create("Susan Miller"), SP, Female, "susan@miller.com")

	andy_thomas := SpManager{}.Create(*Name{}.Create("Andy Thomas"), Manager, "andy@thomas.com")

	session1 := Session{}.Create("11/15/2022", "11:00AM", "1H", "Anderson")

	andy_thomas.AssignedPatients = append(andy_thomas.AssignedPatients, bob_marcs, susan_miller)
	andy_thomas.SessionsManaged = append(andy_thomas.SessionsManaged, session1)

	bob_marcs.SessionsAvailable = append(bob_marcs.SessionsAvailable, session1.Info())
	bob_marcs.SessionsAssigned = append(bob_marcs.SessionsAssigned, session1.Info())

	susan_miller.SessionsAvailable = append(susan_miller.SessionsAvailable, session1.Info())
	susan_miller.SessionsAssigned = append(susan_miller.SessionsAssigned, session1.Info())

	session1.PatientsAvailable = append(session1.PatientsAvailable, bob_marcs)
	session1.PatientsAvailable = append(session1.PatientsAvailable, susan_miller)
	session1.PatientsAssigned = append(session1.PatientsAssigned, bob_marcs)
	session1.PatientsAssigned = append(session1.PatientsAssigned, susan_miller)

	SpUsersBox := append(make(SpUsers, 0, 2), bob_marcs, susan_miller)
	SpManagersBox := append(make(SpManagers, 0, 1), andy_thomas)
	SessionsBox := append(make(SpSessions, 0, 1), session1)
	HospitalCalendar := HospitalCalendar{
		Users:    SpUsersBox,
		Managers: SpManagersBox,
		Sessions: SessionsBox,
	}

	output, err := json.MarshalIndent(HospitalCalendar, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:6600",
	}
	http.HandleFunc("/", index)
	server.ListenAndServe()
}
