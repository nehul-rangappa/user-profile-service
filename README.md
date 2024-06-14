# gigawrks-user-service

This backend service is designed for user profile management supporting a large diversity spread across wide range of countries.

## Technology
* Language: Go v1.22
* Database: MySQL
* Framework: GIN
* Packages: GORM
* Design Pattern: MVC
* Test Driven Development (TDD)

## Features
* User Signup
* User Login 
* Update user profile
* View User Profile
* Delete User Profile
* Retrieve countries information from external client RestCountries API and store it
* View all the available countries with necessary information
* Secure Authentication and Authorization using JWT tokens

Please check the swagger API documentation using `openapi.yaml` for complete details of the APIs

Postman collection is available in the branch `postman-test-collection`

## Project Setup
* Install the necessary requirements stated in technology section
* Clone the repository
* Setup the database and use the schema.sql to create tables if needed
* Change the environment variables in .env
* Run the application using `go run main.go`
* Consume the APIs in a web application or can be tested in Postman


## Project Structure
├── controllers\
│ ├── user.go\
│ ├── user_test.go\
│ ├── country.go\
│ ├── country_test.go\
│ ├── errors.go\
├── models\
│ ├── user.go\
│ ├── user_test.go\
│ ├── country.go\
│ ├── country_test.go\
│ ├── interfaces.go\
│ ├── mock_interfaces.go\
├── main.go\
├── schema.sql\
├── openapi.yaml\
├── .env\
├── README.md

Currently, view component is not used in this service based on our use case. It can be added if needed for your use case

## Future Enhancements
* API Key can be used for making the `/rest-countries` API secure and protected such that it can only be used by administrator.
* 3 Layered architecture or design pattern can be considered, as it is a better practice in Golang with Test Driven Development(TDD) with separation of concerns. It includes a handler layer, service layer and repository layer.
* Support API with different type of filters and optimize database queries if needed.
* More unit tests can be with increased coverage of entire code.
* Few functions handling JWT Tokens and password hash can be maintained in a separate directory `utils` which makes it cleaner.
* More custom defined errors with much better flexibility and messages can be used
* Performance testing of APIs and further optimization if necessary.


## Acknowledgements
* Thanks to RestCountries for providing the country data API.