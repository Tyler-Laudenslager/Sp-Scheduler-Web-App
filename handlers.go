package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	httpRedirectResponse = http.StatusFound
)

func formatTitle(title string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(title)), "")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	}
	t, _ := template.ParseFiles("templates/html-boilerplate.html", "templates/login-content.html")
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
		spuser_records, err := GetAllSpUserRecords(db)
		spmanager.AssignedPatients = spuser_records
		if err != nil {
			fmt.Println("Error Get All User records in dashboard: ", err)
		}
		dashboard_content.Role = "Manager"
		dashboard_content.User = spmanager
		isSpManager = true
	} else {
		spuser.SessionsPool = make([]*SessionInfo, 0)
		sessions_viewed := append(spuser.SessionsAvailable, spuser.SessionsUnavailable...)
		sessions_viewed = append(sessions_viewed, spuser.SessionsAssigned...)
		for _, session_info := range session_records {
			viewed_session := false
			for _, session_viewed_info := range sessions_viewed {
				if session_info.Title == session_viewed_info.Title {
					viewed_session = true
					break
				}
			}
			if !viewed_session {
				spuser.SessionsPool = append(spuser.SessionsPool, session_info)
			} else {
				continue
			}
		}

		err = spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating record")
		}
		dashboard_content.Role = "Standardized Patient"
		dashboard_content.User = spuser
	}
	funcMap := template.FuncMap{"formatTitle": formatTitle}
	t = template.New("templates/html-boilerplate.html").Funcs(funcMap)
	if !isSpManager {
		t, _ = t.ParseFiles("templates/html-boilerplate.html", "templates/dashboard-content.html", "templates/session-content-available.html", "templates/user-settings.html")
	} else {
		t, _ = t.ParseFiles("templates/html-boilerplate.html", "templates/dashboard-content-manager.html", "templates/session-content-manager.html")
	}

	t.ExecuteTemplate(w, "html-boilerplate", dashboard_content)
}

