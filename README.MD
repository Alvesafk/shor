# shor

shor is a simple url shortener backend made in **go**, it uses firestore as it's DB.

It has the following routes.
```go
// Get a hello world.
r.Get("/", conn.GetHelloWorld)

// Post the URL that you want to save on the DB, send it as a JSON, with a "url" field.
// It will return a JSON with the ShortURL.
r.Post("/shor", conn.PostURL)

// Get the URL to redirect you to the long url.
r.Get("/shor/{shortUrl}", conn.GetURL)

// Get the URL struct.
r.Get("/shor/{shortUrl}/stats", conn.GetURL)

// Put updates the URL struct with a new long URL sent via JSON.
r.Put("/shor/{shortUrl}", conn.UpdateURL)

// Deletes a URL struct.
r.Delete("/shor/{shortUrl}", conn.DeleteURL)
```

Build it or run it via `go run cmd/api/main.go`.

You need a service key of the firestore.
