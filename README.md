# Restaurant Management API

A comprehensive REST API for restaurant management built with Go, Gin framework, and MongoDB. This system provides complete CRUD operations for managing restaurants, including user authentication, menu management, order processing, and billing.

## 🚀 Features

- **User Authentication & Authorization**: JWT-based authentication system
- **Menu Management**: Create, read, update, and delete menu items and food categories
- **Order Management**: Handle customer orders and order items
- **Table Management**: Manage restaurant table assignments
- **Invoice Generation**: Automated billing and invoice generation
- **Database Integration**: MongoDB with Docker support
- **RESTful API**: Clean, well-structured REST endpoints
- **Middleware**: Authentication middleware for protected routes

## 🛠️ Tech Stack

- **Backend**: Go 1.23.4
- **Web Framework**: Gin (v1.11.0)
- **Database**: MongoDB 7.0
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator/v10
- **Containerization**: Docker & Docker Compose



## 🔧 Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/parthavpovil/restraunt-management.git
cd restraunt-management
```

### 2. Set Up Environment Variables
Create a `.env` file in the root directory:
```env
PORT=8000
SECRET_KEY=your-super-secret-jwt-key
MONGODB_URL=mongodb://admin:password123@localhost:27017/restraunt?authSource=admin
```

### 3. Start MongoDB with Docker
```bash
docker-compose up -d
```

### 4. Install Go Dependencies
```bash
go mod download
```

### 5. Run the Application
```bash
go run main.go
```

The API will be available at `http://localhost:8000`

## 📚 API Documentation

### Authentication

#### Sign Up
```http
POST /users/signup
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "securepassword",
  "phone": "+1234567890",
  "avatar": "https://example.com/avatar.jpg"
}
```

#### Login
```http
POST /users/login
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "uid": "user_id_here"
}
```

### User Management

#### Get All Users (Protected)
```http
GET /users
Authorization: Bearer <access_token>
```

#### Get User by ID (Protected)
```http
GET /users/:user_id
Authorization: Bearer <access_token>
```

### Food Management

#### Get All Foods
```http
GET /foods
```

#### Get Food by ID
```http
GET /foods/:food_id
```

#### Create Food Item (Protected)
```http
POST /foods
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Margherita Pizza",
  "price": 15.99,
  "food_image": "https://example.com/pizza.jpg",
  "menu_id": "menu_id_here"
}
```

#### Update Food Item (Protected)
```http
PATCH /foods/:food_id
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Updated Pizza Name",
  "price": 17.99
}
```

### Menu Management

#### Get All Menus
```http
GET /menus
```

#### Get Menu by ID
```http
GET /menus/:menu_id
```

#### Create Menu (Protected)
```http
POST /menus
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Dinner Menu",
  "category": "Main Course",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z"
}
```

#### Update Menu (Protected)
```http
PATCH /menus/:menu_id
Authorization: Bearer <access_token>
```

### Table Management

#### Get All Tables
```http
GET /tables
```

#### Get Table by ID
```http
GET /tables/:table_id
```

#### Create Table (Protected)
```http
POST /tables
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "number_of_guests": 4,
  "table_number": 12
}
```

#### Update Table (Protected)
```http
PATCH /tables/:table_id
Authorization: Bearer <access_token>
```

### Order Management

#### Get All Orders
```http
GET /orders
```

#### Get Order by ID
```http
GET /orders/:order_id
```

#### Create Order (Protected)
```http
POST /orders
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "order_date": "2024-01-15T12:00:00Z",
  "table_id": "table_id_here"
}
```

#### Update Order (Protected)
```http
PATCH /orders/:order_id
Authorization: Bearer <access_token>
```

### Order Items Management

#### Get All Order Items
```http
GET /orderItems
```

#### Get Order Items by Order ID
```http
GET /orderItems-order/:order_id
```

#### Get Order Item by ID
```http
GET /orderItems/:orderItem_id
```

#### Create Order Item (Protected)
```http
POST /orderItems
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "quantity": 2,
  "unit_price": 15.99,
  "food_id": "food_id_here",
  "order_id": "order_id_here"
}
```

#### Update Order Item (Protected)
```http
PATCH /orderItems/:orderItem_id
Authorization: Bearer <access_token>
```

### Invoice Management

#### Get All Invoices
```http
GET /invoices
```

