package main

import (
	"api-clean/features/todo"
	"api-clean/internal/config"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	// Protected routes
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		// r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		// 	_, claims, _ := jwtauth.FromContext(r.Context())
		// 	w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		// })
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/api/todo", todo.Routes(configuration))
		})
	})

	// Public routes
	router.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response := make(map[string]string)
			response["message"] = "welcome anonymous"
			render.JSON(w, r, response) // Return some demo response
		})
		r.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
			response := make(map[string]string)
			response["user_id"] = r.FormValue("user_id")
			if len(response["user_id"]) == 0 {
				response["message"] = "user_id needed"
			} else {
				_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": response["user_id"]})
				response["message"] = "success"
				response["code"] = tokenString
			}
			render.JSON(w, r, response) // Return some demo response
		})
	})

	return router
}
