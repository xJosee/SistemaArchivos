package Analisis

import (
	"fmt"
	"strings"
)

func Analizar(comandos string) {
	Comandos := strings.SplitAfter(comandos, " ")
	for i := 0; i < len(Comandos); i++ {
		switch Comandos[i] {
		case "mkdisk":
			fmt.Println("Si jalo")
		}
	}

}
