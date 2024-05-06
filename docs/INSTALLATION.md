# Installation Guide

This guide will walk you through the installation process for setting up a Golang project running on Ubuntu and using PostgreSQL as the database. The project is hosted on GitHub at [UB-CMPS4131/UB-Campus-Safety](https://github.com/UB-CMPS4131/UB-Campus-Safety).

## Prerequisites

Before starting the installation process, ensure you have the following prerequisites installed on your system:

- Ubuntu operating system (preferably the latest version)
- Golang installed on your system
- PostgreSQL installed and configured

## Step 1: Clone the Repository

```bash
git clone https://github.com/UB-CMPS4131/UB-Campus-Safety.git
```

## Step 2: Set Up PostgreSQL

1. Install PostgreSQL on your Ubuntu system if not already installed:
   
   ```bash
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   ```

2. Log into the PostgreSQL interactive terminal as the default user, postgres:
   
   ```bash
   sudo -u postgres psql
   ```

3. Create a new database for the project:
   
   ```sql
   CREATE DATABASE ub_campus_safety;
   ```

4. Create a new user and grant them privileges on the database:
   
   ```sql
   CREATE USER ub_campus_safety WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE ub_campus_safety TO ub_campus_safety;
   ```

5. Exit the PostgreSQL interactive terminal:
   
   ```sql
   \q
   ```

## Step 3: Configure Environment Variables

1. Navigate to the project directory:
   
   ```bash
   cd UB-Campus-Safety
   ```

2. Edit the `main.go` file located inside the cmd/web folder and add the following environment variables:
   
   ``` main.go
   const (
		host     = "localhost"
		port     = 5432
		user     = "ub_campus_safety"
		password = "your_password"
		dbname   = "ub_campus_safety"
	)
   ```

   Replace `your_password` with the password you set for the PostgreSQL user.

## Step 4: Configure Database 

1. In the root directory of the project, Login into your database :
   
   ```bash
   PGPASSWORD=your_password psql -h localhost -U ub_campus_safety ub_campus_safety
   ```
   Replace `your_password` with the password you set for the PostgreSQL user.

2. Run the sql file inside of the project:
   
   ```bash
   \i sql.sql
   ```
