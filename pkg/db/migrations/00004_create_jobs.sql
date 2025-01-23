-- +goose Up
-- +goose StatementBegin
INSERT INTO job_table(company_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches)
VALUES(1, 'SDE', 'Internship', 24.5, 'Open Dream', '2025-06-12 15:04:05', 8.5, 2027, '{"ISE", "CSE", "CY", "DS"}'),
(2, 'SRE', 'FTE', 12, 'Open Dream', '2025-07-12 15:04:05', 7.5, 2026, '{"ISE", "CSE", "CY", "DS", "ECE", "EIE", "EEE", "ETE", "CV", "ME"}'),
(2, 'SDE', 'FTE + Internship', 64, 'Open Dream', '2025-03-12 15:04:05', 9, 2026, '{"ISE", "CSE", "CY", "DS"}'),
(3, 'Mechanical Engineer', 'FTE', 5.5, 'Dream', '2025-05-12 15:04:05', 8.5, 2026, '{"ME"}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM student_offer_table;
DELETE FROM job_table where job_id <= 4;
-- +goose StatementEnd
