DROP TABLE IF EXISTS ROLES CASCADE;
DROP TABLE IF EXISTS PersonnelInfoTable CASCADE;
DROP TABLE IF EXISTS LOGIN CASCADE;
DROP TABLE IF EXISTS Report;
DROP TABLE IF EXISTS log;

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
    password VARCHAR(255) NOT NULL,
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

CREATE TABLE notification_seen (
    notification_id INTEGER REFERENCES notification(notification_id),
    user_id INTEGER REFERENCES personnelinfotable(id),
    seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (notification_id, user_id)
);

CREATE TABLE notification (
    notification_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES personnelinfotable(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO ROLES (role_name) VALUES ('admin'), ('student'), ('guard');

-- Inserting data with a URL for the image
INSERT INTO PersonnelInfoTable (image, Fname, Mname, LName, DOB, GENDER) 
VALUES ('https://w7.pngwing.com/pngs/340/946/png-transparent-avatar-user-computer-icons-software-developer-avatar-child-face-heroes-thumbnail.png', 
        'John', 
        'Alberto', 
        'Doe', 
        '1990-01-01', 
        'Male');



INSERT INTO LOGIN (memberID, username, password, Role, pastPasswords)
VALUES 
    (1, 'john_doe', '12345678', 1, '');
