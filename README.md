# Go Zarinpal Payment Service

A powerful and secure backend service built in Go (Golang) that provides a complete system for user authentication (JWT) and integration with the Zarinpal online payment gateway. This project is fully containerized with Docker for easy setup and deployment.

## Features

* **Full User Authentication:** Secure user registration (`bcrypt` hashing) and login (`JWT` tokens) .
* **Secure JWT Middleware:** All sensitive routes are protected by JWT authentication middleware .
* **Zarinpal Payment Integration:** Complete payment lifecycle:
    * Request payment and generate a gateway URL.
    * Handle and verify the payment callback from Zarinpal.
* **User Payment History:** Authenticated users can retrieve their own transaction history.
* **User Profile Management:** Authenticated users can view and manage their profile.
* **3-Tier Architecture:** Clean separation of concerns (Controller, Service, Repository).
* **Dockerized:** Fully containerized with `Dockerfile` and `docker-compose.yml` for a reproducible environment that includes the Go app and the MySQL database.
* **Configuration-Based:** All settings are managed via `.env` files.
## Technology Stack

* **Go (Golang)**
* **Gin:** High-performance HTTP web framework.
* **MySQL:** Database for storing user and payment data.
* **Docker & Docker Compose:** For containerization and service orchestration.
* **`go-sql-driver/mysql`:** MySQL driver for Go.
* **`golang-jwt/jwt`:** For generating and validating JWTs.
* **`golang.org/x/crypto/bcrypt`:** For password hashing.
* **`joho/godotenv`:** For loading `.env` files.

## Getting Started

You can run this project in two ways. The recommended method is using Docker.

### Method 1: Run with Docker (Recommended)

This is the simplest way to get the entire environment (Go app + MySQL database) up and running.

**Prerequisites:**
* Docker
* Docker Compose

**Steps:**

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/AryaTabani/Zarrinpal.git
    cd your-project-name
    ```

2.  **Create your environment file:**
    Copy the example file to `.env` (This file is ignored by Git).
    ```sh
    cp .env.example .env
    ```

3.  **Edit your `.env` file:**
    Open `.env` and fill in your details. **Crucially**, for Docker, the `DB_HOST` must be the name of the database service defined in `docker-compose.yml`.

    ```ini
    # .env file for Docker
    DB_USER=your_user
    DB_PASSWORD=your_strong_password
    DB_NAME=zarinpal_db
    DB_HOST=db  # <-- MUST be 'db' (the service name in docker-compose)
    DB_PORT=3306

    JWT_SECRET_KEY=your_super_secret_jwt_key
    ZARINPAL_MERCHANT_ID=your_zarinpal_merchant_id
    APP_BASE_URL=http://localhost:8080
    ```

4.  **Build and run with Docker Compose:**
    This command will build the Go application image, pull the MySQL image, and start both containers.
    ```sh
    docker-compose up --build
    ```

The application is now running on `http://localhost:8080`.

### Method 2: Run Locally (Manual Setup)

**Prerequisites:**
* Go 1.18+
* A running MySQL instance (on your local machine)

**Steps:**

1.  **Clone and enter the repository:**
    ```sh
    git clone https://github.com/AryaTabani/Zarrinpal.git
    cd your-project-name
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Create and edit your `.env` file:**
    Copy `.env.example` to `.env`. This time, `DB_HOST` must point to your local MySQL server.

    ```ini
    # .env file for Local/Manual setup
    DB_USER=your_local_db_user
    DB_PASSWORD=your_local_db_password
    DB_NAME=zarinpal_db
    DB_HOST=localhost  # <-- Use 'localhost' or '127.0.0.1'
    DB_PORT=3306

    JWT_SECRET_KEY=your_super_secret_jwt_key
    ZARINPAL_MERCHANT_ID=your_zarinpal_merchant_id
    APP_BASE_URL=http://localhost:8080
    ```

4.  **Run the application:**
    ```sh
    go run main.go
    ```

The server will start on `http://localhost:8080` [cite: `main.go`].

## Environment Variables

Create a `.env.example` file in your root directory so others know what to set up.

```ini
# .env.example
# Copy this file to .env and fill in your values.

# Database Configuration
# For Docker, use DB_HOST=db
# For Local, use DB_HOST=localhost
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=3306

# JWT Configuration
JWT_SECRET_KEY=

# Zarinpal Configuration
ZARINPAL_MERCHANT_ID=

# Application Configuration
# This is used for the Zarinpal callback URL
APP_BASE_URL=http://localhost:8080
````

## API Endpoints

All endpoints are prefixed by your `APP_BASE_URL`.

### Authentication

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/register` | Registers a new user [cite: `userController.go`]. |
| `POST` | `/login` | Logs in a user and returns a JWT [cite: `userController.go`]. |

### User (Authentication Required)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/me` | Gets the profile information of the currently logged-in user. |
| `PUT` | `/me` | Updates the profile information of the currently logged-in user. |
| `GET` | `/payments/history` | Gets a list of all payments made by the currently logged-in user. |

### Payments (Authentication Required)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/payment/request` | Creates a new payment request with Zarinpal and returns a payment URL [cite: `paymentController.go`]. |
| `GET` | `/payment/callback` | The callback URL that Zarinpal will redirect to after payment. It verifies the payment and updates its status [cite: `paymentController.go`]. |


```
```
