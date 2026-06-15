# Vuelang — Windows 11 Installation Guide

> **Author**: Milan Madusanka, Associate TechOps Engineer  
> **Time Required**: 20–30 minutes  
> **Prerequisites**: Basic command line knowledge

---

## Required Software

| Software | Version | Download |
|----------|---------|----------|
| Go | 1.25.0+ | https://go.dev/dl/ |
| Node.js | 18 LTS+ | https://nodejs.org |
| MySQL | 8.0 | https://dev.mysql.com/downloads/installer/ |
| Git | Latest | https://git-scm.com/download/win |
| Make | Latest | Via Chocolatey (see Step 5) |

---

## Step 1 — Install Git

1. Go to **https://git-scm.com/download/win**
2. Download and run the installer
3. Keep all default options and click **Next** through to the end
4. After installation, open **Git Bash** (search it in the Start menu)

> ⚠️ **Important**: Run **all commands in this guide inside Git Bash** — not PowerShell or CMD.

---

## Step 2 — Install Go

1. Go to **https://go.dev/dl/**
2. Download `go1.25.x.windows-amd64.msi`
3. Run the installer — keep all defaults
4. **Restart Git Bash**, then verify:

```bash
go version
# Expected: go version go1.25.x windows/amd64
```

---

## Step 3 — Install Node.js

1. Go to **https://nodejs.org**
2. Download the **LTS** version (left green button)
3. Run the installer — keep all defaults
4. Verify:

```bash
node -v   # v22.x.x
npm -v    # 10.x.x
```

---

## Step 4 — Install MySQL 8.0

1. Go to **https://dev.mysql.com/downloads/installer/**
2. Download **MySQL Installer for Windows**
3. During setup, choose **Developer Default**
4. Set a **root password** during installation — **do not forget this password**
5. MySQL service starts automatically after install

Verify MySQL is running (in PowerShell):

```powershell
Get-Service -Name "MySQL*"
```

Should show `Running`. If it shows `Stopped`:

```powershell
Start-Service -Name "MySQL80"
```

---

## Step 5 — Install Make

Open **PowerShell as Administrator** (right-click Start → Windows Terminal (Admin)):

```powershell
# Install Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Make
choco install make -y
```

**Close PowerShell and restart Git Bash**, then verify:

```bash
make --version
```

---

## Step 6 — Create the MySQL Database

Open **MySQL Command Line Client** from the Start menu, enter your root password, then run:

```sql
CREATE DATABASE vuelang CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

---

## Step 7 — Clone the Repository

In Git Bash:

```bash
git clone https://github.com/milanz247/vuelang.git
cd vuelang
```

---

## Step 8 — Configure Environment

```bash
cp .env.example .env
notepad .env
```

Fill in the values:

```env
PORT=9090
ENV=development

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_root_password
DB_NAME=vuelang

JWT_SECRET=paste_a_long_random_string_here
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7

CORS_ALLOWED_ORIGINS=http://localhost:9090
DB_SEED=true
```

To generate a secure JWT secret, run in Git Bash:

```bash
openssl rand -base64 64
```

Copy the output and paste it as the value of `JWT_SECRET`.

> ⚠️ **Use `PORT=9090`** — Port 8080 is commonly reserved by Windows (Hyper-V) and will cause a bind error.

> ⚠️ **Use `DB_HOST=127.0.0.1`** — Not `localhost`. On Windows these can behave differently with MySQL.

---

## Step 9 — Install Dependencies

```bash
make install
```

This installs Air (Go hot-reload), all npm packages, and the Vuelang CLI tool.

---

## Step 10 — Build the Vue Frontend

```bash
cd ui && npm run build && cd ..
```

> ⚠️ **This step is mandatory.** The Go binary embeds the Vue frontend at compile time. Without this build, `make migrate` and `make dev` will both fail with `embed.go: no matching files found`.

---

## Step 11 — Fix Air Config for Windows

Open `.air.toml`:

```bash
notepad .air.toml
```

Find these two lines and add `.exe` to both:

```toml
# Change this:
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ."

# To this:
bin = "./tmp/main.exe"
cmd = "go build -o ./tmp/main.exe ."
```

Save and close.

---

## Step 12 — Run Database Migrations

```bash
make migrate
```

This creates all required tables in the `vuelang` database.

---

## Step 13 — Seed Development Data

```bash
make seed
```

This creates the default roles and test user accounts.

---

## Step 14 — Start the Dev Server

```bash
make dev
```

You should see:

```
  ┌──────────────────────────────────────────────────┐
  │              Vuelang V2  DEV                     │
  │   App  →  http://localhost:9090                  │
  │   .go  →  Air rebuilds  (<1s)                    │
  │   .vue →  Vite HMR  (instant)                   │
  └──────────────────────────────────────────────────┘
```

Open **http://localhost:9090** in your browser. 🎉

---

## Test Login Accounts (Development Only)

| Email | Password | Role |
|-------|----------|------|
| `superadmin@vuelang.dev` | `password123` | Super Admin |
| `admin@vuelang.dev` | `password123` | Admin |
| `demo@vuelang.dev` | `password123` | User |

> ⚠️ Never use these credentials in production.

---

## Quick Reference — Correct Order

```bash
# 1. Clone
git clone https://github.com/milanz247/vuelang.git
cd vuelang

# 2. Configure environment
cp .env.example .env
notepad .env               # fill in DB password, JWT secret, PORT=9090

# 3. Install dependencies
make install

# 4. Build Vue frontend (mandatory)
cd ui && npm run build && cd ..

# 5. Fix Air for Windows
notepad .air.toml          # add .exe to bin and cmd lines

# 6. Database setup
make migrate
make seed

# 7. Start
make dev
# → http://localhost:9090
```

---

## Troubleshooting

### `embed.go: no matching files found`
The Vue frontend has not been built yet. Fix:
```bash
cd ui && npm run build && cd ..
```

### `make: command not found`
Run `choco install make -y` in PowerShell (Admin), then restart Git Bash.

### `air: command not found` or `vuelang: command not found`
The Go bin folder is not in your PATH. Fix in PowerShell:
```powershell
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Users\$env:USERNAME\go\bin", "User")
```
Restart Git Bash.

### Port 8080 — "An attempt was made to access a socket in a way forbidden"
Windows (Hyper-V) reserves port 8080 at the kernel level. Set `PORT=9090` in `.env` — this is the recommended port for local development on Windows.

### MySQL — Error 1045: Access Denied
The password in `.env` does not match MySQL. Open MySQL Command Line Client, verify your root password, and update `DB_PASSWORD` in `.env`. Make sure `DB_HOST=127.0.0.1` (not `localhost`).

### MySQL service not running
```powershell
Start-Service -Name "MySQL80"
```

### Air — "CMD will not recognize non .exe file"
Update `.air.toml` to add `.exe` to the `bin` and `cmd` lines (see Step 11).

---

## Useful Make Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start dev server with hot-reload for Go and Vue |
| `make build` | Build optimized production binary → `dist/vuelang.exe` |
| `make run` | Run the production binary |
| `make migrate` | Run all pending database migrations |
| `make seed` | Seed development data (roles and demo users) |
| `make test` | Run all tests with race detection and coverage |
| `make clean` | Remove build artifacts and compiled frontend |

---

*For full framework documentation see [README.md](README.md). For production deployment security requirements see [SECURITY_CHECKLIST.md](SECURITY_CHECKLIST.md).*