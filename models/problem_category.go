package models

import (
	"gorm.io/gorm"
)

type ProblemCategory struct {
	gorm.Model
	ProblemId     uint           `gorm:"column:problem_id;type:varchar(36);" json:"problem_id"`       //问题的唯一标识
	CategoryId    uint           `gorm:"column:category_id;type:varchar(36);" json:"category_id"`     //分类的唯一标识
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id;" json:"category_basic"` // 关联分类的基础信息表}
}

func (table *ProblemCategory) TableName() string {
	return "problem_category"
}
