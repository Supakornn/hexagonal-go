package appinfoHandlers

import (
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface {
}

type appinfoHandler struct {
	cfg            config.Iconfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.Iconfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{}
}
