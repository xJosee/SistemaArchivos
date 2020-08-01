package main

import (
	"fmt"

	Analisis "./analisis"
)

func main() {
	menu()
}

func menu() {
	menu := "Bienvenido a la consola de comandos\n>> "
	fmt.Print(menu)

	var comando string
	fmt.Scanln(&comando)
	Analisis.Scanner(comando)
}