func createsession(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	starttime := r.PostFormValue("starttime")
	endtime := r.PostFormValue("endtime")
	location := r.PostFormValue("location")
	description := r.PostFormValue("description")
	patientsneeded, err := strconv.Atoi(r.PostFormValue("patientsneeded"))
	if err != nil {
		fmt.Println("Error converting patients needed to integer")
	}
	newSession := Session{}.Create(title, date, starttime, endtime, location, description)
	newSession.Information.ShowSession = false
	newSession.PatientsNeeded = patientsneeded
	err = newSession.MakeRecord(db)
	if err != nil {
		fmt.Println("Error in Create Session Make Record : ", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func updatesession(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	starttime := r.PostFormValue("starttime")
	endtime := r.PostFormValue("endtime")
	location := r.PostFormValue("location")
	description := r.PostFormValue("description")
	newtitle := r.PostFormValue("newtitle")
	newdate := r.PostFormValue("newdate")
	newstarttime := r.PostFormValue("newstarttime")
	newendtime := r.PostFormValue("newendtime")
	newlocation := r.PostFormValue("newlocation")
	newdescription := r.PostFormValue("newdescription")
	newpatientsneeded, err := strconv.Atoi(r.PostFormValue("newpatientsneeded"))
	if err != nil {
		fmt.Println("Error converting patients needed to integer")
	}
	sessionInfo := SessionInfo{
		Title:       title,
		Date:        date,
		StartTime:   starttime,
		EndTime:     endtime,
		Location:    location,
		Description: description,
	}
	foundSession, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error in Get Session Record in Update Session : ", err)
	}
	foundSession.Information.Title = newtitle
	foundSession.Information.Date = newdate
	foundSession.Information.StartTime = newstarttime
	foundSession.Information.EndTime = newendtime
	foundSession.Information.Location = newlocation
	foundSession.Information.Description = newdescription
	foundSession.PatientsNeeded = newpatientsneeded

	err = foundSession.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error in Update Session Make Record : ", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func assignsp(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	starttime := r.PostFormValue("starttime")
	endtime := r.PostFormValue("endtime")
	location := r.PostFormValue("location")
	description := r.PostFormValue("description")
	sessionInfo := SessionInfo{
		Title:       title,
		Date:        date,
		StartTime:   starttime,
		EndTime:     endtime,
		Location:    location,
		Description: description,
	}
	foundSession, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error getting record in database", err)
		return
	}
	usersToRemove := make([]string, 0)
	for i := 0; i < len(foundSession.PatientsAvailable); i++ {
		patient := *foundSession.PatientsAvailable[i]
		if r.PostFormValue(patient.Username) == "true" {
			foundSession.PatientsAssigned = append(foundSession.PatientsAssigned, &patient)
			usersToRemove = append(usersToRemove, patient.Username)
		}
	}

	for _, username := range usersToRemove {
		for i, su := range foundSession.PatientsAvailable {
			if su.Username == username {
				foundSession.PatientsAvailable = append(foundSession.PatientsAvailable[:i], foundSession.PatientsAvailable[i+1:]...)
			}
		}
	}

	for _, spuser := range foundSession.PatientsAssigned {
		username := spuser.Username
		spuserRecord, err := GetSpUserRecord(username, db)
		if err != nil {
			fmt.Println("Error Getting Record: ", err)
			return
		}
		spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned, foundSession.Information)
		if len(spuserRecord.SessionsAvailable) > 0 {
			for i := 0; i < len(spuserRecord.SessionsAvailable); i++ {
				if spuserRecord.SessionsAvailable[i].Title == foundSession.Information.Title {
					spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable[:i], spuserRecord.SessionsAvailable[i+1:]...)
				}
			}
		}
		if len(spuserRecord.SessionsPool) > 0 {
			for i := 0; i < len(spuserRecord.SessionsPool); i++ {
				if spuserRecord.SessionsPool[i].Title == foundSession.Information.Title {
					spuserRecord.SessionsPool = append(spuserRecord.SessionsPool[:i], spuserRecord.SessionsPool[i+1:]...)
				}
			}
		}
		if len(spuserRecord.SessionsUnavailable) > 0 {
			for i := 0; i < len(spuserRecord.SessionsUnavailable); i++ {
				if spuserRecord.SessionsUnavailable[i].Title == foundSession.Information.Title {
					spuserRecord.SessionsUnavailable = append(spuserRecord.SessionsUnavailable[:i], spuserRecord.SessionsUnavailable[i+1:]...)
				}
			}
		}
		spuserRecord.TotalSessionsAssigned = spuserRecord.TotalSessionsAssigned + 1
		fmt.Println(spuserRecord.TotalSessionsAssigned)
		err = spuserRecord.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error Updating Record: ", err)
			return
		}
	}

	err = foundSession.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error updating record in assign sp", err)
		return
	}

	http.Redirect(w, r, "/dashboard", httpRedirectResponse)

}

func deletesession(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	starttime := r.PostFormValue("starttime")
	endtime := r.PostFormValue("endtime")
	location := r.PostFormValue("location")
	description := r.PostFormValue("description")
	sessionInfo := SessionInfo{
		Title:       title,
		Date:        date,
		StartTime:   starttime,
		EndTime:     endtime,
		Location:    location,
		Description: description,
	}
	foundSession, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error getting record in database", err)
	}
	err = foundSession.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error deleting record in database", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}
