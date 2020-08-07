package comandos

import (
	"encoding/binary"
	"os"
	"os/exec"

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
type MBR struct {
	mbrSize          int
	mbrFechaCreacion string
	mbrDiskSignature int
	diskFit          byte
	mbrParticion     [4]Particion
}

//EBR is...
type EBR struct {
	partStatus byte
	partFit    byte
	partStart  int
	partSize   int
	partNext   int
	partName   [16]byte
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) {
	var Disco MBR
	Disco.mbrSize = size
	Disco.diskFit = fit
	Disco.mbrDiskSignature = 100
	Disco.mbrFechaCreacion = "06/08/2020"
	//Se inicializan las particiones en el MBR
	for p := 0; p < 4; p++ {
		Disco.mbrParticion[p].partStatus = '0'
		Disco.mbrParticion[p].partType = '0'
		Disco.mbrParticion[p].partFit = '0'
		Disco.mbrParticion[p].partSize = 0
		Disco.mbrParticion[p].partStart = -1
		//strcpy(Disco.mbrParticion[p].part_name, "")
	}
	writeFile(path+name+".disk", CalcularSize(size, unit))
	writeFile(path+name+"Raid.disk", CalcularSize(size, unit))
}

//WriteFile is ...
func writeFile(path string, size int) {
	var Disco MBR
	Disco.mbrSize = 10
	Disco.diskFit = 'f'
	Disco.mbrDiskSignature = 100
	Disco.mbrFechaCreacion = "06/08/2020"
	//Se inicializan las particiones en el MBR
	for p := 0; p < 4; p++ {
		Disco.mbrParticion[p].partStatus = '0'
		Disco.mbrParticion[p].partType = '0'
		Disco.mbrParticion[p].partFit = '0'
		Disco.mbrParticion[p].partSize = 0
		Disco.mbrParticion[p].partStart = -1
		//strcpy(Disco.mbrParticion[p].part_name, "")
	}
	file, _ := os.Create(path)
	for i := 0; i < size; i++ {
		var s uint8 = uint8(i)
		binary.Write(file, binary.LittleEndian, s)
	}
	file.Seek(0, 0)
	binary.Write(file, binary.LittleEndian, Disco)
	file.Close()
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
