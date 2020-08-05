package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	parser "./Analisis"
)

type payload struct {
	One   float32
	Two   float64
	Three uint32
}

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

func main() {
	Menu()
}

func Menu() {
	Menu := "Bienvenido a la consola de comandos\n>> "
	fmt.Print(Menu)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	Comando := scanner.Text()
	parser.Analizar(Comando)
	comandoMKDISK(22, 'F', 'K', "/home/jose/Escritorio/", "test")
}

/*
 * FUNCIONES UTILIZADAS PARA EL COMANDO MKDISK
 */

func comandoMKDISK(size int, fit byte, unit byte, path string, name string) {
	err := ioutil.WriteFile(path+name+".disk", []byte(""), 0755)
	if err == nil {
		writeFile(path+name+".disk", CalcularSize(size, unit))
		CrearRaid(size, fit, unit, path, name)
	}
}
func CrearRaid(size int, fit byte, unit byte, path string, name string) {
	err := ioutil.WriteFile(path+name+"Raid.disk", []byte(""), 0755)
	if err == nil {
		writeFile(path+name+"Raid.disk", CalcularSize(size, unit))
	}
}
func writeFile(path string, size int) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {

		s := &payload{
			r.Float32(),
			r.Float64(),
			r.Uint32(),
		}
		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.BigEndian, s)
		//b :=bin_buf.Bytes()
		//l := len(b)
		//fmt.Println(l)
		writeNextBytes(file, bin_buf.Bytes())

	}
}
func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}
func CalcularSize(size int, unit byte) int {
	if unit == 'M' {
		return 63 * size * 1000
	} else if unit == 'K' {
		return 63 * size
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

/*
 * FUNCIONES UTILIZADAS PARA EL COMANDO RMDISK
 */
