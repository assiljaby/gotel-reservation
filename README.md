# 🏨 Hotel Reservation JSON API Server

This project is a RESTful API server built for managing hotel reservations, leveraging Go for backend services and MongoDB for data storage. The system is designed to handle the entire reservation lifecycle, from room availability checks to booking confirmations.

## Key Features
- **Go Backend:** High-performance backend built with Go for handling reservation requests efficiently.
- **MongoDB:** Flexible and scalable NoSQL database for storing room, customer, and reservation data.
- **RESTful API:** JSON-based API following REST principles for easy integration with frontend applications.
- **Room Availability and Booking:** Endpoints for checking room availability and making reservations.
- **Customer Management:** Manage customer data, including registration and history of bookings.

## Getting Started
1. **Clone the repository:**
   ```bash
   git clone https://github.com/assiljaby/gotel-reservation
   ```
2. **Install dependencies:**
   ```bash
   go mod tidy
   ```
3. **Start MongoDB as a container:**
   ```bash
   docker run --name mongodb -d -p 27017:27017 mongodb/mongodb-community-server:6.0-ubi8
   ```
4. **Seed db:**
   ```bash
   make seed
   ```
5. **Test units:**
   ```bash
   make test
   ```
5. **Run the API server:**
   ```bash
   make run
   ```
   
## Example API Usage

### Check Bookings
To check the all the booking, you will need an admin token, which you can get from the seed script:

**Request:**
```bash
GET /api/v1/admin/bookings
```
**Response:**
```json
[
    {
        "id": "66dad1a5daf626e5a50ff6e5",
        "userID": "66dad1a5daf626e5a50ff6e1",
        "roomID": "66dad1a5daf626e5a50ff6e4",
        "fromDate": "2024-09-06T09:55:49.967Z",
        "tillDate": "2024-09-11T09:55:49.967Z",
        "canceled": false
    }
]
```

## TODO

- [x] Init
- [x] Set Makefile
- [x] Initialize DB
- [x] Set Users CRUD
- [x] Set User Validation
- [x] Test User API
- [x] Set Hotel API
- [x] Set JWT Auth
- [x] Set Booking API
- [x] Set Booking Validiation
- [x] Admin Auth
- [x] Set Error Handling Middleware
- [x] Clean Up