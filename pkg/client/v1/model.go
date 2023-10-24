package v1

import "time"

// GameInfo describes the information about a game.
type GameInfo struct {
	ID          int64        // ID is the ID identifying the game in Freestuff
	Title       string       // Title is the title of the game
	Description string       // Description is a textual description of the game
	Until       *time.Time   // Until describes the point at which the game will no longer be free; if nil, the game will be free forever
	Kind        string       // Kind is the type of product (e.g., "free")
	URLs        GameInfoURLs // URLs describe the locations at which the game's information is viewable
	Store       string       // Store is the identifier of the store at which the game's information is available

}

// GameInfoURLs describes URLs for a game.
type GameInfoURLs struct {
	Default string  // Default is the default URL to be used to view the game
	Browser string  // Browser is the URL to which one can navigate a browser to look at the game
	Org     string  // Org is the original URL, bypassing any kind of analytics or referrer information
	Client  *string // Client is the optional URL that, if applicable, opens the game in the relevant client app
}
