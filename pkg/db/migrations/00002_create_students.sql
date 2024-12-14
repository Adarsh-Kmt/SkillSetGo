-- +goose Up
-- +goose StatementBegin
INSERT INTO student_table(usn, name, branch, cgpa, num_active_backlogs, email_id, counsellor_email_id)
VALUES('1RV22IS002', 'Adarsh Kamath', 'ISE', 9.48, 0, 'adarshkamath.is22@rvce.edu.in', 'abcd@gmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM student_table where student_id <= 1;
-- +goose StatementEnd
