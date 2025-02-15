CREATE TABLE student_table(
    student_id SERIAL PRIMARY KEY,
    usn VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(15) NOT NULL,
    name VARCHAR(255) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    batch INT NOT NULL,
    cgpa REAL NOT NULL,
    num_active_backlogs INT NOT NULL, 
    email_id VARCHAR(255) NOT NULL, 
    counsellor_email_id VARCHAR(255) NOT NULL
);

CREATE TABLE company_table(
    company_id SERIAL PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL, 
    poc_name VARCHAR(255) NOT NULL,
    poc_phno VARCHAR(15) NOT NULL,
    industry VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE job_table(
    job_id SERIAL PRIMARY KEY,
    company_id INT NOT NULL,
    job_role VARCHAR(255) NOT NULL,
    job_type VARCHAR(255) NOT NULL,
    job_description TEXT NOT NULL,
    ctc REAL NOT NULL,
    salary_tier VARCHAR(15) NOT NULL,
    apply_by_date TIMESTAMP NOT NULL,
    cgpa_cutoff REAL NOT NULL,
    eligible_batch INT NOT NULL,
    eligible_branches VARCHAR[] NOT NULL,
    FOREIGN KEY(company_id) REFERENCES company_table(company_id)
);

CREATE TABLE student_offer_table(
    student_id INT NOT NULL,
    job_id INT NOT NULL,
    action VARCHAR(15) NOT NULL,
    action_date TIMESTAMP,
    act_by_date TIMESTAMP NOT NULL,
    FOREIGN KEY(student_id) REFERENCES student_table(student_id),
    FOREIGN KEY (job_id) REFERENCES job_table(job_id)
);

CREATE TABLE student_job_application_table(
    student_id INT NOT NULL,
    job_id INT NOT NULL,
    applied_on_date TIMESTAMP NOT NULL,
    FOREIGN KEY(student_id) REFERENCES student_table(student_id),
    FOREIGN KEY(job_id) REFERENCES job_table(job_id)
);

CREATE TABLE student_job_interview_table(
    student_id INT NOT NULL,
    job_id INT NOT NULL,
    venue VARCHAR(50) NOT NULL,
    interview_date TIMESTAMP NOT NULL,
    FOREIGN KEY(student_id) REFERENCES student_table(student_id),
    FOREIGN KEY(job_id) REFERENCES job_table(job_id)
);

