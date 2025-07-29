# Private Blockchain Corporate Data Management System

A secure corporate data management system utilizing private blockchain technology. It prevents data tampering, ensures transparency, implements fine-grained access controls, and securely manages critical corporate data.

## Features

- **Private Blockchain**: All data operations are recorded on the blockchain to prevent tampering
- **Role-Based Access Control (RBAC)**: Four-tier permission management (Admin, Manager, Employee, Guest)
- **Data Encryption**: AES-256 encryption at rest and TLS communication
- **Audit Logging**: All operations are recorded in detail for compliance requirements
- **Version Control**: Automatic tracking of document change history
- **Fine-Grained Permissions**: Detailed access control at file/folder level

## Technology Stack

### Backend
- **Go 1.21+** - Main application language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Primary database
- **Redis** - Cache and session management
- **JWT** - Authentication tokens

### Infrastructure
- **Docker & Docker Compose** - Containerization
- **Nginx** - Reverse proxy and load balancer

## Project Structure

```
├── cmd/                    # Application entry points
│   └── server/            # Main server
├── internal/              # Internal packages
│   ├── api/              # API related
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Middleware
│   │   └── routes/       # Routing
│   ├── blockchain/       # Blockchain implementation
│   ├── config/           # Configuration management
│   ├── database/         # Database related
│   │   └── models/       # Data models
│   ├── security/         # Security features
│   │   ├── auth/         # Authentication
│   │   ├── crypto/       # Encryption
│   │   └── rbac/         # Access control
│   └── services/         # Business logic
├── deployments/          # Deployment configuration
├── docs/                 # Documentation
└── tests/                # Tests
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (if not using Docker)

### 1. Clone Repository

```bash
git clone https://github.com/nshmdayo/in-house-datamanagement-system-sample.git
cd in-house-datamanagement-system-sample
```

### 2. Setup Development Environment

```bash
make setup-dev
```

### 3. Configure Environment Variables

```bash
cp .env.example .env
# Edit .env file to adjust settings
```

### 4. Install Dependencies

```bash
make deps
```

### 5. Start with Docker Compose

```bash
make docker-compose-up
```

### 6. Start for Local Development

```bash
# Ensure database is running
make dev  # Start development server with hot reload
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/profile` - Get profile

### User Management
- `GET /api/v1/users` - Get user list
- `POST /api/v1/users` - Create user (Admin only)
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (Admin only)

### Document Management
- `GET /api/v1/documents` - Get document list
- `POST /api/v1/documents` - Create document
- `GET /api/v1/documents/:id` - Get document details
- `PUT /api/v1/documents/:id` - Update document
- `DELETE /api/v1/documents/:id` - Delete document

### Blockchain
- `GET /api/v1/blockchain/blocks` - Get block list
- `GET /api/v1/blockchain/transactions/:id` - Get transaction details
- `POST /api/v1/blockchain/verify` - Data integrity verification

### Audit Logs
- `GET /api/v1/audit/logs` - Get audit logs (Admin/Manager only)
- `GET /api/v1/audit/statistics` - Get statistics

## Development Commands

```bash
# Start development server
make dev

# Build
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint
make lint

# Format code
make fmt

# Security scan
make security

# Build Docker image
make docker-build

# Database migration
make db-migrate

# Database seed
make db-seed
```

## Security Features

### Authentication & Authorization
- JWT-based authentication
- Role-Based Access Control (RBAC)
- Account lockout on failed login attempts
- Session management

### Data Protection
- AES-256 encryption at rest
- TLS 1.3 encryption in transit
- bcrypt password hashing

### Audit & Logging
- Detailed logging of all operations
- Security event tracking
- IP address and User Agent recording

## Blockchain Features

### Characteristics
- Private blockchain implementation
- Proof of Work consensus
- Merkle Tree for efficient verification
- Automatic data integrity verification

### Recorded Operations
- Document creation, updates, and deletion
- Access permission changes
- User operation history

## Testing

```bash
# Run all tests
make test

# Generate coverage report
make test-coverage

# Benchmark tests
make bench

# Load testing
make load-test
```

## Deployment

### Docker Compose (Recommended)

```bash
# Start in production
make docker-compose-up

# View logs
make docker-compose-logs

# Stop
make docker-compose-down
```

### Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/k8s/
```

## Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### Metrics

Application metrics are exposed at the `/metrics` endpoint.

## Security Considerations

- Regular dependency updates
- Security patch application
- Regular audit log review
- Regular backup execution

## License

This project is released under the MIT License.

## Contribution

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add some amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## Support

If you have questions or issues, please let us know through GitHub Issues.