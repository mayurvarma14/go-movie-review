# Go Movie Review API

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://example.com/build-status)
[![Go Version](https://img.shields.io/badge/go-1.23-blue.svg)](https://go.dev)
[![Docker Ready](https://img.shields.io/badge/docker-ready-blueviolet.svg)](https://www.docker.com/)

A RESTful API built with Go, Gin Gonic, and MongoDB for managing movie reviews. This API provides functionalities for user authentication, genre and movie management, and submitting movie reviews. It's designed with clear separation of roles (Admin and User) and JWT-based authentication for secure access.

## ✨ Features

*   **User Authentication & Authorization:** Secure signup, login, and JWT-based authentication with Admin and User roles.
*   **Genre Management:** Admins can create, read, update, and delete movie genres.
*   **Movie Management:** Admins can create, read, update, delete movies, and users can search and filter movies.
*   **Movie Reviews:** Users (and Admins) can submit reviews for movies, and view reviews for specific movies or all reviews by a user.
*   **Pagination:** Implemented for fetching lists of users, genres, and movies.
*   **Search & Filter:** Movies can be searched by name and filtered by genre.
*   **Dockerized:** Easy setup and deployment with Docker and Docker Compose.

## 🛠️ Tech Stack

*   **Go:** Backend programming language
*   **Gin Gonic:** Web framework for Go
*   **MongoDB:** Database
*   **JWT (JSON Web Tokens):** For authentication and authorization
*   **Docker & Docker Compose:** Containerization

## 🚀 Getting Started

### Prerequisites

*   [Docker](https://www.docker.com/get-started/) (Recommended for easy setup)
*   [Go](https://go.dev/dl/) (If you want to run the application locally without Docker)

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone git@github.com:mayurvarma14/go-movie-review.git
    cd go-movie-review
    ```

2.  **Environment Configuration:**
    *   Copy `sample.env` to `.env` and update the environment variables with your MongoDB credentials and secret key.
    *   Alternatively, environment variables can be set in `docker.env` for Docker Compose.

3.  **Run with Docker Compose (Recommended):**
    ```bash
    docker-compose up --build
    ```
    The API will be accessible at `http://localhost:8080`.

4.  **Run Locally (Optional - requires Go setup):**
    ```bash
    go run main.go
    ```
    Ensure MongoDB is running and accessible based on your `.env` configuration.

### API Endpoints

Explore the API endpoints using the provided `demo.http` file. You can use REST client extensions in VS Code or other tools to execute these requests. Key endpoints include:

*   `/users/signup`, `/users/login`: User registration and login.
*   `/users`: Get all users (Admin only), `/users/{user_id}`: Get a specific user.
*   `/genres`: Genre management endpoints (Admin for create, update, delete).
*   `/movies`: Movie management endpoints (Admin for create, update, delete, User for search/filter).
*   `/reviews`: Review management endpoints (User for add, Owner/Admin for delete).

### Authentication

*   **Admin User:** Has full access to manage genres, movies, and users.
*   **Regular User:** Can access movies, genres, and add/manage their own reviews.
*   **JWT Bearer Token:**  Required for protected endpoints. Obtain tokens after login and include them in the `Authorization` header as `Bearer <token>`.

## 📝 Demo Requests

Refer to the `demo.http` file for example requests to test the API functionalities, including user creation, login, genre/movie management, and review submissions.


