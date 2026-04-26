# SIPELAN Backend Context

## Architecture
- **Language:** Go
- **Web Framework:** Chi Router (`github.com/go-chi/chi/v5`)
- **ORM:** GORM (`gorm.io/gorm`)
- **Database:** SQLite (`sipelan.db`)
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`)

## Authentication Logic
- **Secret Key:** Loaded from `JWT_SECRET_KEY` in `.env`.
- **Token Duration:** 24 hours.
- **Middleware:** `AuthMiddleware` in `common/middleware.go`.
    - Handles case-insensitive "Bearer " prefix stripping.
    - Validates JWT claims (must include a valid person ID).
    - Verifies user existence in the database.
    - Injects the authenticated `models.Person` into the request context.
- **Error Handling:** Returns `401 Unauthorized` with specific messages for:
    - Token not provided.
    - Invalid format.
    - Validation failure (signature/claims).
    - Expired token.
    - User not found in DB.

## Recent Changes
- Improved `AuthMiddleware` robustness for token stripping and whitespace handling.
- Enhanced 401 error messages to provide clearer debugging information.
