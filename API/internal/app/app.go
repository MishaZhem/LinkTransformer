package app

import (
	grpcAnalyticsService "LinkTransformer/internal/ports/grpc/AnalyticsService"
	grpcLinkService "LinkTransformer/internal/ports/grpc/LinkService"
	"encoding/json"

	"context"
	"errors"
)

type Program struct {
	linkService      grpcLinkService.LinkServiceClient
	analyticsService grpcAnalyticsService.AnalyticsServiceClient
}

type App interface {
	GenerateLink(ctx context.Context, url string) (string, error)
	RedirectLink(ctx context.Context, url string) (string, error)
	GetStatistics(ctx context.Context, url string) (string, error)
	GetTotalClicks(ctx context.Context, url string) (int64, error)
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

func NewApp(linkService grpcLinkService.LinkServiceClient, analyticsService grpcAnalyticsService.AnalyticsServiceClient) App {
	return &Program{linkService: linkService, analyticsService: analyticsService}
}

func (r *Program) GenerateLink(ctx context.Context, url string) (string, error) {
	link, err := r.linkService.GenerateLink(ctx, &grpcLinkService.LinkRequest{
		Url: url,
	})
	if err != nil {
		return "", err
	}
	return link.Url, nil
}

func (r *Program) RedirectLink(ctx context.Context, key string) (string, error) {
	link, err := r.linkService.RedirectLink(ctx, &grpcLinkService.LinkRequest{
		Url: key,
	})
	if err != nil {
		return "", err
	}
	return link.Url, nil
}

func (r *Program) GetStatistics(ctx context.Context, url string) (string, error) {
	events, err := r.analyticsService.GetStatistics(ctx, &grpcAnalyticsService.LinkRequest{
		Url: url,
	})
	if err != nil {
		return "", err
	}
	result := make([]*ClickEvent, 0)
	for _, row := range events.List {
		result = append(result, statisticsResponseToRow(row))
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	return jsonString, nil
}

func (r *Program) GetTotalClicks(ctx context.Context, url string) (int64, error) {
	clicks, err := r.analyticsService.GetTotalClicks(ctx, &grpcAnalyticsService.LinkRequest{
		Url: url,
	})
	if err != nil {
		return 0, err
	}

	return clicks.TotalClicks, nil
}

type ClickEvent struct {
	LinkKey   string `json:"link_key"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Time      string `json:"time"`
}

func statisticsResponseToRow(row *grpcAnalyticsService.StatisticsResponse) *ClickEvent {
	return &ClickEvent{
		LinkKey:   row.Url,
		IP:        row.IP,
		UserAgent: row.UserAgent,
		Time:      row.Time,
	}
}