func signupavailable(w http.ResponseWriter, r *http.Request) {
	duplicate := false
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)
	if err != nil {
		fmt.Println("Error: GetSpUserRecord in signupavailable", err)
	}
	sessionInfo := SessionInfo{
		Title:       r.PostFormValue("Title"),
		Date:        r.PostFormValue("Date"),
		StartTime:   r.PostFormValue("StartTime"),
		EndTime:     r.PostFormValue("EndTime"),
		Location:    r.PostFormValue("Location"),
		Description: r.PostFormValue("Description"),
	}

	availableSessionRecord, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	for _, su := range availableSessionRecord.PatientsAvailable {
		if su.Name.First == spuser.Name.First && su.Name.Last == spuser.Name.Last {
			duplicate = true
		}
	}
	if !duplicate {
		availableSessionRecord.PatientsAvailable = append(availableSessionRecord.PatientsAvailable, &spuser)
		err = availableSessionRecord.UpdateRecord(db)
	}
	if err != nil {
		fmt.Println("Error updating session record", err)
	}
	availableSessionRecord, err = GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	duplicate = false
	if spuser.SessionsAvailable != nil {
		for i := 0; i < len(spuser.SessionsAvailable); i++ {
			if *availableSessionRecord.Information == *spuser.SessionsAvailable[i] {
				duplicate = true
			}
		}
	}
	if !duplicate {
		for i := 0; i < len(spuser.SessionsPool); i++ {
			if spuser.SessionsPool[i].Title == sessionInfo.Title {
				spuser.SessionsPool = append(spuser.SessionsPool[:i], spuser.SessionsPool[i+1:]...)
			}
		}
		for i := 0; i < len(spuser.SessionsUnavailable); i++ {
			if spuser.SessionsUnavailable[i].Title == sessionInfo.Title {
				spuser.SessionsUnavailable = append(spuser.SessionsUnavailable[:i], spuser.SessionsUnavailable[i+1:]...)
			}
		}
		spuser.SessionsAvailable = append(spuser.SessionsAvailable, availableSessionRecord.Information)
		spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating record in signupavailable: ", err)
		}
		spuser, err = GetSpUserRecord(session.Values["username"].(string), db)
		if err != nil {
			fmt.Println("Error: GetSpUserRecord in signupavailable", err)
		}

	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func signupnotavailable(w http.ResponseWriter, r *http.Request) {
	duplicate := false
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)
	if err != nil {
		fmt.Println("Error: GetSpUserRecord in signupavailable", err)
	}
	sessionInfo := SessionInfo{
		Title:       r.PostFormValue("Title"),
		Date:        r.PostFormValue("Date"),
		StartTime:   r.PostFormValue("StartTime"),
		EndTime:     r.PostFormValue("EndTime"),
		Location:    r.PostFormValue("Location"),
		Description: r.PostFormValue("Description"),
	}

	notAvailableSessionRecord, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	if spuser.SessionsUnavailable != nil {
		for i := 0; i < len(spuser.SessionsUnavailable); i++ {
			if *notAvailableSessionRecord.Information == *spuser.SessionsUnavailable[i] {
				duplicate = true
			}
		}
	}
	if !duplicate {
		for i := 0; i < len(spuser.SessionsPool); i++ {
			if spuser.SessionsPool[i].Title == sessionInfo.Title {
				spuser.SessionsPool = append(spuser.SessionsPool[:i], spuser.SessionsPool[i+1:]...)
			}
		}
		for i := 0; i < len(spuser.SessionsAvailable); i++ {
			if spuser.SessionsAvailable[i].Title == sessionInfo.Title {
				spuser.SessionsAvailable = append(spuser.SessionsAvailable[:i], spuser.SessionsAvailable[i+1:]...)
			}
		}
		spuser.SessionsUnavailable = append(spuser.SessionsUnavailable, notAvailableSessionRecord.Information)
		spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating record in signupavailable: ", err)
		}
		spuser, err = GetSpUserRecord(session.Values["username"].(string), db)
		if err != nil {
			fmt.Println("Error: GetSpUserRecord in signupavailable", err)
		}
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func changeemail(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)
	if err != nil {
		fmt.Println("Error: GetSpUserRecord in signupavailable", err)
	}
	newemail := r.PostFormValue("newemail")
	spuser.Email = newemail
	err = spuser.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error updating user record in change email handler : ", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func changepassword(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)
	if err != nil {
		fmt.Println("Error: GetSpUserRecord in signupavailable", err)
	}
	newPassword := r.PostFormValue("newpassword")
	newPasswordConfirmed := r.PostFormValue("newpasswordconfirmed")
	if newPassword == newPasswordConfirmed {
		hashedPassword, err := HashPassword(newPassword)
		if err != nil {
			fmt.Println("Error Hashing Password")
			http.Redirect(w, r, "/dashboard", httpRedirectResponse)
		}
		spuser.Password = hashedPassword
	}
	err = spuser.UpdateRecord(db)
	if err != nil {
		fmt.Println("error updating user record in change password", err)
	}
	fmt.Println("New Password is :", newPassword)
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)

}

