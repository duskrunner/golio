package lol

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/internal"
	"github.com/KnutZuidema/golio/internal/mock"
)

func TestMatchClient_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *Matchlist
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &Matchlist{},
			doer: mock.NewJSONMockDoer(Matchlist{}, 200),
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
			want: &Matchlist{},
			doer: mock.NewRateLimitDoer(Matchlist{}),
		},
		{
			name: "unavailable once",
			want: &Matchlist{},
			doer: mock.NewUnavailableOnceDoer(Matchlist{}),
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
			got, err := (&MatchClient{c: client}).List("id", 0, 1)
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMatchClient_ListStream(t *testing.T) {
	t.Parallel()
	count := 0
	tests := []struct {
		name    string
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			doer: &mock.Doer{
				Custom: func(r *http.Request) (*http.Response, error) {
					if count == 0 {
						count++
						return mock.NewJSONMockDoer(Matchlist{
							Matches: make([]*MatchReference, 100),
						}, 200).Do(r)
					}
					return mock.NewJSONMockDoer(Matchlist{}, 200).Do(r)
				},
			},
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
			doer: mock.NewRateLimitDoer(Matchlist{}),
		},
		{
			name: "unavailable once",
			doer: mock.NewUnavailableOnceDoer(Matchlist{}),
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
			got := (&MatchClient{c: client}).ListStream("id")
			for res := range got {
				if res.Error != nil && tt.wantErr != nil {
					require.Equal(t, res.Error, tt.wantErr)
					break
				} else if res.Error != nil {
					require.Equal(t, res.Error, io.EOF)
					return
				}
			}
		})
	}
}

func TestMatchClient_Get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *Match
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &Match{},
			doer: mock.NewJSONMockDoer(Match{}, 200),
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
			want: &Match{},
			doer: mock.NewRateLimitDoer(Match{}),
		},
		{
			name: "unavailable once",
			want: &Match{},
			doer: mock.NewUnavailableOnceDoer(Match{}),
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
			got, err := (&MatchClient{c: client}).Get(1)
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMatchClient_GetTimeline(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *MatchTimeline
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &MatchTimeline{},
			doer: mock.NewJSONMockDoer(MatchTimeline{}, 200),
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
			want: &MatchTimeline{},
			doer: mock.NewRateLimitDoer(MatchTimeline{}),
		},
		{
			name: "unavailable once",
			want: &MatchTimeline{},
			doer: mock.NewUnavailableOnceDoer(MatchTimeline{}),
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
			got, err := (&MatchClient{c: client}).GetTimeline(0)
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMatchClient_ListIDsByTournamentCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    []int
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: []int{},
			doer: mock.NewJSONMockDoer([]int{}, 200),
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
			want: []int{},
			doer: mock.NewRateLimitDoer([]int{}),
		},
		{
			name: "unavailable once",
			want: []int{},
			doer: mock.NewUnavailableOnceDoer([]int{}),
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
			got, err := (&MatchClient{c: client}).ListIDsByTournamentCode("tournamentCode")
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMatchClient_GetForTournament(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    *Match
		doer    internal.Doer
		wantErr error
	}{
		{
			name: "get response",
			want: &Match{},
			doer: mock.NewJSONMockDoer(Match{}, 200),
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
			want: &Match{},
			doer: mock.NewRateLimitDoer(Match{}),
		},
		{
			name: "unavailable once",
			want: &Match{},
			doer: mock.NewUnavailableOnceDoer(Match{}),
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
			got, err := (&MatchClient{c: client}).GetForTournament(0, "tournamentCode")
			require.Equal(t, err, tt.wantErr, fmt.Sprintf("want err %v, got %v", tt.wantErr, err))
			if tt.wantErr == nil {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
