package domain

import (
	"fmt"
	"time"
)

const (
	ISO8601 = "2006-01-02"
)

// ProductRecord represents a table that stores the different prices Product had over time.
type ProductRecord struct {
	ID             int       `json:"id"`
	LastUpdateDate MySqlTime `json:"last_update_date"`
	PurchasePrice  float32   `json:"purchase_price"`
	SalePrice      float32   `json:"sale_price"`
	ProductID      int       `json:"product_id"`
}

type MySqlTime struct {
	time.Time
}

func (mst MySqlTime) MarshalJSON() ([]byte, error) {
	lastUpdateDate := fmt.Sprintf("\"%s\"", mst.Format(ISO8601))
	return []byte(lastUpdateDate), nil
}

func (mst *MySqlTime) UnmarshalJSON(data []byte) (err error) {
	mst.Time, err = time.Parse(ISO8601, string(data[1:len(data)-1]))
	return err
}
