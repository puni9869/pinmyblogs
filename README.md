# 📌 pinmyblogs — Open-Source Blog Bookmarking & Read-Later App (Go)

First commit
**pinmyblogs** is an open-source **blog bookmarking**, **read-later**, and **content organization** application built
with **Go (Golang)**.

Save blog links, extract metadata, organize reading lists, and revisit content distraction-free — fast, secure, and
self-hosted.

> Ideal for developers, writers, and knowledge workers who want a simple, privacy-friendly alternative to hosted
> bookmark tools.

---

## 📊 Project Status & Badges

[![Check](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml/badge.svg)](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/puni9869/pinmyblogs)
![License](https://img.shields.io/github/license/puni9869/pinmyblogs)
![Stars](https://img.shields.io/github/stars/puni9869/pinmyblogs?style=social)

---

## 🔍 What is pinmyblogs?

**pinmyblogs** is a lightweight **self-hosted bookmark manager** focused on:

- Saving blog URLs
- Extracting page metadata (title, favicon, etc.)
- Organizing blogs for later reading
- Running reliably in local or production environments

It is built with **performance, simplicity, and security** in mind.

---

## ✨ Features

- 🔖 Save blog & article URLs
- 🗂️ Organize reading lists
- 🕷️ Automatic metadata scraping
- ⚡ High-performance Go backend
- 🧩 Clean, extensible architecture
- 🏠 Self-hosted & privacy-friendly

---

## 🧱 Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Gin
- **Database:** PostgreSQL or SQLite
- **Frontend:** HTML templates + Tailwind CSS
- **Build Tooling:**  + Air (hot reload)

---

## 🗄️ Database Support

pinmyblogs supports **multiple SQL databases**:

### ✅ PostgreSQL (Recommended for Production)

- High concurrency
- Strong data integrity
- Crash-safe & scalable

### ✅ SQLite (Development & Testing)

- Zero-config
- File-based
- Best for local usage

```text
Production  → PostgreSQL
Development → SQLite
````

---

## 🚀 Getting Started

### 🧰 Prerequisites

* Go (latest stable version)
* PostgreSQL (optional, for production)
* SQLite (optional, for development)
* Make (recommended)

---

### 📥 Installation

```bash
git clone https://github.com/puni9869/pinmyblogs.git
cd pinmyblogs
```

---

### ⚙️ Environment Configuration

### ▶️ Run the Application

```bash
make server
```

Uses **Air** for automatic reload on code changes.

---

## 🧪 Testing & Quality Checks

Run tests:

```bash
make test
```

Lint & static analysis:

```bash
make lint
make vet
make govulncheck
```

---

## 🗂️ Project Structure

```
.
├── cmd/                 # Application entrypoints
├── handlers/            # HTTP handlers
├── middleware/          # CSP & security headers
├── models/              # Database models
├── pkg/                 # Shared packages (scraping, utils)
├── templates/           # HTML templates
├── frontend/            # Static assets
├── types/               # Forms & shared types
└── Makefile             # Dev & build commands
```

---

## 🔐 Security

pinmyblogs includes production-grade security defaults:

* Strict Content Security Policy (CSP)
* Secure HTTP headers
* Clickjacking protection
* MIME-type sniffing prevention

Designed to be safe by default.

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Add tests where applicable
4. Submit a pull request

See issues for ideas and improvements.

---

## 🐛 Issues & Feature Requests

Found a bug or have a feature idea?

👉 [https://github.com/puni9869/pinmyblogs/issues](https://github.com/puni9869/pinmyblogs/issues)

---

## 📜 License

This project is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.

---

## ⭐ Why pinmyblogs?

* Open-source
* Self-hosted
* Developer-friendly
* Privacy-focused
* Written in Go

If you find this project useful, please ⭐ star the repository — it helps others discover it!

