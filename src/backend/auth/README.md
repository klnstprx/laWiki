**OAuth 2.0 Flow Implementation:**

1. **Login Handler:**

   - Generates a random `state` string to prevent CSRF attacks.
   - Stores the `state` in a cookie.
   - Redirects the user to Google's OAuth 2.0 consent screen using `oauthConfig.AuthCodeURL(state)`.

2. **Callback Handler:**
   - Validates the `state` parameter by comparing it to the one stored in the cookie.
   - Exchanges the authorization code for an access token using `oauthConfig.Exchange`.
   - Retrieves user information from Google using the access token.
   - Generates a JWT token containing user information and an expiration time.
   - Stores the JWT token in a secure, HTTP-only cookie.

**JWT Token Generation and Verification:**

- **Token Creation:**

  - User information is unmarshaled into a `User` struct.
  - JWT claims are created, including the user's email, name, and expiration time.
  - The token is signed using the HMAC SHA256 algorithm with a secret key stored in the `JWT_SECRET` environment variable.

- **Token Maintenance:**
  - The client stores the JWT token in a secure cookie.
  - The token includes an expiration time (`exp` claim), after which it becomes invalid.
  - There is no need to maintain server-side sessions since JWT tokens are stateless.

**Authentication Maintenance and Expiration:**

- Authentication is maintained by including the JWT token in subsequent requests, either in a cookie or an `Authorization` header.
- The token's expiration time defines when authentication expires.
- After expiration, the user must re-authenticate by repeating the OAuth flow.
- Optionally, implement refresh tokens or re-authentication strategies for a better user experience.

**Protected Routes:**

- Example of a protected endpoint (`/protected`) demonstrates how to verify the JWT token.
- The handler extracts the token from the request, verifies its signature and expiration, and allows access if valid.

**Security Considerations:**

- Use HTTPS to ensure secure transmission of tokens.
- Set cookies with `HttpOnly` and `Secure` flags.
- Regularly rotate JWT secret keys and handle token revocation if necessary.

**Environment Variables Required:**

- `PORT`: Port on which the auth service runs (default is `8002`).
- `GOOGLE_OAUTH_CLIENT_ID`: Google OAuth 2.0 Client ID.
- `GOOGLE_OAUTH_CLIENT_SECRET`: Google OAuth 2.0 Client Secret.
- `GOOGLE_OAUTH_REDIRECT_URL`: URL to which Google redirects after authentication (e.g., `https://yourdomain.com/callback`).
- `JWT_SECRET`: Secret key used to sign JWT tokens.
