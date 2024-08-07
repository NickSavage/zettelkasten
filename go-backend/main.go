package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"go-backend/models"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/stripe/stripe-go"
)

var s *Server

type Server struct {
	db             *sql.DB
	s3             *s3.Client
	testing        bool
	jwt_secret_key []byte
	stripe_key     string
	mail           *MailClient
	TestInspector  *TestInspector
}

type MailClient struct {
	Host     string
	Password string
}

type TestInspector struct {
	EmailsSent    int
	FilesUploaded int
}

func admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("current_user").(int)
		user, err := s.QueryUser(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		if !user.IsAdmin {
			http.Error(w, "Access denied", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenStr = tokenStr[len("Bearer "):]

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return s.jwt_secret_key, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), "current_user", claims.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type Email struct {
	Subject   string `json:"subject"`
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
}

func (s *Server) SendEmail(subject, recipient, body string) error {
	if s.testing {
		s.TestInspector.EmailsSent += 1
		return nil
	}
	email := Email{
		Subject:   subject,
		Recipient: recipient,
		Body:      body,
	}

	// Convert email struct to JSON

	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}
	go func() {

		// Create a new request
		req, err := http.NewRequest("POST", s.mail.Host+"/api/send", bytes.NewBuffer(emailJSON))
		if err != nil {
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", s.mail.Password)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to send email: %s", resp.Status)
			return
		}
	}()
	return nil
}

func addProtectedRoute(r *mux.Router, path string, handler http.HandlerFunc, method string) *mux.Route {
	return r.HandleFunc(path, jwtMiddleware(logRoute(handler))).Methods(method)

}

func addRoute(r *mux.Router, path string, handler http.HandlerFunc, method string) *mux.Route {
	return r.HandleFunc(path, logRoute(handler)).Methods(method)
}

func main() {
	file, err := openLogFile(os.Getenv("ZETTEL_BACKEND_LOG_LOCATION"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	s = &Server{}

	dbConfig := models.DatabaseConfig{}
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Password = os.Getenv("DB_PASS")
	dbConfig.DatabaseName = os.Getenv("DB_NAME")

	db, err := ConnectToDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	s.db = db
	s.runMigrations()
	s.s3 = s.createS3Client()

	s.stripe_key = os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	s.mail = &MailClient{
		Host:     os.Getenv("MAIL_HOST"),
		Password: os.Getenv("MAIL_PASSWORD"),
	}
	s.jwt_secret_key = []byte(os.Getenv("SECRET_KEY"))

	r := mux.NewRouter()
	addProtectedRoute(r, "/api/auth", s.CheckTokenRoute, "GET")
	addRoute(r, "/api/login", s.LoginRoute, "POST")
	addRoute(r, "/api/reset-password", s.ResetPasswordRoute, "POST")
	addRoute(r, "/api/email-validate", s.ValidateEmailRoute, "POST")
	addRoute(r, "/api/request-reset", s.RequestPasswordResetRoute, "POST")

	addProtectedRoute(r, "/api/files", s.GetAllFilesRoute, "GET")
	addProtectedRoute(r, "/api/files/upload", s.UploadFileRoute, "POST")
	addProtectedRoute(r, "/api/files/{id}", s.GetFileMetadataRoute, "GET")
	addProtectedRoute(r, "/api/files/{id}", s.EditFileMetadataRoute, "PATCH")
	addProtectedRoute(r, "/api/files/{id}", s.DeleteFileRoute, "DELETE")
	addProtectedRoute(r, "/api/files/download/{id}", s.DownloadFileRoute, "GET")

	addProtectedRoute(r, "/api/cards", s.GetCardsRoute, "GET")
	addProtectedRoute(r, "/api/cards", s.CreateCardRoute, "POST")
	addProtectedRoute(r, "/api/next", s.NextIDRoute, "POST")
	addProtectedRoute(r, "/api/cards/{id}", s.GetCardRoute, "GET")
	addProtectedRoute(r, "/api/cards/{id}", s.UpdateCardRoute, "PUT")
	addProtectedRoute(r, "/api/cards/{id}", s.DeleteCardRoute, "DELETE")

	addProtectedRoute(r, "/api/users/{id}", s.GetUserRoute, "GET")
	addProtectedRoute(r, "/api/users/{id}", s.UpdateUserRoute, "PUT")
	addProtectedRoute(r, "/api/users", s.GetUsersRoute, "GET")
	addRoute(r, "/api/users", s.CreateUserRoute, "POST")
	addProtectedRoute(r, "/api/users/{id}/subscription", s.GetUserSubscriptionRoute, "GET")
	addProtectedRoute(r, "/api/current", s.GetCurrentUserRoute, "GET")
	addProtectedRoute(r, "/api/admin", s.GetUserAdminRoute, "GET")

	addProtectedRoute(r, "/api/tasks/{id}", s.GetTaskRoute, "GET")
	addProtectedRoute(r, "/api/tasks", s.GetTasksRoute, "GET")
	addProtectedRoute(r, "/api/tasks", s.CreateTaskRoute, "POST")
	addProtectedRoute(r, "/api/tasks/{id}", s.UpdateTaskRoute, "PUT")
	addProtectedRoute(r, "/api/tasks/{id}", s.DeleteTaskRoute, "DELETE")

	addRoute(r, "/api/billing/create_checkout_session", s.CreateCheckoutSession, "POST")
	addRoute(r, "/api/billing/success", s.GetSuccessfulSessionData, "GET")
	addRoute(r, "/api/webhook", s.HandleWebhook, "POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ZETTEL_URL")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"authorization", "content-type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}
