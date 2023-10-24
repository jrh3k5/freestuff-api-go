package v1

import "context"

// GameCategory is an enumeration of supported game categories by the API.
type GameCategory string

const (
	// GameCategoryAll describes all games known to the API
	GameCategoryAll GameCategory = "all"
	// GameCategoryApproved describes all games that have been manually approved for listing
	GameCategoryApproved GameCategory = "approved"
	// GameCategoryFree describes all known free games
	GameCategoryFree GameCategory = "free"
)

type Client interface {
	// GetGameIDs gets the IDs of the games that fall within the given GameCategory.
	GetGameIDs(ctx context.Context, gameCategory GameCategory) ([]int64, error)

	// GetGameInfo gets game information for the given IDs.
	GetGameInfo(ctx context.Context, gameIDs []int64) ([]GameInfo, error)

	// Ping pings the API.
	Ping(ctx context.Context) error
}
