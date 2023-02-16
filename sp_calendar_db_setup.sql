DROP TABLE SpUsers;
DROP TABLE SpManagers;
DROP TABLE Sessions;


CREATE TABLE SpUsers (
	id SERIAL PRIMARY KEY,
	name JSONB NOT NULL,
	username varchar(30),
	role INT4 NOT NULL,
	totalsessionsassigned INT4,
	sessionspool JSONB,
	sessionssorted JSONB,
	sessionsavailable JSONB,
	sessionsunavailable JSONB,
	sessionsselected JSONB,
	sessionsassigned JSONB,
	password varchar(100),
	email varchar(50) NOT NULL
);

CREATE TABLE SpManagers (
	id SERIAL PRIMARY KEY,
	name JSONB NOT NULL,
	username varchar(30) NOT NULL,
	role INT NOT NULL,
	password varchar(100),
	email varchar(50) NOT NULL,
	assignedpatients JSONB,
	sessionsmanaged JSONB,
	sessionsunmanaged JSONB
);

CREATE TABLE Sessions (
	id SERIAL PRIMARY KEY,
	title varchar(75) NOT NULL,
 	date varchar(50) NOT NULL,
	starttime varchar(50) NOT NULL,
	endtime varchar(50) NOT NULL,
	location varchar(50) NOT NULL,
	description varchar(300),
	status varchar(40),
	createddate varchar(20),
	expireddate varchar(20),
	CheckMarkAssigned BOOL NOT NULL,
	ShowSession BOOL NOT NULL,
	Instructors JSONB,
	PatientsNeeded int4 NOT NULL,
	PatientsAvailable JSONB,
	PatientsAssigned JSONB,
	PatientsSelected JSONB,
	PatientsUnavailable JSONB,
	PatientsNoResponse JSONB
);