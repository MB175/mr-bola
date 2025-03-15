# Mr. Bola - Educational API Security Project

<img src="icon.webp" height=300px>

## Overview
Mr. Bola is an educational project written in Go, designed as a workshop template to teach **Object-Level Authorization** in API security. It provides a simple notes management API with authentication using **JWT tokens**. The project is intentionally designed with a common security vulnerability, making it an ideal hands-on challenge for learning API security.

## Challenge
The API correctly implements authentication but **fails to properly enforce authorization** on note retrieval. The objective is to identify and fix this security flaw, ensuring that users can only access their own notes.


## Running the Integration Test
To verify if the challenge has been solved, an integration test is included. This test starts the API, performs authentication, creates notes, and attempts unauthorized access.

A Makefile command has been provided to quickly run the integration test:
```shell
make audit
```

## Educational Purpose
This project helps developers understand:
- **How authentication works with JWT tokens**
- **The importance of object-level authorization**
- **How to identify and fix access control vulnerabilities**
- **Best practices for securing APIs**

## API Endpoints

### Authentication
#### `POST /auth`
Generates a JWT token for a user.
##### Request:
```json
{
  "username": "userA"
}
```
##### Response:
```json
{
  "token": "your.jwt.token"
}
```

### Create a Note
#### `POST /notes`
Creates a new note. Requires authentication.
##### Request:
```json
{
  "content": "This is my note."
}
```
##### Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "owner": "userA",
  "content": "This is my note."
}
```

### Retrieve a Note
#### `GET /notes/{id}`
Fetches a note by its ID.
- **Expected behavior**: Only the note owner should be able to retrieve it.
- **Current behavior**: Any authenticated user can access any note.

##### Response (if unauthorized user accesses the note):
```json
{
  "error": "Forbidden"
}
```

## Running the Project
### Prerequisites
- Go installed (`go version`)
- SQLite3

### Setup
1. Clone the repository:
   ```sh
   git clone https://github.com/MB175/mr-bola
   cd mr-bola
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the API:
   ```sh
   go run cmd/server/main.go
   ```

### Testing the API
Use **IntelliJ HTTP Client**, **Postman**, or **cURL**:
```sh
curl -X POST http://localhost:8080/auth -d '{"username": "userA"}' -H "Content-Type: application/json"
```

### Solution
<details>
  <summary>Click to reveal</summary>

https://gist.github.com/MB175/cb21487017b08e561fc8837eda759410
</details>


## License
This project is open-source and available for educational purposes.

