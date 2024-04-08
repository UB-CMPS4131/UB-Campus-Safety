
# Deployment Guide

This guide will walk you through the deployment process for the UB Campus Safety Golang application with PostgreSQL as the database on an Ubuntu server.

## Prerequisites

- Ubuntu server with internet access
- Git installed on the server
- Golang installed on the server
- PostgreSQL installed and configured on the server
- GitHub repository access

## Steps

. **Build and Run the Application:**
   - Build the Golang application:
     ```bash
     go build ./cmd/web/
     ```
   - Run the application:
     ```bash
     go run ./cmd/web/
     ```
    - Run test cases:
     ```bash
     go test -v ./...
     ```
     - Run test cases with coverage:
     ```bash
     go test -v -cover ./...
     ```

2. **Access the Application:**
   - The application should now be running. Access it using a web browser at `http://localhost:8080`.

## Additional Notes

- Ensure that your firewall settings allow traffic on port 8080 if you're accessing the application remotely.

