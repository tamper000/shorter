# Shorter — URL Shortener

> A simple and experimental URL shortening service written in Go

## Description

**Shorter** is a lightweight web service that allows you to shorten long URLs using a minimal web interface.  
The project was created as an experimental implementation of a URL shortener with support for multiple databases and Redis caching.

---

## 🔧 Features

- ✅ Shorten URLs via a web interface
- ✅ Cache popular links using Redis
- ✅ Support for two databases: PostgreSQL or MySQL (selected via build tags)
- ✅ Authentication using JWT

---

## 🛠️ Technologies Used

- **Go (Golang)** — primary programming language
- **Redis** — caching of URLs
- **PostgreSQL / MySQL** — data storage (chosen at build time)
- **JWT** — API protection
- **HTML/JS** — minimal frontend UI

---

## 🧪 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/tamper000/shorter.git
   cd shorter```

2. Install dependencies:
   ```bash
   go get ./...
    ```

3. Create a config file:
   ```bash
   cp config.example.yaml config.yaml
   ```

4. Open `config.yaml` and adjust settings to match your environment:
   - Database settings (MySQL or PostgreSQL)
   - Redis configuration
   - Server and JWT settings

5. Run the app:
   ```bash
   go run cmd/main.go
   ```

6. Open in browser:
   ```
   http://localhost:8080
   ```

---

## 🧪 Testing

To run tests:
```bash
go test ./...
```

---

## 📁 Configuration

Example config: `config.example.yaml`.  
The project uses a YAML-based configuration file where you can set:

- MySQL or PostgreSQL settings
- Redis parameters
- JWT configuration
- Server port

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).
