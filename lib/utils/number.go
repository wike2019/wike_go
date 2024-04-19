package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

// 数学计算
func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}

// AddDecimal 加法
func AddDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Add(d2)
}

// SubDecimal 减法
func SubDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Sub(d2)
}

// MulDecimal 乘法
func MulDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Mul(d2)
}

// DivDecimal 除法
func DivDecimal(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Div(d2)
}

// IntDecimal 转换成int
func IntDecimal(d decimal.Decimal) int64 {
	return d.IntPart()
}

// DecimalFloat 转换成float
func DecimalFloat(d decimal.Decimal) float64 {
	f, exact := d.Float64()
	if !exact {
		return f
	}
	return 0
}

// FloatDecimal 浮点型到decimal
func FloatDecimal(d float64) decimal.Decimal {
	return decimal.NewFromFloat(d)
}