func toggleshowsession(w http.ResponseWriter, r *http.Request) {
	sessionInfo := SessionInfo{
		Title:       r.PostFormValue("title"),
		Date:        r.PostFormValue("date"),
		StartTime:   r.PostFormValue("starttime"),
		EndTime:     r.PostFormValue("endtime"),
		Location:    r.PostFormValue("location"),
		Description: r.PostFormValue("description"),
	}

	availableSessionRecord, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	if availableSessionRecord.Information.ShowSession {
		availableSessionRecord.Information.ShowSession = false
	} else {
		availableSessionRecord.Information.ShowSession = true
	}
	allSpUsers, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error Getting all SP User records: ", err)
	}
	for _, su := range allSpUsers {
		allSessions := append(su.SessionsAssigned, su.SessionsAvailable...)
		allSessions = append(allSessions, su.SessionsUnavailable...)
		allSessions = append(allSessions, su.SessionsPool...)
		for _, si := range allSessions {
			if availableSessionRecord.Information.Title == si.Title {
				si.ShowSession = availableSessionRecord.Information.ShowSession
			}
		}
		su.UpdateRecord(db)
	}
	err = availableSessionRecord.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error Updating Session Record in Toggle Show Session : ", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	username := r.PostFormValue("userid")
	password := r.PostFormValue("password")
	spuser, err := GetSpUserRecord(username, db)
	if err != nil {
		spmanager, err := GetSpManagerRecord(username, db)
		if err != nil {
			http.Redirect(w, r, "/login", httpRedirectResponse)
		}
		if !CheckPasswordHash(password, spmanager.Password) {
			http.Redirect(w, r, "/login", httpRedirectResponse)
		}
		session.Values["authenticated"] = true
		session.Values["username"] = spmanager.Username
	} else {
		if !CheckPasswordHash(password, spuser.Password) {
			http.Redirect(w, r, "/login", httpRedirectResponse)
		}
		session.Values["authenticated"] = true
		session.Values["username"] = spuser.Username
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
	}
	session.Save(r, w)
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func sendjson(w http.ResponseWriter, r *http.Request) {
	bob_marcs := SpUser{}.Create(*Name{}.Create("Bob Marcs"), "bmarcs", SP, "bob@marcs.com")
	susan_miller := SpUser{}.Create(*Name{}.Create("Susan Miller"), "smiller", SP, "susan@miller.com")

	andy_thomas := SpManager{}.Create(*Name{}.Create("Andy Thomas"), Manager, "andy@thomas.com")

	session1 := Session{}.Create("Anderson Clinical Nurse Session", "11/15/2022", "11:00AM", "12:00PM", "Anderson", "Check-Up")

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

func createSPRecord(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	split_name := strings.Split(strings.ToLower(name), " ")
	first_initial := string(split_name[0][0])
	last_name := string(split_name[1])
	username := first_initial + last_name
	spuser := SpUser{}.Create(*Name{}.Create(name), username, SP, email)
	hashedPassword, err := HashPassword(password)
	spuser.Password = hashedPassword
	if err != nil {
		fmt.Println("Password Hash Gone Wrong in Create SP Record: ", err)
	}
	err = spuser.MakeRecord(db)
	if err != nil {
		fmt.Println("Error creating record in database in CreateSPRecord: ", err)
	}
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}

func deleteSPRecord(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	spuser, err := GetSpUserRecord(username, db)
	if err != nil {
		fmt.Println("Error Getting SP User Record: ", err)
		return
	}
	err = spuser.DeleteRecord(db)
	if err != nil {
		fmt.Println("Error deleting SP user record: ", err)
		return
	}
	fmt.Println("Deleted SP User Record: ", username)
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}
