# API Definition

## API Definition (User)

### Create User

#### Request

```
Post api/v1/users
```

```json
{
  "username": "beka_birhanu",
  "password": "************"
}
```

#### Response

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

## API Definition (Authentication)

### Sign in

#### Request

```
Post api/v1/auth/signIn
```

```json
{
  "username": "beka_birhanu"
  "password": "************"
}
```

#### Response

```
200 Ok
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

### Sign out

#### Request

**Headers**

```
Cookie: token=<token_value>
```

```
Post api/v1/auth/signOut
```

#### Response

```
204 No Content
```

**Headers**

```
Set-Cookie: token=; HttpOnly; Secure; Max-Age=0
```

## API Definition (Expense)

### Create Expense

#### Request

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

#### Response

```
201 Created
```

```
Location: {{host}/api/v1/users/{{userId}}/expenses}
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "description": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00Z"
}
```

### Get Expenses

#### Get Bulk Request

**Headers**

```
Cookie: token=<token_value>
```

```
GET api/v1/users/{{userId}}/expenses?cursor={base64_string_from_previous_result}&limit={yourPart}&sortField={yourPart}&filterField={yourPart}&filterValue={yourPart}&sortOrder={yourPart}
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
      "description": "Car Repair",
      "amount": 350.5,
      "date": "2024-02-15T08:00:00Z"
    },
    {
      "id": "3f91f017-af32-46b3-9c53-47adb1314c9a",
      "description": "Groceries",
      "amount": 75.3,
      "date": "2024-03-10T08:00:00Z"
    },
    ...
  ],
  "cursor": "base64_string"
}
```

#### Get One Request

**Headers**

```
Cookie: token=<token_value>
```

```
GET api/v1/users/{{userId}}/expenses/{{id}}
```

#### Response

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

### Update Expense

#### Request

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

#### Response

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

### Delete Expense

#### Request

**Headers**

```
Cookie: token=<token_value>
```

```
DELETE api/v1/users/{{userId}}/expenses/{{id}}
```

#### Response

```
204 No Content
```
