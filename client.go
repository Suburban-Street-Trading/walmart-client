package walmart

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Client struct {
	baseURL, clientID, clientSecret, accessToken string
	tokenExpiresAt                               int64
}

// NewClient creates a new Walmart API client.
func NewClient(baseURL, clientID, clientSecret string) *Client {
	return &Client{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *Client) acquireAccessToken(ctx context.Context) error {

	endpoint := "/v3/token"

	var resp AccessTokenResponse

	err := requests.URL(c.baseURL+endpoint).
		BasicAuth(c.clientID, c.clientSecret).
		ContentType("application/x-www-form-urlencoded").
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Param("grant_type", "client_credentials").
		Method(http.MethodPost).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return fmt.Errorf("could not acquire access token: %w", err)
	}

	c.accessToken = resp.AccessToken
	c.tokenExpiresAt = time.Now().Unix() + int64(resp.ExpiresIn)

	return nil
}

func (c *Client) ensureAccessToken(ctx context.Context) error {

	if c.accessToken == "" || c.tokenExpiresAt <= time.Now().Unix() {
		if err := c.acquireAccessToken(ctx); err != nil {
			return fmt.Errorf("could not ensure access token: %w", err)
		}
	}

	return nil
}

func (c *Client) GetAllReleasedOrders(ctx context.Context) ([]Order, error) {

	err := c.ensureAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get all released orders: %w", err)
	}

	endpoint := "/v3/orders/released"

	var resp GetOrdersResponse

	err = requests.URL(c.baseURL+endpoint).
		Transport(requests.Record(nil, "walmart/requests")).
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Header("WM_SEC.ACCESS_TOKEN", c.accessToken).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not get all released orders: %w", err)
	}

	return resp.List.Elements.Order, nil
}

func (c *Client) AcknowledgeOrder(ctx context.Context, purchaseOrderID string) (*Order, error) {

	err := c.ensureAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not acknowledge order: %w", err)
	}

	endpoint := fmt.Sprintf("/v3/orders/%s/acknowledge", purchaseOrderID)

	var resp struct {
		Order Order `json:"order"`
	}

	err = requests.URL(c.baseURL+endpoint).
		ContentType("application/json").
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Header("WM_SEC.ACCESS_TOKEN", c.accessToken).
		Method(http.MethodPost).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not acknowledge order: %w", err)
	}

	return &resp.Order, nil
}

type GetOrdersOpts struct {
	SKU              string `url:"sku,omitempty"`
	PurchaseOrderID  string `url:"purchaseOrderId,omitempty"`
	Status           string `url:"status,omitempty"`
	CreatedStartDate string `url:"createdStartDate,omitempty"`
	CreatedEndDate   string `url:"createdEndDate,omitempty"`
}

func (c *Client) GetOrders(ctx context.Context, opts GetOrdersOpts) ([]Order, error) {

	err := c.ensureAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get orders: %w", err)
	}

	endpoint := "/v3/orders"

	queryParams, err := query.Values(opts)
	if err != nil {
		return nil, fmt.Errorf("could not get orders: %w", err)
	}

	var resp GetOrdersResponse

	err = requests.URL(c.baseURL+endpoint+"?"+queryParams.Encode()).
		Transport(requests.Record(nil, "walmart/requests")).
		ContentType("application/json").
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Header("WM_SEC.ACCESS_TOKEN", c.accessToken).
		ToJSON(&resp).
		Fetch(ctx)

	return resp.List.Elements.Order, nil

}

func (c *Client) GetSupportedCarriers(ctx context.Context) ([]Carrier, error) {

	err := c.ensureAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get supported carriers: %w", err)
	}

	endpoint := "/v3/shipping/labels/carriers"

	var resp SupportedCarriersResponse

	err = requests.URL(c.baseURL+endpoint).
		ContentType("application/json").
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Header("WM_SEC.ACCESS_TOKEN", c.accessToken).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not get supported carriers: %w", err)
	}

	return resp.Carriers, nil

}

func (c *Client) ShipOrderLines(ctx context.Context, purchaseOrderID string, shipment ShipOrderLinesRequest) (*Order, error) {

	err := c.ensureAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not ship order lines: %w", err)
	}

	endpoint := fmt.Sprintf("/v3/orders/%s/shipping", purchaseOrderID)

	var resp struct {
		Order Order `json:"order"`
	}

	err = requests.URL(c.baseURL+endpoint).
		Transport(requests.Record(nil, "walmart/requests")).
		ContentType("application/json").
		Accept("application/json").
		Header("WM_QOS.CORRELATION_ID", uuid.New().String()).
		Header("WM_SVC.NAME", "golang-backend").
		Header("WM_SEC.ACCESS_TOKEN", c.accessToken).
		BodyJSON(shipment).
		Method(http.MethodPost).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not ship order lines: %w", err)
	}

	return &resp.Order, nil

}
