package Analisis

import (
	"fmt"
	"strings"

	comandos "../Comandos"
)

var size string
var path string
var fit string

func Analizar(comandos string) {
	Comandos := strings.Split(comandos, " ")
	VerificarComando(Comandos)
}

func VerificarComando(listaComandos []string) {
	if listaComandos[0] == "mkdisk" {
		for i := 1; i < len(listaComandos); i++ {
			Atributos := strings.Split(listaComandos[i], "=")
			switch Atributos[0] {
			case "-size":
				size = Atributos[1]
			case "-path":
				path = Atributos[1]
			case "-fit":
				fit = Atributos[1]
			}
		}
		comandos.MKDISK(22, 'F', 'K', "/home/jose/Escritorio/", "test")
		fmt.Println("Comando mkdisk ejecutado correctamente")
	} else if listaComandos[0] == "rmdisk" {
		fmt.Println("Comando rmdisk")
	}
}
