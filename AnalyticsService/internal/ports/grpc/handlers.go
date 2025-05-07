package grpc

import (
	"AnalyticsService/internal/app"
	"AnalyticsService/internal/ports/kafka"
	context "context"
	"errors"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedAnalyticsServiceServer
	app app.App
}

func Register(gRPC *grpc.Server, app app.App) {
	RegisterAnalyticsServiceServer(gRPC, &Server{app: app})
}

func (s *Server) GetStatistics(ctx context.Context, req *LinkRequest) (*ListStatisticsResponse, error) {
	stats, err := s.app.GetStatistics(ctx, req.Url)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	result := make([]*StatisticsResponse, 0)
	for _, row := range stats {
		result = append(result, rowToStatisticsResponse(row))
	}
	return &ListStatisticsResponse{List: result}, nil
}

func (s *Server) GetTotalClicks(ctx context.Context, req *LinkRequest) (*GetTotalClicksResponse, error) {
	clicks, err := s.app.GetTotalClicks(ctx, req.Url)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	return &GetTotalClicksResponse{TotalClicks: clicks}, nil
}

func getStatusByError(err error) codes.Code {
	switch {
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, app.ErrBadRequest):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

func rowToStatisticsResponse(row kafka.ClickEvent) *StatisticsResponse {
	return &StatisticsResponse{
		Url:       row.LinkKey,
		IP:        row.IP,
		UserAgent: row.UserAgent,
		Time:      row.Time.String(),
	}
}
