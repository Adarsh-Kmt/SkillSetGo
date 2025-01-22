-- +goose Up
-- +goose StatementBegin
INSERT INTO student_table(usn, password, name, branch, batch, cgpa, num_active_backlogs, email_id, counsellor_email_id)
VALUES('1RV22IS002', '12345','Adarsh Kamath', 'ISE', 2026, 9.48, 0, 'adarshkamath.is22@rvce.edu.in', 'abcd@gmail.com'),
      ('1RV22CH080', '54321','Subramanian Iyer', 'CH', 2026, 6.5, 1, 'siyer.ch22@rvce.edu.in', 'dcba@gmail.com'),
      ('1RV22CS072', '98765','Aathmaram Tukuram Bhide', 'CSE', 2026, 7.85, 0, 'atbhide.cs22@rvce.edu.in', 'pqrs@gmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM student_table where student_id <= 3;
-- +goose StatementEnd
