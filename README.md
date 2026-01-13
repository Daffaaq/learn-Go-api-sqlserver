# Go SQLServer CRUD API

RESTful API sederhana untuk manajemen **User** menggunakan **Go**, **SQL Server**, dan **Chi Router**.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Pagination](#pagination)
- [Validation](#validation)
- [Graceful Shutdown](#graceful-shutdown)
- [License](#license)

---

## Features

- CRUD User (`Create`, `Read`, `Update`, `Delete`)
- Pagination support (`limit` & `offset`)
- Email uniqueness validation
- Input validation using [go-playground/validator](https://github.com/go-playground/validator)
- Graceful server shutdown
- JSON consistent response with error handling
- Configurable via environment variables

---

## Tech Stack

- **Go** (>= 1.21)
- **SQL Server**
- **Chi Router** for routing
- **go-playground/validator** for struct validation

---
