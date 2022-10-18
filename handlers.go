package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

const (
	httpRedirectResponse = http.StatusFound
)

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	}
	t, _ := template.ParseFiles("html-boilerplate.html", "login-content.html")
	t.ExecuteTemplate(w, "html-boilerplate", "")
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login", httpRedirectResponse)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	var t *template.Template
	var spmanager SpManager
	isSpManager := false
	dashboard_content := DashboardContent{
		Date: time.Now().Format("01-02-2006"),
	}
	session_records, err := GetAllSessionInfoRecords(db)
	if err != nil {
		fmt.Println("Error Get All Session Records: ", err)
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)

	if err != nil {
		spmanager, err = GetSpManagerRecord(session.Values["username"].(string), db)
		if err != nil {
			fmt.Println("Error Get User record in dashboard: ", err)
			return
		}
		session_records_manager, err := GetAllSessionRecords(db)
		if err != nil {
			fmt.Println("Error Get All Session records in dashboard: ", err)
			return
		}
		spmanager.SessionsUnmanaged = session_records_manager
		dashboard_content.User = spmanager
		isSpManager = true
	} else {
		spuser.SessionsAvailable = session_records
		dashboard_content.User = spuser
	}
	if !isSpManager {
		t, _ = template.ParseFiles("html-boilerplate.html", "dashboard-content.html", "session-content.html")
	} else {
		t, _ = template.ParseFiles("html-boilerplate.html", "dashboard-content-manager.html", "session-content-manager.html")
	}

	t.ExecuteTemplate(w, "html-boilerplate", dashboard_content)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	var spmanager SpManager
	isSpManager := false
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	username := r.PostFormValue("userid")
	password := r.PostFormValue("password")
	spuser, err := GetSpUserRecord(username, db)
	if err != nil {
		spmanager, err = GetSpManagerRecord(username, db)
		if err != nil {
			http.Redirect(w, r, "/login", httpRedirectResponse)
		}
		isSpManager = true
	}
	if !isSpManager && !CheckPasswordHash(password, spuser.Password) {
		http.Redirect(w, r, "/login", httpRedirectResponse)
	} else {
		if isSpManager && CheckPasswordHash(password, spmanager.Password) {
			session.Values["authenticated"] = true
			session.Values["username"] = spmanager.Username
		} else {
			session.Values["authenticated"] = true
			session.Values["username"] = spuser.Username
		}
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
	}
	session.Save(r, w)
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func sendjson(w http.ResponseWriter, r *http.Request) {
	bob_marcs := SpUser{}.Create(*Name{}.Create("Bob Marcs"), SP, Male, "bob@marcs.com")
	susan_miller := SpUser{}.Create(*Name{}.Create("Susan Miller"), SP, Female, "susan@miller.com")

	andy_thomas := SpManager{}.Create(*Name{}.Create("Andy Thomas"), Manager, "andy@thomas.com")

	session1 := Session{}.Create("11/15/2022", "11:00AM", "1H", "Anderson", "Check-Up")

	andy_thomas.AssignedPatients = append(andy_thomas.AssignedPatients, bob_marcs, susan_miller)
	andy_thomas.SessionsManaged = append(andy_thomas.SessionsManaged, session1)

	bob_marcs.SessionsAvailable = append(bob_marcs.SessionsAvailable, session1.Information)
	bob_marcs.SessionsAssigned = append(bob_marcs.SessionsAssigned, session1.Information)

	susan_miller.SessionsAvailable = append(susan_miller.SessionsAvailable, session1.Information)
	susan_miller.SessionsAssigned = append(susan_miller.SessionsAssigned, session1.Information)

	session1.PatientsAvailable = append(session1.PatientsAvailable, bob_marcs)
	session1.PatientsAvailable = append(session1.PatientsAvailable, susan_miller)
	session1.PatientsAssigned = append(session1.PatientsAssigned, bob_marcs)
	session1.PatientsAssigned = append(session1.PatientsAssigned, susan_miller)

	SpUsersBox := append(make(SpUsersBox, 0, 2), bob_marcs, susan_miller)
	SpManagersBox := append(make(SpManagersBox, 0, 1), andy_thomas)
	SessionsBox := append(make(SpSessionsBox, 0, 1), session1)

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
