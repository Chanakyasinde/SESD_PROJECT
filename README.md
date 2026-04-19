# Inventory Management System

Inventory Management System is a full-stack web application for managing products, stock, and orders with role-based access control.

## Project Overview

- Architecture: Clean/Layered backend with separated domain, use case, repository, and HTTP layers.
- Backend design approach: Built using OOP concepts and SOLID principles for maintainability and scalability.
- Backend: Go, Gin, MongoDB, JWT authentication.
- Frontend: React, TypeScript, Vite.
- Roles: Admin, Staff, Customer.
- Core modules: Authentication, product management, order management, and stock logging.

### Production

- Frontend: https://inventory-os-three.vercel.app/
- Backend API base: https://sesd-project-backend.onrender.com/api/v1

## API Overview

Base path: `/api/v1`

### Public Endpoints

- `POST /auth/register`
- `POST /auth/login`

### Protected Endpoints (JWT Required)

#### Users

- `GET /users/customers` (Admin, Staff)

#### Products

- `GET /products` (All authenticated users)
- `GET /products/:id` (All authenticated users)
- `POST /products` (Admin)
- `PUT /products/:id` (Admin)
- `DELETE /products/:id` (Admin)

#### Orders

- `GET /orders` (All authenticated users)
- `GET /orders/:id` (All authenticated users)
- `POST /orders` (Admin, Staff, Customer)
- `PATCH /orders/:id/status` (Admin, Staff)

### Health Check

- `GET /health`
