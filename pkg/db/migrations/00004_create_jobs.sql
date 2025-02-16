-- +goose Up
-- +goose StatementBegin
INSERT INTO job_table(company_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches, job_description)
VALUES(1, 'SDE', 'Internship', 26.5, 'Internship', '2025-06-12 15:04:05', 8.5, 2026, '{"ISE", "CSE", "CY", "DS"}', 'Backend Software Development Engineer'),
      (2, 'SDE', 'Internship', 21.5, 'Internship', '2025-06-12 15:04:05', 8.5, 2026, '{"ISE", "CSE", "CY", "DS"}', 'Frontend Software Development Engineer'),
(1, 'SRE', 'FTE', 20, 'Open Dream', '2025-07-12 15:04:05', 7.5, 2026, '{"ISE", "CSE", "CY", "DS", "ECE", "EIE", "EEE", "ETE", "CV"}', 'Site Reliability Engineer'),
      (2, 'ML Engineer', 'FTE', 27, 'Open Dream', '2025-02-20 15:04:05', 7, 2026, '{"ISE", "CSE"}', 'Machine Learning Engineer'),
      (1, 'SRE', 'FTE', 7.5, 'Dream', '2025-06-12 15:04:05', 8.5, 2026, '{"ISE", "CSE", "CY", "DS"}', 'Site Reliability Engineer'),
      (2, 'QA Engineer', 'FTE', 8.5, 'Dream', '2025-06-12 15:04:05', 8.5, 2026, '{"ISE", "CSE", "CY", "DS"}', 'QA engineer job to ensure quality assurance'),
(3, 'Mechanical Engineer', 'FTE', 5.5, 'Dream', '2025-05-12 15:04:05', 8.5, 2026, '{"ME"}', 'Mechanical Engineer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM student_job_interview_table;
DELETE FROM student_job_application_table;
DELETE FROM student_offer_table;
DELETE FROM job_table;
-- +goose StatementEnd
