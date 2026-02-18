# Use Case Diagram — Inventory & Order Management System

## Overview

This use case diagram represents the interactions between system actors (Admin, Staff, Customer) and the Inventory & Order Management System.

It highlights role-based access and core system functionalities.

---

```mermaid
flowchart LR

    %% Actors
    Admin((Admin))
    Staff((Staff))
    Customer((Customer))

    %% System Boundary
    subgraph Inventory_Order_System

        UC1[Login / Authenticate]
        UC2[Manage Products]
        UC3[Manage Categories]
        UC4[View Reports]
        UC5[Create Order]
        UC6[Update Order Status]
        UC7[Check Stock Availability]
        UC8[Browse Products]
        UC9[Place Order]
        UC10[Track Order]
        UC11[View Orders]

    end

    %% Admin Use Cases
    Admin --> UC1
    Admin --> UC2
    Admin --> UC3
    Admin --> UC4
    Admin --> UC11

    %% Staff Use Cases
    Staff --> UC1
    Staff --> UC5
    Staff --> UC6
    Staff --> UC7
    Staff --> UC11

    %% Customer Use Cases
    Customer --> UC1
    Customer --> UC8
    Customer --> UC9
    Customer --> UC10
```

---

### Flow Summary

| Actor | Responsibilities | System Capability |
| :--- | :--- | :--- |
| **Admin** | Full system control | Product management, reporting, oversight |
| **Staff** | Operational management | Order processing and stock validation |
| **Customer** | Purchase interaction | Browse products, place and track orders |
| **System** | Enforces business rules | Authentication, stock validation, lifecycle control |

