from flask import Flask, render_template, request, jsonify, redirect, url_for, session
from dotenv import load_dotenv
import os
import requests
from datetime import datetime

load_dotenv()

app = Flask(__name__)
app.config['SECRET_KEY'] = os.getenv('FLASK_SECRET_KEY', 'dev-key-123')

# Go backend API URL
API_URL = 'http://localhost:8080'  # Adjust this to your Go server port

def get_auth_header(token=None):
    """Get the authorization header with the correct token format"""
    if not token:
        token = session.get('access_token') or session.get('company_access_token')
    if not token:
        return None
    # Don't modify the token, send it exactly as received from the API
    return {'Auth': token}

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        data = {
            'usn': request.form.get('usn'),
            'password': request.form.get('password')
        }
        try:
            print(f"Sending login data to API: {data}")  # Debug print
            response = requests.post(f"{API_URL}/student/login", json=data)
            print(f"API Response: {response.status_code} - {response.text}")  # Debug print
            
            if response.status_code == 200:
                response_data = response.json()
                token = response_data.get('access_token')
                if not token:
                    return render_template('login.html', error="Invalid response from server")
                
                # Store the token exactly as received
                session['access_token'] = token
                print(f"Stored token in session: {token}")  # Debug print
                return redirect(url_for('dashboard'))
            else:
                error_msg = response.json().get('error', 'Invalid credentials')
                return render_template('login.html', error=error_msg)
        except requests.exceptions.RequestException as e:
            print(f"Login error: {str(e)}")  # Debug print
            return render_template('login.html', error="Server error")
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
    
    try:
        headers = get_auth_header()
        if not headers:
            return redirect(url_for('login'))
            
        print(f"Using headers for API calls: {headers}")  # Debug print
        
        # Get job offers and applications first to track applied jobs
        offers_response = requests.get(f"{API_URL}/student/offer", headers=headers)
        print(f"Offers API Response Status: {offers_response.status_code}")  # Debug print
        print(f"Raw Offers Response: {offers_response.text}")  # Debug print
        
        offers = []
        applied_job_ids = set()  # Track applied job IDs
        applied_jobs = []
        if offers_response.status_code == 200:
            try:
                offers_data = offers_response.json()
                print(f"Parsed offers_data: {offers_data}")  # Debug print
                
                # Get offers from the response
                if isinstance(offers_data, dict):
                    offers = offers_data.get('offers', [])
                elif isinstance(offers_data, list):
                    offers = offers_data
                else:
                    offers = []
                
                print(f"Extracted offers: {offers}")  # Debug print
                
                # Create a list of applied jobs and track job IDs
                for offer in offers:
                    print(f"Processing offer: {offer}")  # Debug print
                    job_id = offer.get('job_id')
                    print(f"Found job_id: {job_id}")  # Debug print
                    
                    if job_id:
                        applied_job_ids.add(job_id)  # Add to set for quick lookup
                        
                    # Convert action_date to a formatted string if it exists
                    action_date = offer.get('action_date')
                    if action_date:
                        try:
                            # Try to parse the date if it's a datetime string
                            date_obj = datetime.strptime(action_date, '%Y-%m-%dT%H:%M:%SZ')
                            action_date = date_obj.strftime('%Y-%m-%d')
                        except (ValueError, TypeError):
                            try:
                                # Try another common format
                                date_obj = datetime.strptime(action_date, '%Y-%m-%d %H:%M:%S')
                                action_date = date_obj.strftime('%Y-%m-%d')
                            except (ValueError, TypeError):
                                # If parsing fails, use the original string
                                pass
                    
                    # Get the status based on the action field
                    action = offer.get('action', '').upper()
                    if not action:
                        status = 'PENDING'
                    else:
                        status = action
                    
                    applied_job = {
                        'job_id': job_id,
                        'job_role': offer.get('job_role'),
                        'company_name': offer.get('company_name'),
                        'salary_tier': offer.get('salary_tier'),
                        'status': status,
                        'applied_date': action_date or 'Not Available',
                        'ctc': offer.get('ctc')
                    }
                    print(f"Created applied_job: {applied_job}")  # Debug print
                    applied_jobs.append(applied_job)
                
                print(f"Final applied_job_ids: {applied_job_ids}")  # Debug print
                print(f"Final applied_jobs: {applied_jobs}")  # Debug print
            except Exception as e:
                print(f"Error parsing offers: {str(e)}")  # Debug print
                print(f"Error type: {type(e)}")  # Debug print
                print(f"Error args: {e.args}")  # Debug print
                offers = []
                applied_jobs = []
        elif offers_response.status_code == 401:  # Unauthorized
            print("Unauthorized response, clearing session")  # Debug print
            session.clear()
            return redirect(url_for('login'))
            
        # Get all jobs
        response = requests.get(f"{API_URL}/job", headers=headers)
        print(f"Jobs API Response Status: {response.status_code}")  # Debug print
        print(f"Raw Jobs Response: {response.text}")  # Debug print
        
        jobs = []
        if response.status_code == 200:
            try:
                response_data = response.json()
                jobs_data = response_data.get('jobs', [])
                print(f"Number of jobs found: {len(jobs_data)}")  # Debug print
                
                # Process each job to ensure all required fields are present
                processed_jobs = []
                for job in jobs_data:
                    print(f"Processing job: {job}")  # Debug print
                    job_id = job.get('job_id')
                    print(f"Checking job_id {job_id} against applied_job_ids {applied_job_ids}")  # Debug print
                    
                    # Check if this job has been applied to
                    has_applied = job_id in applied_job_ids
                    print(f"Job {job_id} has_applied: {has_applied}")  # Debug print
                    
                    # Convert apply_by_date to a formatted string if it exists
                    apply_by_date = job.get('apply_by_date')
                    if apply_by_date:
                        try:
                            # Try to parse the date if it's a datetime string
                            date_obj = datetime.strptime(apply_by_date, '%Y-%m-%dT%H:%M:%SZ')
                            apply_by_date = date_obj.strftime('%Y-%m-%d')
                        except (ValueError, TypeError):
                            try:
                                # Try another common format
                                date_obj = datetime.strptime(apply_by_date, '%Y-%m-%d %H:%M:%S')
                                apply_by_date = date_obj.strftime('%Y-%m-%d')
                            except (ValueError, TypeError):
                                # If parsing fails, use the original string
                                pass
                    
                    processed_job = {
                        'job_id': job_id,
                        'job_role': job.get('job_role'),
                        'company_name': job.get('company_name'),
                        'industry': job.get('industry'),
                        'salary_tier': job.get('salary_tier'),
                        'ctc': job.get('ctc'),
                        'cgpa_cutoff': job.get('cgpa_cutoff'),
                        'apply_by_date': apply_by_date,
                        'can_apply': job.get('can_apply', True),  # Default to True if not specified
                        'has_applied': has_applied  # Set based on our tracking
                    }
                    
                    print(f"Created processed_job: {processed_job}")  # Debug print
                    
                    # If already applied, set can_apply to False
                    if has_applied:
                        processed_job['can_apply'] = False
                    
                    # Only add jobs that have required fields
                    if processed_job['job_id'] is not None and processed_job['job_role'] and processed_job['company_name']:
                        processed_jobs.append(processed_job)
                
                jobs = processed_jobs
                print(f"Final processed jobs: {jobs}")  # Debug print
            except Exception as e:
                print(f"Error parsing jobs: {str(e)}")  # Debug print
                print(f"Error type: {type(e)}")  # Debug print
                print(f"Error args: {e.args}")  # Debug print
                jobs = []
        elif response.status_code == 401:  # Unauthorized
            print("Unauthorized response, clearing session")  # Debug print
            session.clear()
            return redirect(url_for('login'))
        else:
            print(f"Error fetching jobs: {response.status_code} - {response.text}")  # Debug print
        
        print(f"Final data being sent to template:")  # Debug print
        print(f"Jobs: {jobs}")  # Debug print
        print(f"Offers: {offers}")  # Debug print
        print(f"Applied Jobs: {applied_jobs}")  # Debug print
        
        return render_template('dashboard.html', jobs=jobs, offers=offers, applied_jobs=applied_jobs)
    except requests.exceptions.RequestException as e:
        print(f"Dashboard error: {str(e)}")  # Debug print
        return render_template('dashboard.html', error="Failed to fetch data", jobs=[], offers=[], applied_jobs=[])

