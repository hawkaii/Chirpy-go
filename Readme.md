# Chirpy-Go

Chirpy-Go is a Twitter-esque backend application built with Golang, PostgreSQL, and JWT authentication. This project implements secure login, advanced routing, middleware integration, and a refresh token system for seamless authentication flow. Additionally, it features an extended subscription model through the Chirpy-Red module, powered by Polka Key.

---

## Features

- **File Serving**: Serve static files efficiently using `http.StripPrefix` and `http.FileServer`.
- **Health Checks**: Readiness endpoints for monitoring application status.
- **User Authentication**: Secure login with JWT (JSON Web Tokens) and refresh token system for enhanced authentication.
- **Database Management**: PostgreSQL for efficient and reliable data handling.
- **Advanced Routing and Middleware**: Designed for scalable and maintainable APIs using native Golang `net/http`.
- **Subscription Model**: Integrated with Chirpy-Red, enabling subscriptions using Polka Key.
- **Admin Metrics**: Access metrics to monitor server performance.

---

## Prerequisites

- Golang 1.19 or later
- PostgreSQL 15 or later

---

## Getting Started

### Clone the Repository

```bash
$ git clone https://github.com/hawkaii/Chirpy-go.git
$ cd Chirpy-go
```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
JWT_SECRET=dummy_jwt_secret_key
POLKA_SECRET=dummy_polka_secret_key
DB_URL="postgres://fake_user:fake_password@fake_host:5432/fake_database?sslmode=disable"
```

### Install Dependencies

```bash
$ go mod tidy
```

### Run Database Migrations

Ensure PostgreSQL is running, and then use `sqlc` and `goose` to handle migrations:

1. Generate SQL code for queries using `sqlc`:

```bash
$ sqlc generate
```

2. Apply migrations using `goose`:

```bash
$ goose -dir ./migrations postgres "postgres://fake_user:fake_password@fake_host:5432/fake_database?sslmode=disable" up
```

### Start the Application

Run the backend server:

```bash
$ go run main.go
```

The application will be accessible at `http://localhost:8080`.

---

## API Endpoints

### Static File Server

- **GET** `/app/*` - Serve static files from the `filepathRoot` directory.

### Health and Reset

- **GET** `/api/healthz` - Check server readiness.
- **GET** `/api/reset` - Reset server state (for testing purposes).

### Chirps (Tweets)

- **POST** `/api/chirps` - Create a new chirp.
- **GET** `/api/chirps` - Fetch all chirps.
- **GET** `/api/chirps/{chirpID}` - Retrieve a chirp by ID.
- **DELETE** `/api/chirps/{chirpID}` - Delete a chirp by ID.

### Users

- **POST** `/api/users` - Create a new user.
- **PUT** `/api/users` - Update user information.

### Authentication

- **POST** `/api/login` - Authenticate a user and return access and refresh tokens.
- **POST** `/api/refresh` - Generate a new access token using a valid refresh token.
- **POST** `/api/revoke` - Revoke tokens for logout.

### Subscriptions (Chirpy-Red)

- **POST** `/api/polka/webhooks` - Handle subscription-related webhooks.

### Metrics

- **GET** `/admin/metrics` - Retrieve server performance metrics.

---

## Technologies Used

- **Language**: Golang
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Routing**: Native Golang `net/http`
- **Middleware**: Custom-built using native Golang
- **Subscription Model**: Polka Key integration
- **Migrations**: `sqlc` and `goose`

---

## Contribution

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m 'Add a new feature'`).
4. Push to the branch (`git push origin feature-name`).
5. Open a Pull Request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or suggestions, please reach out:

- **GitHub**: [hawkaii](https://github.com/hawkaii)
- **Email**: [your_email@example.com](mailto:your_email@example.com)

---

Enjoy using Chirpy-Go! 🚀


