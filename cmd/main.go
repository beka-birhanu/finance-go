package main

import (
	"fmt"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/beka-birhanu/finance-go/api/graph"
	"github.com/beka-birhanu/finance-go/api/middleware"
	ratelimiter "github.com/beka-birhanu/finance-go/api/rate_limiter"
	api "github.com/beka-birhanu/finance-go/api/rest"
	"github.com/beka-birhanu/finance-go/api/rest/expense"
	"github.com/beka-birhanu/finance-go/api/rest/user"
	"github.com/beka-birhanu/finance-go/api/router"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	"github.com/beka-birhanu/finance-go/config"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/hash"
	"github.com/beka-birhanu/finance-go/infrastructure/jwt"
	expenserepo "github.com/beka-birhanu/finance-go/infrastructure/repository/expense"
	userrepo "github.com/beka-birhanu/finance-go/infrastructure/repository/user"
	timeservice "github.com/beka-birhanu/finance-go/infrastructure/time_service"
	"golang.org/x/time/rate"
)

// Global variables to hold database configuration
var (
	serverPort = config.Envs.ServerPort
	rateLimit  = config.Envs.APIRate
	rateBurst  = config.Envs.RateBurst
	dbUser     = config.Envs.DBUser
	dbPassword = config.Envs.DBPassword
	dbName     = config.Envs.DBName
	dbHost     = config.Envs.DBHost
	dbPort     = config.Envs.DBPort
)

func main() {
	// Connect to the database
	database := db.Connect(db.Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbHost:     dbHost,
		DbPort:     dbPort,
	})

	// Initialize services and repositories
	timeService := timeservice.New()
	userRepository := userrepo.New(database)
	expenseRepository := expenserepo.New(database)
	jwtService := initializeJWTService(timeService)
	hashService := hash.SingletonService()
	ipRateLimiter := ratelimiter.NewIPRateLimiter(rate.Limit(rateLimit), rateBurst, timeService)

	// Initialize middlewares
	authorizationMiddleware := middleware.Authorization(jwtService, true)
	populateClaimsMiddleware := middleware.Authorization(jwtService, false)
	rateLimitingMiddleware := middleware.RateLimitMiddleware(ipRateLimiter)

	// Initialize command and query handlers
	userRegisterCommandHandler := initializeUserRegisterHandler(userRepository, jwtService, hashService, timeService)
	userLoginQueryHandler := initializeUserLoginQueryHandler(userRepository, jwtService, hashService)
	addExpenseHandler := initializeAddExpenseHandler(userRepository, timeService)
	getExpenseHandler := initializeGetExpenseHandler(expenseRepository)
	getExpensesHandler := initializeGetExpensesHandler(expenseRepository)
	patchExpenseHandler := initializePatchExpenseHandler(expenseRepository)

	userHandler := user.NewHandler(user.Config{
		UserRepository:  userRepository,
		RegisterHandler: userRegisterCommandHandler,
		LoginHandler:    userLoginQueryHandler,
	})

	// Expense routes
	expenseHandler := expense.NewHandler(expense.Config{
		AddHandler:         addExpenseHandler,
		GetHandler:         getExpenseHandler,
		PatchHandler:       patchExpenseHandler,
		GetMultipleHandler: getExpensesHandler,
	})

	resolver := graph.NewResolver(graph.ResolverConfig{
		GetExpenseHandler:         getExpenseHandler,
		AddExpenseHandler:         addExpenseHandler,
		PatchExpenseHandler:       patchExpenseHandler,
	})

	graphHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Create and run the server
	server := router.NewRouter(router.Config{
		Addr:                     fmt.Sprintf(":%s", serverPort),
		RestfullControllers:      []api.IController{userHandler, expenseHandler},
		GraphQlController:        graphHandler,
		AuthorizationMiddleware:  authorizationMiddleware,
		PopulateClaimsMiddleware: populateClaimsMiddleware,
		RateLimitMiddleware:      rateLimitingMiddleware,
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initializePatchExpenseHandler(expenseRepository *expenserepo.Repository) *expensecmd.PatchHandler {
	return expensecmd.NewPatchHandler(expenseRepository)
}

func initializeGetExpenseHandler(expenseRepository *expenserepo.Repository) *expensqry.GetHandler {
	return expensqry.NewGetHandler(expenseRepository)
}

func initializeGetExpensesHandler(expenseRepository *expenserepo.Repository) *expensqry.GetMultipleHandler {
	return expensqry.NewGetMultipleHandler(expenseRepository)
}

// initializeJWTService initializes and returns a new JWT service instance.
func initializeJWTService(timeService *timeservice.Service) *jwt.Service {
	return jwt.New(
		jwt.Config{
			SecretKey:   config.Envs.JWTSecret,
			Issuer:      config.Envs.ServerHost,
			ExpTime:     time.Duration(config.Envs.JWTExpirationInSeconds) * time.Second,
			TimeService: timeService,
		})
}

// initializeUserRegisterHandler initializes and returns a new user register command handler.
func initializeUserRegisterHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service, timeService *timeservice.Service) *registercmd.Handler {
	return registercmd.NewHandler(registercmd.Config{
		UserRepo: userRepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
		TimeSvc:  timeService,
	})
}

// initializeUserLoginQueryHandler initializes and returns a new user login query handler.
func initializeUserLoginQueryHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service) *loginqry.Handler {
	return loginqry.NewHandler(loginqry.Config{
		UserRepository: userRepo,
		JwtService:     jwtService,
		HashService:    hashService,
	})
}

// initializeAddExpenseHandler initializes and returns a new add expense command handler.
func initializeAddExpenseHandler(userRepo *userrepo.Repository, timeService *timeservice.Service) *expensecmd.AddHandler {
	return expensecmd.NewAddHandler(expensecmd.Config{
		UserRepository: userRepo,
		TimeService:    timeService,
	})
}
