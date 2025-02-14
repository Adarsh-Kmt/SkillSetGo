from flask import Flask, render_template, request, jsonify, redirect, url_for, session, flash
from dotenv import load_dotenv
import os
import requests
from datetime import datetime
import base64
import json

load_dotenv()

app = Flask(__name__)
app.config['SECRET_KEY'] = os.getenv('FLASK_SECRET_KEY', 'dev-key-123')

# Go backend API URL
API_URL = 'http://localhost:8080'  # Adjust this to your Go server port

def get_user_id_from_token(token):
    """Extract user_id from JWT token"""
    try:
        # Split the token into parts
        parts = token.split('.')
        if len(parts) != 3:
            return None
        
        # Decode the payload (middle part)
        payload = parts[1]
        payload += '=' * ((4 - len(payload) % 4) % 4)
        
        # Decode the payload
        decoded = base64.b64decode(payload)
        payload_data = json.loads(decoded)
        
        # Extract user_id from claims
        return payload_data.get('sub')  # JWT typically uses 'sub' for subject/user_id
    except Exception as e:
        print(f"Error extracting user_id from token: {str(e)}")
        return None

def get_auth_header():
    if 'access_token' in session:
        return {'Auth': session['access_token']}
    elif 'company_access_token' in session:
        return {'Auth': session['company_access_token']}
    return None

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/login', methods=['GET', 'POST'])
def login():
    try:
        print("\n=== Login Route ===")
        print(f"Request Method: {request.method}")
        print(f"Form Data: {request.form}")
        print(f"Headers: {dict(request.headers)}")
        print(f"Current Session: {session}")
        
        if request.method == 'POST':
            try:
                # Get form data
                usn = request.form.get('usn', '').strip().upper()  # Convert to uppercase
                password = request.form.get('password', '').strip()
                
                print(f"USN: {usn}")
                print(f"Password length: {len(password) if password else 0}")
                
                if not usn or not password:
                    flash('Please enter both USN and password', 'error')
                    return render_template('login.html')
                
                # Prepare login data
                login_data = {
                    'usn': usn,
                    'password': password
                }
                
                # Make login request
                print(f"Making request to {API_URL}/student/login")
                print(f"Request data: {login_data}")
                
                try:
                    print("Sending request...")
                    response = requests.post(
                        f"{API_URL}/student/login",
                        json=login_data,
                        headers={'Content-Type': 'application/json'},
                        timeout=10
                    )
                    print("Request sent successfully")
                    print(f"Response status: {response.status_code}")
                    print(f"Response headers: {dict(response.headers)}")
                    print(f"Response body: {response.text}")
                    
                    if response.status_code == 200:
                        try:
                            response_data = response.json()
                            print(f"Response data: {response_data}")
                            token = response_data.get('access_token')
                            
                            if not token:
                                print("No access_token in response")
                                flash('Invalid server response', 'error')
                                return render_template('login.html')
                            
                            # Store token in session
                            session.clear()
                            session['access_token'] = token
                            session['user_type'] = 'student'
                            
                            print(f"Token stored in session: {token[:10]}...")
                            print(f"Updated session: {session}")
                            
                            return redirect(url_for('dashboard'))
                        except json.JSONDecodeError as e:
                            print(f"Failed to parse response JSON: {str(e)}")
                            print(f"Raw response text: {response.text}")
                            flash('Invalid server response format', 'error')
                            return render_template('login.html')
                    else:
                        error_message = "Invalid credentials"
                        try:
                            error_data = response.json()
                            if isinstance(error_data, dict) and 'error' in error_data:
                                error_message = error_data['error']
                        except Exception as e:
                            print(f"Error parsing error response: {str(e)}")
                            print(f"Raw error response: {response.text}")
                        
                        print(f"Login failed: {error_message}")
                        flash(error_message, 'error')
                        return render_template('login.html')
                except requests.RequestException as e:
                    print(f"Request failed: {str(e)}")
                    print(f"Request exception type: {type(e)}")
                    flash('Failed to connect to server', 'error')
                    return render_template('login.html')
            except Exception as e:
                print(f"Error in POST handling: {str(e)}")
                print(f"Exception type: {type(e)}")
                import traceback
                print(f"Traceback: {traceback.format_exc()}")
                flash('Error processing login request', 'error')
                return render_template('login.html')
    except Exception as e:
        import traceback
        print(f"Exception in login route: {str(e)}")
        print(f"Exception type: {type(e)}")
        print(f"Traceback: {traceback.format_exc()}")
        flash('An unexpected error occurred', 'error')
        return render_template('login.html')
    
    return render_template('login.html')

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        try:
            # Get all required fields
            name = request.form.get('name')
            usn = request.form.get('usn')
            password = request.form.get('password')
            branch = request.form.get('branch')
            cgpa = request.form.get('cgpa')
            email = request.form.get('email')
            batch = request.form.get('batch')
            counsellor_email_id = request.form.get('counsellor_email_id')
            num_of_backlogs = request.form.get('num_of_backlogs')

            print(f"Received registration data: {request.form}")  # Debug print

            # Validate required fields
            if not all([name, usn, password, branch, cgpa, email, batch, counsellor_email_id]):
                missing_fields = [field for field, value in {
                    'name': name, 'usn': usn, 'password': password,
                    'branch': branch, 'cgpa': cgpa, 'email': email,
                    'batch': batch, 'counsellor_email_id': counsellor_email_id
                }.items() if not value]
                return render_template('register.html', 
                    error=f"Missing required fields: {', '.join(missing_fields)}")

            # Convert numeric fields with validation
            try:
                cgpa = float(cgpa)
                if not 0 <= cgpa <= 10:
                    return render_template('register.html', 
                        error="CGPA must be between 0 and 10")
            except ValueError:
                return render_template('register.html', 
                    error="Invalid CGPA value")

            try:
                batch = int(batch)
                if batch < 2026:
                    return render_template('register.html', 
                        error="Batch must be 2026 or later")
            except ValueError:
                return render_template('register.html', 
                    error="Invalid batch value")

            try:
                num_of_backlogs = int(num_of_backlogs or 0)
                if num_of_backlogs < 0:
                    return render_template('register.html', 
                        error="Number of backlogs cannot be negative")
            except ValueError:
                return render_template('register.html', 
                    error="Invalid number of backlogs")

            # Validate USN format
            if not usn.startswith('1RV'):
                return render_template('register.html', 
                    error="USN must start with 1RV")

            # Validate email format
            if not email.endswith('@rvce.edu.in'):
                return render_template('register.html', 
                    error="Email must end with @rvce.edu.in")

            # Create request data - ensure all numeric values are properly typed
            data = {
                'name': str(name),
                'usn': str(usn).upper(),  # Ensure USN is uppercase
                'password': str(password),
                'branch': str(branch).upper(),  # Ensure branch is uppercase
                'cgpa': float(cgpa),  # Ensure it's a float
                'email': str(email).lower(),  # Ensure email is lowercase
                'batch': int(batch),  # Ensure it's an integer
                'counsellor_email_id': str(counsellor_email_id).lower(),  # Ensure email is lowercase
                'num_of_backlogs': int(num_of_backlogs)  # Ensure it's an integer
            }

            print(f"Sending registration data to API: {data}")  # Debug print

            # Make API request with proper headers
            try:
                headers = {'Content-Type': 'application/json'}
                response = requests.post(
                    f"{API_URL}/student/register", 
                    json=data,  # Use json parameter to automatically handle JSON encoding
                    headers=headers
                )
                print(f"API Response Status: {response.status_code}")  # Debug print
                print(f"API Response Headers: {response.headers}")  # Debug print
                print(f"API Response Body: {response.text}")  # Debug print
                
                if response.status_code == 200:
                    return redirect(url_for('login'))
                else:
                    try:
                        error_data = response.json()
                        error_msg = error_data.get('error', 'Registration failed')
                        if 'internal server error' in error_msg.lower():
                            error_msg = "Server error: Please try again later or contact support"
                    except ValueError:
                        error_msg = f"Registration failed: {response.text}"
                    return render_template('register.html', error=error_msg)
            except requests.exceptions.RequestException as e:
                print(f"API Request Error: {str(e)}")  # Debug print
                return render_template('register.html', 
                    error=f"Failed to connect to server: {str(e)}")

        except Exception as e:
            print(f"Unexpected error: {str(e)}")  # Debug print
            return render_template('register.html', 
                error=f"An unexpected error occurred: {str(e)}")

    return render_template('register.html')

