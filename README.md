# TODO list of writing own TODO API

## Project setup
- [x] Upload code to GitHub
- [x] Add `.gitignore`
- [x] Add README with how to run the project
- [x] Create flags for app  
      https://gobyexample.com/command-line-flags
- [x] Use config library  
      https://github.com/spf13/viper

## Configuration
Create a `.env` file in the project root based on `.env.example`.

## Architecture
- [x] Make layers for application:
  - [x] Database (repository)
  - [x] Transport (HTTP)
  - [x] Service (domain / business logic)
- [x] Use interfaces for the database layer to decouple dependencies

## Data and database
- [x] Use UUID for keys  
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
