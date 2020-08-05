package main

import (
	"bufio"
	"fmt"
	"os"

	parser "./Analisis"
)

func main() {
	Menu()
}

//Menu is...
func Menu() {
	Menu := "Bienvenido a la consola de comandos\n>> "
	fmt.Print(Menu)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	Comando := scanner.Text()
	parser.Analizar(Comando)
}