@app.route('/dashboard')
def dashboard():
    if 'access_token' not in session:
        return redirect(url_for('login'))

    headers = get_auth_header()
    if not headers:
        return redirect(url_for('login'))

    try:
        # Get student ID from token
        token = session['access_token']
        token_parts = token.split('.')
        if len(token_parts) > 1:
            import base64
            import json
            payload = json.loads(base64.b64decode(token_parts[1] + '=' * (-len(token_parts[1]) % 4)).decode('utf-8'))
            student_id = payload.get('id')  # Using dict.get() instead of attribute access
            
            if not student_id:
                print("No student ID found in token payload:", payload)
                flash('Invalid session', 'error')
                return redirect(url_for('login'))
        else:
            flash('Invalid token format', 'error')
            return redirect(url_for('login'))

        print(f"Student ID from token: {student_id}")

        # Get student profile
        print("Fetching profile...")
        profile_response = requests.get(f"{API_URL}/student/{student_id}/profile", headers=headers)
        print(f"Profile response status: {profile_response.status_code}")
        print(f"Profile response: {profile_response.text}")
        if profile_response.status_code == 200:
            profile = profile_response.json().get('profile', {})
        else:
            print(f"Error fetching profile: {profile_response.text}")
            profile = {}

        # Get applied jobs
        print("Fetching applied jobs...")
        applied_response = requests.get(f"{API_URL}/student/job/apply", headers=headers)
        print(f"Applied jobs response status: {applied_response.status_code}")
        print(f"Applied jobs response: {applied_response.text}")
        if applied_response.status_code == 200:
            applied_jobs = applied_response.json().get('jobs', [])
        else:
            print(f"Error fetching applied jobs: {applied_response.text}")
            applied_jobs = []

        # Get job offers
        print("Fetching job offers...")
        offers_response = requests.get(f"{API_URL}/student/job/offer", headers=headers)
        print(f"Offers response status: {offers_response.status_code}")
        print(f"Offers response: {offers_response.text}")
        if offers_response.status_code == 200:
            job_offers = offers_response.json().get('offers', [])
        else:
            print(f"Error fetching job offers: {offers_response.text}")
            job_offers = []

        # Get available jobs
        print("Fetching available jobs...")
        jobs_response = requests.get(f"{API_URL}/student/job", headers=headers)
        print(f"Jobs response status: {jobs_response.status_code}")
        print(f"Jobs response: {jobs_response.text}")
        if jobs_response.status_code == 200:
            jobs = jobs_response.json().get('jobs', [])
        else:
            print(f"Error fetching jobs: {jobs_response.text}")
            jobs = []

        print("Rendering template with data:")
        print(f"Profile: {profile}")
        print(f"Applied jobs: {applied_jobs}")
        print(f"Job offers: {job_offers}")
        print(f"Jobs: {jobs}")

        return render_template('dashboard.html', 
                             profile=profile,
                             applied_jobs=applied_jobs,
                             job_offers=job_offers,
                             is_placed=any(job.get('status') == 'Accepted' for job in applied_jobs),
                             jobs=jobs)

    except Exception as e:
        print(f"Error in student dashboard: {str(e)}")
        print(f"Token payload: {payload if 'payload' in locals() else 'Not available'}")
        import traceback
        print(f"Traceback: {traceback.format_exc()}")
        flash('Error loading dashboard data', 'error')
        return render_template('dashboard.html', 
                             profile={},
                             applied_jobs=[],
                             job_offers=[],
                             is_placed=False,
                             jobs=[])