@app.route('/apply/<int:job_id>', methods=['POST'])
def apply_job(job_id):
    if 'access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    headers = get_auth_header()
    try:
        # First check if already applied by checking offers
        offers_response = requests.get(f"{API_URL}/student/offer", headers=headers)
        if offers_response.status_code == 200:
            offers_data = offers_response.json()
            offers = offers_data.get('offers', [])
            
            # Check if already applied to this job
            for offer in offers:
                if offer.get('job_id') == job_id:
                    return jsonify({'success': False, 'error': 'You have already applied to this job'}), 400
        
        # If not already applied, proceed with application
        print(f"Sending application request for job {job_id}")  # Debug print
        response = requests.post(f"{API_URL}/student/apply/{job_id}", headers=headers)
        print(f"Application Response: {response.status_code} - {response.text}")  # Debug print
        
        if response.status_code == 200:
            return jsonify({'success': True, 'message': 'Successfully applied to job'})
        else:
            error_msg = response.json().get('error', 'Failed to apply for job')
            return jsonify({'success': False, 'error': error_msg}), response.status_code
    except requests.exceptions.RequestException as e:
        print(f"Application error: {str(e)}")  # Debug print
        return jsonify({'success': False, 'error': 'Failed to apply for job'}), 500

@app.route('/offer/action', methods=['POST'])
def offer_action():
    if 'access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    data = request.json
    headers = get_auth_header()
    try:
        response = requests.put(f"{API_URL}/student/offer", headers=headers, json=data)
        return jsonify({'success': response.status_code == 200}), response.status_code
    except requests.exceptions.RequestException:
        return jsonify({'success': False, 'error': 'Server error'}), 500

