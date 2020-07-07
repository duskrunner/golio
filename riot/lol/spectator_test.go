package lol

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/internal"
	"github.com/KnutZuidema/golio/internal/mock"
)

func TestSpectatorClient_ListFeatured(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *FeaturedGames
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &FeaturedGames{},
			doer: mock.NewJSONMockDoer(FeaturedGames{}, 200),
		},
		{
			name: "unknown error status",
			wantErr: api.Error{
				Message:    "unknown error reason",
				StatusCode: 999,
			},
			doer: mock.NewStatusMockDoer(999),
		},
		{
			name:    "not found",
			wantErr: api.ErrNotFound,
			doer:    mock.NewStatusMockDoer(http.StatusNotFound),
		},
		{
			name: "rate limited",
			want: &FeaturedGames{},
			doer: mock.NewRateLimitDoer(FeaturedGames{}),
		},
		{
			name: "unavailable once",
			want: &FeaturedGames{},
			doer: mock.NewUnavailableOnceDoer(FeaturedGames{}),
		},
		{
			name:    "unavailable twice",
			wantErr: api.ErrServiceUnavailable,
			doer:    mock.NewStatusMockDoer(http.StatusServiceUnavailable),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := internal.NewClient(api.RegionEuropeWest, "API_KEY", tt.doer, logrus.StandardLogger())
			got, err := (&SpectatorClient{c: client}).ListFeatured()
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestSpectatorClient_GetCurrent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *GameInfo
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &GameInfo{},
			doer: mock.NewJSONMockDoer(GameInfo{}, 200),
		},
		{
			name: "unknown error status",
			wantErr: api.Error{
				Message:    "unknown error reason",
				StatusCode: 999,
			},
			doer: mock.NewStatusMockDoer(999),
		},
		{
			name:    "not found",
			wantErr: api.ErrNotFound,
			doer:    mock.NewStatusMockDoer(http.StatusNotFound),
		},
		{
			name: "rate limited",
			want: &GameInfo{},
			doer: mock.NewRateLimitDoer(GameInfo{}),
		},
		{
			name: "unavailable once",
			want: &GameInfo{},
			doer: mock.NewUnavailableOnceDoer(GameInfo{}),
		},
		{
			name:    "unavailable twice",
			wantErr: api.ErrServiceUnavailable,
			doer:    mock.NewStatusMockDoer(http.StatusServiceUnavailable),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := internal.NewClient(api.RegionEuropeWest, "API_KEY", tt.doer, logrus.StandardLogger())
			got, err := (&SpectatorClient{c: client}).GetCurrent("id")
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
