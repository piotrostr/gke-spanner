// Recently, I have seen that you can now create spanner instances with only
// 0 nodes, with only 100 processing units. That means it is very affordable.
package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

var ctx = context.Background()

// The main function pulls the data from the pre-created `some_table` from
// `database` and returns JSON
//
// This is a simple experiment to check the read
// speed and latency
func main() {
	db := "projects/piotrostr-resources/instances/big-db/databases/database"
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "im healthy!\n")
	})

	r.GET("/get-names", func(c *gin.Context) {
		tx := client.ReadOnlyTransaction()
		defer tx.Close()

		iter := tx.Query(
			ctx,
			spanner.NewStatement(`SELECT * FROM some_table`),
		)
		defer iter.Stop()

		var names []string
		i := 0
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			names = append(names, row.String())

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