@app.route('/apply/<int:job_id>', methods=['POST'])
def apply_job(job_id):
    try:
        if 'access_token' not in session:
            flash('Please log in first', 'error')
            return redirect(url_for('login'))
        
        headers = get_auth_header()
        if not headers:
            session.clear()
            flash('Session expired, please log in again', 'error')
            return redirect(url_for('login'))
        
        # Make request to apply for job
        response = requests.post(
            f"{API_URL}/student/job/{job_id}/apply",
            headers=headers
        )
        
        if response.status_code == 200:
            flash('Successfully applied for job!', 'success')
        elif response.status_code == 401:
            session.clear()
            flash('Session expired, please log in again', 'error')
            return redirect(url_for('login'))
        else:
            error_message = "Failed to apply for job"
            try:
                error_data = response.json()
                if isinstance(error_data, dict) and 'error' in error_data:
                    error_message = error_data['error']
            except:
                pass
            flash(error_message, 'error')
            
        return redirect(url_for('dashboard'))
    except Exception as e:
        print(f"Error applying for job: {str(e)}")
        flash('An error occurred while applying for the job', 'error')
        return redirect(url_for('dashboard'))

@app.route('/offer/action', methods=['POST'])
def offer_action():
    if 'access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    data = request.json
    headers = get_auth_header()
    try:
        response = requests.put(f"{API_URL}/student/job/offer", headers=headers, json=data)
        return jsonify({'success': response.status_code == 200}), response.status_code
    except requests.exceptions.RequestException:
        return jsonify({'success': False, 'error': 'Server error'}), 500

