package models

type Spiel struct {
	ProdukteID uint    `gorm:"primaryKey"`
	Konsole    *string `gorm:"type:varchar(100)"`
	Genre      *string `gorm:"type:varchar(100)"`
	Produkt    Produkt `gorm:"foreignKey:ProdukteID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Spiel) TableName() string {
	return "spiel"
}
