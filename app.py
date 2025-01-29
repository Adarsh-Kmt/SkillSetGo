from flask import Flask, render_template, request, jsonify, redirect, url_for, session
from dotenv import load_dotenv
import os
import requests

load_dotenv()

app = Flask(__name__)
app.config['SECRET_KEY'] = os.getenv('FLASK_SECRET_KEY', 'dev-key-123')

# Go backend API URL
API_URL = 'http://localhost:8080'  # Adjust this to your Go server port

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
                session['access_token'] = response.json()['access_token']
                return redirect(url_for('dashboard'))
            else:
                try:
                    error_msg = response.json().get('error', 'Invalid credentials')
                    if isinstance(error_msg, dict):
                        error_msg = '; '.join(f"{k}: {v}" for k, v in error_msg.items())
                except:
                    error_msg = "Invalid credentials"
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
        # Get jobs from the Go backend
        headers = {'Auth': session['access_token']}
        print(f"Fetching jobs with token: {session['access_token']}")  # Debug print
        
        # Get student profile first to check eligibility
        profile_response = requests.get(f"{API_URL}/student/profile", headers=headers)
        print(f"Profile API Response: {profile_response.status_code} - {profile_response.text}")  # Debug print
        
        # Get all jobs
        response = requests.get(f"{API_URL}/job", headers=headers)
        print(f"Jobs API Response: {response.status_code} - {response.text}")  # Debug print
        
        jobs = []
        if response.status_code == 200:
            try:
                jobs = response.json()
                print(f"Number of jobs found: {len(jobs)}")  # Debug print
                if jobs:
                    print(f"Sample job data: {jobs[0]}")  # Debug print
                if not isinstance(jobs, list):
                    print(f"Jobs is not a list, it's a {type(jobs)}")  # Debug print
                    jobs = []
            except Exception as e:
                print(f"Error parsing jobs: {str(e)}")  # Debug print
                jobs = []
        
        # Get job offers
        offers_response = requests.get(f"{API_URL}/student/offer", headers=headers)
        print(f"Offers API Response: {offers_response.status_code} - {offers_response.text}")  # Debug print
        
        offers = []
        if offers_response.status_code == 200:
            try:
                offers_data = offers_response.json()
                offers = offers_data.get('offers', [])
            except:
                offers = []
        
        return render_template('dashboard.html', jobs=jobs, offers=offers)
    except requests.exceptions.RequestException as e:
        print(f"Dashboard error: {str(e)}")  # Debug print
        return render_template('dashboard.html', error="Failed to fetch data", jobs=[], offers=[])

@app.route('/apply/<int:job_id>', methods=['POST'])
def apply_job(job_id):
    if 'access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    headers = {'Auth': session['access_token']}
    try:
        response = requests.post(f"{API_URL}/student/apply/{job_id}", headers=headers)
        return jsonify({'success': response.status_code == 200}), response.status_code
    except requests.exceptions.RequestException:
        return jsonify({'success': False, 'error': 'Server error'}), 500

@app.route('/offer/action', methods=['POST'])
def offer_action():
    if 'access_token' not in session:
        return jsonify({'success': False, 'error': 'Not authenticated'}), 401
    
    data = request.json
    headers = {'Auth': session['access_token']}
    try:
        response = requests.put(f"{API_URL}/student/offer", headers=headers, json=data)
        return jsonify({'success': response.status_code == 200}), response.status_code
    except requests.exceptions.RequestException:
        return jsonify({'success': False, 'error': 'Server error'}), 500

@app.route('/company/login', methods=['GET', 'POST'])
def company_login():
    if request.method == 'POST':
        data = {
            'username': request.form.get('username'),
            'password': request.form.get('password')
        }
        try:
            print(f"Sending company login data to API: {data}")  # Debug print
            response = requests.post(f"{API_URL}/company/login", json=data)
            print(f"API Response: {response.status_code} - {response.text}")  # Debug print
            
            if response.status_code == 200:
                session['company_access_token'] = response.json().get('access_token')
                return redirect(url_for('company_dashboard'))
            else:
                try:
                    error_msg = response.json().get('error', 'Invalid credentials')
                    if isinstance(error_msg, dict):
                        error_msg = '; '.join(f"{k}: {v}" for k, v in error_msg.items())
                except:
                    error_msg = "Invalid credentials"
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
                try:
                    error_msg = response.json().get('error', 'Registration failed')
                    if isinstance(error_msg, dict):
                        error_msg = '; '.join(f"{k}: {v}" for k, v in error_msg.items())
                except:
                    error_msg = "Registration failed"
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
        headers = {'Auth': session['company_access_token']}
        print(f"Fetching company jobs with token: {session['company_access_token']}")  # Debug print
        response = requests.get(f"{API_URL}/company/jobs", headers=headers)
        print(f"Company Jobs API Response: {response.status_code} - {response.text}")  # Debug print
        
        jobs = []
        if response.status_code == 200:
            try:
                jobs = response.json()
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
        headers = {'Auth': session['company_access_token']}
        
        # Get form data
        eligible_branches = request.form.getlist('eligible_branches')
        ctc = float(request.form.get('ctc', 0))
        
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
            'apply_by_date': request.form.get('apply_by_date'),
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
