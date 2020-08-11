package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"unsafe"

	"github.com/fatih/color"
)

type payload struct {
	One   float32
	Two   float64
	Three uint32
}

//Particion is...
type Particion struct {
	partStatus byte
	partType   byte
	partFit    byte
	partStart  int
	partSize   int
	partName   [16]byte
}

//MBR is...
type MBR struct { //22
	Size          uint8
	FechaCreacion [20]byte
	DiskSignature uint8
	DiskFit       byte
	//Particion     [4]Particion
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) {
	writeFile()
	readFile()
}

//writeFile is...
func writeFile() {
	file, err := os.Create("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//primer structsegundostruct

	disco2 := MBR{}
	disco2.Size = 50
	disco2.DiskSignature = 10
	disco2.DiskFit = 'F'

	s1 := &disco2

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	writeNextBytes(file, binario2.Bytes())

}

//writeNextBytes is...
func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}

//readFile is...
func readFile() {

	file, err := os.Open("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := MBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m.Size)
	fmt.Println(m.DiskSignature)
	fmt.Println(string(m.DiskFit))

}

//readNextBytes is...
func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//CalcularSize is ...
func CalcularSize(size int, unit byte) int {
	if unit == 'M' {
		return size * 1024 * 1024
	} else if unit == 'K' {
		return size * 1024
	}
	return 0
}

//VerificarRuta is ...
func VerificarRuta(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//RMDISK is...
func RMDISK(path string) bool {
	if VerificarRuta(path) {
		//Logica para eliminar el disco
		app := "rm"
		cmd := exec.Command(app, path)
		cmd.Output()
	} else {
		ErrorMessage("[RMDISK] -> El disco que desea eliminar no existe")
		return false
	}
	return true
}

//ErrorMessage is...
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
