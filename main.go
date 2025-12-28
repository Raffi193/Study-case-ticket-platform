// main.go
package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Models
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	EventDate        time.Time `json:"event_date"`
	TotalCapacity    int       `json:"total_capacity"`
	AvailableTickets int       `json:"available_tickets"`
	Price            float64   `json:"price"`
	OrganizerID      int       `json:"organizer_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type Booking struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	EventID       int       `json:"event_id"`
	TicketCount   int       `json:"ticket_count"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"` // pending, confirmed, cancelled
	BookingDate   time.Time `json:"booking_date"`
	PaymentStatus string    `json:"payment_status"` // unpaid, paid
}

// Initialize Database
func initDB() {
    var err error
    dsn := "root:@tcp(127.0.0.1:3306)/ticket-platform?parseTime=true"                                              
    
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    
    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database:", err)
    }
    
    log.Println("✅ Database connected successfully!")
}

// ========== USER ENDPOINTS ==========

// Create User
func createUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (name, email, phone) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Phone,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// Get all users
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email, phone, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Get user by id
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User

	err := db.QueryRow(
		"SELECT id, name, email, phone, created_at FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Update User
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(
		"UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?",
		user.Name, user.Email, user.Phone, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Delete user
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ========== EVENT ENDPOINTS ==========

// Create Event
func createEvent(c *gin.Context) {
	var event Event
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set available tickets sama dengan total capacity
	event.AvailableTickets = event.TotalCapacity

	result, err := db.Exec(
		`INSERT INTO events (name, description, location, event_date, total_capacity, 
		available_tickets, price, organizer_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Name, event.Description, event.Location, event.EventDate,
		event.TotalCapacity, event.AvailableTickets, event.Price, event.OrganizerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	event.ID = int(id)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"data":    event,
	})
}

// Get All Events
func getEvents(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, name, description, location, event_date, total_capacity, 
		available_tickets, price, organizer_id, created_at FROM events
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.Location,
			&event.EventDate, &event.TotalCapacity, &event.AvailableTickets,
			&event.Price, &event.OrganizerID, &event.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		events = append(events, event)
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

// Get Event by ID
func getEventByID(c *gin.Context) {
	id := c.Param("id")
	var event Event

	err := db.QueryRow(`
		SELECT id, name, description, location, event_date, total_capacity, 
		available_tickets, price, organizer_id, created_at 
		FROM events WHERE id = ?`, id,
	).Scan(
		&event.ID, &event.Name, &event.Description, &event.Location,
		&event.EventDate, &event.TotalCapacity, &event.AvailableTickets,
		&event.Price, &event.OrganizerID, &event.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}

// Update Event
func updateEvent(c *gin.Context) {
	id := c.Param("id")
	var event Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
		UPDATE events SET name = ?, description = ?, location = ?, 
		event_date = ?, total_capacity = ?, price = ? WHERE id = ?`,
		event.Name, event.Description, event.Location,
		event.EventDate, event.TotalCapacity, event.Price, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// Delete Event
func deleteEvent(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM events WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// ========== BOOKING ENDPOINTS (ANTI OVERSELLING) ==========

// Create Booking with Transaction to Prevent Overselling
func createBooking(c *gin.Context) {
	var booking Booking
	if err := c.BindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start Transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Lock row untuk mencegah race condition (FOR UPDATE)
	var availableTickets int
	var price float64
	err = tx.QueryRow(`
		SELECT available_tickets, price FROM events 
		WHERE id = ? FOR UPDATE`, booking.EventID,
	).Scan(&availableTickets, &price)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if enough tickets available
	if availableTickets < booking.TicketCount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "Not enough tickets available",
			"available_tickets": availableTickets,
			"requested":         booking.TicketCount,
		})
		return
	}

	// Calculate total harga
	booking.TotalPrice = float64(booking.TicketCount) * price
	booking.Status = "pending"
	booking.PaymentStatus = "unpaid"

	// Insert booking
	result, err := tx.Exec(`
		INSERT INTO bookings (user_id, event_id, ticket_count, total_price, status, payment_status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		booking.UserID, booking.EventID, booking.TicketCount,
		booking.TotalPrice, booking.Status, booking.PaymentStatus,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingID, _ := result.LastInsertId()
	booking.ID = int(bookingID)

	// Update available tickets
	_, err = tx.Exec(`
		UPDATE events SET available_tickets = available_tickets - ? 
		WHERE id = ?`, booking.TicketCount, booking.EventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"data":    booking,
	})
}

// Get All Bookings
func getBookings(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, user_id, event_id, ticket_count, total_price, 
		status, payment_status, booking_date FROM bookings
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.EventID,
			&booking.TicketCount, &booking.TotalPrice, &booking.Status,
			&booking.PaymentStatus, &booking.BookingDate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// Get Booking by ID
func getBookingByID(c *gin.Context) {
	id := c.Param("id")
	var booking Booking

	err := db.QueryRow(`
		SELECT id, user_id, event_id, ticket_count, total_price, 
		status, payment_status, booking_date FROM bookings WHERE id = ?`, id,
	).Scan(
		&booking.ID, &booking.UserID, &booking.EventID,
		&booking.TicketCount, &booking.TotalPrice, &booking.Status,
		&booking.PaymentStatus, &booking.BookingDate,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

// Cancel Booking (Return tickets to available pool)
func cancelBooking(c *gin.Context) {
	id := c.Param("id")

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Get booking details
	var eventID, ticketCount int
	var status string
	err = tx.QueryRow(`
		SELECT event_id, ticket_count, status FROM bookings 
		WHERE id = ? FOR UPDATE`, id,
	).Scan(&eventID, &ticketCount, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking already cancelled"})
		return
	}

	// Update booking status
	_, err = tx.Exec("UPDATE bookings SET status = 'cancelled' WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return tickets to event
	_, err = tx.Exec(`
		UPDATE events SET available_tickets = available_tickets + ? 
		WHERE id = ?`, ticketCount, eventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

// Confirm Payment
func confirmPayment(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec(`
		UPDATE bookings 
		SET payment_status = 'paid', status = 'confirmed' 
		WHERE id = ?`, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment confirmed successfully"})
}

// Main Function
func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Ticket Platform API!",
			"status":  "running",
		})
	})

	// User Routes
	router.POST("/api/users", createUser)
	router.GET("/api/users", getUsers)
	router.GET("/api/users/:id", getUserByID)
	router.PUT("/api/users/:id", updateUser)
	router.DELETE("/api/users/:id", deleteUser)

	// Event Routes
	router.POST("/api/events", createEvent)
	router.GET("/api/events", getEvents)
	router.GET("/api/events/:id", getEventByID)
	router.PUT("/api/events/:id", updateEvent)
	router.DELETE("/api/events/:id", deleteEvent)

	// Booking Routes
	router.POST("/api/bookings", createBooking)
	router.GET("/api/bookings", getBookings)
	router.GET("/api/bookings/:id", getBookingByID)
	router.PUT("/api/bookings/:id/cancel", cancelBooking)
	router.PUT("/api/bookings/:id/confirm-payment", confirmPayment)

	log.Println("Server running on port 8080...")
	router.Run(":8080")
}
