-- +goose Up
-- +goose StatementBegin
INSERT INTO company_table(company_name, poc_name, poc_phno, industry)
VALUES('Microsoft', 'Harish', '1234567890', 'Software'),
('Confluent', 'Shwetha', '9876543210', 'Software'),
('Mercedes', 'Govind', '1122334455', 'Automobile');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM company_table where company_id <= 3;
-- +goose StatementEnd
