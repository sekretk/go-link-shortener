package link

import (
	"go/adv-demo/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link)

	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil

}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	searchResult := repo.Database.DB.First(&link, "hash = ?", hash)

	if searchResult.Error != nil {
		return nil, searchResult.Error
	}

	return &link, nil
}

func (repo *LinkRepository) IsUniqueHash(hash string) (bool, error) {
	var resultCount int64
	result := repo.Database.DB.Model(&Link{}).Where("hash = ?", hash).Count(&resultCount)
	return resultCount == 0, result.Error
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	return link, result.Error
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Link{}, id)
	return result.Error
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, id)
	return &link, result.Error
}

func (repo *LinkRepository) Count() int64 {
	var count int64

	repo.Database.Table("links").Where("deleted_at is NULL").Count(&count)

	return count
}

func (repo *LinkRepository) GetList(limit, offset uint) []Link {
	var links []Link

	query := repo.Database.
		Table("links").
		Where("deleted_at is NULL").
		Session(&gorm.Session{})

	query.Order("id desc").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&links)

	return links
}
