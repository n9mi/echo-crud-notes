package domain

type Category struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null"`
}
