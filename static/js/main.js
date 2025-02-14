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

    // Handle job filters
    const jobFilters = document.querySelectorAll('.job-filter');
    jobFilters.forEach(filter => {
        filter.addEventListener('change', function() {
            // Get selected filters
            const selectedSalaryTiers = Array.from(document.querySelectorAll('input[name="salary-tier"]:checked')).map(cb => cb.value);
            
            // Filter jobs based on selected criteria
            const jobs = document.querySelectorAll('.job-card');
            jobs.forEach(job => {
                const salaryTier = job.dataset.salaryTier;
                const shouldShow = (selectedSalaryTiers.length === 0 || selectedSalaryTiers.includes(salaryTier));
                job.style.display = shouldShow ? '' : 'none';
            });
        });
    });

    // Add form submit handlers for job applications
    const jobApplicationForms = document.querySelectorAll('.job-application-form');
    jobApplicationForms.forEach(form => {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            try {
                const response = await fetch(form.action, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        job_id: form.querySelector('input[name="job_id"]').value
                    }),
                    credentials: 'include'
                });
                
                if (response.ok) {
                    showSuccess('Application submitted successfully!');
                    // Disable the apply button
                    form.querySelector('button[type="submit"]').disabled = true;
                } else {
                    const data = await response.json();
                    showError(data.error || 'Failed to submit application');
                }
            } catch (error) {
                showError('Failed to submit application');
            }
        });
    });
});
