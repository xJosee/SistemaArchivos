package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Particion struct {
	part_status byte
	part_type   byte
	part_fit    byte
	part_start  int
	part_size   int
	part_name   [16]byte
}

type MBR struct {
	mbr_size           int
	mbr_fecha_creacion string
	mbr_disk_signature int
	disk_fit           byte
	mbr_particion      [4]Particion
}
type EBR struct {
	part_status byte
	part_fit    byte
	part_start  int
	part_size   int
	part_next   int
	part_name   [16]byte
}

/*type Comando struct {
	Nombre    string
	Atributos []Atributo
}

type Atributo struct {
	Nombre string
	Valor  string
}*/

func main() {
	Menu()
}

func Menu() {
	Menu := "Bienvenido a la consola de comandos\n>> "
	fmt.Print(Menu)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	Comando := scanner.Text()
	//Peticion(Comando)
	recorrerAST(Comando)
}

func Peticion(comando string) {
	url := "http://localhost:2020"

	var jsonStr = []byte(`{"comando":"` + comando + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	AST := string(body)
	recorrerAST(AST)
}
func recorrerAST(ast string) {
	comandoMKDISK(10, 'F', 'K', "/home/jose/Escritorio/test.disk")
}
func comandoMKDISK(size int, fit byte, unit byte, path string) {
	var Disco MBR
	Disco.mbr_disk_signature = 15
	Disco.disk_fit = fit
	Disco.mbr_size = CalcularSize(size, unit)
}

func CalcularSize(size int, unit byte) int {
	if unit == 'M' {
		return size * 1024 * 1024
	} else if unit == 'K' {
		return size * 1024
	}
	return 0
}

func VerificarRuta(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
