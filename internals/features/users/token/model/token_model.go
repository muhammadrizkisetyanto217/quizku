package model

type Token struct {
    ID    int    `gorm:"primaryKey;autoIncrement" json:"id"`
    Token string `gorm:"type:text;not null" json:"token"`
}

func (Token) TableName() string {
    return "tokens"
}
