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
    """Get authorization header with token"""
    token = None
    if 'company_access_token' in session:
        token = session['company_access_token']
    elif 'access_token' in session:
        token = session['access_token']
    
    if not token:
        return None
    return {'Auth': token}

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
    try:
        print("\n=== Dashboard Route ===")
        print(f"Current Session: {session}")
        
        # Check if user is logged in
        if 'access_token' not in session:
            print("No access_token in session")
            flash('Please log in first', 'error')
            return redirect(url_for('login'))
        
        # Get auth header
        headers = get_auth_header()
        print(f"Auth Headers: {headers}")
        
        if not headers:
            print("No auth headers")
            session.clear()
            flash('Session expired, please log in again', 'error')
            return redirect(url_for('login'))
        
        # Extract user ID from token
        try:
            token = session['access_token']
            payload = token.split('.')[1]
            # Add padding
            payload += '=' * ((4 - len(payload) % 4) % 4)
            decoded = base64.b64decode(payload)
            payload_data = json.loads(decoded)
            user_id = payload_data.get('id')
            
            if not user_id:
                print("No user ID in token")
                session.clear()
                flash('Invalid session', 'error')
                return redirect(url_for('login'))
                
            print(f"User ID from token: {user_id}")
        except Exception as e:
            print(f"Error decoding token: {str(e)}")
            session.clear()
            flash('Invalid session', 'error')
            return redirect(url_for('login'))
        
        # Get student profile
        profile_url = f"{API_URL}/student/{user_id}/profile"
        print(f"Fetching profile from: {profile_url}")
        print(f"Using headers: {headers}")
        
        profile_response = requests.get(profile_url, headers=headers)
        print(f"Profile Response Status: {profile_response.status_code}")
        print(f"Profile Response Body: {profile_response.text}")
        
        if profile_response.status_code == 401:
            print("Unauthorized access to profile")
            session.clear()
            flash('Session expired, please log in again', 'error')
            return redirect(url_for('login'))
        
        if profile_response.status_code != 200:
            print(f"Error fetching profile: {profile_response.status_code}")
            flash('Error loading profile', 'error')
            return render_template('dashboard.html', error="Failed to load profile")
        
        try:
            response_data = profile_response.json()
            profile_data = response_data.get('profile', {})
            print(f"Profile data: {profile_data}")
            
            # Map backend field names to frontend field names
            profile_data['email'] = profile_data.get('email_id')
            profile_data['num_of_backlogs'] = profile_data.get('num_active_backlogs')
            profile_data['graduation_year'] = profile_data.get('batch')
            
            # Ensure all required fields are present
            required_fields = ['name', 'usn', 'email', 'cgpa', 'branch', 'batch', 'num_of_backlogs', 'counsellor_email_id']
            for field in required_fields:
                if not profile_data.get(field):
                    print(f"Missing or empty field: {field}")
                    profile_data[field] = "Not available"
                    
        except json.JSONDecodeError as e:
            print(f"Error decoding profile response: {str(e)}")
            flash('Error loading profile', 'error')
            return render_template('dashboard.html', error="Failed to parse profile data")
        
        # Get jobs data
        print("Fetching jobs data...")
        jobs_response = requests.get(f"{API_URL}/student/job", headers=headers)
        print(f"Jobs Response Status: {jobs_response.status_code}")
        print(f"Jobs Response Body: {jobs_response.text}")
        
        jobs_data = []
        if jobs_response.status_code == 200:
            try:
                response_json = jobs_response.json()
                print(f"Jobs response JSON: {response_json}")
                jobs_data = response_json.get('jobs', [])
                if not jobs_data and isinstance(response_json, list):
                    jobs_data = response_json  # Handle case where response is direct array
                print(f"Jobs data: {jobs_data}")
                
                # Get applied jobs to mark which jobs the user has applied to
                applied_response = requests.get(f"{API_URL}/student/job/apply", headers=headers)
                print(f"Applied jobs response: {applied_response.status_code} - {applied_response.text}")
                
                if applied_response.status_code == 200:
                    try:
                        applied_json = applied_response.json()
                        applied_data = applied_json.get('jobs', []) or []  # Use empty list if null
                        if not applied_data and isinstance(applied_json, list):
                            applied_data = applied_json  # Handle case where response is direct array
                        applied_job_ids = {job.get('job_id') for job in applied_data if job and job.get('job_id')}
                        print(f"Applied job IDs: {applied_job_ids}")
                        
                        # Mark jobs as applied or not
                        for job in jobs_data:
                            job_id = job.get('job_id')
                            job['has_applied'] = job_id in applied_job_ids
                            
                            # Check if user meets CGPA requirement
                            try:
                                user_cgpa = float(profile_data.get('cgpa', 0))
                                job_cgpa = float(job.get('cgpa_cutoff', 0))
                                job['can_apply'] = (
                                    not job['has_applied'] and 
                                    user_cgpa >= job_cgpa
                                )
                            except (ValueError, TypeError) as e:
                                print(f"Error comparing CGPA for job {job_id}: {str(e)}")
                                job['can_apply'] = not job['has_applied']  # Default to allowing apply if CGPA comparison fails
                    except json.JSONDecodeError as e:
                        print(f"Error decoding applied jobs response: {str(e)}")
                
            except json.JSONDecodeError as e:
                print(f"Error decoding jobs response: {str(e)}")
                flash('Error loading jobs', 'error')
        elif jobs_response.status_code == 401:
            session.clear()
            flash('Session expired, please log in again', 'error')
            return redirect(url_for('login'))
        else:
            flash('Error loading jobs', 'error')
        
        return render_template('dashboard.html', 
                             profile=profile_data,
                             jobs=jobs_data)
                             
    except Exception as e:
        import traceback
        print(f"Exception in dashboard route: {str(e)}")
        print(f"Traceback: {traceback.format_exc()}")
        flash('An unexpected error occurred', 'error')
        return redirect(url_for('login'))

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

