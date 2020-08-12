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
	PartStart  int32
	PartSize   int32
	PartName   [16]byte
}

//MBR is...
type MBR struct { //22
	Size          int32
	FechaCreacion [10]byte
	DiskSignature int32
	DiskFit       byte
	Particion     [4]Particion
}

//EBR is...
type EBR struct { //22
	//TODO : Atributos EBR
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) {
	Disco := MBR{}
	Disco.Size = int32(CalcularSize(size, unit))
	Disco.DiskSignature = 10
	Disco.DiskFit = 'F'
	copy(Disco.FechaCreacion[:], "11/08/2020")
	for p := 0; p < 4; p++ {
		Disco.Particion[p].PartStatus = '0'
		Disco.Particion[p].PartType = '0'
		Disco.Particion[p].PartFit = '0'
		Disco.Particion[p].PartSize = 0
		Disco.Particion[p].PartStart = -1
		copy(Disco.Particion[p].PartName[:], "")
	}
	writeFile(path+name+".dsk", CalcularSize(size, unit), Disco)
	readFile(path + name + ".dsk")
	writeFile(path+name+"_raid.dsk", CalcularSize(size, unit), Disco)
}

//writeFile is...
func writeFile(path string, size int, Disco MBR) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
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
func readFile(path string) MBR {

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

	return m

}

//getFile is...
func getFile(path string) *os.File {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	return file
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
		app := "rm"
		cmd := exec.Command(app, path)
		cmd.Output()
	} else {
		ErrorMessage("[RMDISK] -> El disco que desea eliminar no existe")
		return false
	}
	return true
}

//FDISK is...
func FDISK(size int, unit byte, path string, Type byte, fit byte, delete string, name string, add int) {
	if Type == 'p' {
		CrearParticionPrimaria(path, CalcularSize(size, unit))
	} else if Type == 'e' {
		CrearParticionExtendida()
	} else if Type == 'l' {
		CrearParticionLogica()
	} else if delete != "" {
		EliminarParticion()
	} else if add != 0 {
		AgregarQuitarEspacio()
	}
}

//CrearParticionPrimaria is...
func CrearParticionPrimaria(path string, size int) {
	fmt.Println("Size : ", size)
	buffer := '1'
	var mbr MBR
	fmt.Println(buffer, mbr)
	if VerificarRuta(path) {
		Bandera := false
		num := 0
		fmt.Println(Bandera, num)
		mbr = readFile(path)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartStart == -1 || (mbr.Particion[i].PartStatus == '1' && mbr.Particion[i].PartSize >= int32(size)) {
				Bandera = true
				num = i
				break
			}
		}

		if Bandera {
			//Veroficar el espacio libre del disco
			espacioUsado := 0
			for i := 0; i < 4; i++ {
				if mbr.Particion[i].PartStatus != '1' {
					espacioUsado += int(mbr.Particion[i].PartSize)
				}
			}
			EspacioLibre := (int(mbr.Size) - espacioUsado)

			fmt.Println("EspacioDisponible : ", EspacioLibre)
			fmt.Println("EspacioRequerido : ", size)

			if EspacioLibre >= size {
				if ParticionExist(path) {

				}
			}
		}
	}

}

//ParticionExist is...
func ParticionExist(path string) bool {
	extendida := -1
	if VerificarRuta(path) {
		mbr := readFile(path)
		for i := 0; i < 4; i++ {
			if fmt.Sprint(mbr.Particion[i].PartName) == "Name" {
				return true
			} else if mbr.Particion[i].PartType == 'E' {
				extendida = i
			}
		}
		if extendida != -1 {
			File := getFile(path)
			File.Seek(int64(mbr.Particion[extendida].PartStart), 0)
			//TODO : Seguir con la parte de particion exist
		}
	}

	return true
}

//CrearParticionLogica is...
func CrearParticionLogica() {

}

//CrearParticionExtendida is...
func CrearParticionExtendida() {

}

//EliminarParticion is...
func EliminarParticion() {

}

//AgregarQuitarEspacio is...
func AgregarQuitarEspacio() {

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
