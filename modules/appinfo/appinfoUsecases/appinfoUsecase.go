package appinfoUsecases

import (
	"github.com/Supakornn/hexagonal-go/modules/appinfo"
	"github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{appinfoRepository: appinfoRepository}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.appinfoRepository.InsertCategory(req); err != nil {
		return err
	}

	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appinfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}

	return nil
}
