package router

import "net/http"

type Middleware func(http.Handler) http.Handler

// Router struct to hold registered middleware and routes
type Router struct {
	middleware []Middleware
	routes     map[string]http.Handler
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		middleware: []Middleware{},
		routes:     make(map[string]http.Handler),
	}
}

// Use registers a middleware to be used by the router
func (r *Router) Use(m Middleware) {
	r.middleware = append(r.middleware, m)
}

// Handle registers a route with the given method and path to be handled by the provided handler
func (r *Router) Handle(method, path string, handler http.Handler) {
	r.routes[method+path] = handler
}

// ServeHTTP implements the http.Handler interface and applies all registered middleware to incoming requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create a new http.Handler for the route
	handler := r.routes[req.Method+req.URL.Path]
	// Apply middleware to the handler in reverse order
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	// Serve the request
	handler.ServeHTTP(w, req)
}
