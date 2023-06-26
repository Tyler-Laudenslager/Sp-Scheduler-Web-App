#! /bin/bash

echo "Destroy Database"
psql -U postgres -f sp_calendar_db_setup.sql -d sp_calendar
echo "Database Restored"
echo "Start Building App...."
go build -o sp_web_app
echo "Build Finished...."
