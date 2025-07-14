package models

type Filmserie struct {
	ProdukteID uint    `gorm:"primaryKey"`
	Art        *string `gorm:"type:enum_filmserie_art"`
	Genre      *string `gorm:"type:varchar(100)"`
	Produkt    Produkt `gorm:"foreignKey:ProdukteID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Filmserie) TableName() string {
	return "filmserie"
}
