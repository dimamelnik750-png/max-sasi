# TODO list of writing own TODO API

## Project setup
- [x] Upload code to GitHub
- [x] Add `.gitignore`
- [ ] Add README with how to run the project
- [x] Create flags for app  
      https://gobyexample.com/command-line-flags
- [ ] Use config library  
      https://github.com/spf13/viper

## Architecture
- [x] Make layers for application:
  - [ ] Database (repository)
  - [ ] Transport (HTTP)
  - [ ] Service (domain / business logic)
- [ ] Use interfaces for the database layer to decouple dependencies

## Data and database
- [ ] Use UUID for keys  
      https://github.com/google/uuid
- [ ] Connect PostgreSQL using:
  - [ ] Standard library: `database/sql`
  - [ ] Driver: `lib/pq`  
        https://github.com/lib/pq
- [ ] Create table for TODO items
- [ ] Add basic database migrations

## API
- [ ] Create CRUD endpoints for TODO items
- [ ] Add request validation
- [ ] Return proper HTTP status codes
- [ ] Handle errors in a consistent way
- [ ] Working with JSON

## Documentation
- [ ] Add Swagger documentation

## Testing
- [ ] Write unit tests
- [ ] Write integration tests

## Reliability
- [ ] Implement graceful shutdown using `context` package  
      https://gobyexample.com/context
- [ ] Add logging for requests and errors

## Deployment
- [ ] Create Dockerfile
- [ ] Add Docker Compose for app + PostgreSQL