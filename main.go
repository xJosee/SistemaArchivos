package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	parser "./Analisis"
	"github.com/fatih/color"
)

func main() {
	color.HiCyan("Bienvenido a la consola de comandos")
	Menu()
}

//Menu is...
func Menu() {
	for {
		fmt.Print(">> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Comando := scanner.Text()
		if strings.ToLower(Comando) == "salir" {
			break
		}
		parser.Analizar(Comando)
	}
}