@app.route('/company/login', methods=['GET', 'POST'])
def company_login():
    if request.method == 'POST':
        try:
            data = {
                'username': request.form.get('username'),
                'password': request.form.get('password')
            }
            print(f"Sending company login data to API: {data}")  # Debug print
            
            response = requests.post(f"{API_URL}/company/login", json=data)
            print(f"API Response: {response.status_code} - {response.text}")  # Debug print
            
            if response.status_code == 200:
                token = response.json().get('access_token')
                if not token:
                    return render_template('company_login.html', error="Invalid response from server")
                session['company_access_token'] = token
                return redirect(url_for('company_dashboard'))
            else:
                error_msg = response.json().get('error', 'Invalid credentials')
                return render_template('company_login.html', error=error_msg)
        except requests.exceptions.RequestException as e:
            print(f"Company login error: {str(e)}")  # Debug print
            return render_template('company_login.html', error="Server error")
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
    if 'company_access_token' not in session:
        return redirect(url_for('company_login'))
    
    try:
        # Get company's jobs from the Go backend
        headers = get_auth_header()
        print(f"Fetching company jobs with token: {session['company_access_token']}")  # Debug print
        response = requests.get(f"{API_URL}/job", headers=headers)  # Changed from /company/jobs to /job
        print(f"Company Jobs API Response: {response.status_code} - {response.text}")  # Debug print
        
        jobs = []
        if response.status_code == 200:
            try:
                response_data = response.json()
                jobs = response_data.get('jobs', [])  # Extract jobs from the response
                if not isinstance(jobs, list):
                    jobs = []
            except:
                jobs = []
        
        return render_template('company_dashboard.html', jobs=jobs)
    except requests.exceptions.RequestException as e:
        print(f"Company dashboard error: {str(e)}")  # Debug print
        return render_template('company_dashboard.html', error="Failed to fetch data")

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
            # Parse the datetime from the form
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

@app.route('/logout')
def logout():
    session.pop('access_token', None)
    session.pop('company_access_token', None)
    return redirect(url_for('index'))

if __name__ == '__main__':
    app.run(debug=True)
