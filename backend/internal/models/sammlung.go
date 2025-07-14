package models

type Sammlung struct {
	ID        uint      `gorm:"primaryKey"` // Auto-increment -> uint
	WebuserID string    `gorm:"column:webuser_id;not null;type:varchar(255)"`
	Name      *string   `gorm:"type:varchar(255)"`
	Webuser   Webuser   `gorm:"foreignKey:WebuserID;references:ID;constraint:OnDelete:CASCADE"`
	Produkte  []Produkt `gorm:"many2many:sammlung_produkte;"`
}

func (Sammlung) TableName() string {
	return "sammlung"
}