@app.route('/company/login', methods=['GET', 'POST'])
def company_login():
    if request.method == 'POST':
        try:
            data = {
                'username': request.form['username'],
                'password': request.form['password']
            }
            print("Sending company login data to API:", data)
            
            response = requests.post(f"{API_URL}/company/login", json=data)
            print(f"API Response: {response.status_code} - {response.text}")
            
            if response.status_code == 200:
                token = response.json().get('access_token')
                if token:
                    session['company_access_token'] = token
                    session['user_type'] = 'company'
                    return redirect(url_for('company_dashboard'))
                else:
                    flash('Invalid response from server', 'error')
            else:
                flash('Invalid credentials', 'error')
                
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

@app.route('/company/dashboard')
def company_dashboard():
    print("\n=== Company Dashboard Route ===")
    print("Current Session:", session)
    
    if 'company_access_token' not in session:
        print("No company_access_token in session")
        flash('Please log in first', 'error')
        return redirect(url_for('company_login'))
        
    headers = get_auth_header()
    print("Auth Headers:", headers)
    
    if not headers:
        flash('Session expired', 'error')
        return redirect(url_for('company_login'))
        
    # Get company jobs
    response = requests.get(f"{API_URL}/company/job", headers=headers)
    print(f"Jobs Response Status: {response.status_code}")
    print(f"Jobs Response Body: {response.text}")
    
    if response.status_code == 401:
        session.clear()
        flash('Session expired, please log in again', 'error')
        return redirect(url_for('company_login'))
        
    if response.status_code != 200:
        print("Unauthorized access")
        flash('Failed to load jobs', 'error')
        return redirect(url_for('company_login'))
        
    try:
        jobs_data = response.json()
        # Get company name from profile
        profile_response = requests.get(f"{API_URL}/company/profile", headers=headers)
        company_name = "Company"  # Default name
        
        if profile_response.status_code == 200:
            company_name = profile_response.json().get('name', 'Company')
            
        return render_template('company_dashboard.html', company_name=company_name)
        
    except Exception as e:
        print(f"Error rendering dashboard: {str(e)}")
        flash('An error occurred', 'error')
        return redirect(url_for('company_login'))

