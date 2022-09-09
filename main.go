// Recently, I have seen that you can now create spanner instances with only
// 0 nodes, with only 100 processing units. That means it is very affordable.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

type Name struct {
	Id        int64  `spanner:"id" json:"id"`
	FirstName string `spanner:"first_name" json:"first_name"`
}

var Config = map[string]string{
	"instance": "database",
	"project":  "piotrostr-resources",
	"database": "database",
	"table":    "some_table",
}

var SpannerURL = fmt.Sprintf(
	"projects/%s/instances/%s/databases/%s",
	Config["project"],
	Config["instance"],
	Config["database"],
)

func CreateTable() {
	admin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer admin.Close()
	op, err := admin.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database: SpannerURL,
		Statements: []string{
			`CREATE TABLE some_table (
          id INT64,
          first_name STRING(100)
       ) PRIMARY KEY (id)`,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := op.Wait(ctx); err != nil {
		if err.Error() == "rpc error: code = FailedPrecondition desc = Duplicate name in schema: some_table." {
			log.Println("Table already exists")
			return
		} else {
			fmt.Print(err)
			log.Fatal(err)
		}
	}
	fmt.Println("Table some_table was created")
}

// The main function pulls the data from the pre-created `some_table` from
// `database` and returns JSON
//
// This is a simple experiment to check the read speed and latency
func main() {
	shallCreateTable := flag.Bool("create-table", false, "Create table")
	flag.Parse()

	if *shallCreateTable {
		CreateTable()
	}

	client, err := spanner.NewClient(ctx, SpannerURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"config": Config,
		})
	})

	r.POST("/add-names", func(c *gin.Context) {
		err := AddNames(client, Config)
		if err == nil {
			c.JSON(http.StatusCreated, gin.H{
				"status": "ok",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err,
			})
		}
	})

	r.GET("/get-names", func(c *gin.Context) {
		names := GetNames(client, Config)
		c.JSON(http.StatusOK, names)
	})

	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
