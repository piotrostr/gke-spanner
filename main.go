// Recently, I have seen that you can now create spanner instances with only
// 0 nodes, with only 100 processing units. That means it is very affordable.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

var ctx = context.Background()

type Name struct {
	Id        int64  `spanner:"id" json:"id"`
	FirstName string `spanner:"first_name" json:"first_name"`
}

// The main function pulls the data from the pre-created `some_table` from
// `database` and returns JSON
//
// This is a simple experiment to check the read speed and latency
func main() {
	config := map[string]string{
		"instance": "big-db",
		"project":  "piotrostr-resources",
		"database": "database",
		"table":    "some_table",
	}
	db := fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		config["project"],
		config["instance"],
		config["database"],
	)
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"config": config,
		})
	})

	r.GET("/get-names", func(c *gin.Context) {
		tx := client.ReadOnlyTransaction()
		defer tx.Close()

		iter := tx.Query(
			ctx,
			spanner.NewStatement(
				fmt.Sprintf(
					`SELECT * FROM %s`,
					config["table"],
				),
			),
		)
		defer iter.Stop()

		i := 0
		var names []Name
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				log.Println(err)
			}

			var ptr Name
			err = row.ToStruct(&ptr)
			if err != nil {
				log.Println(err)
			}

			names = append(names, ptr)

			// the i is in order not to pull too much data
			i += 1
			if i > 100 {
				break
			}

		}
		c.JSON(http.StatusOK, names)
	})

	r.Run(":8080")
}
