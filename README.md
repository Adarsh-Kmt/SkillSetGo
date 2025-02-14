# SkillSetGo: Student Job Platform

## 🚀 Project Overview

SkillSetGo is a comprehensive student job application platform designed to bridge the gap between students and potential employers. The platform offers a seamless experience for job seekers and companies, providing robust features for job application, tracking, and management.

## 🛠 Tech Stack

### Backend
- **Language**: Go (Golang)
- **Web Framework**: Gorilla Mux
- **ORM**: SQLC
- **Authentication**: JWT (JSON Web Tokens)

### Frontend
- **Language**: Python
- **Web Framework**: Flask
- **Frontend**: Bootstrap, Vanilla JavaScript
- **Template Engine**: Jinja2

### Database
- **Database**: PostgreSQL
- **Hosting**: Local development (localhost:8087)

### Additional Tools
- **Resume Parsing**: Groq AI-powered Resume Matcher
- **PDF Parsing**: PyPDF2

## 🌟 Key Features

### Student Features
- Job Discovery
- Job Application Tracking
- Profile Management
- Resume Parsing and Scoring

### Company Features
- Job Posting
- Applicant Management
- Offer Management
- Candidate Screening

## 🔍 Innovative Components

### Resume Parser
The Resume Parser is an AI-powered tool that helps students understand how well their resume matches a specific job description.

#### Key Functionalities
- PDF Resume Upload
- Job Description Input
- AI-Powered Matching
- Detailed Scoring Mechanism

#### How It Works
1. Upload a PDF resume
2. Paste a job description
3. Get an AI-generated compatibility score
4. Receive insights on resume strengths and weaknesses

### Authentication System
- Secure JWT-based authentication
- Separate login flows for students and companies
- Session management
- Role-based access control

## 📦 Project Structure

```
SkillSetGo/
│
├── pkg/                # Go backend packages
│   ├── handler/        # API request handlers
│   ├── db/             # Database queries and models
│   └── middleware/     # Authentication middleware
│
├── templates/          # HTML templates
│   ├── index.html
│   ├── dashboard.html
│   └── company_dashboard.html
│
├── static/             # Static assets
│   ├── css/
│   └── js/
│
├── innovative component.py  # Resume parsing tool
├── app.py              # Flask application
└── requirements.txt    # Python dependencies
```

## 🚦 Getting Started

### Prerequisites
- Python 3.9+
- Go 1.20+
- PostgreSQL
- Groq API Key (for Resume Parser)

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/SkillSetGo.git
cd SkillSetGo
```

2. Install Python dependencies
```bash
pip install -r requirements.txt
```

3. Set up PostgreSQL database
4. Configure environment variables
5. Run the application
```bash
python app.py
```

## 🔐 Environment Configuration

Create a `.env` file with the following variables:
```
DATABASE_URL=postgresql://username:password@localhost:8087/skillsetgo
JWT_SECRET=your_secret_key
GROQ_API_KEY=your_groq_api_key
```

## 🧪 Testing

- Unit Tests: Located in `tests/` directory
- Integration Tests: Covers API endpoints and database interactions

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🙌 Acknowledgements

- Groq AI for Resume Parsing
- Bootstrap for Frontend Design
- Flask and Go Communities

## 📞 Support

For issues or questions, please open a GitHub issue or contact support@skillsetgo.com
