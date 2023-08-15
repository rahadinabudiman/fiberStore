# fiberStore - Online Store Application using Golang and Fiber

**fiberStore** is a simple online store application developed using the Go programming language and the Fiber web framework. This application allows users to browse products, add them to their cart, and view their cart's contents.

## ERD fiberStore

![ERD fiberStore](https://i.imgur.com/HAptHQi.png)

## Features

- Browse Products: Users can view a list of available products with their details.
- Browse by Category: Customers can view product lists by product category.
- Add to cart: Customers can add products to their shopping cart.
- View cart: Customers can see a list of products that have been added to the shopping cart.
- Remove from cart: Customers can delete products from the shopping cart.
- Checkout and Payment: Customers can proceed to checkout and make payment TransactionDetails.
- Latest Transactions: Customers can view a list of their latest transactions.
- User Authentication: Customers can create accounts, log in, and manage their profiles.

## Technologies Used

- Go (Golang): The backend of the application is developed using the Go programming language.
- Fiber: Fiber is a fast and lightweight web framework for Go that's used to build the web server.
- Cloudinary: Cloudinary is used for handling image and media uploads.
- Validator: Validator is used for data validation.
- JWT: JWT is used for implementing JSON Web Tokens for authentication.
- Godotenv: used for managing environment variables.
- Crypto: is used for cryptographic functions.
- GORM: The GORM ORM library is used for interacting with the MySQL database.

## Installation

NOTE : if you running on local machine, don't forget to create MySQL Database which name is fiberstore. After successfully running, table will auto migration.

1. Clone this repository using git to your local machine

```
git clone https://github.com/rahadinabudiman/fiberStore.git
```

2. go to folder fiberStore

```
cd fiberStore
```

3. Install the required dependencies using Go:

```
go get
```

4. Rename .env.example file:

```
cp .env.example .env
```

5. Open the .env file using a text editor and configure the necessary environment variables, such as database connection settings and Cloudinary credentials.
6. Build and run the application:

```
go run main.go
```

## Installation With Docker

1. Pull image repsitory to your local machine

```
docker pull kmoonkinan/fiberstore:latest
```

2. Run Project with docker-compose

```
docker-compose up
```

## Usage

1. Open your web browser or Postman and navigate to http://localhost:{PORT} to access the fiberStore application.
2. Browse through the available products, add them to your cart, view your cart's contents, manage your cart, and proceed to checkout.
3. Use the authentication features to register and log in as a customer.

## API Documentation

After running project in local machine, you can access `/api/v1/swagger/index.html` or you can access to deployed link : [Deployed Link By AWS](http://54.179.176.114:1309/api/v1/swagger/index.html)

## Usage with Postman

You can also test using Postman provided in the `.postman_collection` and `.postman_environment` files to facilitate testing.

## Demo

1. Account Customer

```
username: customer
password: fiberstore
```

2. Account Admin

```
username: admin
password: fiberstore
```

## Open Endpoints

Open endpoints require no Authentication.

- Login : `POST /api/v1/login`
- Register : `POST /api/v1/register`
- Get Product By Id : `POST /api/v1/product/findById?id=1`
- Get All Product : `POST /api/v1/product?page=1&limit=10`
- Get All Product By Name : `POST /api/v1/product/findByName?name=adeuy&page=1&limit=2`
- Get All Product By Category : `POST /api/v1/product/findByCategory?category=DanGdut&page=1&limit=10`
