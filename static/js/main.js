// Utility functions for API calls and token management
const API_URL = 'http://localhost:8080';

function getToken() {
    return localStorage.getItem('access_token');
}

function setToken(token) {
    if (token && token.startsWith('Bearer ')) {
        token = token.substring(7);
    }
    localStorage.setItem('access_token', token);
}

function clearToken() {
    localStorage.removeItem('access_token');
}

async function apiCall(endpoint, method = 'GET', data = null) {
    const headers = {
        'Content-Type': 'application/json'
    };
    
    const token = getToken();
    if (token) {
        headers['Auth'] = `Bearer ${token}`;
    }

    try {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method,
            headers,
            body: data ? JSON.stringify(data) : null,
            credentials: 'include' // Include cookies for session management
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'API call failed' }));
            throw new Error(errorData.error || 'API call failed');
        }

        return await response.json();
    } catch (error) {
        console.error('API call error:', error);
        showError(error.message);
        throw error;
    }
}

// UI Helper functions
function showError(message) {
    const errorDiv = document.getElementById('error-message');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
        setTimeout(() => {
            errorDiv.style.display = 'none';
        }, 5000);
    }
}

function showSuccess(message) {
    const successDiv = document.getElementById('success-message');
    if (successDiv) {
        successDiv.textContent = message;
        successDiv.style.display = 'block';
        setTimeout(() => {
            successDiv.style.display = 'none';
        }, 5000);
    }
}

// Student Dashboard Functions
function getSelectedFilters(filterType) {
    const filters = document.querySelectorAll(`input[name="${filterType}"]:checked`);
    return Array.from(filters).map(filter => filter.value);
}

async function handleJobApplication(event, jobId) {
    event.preventDefault();
    try {
        await apiCall(`/student/job/${jobId}/apply`, 'POST');
        showSuccess('Successfully applied for the job!');
    } catch (error) {
        console.error('Error applying for job:', error);
        showError('Failed to apply for the job');
    }
}

// Company Dashboard Functions
function loadCompanyJobs() {
    const jobsContainer = document.getElementById('company-jobs-container');
    if (!jobsContainer) return; // Not on company dashboard

    fetch('/company/job')
        .then(response => response.json())
        .then(data => {
            const jobs = data.jobs || [];
            jobsContainer.innerHTML = jobs.length ? jobs.map(job => `
                <div class="card mb-3 job-card">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start">
                            <div>
                                <h5 class="card-title">${job.job_role}</h5>
                                <h6 class="card-subtitle mb-2 text-muted">${job.company_name}</h6>
                            </div>
                            <span class="badge bg-primary">${job.salary_tier}</span>
                        </div>
                        <div class="row mt-3">
                            <div class="col-md-6">
                                <p><strong>CTC:</strong> ${job.ctc} LPA</p>
                                <p><strong>CGPA Cutoff:</strong> ${job.cgpa_cutoff}</p>
                            </div>
                            <div class="col-md-6">
                                <p><strong>Apply By:</strong> ${job.apply_by_date}</p>
                                <p><strong>Industry:</strong> ${job.industry}</p>
                            </div>
                        </div>
                        <div class="mt-3">
                            <button class="btn btn-primary" onclick="viewApplicants(${job.job_id})">
                                View Applicants
                            </button>
                        </div>
                    </div>
                </div>
            `).join('') : '<p class="text-center">No jobs posted yet</p>';
        })
        .catch(error => {
            console.error('Error loading jobs:', error);
            jobsContainer.innerHTML = '<div class="alert alert-danger">Failed to load jobs</div>';
        });
}

function viewApplicants(jobId) {
    const modal = new bootstrap.Modal(document.getElementById('applicationsModal'));
    const container = document.getElementById('applications-container');
    
    container.innerHTML = '<div class="text-center"><div class="spinner-border" role="status"></div></div>';
    modal.show();
    
    fetch(`/company/job/${jobId}/applicants`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Applicants data:', data); // Debug log
            const applicants = data.profiles || [];
            container.innerHTML = applicants.length ? applicants.map(applicant => `
                <div class="card mb-3">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start">
                            <div>
                                <h5 class="card-title">${applicant.name}</h5>
                                <h6 class="card-subtitle mb-2 text-muted">${applicant.usn}</h6>
                            </div>
                            <div>
                                <span class="badge bg-info">CGPA: ${applicant.cgpa}</span>
                            </div>
                        </div>
                        <div class="mt-3">
                            <button class="btn btn-sm btn-outline-primary me-2" onclick="viewProfile('${applicant.usn}')">
                                View Profile
                            </button>
                            <button class="btn btn-sm btn-success me-2" onclick="offerJob(${jobId}, '${applicant.usn}')">
                                Offer Job
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="rejectApplicant(${jobId}, '${applicant.usn}')">
                                Reject
                            </button>
                        </div>
                    </div>
                </div>
            `).join('') : '<p class="text-center">No applications yet</p>';
        })
        .catch(error => {
            console.error('Error loading applicants:', error);
            container.innerHTML = '<div class="alert alert-danger">Failed to load applicants</div>';
        });
}

