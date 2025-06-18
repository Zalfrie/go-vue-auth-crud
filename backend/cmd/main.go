package main

import (
	"go-vue-auth-crud/config"
	"go-vue-auth-crud/routes"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	_ "go-vue-auth-crud/docs" // Swagger generated docs
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go-Vue Auth CRUD API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()
	db := config.ConnectDB(cfg)
	defer db.Close()

	// CLI commands: migrate, seed, serve
	switch os.Args[1] {
		case "migrate":
			// migrations run in ConnectDB
			fmt.Println("Migration selesai")
			return
		case "seed":
			if err := seed.Run(db); err != nil {
				log.Fatal("Seed error:", err)
			}
			fmt.Println("Seed data selesai")
			return
		case "serve":
			// start HTTP server
		default:
			log.Fatal("Unknown command. Use migrate, seed, or serve.")
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db, cfg)

	// Swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swagFiles.Handler))

	addr := ":" + cfg.AppPort
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}