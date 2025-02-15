-- +goose Up
-- +goose StatementBegin
INSERT INTO company_table(company_name, poc_name, poc_phno, industry, username, password)
VALUES('Microsoft', 'Harish', '1234567890', 'Software', 'softmicro', '1234'),
('Confluent', 'Shwetha', '9876543210', 'Software', 'fluentcon', '1234'),
('Mercedes', 'Govind', '1122334455', 'Automobile', 'cedesmer', '1234');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM company_table;
-- +goose StatementEnd
