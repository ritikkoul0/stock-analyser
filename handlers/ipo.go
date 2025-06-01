package handlers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"stock-analyser/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Base URL for Groww IPO API
const growAPIBaseURL = "https://groww.in/v1/api/primaries/v1/ipo/"

// Data models for IPOs

type UPCOMINGIPO struct {
	Symbol            string  `json:"symbol"`
	SearchID          string  `json:"searchId"`
	CompanyName       string  `json:"companyName"`
	IsSME             bool    `json:"isSme"`
	BidStartTimestamp *int64  `json:"bidStartTimestamp,omitempty"` // nullable field
	LogoURL           string  `json:"logoUrl"`
	DocumentURL       *string `json:"documentUrl,omitempty"` // nullable field
}

type CLOSEDIPO struct {
	Symbol              string  `json:"symbol"`
	SearchID            string  `json:"searchId"`
	CompanyName         string  `json:"companyName"`
	IsSME               bool    `json:"isSme"`
	IssuePrice          float64 `json:"issuePrice"`
	ListingPrice        float64 `json:"listingPrice"`
	ListingTimestamp    int64   `json:"listingTimestamp"`
	IsListed            bool    `json:"isListed"`
	LogoURL             string  `json:"logoUrl"`
	OverallSubscription float64 `json:"overallSubscription"`
	ListingReturn       float64 `json:"listingReturn"`
}

type Category struct {
	Category           string  `json:"category"`
	BidCutOffTimestamp int64   `json:"bidCutOffTimestamp"`
	CategoryLabel      string  `json:"categoryLabel"`
	CategorySubText    string  `json:"categorySubText"`
	LotSize            int     `json:"lotSize"`
	MinBidQuantity     int     `json:"minBidQuantity"`
	MinPrice           float64 `json:"minPrice"`
	MaxPrice           float64 `json:"maxPrice"`
}

type OPENIPO struct {
	Symbol              string     `json:"symbol"`
	CompanyCode         *float64   `json:"companyCode,omitempty"` // nullable
	SearchID            string     `json:"searchId"`
	ISIN                *string    `json:"isin,omitempty"` // nullable
	CompanyName         string     `json:"companyName"`
	IsSME               bool       `json:"isSme"`
	IssuePrice          *float64   `json:"issuePrice,omitempty"`        // nullable
	ListingPrice        *float64   `json:"listingPrice,omitempty"`      // nullable
	ListingTimestamp    *int64     `json:"listingTimestamp,omitempty"`  // nullable
	IsListed            *bool      `json:"isListed,omitempty"`          // nullable
	BidStartTimestamp   *int64     `json:"bidStartTimestamp,omitempty"` // nullable
	BidEndTimestamp     *int64     `json:"bidEndTimestamp,omitempty"`   // nullable
	OverallSubscription float64    `json:"overallSubscription"`
	ListingReturn       *float64   `json:"listingReturn,omitempty"` // nullable
	LogoURL             string     `json:"logoUrl"`
	IsPreApply          *bool      `json:"isPreApply,omitempty"` // nullable
	TickSize            *int       `json:"tickSize,omitempty"`   // nullable
	Categories          []Category `json:"categories,omitempty"` // nullable slice
}

// GetUpcomingIPOs handles fetching upcoming IPOs and returning JSON response
func GetUpcomingIPOs(ctx *gin.Context) {
	url := growAPIBaseURL + "upcoming"

	type IPOResponse struct {
		IPOList []UPCOMINGIPO `json:"ipoList"`
	}

	response, err := fetchIPOData[IPOResponse](ctx.Request.Context(), url)
	if err != nil {
		logger.Info("GetUpcomingIPOs error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.IPOList)
}

// GetOpenIPOs handles fetching open IPOs and returning JSON response
func GetOpenIPOs(ctx *gin.Context) {
	url := growAPIBaseURL + "open"

	type IPOResponse struct {
		IPOList []OPENIPO `json:"ipoList"`
	}

	response, err := fetchIPOData[IPOResponse](ctx.Request.Context(), url)
	if err != nil {
		logger.Info("GetOpenIPOs error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.IPOList)
}

// GetClosedIPOs handles fetching closed IPOs and returning JSON response
func GetClosedIPOs(ctx *gin.Context) {
	url := growAPIBaseURL + "closed"

	type IPOResponse struct {
		IPOList []CLOSEDIPO `json:"ipoList"`
	}

	response, err := fetchIPOData[IPOResponse](ctx.Request.Context(), url)
	if err != nil {
		logger.Info("GetClosedIPOs error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.IPOList)
}

// HTTP client with configurable TLS verification
var client *http.Client

func init() {
	skipVerify := os.Getenv("SKIP_SSL_VERIFY") == "true"
	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipVerify,
			},
		},
	}
}

// fetchIPOData fetches JSON from given URL and decodes into type T.
// Uses context for cancellation and timeout.
func fetchIPOData[T any](ctx context.Context, url string) (T, error) {
	var result T

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to fetch IPO data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("GROWW API returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return result, nil
}
