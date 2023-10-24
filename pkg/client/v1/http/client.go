package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	v1 "github.com/jrh3k5/freestuff-api-go/pkg/client/v1"
	freestuffhttp "github.com/jrh3k5/freestuff-api-go/pkg/http"
)

const (
	baseURL = "https://api.freestuffbot.xyz/v1"
)

// HTTPClient is an HTTP-backed implementation of v1.Client
type HTTPClient struct {
	apiKey string
	doer   freestuffhttp.Doer
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient(apiKey string, doer freestuffhttp.Doer) *HTTPClient {
	return &HTTPClient{
		apiKey: apiKey,
		doer:   doer,
	}
}

func (h *HTTPClient) GetGameIDs(ctx context.Context, gameCategory v1.GameCategory) ([]int64, error) {
	response, err := h.executeGet(ctx, "/games/"+string(gameCategory))
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to get games under category '%s': %w", string(gameCategory), err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, h.getResponseError(response)
	}

	var gameList *gameListResponse
	if decodeErr := json.NewDecoder(response.Body).Decode(&gameList); decodeErr != nil {
		return nil, fmt.Errorf("failed to decode game list response body: %w", decodeErr)
	}

	return gameList.Data, nil
}

func (h *HTTPClient) GetGameInfo(ctx context.Context, gameIDs []int64) ([]v1.GameInfo, error) {
	gameIDsStrings := make([]string, len(gameIDs))
	for gameIndex, gameID := range gameIDs {
		gameIDsStrings[gameIndex] = fmt.Sprintf("%d", gameID)
	}
	response, err := h.executeGet(ctx, "/game/"+strings.Join(gameIDsStrings, "+")+"/info")
	if err != nil {
		return nil, fmt.Errorf("failed to execute retrieval of game info for %d game IDs [%s]: %w", len(gameIDs), strings.Join(gameIDsStrings, ", "), err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var gameInfosResponse *gameInfoResponse
	if decodeErr := json.NewDecoder(response.Body).Decode(&gameInfosResponse); decodeErr != nil {
		return nil, fmt.Errorf("failed to decode reponse for reading %d game IDs [%s]: %w", len(gameIDs), strings.Join(gameIDsStrings, ", "), decodeErr)
	}

	gameInfos := make([]v1.GameInfo, 0, len(gameInfosResponse.Data))
	for _, gameInfoResponse := range gameInfosResponse.Data {
		var until *time.Time
		if untilSecs := gameInfoResponse.Until; untilSecs != nil {
			timeVal := time.Unix(*untilSecs, 0)
			until = &timeVal
		}

		gameInfo := v1.GameInfo{
			ID:          gameInfoResponse.ID,
			Title:       gameInfoResponse.Title,
			Description: gameInfoResponse.Description,
			Until:       until,
			Kind:        gameInfoResponse.Kind,
			Store:       gameInfoResponse.Store,
			URLs: v1.GameInfoURLs{
				Default: gameInfoResponse.URLs.Default,
				Browser: gameInfoResponse.URLs.Browser,
				Org:     gameInfoResponse.URLs.Org,
				Client:  &gameInfoResponse.URLs.Client,
			},
		}

		gameInfos = append(gameInfos, gameInfo)
	}

	return gameInfos, nil
}

func (h *HTTPClient) Ping(ctx context.Context) error {
	pingResponse, err := h.executeGet(ctx, "/ping")
	if err != nil {
		return fmt.Errorf("failed to execute ping: %w", err)
	}
	defer func() {
		_ = pingResponse.Body.Close()
	}()

	if pingResponse.StatusCode != http.StatusOK {
		return h.getResponseError(pingResponse)
	}

	return nil
}

func (h *HTTPClient) executeGet(ctx context.Context, requestURI string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build GET request to request URI '%s': %w", requestURI, err)
	}
	request.Header.Set("Authorization", "Basic "+h.apiKey)

	response, err := h.doer.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request against request URI '%s': %w", requestURI, err)
	}
	return response, nil
}

func (h *HTTPClient) getResponseError(response *http.Response) error {
	var errorResponse *errorResponse
	if decodeErr := json.NewDecoder(response.Body).Decode(&errorResponse); decodeErr != nil {
		return fmt.Errorf("unexpected response status code %d; however, the error could not be obtained from the response, so this error is incomplete", response.StatusCode)
	}

	return fmt.Errorf("unexpected response status code %d; response message was: '%s'", response.StatusCode, errorResponse.Message)
}
