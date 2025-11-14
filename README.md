# Go Link Shortener

A robust URL shortening service built with Go, featuring user authentication, click tracking, and comprehensive statistics. This application provides a RESTful API for creating, managing, and tracking shortened URLs with real-time analytics.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Folder Structure](#folder-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Testing](#testing)
- [TODO List](#todo-list)

## 🔍 Overview

This link shortener service allows users to:
- Register and authenticate using JWT tokens
- Create shortened URLs with automatically generated 6-character hashes
- Track clicks and visits with daily statistics
- View analytics grouped by day or month
- Manage their shortened links (create, update, delete, list)

The application follows clean architecture principles with separation of concerns across handlers, services, repositories, and models.

## ✨ Features

### Authentication & Authorization
- **User Registration**: Create new accounts with email and password
- **Login System**: JWT-based authentication
- **Password Security**: Bcrypt hashing for password storage
- **Protected Routes**: Middleware-based authorization for sensitive endpoints

### Link Management
- **Create Short Links**: Generate unique 6-character hash for any URL
- **Hash Uniqueness**: Automatic collision detection and regeneration
- **Update Links**: Modify existing short links
- **Delete Links**: Soft delete with GORM
- **List Links**: Paginated link listing with offset/limit support
- **URL Validation**: Input validation using go-playground/validator

### URL Redirection
- **Fast Redirects**: Efficient hash-based lookup
- **Event Tracking**: Asynchronous click tracking via EventBus
- **HTTP 307 Redirects**: Temporary redirect status for proper HTTP semantics

### Analytics & Statistics
- **Click Tracking**: Daily click aggregation per link
- **Time-based Analytics**: Query statistics by date range
- **Flexible Grouping**: View data grouped by day or month
- **Real-time Processing**: Event-driven architecture for stat collection

### Infrastructure
- **CORS Support**: Configurable cross-origin resource sharing
- **Request Logging**: Automatic request/response logging middleware
- **JSON API**: Consistent JSON request/response handling
- **Error Handling**: Structured error responses
- **Database Migrations**: Auto-migration support with GORM

## 🏗️ Architecture

The application follows a layered architecture pattern:

```
┌─────────────────────────────────────────┐
│         HTTP Handlers Layer             │
│  (Request validation, Response format)  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Service Layer                   │
│  (Business logic, Authentication)       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Repository Layer                │
│  (Database operations, GORM)            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         PostgreSQL Database             │
└─────────────────────────────────────────┘
```

**Event-Driven Architecture:**
```
Link Visit → Publish Event → EventBus → StatService → Save to DB
```

## 🛠️ Tech Stack

### Core
- **Go 1.25.1**: Primary programming language
- **PostgreSQL 16.4**: Relational database
- **GORM**: ORM for database operations

### Libraries & Frameworks
- **net/http**: Standard library HTTP server
- **golang-jwt/jwt/v5**: JWT token generation and validation
- **go-playground/validator/v10**: Request validation
- **bcrypt**: Password hashing
- **godotenv**: Environment variable management

### DevOps
- **Docker Compose**: PostgreSQL containerization
- **sqlmock**: Database testing

## 📁 Folder Structure

The project follows a clean architecture with organized folders:

```
go-link-shortener/
├── cmd/                    # Application entry points and main executables
├── configs/                # Configuration management and loaders
├── internal/               # Private application modules
│   ├── auth/              # Authentication and authorization
│   ├── link/              # Link shortening core logic
│   ├── stat/              # Statistics and analytics
│   └── user/              # User management
├── pkg/                    # Reusable packages and utilities
│   ├── db/                # Database connection and setup
│   ├── di/                # Dependency injection interfaces
│   ├── event/             # Event bus implementation
│   ├── jwt/               # JWT token utilities
│   ├── middleware/        # HTTP middleware components
│   ├── request/           # Request handling utilities
│   └── response/          # Response formatting utilities
└── migrations/             # Database migration scripts
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.1 or higher
- Docker & Docker Compose
- PostgreSQL 16.4 (via Docker)

### Run Flow

Quick start guide to get the application running:

* Run `docker compose up -d` to run Postgres
* Run `go run migrations/auto.go` to apply migrations
* Run `go run cmd/main.go` to run main web server

The server will start on `http://localhost:8081`

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-link-shortener
```

2. **Create `.env` file**
```bash
cat > .env << EOF
DSN=host=localhost user=postgres password=my_pass dbname=postgres port=5432 sslmode=disable TimeZone=UTC
Secret=your-secret-key-here-change-in-production
EOF
```

3. **Start PostgreSQL**
```bash
docker-compose up -d
```

4. **Run migrations**
```bash
go run migrations/auto.go
```

5. **Install dependencies**
```bash
go mod download
```

6. **Run the application**
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8081`

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/auth/...
go test ./pkg/jwt/...
```

## ⚙️ Configuration

Configuration is loaded from environment variables via `.env` file:

| Variable | Description | Example |
|----------|-------------|---------|
| `DSN` | PostgreSQL connection string | `host=localhost user=postgres password=my_pass dbname=postgres port=5432 sslmode=disable` |
| `Secret` | JWT signing secret | `your-secret-key-min-32-chars` |

## 🧪 Testing

The project includes unit tests for critical components:

- **Auth Handler Tests** (`cmd/auth_test.go`)
- **Auth Service Tests** (`internal/auth/service_test.go`)
- **Auth Handler Tests** (`internal/auth/handler_test.go`)
- **JWT Tests** (`pkg/jwt/jwt_test.go`)

Test coverage includes:
- User registration and login flows
- JWT token creation and validation
- Password hashing and verification
- Error handling scenarios

## 📝 TODO List

### 🚨 Critical Issues & Security

1. **Security: JWT Secret Management**
   - ❌ JWT secret is loaded from env but has no validation
   - ❌ No minimum secret length requirement
   - ⚠️ Need to implement secret rotation mechanism
   - **Priority:** HIGH

2. **Security: Password Salt Missing**
   - ❌ Comment in `auth/service.go:29` says "TODO: add salt"
   - ❌ Using bcrypt default cost without salt configuration
   - **Priority:** HIGH

3. **Security: No Rate Limiting**
   - ❌ No rate limiting on login/register endpoints (brute force vulnerability)
   - ❌ No rate limiting on link creation (spam vulnerability)
   - ❌ No rate limiting on redirect endpoint (DDoS vulnerability)
   - **Priority:** HIGH

4. **Security: CORS Configuration**
   - ❌ CORS allows ANY origin (`Access-Control-Allow-Origin: *` equivalent)
   - ❌ Should restrict to specific domains in production
   - **Priority:** MEDIUM

5. **Security: No HTTPS Enforcement**
   - ❌ Server runs on HTTP only (port 8081)
   - ❌ Need TLS/HTTPS support for production
   - **Priority:** HIGH for production

6. **Security: SQL Injection Prevention**
   - ✅ Using GORM with parameterized queries (good)
   - ⚠️ Need to audit all raw SQL queries in stats repository
   - **Priority:** MEDIUM

### 🐛 Bugs & Code Issues

7. **Bug: Empty JSON Field Tag**
   - ❌ `stat/model.go:12` has invalid JSON tag: `Date   datatypes.Date \`json:"Date\``
   - Missing closing backtick and should be lowercase
   - **Priority:** MEDIUM

8. **Bug: Incorrect HTTP Status Codes**
   - ❌ `request/handle.go:13,20` returns 402 (Payment Required) for validation errors
   - Should be 400 (Bad Request)
   - **Priority:** LOW

9. **Bug: EventBus Channel Not Closed**
   - ❌ EventBus channel is never closed, potential goroutine leak
   - ❌ No graceful shutdown mechanism
   - **Priority:** MEDIUM

10. **Bug: No Context Cancellation**
    - ❌ `main.go:52` starts goroutine without context
    - ❌ Cannot gracefully stop StatService
    - **Priority:** MEDIUM

11. **Bug: Unsafe Random Number Generation**
    - ❌ `link/model.go:35` uses `rand.Intn` without seeding
    - ❌ Not cryptographically secure
    - ❌ May generate predictable hashes
    - **Priority:** HIGH

12. **Bug: File Name Has Space**
    - ❌ `pkg/request/validate .go` has space in filename
    - Should be `validate.go`
    - **Priority:** LOW

13. **Bug: Inconsistent Error Handling**
    - ❌ Some functions return errors, others panic
    - ❌ `db.NewDb` panics on error instead of returning error
    - **Priority:** MEDIUM

14. **Bug: Missing Authorization Header Deletion**
    - ❌ `middleware/auth.go:39` tries to delete header AFTER serving
    - This doesn't prevent the header from being sent
    - **Priority:** LOW

### 📊 Data & Database Issues

15. **Database: No User-Link Relationship**
    - ❌ Links are not associated with users
    - ❌ Anyone can delete/update any link
    - ❌ No ownership tracking
    - **Priority:** HIGH

16. **Database: Missing Indexes**
    - ❌ No index on `stats.link_id`
    - ❌ No index on `stats.date`
    - ❌ No composite index on `(link_id, date)`
    - **Priority:** MEDIUM

17. **Database: No Transaction Support**
    - ❌ Critical operations not wrapped in transactions
    - ❌ Potential data inconsistency issues
    - **Priority:** MEDIUM

18. **Database: Connection Pool Not Configured**
    - ❌ No max connections, idle connections, or timeouts configured
    - ❌ May cause connection exhaustion under load
    - **Priority:** MEDIUM

19. **Database: Migration Strategy Missing**
    - ❌ Using AutoMigrate (not suitable for production)
    - ❌ No versioned migrations
    - ❌ No rollback capability
    - **Priority:** HIGH for production

20. **Data: No Soft Delete for Stats**
    - ❌ Stats use soft delete (deleted_at) but shouldn't
    - ❌ Historical data should never be deleted
    - **Priority:** LOW

### 🔧 Architecture & Code Quality

21. **Architecture: Tight Coupling**
    - ❌ Handlers directly instantiate dependencies
    - ⚠️ Partial DI interfaces in `pkg/di` but incomplete
    - ❌ No proper dependency injection container
    - **Priority:** MEDIUM

22. **Architecture: No Service Interface**
    - ❌ Services don't implement interfaces
    - ❌ Difficult to test and mock
    - **Priority:** MEDIUM

23. **Architecture: EventBus Lacks Features**
    - ❌ No multiple subscribers support
    - ❌ No event filtering/routing
    - ❌ No error handling for failed events
    - ❌ No event replay capability
    - **Priority:** MEDIUM

24. **Code Quality: Debug Print Statements**
    - ❌ Multiple `fmt.Println` debug statements throughout code
    - Should use proper structured logging
    - **Priority:** LOW

25. **Code Quality: Magic Numbers**
    - ❌ Hash length hardcoded as 6 in multiple places
    - ❌ No constants defined
    - **Priority:** LOW

26. **Code Quality: Missing Documentation**
    - ❌ No godoc comments on exported functions
    - ❌ No package documentation
    - **Priority:** LOW

27. **Code Quality: Inconsistent Naming**
    - ❌ `configs` vs `config` package naming
    - ❌ `hanshedPwd` typo in `auth/service.go:29`
    - **Priority:** LOW

### 🎯 Features & Improvements

28. **Feature: Custom Short Links**
    - ⭐ Allow users to specify custom hash instead of random
    - ⭐ Need validation for custom hash (length, characters, reserved words)
    - **Priority:** MEDIUM

29. **Feature: Link Expiration**
    - ⭐ Add expiration date for short links
    - ⭐ Background job to clean expired links
    - **Priority:** MEDIUM

30. **Feature: Link Password Protection**
    - ⭐ Optional password for accessing links
    - **Priority:** LOW

31. **Feature: QR Code Generation**
    - ⭐ Generate QR codes for shortened URLs
    - **Priority:** LOW

32. **Feature: Link Preview/Metadata**
    - ⭐ Store page title, description, favicon
    - ⭐ Show preview before redirect (optional)
    - **Priority:** LOW

33. **Feature: Link Categories/Tags**
    - ⭐ Allow users to organize links
    - **Priority:** LOW

34. **Feature: Analytics Enhancement**
    - ⭐ Track referrer, user agent, geolocation
    - ⭐ Track unique vs total clicks
    - ⭐ Export statistics as CSV/JSON
    - **Priority:** MEDIUM

35. **Feature: API Versioning**
    - ⭐ Implement `/v1/` prefix for API routes
    - ⭐ Plan for backwards compatibility
    - **Priority:** MEDIUM

### 🧪 Testing & Quality Assurance

36. **Testing: Low Test Coverage**
    - ❌ Only auth and JWT packages have tests
    - ❌ No tests for link, stat, user modules
    - ❌ No integration tests
    - ❌ No end-to-end tests
    - **Priority:** HIGH

37. **Testing: No Benchmark Tests**
    - ❌ No performance benchmarks
    - ❌ Cannot measure optimization impact
    - **Priority:** LOW

38. **Testing: No Test Database Setup**
    - ❌ Tests might use production database
    - ❌ Need separate test database configuration
    - **Priority:** MEDIUM

### 📖 Documentation & DevOps

39. **Documentation: Missing API Documentation**
    - ❌ No Swagger/OpenAPI specification
    - ❌ No Postman collection
    - **Priority:** MEDIUM

40. **Documentation: No Architecture Diagrams**
    - ⭐ Would benefit from sequence diagrams
    - ⭐ Database ERD diagram
    - **Priority:** LOW

41. **DevOps: No CI/CD Pipeline**
    - ❌ No automated testing on commits
    - ❌ No automated deployments
    - **Priority:** MEDIUM

42. **DevOps: No Health Check Endpoint**
    - ❌ No `/health` or `/ping` endpoint
    - ❌ Cannot monitor service status
    - **Priority:** MEDIUM

43. **DevOps: No Metrics/Monitoring**
    - ❌ No Prometheus metrics
    - ❌ No application performance monitoring
    - ❌ No error tracking (Sentry, etc.)
    - **Priority:** MEDIUM

44. **DevOps: No Docker Image for App**
    - ❌ Only postgres in docker-compose
    - ❌ Application not containerized
    - **Priority:** MEDIUM

45. **DevOps: No Environment Separation**
    - ❌ No dev/staging/production configs
    - ❌ Single .env file for all environments
    - **Priority:** MEDIUM

46. **DevOps: No Logging Strategy**
    - ❌ Logs go to stdout with basic fmt.Println
    - ❌ No structured logging (JSON format)
    - ❌ No log levels (debug, info, warn, error)
    - ❌ No log aggregation setup
    - **Priority:** MEDIUM

### 🔄 Operational Issues

47. **Operations: No Graceful Shutdown**
    - ❌ Server doesn't handle SIGTERM/SIGINT
    - ❌ In-flight requests may be dropped
    - ❌ EventBus goroutine won't stop cleanly
    - **Priority:** HIGH

48. **Operations: No Request ID Tracking**
    - ❌ Cannot trace requests through logs
    - ❌ Difficult to debug issues
    - **Priority:** MEDIUM

49. **Operations: No Backup Strategy**
    - ❌ No database backup documentation
    - ❌ No disaster recovery plan
    - **Priority:** HIGH for production

50. **Operations: Hard-coded Port**
    - ❌ Server port 8081 is hard-coded in `main.go:62`
    - Should be configurable via environment variable
    - **Priority:** LOW

### 🎨 API & Usability

51. **API: Inconsistent Response Format**
    - ⚠️ Some endpoints return data, others return null
    - ⚠️ Error responses not standardized
    - **Priority:** MEDIUM

52. **API: No Pagination Standards**
    - ❌ GetList uses limit/offset without defaults
    - ❌ No total count in pagination response
    - ❌ No "next page" URL
    - **Priority:** LOW

53. **API: No Request Validation Messages**
    - ❌ Validation errors return raw validator output
    - ❌ Not user-friendly
    - **Priority:** LOW

54. **API: Missing PATCH vs PUT Semantics**
    - ❌ PATCH endpoint requires all fields
    - Should support partial updates
    - **Priority:** LOW

55. **Usability: No Link Click Cooldown**
    - ❌ Same user can inflate stats by repeated clicks
    - ⭐ Consider tracking unique visitors (IP, session)
    - **Priority:** MEDIUM

## 📄 License

This project is open source and available for educational purposes.

## 👤 Author

Built as a Go playground project for learning and demonstration purposes.

---

**Note:** This is a development/learning project. Before deploying to production, address the critical security issues listed in the TODO section, particularly:
- JWT secret management
- Rate limiting
- HTTPS/TLS support
- User-link ownership
- Comprehensive testing
- Proper migration strategy
