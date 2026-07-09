# Go Web App Test Environment Security Checklist

A practical checklist for a **Go web application** hosted on **Cloud Run (or Render)** with a **PostgreSQL database**.

---

## 1. Never Commit Secrets

Never store passwords, API keys, or database connection strings in your source code or commit them to GitHub.

**Good Practice**
- Store secrets as environment variables.
- Use a secret manager if one is available.
- Keep `.env` files out of version control.

Example:

```go
conn := os.Getenv("DATABASE_URL")
```

---

## 2. Use a Separate Test Database

Keep testing and production data completely separate.

```
Production
└── production-db

Testing
└── testing-db
```

Benefits:
- Prevents accidental data loss.
- Allows developers to experiment safely.
- Makes database resets easy.

---

## 3. Use HTTPS

Always serve your application over HTTPS.

Fortunately, Cloud Run and Render provide HTTPS automatically, so no additional configuration is typically required.

---

## 4. Require Authentication

If the application isn't intended for public use yet, require users to log in.

Even a basic authentication system is much safer than leaving administrative pages publicly accessible.

Examples of protected pages:

- `/admin`
- `/dashboard`
- `/settings`

---

## 5. Do Not Expose the Database

Users should never connect directly to your database.

Recommended architecture:

```
Internet
    │
    ▼
Cloud Run / Render
    │
    ▼
PostgreSQL
```

Only the Go application should communicate with the database.

---

## 6. Keep Dependencies Updated

Regularly update your Go dependencies.

```bash
go get -u ./...
go mod tidy
```

After updating, run your tests to verify everything still works correctly.

---

## 7. Validate User Input

Always use parameterized SQL queries.

❌ Avoid:

```go
db.Exec("SELECT * FROM users WHERE name='" + username + "'")
```

✅ Use:

```go
db.Query(
    "SELECT * FROM users WHERE name=$1",
    username,
)
```

Parameterized queries protect your application from SQL injection attacks.

---

## 8. Hide Debug Information

Do not expose detailed error messages to users.

Instead:

- Log detailed errors on the server.
- Return a generic error message to the client.

Good:

```
500 Internal Server Error
```

Bad:

```
Database login failed:
postgres://admin:password123@...
```

---

## 9. Use a Limited-Permission Database Account

Create a dedicated database user for the application.

The application account should only have the permissions it actually needs, such as:

- SELECT
- INSERT
- UPDATE
- DELETE

Avoid using a database administrator account.

---

## 10. Use Fake Data

Avoid storing real personal or sensitive information in your testing environment.

Instead, use generated or anonymized data.

Examples of data to avoid:

- Social Security numbers
- Credit card numbers
- Real passwords
- Medical information
- Personal contact information

---

# Minimum Viable Security

For a class project or internal testing environment, completing the following checklist provides a solid security baseline.

- [ ] HTTPS enabled
- [ ] Secrets stored as environment variables
- [ ] Separate test database
- [ ] Authentication for non-public features
- [ ] Parameterized SQL queries
- [ ] Database not publicly accessible
- [ ] No sensitive real-world data stored

---

# Not a Priority for a Class Project

Unless specifically required, these topics are generally unnecessary for a testing environment:

- Web Application Firewalls (WAF)
- DDoS protection
- Intrusion Detection Systems (IDS)
- Kubernetes security
- Penetration testing
- Zero Trust networking

These become more important when deploying a production application with real users.

---

# Summary

For a Go web application hosted on Cloud Run (or Render) with PostgreSQL, the recommended approach is:

```
GitHub
    │
    ▼
Cloud Run / Render
    │
    ▼
PostgreSQL
```

Follow these principles:

- Keep secrets out of source control.
- Separate testing and production environments.
- Require authentication.
- Use HTTPS.
- Parameterize all SQL queries.
- Restrict database access.
- Use fake data for testing.

Following these practices provides a secure and maintainable testing environment without adding unnecessary complexity.