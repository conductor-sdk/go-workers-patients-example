# go-workers-patients-example
Set environment variables: 
```
KEY 
SECRET
CONDUCTOR_SERVER_URL
```
Make sure your application has permission to execute workflow, 
poll and execute tasks. <br>
The quickest way to set this up is to grant Unrestricted Worker, 
Workflow Manager roles.

````
curl -X POST localhost:8083/ --data '{"localTable":"patients1","externalTable":"patients2","dob":"1983-05-21","LocalDBConnectionString":"user=postgres dbname=test host=localhost sslmode=disable","last_name":"Smith","first_name":"John","ExternalDBConnectionString":"user=postgres dbname=test host=localhost sslmode=disable"}'
````

Example input above is for local postgres database `test` with user `postgres`, 
no password, no SSL, `patients1` (meant to serve as local) and `patients2` tables (meant to serve as external) in `public` schema. <br>

table setup for successful run:
```
create table patients1 (
    first_name text,
    last_name text,
    dob date,
    family_doctor_assigned bool
);

create table patients2 (
    first_name text,
    last_name text,
    dob date,
    family_doctor_assigned bool
);
INSERT INTO patients1 VALUES ('John', 'Smith', '1983-05-21'::date, false);
INSERT INTO patients2 VALUES ('John', 'Smith', '1983-05-21'::date, true);
```
