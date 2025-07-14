package models

type SammlungProdukt struct {
	SammlungID uint `gorm:"primaryKey"`                   // Teil des zusammengesetzten PK
	ProduktID  uint `gorm:"primaryKey;column:produkt_id"` // Teil des zusammengesetzten PK, Spaltenname beachten
	// Optional: Relationen zurück, falls benötigt
	// Sammlung   Sammlung `gorm:"foreignKey:SammlungID"`
	// Produkt    Produkt  `gorm:"foreignKey:ProduktID"`
}

func (SammlungProdukt) TableName() string {
	// Name der Verknüpfungstabelle, wie von Prisma (oder dir) definiert
	return "sammlung_produkte"
}
