package models

type Produkt struct {
	ID         uint       `gorm:"primaryKey"`
	Name       string     `gorm:"not null;type:varchar(255)"`
	Nummer     *int       // Nullable Int -> *int
	Art        string     `gorm:"not null;type:varchar(255)"`   // Diskriminator-Spalte
	Sammlungen []Sammlung `gorm:"many2many:sammlung_produkte;"` // Many-to-Many Beziehung zu Sammlung
	// Keine direkten Felder für Buch, Manga etc. hier. Abfrage erfolgt separat.
}

func (Produkt) TableName() string {
	return "produkte"
}
