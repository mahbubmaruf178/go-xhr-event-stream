package main

import (
	"fmt"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Get the request context to detect client disconnection
	conn := r.Context()

	// Start a loop to send SSE messages
	for i := 1; i <= 2000; i++ {
		// Check if the client has closed the connection
		select {
		case <-conn.Done():
			// If client has closed the connection, break the loop
			fmt.Println("Client disconnected, stopping the loop")
			return
		default:
			// Send an SSE message
			fmt.Fprintf(w, "data: Message %d\n\n", i)
			// Flush the response to ensure the client receives the message immediately
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			// Simulate some delay between messages
			time.Sleep(1 * time.Second)
			fmt.Println("Message", i, "sent")
		}
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)
	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files like HTML
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
