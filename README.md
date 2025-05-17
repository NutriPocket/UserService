# User Service

Routes:

    - /auth
      - POST /register
      - POST /login
      - POST /logout
    - /users
      - GET /
      - GET /:username

Build & Run

## Compose

### Run it

-   If you use another PORT than 8080, change it in docker-compose.yaml

```
docker-compose up --build # Build and run app + postgresql
docker-compose up app # Run app
docker-compose up database # Run postgresql
```

### Free resources

```
docker-compose down # Down services
docker-compose down --volumes # Removes services and volumes (postgresql persisted data)
```
