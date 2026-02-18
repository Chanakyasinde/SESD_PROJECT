# Class Diagram — Inventory & Order Management System

## Overview

This class diagram represents the object-oriented design of the Inventory & Order Management System.  
It illustrates core entities, their attributes, methods, and relationships following OOP principles and clean architecture.

---

```mermaid
classDiagram

    class User {
        +UUID id
        +String name
        +String email
        +String passwordHash
        +Role role
        +Date createdAt
        +Date updatedAt
        +login()
        +logout()
    }

    class Admin {
        +createProduct()
        +updateProduct()
        +deleteProduct()
        +viewReports()
    }

    class Staff {
        +createOrder()
        +updateOrderStatus()
        +checkStock()
    }

    class Customer {
        +browseProducts()
        +placeOrder()
        +trackOrder()
    }

    class Category {
        +UUID id
        +String name
        +String description
    }

    class Product {
        +UUID id
        +String name
        +String description
        +Decimal price
        +int stockQuantity
        +int lowStockThreshold
        +int version
        +updateStock()
        +isStockAvailable()
    }

    class Order {
        +UUID id
        +OrderStatus status
        +Decimal totalAmount
        +Date createdAt
        +Date updatedAt
        +calculateTotal()
        +updateStatus()
        +cancelOrder()
    }

    class OrderItem {
        +UUID id
        +int quantity
        +Decimal unitPrice
        +Decimal subtotal
        +calculateSubtotal()
    }

    class StockLog {
        +UUID id
        +int quantityChanged
        +ChangeType changeType
        +Date createdAt
        +logChange()
    }

    %% Inheritance
    User <|-- Admin
    User <|-- Staff
    User <|-- Customer

    %% Relationships
    User "1" --> "0..*" Order : places
    Order "1" --> "1..*" OrderItem : contains
    Product "1" --> "0..*" OrderItem : included_in
    Category "1" --> "0..*" Product : categorizes
    Product "1" --> "0..*" StockLog : generates
    User "1" --> "0..*" StockLog : performs
```

---

### Flow Summary

| Layer | Responsibility | OOP Concept |
| :--- | :--- | :--- |
| **User Hierarchy** | Admin, Staff, Customer inherit from User | **Inheritance**, **Polymorphism** |
| **Product & Inventory** | Handles stock validation and updates | **Encapsulation** |
| **Order Lifecycle** | Manages order states and transitions | **Abstraction**, **State Logic** |
| **Stock Logging** | Tracks inventory changes | **Audit Responsibility** |

