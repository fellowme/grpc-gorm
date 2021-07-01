package base_model

type BaseModel struct {
	Id         int32 `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	CreateTime int64 `gorm:"comment:'创建时间'" json:"create_time"`
	UpdateTime int64 `gorm:"comment:'更新时间'" json:"update_time"`
	IsDeleted  bool  `gorm:"comment:'是否删除0未删除 1删除'" json:"is_deleted"`
}
