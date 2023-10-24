package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	v1 "github.com/jrh3k5/freestuff-api-go/pkg/client/v1"
	api "github.com/jrh3k5/freestuff-api-go/pkg/client/v1/http"
)

func main() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	apiKey := os.Getenv("FREESTUFF_API_KEY")
	if apiKey == "" {
		panic("FREESTUFF_API_KEY must be set")
	}

	freestuffClient := api.NewHTTPClient(apiKey, http.DefaultClient)

	if pingErr := freestuffClient.Ping(ctx); pingErr != nil {
		panic(fmt.Sprintf("failed to ping API: %v", pingErr))
	}

	freeGameIDs, err := freestuffClient.GetGameIDs(ctx, v1.GameCategoryFree)
	if err != nil {
		panic(fmt.Sprintf("failed to get free game IDs: %v", err))
	} else if len(freeGameIDs) == 0 {
		panic("No game IDs returned; this may indicate a deserialization issue")
	}

	gameInfos, err := freestuffClient.GetGameInfo(ctx, freeGameIDs)
	if err != nil {
		panic(fmt.Sprintf("failed to get game info: %v", err))
	} else if len(gameInfos) != len(freeGameIDs) {
		panic(fmt.Sprintf("asked for %d IDs, got %d games", len(freeGameIDs), len(gameInfos)))
	}

	for _, gameInfo := range gameInfos {
		if gameInfo.Title == "" {
			panic(fmt.Sprintf("Game ID %d has no title; this may indicate a deserialization issue", gameInfo.ID))
		}
	}
}