@app.route('/company/dashboard')
def company_dashboard():
    if 'company_access_token' not in session:
        return redirect(url_for('company_login'))

    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        # Get published jobs
        print("Fetching published jobs...")
        jobs_response = requests.get(f"{API_URL}/company/job", headers=headers)
        print(f"Jobs response: {jobs_response.status_code} - {jobs_response.text}")
        
        jobs = []
        if jobs_response.status_code == 200:
            jobs = jobs_response.json().get('jobs', [])
            
            # Get application counts for each job
            for job in jobs:
                app_response = requests.get(f"{API_URL}/company/job/{job['job_id']}/application", headers=headers)
                if app_response.status_code == 200:
                    applications = app_response.json().get('applications', [])
                    job['application_count'] = len(applications)
                else:
                    job['application_count'] = 0

        # Get company name from token
        token = session['company_access_token']
        token_parts = token.split('.')
        if len(token_parts) > 1:
            import base64
            import json
            payload = json.loads(base64.b64decode(token_parts[1] + '=' * (-len(token_parts[1]) % 4)).decode('utf-8'))
            company_name = payload.get('company_name', 'Company')
        else:
            company_name = 'Company'

        return render_template('company_dashboard.html', 
                             company_name=company_name,
                             jobs=jobs)

    except Exception as e:
        print(f"Error in company dashboard: {str(e)}")
        import traceback
        print(f"Traceback: {traceback.format_exc()}")
        flash('Error loading dashboard data', 'error')
        return redirect(url_for('company_login'))

@app.route('/company/login', methods=['GET', 'POST'])
def company_login():
    if request.method == 'POST':
        data = {
            'username': request.form['username'],
            'password': request.form['password']
        }
        
        try:
            response = requests.post(f"{API_URL}/company/login", json=data)
            print(f"Company login response: {response.status_code} - {response.text}")
            
            if response.status_code == 200:
                token_data = response.json()
                session.clear()  # Clear any existing session
                session['company_access_token'] = token_data.get('access_token')
                session['user_type'] = 'company'
                return redirect(url_for('company_dashboard'))
            else:
                error_msg = response.json().get('error', 'Login failed')
                flash(error_msg, 'error')
                
        except Exception as e:
            print(f"Error during company login: {str(e)}")
            flash('An error occurred during login', 'error')
            
    return render_template('company_login.html')

