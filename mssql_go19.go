// license

// TODO +build go1.9
// TODO support other data types here.

// +build go1.8

package mssql

import (
	"database/sql"
	"database/sql/driver"

	// "github.com/cockroachdb/apd"
)

func (c *MssqlConn) NamedValueCheck(nv *driver.NamedValue) error {
	switch v := nv.Value.(type) {
	case sql.Out:
		if c.outs == nil {
			c.outs = make(map[string]interface{})
		}
		c.outs[nv.Name] = v.Value

		// Unwrap the Out value and check the inner value.
		lnv := *nv
		lnv.Value = v.Value
		err := c.NamedValueCheck(&lnv)
		if err != nil {
			if err != driver.ErrSkip {
				return err
			}
			lnv.Value, err = driver.DefaultParameterConverter.ConvertValue(lnv.Value)
			if err != nil {
				return err
			}
		}
		nv.Value = sql.Out{Value: lnv.Value}
		return nil
	// case *apd.Decimal:
	// 	return nil
	default:
		return driver.ErrSkip
	}
}
