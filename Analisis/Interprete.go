package analisis

import (
	"strconv"
	"strings"

	comandos "../Comandos"
	"github.com/fatih/color"
)

var size int
var path string
var fit string = "ff"
var unit string = "k"
var name string

//Analizar is...
func Analizar(comandos string) {
	Comandos := strings.Split(comandos, " ")
	VerificarComando(Comandos)
}

//VerificarComando is...
func VerificarComando(listaComandos []string) {

	if strings.ToLower(listaComandos[0]) == "mkdisk" {

		if VerificarParametros(listaComandos) {
			comandos.MKDISK(size, fit[0], unit[0], path, name)
			SuccessMessage("[MKDISK] -> Comando ejecutado correctamente")
		}

	} else if strings.ToLower(listaComandos[0]) == "rmdisk" {

		if VerificarParametros(listaComandos) {
			//comandos.RMDISK(path)
			SuccessMessage("[RMDISK] -> Comando ejecutado correctamente")
		}

	} else if strings.ToLower(listaComandos[0]) == "mount" {

		if VerificarParametros(listaComandos) {
			//comandos.MOUNT()
			SuccessMessage("[MOUNT] -> Comando ejecutado correctamente")
		}

	} else if strings.ToLower(listaComandos[0]) == "unmount" {

		if VerificarParametros(listaComandos) {
			//comandos.UNMOUNT(path)
			SuccessMessage("[UNMOUNT] -> Comando ejecutado correctamente")
		}

	} else if strings.ToLower(listaComandos[0]) == "fdisk" {

		if VerificarParametros(listaComandos) {
			//comandos.FDISK(path)
			SuccessMessage("[FDISK] -> Comando ejecutado correctamente")
		}

	} else {
		ErrorMessage("[CONSOLA] -> Comando [" + listaComandos[0] + "] incorrecto")
	}

}

//VerificarParametros is...
func VerificarParametros(listaComandos []string) bool {
	for i := 1; i < len(listaComandos); i++ {
		Paramatros := strings.Split(listaComandos[i], "=")
		switch strings.ToLower(Paramatros[0]) {
		case "-size":
			Size, _ := strconv.Atoi(Paramatros[1]) //Convirtiendo el size a string
			size = Size
		case "-path":
			path = Paramatros[1]
		case "-fit":
			fit = Paramatros[1]
		case "-unit":
			unit = Paramatros[1]
		case "-name":
			name = Paramatros[1]
		default:
			ErrorMessage("[CONSOLA] -> Parametro [" + Paramatros[0] + "] incorrecto")
			return false
		}
	}

	return true
}

//ErrorMessage is..
func ErrorMessage(message string) {
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	boldRed.Println(message)
}

//SuccessMessage is...
func SuccessMessage(message string) {
	yellow := color.New(color.FgHiYellow)
	boldyellow := yellow.Add(color.Bold)
	boldyellow.Println(message)
}
