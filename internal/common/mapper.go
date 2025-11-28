package common

import (
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

func ToUserResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToSimpleDepartmentResponse(user.Department),
	}
}

func ToDepartmentData(department *model.Department) *types.DepartmentData {
	if department == nil {
		return nil
	}

	return &types.DepartmentData{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
	}
}

func ToSimpleDepartmentResponse(department *model.Department) *types.SimpleDepartmentResponse {
	if department == nil {
		return nil
	}

	return &types.SimpleDepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
	}
}

func ToRoomResponse(room *model.Room) *types.RoomResponse {
	if room == nil {
		return nil
	}

	return &types.RoomResponse{
		ID:        room.ID,
		Name:      room.Name,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
		CreatedBy: ToBasicUserResponse(room.CreatedBy),
		UpdatedBy: ToBasicUserResponse(room.UpdatedBy),
		RoomType:  ToSimpleRoomTypeResponse(room.RoomType),
		Floor:     ToFloorResponse(room.Floor),
	}
}

func ToRoomsResponse(rooms []*model.Room) []*types.RoomResponse {
	if len(rooms) == 0 {
		return make([]*types.RoomResponse, 0)
	}

	roomsRes := make([]*types.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomsRes = append(roomsRes, ToRoomResponse(room))
	}

	return roomsRes
}

func ToUserData(user *model.User) *types.UserData {
	return &types.UserData{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToDepartmentData(user.Department),
	}
}

func ToSimpleUserResponse(user *model.User) *types.SimpleUserResponse {
	return &types.SimpleUserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToSimpleDepartmentResponse(user.Department),
	}
}

func ToSimpleUsersResponse(users []*model.User) []*types.SimpleUserResponse {
	if len(users) == 0 {
		return make([]*types.SimpleUserResponse, 0)
	}

	usersRes := make([]*types.SimpleUserResponse, 0, len(users))
	for _, user := range users {
		usersRes = append(usersRes, ToSimpleUserResponse(user))
	}

	return usersRes
}

func ToDepartmentResponse(department *model.Department) *types.DepartmentResponse {
	return &types.DepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
		Description: department.Description,
		CreatedAt:   department.CreatedAt,
		UpdatedAt:   department.UpdatedAt,
		CreatedBy:   ToBasicUserResponse(department.CreatedBy),
		UpdatedBy:   ToBasicUserResponse(department.UpdatedBy),
		StaffCount:  department.StaffCount,
	}
}

func ToDepartmentsResponse(departments []*model.Department) []*types.DepartmentResponse {
	if len(departments) == 0 {
		return make([]*types.DepartmentResponse, 0)
	}

	departmentsRes := make([]*types.DepartmentResponse, 0, len(departments))
	for _, department := range departments {
		departmentsRes = append(departmentsRes, ToDepartmentResponse(department))
	}

	return departmentsRes
}

