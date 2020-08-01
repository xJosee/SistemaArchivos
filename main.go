package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
	mbr_tamano         int
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
	var comando string
	fmt.Scanln(&comando)
	Peticion(comando)
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
	fmt.Println("response Body:", string(body))
}

func recorrerAST() {

}

func comandoMKDISK(size int, fit string, unit string, path string) {

}