@app.route('/student/<string:usn>/profile')
def get_student_profile(usn):
    try:
        if 'access_token' not in session and 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        # Get student ID from applicants list
        response = requests.get(f"{API_URL}/student/{usn}", headers=headers)
        print(f"Student ID response for {usn}:", response.status_code, response.text)
        
        if response.status_code != 200:
            return jsonify({'error': 'Student not found'}), response.status_code
            
        student_id = response.json().get('id')
        if not student_id:
            return jsonify({'error': 'Student ID not found'}), 400
            
        # Now get the profile using student ID
        response = requests.get(f"{API_URL}/student/{student_id}/profile", headers=headers)
        print(f"Profile response for student {student_id}:", response.status_code, response.text)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to load profile'}), response.status_code
            
        return response.json()
        
    except Exception as e:
        print(f"Error getting student profile: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/job/<int:job_id>/applicants')
def get_job_applicants(job_id):
    try:
        if 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        response = requests.get(f"{API_URL}/company/job/{job_id}/applicants", headers=headers)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to load applicants'}), response.status_code
            
        return response.json()
        
    except Exception as e:
        print(f"Error getting job applicants: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/job/offer', methods=['POST'])
def offer_job():
    try:
        if 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        data = request.get_json()
        response = requests.post(f"{API_URL}/company/job/offer", 
                               headers=headers,
                               json=data)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to offer job'}), response.status_code
            
        return jsonify({'message': 'Job offered successfully'})
        
    except Exception as e:
        print(f"Error offering job: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/job/reject', methods=['POST'])
def reject_applicant():
    try:
        if 'company_access_token' not in session:
            return jsonify({'error': 'Please log in first'}), 401
            
        headers = get_auth_header()
        if not headers:
            return jsonify({'error': 'Session expired'}), 401
            
        data = request.get_json()
        response = requests.post(f"{API_URL}/company/job/reject", 
                               headers=headers,
                               json=data)
        
        if response.status_code == 401:
            session.clear()
            return jsonify({'error': 'Session expired'}), 401
            
        if response.status_code != 200:
            return jsonify({'error': 'Failed to reject applicant'}), response.status_code
            
        return jsonify({'message': 'Applicant rejected successfully'})
        
    except Exception as e:
        print(f"Error rejecting applicant: {str(e)}")
        return jsonify({'error': 'An unexpected error occurred'}), 500

@app.route('/company/post-job', methods=['POST'])
def post_job():
    if 'company_access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    try:
        headers = get_auth_header()
        
        # Get form data
        eligible_branches = request.form.getlist('eligible_branches')
        ctc = float(request.form.get('ctc', 0))
        
        # Format the date properly
        apply_by_date = request.form.get('apply_by_date')
        try:
            # Try to parse the datetime from the form
            date_obj = datetime.fromisoformat(apply_by_date.replace('Z', '+00:00'))
            # Format it as required by the Go backend
            formatted_date = date_obj.strftime('%Y-%m-%d %H:%M:%S')
        except Exception as e:
            print(f"Date parsing error: {str(e)}")
            return render_template('company_dashboard.html', error="Invalid date format")
        
        # Determine salary tier based on CTC
        if ctc >= 20:
            salary_tier = 'Open Dream'
        elif ctc >= 10:
            salary_tier = 'Dream'
        else:
            salary_tier = 'Mass Recruitment'
            
        data = {
            'job_role': request.form.get('job_role'),
            'job_type': request.form.get('job_type'),
            'ctc': ctc,
            'salary_tier': salary_tier,
            'apply_by_date': formatted_date,
            'cgpa_cutoff': float(request.form.get('cgpa_cutoff', 0)),
            'eligible_batch': int(request.form.get('eligible_batch', 0)),
            'eligible_branches': eligible_branches
        }
        
        print(f"Sending job post data: {data}")  # Debug print
        response = requests.post(f"{API_URL}/job", json=data, headers=headers)
        print(f"Job post response: {response.status_code} - {response.text}")  # Debug print
        
        if response.status_code == 200:
            return redirect(url_for('company_dashboard'))
        else:
            try:
                error_msg = response.json().get('error', 'Failed to post job')
                if isinstance(error_msg, dict):
                    error_msg = '; '.join(f"{k}: {v}" for k, v in error_msg.items())
            except:
                error_msg = "Failed to post job"
            return render_template('company_dashboard.html', error=error_msg)
            
    except Exception as e:
        print(f"Job posting error: {str(e)}")  # Debug print
        return render_template('company_dashboard.html', error=f"Error: {str(e)}")

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

if __name__ == '__main__':
    app.run(debug=True)
