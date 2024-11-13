
# Proxy Load Balancer

## Overview

This project implements a load balancer using the **Go** programming language and the Gin web framework. It forwards incoming requests to a set of backend servers using a round-robin strategy, with health checks to ensure only healthy servers receive traffic. This implementation also includes query parameter tracking for each forwarded request.

---

## How It Works

1. **Proxies Setup:** 
   - Backend servers are defined in a list (`targetServers`).
   - Each server is wrapped in a `Proxy` object, which includes a reverse proxy and a health-checking mechanism.

2. **Health Checks:** 
   - The `HealthChecker` runs in a separate goroutine, periodically checking the `/health` endpoint of each backend server.

3. **Request Handling:** 
   - Incoming requests are routed to the next healthy backend server using a round-robin strategy.
   - A query parameter (`reqCounter`) is added to track requests.

4. **Gin Web Framework:** 
   - The application is built on the Gin framework, allowing for efficient request handling and routing.

---

## Installation and Usage

### Prerequisites

- Go (>=1.16)
- Gin web framework (`go get -u github.com/gin-gonic/gin`)

### Steps

1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```
2. Navigate to the project directory:
   ```bash
   cd proxy-load-balancer
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
4. Run the backend servers:
     ```
     go run backend/main.go 8080
     go run backend/main.go 8081
     go run backend/main.go 8082
     ```
5. Send curl command:
    ```
        curl http://localhost/ping
    ```
6. Test using ab
    ```
        ab -n 100000 -c 1 http://localhost/ping 
    ```
7. Test using wrk
    ```
        wrk -t1 -c5 -d60 http://localhost/ping
    ```

---

## License

This project is licensed under the MIT License.

---

## Acknowledgments

- Inspired the solution from [Coding Challenges by John Crickett](https://codingchallenges.fyi/challenges/challenge-load-balancer)
