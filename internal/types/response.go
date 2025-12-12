package types

import "time"

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UploadPresignedURLResponse struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type ViewPresignedURLResponse struct {
	Url string `json:"url"`
}

type UserResponse struct {
	ID         int64                     `json:"id"`
	Username   string                    `json:"username"`
	Email      string                    `json:"email"`
	Phone      string                    `json:"phone"`
	Role       string                    `json:"role"`
	IsActive   bool                      `json:"is_active"`
	FirstName  string                    `json:"first_name"`
	LastName   string                    `json:"last_name"`
	CreatedAt  time.Time                 `json:"created_at"`
	Department *SimpleDepartmentResponse `json:"department"`
}

type DepartmentResponse struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	DisplayName string             `json:"display_name"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatedBy   *BasicUserResponse `json:"created_by"`
	UpdatedBy   *BasicUserResponse `json:"updated_by"`
	StaffCount  int64              `json:"staff_count"`
}

type SimpleUserResponse struct {
	ID         int64                     `json:"id"`
	FirstName  string                    `json:"first_name"`
	LastName   string                    `json:"last_name"`
	Role       string                    `json:"role"`
	IsActive   bool                      `json:"is_active"`
	CreatedAt  time.Time                 `json:"created_at"`
	Department *SimpleDepartmentResponse `json:"department"`
}

type SimpleBookingResponse struct {
	ID            int64     `json:"id"`
	BookingNumber string    `json:"booking_number"`
	GuestFullName string    `json:"guest_name"`
	BookedOn      time.Time `json:"booked_on"`
	CheckIn       time.Time `json:"check_in"`
	CheckOut      time.Time `json:"check_out"`
	Source        string    `json:"source"`
}

type BasicUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SimpleDepartmentResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type MetaResponse struct {
	Total      uint64 `json:"total"`
	Page       uint32 `json:"page"`
	Limit      uint32 `json:"limit"`
	TotalPages uint16 `json:"total_pages"`
	HasPrev    bool   `json:"has_prev"`
	HasNext    bool   `json:"has_next"`
}

type ServiceTypeResponse struct {
	ID           int64                     `json:"id"`
	Name         string                    `json:"name"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	CreatedBy    *BasicUserResponse        `json:"created_by"`
	UpdatedBy    *BasicUserResponse        `json:"updated_by"`
	Department   *SimpleDepartmentResponse `json:"department"`
	ServiceCount int64                     `json:"service_count"`
}

type SimpleServiceTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SimpleServiceTypeWithBasicServices struct {
	ID       int64                   `json:"id"`
	Name     string                  `json:"name"`
	Slug     string                  `json:"slug"`
	Services []*BasicServiceResponse `json:"services"`
}

type SimpleServiceImageResponse struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type ServiceImageResponse struct {
	ID          int64  `json:"id"`
	Key         string `json:"key"`
	IsThumbnail bool   `json:"is_thumbnail"`
	SortOrder   uint32 `json:"sort_order"`
}

type BasicServiceResponse struct {
	ID          int64                       `json:"id"`
	Name        string                      `json:"name"`
	Slug        string                      `json:"slug"`
	Price       float64                     `json:"price"`
	IsActive    bool                        `json:"is_active"`
	ServiceType *SimpleServiceTypeResponse  `json:"service_type"`
	Thumbnail   *SimpleServiceImageResponse `json:"thumbnail"`
}

type SimpleServiceResponse struct {
	ID            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Price         float64                    `json:"price"`
	Description   string                     `json:"description"`
	ServiceType   *SimpleServiceTypeResponse `json:"service_type"`
	ServiceImages []*ServiceImageResponse    `json:"images"`
}

type ServiceResponse struct {
	ID            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Price         float64                    `json:"price"`
	IsActive      bool                       `json:"is_active"`
	Description   string                     `json:"description"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	ServiceType   *SimpleServiceTypeResponse `json:"service_type"`
	CreatedBy     *BasicUserResponse         `json:"created_by"`
	UpdatedBy     *BasicUserResponse         `json:"updated_by"`
	ServiceImages []*ServiceImageResponse    `json:"images"`
}

type RequestTypeResponse struct {
	ID         int64                     `json:"id"`
	Name       string                    `json:"name"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
	CreatedBy  *BasicUserResponse        `json:"created_by"`
	UpdatedBy  *BasicUserResponse        `json:"updated_by"`
	Department *SimpleDepartmentResponse `json:"department"`
}

type RoomTypeResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedBy *BasicUserResponse `json:"created_by"`
	UpdatedBy *BasicUserResponse `json:"updated_by"`
	RoomCount int64              `json:"room_count"`
}

