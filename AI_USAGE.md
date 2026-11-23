# AI Usage Disclosure

## Purpose

AI assistance was used as a collaborative development aid while building this take-home project. It was used to explore options, explain unfamiliar Go and React concepts, propose implementation steps, review edge cases, and accelerate routine development work.

This document intentionally describes the nature of the collaboration instead of reproducing a private conversation verbatim.

## How AI was used

### Planning and scope

- Turned the assignment brief into a small, time-boxed implementation plan.
- Compared implementation options and kept the scope focused on the required four operations.
- Evaluated trade-offs such as `float64` versus decimal arithmetic, Docker versus additional product functionality, and a lightweight hexagonal architecture versus unnecessary framework complexity.

### Architecture and design review

- Discussed a feature-based React structure and a lightweight hexagonal Go backend.
- Refined the user-requested hexagonal architecture into a proportionate implementation for a small calculator API rather than applying a heavyweight framework.
- Used SOLID principles as an explicit user-requested design reference. The Strategy pattern for calculator operations and a small resolver interface were selected to demonstrate single responsibility, open/closed design, and dependency inversion without adding unnecessary abstractions.
- Reviewed naming conventions, input constraints, error handling, request IDs, structured logging, and API documentation.

### Implementation support

- Assisted with generating and refining frontend and backend code.
- Helped create unit-test scenarios for normal paths, validation failures, malformed requests, division by zero, non-finite results, request-size limits, and UI states.
- Helped draft the OpenAPI contract and project documentation.


## How the key decisions were reached

The user set the architectural direction: Go for the backend, React with TypeScript for the frontend, SOLID principles as a design reference, and a hexagonal architecture that would remain readable for a small take-home assignment.

AI then helped evaluate the trade-offs and translate those constraints into concrete choices:

- A single generic calculation endpoint supports the selected operations while keeping the API compact.
- Independent calculation strategies allow a new operation to be added without changing existing operation implementations.
- A small application-facing resolver interface separates the use case from the in-memory strategy registry.
- The HTTP adapter is responsible for JSON, request validation, status codes, request IDs, logging, and recovery; the domain remains independent of HTTP concerns.
- The result is intentionally lightweight: no database, framework, event bus, dependency-injection container, or speculative ports were added.

The user reviewed and approved these decisions as the implementation progressed.

## Human decisions and review

The following decisions were explicitly reviewed and chosen for this project:

- Go and React with TypeScript were selected to match the assignment.
- Hexagonal architecture was selected by the user, with SOLID principles as a reference for the implementation.
- The mandatory scope is addition, subtraction, multiplication, and division; advanced operations remain intentionally out of scope because they are optional.
- The API is the source of truth for validation and calculation rules.
- Input limits are 20 operands, a maximum absolute operand of `1e14`, up to 10 decimal places, and a 1 MiB request body.
- `float64` is acceptable for this assignment, with its precision limitations documented in the README.
- Docker was added as an optional reproducible way to run the complete application; native local development remains supported.

All generated suggestions were reviewed, adapted when necessary, and validated with the project's test, lint, and build commands before being treated as part of the solution.
