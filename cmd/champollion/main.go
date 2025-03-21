package main

import (
	"fmt"

	"github.com/PanacheTechnologies/Champollion/internal/search/searxng"
	"github.com/PanacheTechnologies/Champollion/pkg/env"
)

func main() {
	searxngURL := env.GetVar("SEARXNG_URL", "http://localhost:8080")
	client := searxng.NewClient(searxngURL)

	fmt.Println(client.Search("golang", nil))
}
