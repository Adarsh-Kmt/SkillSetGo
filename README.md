# SkillSetGo

## Project Overview
SkillSetGo is a comprehensive platform connecting students with companies for job opportunities. The platform provides:
- Company registration and job posting
- Student profiles and job applications
- Job offer management

## Key Features
### For Companies:
- Register and manage company profile
- Post new job opportunities
- View and manage posted jobs
- Track student applications

### For Students:
- Create and update profile
- Browse available jobs based on eligibility
- Apply for jobs with one click
- Track application status
- Manage job offers

## Getting Started
### Prerequisites
- Go 1.20+
- PostgreSQL 14+
- Node.js 18+
- Python 3.8+
- Flask

### Installation
1. Clone the repository
2. Set up environment variables:
   ```
   JWT_PRIVATE_KEY=your_jwt_secret
   FLASK_SECRET_KEY=your_flask_secret
   ```
3. Run database migrations
4. Start the backend server:
   ```
   go run main.go
   ```
5. Start the frontend application:
   ```
   python app.py
   ```

## API Endpoints
### Authentication
- POST /student/login: Student login
- POST /company/login: Company login

### Job Management
- GET /job: Get all jobs (filtered by student eligibility)
- POST /job: Create a new job posting
- GET /company/jobs: Get company's posted jobs

### Student Operations
- POST /student/apply/{job-id}: Apply for a job
- GET /student/offer: Get student's job offers
- PUT /student/offer: Update job offer status

## Database Schema Overview
### Main Tables
- `company_table`: Stores company information
- `student_table`: Stores student information
- `job_table`: Stores job postings
- `student_job_application_table`: Tracks job applications
- `student_offer_table`: Manages job offers

## Features
### Job Application System
Students can apply for jobs through a simple one-click process:
- Automatic eligibility checking based on:
  - CGPA requirements
  - Branch/department matching
  - Batch year
- Real-time application status updates
- Prevention of duplicate applications
- Immediate feedback on application success/failure

## Error Handling
The API uses consistent error responses with:
- HTTP status code
- Error message
- Timestamp

## Deployment
### Production
- Use Docker containers
- Configure reverse proxy
- Enable HTTPS

### Monitoring
- Prometheus for metrics
- Grafana for visualization
- Sentry for error tracking

Create a pull request before merging anything to master.
