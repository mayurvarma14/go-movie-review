# 🎬 Go Movie Review API 🍿

[![Go](https://img.shields.io/badge/Go-1.x-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Web_Framework-red.svg)](https://gin-gonic.com/)
[![MongoDB](https://img.shields.io/badge/MongoDB-Database-green.svg)](https://www.mongodb.com/)
[![Docker](https://img.shields.io/badge/Docker-Containerization-blue.svg)](https://www.docker.com/)
[![JWT](https://img.shields.io/badge/JWT-Authentication-purple.svg)](https://jwt.io/)

A RESTful API for managing movie reviews, built with Go, Gin, and MongoDB.  This project provides endpoints for users to sign up, log in, create/manage movie genres, add/update movies, and submit/view reviews.  It uses JWT for authentication and Docker for easy deployment.

## ✨ Features

*   **User Authentication:**
    *   Sign up and login functionality.
    *   JWT-based authentication for secure access.
    *   Password hashing using bcrypt.
    *   User roles (ADMIN, USER).
*   **Movie Management:**
    *   CRUD operations for movies (Create, Read, Update, Delete).
    *   Search movies by name and filter by genre.
    *    Admin-only access for movie creation and deletion.
*   **Genre Management:**
    *   CRUD operations for movie genres.
    *   Admin-only access.
*   **Review System:**
    *   Users can add reviews for movies.
    *   View reviews for a specific movie.
    *   Delete reviews (own reviews, potentially).
    *   View all reviews by a specific user.
*   **Pagination:** Get all users, genres and movies endpoints support pagination.
*   **Dockerized:** Easy setup and deployment using Docker Compose.
*   **Environment Variables:** Configuration via `.env` file.
*   **Input Validation:**  Uses `validator` package for robust request validation.

## 🚀 Getting Started

### Prerequisites

*   [Go](https://golang.org/dl/) (1.x or later)
*   [Docker](https://www.docker.com/products/docker-desktop)
*   [Docker Compose](https://docs.docker.com/compose/install/)
*   A text editor or IDE (VS Code recommended)

### Installation & Setup

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd go-movie-review
    ```

2.  **Create and configure a `.env` file:**

    ```bash
    cp sample.env .env
    ```
    Then, open `.env` and fill in the required values.  Replace the placeholders with your actual credentials.

3.  **Build and run with Docker Compose:**

    ```bash
    docker-compose up --build
    ```
    This command builds the MongoDB container, initializes the database with the script, and starts the Go application.  The `--build` flag ensures that the Go application is rebuilt if you make changes to the code.  The application will be accessible at `http://localhost:8080`.  The database will be accessible at `mongodb://localhost:27017`.

4.  **(Optional) Run without Docker (for development):**

    *   Make sure you have a MongoDB instance running (either locally or remotely).  Update the `.env` file with the correct connection details.

    *   Install Go dependencies:

        ```bash
        go mod download
        ```

    *   Run the application:
        ```bash
        go run main.go
        ```
        or if you use air for hot-reloading.
        ```bash
        air
        ```
### Testing

The `demo.http` file provides a series of HTTP requests that you can use to test the API endpoints. You can use an API client like [REST Client for VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client), Postman, or Insomnia to execute these requests.  Make sure to replace placeholders (like user IDs and tokens) with actual values obtained during your testing. The `demo.http` contains examples for:

*   User signup and login.
*   Retrieving user information.
*   Creating, getting, updating, and deleting genres.
*   Creating, getting, updating, searching, filtering, and deleting movies.
*   Creating, getting, and deleting reviews.

## 📁 Project Structure

```
├── controllers/        # Controllers handle API logic
│   ├── genreController.go
│   ├── movieController.go
│   ├── reviewController.go
│   └── userController.go
├── database/           # Database connection setup
│   └── databaseConnection.go
├── helpers/            # Utility functions (token generation, auth helpers)
│   ├── authHelper.go
│   └── tokenHelper.go
├── internals/
│   └── config/
│       └── config.go # Loads environment variables
├── middleware/         # Authentication middleware
│   └── authMiddleware.go
├── models/             # Data models (User, Movie, Genre, Review)
│   ├── genreModel.go
│   ├── movieModel.go
│   ├── reviewModel.go
│   └── userModel.go
├── routes/             # API route definitions
│   ├── authRoutes.go
│   ├── genreRoutes.go
│   ├── movieRoutes.go
│   ├── reviewRoutes.go
│   └── userRoutes.go
├── .env                # Environment variables (IMPORTANT: Keep this private)
├── demo.http           # Example HTTP requests for testing
├── docker-compose.yml  # Docker Compose configuration
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── init-mongo.js       # MongoDB initialization script
├── main.go             # Main application entry point
└── sample.env          # Example .env file
```

## 🛠️ Technologies

*   **Go:** Programming language.
*   **Gin:**  Web framework.
*   **MongoDB:** NoSQL database.
*   **Docker:** Containerization.
*   **JWT (JSON Web Tokens):**  Authentication.
*   **bcrypt:** Password hashing.
*   **go-playground/validator/v10:**  Request validation.
*  **joho/godotenv:** For loading environment variables.


## 🤝 Contributing

Contributions are welcome!  Please follow these steps:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature-name`).
3.  Make your changes.
4.  Commit your changes (`git commit -m "Add your commit message"`).
5.  Push to the branch (`git push origin feature/your-feature-name`).
6.  Open a pull request.
