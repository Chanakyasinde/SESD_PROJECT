# ER Diagram — Inventory & Order Management System

## Overview

This ER diagram represents the core database structure for the Inventory & Order Management System, including users, products, orders, and stock tracking with strong data integrity and workflow control.

---

```mermaid
erDiagram

    USERS {
        uuid id PK
        varchar name
        varchar email UK
        varchar password_hash
        enum role "ADMIN | STAFF | CUSTOMER"
        timestamp created_at
        timestamp updated_at
    }

    CATEGORIES {
        uuid id PK
        varchar name UK
        text description
    }

    PRODUCTS {
        uuid id PK
        varchar name
        text description
        decimal price
        int stock_quantity
        int low_stock_threshold
        int version "Optimistic Lock"
        uuid category_id FK
        timestamp updated_at
    }

    ORDERS {
        uuid id PK
        uuid user_id FK
        decimal total_amount
        enum status "PENDING | CONFIRMED | SHIPPED | DELIVERED | CANCELLED"
        timestamp created_at
        timestamp updated_at
    }

    ORDER_ITEMS {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        int quantity
        decimal unit_price
        decimal subtotal
    }

    STOCK_LOGS {
        uuid id PK
        uuid product_id FK
        uuid order_id FK
        uuid changed_by FK
        int quantity_changed
        enum change_type "DEDUCTION | RESTOCK | ADJUSTMENT"
        timestamp created_at
    }

    USERS ||--o{ ORDERS : places
    USERS ||--o{ STOCK_LOGS : performs
    CATEGORIES ||--o{ PRODUCTS : categorizes
    ORDERS ||--o{ ORDER_ITEMS : contains
    PRODUCTS ||--o{ ORDER_ITEMS : included_in
    PRODUCTS ||--o{ STOCK_LOGS : has_history
```

---

### Flow Summary

| Phase | Description | Key Concepts |
| :--- | :--- | :--- |
| **1. Data Integrity** | `USERS` stores `password_hash` instead of plain text passwords. | **Hashing**, **Security Best Practices** |
| **2. Concurrency Control** | `PRODUCTS.version` enables optimistic locking during stock updates. | **Optimistic Locking**, **Versioning** |
| **3. Order Lifecycle** | `ORDERS.status` manages structured order progression. | **State Persistence**, **Enum Enforcement** |
| **4. Audit Trail** | `STOCK_LOGS` records every inventory change for accountability. | **Audit Logging**, **Traceability** |