type OrderRoomResponse struct {
	ID        int64                  `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedBy *BasicUserResponse     `json:"created_by"`
	UpdatedBy *BasicUserResponse     `json:"updated_by"`
	Room      *SimpleRoomResponse    `json:"room"`
	Booking   *SimpleBookingResponse `json:"booking"`
}

type BasicBookingResponse struct {
	ID            int64     `json:"id"`
	BookingNumber string    `json:"booking_number"`
	CheckIn       time.Time `json:"check_in"`
	CheckOut      time.Time `json:"check_out"`
}

type BookingResponse struct {
	ID                 int64                     `json:"id"`
	BookingNumber      string                    `json:"booking_number"`
	GuestFullName      string                    `json:"guest_full_name"`
	GuestEmail         string                    `json:"guest_email"`
	GuestPhone         string                    `json:"guest_phone"`
	CheckIn            time.Time                 `json:"check_in"`
	CheckOut           time.Time                 `json:"check_out"`
	RoomType           string                    `json:"room_type"`
	RoomNumber         uint32                    `json:"room_number"`
	GuestNumber        string                    `json:"guest_number"`
	BookedOn           time.Time                 `json:"booked_on"`
	Source             string                    `json:"source"`
	TotalNetPrice      float64                   `json:"total_net_price"`
	TotalSellPrice     float64                   `json:"total_sell_price"`
	PromotionName      string                    `json:"promotion_name"`
	MealPlan           string                    `json:"meal_plan"`
	BookingPreferences string                    `json:"booking_references"`
	BookingConditions  string                    `json:"booking_conditions"`
	OrderRooms         []*BasicOrderRoomResponse `json:"order_rooms"`
}

type SourceResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FloorResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SimpleRoomTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type BasicRoomWithBasicOrderRoomsResponse struct {
	ID         int64                                     `json:"id"`
	Name       string                                    `json:"name"`
	OrderRooms []*BasicOrderRoomWithBasicBookingResponse `json:"order_rooms"`
}

type SimpleRoomResponse struct {
	ID       int64                   `json:"id"`
	Name     string                  `json:"name"`
	RoomType *SimpleRoomTypeResponse `json:"room_type"`
	Floor    string                  `json:"floor"`
}

type RoomResponse struct {
	ID        int64                   `json:"id"`
	Name      string                  `json:"name"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	CreatedBy *BasicUserResponse      `json:"created_by"`
	UpdatedBy *BasicUserResponse      `json:"updated_by"`
	RoomType  *SimpleRoomTypeResponse `json:"room_type"`
	Floor     string                  `json:"floor"`
	InUse     bool                    `json:"in_use"`
}

type SimpleOrderServiceResponse struct {
	ID           int64                 `json:"id"`
	Service      *BasicServiceResponse `json:"service"`
	Quantity     uint32                `json:"quantity"`
	TotalPrice   float64               `json:"total_price"`
	Status       string                `json:"status"`
	CreatedAt    time.Time             `json:"created_at"`
	GuestNote    *string               `json:"guest_note"`
	StaffNote    *string               `json:"staff_note"`
	CancelReason *string               `json:"cancel_reason"`
	RejectReason *string               `json:"reject_reason"`
}

type BasicOrderServiceResponse struct {
	ID         int64     `json:"id"`
	Service    string    `json:"service"`
	Room       string    `json:"room"`
	Quantity   uint32    `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderServiceResponse struct {
	ID           int64                   `json:"id"`
	Service      *BasicServiceResponse   `json:"service"`
	OrderRoom    *BasicOrderRoomResponse `json:"order_room"`
	Quantity     uint32                  `json:"quantity"`
	TotalPrice   float64                 `json:"total_price"`
	Status       string                  `json:"status"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	GuestNote    *string                 `json:"guest_note"`
	StaffNote    *string                 `json:"staff_note"`
	CancelReason *string                 `json:"cancel_reason"`
	RejectReason *string                 `json:"reject_reason"`
	UpdatedBy    *BasicUserResponse      `json:"updated_by"`
}

type BasicOrderRoomResponse struct {
	ID   int64               `json:"id"`
	Room *SimpleRoomResponse `json:"room"`
}

type BasicOrderRoomWithBasicBookingResponse struct {
	ID      int64                 `json:"id"`
	Booking *BasicBookingResponse `json:"booking"`
}

type SimpleOrderRoomResponse struct {
	ID      int64                  `json:"id"`
	Room    *SimpleRoomResponse    `json:"room"`
	Booking *SimpleBookingResponse `json:"booking"`
}

type NotificationStaffResponse struct {
	ID     int64     `json:"id"`
	ReadAt time.Time `json:"read_at"`
}

type SimpleNotificationResponse struct {
	ID        int64                      `json:"id"`
	Type      string                     `json:"type"`
	Content   string                     `json:"content"`
	ContentID int64                      `json:"content_id"`
	Receiver  string                     `json:"receiver"`
	IsRead    bool                       `json:"is_read"`
	ReadAt    *time.Time                 `json:"read_at"`
	CreatedAt time.Time                  `json:"created_at"`
	StaffRead *NotificationStaffResponse `json:"staff_read"`
}

type SimpleRequestTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SimpleRequestResponse struct {
	ID          int64                      `json:"id"`
	RequestType *SimpleRequestTypeResponse `json:"request_type"`
	Content     string                     `json:"content"`
	Status      string                     `json:"status"`
	CreatedAt   time.Time                  `json:"created_at"`
}

