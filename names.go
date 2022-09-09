package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func AddNames(client *spanner.Client, config map[string]string) error {
	table := config["table"]
	m := []*spanner.Mutation{
		spanner.Insert(table, []string{"id", "first_name"}, []interface{}{1, "John"}),
		spanner.Insert(table, []string{"id", "first_name"}, []interface{}{2, "Jane"}),
		spanner.Insert(table, []string{"id", "first_name"}, []interface{}{3, "Bob"}),
		spanner.Insert(table, []string{"id", "first_name"}, []interface{}{4, "Alice"}),
	}
	_, err := client.Apply(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func GetNames(client *spanner.Client, config map[string]string) ([]Name, error) {
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
			return nil, err
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
	return names, nil
}
