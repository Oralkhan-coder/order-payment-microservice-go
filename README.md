# Order and Payment Microservice System

A robust, decoupled microservice architecture built with **Go** following **Clean Architecture** principles. This project implements a reliable Order processing flow integrated with a Payment service.

## 🚀 Architecture Overview

The system consists of two independent microservices communicating over REST:

- **Order Service (Port 8080)**: Manages order lifecycles (Pending, Paid, Failed, Cancelled).
- **Payment Service (Port 8081)**: Processes authorizations and declines.

```mermaid
flowchart TD
    Client["Client"]

    subgraph OrderSvc["Order Service :8080"]
        OH["HTTP Handler"]
        OS["Order Use Case"]
        OR["Order Repository"]
    end

    subgraph PaymentSvc["Payment Service :8081"]
        PH["HTTP / gRPC Handler"]
        PS["Payment Use Case"]
        PR["Payment Repository"]
        MQ["RabbitMQ Publisher"]
    end

    subgraph NotifSvc["Notification Service"]
        NC["RabbitMQ Consumer"]
        NR["Idempotency Repository"]
    end

    ODB[("PostgreSQL\norders_db")]
    PDB[("PostgreSQL\npayments_db")]
    RabbitMQ[["RabbitMQ"]]

    Client -->|"REST"| OH
    OH --> OS
    OS --> OR
    OR --> ODB
    OS -->|"HTTP POST /payments"| PH
    PH --> PS
    PS --> PR
    PR --> PDB
    PS --> MQ
    MQ --> RabbitMQ
    RabbitMQ --> NC
    NC --> NR
```

### Key Design Principles (Best-Case Design)
- **Thin Handlers**: Business logic and state transitions reside in the Use Case layer.
- **Dependency Injection**: Manual DI implemented at the Composition Root (`main.go`).
- **Interfaces (Ports)**: Decoupled repositories and outbound HTTP clients for testability.
- **Separate Bounded Contexts**: Each service owns its own database schema and repository implementation. No shared "models" package.

---

## 🛠 Business Rules & Requirements

### 1. Financial Accuracy
- All monetary values use `int64` (representing subunits like cents). Floating-point arithmetic (`float64`) is strictly avoided to prevent precision errors.

### 2. Order Invariants
- **Positive Amount**: Orders must have an amount > 0.
- **Non-Cancellable "Paid" Orders**: Once an order is marked as "Paid", it can no longer be "Cancelled".

### 3. Payment Limits
- **Authorization Limit**: Any payment exceeds **100,000 cents** (1,000 units) will be automatically **Declined** by the Payment Service.

### 4. Service Interaction & Resilience
- **Timeouts**: The Order Service uses a custom `http.Client` with a hard **2-second timeout** for all payment requests.
- **Error Handling**: 
  - If the Payment Service is unavailable or times out, the Order Service returns a **503 Service Unavailable** error.
  - The Order is marked as **Failed** in the database when the payment service call fails.
  - No hanging requests: The system fails fast to provide a better user experience.

---

## 🚦 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL

### Database Setup
Both services use migrations to manage their schema. Initialize your PostgreSQL databases (default: `done_db` for Order and a separate one for Payment if configured).

### Running the Services

1. **Start Payment Service**:
   ```bash
   cd payment-service
   go run cmd/main.go
   ```

2. **Start Order Service**:
   ```bash
   cd order-service
   go run cmd/main.go
   ```

---

## 📡 API Endpoints

### Order Service (`localhost:8080`)
- `POST /orders`: Create a new order (triggers payment authorization).
- `GET /orders/:id`: Retrieve order details.
- `DELETE /orders/:id`: Cancel a pending order.

### Payment Service (`localhost:8081`)
- `POST /payments`: Process a payment (Internal call from Order Service).
- `GET /payments/:order_id`: Check payment status by Order ID.

---

## 🧪 Error Handling Logic

| Scenario | Response Status | Order Status |
| :--- | :--- | :--- |
| Payment Success | 200 OK | Paid |
| Payment Declined | 200 OK | Failed |
| Payment Service Down | **503 Service Unavailable** | Failed |
| Invalid Amount (<= 0) | 400 Bad Request | N/A |
| Cancel "Paid" Order | 409 Conflict | Paid |
