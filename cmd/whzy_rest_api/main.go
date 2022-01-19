package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"net"
	"net/http"
	"os"

	gounits "github.com/docker/go-units"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	_ "github.com/godror/godror"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/moihn/whzycrmgo/oapi"
)

func oapiDocHandler(endpoint string, swagger *openapi3.T) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.EqualFold(r.URL.Path, endpoint) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(swagger)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func logInit(logLevelName string) {
	// Only log the warning severity or above.
	if logLevel, err := logrus.ParseLevel(logLevelName); err != nil {
		logrus.Fatal("Fail to parse logging level: ", logLevel)
	} else {
		logrus.SetLevel(logLevel)
		logrus.SetOutput(os.Stderr)
		logrus.SetReportCaller(true)
		// logrus.SetFormatter(&formatter.InfFormatter{})
	}
}

type Config struct {
	DbConnectString *string `yaml:"DbConnectString"`
	ListenPort      *int    `yaml:"ListenPort"`
}

func main() {
	var port, timeout int
	var dbConnString string
	var bodyMaxStr string
	var logLevelName, configFile string
	var printVersion bool

	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.IntVar(&port, "port", 0, "Optional port on which we listen. Dynamically chosen if not supplied. Useful when testing")
	flag.StringVar(&dbConnString, "dbConnectString", "", "Oracle database connection string.")
	flag.IntVar(&timeout, "timeout", 600, "Timeout for forwarding connection")
	flag.StringVar(&logLevelName, "log-level", "warn", "Log level")
	flag.StringVar(&bodyMaxStr, "bodySizeMax", "512MB", "Max size for request body in huamn readable format, like 10MB")
	flag.BoolVar(&printVersion, "version", false, "Show version of the executable")
	flag.Parse()

	if len(configFile) > 0 {
		configString, err := os.ReadFile(configFile)
		if err != nil {
			logrus.Fatalf("failed to read file %v: %v", configFile, err)
		}
		var config Config
		err = yaml.Unmarshal(configString, &config)
		if err != nil {
			logrus.Fatalf("failed to parse configuration file %v: %v", configFile, err)
		}

		if len(dbConnString) == 0 && config.DbConnectString != nil {
			dbConnString = *config.DbConnectString
		}

		if config.ListenPort != nil && port == 0 {
			port = *config.ListenPort
		}
	}

	if printVersion {
		// fmt.Println("Version: {{x-execVersionString}}")
		return
	}

	logInit(logLevelName)

	bodyMax, err := gounits.FromHumanSize(bodyMaxStr)
	if err != nil {
		logrus.Fatalf("Error in parsing max body size: %v", err)
	}
	swaggerForValidation, _ := oapi.GetSwagger()
	swaggerForUi, err := oapi.GetSwagger()
	if err != nil {
		logrus.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swaggerForValidation.Servers = nil
	swaggerForUi.Servers = nil

	// Clear out the security array in the swagger validation spec, so if the forwareded request
	// carries no authentication information, it will not throw an error.
	// It is still needed by SwaggerUI to render correct authentication control.
	swaggerForValidation.Security = nil

	// open database
	logrus.Debugf("DbConnectionString: %v", dbConnString)
	db, err := sql.Open("godror", dbConnString)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}

	// Create an instance of our handler which satisfies the generated interface
	handler := oapi.NewServiceHandler(timeout, bodyMax, db)

	// This is how you set up a basic chi router
	oapiRouter := chi.NewRouter()

	// Use oapi document serving middleware to handle doc request
	oapiRouter.Use(oapiDocHandler("/api-doc", swaggerForUi))

	// Use our validation middleware to check all requests against the
	// OpenAPI schema. Disable format checking as OAS 3 accept them for
	// docummentation purpose.
	openapi3.SchemaFormatValidationDisabled = true
	oapiRouter.Use(middleware.OapiRequestValidator(swaggerForValidation))

	// We now register our handler above as the handler for the interface
	oapi.HandlerFromMux(handler, oapiRouter)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logrus.Fatal(err)
	}
	port = listener.Addr().(*net.TCPAddr).Port
	logrus.Println("Listening on port", port)

	s := &http.Server{
		Handler: oapiRouter,
	}

	ticker := time.NewTicker(time.Hour)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for {
			select {
			case oscall := <-c:
				logrus.Infof("system call:%+v", oscall)
				ticker.Stop()
				if err := s.Shutdown(context.Background()); err != nil {
					// Error from closing listeners, or context timeout:
					logrus.Fatalf("HTTP server Shutdown: %v", err)
				}
			case t := <-ticker.C:
				logrus.Infof("Ping database at %v", t)
				if err := db.Ping(); err != nil {
					logrus.Fatal(err)
				}
			}
		}
	}()

	// And we serve HTTP until the world ends.
	err = s.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Http server Serve() error: %v", err)
	} else if err == http.ErrServerClosed {
		logrus.Info("Http server is shutdown")
	}
}
