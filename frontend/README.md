# Sezzle Calculator Frontend

React and TypeScript user interface for the Go calculator API. The frontend uses Vite, native CSS, `fetch`, and a feature-based structure; no UI framework or global state library is required for this scope.

## Requirements

- Node.js 26 LTS or newer.
- npm 11 or newer.
- The backend available at `http://localhost:8080`.

## Run in development

Start the backend in one terminal:

```bash
cd ../backend
go run ./cmd/api
```

Then start the frontend in another terminal:

```bash
cd frontend
npm install
npm run dev
```

Open the URL printed by Vite, normally `http://localhost:5173`. The development proxy forwards `/api` calls to `http://localhost:8080`, so CORS configuration is not required locally.

`npm run dev` starts the application. `npm run build` only generates optimized static files in `dist`; it does not run a server. To inspect those files after building, use `npm run preview`.

For the optional containerized setup, run `docker compose up --build` from the repository root and open `http://localhost:3000`. Nginx serves the production build and proxies `/api` requests to the backend container.

## Scripts

```bash
npm run dev             # development server
npm run lint            # static analysis
npm run test -- --run   # non-interactive tests
npm run test:coverage   # coverage report
npm run build           # TypeScript validation and production build
```

## Structure

```text
src/
├── app/                         # application composition and global styles
└── features/calculator/
    ├── api/                     # HTTP client for the backend contract
    ├── components/              # React user interface
    ├── hooks/                   # state and event coordination
    └── model/                   # types, reducer, and pure validation
```

## Input rules

The form keeps operands as text while users edit them. `inputMode="decimal"` improves the mobile keyboard, while TypeScript validates the value before sending it to the API.

- A value may use digits, an optional leading minus sign, and one period decimal separator.
- Thousands separators, decimal commas, exponent notation, letters, and currency symbols are not accepted.
- A maximum of 20 operands can be added for addition, subtraction, and multiplication.
- Division and percentage always use exactly two operands.
- Percentage uses `amount × percentage / 100`; for example, `200` and `15` return `30`.
- Each field is limited to 17 characters, 10 decimal places, and 15 significant digits.
- The accepted numeric range is zero or an absolute value from `0.0000000001` to `100000000000000`.

The backend remains the source of truth. Frontend validation improves feedback and prevents unnecessary requests, but the API validates the same contract independently.

## Request errors

The API client uses `AbortController` with a 10-second timeout. It distinguishes validation/API errors, network failures, timeouts, and malformed successful responses. In development, it logs only error category, HTTP status, API error code, and `X-Request-ID`; it never logs operand values.
