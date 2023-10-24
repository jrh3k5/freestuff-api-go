package http_test

import (
	"context"
	"time"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"net/http"

	_ "embed"

	freestuffhttp "github.com/jrh3k5/freestuff-api-go/pkg/client/v1/http"
)

//go:embed game_info_response_443200.json
var gameInfo443200JSON string

var _ = Describe("Client", func() {
	var apiClient *freestuffhttp.HTTPClient
	var apiKey string

	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()

		httpmock.Activate()
		DeferCleanup(httpmock.DeactivateAndReset)

		apiKey = "test-api-key"

		apiClient = freestuffhttp.NewHTTPClient(apiKey, &http.Client{})
	})

	Context("GetGameInfo", func() {
		It("should successfully deserialize the response", func() {
			httpmock.RegisterResponder(http.MethodGet, "https://api.freestuffbot.xyz/v1/game/443200/info", httpmock.NewStringResponder(http.StatusOK, gameInfo443200JSON))

			gameInfos, err := apiClient.GetGameInfo(ctx, []int64{443200})
			Expect(err).ToNot(HaveOccurred(), "getting the game info should not fail")
			Expect(gameInfos).To(HaveLen(1), "a single game info should be returned")

			gameInfo := gameInfos[0]
			Expect(gameInfo.ID).To(Equal(int64(443200)), "the game info ID should match")
			Expect(gameInfo.Title).To(Equal("The Evil Within"), "the game title should be correct")
			Expect(gameInfo.Description).ToNot(BeEmpty(), "the game description should be present")
			Expect(gameInfo.Until).ToNot(BeNil(), "the game info should have an 'until' timestamp")
			Expect(*gameInfo.Until).To(BeTemporally("==", time.Unix(1698332400, 0)), "the 'until' time should be correct")
			Expect(gameInfo.Kind).To(Equal("game"), "the game should be a free game")
			Expect(gameInfo.URLs.Default).To(Equal("https://redirect.freestuffbot.xyz/game/GYpU#The-Evil-Within"), "the default URL should be correct")
			Expect(gameInfo.URLs.Browser).To(Equal("https://redirect.freestuffbot.xyz/game/GYpU#The-Evil-Within"), "the browser URL should be correct")
			Expect(gameInfo.URLs.Org).To(Equal("https://store.epicgames.com/p/the-evil-within"), "the org URL should be correct")
			Expect(gameInfo.URLs.Client).ToNot(BeNil(), "a client URL should be present")
			Expect(*gameInfo.URLs.Client).To(Equal("https://redirect.freestuffbot.xyz/game/mviK"), "the client URL should be correct")
			Expect(gameInfo.Store).To(Equal("epic"), "the game info store should be correct")
		})
	})
})
