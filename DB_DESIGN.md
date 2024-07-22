# Database Schema Documentation

## 1. Table: Users

### Schema

| Column       | Type     | Constraints      | Description                               |
| ------------ | -------- | ---------------- | ----------------------------------------- |
| Id           | UUID     | Primary Key      | Unique identifier for the user.           |
| Username     | VARCHAR  | Not Null, Unique | Username of the user.                     |
| PasswordHash | VARCHAR  | Not Null         | Hashed password of the user.              |
| CreatedAt    | DATETIME | Not Null         | Timestamp when the user was created.      |
| UpdatedAt    | DATETIME | Not Null         | Timestamp when the user was last updated. |

### Relationships

- **Expenses**: One-to-many relationship with `Expenses`. Each user can have multiple expenses.

## 2. Table: Expenses

### Schema

| Column      | Type         | Constraints                | Description                                  |
| ----------- | ------------ | -------------------------- | -------------------------------------------- |
| Id          | UUID         | Not Null                   | Unique identifier for the expense.           |
| Description | VARCHAR      | Not Null                   | Description of the expense.                  |
| Amount      | DECIMAL      | Not Null, Positive         | Amount of the expense.                       |
| Date        | DATETIME     | Not Null                   | Date when the expense occurred.              |
| UserId      | UUID         | Foreign Key to Users table | Identifier of the user who made the expense. |
| CreatedAt   | DATETIME     | Not Null                   | Timestamp when the expense was created.      |
| UpdatedAt   | DATETIME     | Not Null                   | Timestamp when the expense was last updated. |
| PRIMARY KEY | (Id, UserId) |                            | Composite primary key on `Id` and `UserId`.  |

### Relationships

- **User**: Many-to-one relationship with `Users`. Each expense is linked to a single user.

### Notes

- **UUID** is used as a unique identifier for both `Users` and `Expenses` to ensure global uniqueness.
- **VARCHAR** fields for `Username` and `PasswordHash` in `Users` to store textual data.
- **DATETIME** fields for `CreatedAt` and `UpdatedAt` to track timestamps.
- **DECIMAL** type for `Amount` in `Expenses` to represent monetary values accurately.
- **Foreign Key** relationship from `Expenses` to `Users` to establish the linkage between expenses and users.

### Indexes

- **Users**

  - Unique index on `Username` to enforce uniqueness.

- **Expenses**
  - Composite primary key on `(Id, UserId)` to ensure uniqueness and establish a composite relationship with `Users`.
