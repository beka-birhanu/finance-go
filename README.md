# Go Backend Clean Architecture (Financial Tracker)

This is a back-end for a Financial Tracker app, built to learn Clean Architecture and Domain-Driven Design using vanilla Go. Feel free to use this project as a template for your own Go backend projects.

I appreciate any feedback on the project; it helps everyone, especially me.

## Technologies

- Go
- Postgres
- Docker

## About Me

Hi, I’m Beka Birhanu. I’m currently part of the A2SV training program and work with Clean Architecture on a daily basis.

### How to Run This Project

You can run this Go Backend Clean Architecture project with or without Docker. Here’s how to do both:

- **Clone the project**

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone https://github.com/beka-birhanu/finance-go.git

# Move to the project root directory
cd finance-go
```

#### Run Without Docker

2. Install Go and Postgres if not already installed on your machine.
1. Edit a `.env` with your configuration.
1. Run `make run`.
1. Access the API at `http://localhost:8080`.

#### Run With Docker

2. Install Docker and Docker Compose.
3. Run `docker-compose up -d`.
4. Access the API at `http://localhost:8080`.

### How to Run Tests

```bash
# Run all tests
make test
```

## File Structure

![file structure](./assets/file_structure.png)

### API Documentation

#### User

**Create User**

- **Request**

  ```
  POST api/v1/users
  ```

  ```json
  {
    "username": "beka_birhanu",
    "password": "************"
  }
  ```

- **Response**

  ```
  201 Created
  ```

  ```json
  {
    "id": "00000000-0000-0000-0000-000000000000",
    "username": "beka_birhanu"
  }
  ```

  **Headers**

  ```
  Set-Cookie: token=<token_value>; HttpOnly; Secure
  ```

**Authentication**

- **Sign In**

  - **Request**

    ```
    POST api/v1/auth/signIn
    ```

    ```json
    {
      "username": "beka_birhanu",
      "password": "************"
    }
    ```

  - **Response**

    ```
    200 OK
    ```

    **Headers**

    ```
    Set-Cookie: token=<token_value>; HttpOnly; Secure
    ```

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "username": "beka_birhanu"
    }
    ```

- **Sign Out**

  - **Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    POST api/v1/auth/signOut
    ```

  - **Response**

    ```
    204 No Content
    ```

    **Headers**

    ```
    Set-Cookie: token=; HttpOnly; Secure; Max-Age=0
    ```

**Expense**

- **Create Expense**

  - **Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    POST api/v1/users/{{userId}}/expenses
    ```

    ```json
    {
      "description": "Groceries",
      "amount": 279.7,
      "date": "2024-06-08T08:00:00Z"
    }
    ```

  - **Response**

    ```
    201 Created
    ```

    ```
    Location: {{host}}/api/v1/users/{{userId}}/expenses
    ```

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "description": "Groceries",
      "amount": 279.7,
      "date": "2024-06-08T08:00:00Z"
    }
    ```

- **Get Expenses**

  - **Get Bulk Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    GET api/v1/users/{{userId}}/expenses?cursor={base64_string_from_previous_result}&limit={yourPart}&sortField={yourPart}&filterField={yourPart}&filterValue={yourPart}&sortOrder={yourPart}
    ```

  - **Response**

    ```
    200 OK
    ```

    ```json
    {
      "expenses": [
        {
          "id": "286d7bbf-e6e0-4bfd-b4e0-906a613193db",
          "description": "Car Repair",
          "amount": 350.5,
          "date": "2024-02-15T08:00:00Z"
        },
        {
          "id": "3f91f017-af32-46b3-9c53-47adb1314c9a",
          "description": "Groceries",
          "amount": 75.3,
          "date": "2024-03-10T08:00:00Z"
        }
      ],
      "cursor": "base64_string"
    }
    ```

  - **Get One Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    GET api/v1/users/{{userId}}/expenses/{{id}}
    ```

  - **Response**

    ```
    200 OK
    ```

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "description": "Groceries",
      "amount": 279.7,
      "date": "2024-06-08T08:00:00Z"
    }
    ```

- **Update Expense**

  - **Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    PUT api/v1/users/{{userId}}/expenses/{{id}}
    ```

    ```json
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "description": "Groceries",
      "amount": 279.7,
      "date": "2024-06-08T08:00:00Z"
    }
    ```

  - **Response**

    ```
    204 No Content
    ```

    or

    ```
    201 Created
    ```

    ```
    Location: {{host}}/api/v1/users/{{userId}}/expenses/{{id}}
    ```

- **Delete Expense**

  - **Request**

    **Headers**

    ```
    Cookie: token=<token_value>
    ```

    ```
    DELETE api/v1/users/{{userId}}/expenses/{{id}}
    ```

  - **Response**

    ```
    204 No Content
    ```

### TODO

- Improve based on feedback.
- Add more test cases.
- Collect mocks into a dedicated directory.

If this project helps you in any way, show your support ❤️ by starring this project ✌️

### License

```
MIT License

Copyright (c) 2024 Beka Birhanu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### Contributing to Finance-Go

All pull requests are welcome.

---

Feel free to adjust any details according to your preferences!
