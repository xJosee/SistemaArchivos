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

//Particion is...
type Particion struct {
	PartStatus byte
	PartType   byte
	PartFit    byte
	PartStart  uint8
	PartSize   uint8
	PartName   [16]byte
}

//MBR is...
type MBR struct { //22
	Size          uint8
	FechaCreacion [20]byte
	DiskSignature uint8
	DiskFit       byte
	Particion     [4]Particion
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) {
	Disco := MBR{}
	Disco.Size = 50
	Disco.DiskSignature = 10
	Disco.DiskFit = 'F'
	for p := 0; p < 4; p++ {
		Disco.Particion[p].PartStatus = '0'
		Disco.Particion[p].PartType = '0'
		Disco.Particion[p].PartFit = '0'
		Disco.Particion[p].PartSize = 5
		Disco.Particion[p].PartStart = 0
		//strcpy(Disco.Particion[p].part_name, "")
	}
	writeFile(path+name+".disk", CalcularSize(size, unit), Disco)
	readFile(path + name + ".disk")
	writeFile(path+name+"Raid.disk", CalcularSize(size, unit), Disco)
}

//writeFile is...
func writeFile(path string, size int, Disco MBR) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//primer structsegundostruct

	for i := 0; i < size; i++ {
		var ii uint8 = uint8(0)
		err := binary.Write(file, binary.LittleEndian, ii)
		if err != nil {
			fmt.Println("err!", err)
		}
	}
	file.Seek(0, 0)
	s1 := &Disco
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
func readFile(path string) {

	file, err := os.Open(path)
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
