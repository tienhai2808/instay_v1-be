package implement

import (
	"context"
	"math"
	"time"

	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type dashboardSvcImpl struct {
	userRepo    repository.UserRepository
	roomRepo    repository.RoomRepository
	serviceRepo repository.ServiceRepository
	bookingRepo repository.BookingRepository
	orderRepo   repository.OrderRepository
	requestRepo repository.RequestRepository
	reviewRepo  repository.ReviewRepository
	logger      *zap.Logger
}

func NewDashboardService(
	userRepo repository.UserRepository,
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
	bookingRepo repository.BookingRepository,
	orderRepo repository.OrderRepository,
	requestRepo repository.RequestRepository,
	reviewRepo repository.ReviewRepository,
	logger *zap.Logger,
) service.DashboardService {
	return &dashboardSvcImpl{
		userRepo,
		roomRepo,
		serviceRepo,
		bookingRepo,
		orderRepo,
		requestRepo,
		reviewRepo,
		logger,
	}
}

func (s *dashboardSvcImpl) Overview(ctx context.Context) (*types.DashboardResponse, error) {
	res := &types.DashboardResponse{
		OrderServiceStats:    make([]*types.StatusChartResponse, 0),
		RequestStats:         make([]*types.StatusChartResponse, 0),
		DailyBookingStats:    make([]*types.DailyBookingChartResponse, 0),
		BookingSourceStats:   make([]*types.ChartData, 0),
		ServiceUsageStats:    make([]*types.ChartData, 0),
		PopularRoomTypeStats: make([]*types.PopularRoomTypeChartData, 0),
		RevenueSourceStats:   make([]*types.ChartData, 0),
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		count, err := s.userRepo.Count(ctx)
		if err != nil {
			return err
		}
		res.TotalStaff = count
		return nil
	})

	g.Go(func() error {
		avg, err := s.reviewRepo.AverageRating(ctx)
		if err != nil {
			return err
		}
		res.AverageReviewRating = math.Round(avg*100) / 100
		return nil
	})

	g.Go(func() error {
		data, err := s.bookingRepo.GetBookingCountBySource(ctx)
		if err != nil {
			return err
		}
		res.BookingSourceStats = data
		return nil
	})

	g.Go(func() error {
		data, err := s.serviceRepo.GetServiceUsageStats(ctx)
		if err != nil {
			return err
		}
		res.ServiceUsageStats = data
		return nil
	})

	g.Go(func() error {
		total, err := s.roomRepo.CountRoom(ctx)
		if err != nil {
			return err
		}

		occupied, err := s.roomRepo.CountOccupancyRoom(ctx)
		if err != nil {
			return err
		}

		res.TotalRooms = total
		res.OccupiedRooms = occupied
		return nil
	})

	g.Go(func() error {
		data, err := s.orderRepo.GetPopularRoomTypeStats(ctx)
		if err != nil {
			return err
		}

		var total int64
		for _, item := range data {
			total += item.Count
		}

		if total > 0 {
			for _, item := range data {
				item.Percentage = math.Round((float64(item.Count)/float64(total))*100*100) / 100
			}
		}

		res.PopularRoomTypeStats = data
		return nil
	})

	g.Go(func() error {
		data, err := s.bookingRepo.GetRevenueBySource(ctx)
		if err != nil {
			return err
		}

		var totalRevenue float64
		for _, item := range data {
			totalRevenue += item.Value
		}
		if totalRevenue > 0 {
			for _, item := range data {
				item.Percentage = math.Round((item.Value/totalRevenue)*100*100) / 100
			}
		}
		res.RevenueSourceStats = data
		return nil
	})

	g.Go(func() error {
		count, err := s.serviceRepo.CountService(ctx)
		if err != nil {
			return err
		}
		res.TotalServices = count
		return nil
	})

	g.Go(func() error {
		count, err := s.bookingRepo.CountBooking(ctx)
		if err != nil {
			return err
		}
		res.TotalBookings = count
		return nil
	})

	g.Go(func() error {
		sum, err := s.bookingRepo.SumBookingTotalSellPrice(ctx)
		if err != nil {
			return err
		}
		res.BookingRevenue = sum
		return nil
	})

	g.Go(func() error {
		data, err := s.orderRepo.OrderServiceStatusDistribution(ctx)
		if err != nil {
			return err
		}

		calculatePercentage(data)
		res.OrderServiceStats = data
		return nil
	})

	g.Go(func() error {
		data, err := s.requestRepo.RequestStatusDistribution(ctx)
		if err != nil {
			return err
		}

		calculatePercentage(data)
		res.RequestStats = data
		return nil
	})

	g.Go(func() error {
		minDate, maxDate, err := s.bookingRepo.GetBookingDateRange(ctx)
		if err != nil {
			return err
		}

		rawStats, err := s.bookingRepo.GetDailyStats(ctx)
		if err != nil {
			return err
		}

		statsMap := make(map[string]*types.DailyBookingResult)
		for _, item := range rawStats {
			dateKey := item.Date
			if len(dateKey) > 10 {
				dateKey = dateKey[:10]
			}
			statsMap[dateKey] = item
		}

		var finalResponse []*types.DailyBookingChartResponse

		current := time.Date(minDate.Year(), minDate.Month(), minDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 0, 0, 0, 0, time.UTC)

		for !current.After(end) {
			dateStr := current.Format("2006-01-02")

			stat := &types.DailyBookingChartResponse{
				Date: dateStr,
			}
			if val, ok := statsMap[dateStr]; ok {
				stat.BookingCount = val.BookingCount
				stat.Revenue = val.Revenue
			} else {
				stat.BookingCount = 0
				stat.Revenue = 0
			}

			finalResponse = append(finalResponse, stat)

			current = current.AddDate(0, 0, 1)
		}

		res.DailyBookingStats = finalResponse
		return nil
	})

	if err := g.Wait(); err != nil {
		s.logger.Error("get dashboard failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func calculatePercentage(data []*types.StatusChartResponse) {
	var total int64
	for _, item := range data {
		total += item.Count
	}

	if total == 0 {
		return
	}

	for _, item := range data {
		item.Percentage = math.Round((float64(item.Count)/float64(total))*100*100) / 100
	}
}
