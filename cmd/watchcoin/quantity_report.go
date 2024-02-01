package watchcoin

import (
	"encoding/json"

	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ReportItem struct {
	gorm.Model
	Data string
}

type QuantityReport struct {
	db *gorm.DB
}

func NewQuantityReport(fpath string) (*QuantityReport, error) {
	db, err := gorm.Open(sqlite.Open(fpath), &gorm.Config{})

	report := QuantityReport{
		db: db,
	}

	if err != nil {
		return &report, err
	}

	err = db.AutoMigrate(&ReportItem{}, &xeggexlib.Bids{})
	if err != nil {
		return &report, err
	}

	return &report, nil

}

func (report *QuantityReport) AddReport(data *SumBooks) error {

	dataraw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rep := ReportItem{
		Data: string(dataraw),
	}

	err = report.db.Save(&rep).Error
	return err
}
