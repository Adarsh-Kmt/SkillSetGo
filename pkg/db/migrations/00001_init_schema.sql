-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';
CREATE TABLE student_table(
                              student_id SERIAL PRIMARY KEY,
                              usn VARCHAR(15) UNIQUE NOT NULL,
                              name VARCHAR(255) NOT NULL,
                              branch VARCHAR(255) NOT NULL,
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
                              industry VARCHAR(255) NOT NULL
);

CREATE TABLE job_table(
                          job_id SERIAL PRIMARY KEY,
                          company_id INT NOT NULL,
                          job_role VARCHAR(255) NOT NULL,
                          ctc REAL NOT NULL,
                          salary_tier VARCHAR(15) NOT NULL,
                          apply_by_date TIMESTAMP NOT NULL,
                          cgpa_cutoff REAL NOT NULL,
                          eligible_branches VARCHAR[] NOT NULL,
                          FOREIGN KEY(company_id) REFERENCES company_table(company_id)
);

CREATE TABLE student_offer_table(
                                    student_id INT NOT NULL,
                                    job_id INT NOT NULL,
                                    action VARCHAR(15) NOT NULL,
                                    action_date TIMESTAMP NOT NULL,
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
                                            venue VARCHAR(50) NOT NULL,
                                            interview_date TIMESTAMP NOT NULL,
                                            interview_round INT NOT NULL,
                                            result VARCHAR(15) NOT NULL,
                                            FOREIGN KEY(student_id) REFERENCES student_table(student_id)
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE student_job_interview_table;
DROP TABLE student_job_application_table;
DROP TABLE student_offer_table;
DROP TABLE job_table;
DROP TABLE company_table;
DROP TABLE student_table;
-- +goose StatementEnd
