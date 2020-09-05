package analisis

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	comandos "../Comandos"
	"github.com/fatih/color"
)

var (
	size     int    = 0
	path     string = ""
	fit      string = "ff"
	unit     string = "k"
	name     string = ""
	tipo     string = ""
	delete   string = ""
	add      int    = 0
	id       string = ""
	user     string = ""
	password string = ""
	p        bool   = false
	count    string = ""
	nombre   string = ""
)

//Analizar is...
func Analizar(comandos string) {
	if comandos != "" {
		if !strings.HasPrefix(comandos, "#") {
			Comandos := strings.Split(comandos, " ")
			VerificarComando(Comandos)
		} else {
			Comentario(comandos)
		}
	}
	size = 0
	path = ""
	fit = "ff"
	unit = "k"
	name = ""
	tipo = ""
	delete = ""
	add = 0
	id = ""
	user = ""
	password = ""
	p = false
	count = ""
	nombre = ""
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
				if comandos.MKDISK(size, fit[0], unit[0], path, name) {
					SuccessMessage("[MKDISK] -> Disco creado correctamente")
				} else {
					ErrorMessage("[MKDISK] -> Ya existe un disco con ese nombre")
				}
			}
		} else {
			ErrorMessage("[MKDISK] -> Algo anda mal con un parametro")
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
		} else {
			ErrorMessage("[MKDISK] -> Algo anda mal con un parametro")
		}

	} else if strings.ToLower(listaComandos[0]) == "fdisk" {
		if VerificarParametros(listaComandos) {

			if strings.ToLower(tipo) == "p" {
				//CrearParicionPrimaria
				comandos.CrearParticionPrimaria(path, CalcularSize(size, unit[0]), name, fit[0])
			} else if strings.ToLower(tipo) == "e" {
				//CrearParticionExtendida
				comandos.CrearParticionExtendida(path, CalcularSize(size, unit[0]), name, fit[0])
			} else if strings.ToLower(tipo) == "l" {
				//CrearParticionLogica
				comandos.CrearParticionLogica(path, name, CalcularSize(size, unit[0]), fit[0])
			} else if strings.ToLower(delete) == "fast" || strings.ToLower(delete) == "full" {
				//EliminarParticion
				comandos.EliminarParticion(path, name, delete)
			} else if add != 0 {
				//AgregarQuitarEspacio
				comandos.AgregarQuitarEspacio(path, name, add, unit[0])
			}

		} else {
			ErrorMessage("[MKDISK] -> Algo anda mal con un parametro")
		}

	} else if strings.ToLower(listaComandos[0]) == "mount" {

		if VerificarParametros(listaComandos) {
			if path == "" {
				ErrorMessage("[MOUNT] -> Parametro -path no especificado")
			} else if name == "" {
				ErrorMessage("[MOUNT] -> Parametro -name no especificado")
			} else {
				comandos.MOUNT(path, name)
			}

		}

	} else if strings.ToLower(listaComandos[0]) == "unmount" {

		if VerificarParametros(listaComandos) {
			if id == "" {
				ErrorMessage("[MOUNT] -> Parametro -id no especificado")
			} else {
				comandos.UNMOUNT(id)
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "exec" {

		if VerificarParametros(listaComandos) {
			EXEC(path)
			//SuccessMessage("[EXEC] -> Comando ejecutado correctamente")
		}

	} else if strings.ToLower(listaComandos[0]) == "rep" {

		if VerificarParametros(listaComandos) {

			if nombre == "" {
				ErrorMessage("")
			} else if id == "" {

			} else if path == "" {

			} else {
				if strings.ToLower(nombre) == "mbr" {
					comandos.ReporteEBR(path)
				} else if strings.ToLower(nombre) == "disk" {
					comandos.ReporteDisco(path)
				} else if strings.ToLower(nombre) == "sb" {
					comandos.ReporteSuperBloque(id)
				} else if strings.ToLower(nombre) == "bm_arbdir" {
					comandos.ReporteBMarbdir(path, id)
				} else if strings.ToLower(nombre) == "bm_detdir" {
					comandos.ReporteBMdetdir(path, id)
				} else if strings.ToLower(nombre) == "bm_inode" {
					comandos.ReporteBMinode(path, id)
				} else if strings.ToLower(nombre) == "bm_block" {
					comandos.ReporteBMblock(path, id)
				} else if strings.ToLower(nombre) == "bitacora" {
					comandos.ReporteBitacora(path, id)
				} else if strings.ToLower(nombre) == "tree_file" {
					fmt.Print("Ingresa el nombre de la carpeta : ")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					Carpeta := scanner.Text()
					comandos.ReporteTreeFile(Carpeta, id, path)
				} else if strings.ToLower(nombre) == "tree_directorio" {
					comandos.ReporteDirectorio(path, id)
				} else if strings.ToLower(nombre) == "tree_complete" {
					comandos.ReporteTreeComplete(path, id)
				} else if strings.ToLower(nombre) == "ls" {
					comandos.ReporteLS(path, id)
				}
			}

		}

	} else if strings.ToLower(listaComandos[0]) == "login" {

		if VerificarParametros(listaComandos) {
			if user == "" {
				ErrorMessage("[LOGIN] -> Parametro -usr no definido")
			} else if password == "" {
				ErrorMessage("[LOGIN] -> Parametro -pwd no definido")
			} else if id == "" {
				ErrorMessage("[LOGIN] -> Parametro -id no definido")
			} else {
				comandos.Login(user, password, id)
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "mkfile" {

		if VerificarParametros(listaComandos) {
			if id == "" {
				ErrorMessage("[MKFILE] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[MKFILE] -> Parametro -path no definido")
			} else {
				comandos.MKFILE(id, path, true, size, count)
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "pause" {

		fmt.Println("[CONSOLA] -> Presiona enter para continuar")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

	} else if strings.ToLower(listaComandos[0]) == "mkfs" {
		if VerificarParametros(listaComandos) {
			comandos.MKFS(id)
		}
	} else if strings.ToLower(listaComandos[0]) == "mkdir" {
		if VerificarParametros(listaComandos) {
			if id == "" {
				ErrorMessage("[MKDIR] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[MKDIR] -> Parametro -path no definido")
			} else {
				comandos.ComandoMKDIR(id, path, p)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "1" {
		fmt.Println("")
		fmt.Println(" - mkdisk")
		fmt.Println(" - rmdisk")
		fmt.Println(" - fdisk")
		fmt.Println(" - mount")
		fmt.Println(" - unmount")
		fmt.Println(" - exec")
		fmt.Println(" - rep")
		fmt.Println("")
	} else if strings.ToLower(listaComandos[0]) == "2" {
		fmt.Println("")
		fmt.Println(" - mbr")
		fmt.Println(" - disk")
		fmt.Println(" - sb")
		fmt.Println(" - bm_ardir")
		fmt.Println(" - bm_detdir")
		fmt.Println(" - bm_inode")
		fmt.Println(" - bm_block")
		fmt.Println(" - bitacora")
		fmt.Println(" - tree_file")
		fmt.Println(" - tree_directorio")
		fmt.Println(" - tree_complete")
		fmt.Println(" - ls")
		fmt.Println("")
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
			if Size < 0 {
				return false
			}
			size = Size
		case "-path":
			if strings.Contains(Paramatros[1], "\"") {
				path = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				path = Paramatros[1]
			}
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
		case "-id":
			id = Paramatros[1]
		case "-usr":
			user = Paramatros[1]
		case "-pwd":
			password = Paramatros[1]
		case "-p":
			p = true
		case "-cont":
			count = Paramatros[1]
		case "-nombre":
			nombre = Paramatros[1]
		default:
			ErrorMessage("[CONSOLA] -> Parametro [" + Paramatros[0] + "] incorrecto")
			return false
		}
	}

	return true
}

//EXEC is...
func EXEC(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			Comando("[CONSOLA] -> " + scanner.Text())
			Analizar(scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

//CalcularSize is ...
func CalcularSize(size int, unit byte) int {
	if unit == 'M' || unit == 'm' {
		return size * 1024 * 1024
	} else if unit == 'K' || unit == 'k' {
		return size * 1024
	}
	return 0
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

//Comando is...
func Comando(message string) {
	green := color.New(color.FgHiGreen)
	boldgreen := green.Add(color.Bold)
	boldgreen.Println(message)
}

//Comentario is...
func Comentario(message string) {
	white := color.New(color.FgWhite)
	boldwhite := white.Add(color.Bold)
	boldwhite.Println("[COMENTARIO] -> ", message)
}