@app.route('/company/register', methods=['GET', 'POST'])
def company_register():
    if request.method == 'POST':
        try:
            data = {
                'company_name': request.form.get('company_name'),
                'industry': request.form.get('industry'),
                'poc_name': request.form.get('poc_name'),
                'poc_phno': request.form.get('poc_phno'),
                'username': request.form.get('username'),
                'password': request.form.get('password')
            }
            
            print(f"Sending company registration data: {data}")  # Debug print
            response = requests.post(f"{API_URL}/company/register", json=data)
            print(f"Registration response: {response.status_code} - {response.text}")  # Debug print
            
            if response.status_code == 200:
                return redirect(url_for('company_login'))
            else:
                error_msg = response.json().get('error', 'Registration failed')
                return render_template('company_register.html', error=error_msg)
        except requests.exceptions.RequestException as e:
            print(f"Company registration error: {str(e)}")  # Debug print
            return render_template('company_register.html', error="Server error")
    return render_template('company_register.html')

@app.route('/student/<string:usn>/profile')
def get_student_profile(usn):
    try:
        if 'access_token' not in session and 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        # First get student ID using USN from the database
        response = requests.get(f"{API_URL}/student/id/{usn}", headers=headers)
        print(f"Student ID response for {usn}:", response.status_code, response.text)
        
        if response.status_code != 200:
            return jsonify({'error': 'Student not found'}), response.status_code
            
        student_data = response.json()
        student_id = student_data.get('student_id')
        
        if not student_id:
            return jsonify({'error': 'Student ID not found'}), 400
            
        # Now get the profile using student ID
        profile_response = requests.get(f"{API_URL}/student/{student_id}/profile", headers=headers)
        print(f"Profile response for student {student_id}:", profile_response.status_code, profile_response.text)
        
        if profile_response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if profile_response.status_code != 200:
            return jsonify({'error': 'Failed to load profile'}), profile_response.status_code
            
        return profile_response.json()
        
    except Exception as e:
        print(f"Error getting student profile: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/job/<int:job_id>/applications')
def get_job_applications(job_id):
    if 'company_access_token' not in session:
        return jsonify({'error': 'Please log in first'}), 401
        
    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        # Get applications for this job
        print(f"Fetching applications for job {job_id}...")
        response = requests.get(f"{API_URL}/company/job/{job_id}/application", headers=headers)
        print(f"Applications response: {response.status_code} - {response.text}")
        
        if response.status_code == 200:
            return response.json()
        else:
            return jsonify({'error': 'Failed to fetch applications'}), response.status_code
            
    except Exception as e:
        print(f"Error fetching applications: {str(e)}")
        return jsonify({'error': 'An error occurred'}), 500

@app.route('/company/job/<int:job_id>/offer/<string:usn>', methods=['POST'])
def offer_job(job_id, usn):
    if 'company_access_token' not in session:
        return jsonify({'error': 'Please log in first'}), 401
        
    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        response = requests.post(
            f"{API_URL}/company/job/{job_id}/application/{usn}/offer",
            headers=headers
        )
        
        if response.status_code == 200:
            return jsonify({'message': 'Job offered successfully'})
        else:
            return jsonify({'error': response.json().get('error', 'Failed to offer job')}), response.status_code
            
    except Exception as e:
        print(f"Error offering job: {str(e)}")
        return jsonify({'error': 'An error occurred'}), 500

@app.route('/company/job/<int:job_id>/reject/<string:usn>', methods=['POST'])
def reject_application(job_id, usn):
    if 'company_access_token' not in session:
        return jsonify({'error': 'Please log in first'}), 401
        
    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        response = requests.post(
            f"{API_URL}/company/job/{job_id}/application/{usn}/reject",
            headers=headers
        )
        
        if response.status_code == 200:
            return jsonify({'message': 'Application rejected successfully'})
        else:
            return jsonify({'error': response.json().get('error', 'Failed to reject application')}), response.status_code
            
    except Exception as e:
        print(f"Error rejecting application: {str(e)}")
        return jsonify({'error': 'An error occurred'}), 500

@app.route('/company/post-job', methods=['POST'])
def post_job():
    if 'company_access_token' not in session:
        return redirect(url_for('company_login'))
        
    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        # Get form data
        job_data = {
            'job_role': request.form['job_role'],
            'job_type': request.form['job_type'],
            'ctc': float(request.form['ctc']),
            'salary_tier': request.form['salary_tier'],
            'description': request.form['description'],
            'cgpa_cutoff': float(request.form['cgpa_cutoff']),
            'apply_by_date': request.form['apply_by_date']
        }
        
        print(f"Posting job with data: {job_data}")
        response = requests.post(f"{API_URL}/company/job", headers=headers, json=job_data)
        print(f"Post job response: {response.status_code} - {response.text}")
        
        if response.status_code == 200:
            flash('Job posted successfully!', 'success')
        else:
            error_msg = response.json().get('error', 'Failed to post job')
            flash(f'Failed to post job: {error_msg}', 'error')
            
    except Exception as e:
        print(f"Error posting job: {str(e)}")
        flash('An error occurred while posting the job', 'error')
        
    return redirect(url_for('company_dashboard'))

@app.route('/profile/update', methods=['POST'])
def update_profile():
    if 'access_token' not in session:
        return redirect(url_for('login'))
    
    if 'user_id' not in session:
        return redirect(url_for('login'))
    
    headers = get_auth_header()
    try:
        # Get form data
        profile_data = {
            'name': request.form.get('name'),
            'email': request.form.get('email'),
            'phone': request.form.get('phone'),
            'cgpa': float(request.form.get('cgpa')),
            'branch': request.form.get('branch'),
            'graduation_year': int(request.form.get('graduation_year')),
            'skills': request.form.get('skills'),
            'resume_link': request.form.get('resume_link')
        }
        
        # Update profile
        response = requests.put(
            f"{API_URL}/student/{session['user_id']}/profile",
            headers=headers,
            json=profile_data
        )
        
        if response.status_code == 200:
            flash('Profile updated successfully!', 'success')
        elif response.status_code == 401:
            flash('Please login again', 'error')
            session.clear()
            return redirect(url_for('login'))
        else:
            error_data = response.json()
            flash(error_data.get('error', 'Failed to update profile'), 'error')
        
        return redirect(url_for('dashboard'))
    except requests.exceptions.RequestException as e:
        print(f"Error updating profile: {str(e)}")  # Debug print
        flash('Failed to connect to server', 'error')
        return redirect(url_for('dashboard'))

@app.route('/logout')
def logout():
    session.pop('access_token', None)
    session.pop('company_access_token', None)
    session.pop('user_id', None)
    return redirect(url_for('index'))

@app.route('/company/job')
def get_company_jobs():
    try:
        if 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        response = requests.get(f"{API_URL}/company/job", headers=headers)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to load jobs'}), response.status_code
            
        return response.json()
        
    except Exception as e:
        print(f"Error getting company jobs: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/student/job-offers')
def get_student_job_offers():
    try:
        if 'access_token' not in session:
            print("No access token found")
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            print("No auth headers")
            return jsonify({'error': 'Session expired'}), 401

        print("Fetching job offers from API...")
        response = requests.get(f"{API_URL}/student/offers", headers=headers)  # Updated endpoint
        print(f"API Response Status: {response.status_code}")
        print(f"API Response Body: {response.text}")
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to load job offers'}), response.status_code
            
        offers_data = response.json()
        print(f"Processed offers data: {offers_data}")
        return jsonify(offers_data)
        
    except Exception as e:
        print(f"Error getting job offers: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/student/accept-offer/<int:job_id>', methods=['POST'])
def accept_offer(job_id):
    if 'access_token' not in session:
        return redirect(url_for('login'))
        
    headers = get_auth_header()
    if not headers:
        return redirect(url_for('login'))
        
    try:
        # Send PUT request to update offer status with the correct action value
        response = requests.put(
            f"{API_URL}/student/job/offer",
            headers=headers,
            json={
                "job_id": job_id,
                "action": "ACCEPT"  # Changed from ACCEPTED to ACCEPT
            }
        )
        
        print(f"Accept offer response status: {response.status_code}")
        print(f"Accept offer response: {response.text}")
        
        if response.status_code == 200:
            flash('🎉 Congratulations! You have accepted the job offer!', 'success')
        else:
            error_msg = response.json().get('error', 'Failed to accept offer')
            print(f"Error accepting offer: {error_msg}")
            flash(f'Failed to accept offer: {error_msg}', 'error')
            
    except Exception as e:
        print(f"Error accepting offer: {str(e)}")
        flash('An error occurred while accepting the offer', 'error')
        
    return redirect(url_for('dashboard'))

@app.route('/student/reject-offer/<int:job_id>', methods=['POST'])
def reject_offer(job_id):
    if 'access_token' not in session:
        return redirect(url_for('login'))
        
    headers = get_auth_header()
    if not headers:
        return redirect(url_for('login'))
        
    try:
        # Send PUT request to update offer status with the correct action value
        response = requests.put(
            f"{API_URL}/student/job/offer",
            headers=headers,
            json={
                "job_id": job_id,
                "action": "REJECT"  # Changed from REJECTED to REJECT
            }
        )
        
        print(f"Reject offer response status: {response.status_code}")
        print(f"Reject offer response: {response.text}")
        
        if response.status_code == 200:
            flash('You have rejected the job offer', 'info')
        else:
            error_msg = response.json().get('error', 'Failed to reject offer')
            print(f"Error rejecting offer: {error_msg}")
            flash(f'Failed to reject offer: {error_msg}', 'error')
            
    except Exception as e:
        print(f"Error rejecting offer: {str(e)}")
        flash('An error occurred while rejecting the offer', 'error')
        
    return redirect(url_for('dashboard'))

@app.route('/student/placement-status')
def get_placement_status():
    try:
        if 'access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401

        response = requests.get(f"{API_URL}/student/placement-status", headers=headers)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to get placement status'}), response.status_code
            
        return response.json()
        
    except Exception as e:
        print(f"Error getting placement status: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/confirmed-candidates')
def get_confirmed_candidates():
    if 'company_access_token' not in session:
        return jsonify({'error': 'Please log in first'}), 401
        
    headers = get_auth_header()
    if not headers:
        return jsonify({'error': 'Session expired'}), 401
        
    try:
        # Get confirmed candidates (students who accepted offers)
        response = requests.get(f"{API_URL}/company/confirmed-candidates", headers=headers)
        print(f"Confirmed candidates response: {response.status_code}, {response.text}")
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to load confirmed candidates'}), response.status_code
            
        return response.json()
        
    except Exception as e:
        print(f"Error getting confirmed candidates: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/job/<int:job_id>/applicants')
def view_applicants(job_id):
    if 'company_access_token' not in session:
        return redirect(url_for('company_login'))
        
    headers = {
        'Auth': session['company_access_token']
    }
    
    try:
        # Get job details
        job_response = requests.get(f"{API_URL}/company/job/{job_id}", headers=headers)
        if job_response.status_code != 200:
            flash('Job not found', 'error')
            return redirect(url_for('company_dashboard'))
            
        job = job_response.json()
        if isinstance(job, dict) and 'jobs' in job:
            job = job['jobs'][0] if job['jobs'] else {}
        
        # Get applications for this job
        print(f"Fetching applications for job {job_id}...")
        applications_response = requests.get(f"{API_URL}/company/job/{job_id}/application", headers=headers)
        print(f"Applications response: {applications_response.status_code} - {applications_response.text}")
        
        if applications_response.status_code == 200:
            applications = applications_response.json().get('applications', [])
            return render_template('applicants.html', 
                                 job=job,
                                 applications=applications)
        else:
            flash('Failed to fetch applicants', 'error')
            return redirect(url_for('company_dashboard'))
            
    except Exception as e:
        print(f"Error viewing applicants: {str(e)}")
        import traceback
        print(f"Traceback: {traceback.format_exc()}")
        flash('An error occurred while fetching applicants', 'error')
        return redirect(url_for('company_dashboard'))

if __name__ == '__main__':
    app.run(debug=True)
