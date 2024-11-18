
-- Create ROLES table
CREATE TABLE ROLES (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL
);

-- Create PersonnelInfoTable
CREATE TABLE PersonnelInfoTable (
    id SERIAL PRIMARY KEY,
    image VARCHAR(500), -- Assuming VARCHAR data type for image URLs, adjust size as needed
    Fname VARCHAR(50) NOT NULL,
    Mname VARCHAR(50),
    LName VARCHAR(50) NOT NULL,
    username VARCHAR(50),
    DOB DATE,
    GENDER VARCHAR(10),
    imagedata BYTEA -- Fixed syntax error
);

-- Create LOGIN table
CREATE TABLE LOGIN (
    id SERIAL PRIMARY KEY,
    memberID INTEGER REFERENCES PersonnelInfoTable(id),
    username VARCHAR(50) NOT NULL,
    password bytea NOT NULL,
    role INTEGER REFERENCES ROLES(role_id),
    pastPasswords TEXT -- Assuming pastPasswords will store multiple passwords separated by a delimiter
);


CREATE TABLE Report (
    report_id SERIAL PRIMARY KEY,
    type_of_incident TEXT NOT NULL,
    person_name VARCHAR(150), -- Combined field for first name and last name
    location TEXT NOT NULL,
    description TEXT NOT NULL,
    device_location TEXT,
    imagename VARCHAR(50),
    imagedata BYTEA,
    encodedimagedata TEXT,
    mime_type VARCHAR(50)
);

CREATE TABLE log (
    id SERIAL PRIMARY KEY,
    person_name VARCHAR(150),
    log_date DATE NOT NULL,
    log_time TIME NOT NULL,
    check_type VARCHAR(8) NOT NULL
);
CREATE TABLE notification (
    notification_id SERIAL PRIMARY KEY,
    user_id integer,
    title VARCHAR(100) NOT NULL,
    message text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE contact (
    contact_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    number text NOT NULL,
    email text NOT NULL
);

CREATE TABLE emergency (
    emergency_id SERIAL PRIMARY KEY,
    person_name VARCHAR(150),
    location text NOT NULL,
    message text NOT NULL
);

CREATE TABLE studentcontact (
    contact_id SERIAL PRIMARY KEY,
    LoginID INTEGER REFERENCES LOGIN(id),
    name VARCHAR(100) NOT NULL,
    number text NOT NULL,
    email text NOT NULL
);
CREATE TABLE notification_seen (
    notification_id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    seen_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO ROLES (role_name) VALUES ('admin'), ('student'), ('guard');

-- Inserting data with a URL for the image
INSERT INTO PersonnelInfoTable (image, Fname, Mname, LName, DOB, GENDER) 
VALUES ('', 'abner','adonay','mencia','2003-12-17','Male');

INSERT INTO LOGIN (memberID, username, password, Role, pastPasswords)
VALUES (1, 'abner', '\x243261243132246d75577146447968344a4c76764e495643412f46584f435542456a5868334d7665324e677279304749456631556531303653543543', 1, '');
