// routes/api.go
//
// This is the single file where every API route is declared.
// Same idea as Laravel's routes/api.php.
//
// How to add a new resource:
//  1. Create  app/models/Product.go
//  2. Create  app/controllers/ProductController.go
//  3. Add the routes below — that's it.

package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"go-cloud-erp/app/controllers"
	"go-cloud-erp/app/middleware"
)

func Register(api *gin.RouterGroup, db *sql.DB) {

	// ── Controllers ───────────────────────────────────────────────────────────
	greeting := &controllers.GreetingController{}
	users    := &controllers.UserController{DB: db}

	// ── Public routes ─────────────────────────────────────────────────────────
	api.GET("/greeting", greeting.Index)

	// ── Protected routes (requires Bearer token) ──────────────────────────────
	auth := api.Group("/", middleware.Auth())
	{
		auth.GET   ("/users",        users.Index)
		auth.GET   ("/users/:id",    users.Show)
		auth.POST  ("/users",        users.Store)
		auth.PUT   ("/users/:id",    users.Update)
		auth.DELETE("/users/:id",    users.Destroy)

		// Add new resources here:
		// products := &controllers.ProductController{DB: db}
		// auth.GET   ("/products",        products.Index)
		// auth.GET   ("/products/:id",    products.Show)
		// auth.POST  ("/products",        products.Store)
		// auth.PUT   ("/products/:id",    products.Update)
		// auth.DELETE("/products/:id",    products.Destroy)
	}
}
