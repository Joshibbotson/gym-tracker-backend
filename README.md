handler: HTTP Handlers
Service: Core business logic
Repository: Data access layer

Sessions in web applications are a way to manage and maintain user state across requests. This is especially helpful for authenticated users, so that after they log in, they don’t have to authenticate again on every page they visit. Here’s a general overview of how sessions work, how they compare to JWT tokens, and when they can be used together.

How Sessions Work
User Authentication: When a user logs in, the backend verifies their credentials (e.g., username and password).
Session Creation: If the credentials are valid, the server generates a session, often identified by a unique session ID.
Session Storage: The session ID and associated data (like user ID, roles, etc.) are stored on the server in memory, a database, or a cache like Redis.
Session Cookie: The server sends a session ID to the client, typically as a cookie. This cookie is automatically sent back with each subsequent request to the server.
Session Validation: For each request, the server checks the session ID in the cookie against the session data it has stored to identify the user.

Err: when a cookie is no longer present we must reset the localstorage on the frontend and prompt a relogin.
