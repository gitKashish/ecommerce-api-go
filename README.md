# E-commerce API in Go

This project provides a RESTful API for an e-commerce application, developed in Go. The API supports user authentication, product listing, and cart checkout functionalities.

## Features

- **User Login**
- **User Registration**
- **Product Listing**
- **Cart Checkout**
- **JWT Authentication**
- **MySQL Database Migrations**

## Getting Started

### Prerequisites

- Go 1.22 (utilizes new `net/http` enhancements)
- MySQL
- GNU Make

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/gitKashish/ecommerce-api-go.git
   ```
   
3. Navigate to the project directory:

   ```bash
   cd ecommerce-api-go
   ```
   
4. Install dependencies:

   ```bash
   go mod tidy
   ```

### Database Setup

1. Create a MySQL database:

   ```sql
   CREATE DATABASE ecom;
   ```
   
3. Run the `UP` database migrations:

   ```bash
   make migrate-up
   ```
   
4. Run the `DOWN` database migrations: (Drops tables)

   ```bash
   make migrate-down
   ```
   

### Running the API

1. Set up your environment variables:
  Create a `.env` file in the project's root directory.

    ```env
    PublicHost = dbPublicHost
    Port = portNumber
    DBUser = dbUserName
    DBPassword = dbPassword
    DBAdress = dbPublicHost:portNumber
    DBName = ecom
    JWTExpirationInSeconds = 3600*24*7
    JWTSecret = notSoSecret
    ```
    
3. Build the executable:

   ```bash
   make build
   ```
   
5. Start the server: (Builds & runs)

   ```bash
   make run
   ```
   
7. The API will be available at `http://localhost:8080`.

## API Endpoints

### User Authentication

#### Register a New User

- **Endpoint:** `POST /v1/register`
- **Description:** Register a new user.
- **Request Body:**
  
  ```json
  {
    "firstName": "exampleFirstName",
    "lastName": "exampleLastName",
    "email": "user@example.com",
    "password": "examplePassword"
  }
  ```
  
- **Response:**

  ```json
  {
    "message": "User registered successfully"
  }
  ```

#### User Login

- **Endpoint:** `POST /v1/login`
- **Description:** Log in an existing user.
- **Request Body:**

  ```json
  {
    "username": "exampleUser",
    "password": "examplePassword"
  }
  ```
  
- **Response:**

  ```json
  {
    "token": "jwt-token"
  }
  ```

### Products

#### Get Products List

- **Endpoint:** `GET /v1/products`
- **Description:** Retrieve a list of available products.
- **Response:**

  ```json
  [
    {
      "id": 1,
      "name": "Product 1",
      "description": "Description of product 1",
      "image": "path/to/image1",
      "price": 200.0,
      "quantity" : 6,
      "createdAt" : "2024-06-10T19:18:24Z"
    },
    {
      "id": 2,
      "name": "Product 2",
      "description": "Description of product 2",
      "image": "path/to/image2",
      "price": 150.0,
      "quantity" : 10,
      "createdAt" : "2024-06-08T19:18:24Z"
    }
  ]
  ```

### Cart

#### Checkout

- **Endpoint:** `POST /v1/cart/checkout`
- **Description:** Checkout the cart and create an order.
- **Request Body:**

  ```json
  {
    "cartItems": [
      {
        "productId": 1,
        "quantity": 2
      },
      {
        "productId": 2,
        "quantity": 1
      }
    ]
  }
  ```
  
- **Response:**

  ```json
  {
    "order_id": 14,
    "total_price": 550.0
  }
  ```

## Contributing

Contributions are most welcome! Please fork the repository and create a pull request with your changes.

Happy coding! 🎉

---
