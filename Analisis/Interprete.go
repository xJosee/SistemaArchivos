package analisis

import (
	"strconv"
	"strings"

	comandos "../Comandos"
	"github.com/fatih/color"
)

var (
	size   int
	path   string
	fit    string = "ff"
	unit   string = "k"
	name   string
	tipo   string
	delete string
	add    int
)

//Analizar is...
func Analizar(comandos string) {
	Comandos := strings.Split(comandos, " ")
	VerificarComando(Comandos)
}

//VerificarComando is...
func VerificarComando(listaComandos []string) {

	if strings.ToLower(listaComandos[0]) == "mkdisk" {

		if VerificarParametros(listaComandos) {
			if path == "" {
				ErrorMessage("[MKDISK] -> Parametro -path no especificado")
			} else if name == "" {
				ErrorMessage("[MKDISK] -> Parametro -name no especificado")
			} else if size == 0 {
				ErrorMessage("[MKDISK] -> Parametro -size no especificado")
			} else {
				comandos.MKDISK(size, fit[0], unit[0], path, name)
				SuccessMessage("[MKDISK] -> Comando ejecutado correctamente")
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "rmdisk" {

		if VerificarParametros(listaComandos) {
			if path != "" {
				if comandos.RMDISK(path) {
					SuccessMessage("[RMDISK] -> Comando ejecutado correctamente")
				}
			} else {
				ErrorMessage("[RMDISK] -> Parametro -path no especificado")
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "fdisk" {
		Bandera := true
		if VerificarParametros(listaComandos) {
			if size == 0 {
				ErrorMessage("[FDISK] -> Parametro -size no especificado")
				Bandera = false
			} else if path == "" {
				ErrorMessage("[FDISK] -> Parametro -path no especificado")
				Bandera = false
			} else if name == "" {
				ErrorMessage("[FDISK] -> Parametro -name no especificado")
				Bandera = false
			} else if tipo != "" {
				if tipo != "p" && tipo != "e" && tipo != "l" {
					ErrorMessage("[FDISK] -> Valor del parametro -type incorrecto")
					Bandera = false
				}
			} else if delete != "" {
				if delete != "full" && delete != "fast" {
					ErrorMessage("[FDISK] -> Valor del parametro -delete incorrecto")
					Bandera = false
				}
			}

			if Bandera {
				comandos.FDISK(size, unit[0], path, tipo[0], fit[0], delete, name, add)
				SuccessMessage("[FDISK] -> Comando ejecutado correctamente")
			}
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

	} else if strings.ToLower(listaComandos[0]) == "exec" {

		if VerificarParametros(listaComandos) {
			//comandos.FDISK(path)
			SuccessMessage("[EXEC] -> Comando ejecutado correctamente")
		}

	} else {
		ErrorMessage("[CONSOLA] -> Comando [" + listaComandos[0] + "] incorrecto")
	}

}

//VerificarParametros is...
func VerificarParametros(listaComandos []string) bool {
	for i := 1; i < len(listaComandos); i++ {
		Paramatros := strings.Split(listaComandos[i], "->")
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
		case "-type":
			tipo = strings.ToLower(Paramatros[1])
		case "-delete":
			delete = strings.ToLower(Paramatros[1])
		case "-add":
			Add, _ := strconv.Atoi(Paramatros[1]) //Convirtiendo el size a string
			add = Add
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