type RequestResponse struct {
	ID          int64                      `json:"id"`
	RequestType *SimpleRequestTypeResponse `json:"request_type"`
	OrderRoom   *BasicOrderRoomResponse    `json:"order_room"`
	Content     string                     `json:"content"`
	Status      string                     `json:"status"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
	UpdatedBy   *BasicUserResponse         `json:"updated_by"`
}

type BasicRequestResponse struct {
	ID          int64     `json:"id"`
	RequestType string    `json:"request_type"`
	Room        string    `json:"room"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type SimpleMessageResponse struct {
	ID         int64                   `json:"id"`
	Content    *string                 `json:"content"`
	ImageKey   *string                 `json:"image_key"`
	SenderType string                  `json:"sender_type"`
	Sender     *BasicUserResponse      `json:"sender_id"`
	CreatedAt  time.Time               `json:"created_at"`
	IsRead     bool                    `json:"is_read"`
	ReadAt     *time.Time              `json:"read_at"`
	StaffReads []*MessageStaffResponse `json:"staff_reads"`
}

type MessageResponse struct {
	ID         int64                   `json:"id"`
	Content    *string                 `json:"content"`
	ImageKey   *string                 `json:"image_key"`
	SenderType string                  `json:"sender_type"`
	Sender     *BasicUserResponse      `json:"sender_id"`
	CreatedAt  time.Time               `json:"created_at"`
	IsRead     bool                    `json:"is_read"`
	ReadAt     *time.Time              `json:"read_at"`
	StaffReads []*MessageStaffResponse `json:"staff_reads"`
	ChatID     int64                   `json:"chat_id"`
}

type MessageStaffResponse struct {
	ID     int64              `json:"id"`
	ReadAt time.Time          `json:"read_at"`
	Staff  *BasicUserResponse `json:"staff"`
}

type SimpleChatResponse struct {
	ID          int64                    `json:"id"`
	OrderRoom   *SimpleOrderRoomResponse `json:"order_room"`
	ExpiredAt   time.Time                `json:"expired_at"`
	LastMessage *SimpleMessageResponse   `json:"last_message"`
}

type SimpleChatWithMessageResponse struct {
	ID        int64                    `json:"id"`
	OrderRoom *SimpleOrderRoomResponse `json:"order_room"`
	ExpiredAt time.Time                `json:"expired_at"`
	Messages  []*SimpleMessageResponse `json:"messages"`
}

type BasicChatResponse struct {
	ID          int64                 `json:"id"`
	Code        string                `json:"code"`
	ExpiredAt   time.Time             `json:"expired_at"`
	LastMessage *BasicMessageResponse `json:"last_message"`
}

type BasicChatWithMessageResponse struct {
	ID         int64                     `json:"id"`
	Department *SimpleDepartmentResponse `json:"department"`
	ExpiredAt  time.Time                 `json:"expired_at"`
	Messages   []*BasicMessageResponse   `json:"messages"`
}

type BasicMessageResponse struct {
	ID         int64      `json:"id"`
	Content    *string    `json:"content"`
	ImageKey   *string    `json:"image_key"`
	SenderType string     `json:"sender_type"`
	CreatedAt  time.Time  `json:"created_at"`
	IsRead     bool       `json:"is_read"`
	ReadAt     *time.Time `json:"read_at"`
}

type BasicNotificationResponse struct {
	ID        int64      `json:"id"`
	Type      string     `json:"type"`
	Content   string     `json:"content"`
	ContentID int64      `json:"content_id"`
	Receiver  string     `json:"receiver"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type WSResponse struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type UpdateReadMessagesResponse struct {
	ChatID     int64              `json:"chat_id"`
	ReaderType string             `json:"reader_type"`
	ReadAt     time.Time          `json:"read_at"`
	Reader     *BasicUserResponse `json:"reader"`
}

type SimpleReviewResponse struct {
	ID        int64     `json:"id"`
	Star      uint32    `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewResponse struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	Star        uint32    `json:"star"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderRoomID int64     `json:"order_room_id"`
}

type DashboardResponse struct {
	TotalStaff     int64   `json:"total_staff"`
	TotalRooms     int64   `json:"total_rooms"`
	OccupiedRooms  int64   `json:"occupied_rooms"`
	TotalServices  int64   `json:"total_services"`
	TotalBookings  int64   `json:"total_bookings"`
	BookingRevenue float64 `json:"booking_revenue"`

	AverageReviewRating float64 `json:"average_review_rating"`

	BookingSourceStats   []*ChartData                `json:"booking_source_stats"`
	ServiceUsageStats    []*ChartData                `json:"service_usage_stats"`
	PopularRoomTypeStats []*PopularRoomTypeChartData `json:"popular_room_type_stats"`
	RevenueSourceStats   []*ChartData                `json:"revenue_source_stats"`

	OrderServiceStats []*StatusChartResponse       `json:"order_service_stats"`
	RequestStats      []*StatusChartResponse       `json:"request_stats"`
	DailyBookingStats []*DailyBookingChartResponse `json:"daily_booking_stats"`
}

type StatusChartResponse struct {
	Status     string  `json:"status"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type DailyBookingChartResponse struct {
	Date         string  `json:"date"`
	BookingCount int64   `json:"booking_count"`
	Revenue      float64 `json:"revenue"`
}