function viewProfile(usn) {
    const modal = new bootstrap.Modal(document.getElementById('profileModal'));
    const container = document.getElementById('profile-container');
    
    container.innerHTML = '<div class="text-center"><div class="spinner-border" role="status"></div></div>';
    modal.show();
    
    fetch(`/student/${usn}/profile`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Profile data:', data); // Debug log
            const profile = data.profile || data;
            container.innerHTML = `
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">${profile.name || profile.student_name}</h5>
                        <div class="profile-details mt-3">
                            <div class="row mb-2">
                                <div class="col-4"><strong>USN:</strong></div>
                                <div class="col-8">${profile.usn}</div>
                            </div>
                            <div class="row mb-2">
                                <div class="col-4"><strong>Branch:</strong></div>
                                <div class="col-8">${profile.branch}</div>
                            </div>
                            <div class="row mb-2">
                                <div class="col-4"><strong>CGPA:</strong></div>
                                <div class="col-8">${profile.cgpa}</div>
                            </div>
                            <div class="row mb-2">
                                <div class="col-4"><strong>Batch:</strong></div>
                                <div class="col-8">${profile.batch}</div>
                            </div>
                            <div class="row mb-2">
                                <div class="col-4"><strong>Email:</strong></div>
                                <div class="col-8">${profile.email}</div>
                            </div>
                            <div class="row mb-2">
                                <div class="col-4"><strong>Backlogs:</strong></div>
                                <div class="col-8">${profile.num_active_backlogs || 0}</div>
                            </div>
                            ${profile.skills ? `
                            <div class="row mb-2">
                                <div class="col-4"><strong>Skills:</strong></div>
                                <div class="col-8">
                                    ${profile.skills.map(skill => `<span class="badge bg-secondary me-1">${skill}</span>`).join('')}
                                </div>
                            </div>
                            ` : ''}
                        </div>
                    </div>
                </div>
            `;
        })
        .catch(error => {
            console.error('Error loading profile:', error);
            container.innerHTML = '<div class="alert alert-danger">Failed to load profile. Please try again later.</div>';
        });
}

function offerJob(jobId, usn) {
    if (!confirm('Are you sure you want to offer this job to the student?')) return;
    
    fetch('/company/job/offer', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ job_id: jobId, usn: usn })
    })
    .then(response => response.json())
    .then(data => {
        alert('Job offered successfully!');
        viewApplicants(jobId); // Refresh applicants list
    })
    .catch(error => {
        console.error('Error offering job:', error);
        alert('Failed to offer job');
    });
}

function rejectApplicant(jobId, usn) {
    if (!confirm('Are you sure you want to reject this applicant?')) return;
    
    fetch('/company/job/reject', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ job_id: jobId, usn: usn })
    })
    .then(response => response.json())
    .then(data => {
        alert('Applicant rejected successfully!');
        viewApplicants(jobId); // Refresh applicants list
    })
    .catch(error => {
        console.error('Error rejecting applicant:', error);
        alert('Failed to reject applicant');
    });
}

// Authentication Functions
async function handleLogin(event, type = 'student') {
    event.preventDefault();
    const form = event.target;
    
    try {
        const loginData = type === 'student' ? {
            usn: form.usn.value,
            password: form.password.value
        } : {
            username: form.username.value,
            password: form.password.value
        };

        const endpoint = type === 'student' ? '/student/login' : '/company/login';
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData),
            credentials: 'include'
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ error: 'Login failed' }));
            throw new Error(errorData.error || 'Login failed');
        }

        const data = await response.json();
        setToken(data.access_token);
        
        // Store the token in the session
        await fetch(window.location.origin + endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });

        // Redirect based on login type
        window.location.href = type === 'student' ? '/dashboard' : '/company/dashboard';
    } catch (error) {
        console.error('Login error:', error);
        showError('Login failed. Please check your credentials.');
    }
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function() {
    // Add form submit handlers
    const studentLoginForm = document.getElementById('student-login-form');
    if (studentLoginForm) {
        studentLoginForm.addEventListener('submit', (e) => handleLogin(e, 'student'));
    }

    const companyLoginForm = document.getElementById('company-login-form');
    if (companyLoginForm) {
        companyLoginForm.addEventListener('submit', (e) => handleLogin(e, 'company'));
    }

    // Event listeners for job filters
    const jobFilters = document.querySelectorAll('.btn-check');
    jobFilters.forEach(filter => {
        filter.addEventListener('change', function() {
            document.querySelector('form').submit();
        });
    });

    // Initialize company dashboard if on that page
    if (document.getElementById('company-jobs-container')) {
        loadCompanyJobs();
    }
});
