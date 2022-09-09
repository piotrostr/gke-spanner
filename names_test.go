package main

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/spanner"
)

func TestAddNames(t *testing.T) {
	client, err := spanner.NewClient(context.TODO(), SpannerURL)
	if err != nil {
		t.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()
	_ = AddNames(client, Config)
}

func TestGetNames(t *testing.T) {
	client, err := spanner.NewClient(context.TODO(), SpannerURL)
	if err != nil {
		t.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()
	names := GetNames(client, Config)
	fmt.Println(names)
}
