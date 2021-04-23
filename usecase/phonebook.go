package usecase

import (
	"context"
	"database/sql"

	"github.com/sapawarga/phonebook-service/helper"
	"github.com/sapawarga/phonebook-service/model"
	"github.com/sapawarga/phonebook-service/repository"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// PhoneBook ...
type PhoneBook struct {
	repo   repository.PhoneBookI
	logger kitlog.Logger
}

// NewPhoneBook ...
func NewPhoneBook(repo repository.PhoneBookI, logger kitlog.Logger) *PhoneBook {
	return &PhoneBook{
		repo:   repo,
		logger: logger,
	}
}

// GetList ...
func (pb *PhoneBook) GetList(ctx context.Context, params *model.ParamsPhoneBook) (*model.PhoneBookWithMeta, error) {
	logger := kitlog.With(pb.logger, "method", "GetList")
	var limit, page, offset int64 = 10, 1, 0
	if params.Limit != nil {
		limit = helper.GetInt64FromPointer(params.Limit)
	}
	if params.Page != nil {
		page = helper.GetInt64FromPointer(params.Page)
	}
	offset = (page - 1) * limit

	req := &model.GetListRequest{
		Search:     params.Search,
		RegencyID:  params.RegencyID,
		DistrictID: params.DistrictID,
		VillageID:  params.VillageID,
		Status:     params.Status,
		Limit:      &limit,
		Offset:     &offset,
	}

	resp, err := pb.repo.GetListPhoneBook(ctx, req)
	if err != nil {
		level.Error(logger).Log("error", err)
		return nil, err
	}

	data := make([]*model.Phonebook, 0)

	for _, v := range resp {
		result := &model.Phonebook{
			ID:           v.ID,
			PhoneNumbers: v.PhoneNumbers.String,
			Description:  v.Description.String,
			Name:         v.Name.String,
			Address:      v.Address.String,
			Latitude:     v.Latitude.String,
			Longitude:    v.Longitude.String,
			Status:       v.Status.Int64,
		}
		if v.CategoryID.Valid {
			categoryName, err := pb.repo.GetCategoryNameByID(ctx, v.CategoryID.Int64)
			if err != nil && err != sql.ErrNoRows {
				level.Error(logger).Log("error_get_category", err)
				return nil, err
			}
			result.Category = categoryName
		}
		data = append(data, result)
	}

	total, err := pb.repo.GetMetaDataPhoneBook(ctx, req)
	if err != nil {
		level.Error(logger).Log("error", err)
		return nil, err
	}

	return &model.PhoneBookWithMeta{
		PhoneBooks: data,
		Page:       page,
		Total:      total,
	}, nil
}

// GetDetail ...
func (pb *PhoneBook) GetDetail(ctx context.Context, id int64) (*model.PhonebookDetail, error) {
	logger := kitlog.With(pb.logger, "method", "GetDetail")

	resp, err := pb.repo.GetPhonebookDetailByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("error", err)
		return nil, err
	}

	result := &model.PhonebookDetail{
		ID:             resp.ID,
		Name:           resp.Name.String,
		Address:        resp.Address.String,
		Description:    resp.Description.String,
		PhoneNumbers:   resp.PhoneNumbers.String,
		Latitude:       resp.Latitude.String,
		Longitude:      resp.Longitude.String,
		CoverImagePath: resp.CoverImagePath.String,
		Status:         resp.Status.Int64,
		CreatedAt:      resp.CreatedAt.Time,
		UpdatedAt:      resp.UpdatedAt.Time,
	}

	if resp.CategoryID.Valid {
		categoryName, err := pb.repo.GetCategoryNameByID(ctx, resp.CategoryID.Int64)
		if err != nil {
			level.Error(logger).Log("error_get_category", err)
			return nil, err
		}
		result.CategoryID = resp.CategoryID.Int64
		result.CategoryName = categoryName
	}

	if resp.RegencyID.Valid {
		regencyName, err := pb.repo.GetLocationNameByID(ctx, resp.RegencyID.Int64)
		if err != nil {
			level.Error(logger).Log("error_get_regency", err)
			return nil, err
		}
		result.RegencyID = resp.RegencyID.Int64
		result.RegencyName = regencyName
	}

	if resp.DistrictID.Valid {
		districtName, err := pb.repo.GetLocationNameByID(ctx, resp.DistrictID.Int64)
		if err != nil {
			level.Error(logger).Log("error_get_district", err)
			return nil, err
		}
		result.DistrictID = resp.DistrictID.Int64
		result.DistrictName = districtName
	}

	if resp.VillageID.Valid {
		villageName, err := pb.repo.GetLocationNameByID(ctx, resp.VillageID.Int64)
		if err != nil {
			level.Error(logger).Log("error_get_village", err)
			return nil, err
		}
		result.VillageID = resp.VillageID.Int64
		result.VillageName = villageName
	}
	return result, nil
}

// Insert ...
func (pb *PhoneBook) Insert(ctx context.Context, params *model.AddPhonebook) error {
	logger := kitlog.With(pb.logger, "method", "Insert")
	if params.CategoryID != nil {
		_, err := pb.repo.GetCategoryNameByID(ctx, helper.GetInt64FromPointer(params.CategoryID))
		if err != nil {
			level.Error(logger).Log("error_get_category", err)
			return err
		}
	}

	if err := pb.repo.Insert(ctx, params); err != nil {
		level.Error(logger).Log("error", err)
		return err
	}
	return nil
}

// Update ...
func (pb *PhoneBook) Update(ctx context.Context, params *model.UpdatePhonebook) error {
	// TODO: update phonebook
	logger := kitlog.With(pb.logger, "method", "Update")
	if _, err := pb.repo.GetPhonebookDetailByID(ctx, params.ID); err != nil {
		level.Error(logger).Log("error_get_detail", err)
		return err
	}

	if err := pb.repo.Update(ctx, params); err != nil {
		level.Error(logger).Log("error_update", err)
		return err
	}

	return nil
}

// Delete ...
func (pb *PhoneBook) Delete(ctx context.Context, id int64) error {
	// TODO: delete phonebook
	logger := kitlog.With(pb.logger, "method", "Delete")
	if _, err := pb.repo.GetPhonebookDetailByID(ctx, id); err != nil {
		level.Error(logger).Log("error_get_detail", err)
		return err
	}

	if err := pb.repo.Delete(ctx, id); err != nil {
		level.Error(logger).Log("error_delete", err)
		return err
	}

	return nil
}
