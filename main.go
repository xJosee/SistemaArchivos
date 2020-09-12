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
	blue := color.New(color.FgCyan)
	boldblue := blue.Add(color.Bold)
	boldblue.Println("╔════════════════════════════════════════════════════════════════════════╗")
	boldblue.Println("║                                                                        ║")
	boldblue.Println("║                            B I E N V E N I D O                         ║")
	boldblue.Println("║                    C O N S O L A  D E  C O M A N D O S                 ║")
	boldblue.Println("║                                                                        ║")
	boldblue.Println("╠════════════════════════════════════════════════════════════════════════╣")
	boldblue.Println("║                                                                        ║")
	boldblue.Println("║                        1. Comandos Disponibles                         ║")
	boldblue.Println("║                        2. Reportes Disponibles                         ║")
	boldblue.Println("║                        3. Acerca de...                                 ║")
	boldblue.Println("║                                                                        ║")
	boldblue.Println("╚════════════════════════════════════════════════════════════════════════╝")
	Menu()
}

//Menu is...
func Menu() {
	for {
		fmt.Print("  >> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Comando := scanner.Text()
		if strings.ToLower(Comando) == "salir" {
			return
		}
		parser.Analizar(Comando)
	}
}
