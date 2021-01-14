# Cat-Scribers
[![Go Reference](https://pkg.go.dev/badge/github.com/EtienneBerube/cat-scribers.svg)](https://pkg.go.dev/github.com/EtienneBerube/cat-scribers)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/EtienneBerube/cat-scribers/blob/main/LICENSE)

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)
[![Ask Me Anything !](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)](https://GitHub.com/EtienneBerube)

This project is an implementation of the Shopify Developer Challenge 2021. You can access the prompt by clicking 
[here](https://docs.google.com/document/d/1ZKRywXQLZWOqVOHC4JkF3LqdpO3Llpfk_CkZPR8bjak/edit).

The __Cat-Scibers__ project is a subscription based REST API to subscribe to cat photos. Its purpose is to provide the 
ability to every user to subscribe to other user's cat photos. Not only that, but it is also impossible to submit a 
picture that does not contain a Cat. This project utilizes Google Vision API to ensure that the content of the site 
stays pure with respect to its original purpose.

Currently, payments are taken every first of the month at midnight, or when you subscribe to another user.

## How to run
1) Install Golang (this project was developed with `go 1.14.2 windows/amd64`)
2) Run `go get`
3) Get an API key for Google Vision API
4) Get a MongoDB instance running (the original version uses [MongoDB Atlas](https://www.mongodb.com/cloud/atlas), but a local version would work just as well. 
5) Set the following Environment Variables:
```Env
HTTP_PORT=<Port number>
GOOGLE_APPLICATION_CREDENTIALS=<Path to Google JSON credentials file> 
MONGODB_URL=<MongoDB connection URL>
JWT_SECREt=<JWT secret>
```


## API documentation

### Open
* `GET /ping` - Health
    * Request: `empty`
    * Response: "Pong"
* `POST /login` - Login
    * Request:
      ``` 
        email: string,
        password: string
      ```
      
    * Response: 
        ```
        user_id: string,
        token: string
        ```
* `POST /signup` - Sign up
    * Request:
      ``` 
      name: string,
      email: string,
      password: string,
      subscription_price: number
      ```
      
    * Response:
      ```
        user_id: string,
        token: string
      ```
* `GET /users` - Get all users
    * Request: `empty`
    
    * Response:

        _Array of:_
        ```
        id: string
        name: string
        email: string
        subscriptions: string[]
        photos: string[]
        balance: number
        subscription_price: number
       ```
* `GET /user/:id` - Get user by ID
    * Request: `empty`

    * Response:
        ```
        id: string
        name: string
        email: string
        subscriptions: string[]
        photos: string[]
        balance: number
        subscription_price: number
       ```

### Require to be authenticated
To access those routes, the request must have an Authorization header. Example: `Authorization Bearer <token> `

* `GET /user` - Get current user
    * Request: `empty`

    * Response:
        ```
        id: string
        name: string
        email: string
        subscriptions: string[]
        photos: string[]
        balance: number
        subscription_price: number
       ```
* `PUT /user` - Update current user
    * Request:
        ```
        id: string
        name: string
        email: string
        subscriptions: string[]
        photos: string[]
        balance: number
        subscription_price: number
       ```

    * Response:
        ```
        id: string
        name: string
        email: string
        subscriptions: string[]
        photos: string[]
        balance: number
        subscription_price: number
       ```
* `DELETE /user` - Delete current user
    * Request: `empty`
    * Response: `message: successul`
      
* `POST /user/photo/` - Upload photo
    * Request:
      ```
      name: string
      description: string
      base64: string
      type: string
      ```
  
    * Response:
      ```
      id: string
      ```
* `POST /user/photos/` - Upload multiple photos
    * Request:
      
      _Array of:_
        ```
        name: string
        description: string
        base64: string
        type: string
        ```

    * Response:
      ```
      ids: string[]
      rejected: string[]
      ```
* `DELETE /user/photo/:id` - Delete a photo by ID (current user)
    * Request: `empty`
    
    * Response:
       ```
       message: string
       ```
* `GET /photo/:id/` - Get a photo by ID (must be subscribed to owner)
    * Request: `empty`

    * Response:
       ```
        id: string
        name: string
        description: string
        base64: string
        type: string
        owner_id: string
        ```
* `GET /user/:id/photos/` - Get all photos from user with ID (must be subscribed to owner)
    * Request: 
      * Query Param: `name: string` - a string to search pictures by
    * Response:
      _Array of:_
       ```
        id: string
        name: string
        description: string
        base64: string
        type: string
        owner_id: string
        ```
* `POST /subscribe/:id` - Subscribe current user to user with ID
    * Request: `empty`
    * Response:
       ```
        message: string
        ```
* `DELETE /subscribe/:id` - unsubscribe current user from user with ID
    * Request: `empty`
    * Response:
       ```
        message: string
        ```

## Backlog
* Implement all relevant HTTP methods for a complete REST project 
  (This project implemented the basic methods for the MVP)
    
* Implement optimized queries for specific operations that could be performed in a single aggregation pipeline 
  and would bring greater performances for certain operations.
  
* Fix some obvious problems (ex: Can update the balance of a user)

