# API Documentation

This document describes the RESTful API endpoints provided by the WebCoreGo framework and its modules.

## Table of Contents

1. [Global Endpoints](#global-endpoints)
2. [Module Endpoints](#module-endpoints)
3. [Authentication](#authentication)
4. [Error Responses](#error-responses)
5. [Pagination](#pagination)
6. [Rate Limiting](#rate-limiting)

## Global Endpoints

### Health Check

```
GET /health
```

Returns the health status of the application.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "modules": {
    "module-a": {
      "status": "loaded",
      "version": "1.0.0"
    }
  }
}
```

### Application Info

```
GET /info
```

Returns information about the application and loaded modules.

**Response:**
```json
{
  "name": "WebCoreGo API",
  "version": "1.0.0",
  "environment": "development",
  "modules": [
    {
      "name": "module-a",
      "version": "1.0.0",
      "status": "loaded",
      "routes": [
        "GET /api/v1/module-a/items",
        "POST /api/v1/module-a/items",
        "GET /api/v1/module-a/items/:id",
        "PUT /api/v1/module-a/items/:id",
        "DELETE /api/v1/module-a/items/:id"
      ]
    }
  ]
}
```

## Module Endpoints

### Module A - Items Management

#### Get Items

```
GET /api/v1/module-a/items
```

Retrieve a paginated list of items.

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10) - Number of items per page
- `search` (optional) - Search term
- `status` (optional) - Filter by status

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Sample Item 1",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "name": "Sample Item 2",
      "status": "inactive",
      "created_at": "2024-01-15T10:31:00Z",
      "updated_at": "2024-01-15T10:31:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

#### Create Item

```
POST /api/v1/module-a/items
```

Create a new item.

**Request Body:**
```json
{
  "name": "New Item",
  "status": "active"
}
```

**Response:**
```json
{
  "id": 3,
  "name": "New Item",
  "status": "active",
  "created_at": "2024-01-15T10:32:00Z",
  "updated_at": "2024-01-15T10:32:00Z"
}
```

#### Get Item

```
GET /api/v1/module-a/items/:id
```

Retrieve a specific item by ID.

**Response:**
```json
{
  "id": 1,
  "name": "Sample Item 1",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Update Item

```
PUT /api/v1/module-a/items/:id
```

Update an existing item.

**Request Body:**
```json
{
  "name": "Updated Item",
  "status": "inactive"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Updated Item",
  "status": "inactive",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

#### Delete Item

```
DELETE /api/v1/module-a/items/:id
```

Delete an item by ID.

**Response:**
```json
{
  "message": "Item deleted successfully"
}
```

### Module A - Module Specific Endpoints

#### Module Health

```
GET /api/v1/module-a/health
```

Check the health status of module-a.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

#### Module Info

```
GET /api/v1/module-a/info
```

Get information about module-a.

**Response:**
```json
{
  "name": "module-a",
  "version": "1.0.0",
  "description": "Example module for demonstration",
  "routes": [
    "GET /api/v1/module-a/items",
    "POST /api/v1/module-a/items",
    "GET /api/v1/module-a/items/:id",
    "PUT /api/v1/module-a/items/:id",
    "DELETE /api/v1/module-a/items/:id",
    "GET /api/v1/module-a/health",
    "GET /api/v1/module-a/info"
  ]
}
```

## Authentication

### JWT Authentication

The API uses JWT (JSON Web Token) for authentication. All endpoints except `/health` and `/info` require a valid JWT token in the Authorization header.

**Header:**
```
Authorization: Bearer <token>
```

### Login

```
POST /auth/login
```

Authenticate and receive a JWT token.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### Refresh Token

```
POST /auth/refresh
```

Refresh an expired JWT token using a refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

## Error Responses

All error responses follow a consistent format:

```json
{
  "error": {
    "code": 400,
    "message": "Bad Request",
    "details": "Error details here"
  }
}
```

### Common Error Codes

- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Unprocessable Entity
- `429` - Too Many Requests
- `500` - Internal Server Error
- `503` - Service Unavailable

### Error Response Examples

#### Validation Error

```json
{
  "error": {
    "code": 422,
    "message": "Validation failed",
    "details": {
      "name": "Name is required"
    }
  }
}
```

#### Not Found Error

```json
{
  "error": {
    "code": 404,
    "message": "Resource not found",
    "details": "Item with ID 999 not found"
  }
}
```

#### Authentication Error

```json
{
  "error": {
    "code": 401,
    "message": "Unauthorized",
    "details": "Invalid or expired token"
  }
}
```

## Pagination

All list endpoints support pagination with the following query parameters:

- `page` (default: 1) - Page number
- `page_size` (default: 10, max: 100) - Number of items per page

### Pagination Response Format

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

### Pagination Links

The response includes pagination links for easy navigation:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5,
    "links": {
      "first": "/api/v1/module-a/items?page=1&page_size=10",
      "prev": null,
      "next": "/api/v1/module-a/items?page=2&page_size=10",
      "last": "/api/v1/module-a/items?page=5&page_size=10"
    }
  }
}
```

## Rate Limiting

Rate limiting is applied per authenticated identity (user ID, or the API key
when no user ID is available). It runs after authentication on the protected
route group, so:

- requests with a missing or invalid credential are rejected by authentication
  (401) and never create a rate-limit bucket;
- public routes outside the protected group (health checks, docs) are not
  limited;
- with `auth.type: none` there is no identity, so nothing is limited (a warning
  is logged at startup).

The algorithm is a sliding window counter: no 2x burst at window boundaries.

### Configuration

```yaml
app:
  rate_limit:
    enabled: true
    max: 1000        # requests per window per identity
    window: 1m
    backend: memory  # memory: per process | redis: shared across pods
    fail_open: true  # when redis is unreachable: true lets requests through, false returns 503
```

`backend: redis` uses the loaded `cache:redis` library (requires `redis.host`),
so every pod counts against the same quota. Startup fails when `max <= 0`,
`window <= 0`, or the backend is unknown.

### Rate Limit Headers

Responses on the protected group include:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 995
X-RateLimit-Reset: 2026-10-05T10:01:00Z
```

### Rate Limit Response

When the quota is exhausted the API returns `429`:

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again later."
}
```

## Webhook Support

### Register Webhook

```
POST /webhooks
```

Register a webhook endpoint to receive notifications.

**Request Body:**
```json
{
  "url": "https://your-app.com/webhook",
  "events": ["item.created", "item.updated", "item.deleted"],
  "secret": "your-webhook-secret"
}
```

### Webhook Payload

When an event occurs, the webhook receives:

```json
{
  "event": "item.created",
  "data": {
    "id": 1,
    "name": "Sample Item",
    "status": "active"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## WebSocket Support

### Connect to WebSocket

```
ws://localhost:7272/ws
```

Connect to the WebSocket endpoint for real-time updates.

### WebSocket Events

- `item.created` - New item created
- `item.updated` - Item updated
- `item.deleted` - Item deleted

### WebSocket Message Format

```json
{
  "type": "item.created",
  "data": {
    "id": 1,
    "name": "Sample Item",
    "status": "active"
  }
}
```

## API Versioning

The API uses URL versioning:

- `/api/v1/` - Current stable version
- `/api/v2/` - Future version (when available)

## OpenAPI/Swagger Documentation

API documentation is available at:

```
http://localhost:7272/docs
```

This provides interactive API documentation with:

- Endpoint descriptions
- Request/response examples
- Authentication information
- Try it out functionality

## Best Practices

### 1. Use HTTPS

Always use HTTPS in production to protect data in transit.

### 2. Handle Errors Gracefully

Always check for error responses and handle them appropriately in your client applications.

### 3. Implement Retry Logic

For transient errors, implement retry logic with exponential backoff.

### 4. Use Pagination

Always use pagination when working with list endpoints to avoid performance issues.

### 5. Cache Responses

Use caching for frequently accessed data to improve performance.

### 6. Monitor API Usage

Monitor API usage and set up alerts for unusual activity.

### 7. Keep Tokens Secure

Store JWT tokens securely and implement proper token refresh mechanisms.

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check if the JWT token is valid
   - Verify the Authorization header format
   - Ensure the token hasn't expired

2. **429 Too Many Requests**
   - Wait for the rate limit to reset
   - Implement proper rate limiting in your application
   - Use caching to reduce API calls

3. **500 Internal Server Error**
   - Check server logs for details
   - Verify database connectivity
   - Ensure all required dependencies are available

4. **404 Not Found**
   - Verify the endpoint URL
   - Check if the resource exists
   - Ensure proper permissions

### Debug Mode

Enable debug mode for detailed error information:

```
GET /health?debug=true
```

This returns additional information about the application state and configuration.

## Support

For API support and questions:

1. Check the [module development guide](./module-development.md)
2. Review the [troubleshooting section](#troubleshooting)
3. Contact the development team with specific error messages

## Changelog

### v1.0.0

- Initial API release
- Basic CRUD operations for module-a
- JWT authentication
- Rate limiting
- Pagination support
- WebSocket support
- Webhook support
