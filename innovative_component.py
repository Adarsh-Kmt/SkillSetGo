from flask import Flask, request, render_template_string, jsonify
import PyPDF2
from groq import Groq

app = Flask(__name__)

# Replace with your actual Groq API key
GROQ_API = "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"
# HTML Template for Index Page
HTML_TEMPLATE = """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Resume Matcher</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f8f9fa;
            color: #212529;
        }
        .container {
            max-width: 800px;
            margin: 50px auto;
            background: #ffffff;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,.1);
        }
        h1 {
            margin-bottom: 30px;
            color: #0d6efd;
            text-align: center;
        }
        label {
            font-weight: 500;
            margin-top: 15px;
            display: block;
            color: #212529;
        }
        input, textarea {
            width: 100%;
            padding: 10px;
            margin: 8px 0;
            border: 1px solid #ced4da;
            border-radius: 4px;
            box-shadow: none;
        }
        button {
            width: auto;
            padding: 10px 20px;
            margin: 20px 0;
            background-color: #0d6efd;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #0b5ed7;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Resume and Job Description Matcher</h1>
        <form action="/match" method="POST" enctype="multipart/form-data">
            <label for="resume">Upload Resume (PDF):</label>
            <input type="file" name="resume" id="resume" required>
            
            <label for="jobDescription">Enter Job Description:</label>
            <textarea name="jobDescription" id="jobDescription" rows="5" required></textarea>
            
            <button type="submit">Match</button>
        </form>
    </div>
</body>
</html>
"""

# HTML Template for Result Page
RESULT_TEMPLATE = """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Match Results</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f8f9fa;
            color: #212529;
        }
        .container {
            max-width: 800px;
            margin: 50px auto;
            background: #ffffff;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,.1);
        }
        h1 {
            margin-bottom: 30px;
            color: #0d6efd;
            text-align: center;
        }
        p {
            font-size: 16px;
            line-height: 1.5;
        }
        a {
            color: #0d6efd;
            text-decoration: none;
            font-weight: bold;
            display: inline-block;
            margin-top: 20px;
            padding: 10px 20px;
            background-color: #ffffff;
            border-radius: 4px;
            transition: background-color 0.3s;
        }
        a:hover {
            background-color: #f8f9fa;
        }
        #result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 4px;
            background-color: #f8f9fa;
            border: 1px solid #dee2e6;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Match Results</h1>
        <p><strong>Score:</strong> {{ score }}/10</p>
        <p><strong>Message:</strong> {{ message }}</p>
        <h2>Resume Text:</h2>
        <p>{{ resume_text }}</p>
        <h2>Job Description:</h2>
        <p>{{ job_description }}</p>
        <a href="/">Go Back</a>
    </div>
</body>
</html>
"""

def beautify_text(text):
    client = Groq(api_key=GROQ_API)
    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "system",
                "content": "Given the resume content, extract useful info and return it in a presentable format, dont't use markdown as I will be using it in a HTML page. Return in plain text. Don't use ** for bold, I will be displaying it using render_template of Flask. Use new line sequences"
            },
            {   "role": "user",
                "content": text
            }
        ],
        model="llama3-8b-8192",
    )
    return chat_completion.choices[0].message.content

@app.route('/')
def index():
    return render_template_string(HTML_TEMPLATE)

def extract_text_from_pdf(pdf_file):
    pdf_reader = PyPDF2.PdfReader(pdf_file)
    text = ""
    for page in pdf_reader.pages:
        text += page.extract_text()
    return text

def get_groq_score(resume_text, job_description):
    l=[]
    client = Groq(api_key=GROQ_API)
    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "system",
                "content": "You are a helpful assistant which enables students to identify whether their skills are aligned with the job description. Your reply should be in the form of a feedback when given a resume content and JD. Both are separated by ******. Do not use bold formatting return in plain text. Do not include any score. Limit to 50 words."
            },
            {
                "role": "user",
                "content": resume_text + "******" + job_description
            }
        ],
        model="llama3-8b-8192",
    )
    l.append(chat_completion.choices[0].message.content)
    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "system",
                "content": "You are a helpful assistant which enables students to identify whether their skills are aligned with the job description. Your reply should be in the form of a score out of 10 and only the score. Nothing else. Only the score. Return in plain text"+l[0]+"THis is the feedback provided ensure that the score matches feedback. Return only score"
            },
            {
                "role": "user",
                "content": resume_text + "******" + job_description
            }
        ],
        model="llama3-8b-8192",
    )
    l.append(chat_completion.choices[0].message.content)
    return l

@app.route('/match', methods=['POST'])
def match():
    if 'resume' not in request.files or 'jobDescription' not in request.form:
        return jsonify({'error': 'Please upload a resume and enter a job description.'}), 400

    resume = request.files['resume']
    job_description = request.form['jobDescription']

    try:
        # Extract text from resume PDF
        resume_text = extract_text_from_pdf(resume)
        r=beautify_text(resume_text)
        # Get score from Groqcd
        message = get_groq_score(r, job_description)
        score=message[1]
        m=message[0]
        
        return render_template_string(RESULT_TEMPLATE, score=score,message=m, resume_text=r, job_description=job_description)

    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
