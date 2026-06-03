// Package strava provides the API client for interfacing with Stava data
package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zwinslett/strava-cli-go/model"
)

type Client struct {
	httpClient   *http.Client
	clientId     string
	clientSecret string
	refreshToken string
	accessToken  string
	baseUrl      string
	authUrl      string
}

func NewClient() *Client {
	return &Client{
		httpClient:   &http.Client{},
		clientId:     os.Getenv("STRAVA_CLIENT_ID"),
		clientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		refreshToken: os.Getenv("STRAVA_REFRESH_TOKEN"),
		baseUrl:      "https://www.strava.com/api/v3",
		authUrl:      "https://www.strava.com/oauth/token",
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Client) SetAccessToken(ctx context.Context) error {
	data := url.Values{}
	data.Set("client_id", c.clientId)
	data.Set("client_secret", c.clientSecret)
	data.Set("refresh_token", c.refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.authUrl,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("failed: %d body=%s", resp.StatusCode, body.String())
	}

	var tr tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return err
	}
	c.accessToken = tr.AccessToken
	return nil
}

func (c *Client) doGet(ctx context.Context, url string, output any) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("failed: %d body=%s", resp.StatusCode, body.String())
	}
	return json.NewDecoder(resp.Body).Decode(output)
}

func (c *Client) GetActivityById(ctx context.Context, id int64) (model.DetailedActivity, error) {
	url := fmt.Sprintf("%s/activities/%d", c.baseUrl, id)

	var activity model.DetailedActivity
	err := c.doGet(ctx, url, &activity)
	if err != nil {
		return model.DetailedActivity{}, err
	}
	return activity, nil
}

func (c *Client) GetActivitiesByRange(ctx context.Context, after int64, before int64) ([]model.Activity, error) {
	rawUrl := fmt.Sprintf("%s/athlete/activities", c.baseUrl)
	params := url.Values{}
	params.Set("after", fmt.Sprintf("%d", after))
	params.Set("before", fmt.Sprintf("%d", before))
	params.Set("per_page", "200")
	fullUrl := rawUrl + "?" + params.Encode()

	var activities []model.Activity
	err := c.doGet(ctx, fullUrl, &activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (c *Client) GetRecentActivities(ctx context.Context, perPage int) ([]model.Activity, error) {
	rawUrl := fmt.Sprintf("%s/athlete/activities", c.baseUrl)
	params := url.Values{}
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	fullUrl := rawUrl + "?" + params.Encode()

	var activities []model.Activity
	err := c.doGet(ctx, fullUrl, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (c *Client) GetActivityZones(ctx context.Context, id int64) ([]model.Zones, error) {
	url := fmt.Sprintf("%s/activities/%d/zones", c.baseUrl, id)
	var zones []model.Zones
	err := c.doGet(ctx, url, &zones)
	if err != nil {
		return nil, err
	}
	return zones, nil
}
