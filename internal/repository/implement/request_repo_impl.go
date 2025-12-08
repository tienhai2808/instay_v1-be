package implement

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type requestRepoImpl struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) repository.RequestRepository {
	return &requestRepoImpl{db}
}

func (r *requestRepoImpl) CreateRequestType(ctx context.Context, requestType *model.RequestType) error {
	return r.db.WithContext(ctx).Create(requestType).Error
}

func (r *requestRepoImpl) FindAllRequestTypesWithDetails(ctx context.Context) ([]*model.RequestType, error) {
	var requestTypes []*model.RequestType
	if err := r.db.WithContext(ctx).Preload("Department").Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&requestTypes).Error; err != nil {
		return nil, err
	}

	return requestTypes, nil
}

func (r *requestRepoImpl) RequestStatusDistribution(ctx context.Context) ([]*types.StatusChartResponse, error) {
	var results []*types.StatusChartResponse

	if err := r.db.WithContext(ctx).Model(&model.Request{}).Select("status, COUNT(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}


func (r *requestRepoImpl) FindAllRequestTypes(ctx context.Context) ([]*model.RequestType, error) {
	var requestTypes []*model.RequestType
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&requestTypes).Error; err != nil {
		return nil, err
	}

	return requestTypes, nil
}

func (r *requestRepoImpl) FindRequestTypeByID(ctx context.Context, requestTypeID int64) (*model.RequestType, error) {
	var requestType model.RequestType
	if err := r.db.WithContext(ctx).Where("id = ?", requestTypeID).First(&requestType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &requestType, nil
}

func (r *requestRepoImpl) FindRequestTypeByIDWithDetails(ctx context.Context, requestTypeID int64) (*model.RequestType, error) {
	var requestType model.RequestType
	if err := r.db.WithContext(ctx).Preload("Department.Staffs").Where("id = ?", requestTypeID).First(&requestType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &requestType, nil
}

func (r *requestRepoImpl) UpdateRequestType(ctx context.Context, requestTypeID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.RequestType{}).Where("id = ?", requestTypeID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRequestTypeNotFound
	}

	return nil
}

func (r *requestRepoImpl) DeleteRequestType(ctx context.Context, requestTypeID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", requestTypeID).Delete(&model.RequestType{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRequestTypeNotFound
	}

	return nil
}

func (r *requestRepoImpl) CreateRequest(ctx context.Context, request *model.Request) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *requestRepoImpl) FindRequestByIDWithRequestTypeDetailsAndOrderRoomDetailsTx(tx *gorm.DB, requestID int64) (*model.Request, error) {
	var request model.Request
	if err := tx.Preload("RequestType.Department.Staffs").Preload("OrderRoom.Booking").Where("id = ?", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}

func (r *requestRepoImpl) UpdateRequestTx(tx *gorm.DB, requestID int64, updateData map[string]any) error {
	return tx.Model(&model.Request{}).Where("id = ?", requestID).Updates(updateData).Error
}

func (r *requestRepoImpl) FindAllRequestsByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.Request, error) {
	var requests []*model.Request
	if err := r.db.WithContext(ctx).Preload("RequestType").Where("order_room_id = ?", orderRoomID).Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *requestRepoImpl) FindRequestByIDWithDetails(ctx context.Context, requestID int64) (*model.Request, error) {
	var request model.Request
	if err := r.db.WithContext(ctx).Preload("OrderRoom.Room.RoomType").Preload("OrderRoom.Room.Floor").Preload("UpdatedBy").Preload("RequestType.Department.Staffs").Where("id = ?", requestID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}

func (r *requestRepoImpl) FindAllRequestsWithDetailsPaginated(ctx context.Context, query types.RequestPaginationQuery, departmentID *int64) ([]*model.Request, int64, error) {
	var requests []*model.Request
	var total int64

	db := r.db.WithContext(ctx).Preload("OrderRoom.Room").Preload("RequestType").Model(&model.Request{})
	db = applyRequestFilters(db, query)

	if departmentID != nil {
		db = db.Joins("JOIN request_types rt ON rt.id = requests.request_type_id").
			Where("rt.department_id = ?", *departmentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRequestSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func applyRequestFilters(db *gorm.DB, query types.RequestPaginationQuery) *gorm.DB {
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.From != "" || query.To != "" {
		const layout = "2006-01-02"

		if query.From != "" {
			if parsedFrom, err := time.Parse(layout, query.From); err == nil {
				db = db.Where("request.created_at >= ?", parsedFrom)
			}
		}

		if query.To != "" {
			if parsedTo, err := time.Parse(layout, query.To); err == nil {
				endOfDay := parsedTo.AddDate(0, 0, 1)
				db = db.Where("request.created_at < ?", endOfDay)
			}
		}
	}

	return db
}

func applyRequestSorting(db *gorm.DB, query types.RequestPaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	db = db.Order("created_at " + strings.ToUpper(query.Order))

	return db
}
