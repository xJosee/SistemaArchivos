package analisis

import (
	"fmt"
	"unicode"
)

var cadenaconcatenar byte = 0
var estadoprincipal int = 0

func Scanner(entrada string) {
	for i := 0; i < len(entrada); i++ {
		cadenaconcatenar = entrada[i]
		switch estadoprincipal {
		case 0:
			if unicode.IsDigit(rune(cadenaconcatenar)) {
				fmt.Println("No es una letra")
			}
		}
	}
}
