// Package mode считывает тип сервера
package mode

import (
	// Libraries
	"flag"
)

// Type возвращает тип запускаемого сервера
func Type() *string {
	mode := flag.String("mode", "", "")
	flag.Parse()

	return mode
}
