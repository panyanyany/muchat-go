package model_util

import (
	"gorm.io/gorm"
	"muchat-go/util/gin_util"
)

func Paginate(params *gin_util.PaginationParams) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if params.Limit == 0 {
            params.Limit = 99
        }
        if params.Page == 0 {
            params.Page = 1
        }

        offset := (params.Page - 1) * params.Limit

        return db.Offset(offset).Limit(params.Limit)
    }
}
