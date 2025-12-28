# Platform Tiket Event Online

Backend system untuk platform pemesanan tiket event dengan mekanisme anti-overselling.

 Features
- ✅ CRUD Users, Events, Bookings
- ✅ Anti-overselling mechanism dengan database transaction
- ✅ Row locking untuk prevent race condition
- ✅ RESTful API dengan Gin Framework
- ✅ MySQL database dengan foreign key constraints

 Tech Stack
- **Language:** Go (Golang) 1.20+
- **Framework:** Gin Web Framework
- **Database:** MySQL 8.0+
- **Driver:** go-sql-driver/mysql

 Installation
 Prerequisites
- Go 1.20 or higher
- MySQL 8.0 or higher
- Git

Setup
1. Clone repository:
```bash
git clone https://github.com/Raffi193/ticket-platform.git
cd ticket-platform
```

2. Install dependencies:
```bash
go mod download
```

3. Setup database:
```bash
mysql -u root -p < schema.sql
```

4. Configure environment:
```bash
cp .env.example .env
# Edit .env dengan kredensial database Anda
```

5. Run application:
```bash
go run main.go
```

Server akan running di `http://localhost:8080`

API Endpoints
### Users
- `POST /api/users` - Create user
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Events
- `POST /api/events` - Create event
- `GET /api/events` - Get all events
- `GET /api/events/:id` - Get event by ID
- `PUT /api/events/:id` - Update event
- `DELETE /api/events/:id` - Delete event

### Bookings
- `POST /api/bookings` - Create booking (anti-overselling)
- `GET /api/bookings` - Get all bookings
- `GET /api/bookings/:id` - Get booking by ID
- `PUT /api/bookings/:id/cancel` - Cancel booking
- `PUT /api/bookings/:id/confirm-payment` - Confirm payment

Anti-Overselling Mechanism
Sistem ini menggunakan:
1. **Database Transaction** (BEGIN, COMMIT, ROLLBACK)
2. **Row Locking** (`SELECT ... FOR UPDATE`)
3. **Pessimistic Locking Strategy**
4. **Real-time Validation**

Documentation
Dokumentasi lengkap tersedia di folder `docs/`:
- ERD (Entity Relationship Diagram)
- Flowchart (Booking Process)
- API Documentation

Testing
Test API menggunakan Postman:
1. Import collection dari `postman/Ticket_Platform.postman_collection.json`
2. Setup environment variables
3. Run tests

Developer
- **Name:** M. Rafi ash shiddiqie 
- **Email:** m.raffi1808@gmail.com
- **GitHub:** [@username](https://github.com/Raffi193)

License
MIT License
---
**Developed for GDGoC UNSRI Backend Development Selection 2025**
