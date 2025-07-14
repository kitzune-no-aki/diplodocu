package models

type Manga struct {
	ProdukteID uint    `gorm:"primaryKey"`
	Mangaka    *string `gorm:"type:varchar(255)"`
	Sprache    *string `gorm:"type:varchar(50)"`
	Genre      *string `gorm:"type:varchar(100)"`
	Produkt    Produkt `gorm:"foreignKey:ProdukteID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Manga) TableName() string {
	return "manga"
}
