# Expense Management GraphQL API

Base URL: `{host}/api/graph`

## **Scalars**

- `UUID`: Represents a universally unique identifier.
- `Float32`: Represents a 32-bit floating-point number.
- `Time`: Represents date and time in RFC 3339 format.

---

## **Types**

### **Expense**

| Field         | Type     | Description                                  |
| ------------- | -------- | -------------------------------------------- |
| `id`          | UUID!    | Unique identifier for the expense.           |
| `description` | String!  | Description of the expense.                  |
| `amount`      | Float32! | Amount spent in the expense.                 |
| `date`        | Time!    | Date of the expense.                         |
| `userId`      | UUID!    | User ID associated with the expense.         |
| `createdAt`   | Time!    | Creation timestamp of the expense record.    |
| `updatedAt`   | Time!    | Last update timestamp of the expense record. |

### **PaginatedExpenseResponse**

| Field      | Type          | Description                             |
| ---------- | ------------- | --------------------------------------- |
| `expenses` | `[Expense!]!` | List of expenses matching the criteria. |
| `cursor`   | String        | Cursor for pagination.                  |

---

## **Queries**

### `expense`

Fetch a single expense by `userId` and `id`.

**Request:**

```graphql
query {
  expense(userId: UUID!, id: UUID!): Expense!
}
```

**Response:**
Returns an `Expense` object.

### `expenses`

Fetch a list of expenses based on provided parameters.

**Request:**

```graphql
query {
  expenses(params: GetMultipleInput!): PaginatedExpenseResponse!
}
```

**Response:**
Returns a `PaginatedExpenseResponse` object.

---

## **Mutations**

### `createExpense`

Create a new expense.

**Request:**

```graphql
mutation {
  createExpense(data: CreateExpenseInput!): Expense!
}
```

**Response:**
Returns the created `Expense` object.

### `updateExpense`

Update an existing expense.

**Request:**

```graphql
mutation {
  updateExpense(data: UpdateExpenseInput!): Expense!
}
```

**Response:**
Returns the updated `Expense` object.

---

## **Inputs**

### **GetMultipleInput**

| Field       | Type      | Description                                    |
| ----------- | --------- | ---------------------------------------------- |
| `cursor`    | String    | Cursor for pagination (optional).              |
| `limit`     | Int       | Maximum number of results to fetch (optional). |
| `sortField` | SortField | Field to sort by (optional).                   |
| `sortOrder` | SortOrder | Order to sort (optional).                      |
| `userId`    | UUID!     | User ID associated with expenses.              |

### **CreateExpenseInput**

| Field         | Type     | Description                          |
| ------------- | -------- | ------------------------------------ |
| `description` | String!  | Description of the expense.          |
| `amount`      | Float32! | Amount spent in the expense.         |
| `date`        | Time!    | Date of the expense.                 |
| `userId`      | UUID!    | User ID associated with the expense. |

### **UpdateExpenseInput**

| Field         | Type    | Description                          |
| ------------- | ------- | ------------------------------------ |
| `description` | String  | Updated description (optional).      |
| `amount`      | Float32 | Updated amount (optional).           |
| `date`        | Time    | Updated date (optional).             |
| `userId`      | UUID!   | User ID associated with the expense. |
| `id`          | UUID!   | Unique identifier for the expense.   |

---

## **Enums**

### **SortField**

| Value    | Description             |
| -------- | ----------------------- |
| `amount` | Sort by expense amount. |
| `date`   | Sort by expense date.   |

### **SortOrder**

| Value  | Description       |
| ------ | ----------------- |
| `asc`  | Ascending order.  |
| `desc` | Descending order. |

---

## **Authentication Note**

**Note:**

- For login and logout, use the RESTful APIs, as the GraphQL API doesn't handle authentication directly.

**Reminder:** Replace `{host}` with the actual hostname when making requests to the API.
