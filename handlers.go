package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

func formatDate(date string) string {
	timeT, _ := time.Parse("01/02/2006", date)
	return string(timeT.Format("Monday, January 02, 2006"))

}

func pastSession(date string) bool {
	sessionDate := date
	currentDate := time.Now()

	sessionDateParsed, _ := time.Parse("01/02/2006", sessionDate)

	return sessionDateParsed.Before(currentDate)
}

func notPastSession(date string) bool {
	sessionDate := date
	currentDate := time.Now()

	sessionDateParsed, _ := time.Parse("01/02/2006", sessionDate)

	return sessionDateParsed.After(currentDate)
}

func StatusAssigned(status string) bool {
	return status == "assigned"
}

func StatusNoResponse(status string) bool {
	return status == "noresponse"
}

func StatusUnavailable(status string) bool {
	return status == "unavailable"
}

func StatusAvailable(status string) bool {
	return status == "available"
}

func sortSessionInfoByDate(a []*SessionInfo) []*SessionInfo {
	sort.Sort(SessionInfoContainer(a[:]))
	return a
}

func sortSessionByDate(a []*Session) []*Session {
	sort.Sort(SessionContainer(a[:]))
	return a
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
		Date: time.Now().Format("Monday, January 02, 2006"),
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

		if r.PostFormValue("orderBy") != "" {
			session.Values["orderBy"] = r.PostFormValue("orderBy")
		}
		if r.PostFormValue("orderBy") == "byLocation" {
			sort.Slice(session_records_manager, func(i int, j int) bool {
				return session_records_manager[i].Information.Location < session_records_manager[j].Information.Location
			})
			dashboard_content.ByLocation = true
			dashboard_content.ByDate = false
		}

		if session.Values["orderBy"] == "byLocation" {
			sort.Slice(session_records_manager, func(i int, j int) bool {
				return session_records_manager[i].Information.Location < session_records_manager[j].Information.Location
			})
			dashboard_content.ByLocation = true
			dashboard_content.ByDate = false
		}

		if r.PostFormValue("orderBy") == "byDate" {
			sort.Slice(session_records_manager, func(i int, j int) bool {
				iDate := session_records_manager[i].Information.Date
				jDate := session_records_manager[j].Information.Date

				iParsed, _ := time.Parse("01/02/2006", iDate)
				jParsed, _ := time.Parse("01/02/2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.ByLocation = false
			dashboard_content.ByDate = true
		}

		if session.Values["orderBy"] == "byDate" {
			sort.Slice(session_records_manager, func(i int, j int) bool {
				iDate := session_records_manager[i].Information.Date
				jDate := session_records_manager[j].Information.Date

				iParsed, _ := time.Parse("01/02/2006", iDate)
				jParsed, _ := time.Parse("01/02/2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.ByLocation = false
			dashboard_content.ByDate = true
		}
		session.Save(r, w)
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
		spuser.SessionsSorted = make([]*SessionInfo, 0)
		for _, si := range spuser.SessionsAssigned {
			si.Status = "assigned"
			spuser.SessionsSorted = append(spuser.SessionsSorted, si)
		}
		for _, si := range spuser.SessionsAvailable {
			si.Status = "available"
			spuser.SessionsSorted = append(spuser.SessionsSorted, si)
		}
		for _, si := range spuser.SessionsUnavailable {
			si.Status = "unavailable"
			spuser.SessionsSorted = append(spuser.SessionsSorted, si)
		}
		for _, si := range spuser.SessionsPool {
			si.Status = "noresponse"
			spuser.SessionsSorted = append(spuser.SessionsSorted, si)
		}
		if r.PostFormValue("orderBy") != "" {
			session.Values["orderBy"] = r.PostFormValue("orderBy")
		}
		if r.PostFormValue("orderBy") == "byLocation" {
			sort.Slice(spuser.SessionsSorted, func(i int, j int) bool {
				return spuser.SessionsSorted[i].Location < spuser.SessionsSorted[j].Location
			})
			dashboard_content.ByLocation = true
			dashboard_content.ByDate = false
		}

		if session.Values["orderBy"] == "byLocation" {
			sort.Slice(spuser.SessionsSorted, func(i int, j int) bool {
				return spuser.SessionsSorted[i].Location < spuser.SessionsSorted[j].Location
			})
			dashboard_content.ByLocation = true
			dashboard_content.ByDate = false
		}

		if r.PostFormValue("orderBy") == "byDate" {
			sort.Slice(spuser.SessionsSorted, func(i int, j int) bool {
				iDate := spuser.SessionsSorted[i].Date
				jDate := spuser.SessionsSorted[j].Date

				iParsed, _ := time.Parse("01/02/2006", iDate)
				jParsed, _ := time.Parse("01/02/2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.ByLocation = false
			dashboard_content.ByDate = true
		}

		if session.Values["orderBy"] == "byDate" {
			sort.Slice(spuser.SessionsSorted, func(i int, j int) bool {
				iDate := spuser.SessionsSorted[i].Date
				jDate := spuser.SessionsSorted[j].Date

				iParsed, _ := time.Parse("01/02/2006", iDate)
				jParsed, _ := time.Parse("01/02/2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.ByLocation = false
			dashboard_content.ByDate = true
		}
		session.Save(r, w)
		err = spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating record")
		}
		dashboard_content.Role = "Standardized Patient"
		dashboard_content.User = spuser
	}
	funcMap := template.FuncMap{"formatTitle": formatTitle, "formatDate": formatDate, "sortSessionInfoByDate": sortSessionInfoByDate, "sortSessionByDate": sortSessionByDate, "StatusAssigned": StatusAssigned, "StatusNoResponse": StatusNoResponse, "StatusAvailable": StatusAvailable, "StatusUnavailable": StatusUnavailable, "pastSession": pastSession, "notPastSession": notPastSession}
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
	usersToRemoveAvailable := make([]string, 0)
	usersToRemoveAssigned := make([]string, 0)
	usersToRemoveSelected := make([]string, 0)

	for i := 0; i < len(foundSession.PatientsAvailable); i++ {
		patient := *foundSession.PatientsAvailable[i]
		if r.PostFormValue(patient.Username) == "savedselected" {
			foundSession.PatientsSelected = append(foundSession.PatientsSelected, &patient)
			usersToRemoveAvailable = append(usersToRemoveAvailable, patient.Username)
		}
	}
	for i := 0; i < len(foundSession.PatientsSelected); i++ {
		patient := *foundSession.PatientsSelected[i]
		if r.PostFormValue(patient.Username) == "savedassigned" {
			foundSession.PatientsAssigned = append(foundSession.PatientsAssigned, &patient)
			usersToRemoveSelected = append(usersToRemoveSelected, patient.Username)
		}
	}
	for i := 0; i < len(foundSession.PatientsSelected); i++ {
		patient := *foundSession.PatientsSelected[i]
		if r.PostFormValue(patient.Username) == "removeselected" {
			foundSession.PatientsAvailable = append(foundSession.PatientsAvailable, &patient)
			usersToRemoveSelected = append(usersToRemoveSelected, patient.Username)
		}
	}
	for i := 0; i < len(foundSession.PatientsAssigned); i++ {
		patient := *foundSession.PatientsAssigned[i]
		if r.PostFormValue(patient.Username) == "removeassigned" {
			foundSession.PatientsAvailable = append(foundSession.PatientsAvailable, &patient)
			usersToRemoveAssigned = append(usersToRemoveAssigned, patient.Username)
		}
	}
	if len(usersToRemoveAvailable) > 0 {
		for _, username := range usersToRemoveAvailable {
			for i, su := range foundSession.PatientsAvailable {
				if su.Username == username {
					foundSession.PatientsAvailable = append(foundSession.PatientsAvailable[:i], foundSession.PatientsAvailable[i+1:]...)
				}
			}
		}
	}
	if len(usersToRemoveSelected) > 0 {
		for _, username := range usersToRemoveSelected {
			for i, su := range foundSession.PatientsSelected {
				if su.Username == username {
					foundSession.PatientsSelected = append(foundSession.PatientsSelected[:i], foundSession.PatientsSelected[i+1:]...)
				}
			}
		}
	}
	if len(usersToRemoveAssigned) > 0 {
		for _, username := range usersToRemoveAssigned {
			for i, su := range foundSession.PatientsAssigned {
				if su.Username == username {
					foundSession.PatientsAssigned = append(foundSession.PatientsAssigned[:i], foundSession.PatientsAssigned[i+1:]...)
				}
			}
		}
	}
	if len(foundSession.PatientsAssigned) > 0 {
		for _, spuser := range foundSession.PatientsAssigned {
			username := spuser.Username
			spuserRecord, err := GetSpUserRecord(username, db)
			if err != nil {
				fmt.Println("Error Getting Record: ", err)
				return
			}
			duplicate := false
			for _, si := range spuserRecord.SessionsAssigned {
				if si.Title == foundSession.Information.Title {
					duplicate = true
				}
			}
			if !duplicate {
				spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned, foundSession.Information)
			}
			if len(spuserRecord.SessionsAvailable) > 0 {
				for i := 0; i < len(spuserRecord.SessionsAvailable); i++ {
					if spuserRecord.SessionsAvailable[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable[:i], spuserRecord.SessionsAvailable[i+1:]...)
					}
				}
			}
			if len(spuserRecord.SessionsSelected) > 0 {
				for i := 0; i < len(spuserRecord.SessionsSelected); i++ {
					if spuserRecord.SessionsSelected[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsSelected = append(spuserRecord.SessionsSelected[:i], spuserRecord.SessionsSelected[i+1:]...)
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
			err = spuserRecord.UpdateRecord(db)
			if err != nil {
				fmt.Println("Error Updating Record: ", err)
				return
			}
		}
	}

	if len(foundSession.PatientsAvailable) > 0 {
		for _, spuser := range foundSession.PatientsAvailable {
			username := spuser.Username
			spuserRecord, err := GetSpUserRecord(username, db)
			if err != nil {
				fmt.Println("Error Getting Record: ", err)
				return
			}
			duplicate := false
			for _, si := range spuserRecord.SessionsAvailable {
				if si.Title == foundSession.Information.Title {
					duplicate = true
				}
			}
			if !duplicate {
				spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable, foundSession.Information)
			}
			if len(spuserRecord.SessionsAssigned) > 0 {
				for i := 0; i < len(spuserRecord.SessionsAssigned); i++ {
					if spuserRecord.SessionsAssigned[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned[:i], spuserRecord.SessionsAssigned[i+1:]...)
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
			if len(spuserRecord.SessionsSelected) > 0 {
				for i := 0; i < len(spuserRecord.SessionsSelected); i++ {
					if spuserRecord.SessionsSelected[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsSelected = append(spuserRecord.SessionsSelected[:i], spuserRecord.SessionsSelected[i+1:]...)
					}
				}
			}
			if spuserRecord.TotalSessionsAssigned > 0 {
				spuserRecord.TotalSessionsAssigned = spuserRecord.TotalSessionsAssigned - 1
			}
			err = spuserRecord.UpdateRecord(db)
			if err != nil {
				fmt.Println("Error Updating Record: ", err)
				return
			}
		}
	}

	if len(foundSession.PatientsSelected) > 0 {
		for _, spuser := range foundSession.PatientsSelected {
			username := spuser.Username
			spuserRecord, err := GetSpUserRecord(username, db)
			if err != nil {
				fmt.Println("Error Getting Record: ", err)
				return
			}
			spuserRecord.SessionsSelected = append(spuserRecord.SessionsSelected, foundSession.Information)

			//delete any occurances of session from other session boxes
			if len(spuserRecord.SessionsAvailable) > 0 {
				for i := 0; i < len(spuserRecord.SessionsAvailable); i++ {
					if spuserRecord.SessionsAvailable[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable[:i], spuserRecord.SessionsAvailable[i+1:]...)
					}
				}
			}
			if len(spuserRecord.SessionsAssigned) > 0 {
				for i := 0; i < len(spuserRecord.SessionsAssigned); i++ {
					if spuserRecord.SessionsAssigned[i].Title == foundSession.Information.Title {
						spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned[:i], spuserRecord.SessionsAssigned[i+1:]...)
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
			// end
			if spuserRecord.TotalSessionsAssigned > 0 {
				spuserRecord.TotalSessionsAssigned = spuserRecord.TotalSessionsAssigned - 1
			}
			err = spuserRecord.UpdateRecord(db)
			if err != nil {
				fmt.Println("Error Updating Record: ", err)
				return
			}
		}
	}

	err = foundSession.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error updating record in assign sp", err)
		return
	}
	title = formatTitle(title)
	http.Redirect(w, r, "/dashboard#"+title, httpRedirectResponse)
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
	allSpUsers, _ := GetAllSpUserRecords(db)
	for _, spuser := range allSpUsers {
		for i, si := range spuser.SessionsAssigned {
			if foundSession.Information.Title == si.Title {
				spuser.SessionsAssigned = append(spuser.SessionsAssigned[:i], spuser.SessionsAssigned[i+1:]...)
				break
			}
		}
		for i, si := range spuser.SessionsPool {
			if foundSession.Information.Title == si.Title {
				spuser.SessionsPool = append(spuser.SessionsPool[:i], spuser.SessionsPool[i+1:]...)
				break
			}
		}
		for i, si := range spuser.SessionsAvailable {
			if foundSession.Information.Title == si.Title {
				spuser.SessionsAvailable = append(spuser.SessionsAvailable[:i], spuser.SessionsAvailable[i+1:]...)
				break
			}
		}
		for i, si := range spuser.SessionsUnavailable {
			if foundSession.Information.Title == si.Title {
				spuser.SessionsUnavailable = append(spuser.SessionsUnavailable[:i], spuser.SessionsUnavailable[i+1:]...)
				break
			}
		}

		err = spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating spuser record in delete session: ", err)
			return
		}
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
			if availableSessionRecord.Information.Title == spuser.SessionsAvailable[i].Title &&
				availableSessionRecord.Information.Date == spuser.SessionsAvailable[i].Date &&
				availableSessionRecord.Information.Location == spuser.SessionsAvailable[i].Location &&
				availableSessionRecord.Information.StartTime == spuser.SessionsAvailable[i].StartTime &&
				availableSessionRecord.Information.EndTime == spuser.SessionsAvailable[i].EndTime &&
				availableSessionRecord.Information.Description == spuser.SessionsAvailable[i].Description {
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
	title := formatTitle(availableSessionRecord.Information.Title)

	http.Redirect(w, r, "/dashboard#"+title, httpRedirectResponse)
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
			if notAvailableSessionRecord.Information.Title == spuser.SessionsUnavailable[i].Title &&
				notAvailableSessionRecord.Information.Date == spuser.SessionsUnavailable[i].Date &&
				notAvailableSessionRecord.Information.Location == spuser.SessionsUnavailable[i].Location &&
				notAvailableSessionRecord.Information.StartTime == spuser.SessionsUnavailable[i].StartTime &&
				notAvailableSessionRecord.Information.EndTime == spuser.SessionsUnavailable[i].EndTime &&
				notAvailableSessionRecord.Information.Description == spuser.SessionsUnavailable[i].Description {
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
	title := formatTitle(notAvailableSessionRecord.Information.Title)
	http.Redirect(w, r, "/dashboard#"+title, httpRedirectResponse)
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
	// 60 seconds * 60 minutes = 1 hour max age
	session.Options.MaxAge = 60 * 60 // in seconds
	session.Options.Secure = true

	username := r.PostFormValue("userid")
	username = strings.ToLower(username)
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
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
	}
	_, err := GetSpManagerRecord(session.Values["username"].(string), db)
	if err != nil {
		http.Redirect(w, r, "/login", httpRedirectResponse)
	}
	SpUsersBox := make(SpUsersBox, 0)
	SpManagersBox := make(SpManagersBox, 0)
	SessionsBox := make(SpSessionsBox, 0)

	allSpUsers, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error getting all sp user records in json: ", err)
	}
	allManagers, err := GetAllSpManagerRecords(db)
	if err != nil {
		fmt.Println("Error getting all sp manager records in json: ", err)
	}
	allSessions, err := GetAllSessionRecords(db)
	if err != nil {
		fmt.Println("Error getting all session records in json: ", err)
	}

	for _, spuser := range allSpUsers {
		SpUsersBox = append(SpUsersBox, spuser)
	}

	for _, spmanager := range allManagers {
		SpManagersBox = append(SpManagersBox, spmanager)
	}

	for _, session := range allSessions {
		SessionsBox = append(SessionsBox, session)
	}

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
	spRecords, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Get all sp user records in createSPRecord: ", err)
		return
	}
	duplicate := false
	for _, su := range spRecords {
		fmt.Println(su.Username, username)
		if su.Username == username {
			duplicate = true
		}
	}
	if !duplicate {
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
