package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/files"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/supakornn/hexagonal-go/pkg/utils"
)

type filesHandlerErrCode string

const (
	uploadedErr filesHandlerErrCode = "files-001"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg         config.Iconfig
	fileUsecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.Iconfig, fileUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)
	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(uploadedErr), err.Error()).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(uploadedErr), "invalid file type").Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(uploadedErr),
				fmt.Sprintf("fize size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)

		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(uploadedErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}