#### Get Invoice by ID
```http
GET /invoices/:invoice_id
```

#### Create Invoice (Protected)
```http
POST /invoices
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "invoice_status": "PENDING",
  "payment_method": "CARD",
  "payment_status": "PENDING",
  "payment_due_date": "2024-01-20T23:59:59Z",
  "order_id": "order_id_here"
}
```

#### Update Invoice (Protected)
```http
PATCH /invoices/:invoice_id
Authorization: Bearer <access_token>
```

## 🔐 Authentication

This API uses JWT (JSON Web Tokens) for authentication. After logging in, you'll receive:

- **Access Token**: Use this in the `Authorization` header as `Bearer <token>`
- **Refresh Token**: Use this to get new access tokens when they expire

### Protected Routes

All routes except the following require authentication:
- `POST /users/signup`
- `POST /users/login`
- `GET /foods`
- `GET /foods/:food_id`
- `GET /menus`
- `GET /menus/:menu_id`
- `GET /tables`
- `GET /tables/:table_id`
- `GET /orders`
- `GET /orders/:order_id`
- `GET /orderItems`
- `GET /orderItems/:orderItem_id`
- `GET /orderItems-order/:order_id`
- `GET /invoices`
- `GET /invoices/:invoice_id`

## 🗄️ Database Schema

The application uses MongoDB with the following collections:

- **users**: User account information
- **foods**: Food items with pricing and images
- **menus**: Menu categories and time periods
- **tables**: Restaurant table information
- **orders**: Customer orders
- **orderItems**: Individual items within orders
- **invoices**: Billing and payment information

## 🐳 Docker Support


```

**Services:**
- **MongoDB**: Available at `localhost:27017`
- **Mongo Express**: Web interface at `http://localhost:8081`

### Docker Credentials
- **Username**: admin
- **Password**: password123
- **Database**: restraunt

## 📁 Project Structure

```
restraunt-management/
├── controllers/          # API route handlers
│   ├── userController.go
│   ├── foodController.go
│   ├── menuController.go
│   ├── orderController.go
│   ├── orderitemsController.go
│   ├── tableController.go
│   └── invoiceController.go
├── database/            # Database connection
│   └── databaseConnection.go
├── helpers/             # Utility functions
│   └── tokenHelper.go
├── middleware/          # HTTP middleware
│   └── authMiddleware.go
├── models/              # Data models
│   ├── userModel.go
│   ├── foodModel.go
│   ├── menuModel.go
│   ├── orderModel.go
│   ├── orderItemModel.go
│   ├── tableModel.go
│   ├── invoiceModel.go
│   └── noteModel.go
├── routes/              # Route definitions
│   ├── userRouter.go
│   ├── foodRouter.go
│   ├── menuRouter.go
│   ├── orderRouter.go
│   ├── orderitemRouter.go
│   ├── tableRouter.go
│   └── invoiceRouter.go
├── mongo-init/          # MongoDB initialization
├── docker-compose.yml   # Docker services
├── Dockerfile          # App containerization
├── main.go             # Application entry point
├── go.mod              # Go module file
└── go.sum              # Go dependencies
```

## 🧪 Testing

### Using Postman
1. Import the included `Restraunt.postman_collection.json`
2. Set up environment variables for base URL and tokens
3. Run the collection to test all endpoints

### Using cURL
```bash
# Test signup
curl -X POST http://localhost:8000/users/signup \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","email":"john@example.com","password":"password123","phone":"+1234567890"}'

# Test login
curl -X POST http://localhost:8000/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```



### Environment Variables
- `PORT`: Server port (default: 8000)
- `SECRET_KEY`: JWT signing secret
- `MONGODB_URL`: MongoDB connection string

### MongoDB Configuration
- Default database: `restraunt`
- Authentication required
- Connection pooling enabled

## 📝 Development Notes

### Error Handling
- Proper HTTP status codes
- Detailed error messages
- Logging for debugging

### Security Features
- Password hashing with bcrypt
- JWT token expiration
- Protected route middleware
- Input validation

### Performance
- Database connection pooling
- Efficient MongoDB aggregation
- Pagination support for large datasets


## 👥 Author

**Parthav Povil** - [GitHub](https://github.com/parthavpovil)



- Built with Go and the Gin web framework
- MongoDB for database operations
- JWT for secure authentication
- Docker for containerization

---

For more information or support, please open an issue on GitHub.
