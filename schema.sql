CREATE TABLE users (
    id SERIAL PRIMARY KEY,          
    name TEXT NOT NULL,             
    branch TEXT NOT NULL,           
    cgpa NUMERIC(3, 2) NOT NULL,    
    active_backlogs BOOLEAN NOT NULL,   
    email_id TEXT UNIQUE NOT NULL,  
    usn TEXT UNIQUE NOT NULL,       
    counsellor_name TEXT NOT NULL   
);
