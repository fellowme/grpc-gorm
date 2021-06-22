package base_model

type BaseModel struct {
	Id         int32 `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
	IsDeleted  bool  `gorm:"comment:'是否删除0未删除 1删除'" json:"is_deleted"`
}
