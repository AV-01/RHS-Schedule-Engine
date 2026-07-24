<img src="assets/RHS_Schedule_Engine-Banner.png" alt="Banner" width="100%">

# AV's RHS Schedule Engine

High-performance RESTful API Engine built in Go for managing and querying high school schedules, classes, teachers, and student records.

## Description

RHS Schedule Engine is an experimental, high-performance backend built in Go using the Gin web framework and PostgreSQL (Supabase). Designed to handle complex school scheduling architectures, it processes multi-year student rosters, class schedules, and teacher assignments while delivering low-latency queries and low-level control over API security and request throttling.

The project has these features:

* **JWT & Demo Authentication:** Secures sensitive endpoints with 24-hour JWT token validation, alongside an instant zero-config public demo mode (`Bearer demo-key`).

* **Student & Schedule Management:** Supports fast paginated lookup of student profiles, grade-based filtering, name searching, and multi-year historical schedule retrieval.

* **Teacher & Class Analytics:** Provides specialized lookup endpoints for teacher rosters, course catalogs, and teacher period schedule maps.

* **Production Middleware Engine:** Includes in-memory sliding-window rate limiting (60 req/min limit per client IP), automated login audit logging, CORS configuration, and recovery management.

* **Interactive Swagger Documentation:** Integrated OpenAPI/Swagger UI powered by `swaggo/gin-swagger` for interactive in-browser endpoint testing and schema validation.

Overall, I learned a lot during this project... Building a scheduling engine in Go with Gin and PostgreSQL was fast, clean, and ridiculously performant.

## Demo

In case you don't want to spin up PostgreSQL and run it locally, here's a little demo of how it would look!

<table>
<tr>
<td width="50%" align="center">
<h3>Swagger API Docs</h3>
<img src="assets/RHS_Schedule_Engine-Swagger.png" alt="Swagger Docs" width="100%">
</td>
<td width="50%" align="center">
<h3>Demo Video (Click to watch!)</h3>
<a href="https://www.youtube.com/watch?v=TztVnaJ4xbQ">
<img src="https://img.youtube.com/vi/TztVnaJ4xbQ/0.jpg" alt="Demo Video" width="100%">
</a>
</td>
</tr>
</table>

### Dependencies

* In case it's not obvious, you'll need to run Go code. Make sure you have `go` installed (v1.26+ recommended)!
* `PostgreSQL` or `Supabase` database instance (or run in built-in Demo Mode without a database connection)
* `swag` CLI tool if re-generating OpenAPI documentation (`go install github.com/swaggo/swag/cmd/swag@latest`)
* `Docker` if you prefer to run inside an Alpine containerized environment

### Installing

* Clone the repository from GitHub:
* `git clone https://github.com/AV-01/RHS-Schedule-Engine.git`
* `cd RHS-Schedule-Engine`
* Download dependencies:
* `go mod download`
* Create a `.env` file in the root directory:
* `SUPABASE_DB_URL=postgres://user:password@host:port/dbname`
* `JWT_SECRET=your_secret_key_here`

### Executing program

#### If you are running locally with Go
* Make sure you are in the project root directory `RHS-Schedule-Engine\`
* Open terminal and run: `go run main.go`
* Open your browser and navigate to `http://localhost:8080/swagger/index.html`
* Enjoy the magic!

#### If you are using Docker
* Build the container: `docker build -t rhs-schedule-engine .`
* Run the container: `docker run -p 8080:8080 rhs-schedule-engine`
* Access the API at `http://localhost:8080/api/v1`

> [!WARNING]
> Do not commit real `SUPABASE_DB_URL` strings or `JWT_SECRET` keys to public repositories! Ensure production secrets are properly stored in environment variables.

## Help

Open an issue if you need help!

## Version History

* 1.1.0
* Added JWT authentication & login endpoint
* Added rate limiting middleware (60 requests/min per IP)
* Added login audit logger middleware
* Added Swagger API documentation via gin-swagger
* 1.0.0
* Added student, teacher, and class schedule API endpoints
* Integrated PostgreSQL / Supabase backend connection
* Added data parsing and cleaning scripts for multi-year schedule data
* 0.1.0
* Initial Release

## License

This project is licensed under the MIT License - see the LICENSE.md file for details

## Authors
Made with love by AV