# Cofee Shop - Coffee Shop Management System

## Overview

**Cofee Shop** is a coffee shop management system built with **Go** and **MongoDB**. It follows a **three-layered architecture** (Handlers, Services, Repositories) for clean code structure.

## Features

- Manage **Orders** (Create, Read, Update, Delete)
- Store data in **MongoDB**
- REST API for frontend integration

---

## Installation

### Prerequisites

- **Go** (1.18+)
- **MongoDB** (running instance or cloud)

### Clone Repository

```sh
git clone https://github.com/Ytram88/cofee-shop-mongo.git
cd hot-coffee
```

### Install Dependencies

```sh
go mod tidy
```

Create a `.env` file:

```env
MONGO_USER="cofeeStaff"
MONGO_PASSWORD="cofeeAdmin"
JWT_SECRET="secretJWT123"
JWT_EXPIRATION_IN_SECONDS=60*120
```

### Run Application

```sh
go run ./cmd/myapp/.
```

or if you dont have go installed, use binary in the repository

```sh
docker run -e {envs here}
```

---

## Database Schema (MongoDB Collections)

### \*\*Orders Collection (`orders`)

```json
{
  "_id": ObjectId("...")
  "order_id": "order123",
  "customer_name": "Alice Smith",
  "items": [
    { "product_id": "latte", "quantity": 2 },
    { "product_id": "muffin", "quantity": 1 }
  ],
  "status": "open",
  "created_at": "2023-10-01T09:00:00Z"
}
```

### \*\*Products Collection (`products`)

```json
{
  "_id": ObjectId("...")
  "product_id": "latte",
  "name": "Latte",
  "price": 4.50
}
```

---

## API Endpoints

### **Authorization**

| Method | Endpoint    | Description                  |
|--------|-------------|------------------------------|
| `POST` | `/register` | Creating a new account in db |
| `POST`  | `/login`    | Getting a JWT token          |

### **Orders**

| Method   | Endpoint            | Description        |
| -------- | ------------------- | ------------------ |
| `POST`   | `/orders`           | Create a new order |
| `GET`    | `/orders`           | Get all orders     |
| `GET`    | `/orders/{id}`      | Get order by ID    |
| `PUT`    | `/orders/{id}`      | Update an order    |
| `DELETE` | `/orders/{id}`      | Delete an order    |
| `POST`   | `/orders/{id}/close`| Close an order     |

### **Menu Items**

| Method   | Endpoint        | Description          |
| -------- | -------------- | -------------------- |
| `POST`   | `/menu`        | Add a new menu item  |
| `GET`    | `/menu`        | Get all menu items   |
| `GET`    | `/menu/{id}`   | Get menu item by ID  |
| `PUT`    | `/menu/{id}`   | Update a menu item   |
| `DELETE` | `/menu/{id}`   | Delete a menu item   |

### **Inventory**

| Method   | Endpoint           | Description            |
| -------- | ----------------- | ---------------------- |
| `POST`   | `/inventory`      | Add a new inventory item |
| `GET`    | `/inventory`      | Get all inventory items |
| `GET`    | `/inventory/{id}` | Get inventory item by ID |
| `PUT`    | `/inventory/{id}` | Update an inventory item |
| `DELETE` | `/inventory/{id}` | Delete an inventory item |

### **Aggregation**

| Method   | Endpoint           | Description            |
| -------- | ----------------- | ---------------------- |
| `GET`   | `/reports/total-sales`      |  Get the total sales amount |
| `GET`    | `/reports/popular-items`  | Get a list of popular menu items |
---

