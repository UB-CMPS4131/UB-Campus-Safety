# UB Campus Safety - An Integrated security System

The University of Belize grapples with a multifaceted challenge regarding the safety and security of its faculty, staff, and student body due to the absence of an integrated and user-friendly campus safety application. While the university employs certain safety precautions, they are not optimized to meet its community's dynamic and evolving safety needs. The existing measures lack the seamless ease of use, instantaneous information sharing, and personalized support that a comprehensive safety app could provide. Consequently, there is a growing demand among stakeholders for a centralized platform that enables swift responses to crises, disseminates pertinent safety information, and delivers timely notifications.


[![Build Status](https://github.com/mishal23/virtual-clinic/actions/workflows/django.yml/badge.svg)](https://github.com/UB-CMPS4131/UB-Campus-Safety/blob/main/.github/workflows/CI-Test.yml)
[![Coverage Status](https://img.shields.io/codecov/c/github/mishal23/virtual-clinic.svg)](https://codecov.io/gh/mishal23/virtual-clinic)


## Introduction

- Everything is well documented, please take a look at [docs](./docs) folder.
- All the required UML Diagrams are also drawn.
- Steps to setup the project are mentioned [here](./docs/INSTALLATION.md)
- Steps to deploy are mentioned [here](./docs/DEPLOY.md)

## Features:

- Common Login for all users

### Admin

- Add Doctor/Lab/Chemist
- Archive Users
- Restore Archived Users
- Add/Delete Speciality/Symptoms
- Add Hospitals
- View Activity
- View System Statistics
- View/Send Messages
- Update Profile
- Change Password

### Student

- Create Appointments
- Update Medical Information
- View Prescriptions
- View Medical Tests
- View/Send Messages
- Generate Invoice of Prescription
- Update Profile
- Change Password

### Guard

- Consult Appointments
- View/Update/Generate Prescriptions
- View Medical Information of patients
- Update Profile
- Change Password


## Structure of Repository

- All the documents are in `docs` folder.
- All the UML Diagrams are in `UML Diagrams` folder.
- In the virtualclinic folder
  - `public` folder contains all the templates.
  - `server` folder contains the views (business logic).
  - `testing` folder contains all the tests cases.
  - `virtualclinic` folder contains Django configuration files for the project.
