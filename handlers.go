package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func sessionEqual(session1, session2 *SessionInfo) bool {
	if session1.Title != session2.Title {
		return false
	}
	if session1.Date != session2.Date {
		return false
	}
	if session1.StartTime != session2.StartTime {
		return false
	}
	if session1.EndTime != session2.EndTime {
		return false
	}
	if session1.Location != session2.Location {
		return false
	}
	if session1.Description != session2.Description {
		return false
	}
	return true
}
func formatTitle(title string) string {
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "@", "")
	title = strings.ReplaceAll(title, "(", "")
	title = strings.ReplaceAll(title, ")", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, "|", "")
	title = strings.ReplaceAll(title, "+", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "#", "")
	title = strings.ReplaceAll(title, "%", "")
	title = strings.ReplaceAll(title, "$", "")
	title = strings.ReplaceAll(title, "!", "")
	title = strings.ReplaceAll(title, "^", "")
	title = strings.ReplaceAll(title, "&", "")
	title = strings.ReplaceAll(title, "[", "")
	title = strings.ReplaceAll(title, "]", "")
	title = strings.ReplaceAll(title, "{", "")
	title = strings.ReplaceAll(title, "}", "")
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, ";", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "<", "")
	title = strings.ReplaceAll(title, ">", "")
	return strings.Join(strings.Fields(strings.TrimSpace(title)), "")
}

func formatDate(date string) string {
	timeT, _ := time.Parse("01/02/2006", date)
	return timeT.Format("Monday, January 02, 2006")

}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetSessionArchiveDates(sessions []*Session) []string {
	dates := []string{}
	for _, session := range sessions {
		time, _ := time.Parse("01/02/2006", session.Information.Date)
		date := time.Format("January, 2006")
		// if date in dates continue else add it
		dates = append(dates, date)
		dates = removeDuplicate(dates)
	}
	return dates
}

func GetSessionInfoArchiveDates(sessions []*SessionInfo) []string {
	dates := []string{}
	for _, session := range sessions {
		time, _ := time.Parse("01/02/2006", session.Date)
		date := time.Format("January, 2006")
		// if date in dates continue else add it
		dates = append(dates, date)
		dates = removeDuplicate(dates)
	}
	return dates
}

func CheckExpired(expireddate string) bool {
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	if expireddate != "" {
		expiredDateParsed, _ := time.Parse("01/02/2006", expireddate)
		currentDate := time.Now().In(loc)
		return currentDate.After(expiredDateParsed)
	} else {
		return false
	}
}

func CheckNotExpired(expireddate string) bool {
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	if expireddate != "" {
		expiredDateParsed, _ := time.Parse("01/02/2006", expireddate)
		currentDate := time.Now().In(loc)
		return currentDate.Before(expiredDateParsed)
	} else {
		return false
	}
}

func ExpirationDateSet(expireddate string) bool {
	if expireddate != "" {
		return true
	} else {
		return false
	}
}

func pastSession(date string) bool {
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	sessionDate := date
	currentDate := time.Now().In(loc).AddDate(0, 0, -1)

	sessionDateParsed, _ := time.Parse("01/02/2006", sessionDate)
	return currentDate.After(sessionDateParsed)
}

