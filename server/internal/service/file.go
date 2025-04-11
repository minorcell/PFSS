package service

import (
	"errors"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/minorcell/pfss/internal/model"
	"gorm.io/gorm"
)

// FileService handles file-related operations
type FileService struct {
	db            *gorm.DB
	bucketService *BucketService
}

// NewFileService creates a new file service
func NewFileService(db *gorm.DB, bucketService *BucketService) *FileService {
	return &FileService{
		db:            db,
		bucketService: bucketService,
	}
}

// validateFilePath validates file path format
func validateFilePath(filePath string) error {
	if !strings.HasPrefix(filePath, "/") {
		return errors.New("file path must start with '/'")
	}
	if strings.Contains(filePath, "..") {
		return errors.New("file path cannot contain '..'")
	}
	return nil
}

// CreateFile creates a new file record
func (s *FileService) CreateFile(req *model.FileCreateRequest, userID uint, isRoot bool) (*model.File, error) {
	// Check bucket access
	if _, err := s.bucketService.GetUserBucketPermission(req.BucketID, userID); err != nil {
		return nil, errors.New("permission denied: no access to bucket")
	}

	// Validate file path
	if err := validateFilePath(req.Path); err != nil {
		return nil, err
	}

	// Check if file already exists in the bucket
	var count int64
	if err := s.db.Model(&model.File{}).
		Where("bucket_id = ? AND path = ?", req.BucketID, req.Path).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("file already exists in this path")
	}

	// Create file record
	file := &model.File{
		BucketID:    req.BucketID,
		Name:        req.Name,
		Path:        req.Path,
		ContentType: req.ContentType,
		Size:        req.Size,
		Metadata:    req.Metadata,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := s.db.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

// GetFiles returns a list of files with pagination
func (s *FileService) GetFiles(bucketID uint, userID uint, isRoot bool, page, pageSize int) ([]model.File, int64, error) {
	// Check bucket access
	if _, err := s.bucketService.GetUserBucketPermission(bucketID, userID); err != nil {
		return nil, 0, errors.New("permission denied: no access to bucket")
	}

	var files []model.File
	var total int64
	query := s.db.Model(&model.File{}).Where("bucket_id = ?", bucketID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get files with pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetFileByID returns a file by ID
func (s *FileService) GetFileByID(id uint, userID uint, isRoot bool) (*model.File, error) {
	var file model.File
	if err := s.db.First(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	// Check bucket access
	if _, err := s.bucketService.GetUserBucketPermission(file.BucketID, userID); err != nil {
		return nil, errors.New("permission denied: no access to bucket")
	}

	return &file, nil
}

// UpdateFile updates file information
func (s *FileService) UpdateFile(id uint, req *model.FileUpdateRequest, userID uint, isRoot bool) error {
	// Get file
	file, err := s.GetFileByID(id, userID, isRoot)
	if err != nil {
		return err
	}

	// Check bucket write access
	perm, err := s.bucketService.GetUserBucketPermission(file.BucketID, userID)
	if err != nil || perm.Access == "read" {
		return errors.New("permission denied: requires write access")
	}

	// Validate new file path if provided
	if req.Path != "" {
		if err := validateFilePath(req.Path); err != nil {
			return err
		}
		// Check if new path already exists
		var count int64
		if err := s.db.Model(&model.File{}).
			Where("bucket_id = ? AND path = ? AND id != ?", file.BucketID, req.Path, id).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("file already exists in this path")
		}
	}

	// Update file
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.ContentType != "" {
		updates["content_type"] = req.ContentType
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}
	updates["updated_by"] = userID

	return s.db.Model(file).Updates(updates).Error
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(id uint, userID uint, isRoot bool) error {
	// Get file
	file, err := s.GetFileByID(id, userID, isRoot)
	if err != nil {
		return err
	}

	// Check bucket write access
	perm, err := s.bucketService.GetUserBucketPermission(file.BucketID, userID)
	if err != nil || perm.Access == "read" {
		return errors.New("permission denied: requires write access")
	}

	// Delete file record
	return s.db.Delete(file).Error
}

// GetUploadURL generates a pre-signed URL for file upload
func (s *FileService) GetUploadURL(id uint, userID uint, isRoot bool) (string, error) {
	// Get file
	file, err := s.GetFileByID(id, userID, isRoot)
	if err != nil {
		return "", err
	}

	// Check bucket write access
	perm, err := s.bucketService.GetUserBucketPermission(file.BucketID, userID)
	if err != nil || perm.Access == "read" {
		return "", errors.New("permission denied: requires write access")
	}

	// TODO: Implement storage service integration to generate pre-signed upload URL
	uploadURL := "/api/v1/files/" + path.Join(strconv.FormatUint(uint64(file.BucketID), 10), file.Path)

	return uploadURL, nil
}

// GetDownloadURL generates a pre-signed URL for file download
func (s *FileService) GetDownloadURL(id uint, userID uint, isRoot bool) (string, time.Time, error) {
	// Get file
	file, err := s.GetFileByID(id, userID, isRoot)
	if err != nil {
		return "", time.Time{}, err
	}

	// TODO: Implement storage service integration to generate pre-signed download URL
	downloadURL := "/api/v1/files/" + path.Join(strconv.FormatUint(uint64(file.BucketID), 10), file.Path)
	expiresAt := time.Now().Add(24 * time.Hour) // URL expires in 24 hours

	return downloadURL, expiresAt, nil
}
