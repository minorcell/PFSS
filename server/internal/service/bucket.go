package service

import (
	"errors"
	"strings"
	"time"

	"github.com/minorcell/pfss/internal/model"
	"gorm.io/gorm"
)

// BucketService handles bucket-related operations
type BucketService struct {
	db *gorm.DB
}

// NewBucketService creates a new bucket service
func NewBucketService(db *gorm.DB) *BucketService {
	return &BucketService{db: db}
}

// validateBucketName validates bucket name format
func validateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("bucket name must be between 3 and 63 characters")
	}
	if !strings.HasPrefix(name, "pfss-") {
		return errors.New("bucket name must start with 'pfss-'")
	}
	// Add more validation rules as needed
	return nil
}

// CreateBucket creates a new bucket
func (s *BucketService) CreateBucket(req *model.BucketCreateRequest, ownerID uint) (*model.Bucket, error) {
	// Validate bucket name
	if err := validateBucketName(req.Name); err != nil {
		return nil, err
	}

	// Check if bucket name already exists
	var count int64
	if err := s.db.Model(&model.Bucket{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("bucket name already exists")
	}

	// Create bucket
	bucket := &model.Bucket{
		Name:        req.Name,
		OwnerID:     ownerID,
		Description: req.Description,
	}

	// Start transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create bucket
		if err := tx.Create(bucket).Error; err != nil {
			return err
		}

		// Create owner permission (admin access)
		ownerPerm := model.BucketPermission{
			BucketID: bucket.ID,
			UserID:   ownerID,
			Access:   "admin",
		}
		if err := tx.Create(&ownerPerm).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// GetBuckets returns a list of buckets with pagination
func (s *BucketService) GetBuckets(userID uint, isRoot bool, page, pageSize int) ([]model.Bucket, int64, error) {
	var buckets []model.Bucket
	var total int64
	query := s.db.Model(&model.Bucket{})

	// If not root user, only show buckets they have access to
	if !isRoot {
		query = query.Joins("JOIN bucket_permissions ON buckets.id = bucket_permissions.bucket_id").
			Where("bucket_permissions.user_id = ?", userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get buckets with pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&buckets).Error; err != nil {
		return nil, 0, err
	}

	return buckets, total, nil
}

// GetBucketByID returns a bucket by ID
func (s *BucketService) GetBucketByID(id uint, userID uint, isRoot bool) (*model.Bucket, error) {
	var bucket model.Bucket

	query := s.db.Model(&model.Bucket{})
	if !isRoot {
		// Check if user has access to the bucket
		query = query.Joins("JOIN bucket_permissions ON buckets.id = bucket_permissions.bucket_id").
			Where("bucket_permissions.user_id = ?", userID)
	}

	if err := query.First(&bucket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bucket not found or access denied")
		}
		return nil, err
	}

	return &bucket, nil
}

// UpdateBucket updates bucket information
func (s *BucketService) UpdateBucket(id uint, req *model.BucketUpdateRequest, userID uint, isRoot bool) error {
	// Get bucket
	bucket, err := s.GetBucketByID(id, userID, isRoot)
	if err != nil {
		return err
	}

	// Check permissions
	if !isRoot {
		perm, err := s.GetUserBucketPermission(id, userID)
		if err != nil || perm.Access != "admin" {
			return errors.New("permission denied: requires admin access")
		}
	}

	// Validate new bucket name if provided
	if req.Name != "" {
		if err := validateBucketName(req.Name); err != nil {
			return err
		}
		// Check if new name already exists
		var count int64
		if err := s.db.Model(&model.Bucket{}).Where("name = ? AND id != ?", req.Name, id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("bucket name already exists")
		}
	}

	// Update bucket
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	return s.db.Model(bucket).Updates(updates).Error
}

// DeleteBucket deletes a bucket
func (s *BucketService) DeleteBucket(id uint, userID uint, isRoot bool) error {
	// Get bucket
	bucket, err := s.GetBucketByID(id, userID, isRoot)
	if err != nil {
		return err
	}

	// Check permissions
	if !isRoot {
		perm, err := s.GetUserBucketPermission(id, userID)
		if err != nil || perm.Access != "admin" {
			return errors.New("permission denied: requires admin access")
		}
	}

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete bucket permissions
		if err := tx.Where("bucket_id = ?", id).Delete(&model.BucketPermission{}).Error; err != nil {
			return err
		}

		// Delete bucket
		return tx.Delete(bucket).Error
	})
}

// GetBucketPermissions returns a list of permissions for a bucket
func (s *BucketService) GetBucketPermissions(bucketID uint, userID uint, isRoot bool) ([]model.BucketPermission, error) {
	// Check if user has access to view permissions
	if !isRoot {
		perm, err := s.GetUserBucketPermission(bucketID, userID)
		if err != nil || perm.Access != "admin" {
			return nil, errors.New("permission denied: requires admin access")
		}
	}

	var permissions []model.BucketPermission
	if err := s.db.Where("bucket_id = ?", bucketID).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserBucketPermission gets a user's permission for a bucket
func (s *BucketService) GetUserBucketPermission(bucketID, userID uint) (*model.BucketPermission, error) {
	var perm model.BucketPermission
	err := s.db.Where("bucket_id = ? AND user_id = ?", bucketID, userID).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &perm, nil
}

// UpdateBucketPermissions updates bucket permissions
func (s *BucketService) UpdateBucketPermissions(bucketID uint, permissions []model.BucketPermissionRequest, userID uint, isRoot bool) error {
	// Check if user has admin access
	if !isRoot {
		perm, err := s.GetUserBucketPermission(bucketID, userID)
		if err != nil || perm.Access != "admin" {
			return errors.New("permission denied: requires admin access")
		}
	}

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing permissions except owner's
		if err := tx.Where("bucket_id = ? AND user_id != ?", bucketID, userID).Delete(&model.BucketPermission{}).Error; err != nil {
			return err
		}

		// Add new permissions
		for _, p := range permissions {
			// Skip if trying to modify owner's permission
			if p.UserID == userID {
				continue
			}

			var expiresAt *time.Time
			if p.ExpiresAt != nil {
				t := time.Time(*p.ExpiresAt)
				expiresAt = &t
			}
			perm := model.BucketPermission{
				BucketID:  bucketID,
				UserID:    p.UserID,
				Access:    p.Access,
				ExpiresAt: expiresAt,
			}
			if err := tx.Create(&perm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetBucketStats returns statistics for a bucket
func (s *BucketService) GetBucketStats(bucketID uint, userID uint, isRoot bool) (*model.BucketStats, error) {
	// Check if user has access to the bucket
	if !isRoot {
		if _, err := s.GetUserBucketPermission(bucketID, userID); err != nil {
			return nil, errors.New("permission denied")
		}
	}

	var stats model.BucketStats
	var totalSize int64

	// Get file count
	if err := s.db.Model(&model.File{}).Where("bucket_id = ?", bucketID).Count(&stats.FileCount).Error; err != nil {
		return nil, err
	}

	// Get total size
	if err := s.db.Model(&model.File{}).Where("bucket_id = ?", bucketID).Select("COALESCE(SUM(size), 0)").Scan(&totalSize).Error; err != nil {
		return nil, err
	}
	stats.TotalSize = totalSize

	return &stats, nil
}
