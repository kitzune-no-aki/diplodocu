package models

type Webuser struct {
	ID         string     `gorm:"primaryKey;type:varchar(255);not null"`
	Name       *string    `gorm:"type:varchar(255)"`    // Nullable String -> *string
	Sammlungen []Sammlung `gorm:"foreignKey:WebuserID"` // Ein User hat viele Sammlungen
}

func (Webuser) TableName() string {
	return "webuser" // Exakter Tabellenname
}
