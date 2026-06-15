// vuelang — the Vuelang framework CLI.
//
// Usage:
//
//	vuelang make:model Product
//	vuelang make:controller ProductController
//	vuelang make:middleware AdminOnly
//	vuelang make:migration create_products_table
//	vuelang make:seeder ProductSeeder
//	vuelang serve
//	vuelang version
package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

const version = "2.0.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	cmd := os.Args[1]
	var name string
	if len(os.Args) > 2 {
		name = os.Args[2]
	}

	switch cmd {
	case "make:model":
		mustHaveName(name, "make:model")
		makeModel(name)
	case "make:controller":
		mustHaveName(name, "make:controller")
		makeController(name)
	case "make:middleware":
		mustHaveName(name, "make:middleware")
		makeMiddleware(name)
	case "make:migration":
		mustHaveName(name, "make:migration")
		makeMigration(name)
	case "make:seeder":
		mustHaveName(name, "make:seeder")
		makeSeeder(name)
	case "version", "-v", "--version":
		fmt.Printf("Vuelang CLI v%s\n", version)
	case "serve":
		fmt.Println("Use 'make dev' to start the development server.")
		fmt.Println("Use 'make run' to start the production binary.")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`Vuelang CLI v` + version + `

Usage:
  vuelang <command> [arguments]

Available commands:
  make:model       <Name>         Generate a model file
  make:controller  <Name>         Generate a controller file
  make:middleware  <Name>         Generate a middleware file
  make:migration   <name>         Generate a numbered migration file
  make:seeder      <Name>         Generate a seeder file
  serve                           Print server start instructions
  version                         Show CLI version

Examples:
  vuelang make:model Product
  vuelang make:controller ProductController
  vuelang make:middleware AdminOnly
  vuelang make:migration create_products_table
`)
}

func mustHaveName(name, cmd string) {
	if strings.TrimSpace(name) == "" {
		fmt.Fprintf(os.Stderr, "error: %s requires a name argument\n", cmd)
		os.Exit(1)
	}
}

// ── Generators ────────────────────────────────────────────────────────────────

func makeModel(name string) {
	path := fmt.Sprintf("app/models/%s.go", strings.ToLower(name))
	writeTemplate(path, modelTmpl, map[string]string{"Name": name})
}

func makeController(name string) {
	path := fmt.Sprintf("app/controllers/%s.go", name)
	writeTemplate(path, controllerTmpl, map[string]string{"Name": name})
}

func makeMiddleware(name string) {
	slug := strings.ToLower(name)
	path := fmt.Sprintf("app/middleware/%s.go", slug)
	writeTemplate(path, middlewareTmpl, map[string]string{"Name": name, "Slug": slug})
}

func makeMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	path := fmt.Sprintf("database/migrations/%s_%s.go", timestamp, name)
	writeTemplate(path, migrationTmpl, map[string]string{
		"FuncName": snakeToCamel(name),
		"Name":     name,
	})
	fmt.Printf("✓ Migration registered — add it to database/migrations/runner.go\n")
}

func makeSeeder(name string) {
	path := fmt.Sprintf("database/seeders/%s.go", name)
	writeTemplate(path, seederTmpl, map[string]string{"Name": name})
}

func writeTemplate(path, tmplStr string, data any) {
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "error: %s already exists\n", path)
		os.Exit(1)
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	t := template.Must(template.New("").Parse(tmplStr))
	if err := t.Execute(f, data); err != nil {
		fmt.Fprintf(os.Stderr, "template error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created %s\n", path)
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// ── Templates ─────────────────────────────────────────────────────────────────

const modelTmpl = `package models

import "time"

type {{.Name}} struct {
	ID        uint      ` + "`" + `json:"id"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
`

const controllerTmpl = `package controllers

import (
	"github.com/gin-gonic/gin"
	"vuelang/internal/framework/response"
)

type {{.Name}} struct{}

func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{}
}

func (ctrl *{{.Name}}) Index(c *gin.Context) {
	response.Success(c, nil, "OK")
}
`

const middlewareTmpl = `package middleware

import (
	"github.com/gin-gonic/gin"
	"vuelang/internal/framework/response"
)

// {{.Name}} middleware
func {{.Name}}() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement {{.Name}} logic
		_ = response.ServerError // remove this line
		c.Next()
	}
}
`

const migrationTmpl = `package migrations

import "database/sql"

func {{.FuncName}}(db *sql.DB) error {
	_, err := db.Exec(` + "`" + `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP
			                           ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	` + "`" + `)
	return err
}
`

const seederTmpl = `package seeders

import (
	"context"
	"database/sql"
	"log/slog"
)

func Seed{{.Name}}(db *sql.DB) error {
	ctx := context.Background()
	_ = ctx
	slog.Info("seeding {{.Name}}...")
	// TODO: add seed data
	return nil
}
`
