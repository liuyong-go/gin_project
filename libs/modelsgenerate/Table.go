package modelsgenerate

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}