func ToBasicUserResponse(user *model.User) *types.BasicUserResponse {
	if user == nil {
		return nil
	}

	return &types.BasicUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func ToServiceTypeResponse(serviceType *model.ServiceType) *types.ServiceTypeResponse {
	if serviceType == nil {
		return nil
	}

	return &types.ServiceTypeResponse{
		ID:           serviceType.ID,
		Name:         serviceType.Name,
		CreatedAt:    serviceType.CreatedAt,
		UpdatedAt:    serviceType.UpdatedAt,
		CreatedBy:    ToBasicUserResponse(serviceType.CreatedBy),
		UpdatedBy:    ToBasicUserResponse(serviceType.UpdatedBy),
		Department:   ToSimpleDepartmentResponse(serviceType.Department),
		ServiceCount: serviceType.ServiceCount,
	}
}

func ToSimpleDepartmentsResponse(departments []*model.Department) []*types.SimpleDepartmentResponse {
	if len(departments) == 0 {
		return make([]*types.SimpleDepartmentResponse, 0)
	}

	departmentsRes := make([]*types.SimpleDepartmentResponse, 0, len(departments))
	for _, department := range departments {
		departmentsRes = append(departmentsRes, ToSimpleDepartmentResponse(department))
	}

	return departmentsRes
}

func ToServiceTypesResponse(serviceTypes []*model.ServiceType) []*types.ServiceTypeResponse {
	if len(serviceTypes) == 0 {
		return make([]*types.ServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*types.ServiceTypeResponse, 0, len(serviceTypes))
	for _, serviceType := range serviceTypes {
		serviceTypesRes = append(serviceTypesRes, ToServiceTypeResponse(serviceType))
	}

	return serviceTypesRes
}

func ToSimpleServiceTypeResponse(serviceType *model.ServiceType) *types.SimpleServiceTypeResponse {
	if serviceType == nil {
		return nil
	}

	return &types.SimpleServiceTypeResponse{
		ID:   serviceType.ID,
		Name: serviceType.Name,
		Slug: serviceType.Slug,
	}
}

func ToSimpleServiceImageResponse(image *model.ServiceImage) *types.SimpleServiceImageResponse {
	if image == nil {
		return nil
	}

	return &types.SimpleServiceImageResponse{
		ID:  image.ID,
		Key: image.Key,
	}
}

func ToBasicServiceResponse(service *model.Service) *types.BasicServiceResponse {
	if service == nil {
		return nil
	}

	return &types.BasicServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Slug:        service.Slug,
		Price:       service.Price,
		IsActive:    service.IsActive,
		ServiceType: ToSimpleServiceTypeResponse(service.ServiceType),
		Thumbnail:   ToSimpleServiceImageResponse(service.ServiceImages[0]),
	}
}

func ToSimpleServiceResponse(service *model.Service) *types.SimpleServiceResponse {
	if service == nil {
		return nil
	}

	return &types.SimpleServiceResponse{
		ID:            service.ID,
		Name:          service.Name,
		Price:         service.Price,
		Description:   service.Description,
		ServiceType:   ToSimpleServiceTypeResponse(service.ServiceType),
		ServiceImages: ToServiceImagesResponse(service.ServiceImages),
	}
}

func ToBasicServicesResponse(services []*model.Service) []*types.BasicServiceResponse {
	if len(services) == 0 {
		return make([]*types.BasicServiceResponse, 0)
	}

	servicesRes := make([]*types.BasicServiceResponse, 0, len(services))
	for _, service := range services {
		servicesRes = append(servicesRes, ToBasicServiceResponse(service))
	}

	return servicesRes
}

func ToSimpleServiceTypeWithBasicServices(serviceType *model.ServiceType) *types.SimpleServiceTypeWithBasicServices {
	if serviceType == nil {
		return nil
	}

	return &types.SimpleServiceTypeWithBasicServices{
		ID:       serviceType.ID,
		Name:     serviceType.Name,
		Slug:     serviceType.Slug,
		Services: ToBasicServicesResponse(serviceType.Services),
	}
}

func ToServiceImageResponse(image *model.ServiceImage) *types.ServiceImageResponse {
	if image == nil {
		return nil
	}

	return &types.ServiceImageResponse{
		ID:          image.ID,
		Key:         image.Key,
		IsThumbnail: image.IsThumbnail,
		SortOrder:   image.SortOrder,
	}
}

func ToServiceImagesResponse(images []*model.ServiceImage) []*types.ServiceImageResponse {
	if len(images) == 0 {
		return make([]*types.ServiceImageResponse, 0)
	}

	imageRes := make([]*types.ServiceImageResponse, 0, len(images))
	for _, img := range images {
		imageRes = append(imageRes, ToServiceImageResponse(img))
	}

	return imageRes
}

func ToServiceResponse(service *model.Service) *types.ServiceResponse {
	if service == nil {
		return nil
	}

	return &types.ServiceResponse{
		ID:            service.ID,
		Name:          service.Name,
		Price:         service.Price,
		IsActive:      service.IsActive,
		Description:   service.Description,
		CreatedAt:     service.CreatedAt,
		UpdatedAt:     service.UpdatedAt,
		ServiceType:   ToSimpleServiceTypeResponse(service.ServiceType),
		CreatedBy:     ToBasicUserResponse(service.CreatedBy),
		UpdatedBy:     ToBasicUserResponse(service.UpdatedBy),
		ServiceImages: ToServiceImagesResponse(service.ServiceImages),
	}
}

func ToSimpleServiceTypesResponse(serviceTypes []*model.ServiceType) []*types.SimpleServiceTypeResponse {
	if len(serviceTypes) == 0 {
		return make([]*types.SimpleServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*types.SimpleServiceTypeResponse, 0, len(serviceTypes))
	for _, serviceType := range serviceTypes {
		serviceTypesRes = append(serviceTypesRes, ToSimpleServiceTypeResponse(serviceType))
	}

	return serviceTypesRes
}

func ToRequestTypeResponse(requestType *model.RequestType) *types.RequestTypeResponse {
	if requestType == nil {
		return nil
	}

	return &types.RequestTypeResponse{
		ID:         requestType.ID,
		Name:       requestType.Name,
		CreatedAt:  requestType.CreatedAt,
		UpdatedAt:  requestType.UpdatedAt,
		CreatedBy:  ToBasicUserResponse(requestType.CreatedBy),
		UpdatedBy:  ToBasicUserResponse(requestType.UpdatedBy),
		Department: ToSimpleDepartmentResponse(requestType.Department),
	}
}

func ToRequestTypesResponse(requestTypes []*model.RequestType) []*types.RequestTypeResponse {
	if len(requestTypes) == 0 {
		return make([]*types.RequestTypeResponse, 0)
	}

	requestTypesRes := make([]*types.RequestTypeResponse, 0, len(requestTypes))
	for _, requestType := range requestTypes {
		requestTypesRes = append(requestTypesRes, ToRequestTypeResponse(requestType))
	}

	return requestTypesRes
}

func ToSimpleRequestTypesResponse(requestTypes []*model.RequestType) []*types.SimpleRequestTypeResponse {
	if len(requestTypes) == 0 {
		return make([]*types.SimpleRequestTypeResponse, 0)
	}

	requestTypesRes := make([]*types.SimpleRequestTypeResponse, 0, len(requestTypes))
	for _, requestType := range requestTypes {
		requestTypesRes = append(requestTypesRes, ToSimpleRequestTypeResponse(requestType))
	}

	return requestTypesRes
}

func ToRoomTypeResponse(roomType *model.RoomType) *types.RoomTypeResponse {
	if roomType == nil {
		return nil
	}

	return &types.RoomTypeResponse{
		ID:        roomType.ID,
		Name:      roomType.Name,
		CreatedAt: roomType.CreatedAt,
		UpdatedAt: roomType.UpdatedAt,
		CreatedBy: ToBasicUserResponse(roomType.CreatedBy),
		UpdatedBy: ToBasicUserResponse(roomType.UpdatedBy),
		RoomCount: roomType.RoomCount,
	}
}

func ToRoomTypesResponse(roomTypes []*model.RoomType) []*types.RoomTypeResponse {
	if len(roomTypes) == 0 {
		return make([]*types.RoomTypeResponse, 0)
	}

	roomTypesRes := make([]*types.RoomTypeResponse, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		roomTypesRes = append(roomTypesRes, ToRoomTypeResponse(roomType))
	}

	return roomTypesRes
}

func ToSimpleRoomTypeResponse(roomType *model.RoomType) *types.SimpleRoomTypeResponse {
	if roomType == nil {
		return nil
	}

	return &types.SimpleRoomTypeResponse{
		ID:   roomType.ID,
		Name: roomType.Name,
	}
}

func ToSimpleRoomTypesResponse(roomTypes []*model.RoomType) []*types.SimpleRoomTypeResponse {
	if len(roomTypes) == 0 {
		return make([]*types.SimpleRoomTypeResponse, 0)
	}

	roomTypesRes := make([]*types.SimpleRoomTypeResponse, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		roomTypesRes = append(roomTypesRes, ToSimpleRoomTypeResponse(roomType))
	}

	return roomTypesRes
}

func ToSimpleBookingResponse(booking *model.Booking) *types.SimpleBookingResponse {
	if booking == nil {
		return nil
	}

	return &types.SimpleBookingResponse{
		ID:            booking.ID,
		BookingNumber: booking.BookingNumber,
		GuestFullName: booking.GuestFullName,
		BookedOn:      booking.BookedOn,
		CheckIn:       booking.CheckIn,
		CheckOut:      booking.CheckOut,
	}
}

func ToSimpleBookingsResponse(bookings []*model.Booking) []*types.SimpleBookingResponse {
	if len(bookings) == 0 {
		return make([]*types.SimpleBookingResponse, 0)
	}

	bookingsRes := make([]*types.SimpleBookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		bookingsRes = append(bookingsRes, ToSimpleBookingResponse(booking))
	}

	return bookingsRes
}

func ToBookingResponse(booking *model.Booking) *types.BookingResponse {
	if booking == nil {
		return nil
	}

	return &types.BookingResponse{
		ID:                 booking.ID,
		BookingNumber:      booking.BookingNumber,
		GuestFullName:      booking.GuestFullName,
		GuestEmail:         booking.GuestEmail,
		GuestPhone:         booking.GuestPhone,
		CheckIn:            booking.CheckIn,
		CheckOut:           booking.CheckOut,
		RoomType:           booking.RoomType,
		RoomNumber:         booking.RoomNumber,
		GuestNumber:        booking.GuestNumber,
		BookedOn:           booking.BookedOn,
		Source:             booking.Source,
		TotalNetPrice:      booking.TotalNetPrice,
		TotalSellPrice:     booking.TotalSellPrice,
		PromotionName:      booking.PromotionName,
		MealPlan:           booking.MealPlan,
		BookingPreferences: booking.BookingPreferences,
		BookingConditions:  booking.BookingConditions,
	}
}

func ToFloorResponse(floor *model.Floor) *types.FloorResponse {
	if floor == nil {
		return nil
	}

	return &types.FloorResponse{
		ID:   floor.ID,
		Name: floor.Name,
	}
}

func ToSimpleRoomResponse(room *model.Room) *types.SimpleRoomResponse {
	if room == nil {
		return nil
	}

	return &types.SimpleRoomResponse{
		ID:       room.ID,
		Name:     room.Name,
		RoomType: ToSimpleRoomTypeResponse(room.RoomType),
		Floor:    ToFloorResponse(room.Floor),
	}
}

func ToOrderRoomResponse(orderRoom *model.OrderRoom) *types.OrderRoomResponse {
	if orderRoom == nil {
		return nil
	}

	return &types.OrderRoomResponse{
		ID:        orderRoom.ID,
		CreatedAt: orderRoom.CreatedAt,
		UpdatedAt: orderRoom.UpdatedAt,
		CreatedBy: ToBasicUserResponse(orderRoom.CreatedBy),
		UpdatedBy: ToBasicUserResponse(orderRoom.UpdatedBy),
		Room:      ToSimpleRoomResponse(orderRoom.Room),
		Booking:   ToSimpleBookingResponse(orderRoom.Booking),
	}
}

func ToFloorsResponse(floors []*model.Floor) []*types.FloorResponse {
	if len(floors) == 0 {
		return make([]*types.FloorResponse, 0)
	}

	floorsRes := make([]*types.FloorResponse, 0, len(floors))
	for _, floor := range floors {
		floorsRes = append(floorsRes, ToFloorResponse(floor))
	}

	return floorsRes
}

func ToSimpleOrderServiceResponse(orderService *model.OrderService) *types.SimpleOrderServiceResponse {
	if orderService == nil {
		return nil
	}

	return &types.SimpleOrderServiceResponse{
		ID:           orderService.ID,
		Code:         orderService.Code,
		Service:      ToBasicServiceResponse(orderService.Service),
		Quantity:     orderService.Quantity,
		TotalPrice:   orderService.TotalPrice,
		Status:       orderService.Status,
		GuestNote:    orderService.GuestNote,
		StaffNote:    orderService.StaffNote,
		CancelReason: orderService.CancelReason,
		CreatedAt:    orderService.CreatedAt,
	}
}

func ToBasicOrderServiceResponse(orderService *model.OrderService) *types.BasicOrderServiceResponse {
	if orderService == nil {
		return nil
	}

	return &types.BasicOrderServiceResponse{
		ID:          orderService.ID,
		Code:        orderService.Code,
		ServiceName: orderService.Service.Name,
		RoomName:    orderService.OrderRoom.Room.Name,
		Quantity:    orderService.Quantity,
		TotalPrice:  orderService.TotalPrice,
		Status:      orderService.Status,
		CreatedAt:   orderService.CreatedAt,
	}
}

func ToBasicOrderServicesResponse(orderServices []*model.OrderService) []*types.BasicOrderServiceResponse {
	if len(orderServices) == 0 {
		return make([]*types.BasicOrderServiceResponse, 0)
	}

	orderServicesRes := make([]*types.BasicOrderServiceResponse, 0, len(orderServices))
	for _, orderService := range orderServices {
		orderServicesRes = append(orderServicesRes, ToBasicOrderServiceResponse(orderService))
	}

	return orderServicesRes
}

func ToSimpleOrderRoomResponse(orderRoom *model.OrderRoom) *types.SimpleOrderRoomResponse {
	if orderRoom == nil {
		return nil
	}

	return &types.SimpleOrderRoomResponse{
		ID:   orderRoom.ID,
		Room: ToSimpleRoomResponse(orderRoom.Room),
	}
}

func ToOrderServiceResponse(orderService *model.OrderService) *types.OrderServiceResponse {
	if orderService == nil {
		return nil
	}

	return &types.OrderServiceResponse{
		ID:           orderService.ID,
		Code:         orderService.Code,
		Service:      ToBasicServiceResponse(orderService.Service),
		OrderRoom:    ToSimpleOrderRoomResponse(orderService.OrderRoom),
		Quantity:     orderService.Quantity,
		TotalPrice:   orderService.TotalPrice,
		Status:       orderService.Status,
		CreatedAt:    orderService.CreatedAt,
		UpdatedAt:    orderService.UpdatedAt,
		GuestNote:    orderService.GuestNote,
		StaffNote:    orderService.StaffNote,
		CancelReason: orderService.CancelReason,
		UpdatedBy:    ToBasicUserResponse(orderService.UpdatedBy),
	}
}

func ToSimpleNotificationResponse(notification *model.Notification) *types.SimpleNotificationResponse {
	if notification == nil {
		return nil
	}

	return &types.SimpleNotificationResponse{
		ID:        notification.ID,
		Type:      notification.Type,
		Content:   notification.Content,
		ContentID: notification.ContentID,
		Receiver:  notification.Receiver,
		IsRead:    notification.IsRead,
		ReadAt:    notification.ReadAt,
		CreatedAt: notification.CreatedAt,
		StaffRead: ToNotificationStaffResponse(notification.StaffsRead[0]),
	}
}

func ToSimpleNotificationsResponse(notifications []*model.Notification) []*types.SimpleNotificationResponse {
	if len(notifications) == 0 {
		return make([]*types.SimpleNotificationResponse, 0)
	}

	notificationsRes := make([]*types.SimpleNotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		notificationsRes = append(notificationsRes, ToSimpleNotificationResponse(notification))
	}

	return notificationsRes
}

func ToSimpleRequestTypeResponse(requestType *model.RequestType) *types.SimpleRequestTypeResponse {
	if requestType == nil {
		return nil
	}

	return &types.SimpleRequestTypeResponse{
		ID:   requestType.ID,
		Name: requestType.Name,
		Slug: requestType.Slug,
	}
}

func ToSimpleRequestResponse(request *model.Request) *types.SimpleRequestResponse {
	if request == nil {
		return nil
	}

	return &types.SimpleRequestResponse{
		ID:          request.ID,
		Code:        request.Code,
		Content:     request.Content,
		RequestType: ToSimpleRequestTypeResponse(request.RequestType),
		Status:      request.Status,
		CreatedAt:   request.CreatedAt,
	}
}

func ToSimpleRequestsResponse(requests []*model.Request) []*types.SimpleRequestResponse {
	if len(requests) == 0 {
		return make([]*types.SimpleRequestResponse, 0)
	}

	requestsRes := make([]*types.SimpleRequestResponse, 0, len(requests))
	for _, request := range requests {
		requestsRes = append(requestsRes, ToSimpleRequestResponse(request))
	}

	return requestsRes
}

func ToSimpleOrderServicesResponse(orderServices []*model.OrderService) []*types.SimpleOrderServiceResponse {
	if len(orderServices) == 0 {
		return make([]*types.SimpleOrderServiceResponse, 0)
	}

	orderServicesRes := make([]*types.SimpleOrderServiceResponse, 0, len(orderServices))
	for _, orderService := range orderServices {
		orderServicesRes = append(orderServicesRes, ToSimpleOrderServiceResponse(orderService))
	}

	return orderServicesRes
}

func ToNotificationStaffResponse(notificationStaff *model.NotificationStaff) *types.NotificationStaffResponse {
	if notificationStaff == nil {
		return nil
	}

	return &types.NotificationStaffResponse{
		ID:     notificationStaff.ID,
		ReadAt: notificationStaff.ReadAt,
	}
}

func ToRequestResponse(request *model.Request) *types.RequestResponse {
	if request == nil {
		return nil
	}

	return &types.RequestResponse{
		ID:          request.ID,
		Code:        request.Code,
		RequestType: ToSimpleRequestTypeResponse(request.RequestType),
		OrderRoom:   ToSimpleOrderRoomResponse(request.OrderRoom),
		Content:     request.Content,
		Status:      request.Status,
		CreatedAt:   request.CreatedAt,
		UpdatedAt:   request.UpdatedAt,
		UpdatedBy:   ToBasicUserResponse(request.UpdatedBy),
	}
}

func ToBasicRequestResponse(request *model.Request) *types.BasicRequestResponse {
	if request == nil {
		return nil
	}

	return &types.BasicRequestResponse{
		ID:              request.ID,
		Code:            request.Code,
		RequestTypeName: request.RequestType.Name,
		RoomName:        request.OrderRoom.Room.Name,
		Status:          request.Status,
		CreatedAt:       request.CreatedAt,
	}
}

func ToBasicRequestsResponse(requests []*model.Request) []*types.BasicRequestResponse {
	if len(requests) == 0 {
		return make([]*types.BasicRequestResponse, 0)
	}

	requestsRes := make([]*types.BasicRequestResponse, 0, len(requests))
	for _, request := range requests {
		requestsRes = append(requestsRes, ToBasicRequestResponse(request))
	}

	return requestsRes
}
