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
	id       []string
	user     string = ""
	password string = ""
	p        bool   = false
	count    string = ""
	nombre   string = ""
	grupo    string = ""
	r        bool   = false
	ugo      int    = 0
	file     []string
	rf       bool   = false
	ruta     string = ""
	dest     string = ""
)

//Analizar is...
func Analizar(comandos string) {
	if comandos != "" {
		if !strings.HasPrefix(comandos, "#") {
			Comandos := strings.Split(comandos, " ")
			VerificarComando(Comandos)
		} else {
			//Comentario(comandos)
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
	user = ""
	password = ""
	p = false
	count = ""
	nombre = ""
	grupo = ""
	r = false
	ugo = 0
	file = nil
	rf = false
	ruta = ""
	dest = ""
	id = nil
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
			if len(id) == 0 {
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
				ErrorMessage("[REP] -> Parametro -nombre no defino")
			} else if path == "" {

			} else {
				if strings.ToLower(nombre) == "mbr" {
					comandos.ReporteEBR(path, ruta)
					SuccessMessage("[REP] -> Reporte 'mbr' Generado Correctamente")
				} else if strings.ToLower(nombre) == "disk" {
					comandos.ReporteDisco(path, ruta)
					SuccessMessage("[REP] -> Reporte 'disk' Generado Correctamente")
				} else if strings.ToLower(nombre) == "sb" {
					comandos.ReporteSuperBloque(id[0], path)
					SuccessMessage("[REP] -> Reporte 'sb' Generado Correctamente")
				} else if strings.ToLower(nombre) == "bm_arbdir" {
					comandos.ReporteBMarbdir(path, id[0])
					SuccessMessage("[REP] -> Reporte 'bm_arbdir' Generado Correctamente")
				} else if strings.ToLower(nombre) == "bm_detdir" {
					comandos.ReporteBMdetdir(path, id[0])
					SuccessMessage("[REP] -> Reporte 'bm_detdir' Generado Correctamente")
				} else if strings.ToLower(nombre) == "bm_inode" {
					comandos.ReporteBMinode(path, id[0])
					SuccessMessage("[REP] -> Reporte 'bm_inode' Generado Correctamente")
				} else if strings.ToLower(nombre) == "bm_block" {
					comandos.ReporteBMblock(path, id[0])
					SuccessMessage("[REP] -> Reporte 'bm_block' Generado Correctamente")
				} else if strings.ToLower(nombre) == "bitacora" {
					comandos.ReporteBitacora(path, id[0])
					SuccessMessage("[REP] -> Reporte 'bitacora' Generado Correctamente")
				} else if strings.ToLower(nombre) == "directorio" {
					comandos.ReporteDirectorio(path, id[0])
					SuccessMessage("[REP] -> Reporte 'directorio' Generado Correctamente")
				} else if strings.ToLower(nombre) == "tree_file" {
					fmt.Print("Ingresa la ruta de la carpeta : ")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					Carpeta := scanner.Text()
					comandos.ReporteTreeFile(Carpeta, id[0], path)
					SuccessMessage("[REP] -> Reporte 'tree_file' Generado Correctamente")
				} else if strings.ToLower(nombre) == "tree_directorio" {
					fmt.Print("Ingresa la ruta de la carpeta : ")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					Carpeta := scanner.Text()
					comandos.ReporteTreeDirectorio(Carpeta, path, id[0])
					SuccessMessage("[REP] -> Reporte 'tree_directorio' Generado Correctamente")
				} else if strings.ToLower(nombre) == "tree_complete" {
					comandos.ReporteTreeComplete(path, id[0])
					SuccessMessage("[REP] -> Reporte 'tree_complete' Generado Correctamente")
				} else if strings.ToLower(nombre) == "ls" {
					comandos.ReporteLS(path, id[0])
					SuccessMessage("[REP] -> Reporte 'ls' Generado Correctamente")
				}
			}

		}

	} else if strings.ToLower(listaComandos[0]) == "login" {

		if VerificarParametros(listaComandos) {
			if user == "" {
				ErrorMessage("[LOGIN] -> Parametro -usr no definido")
			} else if password == "" {
				ErrorMessage("[LOGIN] -> Parametro -pwd no definido")
			} else if len(id) == 0 {
				ErrorMessage("[LOGIN] -> Parametro -id no definido")
			} else {
				comandos.Login(user, password, id[0])
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "mkfile" {

		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[MKFILE] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[MKFILE] -> Parametro -path no definido")
			} else {
				comandos.MKFILE(id[0], path, true, size, count, true)
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "pause" {

		fmt.Println("[CONSOLA] -> Presiona enter para continuar")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

	} else if strings.ToLower(listaComandos[0]) == "logout" {

		comandos.Logout()

	} else if strings.ToLower(listaComandos[0]) == "loss" {

		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[System] -> parametro -id no especificado")
			} else {
				comandos.SystemLoss(id[0])
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "recovery" {

		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[System] -> parametro -id no especificado")
			} else {
				comandos.SystemRecovery(id[0])
			}
		}

	} else if strings.ToLower(listaComandos[0]) == "mkfs" {
		if VerificarParametros(listaComandos) {
			comandos.MKFS(id[0])
		}
	} else if strings.ToLower(listaComandos[0]) == "mkdir" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[MKDIR] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[MKDIR] -> Parametro -path no definido")
			} else {
				comandos.ComandoMKDIR(id[0], path, p, true)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "mkgrp" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[MKGRP] -> Parametro -id no definido")
			} else if name == "" {
				ErrorMessage("[MKGRP] -> Parametro -name no definido")
			} else {
				comandos.MKGRP(id[0], name)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "mkusr" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[MKUSR] -> Parametro -id no definido")
			} else if user == "" {
				ErrorMessage("[MKUSR] -> Parametro -user no definido")
			} else if password == "" {
				ErrorMessage("[MKUSR] -> Parametro -password no definido")
			} else if grupo == "" {
				ErrorMessage("[MKUSR] -> Parametro -grupo no definido")
			} else {
				comandos.MKUSR(id[0], user, grupo, password)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "rmgrp" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[RMGRP] -> Parametro -id no definido")
			} else if name == "" {
				ErrorMessage("[RMGRP] -> Parametro -name no definido")
			} else {
				comandos.EliminarGrupo(id[0], name)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "rmusr" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[RMUSR] -> Parametro -id no definido")
			} else if name == "" {
				ErrorMessage("[RMUSR] -> Parametro -user no definido")
			} else {
				comandos.EliminarUsuario(id[0], name)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "chmod" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[CHMOD] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[CHMOD] -> Parametro -path no definido")
			} else if ugo == 0 {
				ErrorMessage("[CHMOD] -> Parametro -ugo no definido")
			} else {
				comandos.CHMOD(id[0], path, ugo, r)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "cat" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[CAT] -> Parametro -id no definido")
			} else if len(file) == 0 {
				ErrorMessage("[CAT] -> Parametro -file no definido")
			} else {
				comandos.ComandoCat(file, id[0])
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "rm" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[RM] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[RM] -> Parametro -path no definido")
			} else {
				comandos.ComandoRM(id[0], path, rf)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "cp" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[CP] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[CP] -> Parametro -path no definido")
			} else if dest == "" {
				ErrorMessage("[CP] -> Parametro -dest no definido")
			} else {
				comandos.ComandoCopy(id[0], path, dest)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "mv" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[MV] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[MV] -> Parametro -path no definido")
			} else if dest == "" {
				ErrorMessage("[MV] -> Parametro -dest no definido")
			} else {
				comandos.ComandoMove(id[0], path, dest)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "ren" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[REN] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[REN] -> Parametro -path no definido")
			} else if name == "" {
				ErrorMessage("[REN] -> Parametro -name no definido")
			} else {
				comandos.ComandoRenombrar(id[0], path, name)
			}
		}
	} else if strings.ToLower(listaComandos[0]) == "edit" {
		if VerificarParametros(listaComandos) {
			if len(id) == 0 {
				ErrorMessage("[EDIT] -> Parametro -id no definido")
			} else if path == "" {
				ErrorMessage("[EDIT] -> Parametro -path no definido")
			} else if count == "" {
				ErrorMessage("[EDIT] -> Parametro -edit no definido")
			} else if size == 0 {
				ErrorMessage("[EDIT] -> Parametro -size no definido")
			} else {
				comandos.ComandoEdit(id[0], size, path, count)
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
	} else if strings.ToLower(listaComandos[0]) == "3" {
		fmt.Println("          Universidad San Carlos de Guatemala")
		fmt.Println("          Ingenieria en Ciencias y Sistemas")
		fmt.Println("          [MIA]Manejo e Implementacion de archivos")
		fmt.Println("          Jose Luis Herrera Martinez")
		fmt.Println("          201807431")
	} else if strings.ToLower(listaComandos[0]) == " " {

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
			if strings.Contains(Paramatros[1], "\"") {
				name = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				name = Paramatros[1]
			}
		case "-type":
			tipo = strings.ToLower(Paramatros[1])
		case "-delete":
			delete = strings.ToLower(Paramatros[1])
		case "-add":
			Add, _ := strconv.Atoi(Paramatros[1]) //Convirtiendo el size a int
			add = Add
		case "-id":
			if strings.Contains(Paramatros[1], "\"") {
				aux := strings.ReplaceAll(Paramatros[1], "\"", "")
				id = append(id, aux)
			} else {
				aux := strings.ReplaceAll(Paramatros[1], "\"", "")
				id = append(id, aux)
			}
		case "-usr":
			if strings.Contains(Paramatros[1], "\"") {
				user = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				user = Paramatros[1]
			}
		case "-pwd":
			if strings.Contains(Paramatros[1], "\"") {
				password = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				password = Paramatros[1]
			}
		case "-p":
			p = true
		case "-cont":
			if strings.Contains(Paramatros[1], "\"") {
				count = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				count = Paramatros[1]
			}
		case "-nombre":
			if strings.Contains(Paramatros[1], "\"") {
				nombre = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				nombre = Paramatros[1]
			}
		case "-grp":
			if strings.Contains(Paramatros[1], "\"") {
				grupo = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				grupo = Paramatros[1]
			}
		case "-ugo":
			Aux, _ := strconv.Atoi(Paramatros[1]) //Convirtiendo el size a int
			ugo = Aux
		case "-r":
			r = true
		case "-file":
			if strings.Contains(Paramatros[1], "\"") {
				aux := strings.ReplaceAll(Paramatros[1], "\"", "")
				file = append(file, aux)
			} else {
				aux := strings.ReplaceAll(Paramatros[1], "\"", "")
				file = append(file, aux)
			}
		case "-rf":
			rf = true
		case "-ruta":
			if strings.Contains(Paramatros[1], "\"") {
				ruta = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				ruta = Paramatros[1]
			}
		case "-dest":
			if strings.Contains(Paramatros[1], "\"") {
				dest = strings.ReplaceAll(Paramatros[1], "\"", "")
			} else {
				dest = Paramatros[1]
			}
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

			if !strings.HasPrefix(scanner.Text(), "#") {
				Comando("[CONSOLA] -> " + scanner.Text())
				Analizar(scanner.Text())
			} else {
				Comentario(scanner.Text())
			}
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
	} else if unit == 'b' || unit == 'B' {
		return size
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
