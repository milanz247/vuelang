# Vuelang

> Built by **Milan Madusanka**, Associate TechOps Engineer

Vuelang is a full-stack MVC web framework that combines a robust Go (Gin) backend with a modern Vue 3 frontend. Inspired by Laravel's elegant architecture, Vuelang provides a familiar MVC structure while leveraging the performance of Go and the reactivity of Vue 3. 

## 🚀 Features

* **Laravel-like MVC Architecture:** Structured with `app/controllers`, `app/models`, and `app/middleware` for clean separation of concerns.
* **Centralized Routing:** All API routes are declared in a single `routes/api.go` file, making it easy to see the entire API surface at a glance.
* **Go Backend:** Built on top of the lightning-fast [Gin](https://gin-gonic.com/) HTTP web framework, with MySQL database integration.
* **Vue 3 Frontend:** Modern frontend powered by Vue 3, [Vite](https://vitejs.dev/), [Tailwind CSS](https://tailwindcss.com/), and [Shadcn Vue](https://www.shadcn-vue.com/).
* **Blazing Fast Development:** Concurrent hot-reloading using [Air](https://github.com/air-verse/air) for the Go backend and Vite HMR for the Vue frontend.
* **Single Binary Deployment:** The `make build` command compiles the Vue frontend and embeds it directly into the Go binary for a seamless, single-file production deployment.

## 📂 Project Structure

```
vuelang/
├── app/                  # Application logic (MVC)
│   ├── controllers/      # Route handlers (e.g., UserController.go)
│   ├── middleware/       # HTTP middleware (e.g., Auth)
│   └── models/           # Data models
├── database/             # Database migrations and seeding
├── internal/             # Core platform code (config, db, server, logger)
├── routes/               # Route definitions
│   └── api.go            # Central API route registry
├── ui/                   # Vue 3 Frontend application
│   ├── src/              # Vue components, views, and assets
│   ├── package.json      # Node.js dependencies
│   └── vite.config.ts    # Vite configuration
├── main.go               # Application entry point
├── Makefile              # Build and run commands
└── .env.example          # Environment variables template
```

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed:
* [Go](https://golang.org/) (v1.25.0 or later)
* [Node.js](https://nodejs.org/) and npm (for the frontend)
* [MySQL](https://www.mysql.com/) database

## 🚦 Getting Started

1. **Clone the repository** and navigate to the project root.
2. **Setup Environment Variables:**
   Copy the example environment file and update your database credentials:
   ```bash
   cp .env.example .env
   ```
3. **Install Dependencies:**
   Run the install command to fetch Go modules, frontend packages, and install `air` for hot-reloading:
   ```bash
   make install
   ```

## 💻 Commands

Vuelang includes a `Makefile` with helpful commands to manage your development and production workflows:

### Development
Start the development server with concurrent hot-reloading for both backend (Air) and frontend (Vite):
```bash
make dev
```
* **Backend API:** `http://localhost:8080` (rebuilds on `.go` changes)
* **Frontend UI:** Proxied to Vite for instant HMR. Open `http://localhost:8080` in your browser.

### Production Build
Build the Vue 3 frontend and compile the Go binary with the UI embedded inside it:
```bash
make build
```
This will generate a single, production-ready binary at `dist/vuelang`.

### Run Production
Run the compiled production binary:
```bash
make run
```

### Cleanup
Remove build artifacts, `dist/` directories, and temporary files:
```bash
make clean
```

## 📖 How Routing Works

Adding a new API resource is straightforward and intuitive:
1. Create a model in `app/models/`.
2. Create a controller in `app/controllers/`.
3. Register the route in `routes/api.go`.

```go
// routes/api.go
users := &controllers.UserController{DB: db}

auth := api.Group("/", middleware.Auth())
{
    auth.GET   ("/users",     users.Index)
    auth.GET   ("/users/:id", users.Show)
    auth.POST  ("/users",     users.Store)
}
```
