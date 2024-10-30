package appinfoUsecases

import appinforepositories "github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"

type IAppinfoUsecase interface {
}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{appinfoRepository: appinfoRepository}
}
