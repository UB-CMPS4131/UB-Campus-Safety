DROP TABLE IF EXISTS ROLES CASCADE;
DROP TABLE IF EXISTS PersonnelInfoTable CASCADE;
DROP TABLE IF EXISTS LOGIN CASCADE;
DROP TABLE IF EXISTS Report;

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
    DOB DATE,
    GENDER VARCHAR(10)
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
    file_path bytea
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

INSERT INTO PersonnelInfoTable (image, Fname, Mname, LName, DOB, GENDER) 
VALUES ('https://img.utdstc.com/icon/f47/819/f47819c8d0663a52a0d637b8d137169661d3033c6921d4811318731b8ed426b0:200', 
        'Alex', 
        'Humberto', 
        'Peraza', 
        '2003-03-12', 
        'Male');

INSERT INTO PersonnelInfoTable (image, Fname, Mname, LName, DOB, GENDER) 
VALUES ('https://cdn.iconscout.com/icon/free/png-256/free-avatar-370-456322.png?f=webp', 
        'Michael', 
        'David', 
        'Brown', 
        '1985-09-20', 
        'Male');


INSERT INTO LOGIN (memberID, username, password, Role, pastPasswords)
VALUES 
    (1, 'john_doe', '12345678', 3, ''),
    (2, '2020152022', '12345678', 2, ''),
    (3, 'mpit', '12345678', 1, '');
