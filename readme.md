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

### Setup Environment Variables
Setup the database for use 
Create three collections
```
Inventory
Menu
Orders
```

Create a `.env` file:

```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=set the name of the DB
```

### Run Application

```sh
go run main.go
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

---

