package av

import (
	"context"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	if NewClient(WithAPIKey(testApiKey)) == nil {
		t.Error("got nil client")
	}
}

func TestClient_StockTimeSeries_buildsUrl(t *testing.T) {
	const (
		expected = "query?apikey=test&datatype=csv&function=TIME_SERIES_DAILY&outputsize=compact&symbol=TEST"
	)
	res := &http.Response{
		Body:       NewBuffCloser(sampleTimeSeriesData),
		StatusCode: http.StatusOK,
	}
	conn := NewResponseConnection(res)
	client := NewClient(WithAPIKey(testApiKey), WithConnection(conn))

	_, _ = client.StockTimeSeries(context.Background(), TimeSeriesDaily, "TEST")

	if conn.endpoint.String() != expected {
		t.Errorf("unexpected url, want %s got %s", expected, conn.endpoint.String())
	}
}

func TestClient_StockTimeSeries_getsResults(t *testing.T) {
	res := &http.Response{
		Body:       NewBuffCloser(sampleTimeSeriesData),
		StatusCode: http.StatusOK,
	}
	conn := NewResponseConnection(res)
	client := NewClient(WithAPIKey(testApiKey), WithConnection(conn))

	result, err := client.StockTimeSeries(context.Background(), TimeSeriesDaily, "TEST")
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
	if result == nil {
		t.Error("nil results")
	}
}
