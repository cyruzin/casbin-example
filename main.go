package main

import (
	"context"
	"fmt"
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	// Initialize connection constants.
	HOST     = "localhost"
	DATABASE = "auth"
	USER     = "admin"
	PASSWORD = "secret"
)

func main() {
	ctx := context.Background()

	connectionInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, USER, PASSWORD, DATABASE,
	)

	// connect to the database first.
	db, err := sqlx.ConnectContext(ctx, "pgx", connectionInfo)
	if err != nil {
		panic(err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Panic(err)
	}

	// Initialize an adapter and use it in a Casbin enforcer:
	// The adapter will use the Postgres table name "casbin_rule_test",
	// the default table name is "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := sqladapter.NewAdapter(db.DB, "postgres", "casbin_rule_test")
	if err != nil {
		log.Panic(err)
	}

	// Create a configuration file. Head to https://casbin.org/editor/.
	// On this example I'm using the restful example.
	e, err := casbin.NewEnforcer("config/restful_model.conf", a)
	if err != nil {
		log.Panic(err)
	}

	// Load the policy from DB.
	if err = e.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
	}

	// Check the permission.
	has, err := e.Enforce("123", "/reports/*", "GET")
	if err != nil {
		log.Println("Enforce failed, err: ", err)
	}
	if !has {
		log.Println("do not have permission")
	}

	// Modify the policy.
	e.AddPolicy("123", "/reports/*", "GET")

	// Save the policy back to DB.
	if err = e.SavePolicy(); err != nil {
		log.Println("SavePolicy failed, err: ", err)
	}

	// Check the permission again.
	has, err = e.Enforce("123", "/reports/*", "GET")
	if err != nil {
		log.Println("Enforce failed, err: ", err)
	}
	if !has {
		log.Println("do not have permission")
	} else {
		log.Println("have permission")
	}

}
