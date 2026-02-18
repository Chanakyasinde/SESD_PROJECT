# Sequence Diagram — Inventory & Order Management System

## Overview

This sequence diagram illustrates the primary workflow of the system:  
**Customer placing an order**, including authentication, stock validation, order creation, stock deduction, and audit logging.

---

```mermaid
sequenceDiagram

    actor Customer
    participant Frontend
    participant AuthMiddleware
    participant OrderController
    participant OrderService
    participant ProductRepository
    participant OrderRepository
    participant StockLogRepository
    participant Database

    Customer->>Frontend: Place Order Request
    Frontend->>AuthMiddleware: Send JWT Token
    AuthMiddleware->>AuthMiddleware: Validate Token
    AuthMiddleware-->>Frontend: Authorized

    Frontend->>OrderController: POST /orders
    OrderController->>OrderService: createOrder()

    OrderService->>ProductRepository: checkStock(productId, quantity)
    ProductRepository->>Database: SELECT stock, version
    Database-->>ProductRepository: Return stock data
    ProductRepository-->>OrderService: Stock Available

    OrderService->>OrderRepository: saveOrder()
    OrderRepository->>Database: INSERT order + order_items
    Database-->>OrderRepository: Success

    OrderService->>ProductRepository: updateStock()
    ProductRepository->>Database: UPDATE stock (Optimistic Lock)
    Database-->>ProductRepository: Success

    OrderService->>StockLogRepository: logStockChange()
    StockLogRepository->>Database: INSERT stock_log
    Database-->>StockLogRepository: Logged

    OrderService-->>OrderController: Order Created
    OrderController-->>Frontend: 201 Created
    Frontend-->>Customer: Order Confirmation
```

---

### Flow Summary

| Step | Description | Backend Concept |
| :--- | :--- | :--- |
| **1. Authentication** | JWT token validated before processing request | Middleware, Security |
| **2. Stock Validation** | Product stock checked before order creation | Business Logic Layer |
| **3. Transactional Save** | Order and order items saved atomically | Database Transaction |
| **4. Concurrency Control** | Stock updated using version check | Optimistic Locking |
| **5. Audit Logging** | Stock changes recorded in `StockLog` | Traceability |
| **6. Response Handling** | Proper HTTP status returned | RESTful API Design |

