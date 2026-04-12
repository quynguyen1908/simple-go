package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang/internal/user"
	"golang/pkg/config"
	"golang/pkg/mailer"

	_ "golang/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBTimezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection established")

	db.AutoMigrate(&user.Role{}, &user.User{}, &user.UserProfile{}, &user.UserToken{})

	mailService := mailer.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	userRepo := user.NewUserRepository(db)

	fmt.Println("Seeding roles...")
	if err := userRepo.SeedRoles(context.Background()); err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	userService := user.NewUserService(userRepo, mailService)
	userHandler := user.NewUserHandler(userService, cfg.AppURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/users/register", userHandler.RegisterHandler)
	mux.HandleFunc("/api/v1/users/confirm-email", userHandler.ConfirmEmailHandler)
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	fmt.Printf("Server is running on port %s\n", cfg.Port)
	fmt.Printf("Swagger docs available at http://localhost:%s/docs/index.html\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
