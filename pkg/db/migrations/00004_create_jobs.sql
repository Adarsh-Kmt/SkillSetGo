-- +goose Up
-- +goose StatementBegin
INSERT INTO job_table(company_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_branches)
VALUES(1, 'SDE', 24.5, 'Open Dream', '2025-06-12 15:04:05', 8.5, '{"ISE", "CSE", "CY", "DS"}'),
(2, 'SRE', 12, 'Open Dream', '2025-07-12 15:04:05', 7.5, '{"ISE", "CSE", "CY", "DS", "ECE", "EIE", "EEE", "ETE", "CV", "ME"}'),
(2, 'SDE', 64, 'Open Dream', '2025-03-12 15:04:05', 9, '{"ISE", "CSE", "CY", "DS"}'),
(3, 'Mechanical Engineer', 5.5, 'Dream', '2025-05-12 15:04:05', 8.5, '{"ME"}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM job_table where job_id <= 4;
-- +goose StatementEnd