func notPastSession(date string) bool {
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	sessionDate := date
	currentDate := time.Now().In(loc).AddDate(0, 0, -1)

	sessionDateParsed, _ := time.Parse("01/02/2006", sessionDate)
	return currentDate.Before(sessionDateParsed)
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

func sortSpUserByLastName(assignedPatients []*SpUser) []*SpUser {
	sort.Slice(assignedPatients, func(i int, j int) bool {
		return assignedPatients[i].Name.Last < assignedPatients[j].Name.Last
	})
	return assignedPatients
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	} else {
		t, _ := template.ParseFiles("templates/html-boilerplate.html", "templates/login-content.html")
		t.ExecuteTemplate(w, "html-boilerplate", "")
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", httpRedirectResponse)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if session.Options.MaxAge != 3600 {
		session.Options.MaxAge = 60 * 60
		session.Options.Secure = true
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	} else {
		var t *template.Template
		var spmanager SpManager
		isSpManager := false
		loc, err := time.LoadLocation("EST")
		if err != nil {
			fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
		}
		dashboard_content := DashboardContent{
			Date: time.Now().In(loc).Format("Monday, January 02, 2006"),
		}
		session_records, err := GetAllSessionInfoRecords(db)
		if err != nil {
			fmt.Println("Error Get All Session Records: ", err)
		}
		spuser, err := GetSpUserRecord(session.Values["username"].(string), db)

		if err != nil {
			// Cannot Find SP User Account Must Be Manager
			spmanager, err = GetSpManagerRecord(session.Values["username"].(string), db)
			if err != nil {
				fmt.Println("Error Get Manager record in dashboard: ", err)
				return
			}
			session_records_manager, err := GetAllSessionRecords(db)
			if err != nil {
				fmt.Println("Error Get All Session records in dashboard: ", err)
				return
			}

			if session.Values["dateFilter"] == nil {
				timenow := time.Now()
				dateFilter := timenow.Format("January, 2006")
				session.Values["dateFilter"] = dateFilter
			}

			if r.PostFormValue("date") != "allsessions" {
				if r.PostFormValue("date") != "" {
					session.Values["dateFilter"] = r.PostFormValue("date")
					dashboard_content.SelectedDate = r.PostFormValue("date")
					session_records_manager_new := make([]*Session, 0)
					for _, s := range session_records_manager {
						time, _ := time.Parse("01/02/2006", s.Information.Date)
						date := time.Format("January, 2006")
						if r.PostFormValue("date") == date {
							session_records_manager_new = append(session_records_manager_new, s)
						}
					}
					session_records_manager = session_records_manager_new
				} else if session.Values["dateFilter"] != nil && session.Values["dateFilter"] != "allsessions" {
					dashboard_content.SelectedDate = session.Values["dateFilter"].(string)
					session_records_manager_new := make([]*Session, 0)
					for _, s := range session_records_manager {
						time, _ := time.Parse("01/02/2006", s.Information.Date)
						date := time.Format("January, 2006")
						if session.Values["dateFilter"] == date {
							session_records_manager_new = append(session_records_manager_new, s)
						}
					}
					session_records_manager = session_records_manager_new
				} else {
					dashboard_content.SelectedDate = "All Sessions"
				}
			} else if r.PostFormValue("date") == "allsessions" {
				dashboard_content.SelectedDate = "All Sessions"
				session.Values["dateFilter"] = "allsessions"
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
			session_records, _ := GetAllSessionRecords(db)
			spuser_records, err := GetAllSpUserRecords(db)
			spmanager.AssignedPatients = spuser_records
			if err != nil {
				fmt.Println("Error Get All User records in dashboard: ", err)
			}
			dashboard_content.Archives = GetSessionArchiveDates(session_records)
			sort.Slice(dashboard_content.Archives, func(i int, j int) bool {
				iDate := dashboard_content.Archives[i]
				jDate := dashboard_content.Archives[j]

				iParsed, _ := time.Parse("January, 2006", iDate)
				jParsed, _ := time.Parse("January, 2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.Role = "Manager"
			dashboard_content.User = spmanager
			isSpManager = true
		} else {
			// Sp User Section of the Dashboard
			spuser.SessionsPool = make([]*SessionInfo, 0)
			sessions_viewed := append(spuser.SessionsAvailable, spuser.SessionsUnavailable...)
			sessions_viewed = append(sessions_viewed, spuser.SessionsAssigned...)
			for _, session_info := range session_records {
				viewed_session := false
				for _, session_viewed_info := range sessions_viewed {
					if sessionEqual(session_info, session_viewed_info) {
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

			if session.Values["dateFilter"] == nil {
				timenow := time.Now()
				dateFilter := timenow.Format("January, 2006")
				session.Values["dateFilter"] = dateFilter
			}
			if r.PostFormValue("date") == "" {
				loc, err := time.LoadLocation("EST")
				if err != nil {
					fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
				}
				timenow := time.Now().In(loc).AddDate(0, 1, 0)
				dateFilter := timenow.Format("January, 2006")
				session.Values["dateFilter"] = dateFilter
				dashboard_content.SelectedDate = dateFilter
				newSessionsSorted := make([]*SessionInfo, 0)
				for _, s := range spuser.SessionsSorted {
					time, _ := time.Parse("01/02/2006", s.Date)
					date := time.Format("January, 2006")
					if dateFilter == date {
						newSessionsSorted = append(newSessionsSorted, s)
					}
				}
				spuser.SessionsSorted = newSessionsSorted

			} else if r.PostFormValue("date") == "currentMonth" {
				loc, err := time.LoadLocation("EST")
				if err != nil {
					fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
				}
				timenow := time.Now().In(loc).AddDate(0, 1, 0)
				dateFilter := timenow.Format("January, 2006")
				session.Values["dateFilter"] = dateFilter
				dashboard_content.SelectedDate = dateFilter
				newSessionsSorted := make([]*SessionInfo, 0)
				for _, s := range spuser.SessionsSorted {
					time, _ := time.Parse("01/02/2006", s.Date)
					date := time.Format("January, 2006")
					if dateFilter == date {
						newSessionsSorted = append(newSessionsSorted, s)
					}
				}
				spuser.SessionsSorted = newSessionsSorted
			} else if r.PostFormValue("date") != "allsessions" {
				if r.PostFormValue("date") != "" {
					session.Values["dateFilter"] = r.PostFormValue("date")
					dashboard_content.SelectedDate = r.PostFormValue("date")
					newSessionsAssigned := make([]*SessionInfo, 0)
					for _, s := range spuser.SessionsAssigned {
						time, _ := time.Parse("01/02/2006", s.Date)
						date := time.Format("January, 2006")
						if r.PostFormValue("date") == date {
							newSessionsAssigned = append(newSessionsAssigned, s)
						}
					}
					spuser.SessionsSorted = newSessionsAssigned
				} else if session.Values["dateFilter"] != nil && session.Values["dateFilter"] != "allsessions" {
					dashboard_content.SelectedDate = session.Values["dateFilter"].(string)
					newSessionsAssigned := make([]*SessionInfo, 0)
					for _, s := range spuser.SessionsAssigned {
						time, _ := time.Parse("01/02/2006", s.Date)
						date := time.Format("January, 2006")
						if session.Values["dateFilter"] == date {
							newSessionsAssigned = append(newSessionsAssigned, s)
						}
					}
					spuser.SessionsSorted = newSessionsAssigned
				} else {
					dashboard_content.SelectedDate = "All Assigned"
				}
			} else if r.PostFormValue("date") == "allsessions" {
				dashboard_content.SelectedDate = "All Assigned"
				spuser.SessionsSorted = spuser.SessionsAssigned
				session.Values["dateFilter"] = "allsessions"
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
			dashboard_content.Archives = GetSessionInfoArchiveDates(spuser.SessionsAssigned)
			sort.Slice(dashboard_content.Archives, func(i int, j int) bool {
				iDate := dashboard_content.Archives[i]
				jDate := dashboard_content.Archives[j]

				iParsed, _ := time.Parse("January, 2006", iDate)
				jParsed, _ := time.Parse("January, 2006", jDate)

				return iParsed.Before(jParsed)
			})
			dashboard_content.Role = "Standardized Patient"
			dashboard_content.User = spuser
		}
		funcMap := template.FuncMap{
			"formatTitle":           formatTitle,
			"formatDate":            formatDate,
			"ExpirationDateSet":     ExpirationDateSet,
			"CheckExpired":          CheckExpired,
			"CheckNotExpired":       CheckNotExpired,
			"sortSessionInfoByDate": sortSessionInfoByDate,
			"sortSessionByDate":     sortSessionByDate,
			"sortSpUserByLastName":  sortSpUserByLastName,
			"StatusAssigned":        StatusAssigned,
			"StatusNoResponse":      StatusNoResponse,
			"StatusAvailable":       StatusAvailable,
			"StatusUnavailable":     StatusUnavailable,
			"pastSession":           pastSession,
			"notPastSession":        notPastSession,
		}
		t = template.New("templates/html-boilerplate.html").Funcs(funcMap)
		if !isSpManager {
			t, _ = t.ParseFiles("templates/html-boilerplate.html", "templates/dashboard-content.html", "templates/session-content-available.html", "templates/user-settings.html")
		} else {
			t, _ = t.ParseFiles("templates/html-boilerplate.html", "templates/dashboard-content-manager.html", "templates/session-content-manager.html")
		}
		t.ExecuteTemplate(w, "html-boilerplate", dashboard_content)
	}
}

func createsession(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	starttime := r.PostFormValue("starttime")
	endtime := r.PostFormValue("endtime")
	location := r.PostFormValue("location")
	description := r.PostFormValue("description")
	patientsneeded, err := strconv.Atoi(r.PostFormValue("patientsneeded"))
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if err != nil {
		fmt.Println("Error converting patients needed to integer")
	}
	newSession := Session{}.Create(title, date, starttime, endtime, location, description)
	loc, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Println("Error in LoadLocation CheckExpirationDate :", err)
	}
	timenow := time.Now().In(loc)
	datetime, _ := time.Parse("01/02/2006", date)
	dateFilter := datetime.Format("January, 2006")
	session.Values["dateFilter"] = dateFilter
	session.Save(r, w)
	newSession.Information.CreatedDate = timenow.Format("01/02/2006")
	newSession.Information.ExpiredDate = ""
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

	allSpUsers, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error Getting all SP User records: ", err)
	}
	// For Every SP record in the database
	for _, su := range allSpUsers {
		// collect all sessions except for assigned ones
		allSessions := append(su.SessionsAssigned, su.SessionsAvailable...)
		allSessions = append(allSessions, su.SessionsUnavailable...)
		allSessions = append(allSessions, su.SessionsPool...)
		allSessions = append(allSessions, su.SessionsAssigned...)
		// find the session needed to be updated
		for _, si := range allSessions {

			if sessionEqual(foundSession.Information, si) {
				si.Title = newtitle
				si.Date = newdate
				si.StartTime = newstarttime
				si.EndTime = newendtime
				si.Location = newlocation
				si.Description = newdescription
			}
		}
		// update the found session record
		su.UpdateRecord(db)
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
	uniqueID := formatTitle(foundSession.Information.Title)
	uniqueID = formatTitle(uniqueID + foundSession.Information.Date)
	uniqueID = formatTitle(uniqueID + foundSession.Information.StartTime)
	uniqueID = formatTitle(uniqueID + foundSession.Information.EndTime)
	uniqueID = formatTitle(uniqueID + foundSession.Information.Location)
	http.Redirect(w, r, "/dashboard#"+uniqueID, httpRedirectResponse)
}

func assignsp(w http.ResponseWriter, r *http.Request) {
	//Get the information for the session
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
	// end information for the session
	if err != nil {
		fmt.Println("Error getting record in database", err)
		return
	}
	// Unclick the check mark on the session
	foundSession.Information.CheckMarkAssigned = false
	// Setup who to remove from each category to place
	// in another category
	usersToRemoveAvailable := make([]string, 0)
	usersToRemoveAssigned := make([]string, 0)
	usersToRemoveSelected := make([]string, 0)

	// Users moved from available to the selected section
	for i := 0; i < len(foundSession.PatientsAvailable); i++ {
		patient := *foundSession.PatientsAvailable[i]
		if r.PostFormValue(patient.Username) == "savedselected" {
			foundSession.PatientsSelected = append(foundSession.PatientsSelected, &patient)
			usersToRemoveAvailable = append(usersToRemoveAvailable, patient.Username)
		}
	}
	// User moved from selected to the assigned section
	// or User moved from selected to the available section
	for i := 0; i < len(foundSession.PatientsSelected); i++ {
		patient := *foundSession.PatientsSelected[i]
		if r.PostFormValue(patient.Username) == "savedassigned" {
			foundSession.PatientsAssigned = append(foundSession.PatientsAssigned, &patient)
			usersToRemoveSelected = append(usersToRemoveSelected, patient.Username)
		}
		if r.PostFormValue(patient.Username) == "removeselected" {
			foundSession.PatientsAvailable = append(foundSession.PatientsAvailable, &patient)
			usersToRemoveSelected = append(usersToRemoveSelected, patient.Username)
		}

	}
	// User moved from assigned to the available section
	for i := 0; i < len(foundSession.PatientsAssigned); i++ {
		patient := *foundSession.PatientsAssigned[i]
		if r.PostFormValue(patient.Username) == "removeassigned" {
			foundSession.PatientsAvailable = append(foundSession.PatientsAvailable, &patient)
			usersToRemoveAssigned = append(usersToRemoveAssigned, patient.Username)
		}
	}
	// Remove all users from remove available because they have been placed elsewhere
	if len(usersToRemoveAvailable) > 0 {
		for _, username := range usersToRemoveAvailable {
			for i, su := range foundSession.PatientsAvailable {
				if su.Username == username {
					foundSession.PatientsAvailable = append(foundSession.PatientsAvailable[:i], foundSession.PatientsAvailable[i+1:]...)
				}
			}
		}
	}
	// Remove all users from remove selected because they have been placed elsewhere
	if len(usersToRemoveSelected) > 0 {
		for _, username := range usersToRemoveSelected {
			for i, su := range foundSession.PatientsSelected {
				if su.Username == username {
					foundSession.PatientsSelected = append(foundSession.PatientsSelected[:i], foundSession.PatientsSelected[i+1:]...)
				}
			}
		}
	}
	// Remove all users from remove assigned because they have been placed elsewhere
	if len(usersToRemoveAssigned) > 0 {
		for _, username := range usersToRemoveAssigned {
			for i, su := range foundSession.PatientsAssigned {
				if su.Username == username {
					foundSession.PatientsAssigned = append(foundSession.PatientsAssigned[:i], foundSession.PatientsAssigned[i+1:]...)
				}
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
				if sessionEqual(si, foundSession.Information) {
					duplicate = true
				}
			}
			if !duplicate {
				spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable, foundSession.Information)

				if len(spuserRecord.SessionsAssigned) > 0 {
					for i := 0; i < len(spuserRecord.SessionsAssigned); i++ {
						if sessionEqual(spuserRecord.SessionsAssigned[i], foundSession.Information) {
							spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned[:i], spuserRecord.SessionsAssigned[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsPool) > 0 {
					for i := 0; i < len(spuserRecord.SessionsPool); i++ {
						if sessionEqual(spuserRecord.SessionsPool[i], foundSession.Information) {
							spuserRecord.SessionsPool = append(spuserRecord.SessionsPool[:i], spuserRecord.SessionsPool[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsUnavailable) > 0 {
					for i := 0; i < len(spuserRecord.SessionsUnavailable); i++ {
						if sessionEqual(spuserRecord.SessionsUnavailable[i], foundSession.Information) {
							spuserRecord.SessionsUnavailable = append(spuserRecord.SessionsUnavailable[:i], spuserRecord.SessionsUnavailable[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsSelected) > 0 {
					for i := 0; i < len(spuserRecord.SessionsSelected); i++ {
						if sessionEqual(spuserRecord.SessionsSelected[i], foundSession.Information) {
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
	}

	if len(foundSession.PatientsSelected) > 0 {
		for _, spuser := range foundSession.PatientsSelected {
			username := spuser.Username
			spuserRecord, err := GetSpUserRecord(username, db)
			if err != nil {
				fmt.Println("Error Getting Record: ", err)
				return
			}
			duplicate := false
			for _, si := range spuserRecord.SessionsSelected {
				if sessionEqual(si, foundSession.Information) {
					duplicate = true
				}
			}
			if !duplicate {
				spuserRecord.SessionsSelected = append(spuserRecord.SessionsSelected, foundSession.Information)

				//delete any occurances of session from other session boxes
				if len(spuserRecord.SessionsAvailable) > 0 {
					for i := 0; i < len(spuserRecord.SessionsAvailable); i++ {
						if sessionEqual(spuserRecord.SessionsAvailable[i], foundSession.Information) {
							spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable[:i], spuserRecord.SessionsAvailable[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsAssigned) > 0 {
					for i := 0; i < len(spuserRecord.SessionsAssigned); i++ {
						if sessionEqual(spuserRecord.SessionsAssigned[i], foundSession.Information) {
							spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned[:i], spuserRecord.SessionsAssigned[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsPool) > 0 {
					for i := 0; i < len(spuserRecord.SessionsPool); i++ {
						if sessionEqual(spuserRecord.SessionsPool[i], foundSession.Information) {
							spuserRecord.SessionsPool = append(spuserRecord.SessionsPool[:i], spuserRecord.SessionsPool[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsUnavailable) > 0 {
					for i := 0; i < len(spuserRecord.SessionsUnavailable); i++ {
						if sessionEqual(spuserRecord.SessionsUnavailable[i], foundSession.Information) {
							spuserRecord.SessionsUnavailable = append(spuserRecord.SessionsUnavailable[:i], spuserRecord.SessionsUnavailable[i+1:]...)
						}
					}
				}
				err = spuserRecord.UpdateRecord(db)
				if err != nil {
					fmt.Println("Error Updating Record: ", err)
					return
				}
			}
		}
	}

	err = foundSession.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error updating record in assign sp", err)
		return
	}
	title = formatTitle(title + date + starttime + endtime + location)
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
			if sessionEqual(availableSessionRecord.Information, spuser.SessionsAvailable[i]) {
				duplicate = true
			}
		}
	}
	if !duplicate {
		for i := 0; i < len(spuser.SessionsPool); i++ {
			// Removed Session From Sessions Pool
			if sessionEqual(spuser.SessionsPool[i], &sessionInfo) {
				spuser.SessionsPool = append(spuser.SessionsPool[:i], spuser.SessionsPool[i+1:]...)
			}
		}
		for i := 0; i < len(spuser.SessionsUnavailable); i++ {
			//Removed Session from Sessions Unavailable
			if sessionEqual(spuser.SessionsUnavailable[i], &sessionInfo) {
				spuser.SessionsUnavailable = append(spuser.SessionsUnavailable[:i], spuser.SessionsUnavailable[i+1:]...)
			}
		}
		//Add session to SessionsAvailable
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
	title := formatTitle(availableSessionRecord.Information.Title + availableSessionRecord.Information.Date +
		availableSessionRecord.Information.StartTime + availableSessionRecord.Information.EndTime +
		availableSessionRecord.Information.Location)

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
			if sessionEqual(notAvailableSessionRecord.Information, spuser.SessionsUnavailable[i]) {
				duplicate = true
			}
		}
	}
	if !duplicate {
		for i := 0; i < len(spuser.SessionsPool); i++ {
			if sessionEqual(spuser.SessionsPool[i], &sessionInfo) {
				spuser.SessionsPool = append(spuser.SessionsPool[:i], spuser.SessionsPool[i+1:]...)
			}
		}
		for i := 0; i < len(spuser.SessionsAvailable); i++ {
			if sessionEqual(spuser.SessionsAvailable[i], &sessionInfo) {
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
	title := formatTitle(notAvailableSessionRecord.Information.Title + notAvailableSessionRecord.Information.Date +
		notAvailableSessionRecord.Information.StartTime + notAvailableSessionRecord.Information.EndTime +
		notAvailableSessionRecord.Information.Location)
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
		spmanager, err := GetSpManagerRecord(session.Values["username"].(string), db)
		if err != nil {
			fmt.Println("Error get spmanager record in changeemail: ", err)
			return
		}
		newemail := r.PostFormValue("newemail")
		spmanager.Email = newemail
		err = spmanager.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating user record in change email handler : ", err)
		}
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	} else {
		newemail := r.PostFormValue("newemail")
		spuser.Email = newemail
		err = spuser.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating user record in change email handler : ", err)
		}
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	}
}

func changepassword(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", httpRedirectResponse)
		return
	}
	spuser, err := GetSpUserRecord(session.Values["username"].(string), db)
	if err != nil {
		spmanager, err := GetSpManagerRecord(session.Values["username"].(string), db)
		if err != nil {
			fmt.Println("Error getting spmanager record: ", err)
		}
		newPassword := r.PostFormValue("newpassword")
		newPasswordConfirmed := r.PostFormValue("newpasswordconfirmed")
		if newPassword == newPasswordConfirmed {
			hashedPassword, err := HashPassword(newPassword)
			if err != nil {
				fmt.Println("Error Hashing Password")
				http.Redirect(w, r, "/dashboard", httpRedirectResponse)
			}
			spmanager.Password = hashedPassword
		}
		err = spmanager.UpdateRecord(db)
		if err != nil {
			fmt.Println("error updating user record in change password", err)
		}
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	} else {
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
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	}

}

func toggleshowsession(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show Session Toggled")
	sessionInfo := SessionInfo{
		Title:       r.PostFormValue("title"),
		Date:        r.PostFormValue("date"),
		StartTime:   r.PostFormValue("starttime"),
		EndTime:     r.PostFormValue("endtime"),
		Location:    r.PostFormValue("location"),
		Description: r.PostFormValue("description"),
	}
	// Get the Session
	availableSessionRecord, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	// Update Session Details
	if availableSessionRecord.Information.ShowSession {
		availableSessionRecord.Information.ShowSession = false
	} else {
		availableSessionRecord.Information.ShowSession = true
	}
	// Update SP users with those session details
	allSpUsers, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Error Getting all SP User records: ", err)
	}
	for _, su := range allSpUsers {
		allSessions := append(su.SessionsAssigned, su.SessionsAvailable...)
		allSessions = append(allSessions, su.SessionsUnavailable...)
		allSessions = append(allSessions, su.SessionsPool...)
		// Find the session in the SP records to update
		for _, si := range allSessions {
			if availableSessionRecord.Information.Title == si.Title {
				si.ShowSession = availableSessionRecord.Information.ShowSession
			}
		}
		// update each SP User
		su.UpdateRecord(db)
	}
	// Update Original Session Record
	err = availableSessionRecord.UpdateRecord(db)
	if err != nil {
		fmt.Println("Error Updating Session Record in Toggle Show Session : ", err)
	}
	uniqueID := formatTitle(availableSessionRecord.Information.Title)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.Date)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.StartTime)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.EndTime)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.Location)
	// send user back to the correct position on the page
	http.Redirect(w, r, "/dashboard#"+uniqueID, httpRedirectResponse)
}

func togglehourglass(w http.ResponseWriter, r *http.Request) {
	//Obtain Session from Session Information
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
	uniqueID := formatTitle(availableSessionRecord.Information.Title)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.Date)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.StartTime)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.EndTime)
	uniqueID = formatTitle(uniqueID + availableSessionRecord.Information.Location)
	if availableSessionRecord.Information.ExpiredDate != "" {
		if CheckExpired(availableSessionRecord.Information.ExpiredDate) {
			availableSessionRecord.Information.ExpiredDate = ""
			availableSessionRecord.Information.ShowSession = false
			err = availableSessionRecord.UpdateRecord(db)
			if err != nil {
				fmt.Println("Error updating expired record in toggle hour glass", err)
			}
			http.Redirect(w, r, "/dashboard#"+uniqueID, httpRedirectResponse)
		}
	} else {
		//fmt.Println("Session Expires: ", availableSessionRecord.Information.ExpiredDate)
		// end Session Information block
		// Load Eastern Standard Time
		loc, err := time.LoadLocation("EST")
		if err != nil {
			fmt.Println("Error loading location time in toggleHourGlass")
		}
		timenow := time.Now().In(loc)
		// end Load of Eastern Standard Time
		// change expiration date of session
		availableSessionRecord.Information.ExpiredDate = timenow.AddDate(0, 0, 5).Format("01/02/2006")
		availableSessionRecord.Information.ShowSession = true
		// get All SP Records from Database
		allSpUsers, err := GetAllSpUserRecords(db)
		if err != nil {
			fmt.Println("Error Getting all SP User records: ", err)
		}
		// For Every SP record in the database
		for _, su := range allSpUsers {
			// collect all sessions except for assigned ones
			allSessions := append(su.SessionsAssigned, su.SessionsAvailable...)
			allSessions = append(allSessions, su.SessionsUnavailable...)
			allSessions = append(allSessions, su.SessionsPool...)
			// find the session needed to be updated
			for _, si := range allSessions {

				if sessionEqual(availableSessionRecord.Information, si) {
					si.ExpiredDate = availableSessionRecord.Information.ExpiredDate
				}
			}
			// update the found session record
			su.UpdateRecord(db)
		}
		//fmt.Println("New Session Expiration: ", availableSessionRecord.Information.ExpiredDate)
		err = availableSessionRecord.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error Updating Session Record in Toggle Show Session : ", err)
		}
		http.Redirect(w, r, "/dashboard#"+uniqueID, httpRedirectResponse)
	}
}

func togglechecksquare(w http.ResponseWriter, r *http.Request) {
	sessionInfo := SessionInfo{
		Title:       r.PostFormValue("title"),
		Date:        r.PostFormValue("date"),
		StartTime:   r.PostFormValue("starttime"),
		EndTime:     r.PostFormValue("endtime"),
		Location:    r.PostFormValue("location"),
		Description: r.PostFormValue("description"),
	}
	//Put this session into all the SPs assigned sessionsAssignedBox
	foundSession, err := GetSessionRecord(&sessionInfo, db)
	if err != nil {
		fmt.Println("Error GetSessionRecord in signupavailable", err)
	}
	if !foundSession.Information.CheckMarkAssigned {
		foundSession.Information.CheckMarkAssigned = true
		foundSession.Information.ShowSession = true
		err = foundSession.UpdateRecord(db)
		if err != nil {
			fmt.Println("Error updating record in togglecheckassign ", err)
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
				if si == foundSession.Information {
					duplicate = true
				}
			}
			if !duplicate {
				spuserRecord.SessionsAssigned = append(spuserRecord.SessionsAssigned, foundSession.Information)

				if len(spuserRecord.SessionsAvailable) > 0 {
					for i := 0; i < len(spuserRecord.SessionsAvailable); i++ {
						if sessionEqual(spuserRecord.SessionsAvailable[i], foundSession.Information) {
							spuserRecord.SessionsAvailable = append(spuserRecord.SessionsAvailable[:i], spuserRecord.SessionsAvailable[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsSelected) > 0 {
					for i := 0; i < len(spuserRecord.SessionsSelected); i++ {
						if sessionEqual(spuserRecord.SessionsSelected[i], foundSession.Information) {
							spuserRecord.SessionsSelected = append(spuserRecord.SessionsSelected[:i], spuserRecord.SessionsSelected[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsPool) > 0 {
					for i := 0; i < len(spuserRecord.SessionsPool); i++ {
						if sessionEqual(spuserRecord.SessionsPool[i], foundSession.Information) {
							spuserRecord.SessionsPool = append(spuserRecord.SessionsPool[:i], spuserRecord.SessionsPool[i+1:]...)
						}
					}
				}
				if len(spuserRecord.SessionsUnavailable) > 0 {
					for i := 0; i < len(spuserRecord.SessionsUnavailable); i++ {
						if sessionEqual(spuserRecord.SessionsUnavailable[i], foundSession.Information) {
							spuserRecord.SessionsUnavailable = append(spuserRecord.SessionsUnavailable[:i], spuserRecord.SessionsUnavailable[i+1:]...)
						}
					}
				}
				spuserRecord.TotalSessionsAssigned = spuserRecord.TotalSessionsAssigned + 1
			}
			err = spuserRecord.UpdateRecord(db)
			if err != nil {
				fmt.Println("Error Updating Record: ", err)
				return
			}
		}
	}
	uniqueID := formatTitle(foundSession.Information.Title)
	uniqueID = formatTitle(uniqueID + foundSession.Information.Date)
	uniqueID = formatTitle(uniqueID + foundSession.Information.StartTime)
	uniqueID = formatTitle(uniqueID + foundSession.Information.EndTime)
	uniqueID = formatTitle(uniqueID + foundSession.Information.Location)

	http.Redirect(w, r, "/dashboard#"+uniqueID, httpRedirectResponse)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionAuthSPCalendar")

	username := r.PostFormValue("userid")
	username = strings.ToLower(username)
	password := r.PostFormValue("password")
	spuser, err := GetSpUserRecord(username, db)
	if err != nil {
		spmanager, err := GetSpManagerRecord(username, db)
		if err != nil {
			t, _ := template.ParseFiles("templates/html-boilerplate.html", "templates/login-content.html")
			t.ExecuteTemplate(w, "html-boilerplate", "")
		} else if !CheckPasswordHash(password, spmanager.Password) {
			http.Redirect(w, r, "/login", httpRedirectResponse)
		} else {
			session.Values["authenticated"] = true
			session.Values["username"] = spmanager.Username
		}
	} else {
		if !CheckPasswordHash(password, spuser.Password) {
			t, _ := template.ParseFiles("templates/html-boilerplate.html", "templates/login-content.html")
			t.ExecuteTemplate(w, "html-boilerplate", "")
		} else {
			session.Values["authenticated"] = true
			session.Values["username"] = spuser.Username
		}
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		t, _ := template.ParseFiles("templates/html-boilerplate.html", "templates/login-content.html")
		t.ExecuteTemplate(w, "html-boilerplate", "")

	} else {
		session.Save(r, w)
		http.Redirect(w, r, "/dashboard", httpRedirectResponse)
	}
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

func sessionbackup(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("BackUpSessions")
	if err == nil {
		data, err := io.ReadAll(file)
		if err == nil {
			SpUsersBox := make(SpUsersBox, 0)
			SpManagersBox := make(SpManagersBox, 0)
			SessionsBox := make(SpSessionsBox, 0)
			HospitalCalendar := HospitalCalendar{
				Users:    SpUsersBox,
				Managers: SpManagersBox,
				Sessions: SessionsBox,
			}
			err = json.Unmarshal(data, &HospitalCalendar)
			if err != nil {
				fmt.Fprintln(w, "Error Uploading File!!")
			}
			for _, session := range HospitalCalendar.Sessions {
				err = session.MakeRecord(db)
				if err != nil {
					fmt.Println(w, "Error Creating Session Record")
				}
			}
			for _, user := range HospitalCalendar.Users {
				err = user.MakeRecord(db)
				if err != nil {
					fmt.Fprintln(w, "Error Creating User Record")
				}
			}
			for _, manager := range HospitalCalendar.Managers {
				err = manager.UpdateRecord(db)
				if err != nil {
					fmt.Fprintln(w, "Error Creating User Record")
				}
			}
			http.Redirect(w, r, "/dashboard", httpRedirectResponse)

		}
	}
}

func createSPRecord(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	username := r.PostFormValue("email")
	spRecords, err := GetAllSpUserRecords(db)
	if err != nil {
		fmt.Println("Get all sp user records in createSPRecord: ", err)
		return
	}
	duplicate := false
	for _, su := range spRecords {
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
	http.Redirect(w, r, "/dashboard", httpRedirectResponse)
}
