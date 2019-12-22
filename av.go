package av

import (
	"context"
	"net/url"
	"time"
)

const (
	// HostDefault is the default host for Alpha Vantage
	HostDefault    = "www.alphavantage.co"
	TimeoutDefault = time.Second * 30
)

const (
	schemeHttps = "https"

	queryApiKey     = "apikey"
	queryDataType   = "datatype"
	queryOutputSize = "outputsize"
	querySymbol     = "symbol"
	queryMarket     = "market"
	queryEndpoint   = "function"
	queryInterval   = "interval"

	valueCompact                 = "compact"
	valueJson                    = "csv"
	valueDigitalCurrencyEndpoint = "DIGITAL_CURRENCY_INTRADAY"

	pathQuery = "query"
)

// Client is a service used to query Alpha Vantage stock data
type Client struct {
	copts clientOptions
}

func defaultClientOptions() clientOptions {
	return clientOptions{
		apiKey: "",
		conn:   NewConnection(),
	}
}

// NewClientConnection creates a new Client with the default Alpha Vantage connection
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		copts: defaultClientOptions(),
	}

	for _, opt := range opts {
		opt.apply(&c.copts)
	}

	return c
}

func (c *Client) Conn() Connection {
	return c.copts.conn
}

// buildRequestPath builds an endpoint URL with the given query parameters
func (c *Client) buildRequestPath(params map[string]string) *url.URL {
	// build our URL
	endpoint := &url.URL{}
	endpoint.Path = pathQuery

	// base parameters
	query := endpoint.Query()
	query.Set(queryApiKey, c.copts.apiKey)
	query.Set(queryDataType, valueJson)
	query.Set(queryOutputSize, valueCompact)

	// additional parameters
	for key, value := range params {
		query.Set(key, value)
	}

	endpoint.RawQuery = query.Encode()

	return endpoint
}

// StockTimeSeriesIntraday queries a stock symbols statistics throughout the day.
// Data is returned from past to present.
func (c *Client) StockTimeSeriesIntraday(ctx context.Context, timeInterval TimeInterval, symbol string) ([]*TimeSeriesValue, error) {
	endpoint := c.buildRequestPath(map[string]string{
		queryEndpoint: timeSeriesIntraday.keyName(),
		queryInterval: timeInterval.keyName(),
		querySymbol:   symbol,
	})
	response, err := c.Conn().Request(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return parseTimeSeriesData(response.Body)
}

// StockTimeSeries queries a stock symbols statistics for a given time frame.
// Data is returned from past to present.
func (c *Client) StockTimeSeries(ctx context.Context, timeSeries TimeSeries, symbol string) ([]*TimeSeriesValue, error) {
	endpoint := c.buildRequestPath(map[string]string{
		queryEndpoint: timeSeries.keyName(),
		querySymbol:   symbol,
	})
	response, err := c.Conn().Request(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return parseTimeSeriesData(response.Body)
}

// DigitalCurrency queries statistics of a digital currency in terms of a physical currency throughout the day.
// Data is returned from past to present.
func (c *Client) DigitalCurrency(ctx context.Context, digital string, physical string) ([]*DigitalCurrencySeriesValue, error) {
	endpoint := c.buildRequestPath(map[string]string{
		queryEndpoint: valueDigitalCurrencyEndpoint,
		querySymbol:   digital,
		queryMarket:   physical,
	})
	response, err := c.Conn().Request(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return parseDigitalCurrencySeriesData(response.Body)
}
