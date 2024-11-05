package appinfoUsecases

import (
	"github.com/Supakornn/hexagonal-go/modules/appinfo"
	appinforepositories "github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error)
}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{appinfoRepository: appinfoRepository}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}

	return category, nil
}
