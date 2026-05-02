---
name: add-api-endpoint
description: Add a new REST API endpoint to CanIFly — backend handler, route registration, service logic, and frontend axios client. Use when the user asks to add an endpoint, expose new backend functionality to the UI, or wire up a new `/api/...` route.
---

# Add API Endpoint

CanIFly serves REST APIs from the Go backend on `http://localhost:42423` and consumes them via axios in the React renderer. New endpoints touch four files in a fixed order. Follow the steps below — do not skip the verification step.

## Inputs you need before starting

- **Route**: e.g. `/api/skill-plans/import`
- **HTTP method(s)**: GET / POST / PUT / DELETE
- **Auth**: does it require a valid session? (Most `/api/*` routes do.)
- **Domain**: which existing handler file does it belong to? (`internal/handlers/`: account, character, config, evedata, fuzzworks, skillplan, auth, websocket, …) — create a new file only if no existing one fits.
- **Request/response shape**: input fields and output JSON.

## Steps

1. **Define the handler** in `internal/handlers/<domain>.go`.
   - Receiver matches existing handlers in that file (e.g. `func (h *SkillPlanHandler) Import() http.HandlerFunc`).
   - Return `http.HandlerFunc`. Decode JSON via `json.NewDecoder(r.Body).Decode(&req)`. Encode responses via `utils.go` helpers if present in the file's domain.
   - Call into `internal/services/<domain>/...` for business logic — keep handlers thin.

2. **Register the route** in `internal/server/router.go`.
   - Add `r.HandleFunc("/api/...", handler.Method()).Methods("POST")` near related routes.
   - If the route requires auth, register it inside the existing authenticated subrouter (look for the pattern used by `/api/session/refresh` and other authenticated routes in the file).

3. **Implement the service logic** in `internal/services/<domain>/`.
   - Add the method to the relevant service struct. Define an interface entry in `internal/services/interfaces/` if the service is consumed via interface elsewhere.
   - Persist via `internal/persist/` stores; do not write files directly from the service.
   - Return errors with context — they will be logged via logrus by the handler.

4. **Add the frontend client** in `renderer/src/api/`.
   - Extend `apiService.jsx` with a new function that calls `apiRequest(...)` from `apiRequest.jsx`.
   - Mirror the request/response types loosely — JS is untyped but keep field names identical to the Go JSON tags.

5. **Wire it into the UI** (only if the user asked for UI changes — otherwise stop here and report the API is ready).
   - Find the page/component under `renderer/src/pages/` or `renderer/src/components/` and call the new api function from a hook in `renderer/src/hooks/`.

## Verification (do not skip)

Run from repo root:

```bash
npm run test:go        # Go tests
npm run test:react     # Vitest
go build ./...         # ensures handler + router compile
```

If you added Go imports, run `go mod tidy`. If you only changed Go code, `npm run dev:go` is enough to smoke-test; if you changed both sides, `npm start`.

## Conventions to respect

- **Linting**: `golangci-lint` for Go, ESLint + Prettier for React. Don't reformat unrelated lines.
- **Logging**: logrus with structured fields, not `fmt.Println`.
- **Errors**: explicit error returns; do not swallow.
- **No global state**: services are constructed in `internal/cmd/` and injected.
- **Trust internal callers**: validate at the HTTP boundary (request decoding) — not between service layers.

## Common pitfalls

- Forgetting to register the route in `router.go` — the handler compiles but 404s at runtime.
- Adding the route outside the authenticated subrouter when it should be authenticated.
- Hardcoding the port. Always go through the existing axios base URL config in `apiRequest.jsx`.
