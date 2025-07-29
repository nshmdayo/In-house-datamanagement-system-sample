# Private Blockchain Corporate Data Management System - Development Guidelines

## Project Overview

This project is a corporate data management system that leverages private blockchain technology. It ensures data integrity prevention, transparency, and access control to securely manage critical corporate data.

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Blockchain**: Hyperledger Fabric or custom implementation
- **Database**: PostgreSQL
- **API**: REST API (Gin Framework)
- **Authentication**: JWT + RBAC (Role-Based Access Control)
- **Encryption**: AES-256, RSA

### Frontend (Future Implementation)
- **Framework**: React/Next.js
- **State Management**: Redux Toolkit
- **UI**: Material-UI

## System Requirements

### Functional Requirements

#### 1. User Management
- User registration and authentication (admin approval required)
- Role-based access control (Admin, Manager, Employee, Guest)
- Profile management
- Login history tracking

#### 2. Data Management
- Document registration, update, and deletion
- Version control
- Metadata management (creator, modification date, category, etc.)
- Encrypted file storage

#### 3. Blockchain Features
- Record all data operations on blockchain
- Tamper detection functionality
- Data integrity verification
- Transaction tracking

#### 4. Access Control
- File/folder-level access permission settings
- Department-based data access control
- Access restrictions based on confidentiality levels

#### 5. Audit Features
- Recording of all operation logs
- Access log visualization
- Compliance report generation

#### 6. Search and Filtering
- Full-text search
- Metadata search
- Advanced filtering capabilities

### Non-Functional Requirements

#### 1. Security
- Data encryption (at rest and in transit)
- Multi-factor authentication (MFA)
- Security log monitoring
- Regular security audits

#### 2. Performance
- API response time: 95%+ of requests within 500ms
- Concurrent connections: 1000 users
- Database response time: within 100ms

#### 3. Availability
- System uptime: 99.9%+
- Automated backups (daily and weekly)
- Disaster recovery plan

#### 4. Scalability
- Horizontal scaling support
- Microservices architecture
- Containerization (Docker)

## Architecture Design

### System Configuration

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Backend       │
│   (React)       │◄──►│   (Nginx)       │◄──►│   (Go)          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                         │
                         ┌─────────────────┐            │
                         │  Blockchain     │◄───────────┤
                         │  Network        │            │
                         └─────────────────┘            │
                                                         │
                         ┌─────────────────┐            │
                         │  PostgreSQL     │◄───────────┘
                         │  Database       │
                         └─────────────────┘
```

### Directory Structure

```
├── cmd/                    # Application entry points
│   └── server/
├── internal/               # Internal packages
│   ├── api/               # API related
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware
│   │   └── routes/        # Routing
│   ├── blockchain/        # Blockchain related
│   │   ├── block/         # Block structure
│   │   ├── chain/         # Chain management
│   │   ├── consensus/     # Consensus
│   │   └── transaction/   # Transaction
│   ├── config/            # Configuration management
│   ├── database/          # Database related
│   │   ├── migrations/    # Migrations
│   │   └── models/        # Data models
│   ├── security/          # Security related
│   │   ├── auth/          # Authentication
│   │   ├── crypto/        # Encryption
│   │   └── rbac/          # Role-Based Access Control
│   └── services/          # Business logic
├── pkg/                   # External packages
├── scripts/               # Scripts
├── deployments/           # Deployment configuration
├── docs/                  # Documentation
└── tests/                 # Tests
```

## Data Models

### Main Entities

#### User
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Role      Role      `json:"role"`
    IsActive  bool      `json:"is_active" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### Document
```go
type Document struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Content     string    `json:"content"`
    FileHash    string    `json:"file_hash" gorm:"unique"`
    Category    string    `json:"category"`
    AccessLevel int       `json:"access_level"`
    CreatedBy   uint      `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### BlockchainRecord
```go
type BlockchainRecord struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    TransactionID string    `json:"transaction_id" gorm:"unique"`
    BlockHash     string    `json:"block_hash"`
    DocumentID    uint      `json:"document_id"`
    Action        string    `json:"action"`
    UserID        uint      `json:"user_id"`
    Timestamp     time.Time `json:"timestamp"`
}
```

## API Design

### Authentication Endpoints
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/logout` - Logout

### User Management
- `GET /api/v1/users` - Get user list
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user details
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Document Management
- `GET /api/v1/documents` - Get document list
- `POST /api/v1/documents` - Create document
- `GET /api/v1/documents/{id}` - Get document details
- `PUT /api/v1/documents/{id}` - Update document
- `DELETE /api/v1/documents/{id}` - Delete document

### Blockchain
- `GET /api/v1/blockchain/blocks` - Get block list
- `GET /api/v1/blockchain/transactions/{id}` - Get transaction details
- `POST /api/v1/blockchain/verify` - Data integrity verification

## Security Guidelines

### Encryption
- Database sensitive data encrypted with AES-256
- Passwords hashed with bcrypt
- API communication uses TLS 1.3

### Authentication & Authorization
- JWT token expiry: 15 minutes
- Refresh token expiry: 7 days
- Fine-grained access control with RBAC

### Input Validation
- Validation of all input data
- SQL injection prevention
- XSS protection

## Testing Strategy

### Unit Tests
- Coverage 80%+
- Test isolation using mocks

### Integration Tests
- API endpoint testing
- Database integration testing

### Security Tests
- Vulnerability scanning
- Penetration testing

## Development Standards

### Code Style
- Code formatting with gofmt
- Linting with golint
- Naming conventions follow Go standards

### Commit Conventions
- Adopt Conventional Commits format
- feat: New features
- fix: Bug fixes
- docs: Documentation
- test: Tests

### Package Management
- Use Go Modules
- Regular dependency updates

## Deployment

### Environments
- Development
- Staging  
- Production

### CI/CD
- Use GitHub Actions
- Automated test execution
- Docker image building
- Automated deployment to K8s

---

Please follow these guidelines to build a secure and reliable corporate data management system.
