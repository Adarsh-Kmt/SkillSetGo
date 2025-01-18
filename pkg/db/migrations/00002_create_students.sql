-- +goose Up
-- +goose StatementBegin
INSERT INTO student_table(usn, name, branch, cgpa, num_active_backlogs, email_id, counsellor_email_id)
VALUES('1RV22IS002', 'Adarsh Kamath', 'ISE', 9.48, 0,'9090909090', 'adarshkamath.is22@rvce.edu.in', 'abcd@gmail.com'),
      ('1RV22CH080', 'Subramanian Iyer', 'CH', 6.5, 1, '9090909091','siyer.ch22@rvce.edu.in', 'dcba@gmail.com'),
      ('1RV22CS072', 'Aathmaram Tukuram Bhide', 'CSE', 7.85, 0,'9090909092', 'atbhide.cs22@rvce.edu.in', 'pqrs@gmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM student_table where student_id <= 3;
-- +goose StatementEnd
