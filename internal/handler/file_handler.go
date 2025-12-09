package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileSvc service.FileService
}

func NewFileHandler(fileSvc service.FileService) *FileHandler {
	return &FileHandler{fileSvc}
}

// UploadPresignedURLs godoc
// @Summary      Get Upload Presigned URLs
// @Description  Tạo một hoặc nhiều URL (presigned) để client upload file lên storage
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        payload                       body      types.UploadPresignedURLsRequest  true  "Danh sách file cần tạo URL (tên file, content type)"
// @Success      200                           {object}  types.APIResponse{data=object{presigned_urls=[]types.UploadPresignedURLResponse}}  "Tạo URL upload thành công"
// @Failure      400                           {object}  types.APIResponse  "Bad Request (validation error)"
// @Failure      500                           {object}  types.APIResponse  "Internal Server Error"
// @Router       /files/presigned-urls/uploads [post]
func (h *FileHandler) UploadPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.UploadPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	presignedURLs, err := h.fileSvc.CreateUploadURLs(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Generate upload presigned urls successfully", gin.H{
		"presigned_urls": presignedURLs,
	})
}

// ViewPresignedURLs godoc
// @Summary      Get View Presigned URLs
// @Description  Tạo một hoặc nhiều URL (presigned) để client xem/tải file từ storage
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        payload                     body      types.ViewPresignedURLsRequest  true  "Danh sách file key cần xem"
// @Success      200                         {object}  types.APIResponse{data=object{presigned_urls=[]types.ViewPresignedURLResponse}}  "Tạo URL xem thành công"
// @Failure      400                         {object}  types.APIResponse  "Bad Request (validation error)"
// @Failure      500                         {object}  types.APIResponse  "Internal Server Error"
// @Router       /files/presigned-urls/views [post]
func (h *FileHandler) ViewPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.ViewPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	presignedURLs, err := h.fileSvc.CreateViewURLs(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Generate view presigned url successfully", gin.H{
		"presigned_url": presignedURLs,
	})
}
