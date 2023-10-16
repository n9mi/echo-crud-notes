package domain

type Note struct {
	ID         int    `gorm:"primaryKey"`
	Title      string `gorm:"type:varchar(100);not null"`
	Body       string `gorm:"type:varchar(255);not null"`
	CategoryID int
	Category   Category
}

type ScanNote struct {
	ID       int
	Title    string
	Body     string
	Category string
}
