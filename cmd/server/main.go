package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/Stuko0/scarlet-backend/gen/proto/entity/v1/entityv1connect"
	"github.com/Stuko0/scarlet-backend/gen/proto/user/v1/userv1connect"
	"github.com/Stuko0/scarlet-backend/gen/proto/wildfire/v1/nrtv1connect"
	"github.com/Stuko0/scarlet-backend/internal/auth"
	"github.com/Stuko0/scarlet-backend/internal/database"
	"github.com/Stuko0/scarlet-backend/internal/domain/entity"
	"github.com/Stuko0/scarlet-backend/internal/domain/user"
	"github.com/Stuko0/scarlet-backend/internal/domain/wildfire/current"
	"github.com/Stuko0/scarlet-backend/internal/domain/wildfire/scrapers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {}

func main() {
	err:= godotenv.Load()
	if err!=nil{log.Fatal("Error loading .env file: ", err)}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatalf("failed to create pool, %v", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err !=nil{
		log.Fatalf("Failed to ping postgresql database: %v", err)
	}

	mongoClient, err:= database.NewMongoClient(ctx)
	if err!=nil{
		log.Printf("failed to connect to MongoDB: %v", err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	wildfireDB:=mongoClient.Database("scarlet")
	redisClient := database.NewRedisClient()
	defer redisClient.Close()
	repo:=wildfire.NewCachedRepository(wildfire.NewMongoRepository(wildfireDB), redisClient)
	nasaScraper:=scrapers.NewFIRMSScraper(os.Getenv("NASA_API_KEY"), "-69.38,-22.53,-57.26,-9.38") // Bolivia bounding box
	weatherScraper:=scrapers.NewOpenMeteoScraper()

	wildfireSvc:=wildfire.NewWildfireNRTService(repo, redisClient,[]scrapers.Scraper{nasaScraper}, weatherScraper,)
	wildfireSvc.TriggerImmediateScrape(ctx)
	go wildfireSvc.RunScrapers(ctx)

	db := &database.Postgres{Pool: pool}
	privateKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	publicKey := []byte(os.Getenv("JWT_PUBLIC_KEY"))
	jwtManager, err := auth.NewJWTManager(privateKey, publicKey, 24*time.Hour)
	if err != nil{log.Fatalf("failed to create JWT manager: %v",err)}

	userRepo:=user.NewUserRepository(db)
	userService:=user.NewUserService(userRepo,jwtManager)

	mux:=http.NewServeMux()
	path, handler := userv1connect.NewUserServiceHandler(userService)
	mux.Handle(path, handler)

	checker:=grpchealth.NewStaticChecker(userv1connect.UserServiceName)
	mux.Handle(grpchealth.NewHandler(checker))

	entityRepo := entity.NewEntitieRepository(db)
	entityService := entity.NewEntityService(entityRepo)
	pathEntities, handlerEntities := entityv1connect.NewEntityServiceHandler(entityService)
	mux.Handle(pathEntities, handlerEntities)

	pathFires, handlerFires := nrtv1connect.NewNRTServiceHandler(
		wildfireSvc,
		connect.WithCompressMinBytes(1024),
	)
	mux.Handle(pathFires, handlerFires)

	server := &http.Server{
		Addr: ":8000",
		Handler: h2c.NewHandler(
			mux, &http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 *time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(){
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed){
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit:= make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err !=nil{
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}