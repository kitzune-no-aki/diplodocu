package models

type Buch struct {
	ProdukteID uint    `gorm:"primaryKey"` // Ist PK und FK zugleich
	Autor      *string `gorm:"type:varchar(255)"`
	Sprache    *string `gorm:"type:varchar(50)"`
	Genre      *string `gorm:"type:varchar(100)"`
	Produkt    Produkt `gorm:"foreignKey:ProdukteID;references:ID;constraint:OnDelete:CASCADE"` // Optional: Referenz zurück zum Basisprodukt
}

func (Buch) TableName() string {
	return "buch"
}
