package services

import (
	"blog/config"
	"blog/models"

	"gorm.io/gorm"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) CreatePost(post *models.BlogPost) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *PostService) UpdatePost(post *models.BlogPost) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(post).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *PostService) DeletePost(id string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.BlogPost{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *PostService) GetPost(id string) (*models.BlogPost, error) {
	var post models.BlogPost
	if err := config.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetAllPosts(page int, limit int) ([]models.BlogPost, int64, error) {
	var posts []models.BlogPost
	var total int64

	offset := (page - 1) * limit

	if err := config.DB.Model(&models.BlogPost{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := config.DB.Offset(offset).Limit(limit).Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *PostService) GetPostBySlug(slug string) (*models.BlogPost, error) {
	var post models.BlogPost
	if err := config.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
