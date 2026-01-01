# 📌 pinmyblogs

Save blogs, organize them, and read later — distraction free.

pinmyblogs is a lightweight bookmarking and blog-saving service built with **Go**, designed to be fast, simple, and
developer-friendly.

---

## 📊 Project Status

[![Check](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml/badge.svg)](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/puni9869/pinmyblogs)
![License](https://img.shields.io/github/license/puni9869/pinmyblogs)
![Stars](https://img.shields.io/github/stars/puni9869/pinmyblogs?style=social)

---

## ✨ Features

- 🔖 Save blog URLs for later reading
- 🗂️ Organize and manage bookmarks
- 🧩 Simple and extensible architecture
- 🕷️ Metadata scraping (title, favicon, etc.)

---

## 🧱 Tech Stack

- **Backend:** Go (Gin)
- **Database:** PostgreSQL or SQLite
- **Frontend:** HTML templates + Tailwind CSS + Javascript
- **Build Tools:** Make + Air

---

## 🗄️ Database

pinmyblogs supports both **PostgreSQL** and **SQLite**.

- **PostgreSQL** is recommended for production due to its robustness,
  concurrency support, and reliability.
- **SQLite** can be used for local development or lightweight testing.

### Supported Databases

- PostgreSQL (production)
- SQLite (development/testing)

---

## 🚀 Getting Started

### 🧰 Prerequisites

- Go (latest stable version)
- PostgreSQL (for production)
- SQLite (optional, for local dev)
- Make (optional but recommended)

---

### 📥 Clone the Repository

```bash
git clone https://github.com/puni9869/pinmyblogs.git
cd pinmyblogs
````

---

### ⚙️ Environment Setup

Set environment for local development:

```bash
export ENVIRONMENT=local
```

### ▶️ Run the Application

With hot reload:

```bash
make server
```

This uses **Air**, so changes are reflected instantly during development.

---

## 🧪 Testing & Quality

Run unit tests:

```bash
make test
```

Run linters:

```bash
make lint
```

Security checks:

```bash
make govulncheck
make vet
```

---

## 🗂️ Project Structure

```
.
├── cmd/                 # Application entrypoints
├── frontend/            # Static frontend assets
├── handlers/            # HTTP handlers
├── middleware/          # Security, CSP, headers
├── models/              # Database models
├── pkg/                 # Shared packages (scraping, utils)
├── templates/           # HTML templates
├── types/               # Shared types & forms
└── Makefile             # Build & dev commands
```

---

## 🔐 Security

pinmyblogs uses:

* Strict Content Security Policy (CSP)
* Secure HTTP headers
* Clickjacking protection
* MIME sniffing prevention

External resources are minimized for better security.

---

## 🤝 Contributing

Contributions are welcome! 🎉

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request

Please include tests and follow existing code style.

---

## 🐛 Issues & Feedback

* Found a bug? → Open an issue
* Have an idea? → Feature requests are welcome

👉 [https://github.com/puni9869/pinmyblogs/issues](https://github.com/puni9869/pinmyblogs/issues)

---

## 📜 License

This project is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.

---

## ❤️ Acknowledgements

Built with love using Go, Js and Tailwind-css and open-source tools.

---

⭐ If you find this project useful, please consider starring the repo!
---
