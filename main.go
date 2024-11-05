package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	"github.com/guarzo/canifly/internal/service/skillplan"
	"github.com/guarzo/canifly/internal/service/skilltype"
	"github.com/guarzo/canifly/internal/utils/crypto"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

var version = "0.0.39"

func main() {
	xlog.Logf("Starting application, version %s", version)
	loadEnv()

	// Initialize application components
	secret := getSecretKey()
	initializeComponents(secret)

	// Set up router and server
	r := setupRouter(secret)
	port := getPort()
	startServer(r, port)
}

// loadEnv loads environment variables from a .env file, the embedded file, or the system environment.
func loadEnv() {
	// First, try to load from an existing .env file in the file system
	if err := godotenv.Load(); err != nil {
		// If the .env file is not found, try to load it from the embedded StaticFiles
		embeddedEnv, err := embed.StaticFiles.Open("static/.env")
		if err != nil {
			log.Println("Failed to load embedded .env file:", err)
			log.Println("Using system environment variables.")
			return
		}
		defer embeddedEnv.Close()

		// Parse the embedded .env file
		envMap, err := godotenv.Parse(embeddedEnv)
		if err != nil {
			log.Println("Failed to parse embedded .env file:", err)
			return
		}

		// Set each environment variable from the parsed map
		for key, value := range envMap {
			os.Setenv(key, value)
		}
	}
}

// getSecretKey retrieves or generates the encryption secret key as a string.
func getSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		key, err := flyHandlers.GenerateSecret()
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		secret = base64.StdEncoding.EncodeToString(key)
		log.Println("Warning: using a generated key for testing only.")
	}
	return secret
}

// initializeComponents sets up the app's core components (crypto, OAuth, skill plans, skill types).
func initializeComponents(secret string) {
	// Decode the base64 secret for crypto initialization
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Fatalf("Failed to decode secret key: %v", err)
	}
	if err := crypto.Initialize(key); err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	clientID, clientSecret, callbackURL := os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_CLIENT_SECRET"), os.Getenv("EVE_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		log.Fatalf("EVE_CLIENT_ID, EVE_CLIENT_SECRET, and EVE_CALLBACK_URL must be set")
	}
	api.InitializeOAuth(clientID, clientSecret, callbackURL)

	if err := skillplan.ProcessSkillPlans(); err != nil {
		log.Fatalf("Failed to load skill plans: %v", err)
	}

	if err := skilltype.LoadSkillTypes(); err != nil {
		log.Fatalf("Failed to load skill types: %v", err)
	}

	if err := embed.LoadStatic(); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	writeable, err := skillplan.GetWritablePlansPath()
	xlog.Logf("Writeable path is: %s", writeable)
}

// getPort returns the port the server should listen on.
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// setupRouter configures and returns the app’s router.
func setupRouter(secret string) *mux.Router {
	sessionStore := flyHandlers.NewSessionService(secret)
	r := mux.NewRouter()

	// Serve static files from the embedded filesystem
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// Utility and user functions
	r.HandleFunc("/callback/", flyHandlers.CallbackHandler(sessionStore))
	r.HandleFunc("/", flyHandlers.HomeHandler(sessionStore))
	r.HandleFunc("/login", flyHandlers.LoginHandler)
	r.HandleFunc("/auth-character", flyHandlers.AuthCharacterHandler)
	r.HandleFunc("/logout", flyHandlers.LogoutHandler(sessionStore))
	r.HandleFunc("/get-skill-plan", flyHandlers.SkillPlanFileHandler)
	r.HandleFunc("/save-skill-plan", flyHandlers.SaveSkillPlanHandler)
	r.HandleFunc("/delete-skill-plan", flyHandlers.DeleteSkillPlanHandler)

	// Admin routes
	r.HandleFunc("/reset-identities", flyHandlers.ResetIdentitiesHandler(sessionStore))

	return r
}

// startServer starts the HTTP server with CORS middleware.
func startServer(r *mux.Router, port string) {
	xlog.Logf("Listening on port %s", port)
	corsOptions := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsOptions)(r)))
}
