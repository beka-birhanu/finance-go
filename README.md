# Financial Tracker

This is a back-end for a Financial-tracker, an app I made to learn React and to learn the basics of CRUD RESTful API in go

## Technologies

1. Go
2. Postgres
3. Docker

## Usage

Simply clone the repository and run the command:

```bash
docker-compose up
```

## API Definition (Expense)

### Create Expense

#### Request

```
POST /expenses
```

```json
{
  "title": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00"
}
```

#### Response

```
201 Created
```

```
Location: {{host}}/expenses/{{id}}
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "title": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00"
}
```

### Get Expenses

#### Get Bulk Request

```
GET /expenses?pageNumber={yourPart}&pageSize={yourPart}&filterField={yourPart}&filterValue={yourPart}&sortField={yourPart}&sortOrder={yourPart}
```

#### Response

```
200 OK
```

```json
{
  "expenses": [
    {
      "id": "286d7bbf-e6e0-4bfd-b4e0-906a613193db",
      "title": "Car Repair",
      "amount": 350.5,
      "date": "2024-02-15T08:00:00"
    },
    {
      "id": "3f91f017-af32-46b3-9c53-47adb1314c9a",
      "title": "Groceries",
      "amount": 75.3,
      "date": "2024-03-10T08:00:00"
    },
    {
      "id": "9faf45e0-38c5-4c27-b9d8-f6b04b719060",
      "title": "Electric Bill",
      "amount": 120.75,
      "date": "2024-04-05T08:00:00"
    },
    {
      "id": "805795d0-f73d-4334-94af-ce4fab83e142",
      "title": "Dinner Out",
      "amount": 60.2,
      "date": "2024-04-18T08:00:00"
    },
    {
      "id": "537126c9-f485-4674-900e-c77ae349c25f",
      "title": "Movie Tickets",
      "amount": 45.0,
      "date": "2024-05-01T08:00:00"
    }
  ],
  "totalCount": 5,
  "pageNumber": 1,
  "pageSize": 10
}
```

#### Get One Request

```
GET /expenses/{{id}}
```

#### Response

```
200 OK
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "title": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00"
}
```

### Update Expense

#### Request

```
PUT /expenses/{{id}}
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "title": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00"
}
```

#### Response

```
204 No Content
```

or

```
201 Created
```

```
Location: {{host}}/expenses/{{id}}
```

### Delete Expense

#### Request

```
DELETE /expenses/{{id}}
```

#### Response

```
204 No Content
```

## API Definition (User)

### Create User

#### Request

```
Post /users
```

```json
{
  "firstName": "Beka",
  "lastName": "Birhanu",
  "email": "romareo@gmail.com",
  "password": "************"
}
```

#### Response

```
201 Created
```

**Headers**

```
Set-Cookie: token
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "firstName": "Beka",
  "lastName": "Birhanu",
  "email": "romareo@gmail.com"
}
```

## API Definition (Authentication)

### Sign in

#### Request

```
Post /auth/signIn
```

```json
{
  "email": "romareo@gmail.com",
  "password": "************"
}
```

#### Response

```
200 Ok
```

**Headers**

```
Set-Cookie: token
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "firstName": "Beka",
  "lastName": "Birhanu",
  "email": "romareo@gmail.com"
}
```

### Sign out

#### Request

```
Post /auth/signOut
```

#### Response

```
204 No Content
```
