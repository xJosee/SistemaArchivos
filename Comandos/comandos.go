package comandos

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"unsafe"

	"C"

	"github.com/fatih/color"
)
import (
	"io"
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
	PartStatus byte
	PartDit    byte
	PartStart  int32
	PartSize   int32
	PartNext   int32
	PartName   [16]byte
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) {
	//TODO : Ponerle todas las validaciones posibles
	//Creando una instancia del struct MBR que representa al disco
	Disco := MBR{}
	Disco.Size = int32(CalcularSize(size, unit))
	Disco.DiskSignature = 10
	Disco.DiskFit = 'F'
	copy(Disco.FechaCreacion[:], "11/08/2020")
	//Inicializando las particiones del Disco
	for p := 0; p < 4; p++ {
		Disco.Particion[p].PartStatus = '0'
		Disco.Particion[p].PartType = '0'
		Disco.Particion[p].PartFit = '0'
		Disco.Particion[p].PartSize = 0
		Disco.Particion[p].PartStart = -1
		copy(Disco.Particion[p].PartName[:], "")
	}
	//Metodo que escribe el disco(archivo)
	writeFile(path+name+".dsk", CalcularSize(size, unit), Disco)
	//Metodo para leer el struct MBR del Disco(archivo)
	readMBR(path + name + ".dsk")
	//Crea una copia del disco (RAID)
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

//readMBR is...
func readMBR(path string) MBR {

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := MBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytesMBR(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m

}

//readNextBytesMBR is...
func readNextBytesMBR(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//readEBR is...
func readEBR(path string) EBR {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := EBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytesEBR(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return m
}

//readNextBytesEBR is...
func readNextBytesEBR(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

//getFile is...
func getFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	return file
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
		CrearParticionPrimaria(path, CalcularSize(size, unit), name, fit)
		readMBR(path)
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
func CrearParticionPrimaria(path string, size int, name string, fit byte) {
	//TODO : Obtener el file
	//TODO : Optimizar y hacer varias pruebas
	File := getFile(path)
	var mbr MBR
	if VerificarRuta(path) {
		Bandera := false
		num := 0
		mbr = readMBR(path)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartStart == -1 || (mbr.Particion[i].PartStatus == '1' && mbr.Particion[i].PartSize >= int32(size)) {
				Bandera = true
				num = i
				break
			}
		}

		if Bandera {
			//Verificar el espacio libre del disco
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
				if !ParticionExist(path, name) {
					if mbr.DiskFit == 'F' || mbr.DiskFit == 'f' { //FIRST FIT
						mbr.Particion[num].PartType = 'P'
						mbr.Particion[num].PartFit = fit
						//start
						if num == 0 {
							mbr.Particion[num].PartStart = int32(unsafe.Sizeof(mbr))
						} else {
							mbr.Particion[num].PartStart = mbr.Particion[num-1].PartStart + mbr.Particion[num-1].PartSize
						}
						mbr.Particion[num].PartSize = int32(size)
						mbr.Particion[num].PartStatus = '0'
						copy(mbr.Particion[num].PartName[:], name)
						//Se guarda de nuevo el MBR
						File.Seek(0, 0)
						writeFile(path, 15360, mbr)
						//Se guardan los bytes de la particion
						File.Seek(int64(mbr.Particion[num].PartStart), 0)
						for i := 0; i < size; i++ {
							File.Write([]byte{1})
						}
						SuccessMessage("[FDISK] -> Particion Primaria creado correctamente")
					}

				}
			}
		}
	}

}

//ParticionExist is...
func ParticionExist(path string, name string) bool {
	extendida := -1
	if VerificarRuta(path) {
		mbr := readMBR(path)
		for i := 0; i < 4; i++ {
			if fmt.Sprint(mbr.Particion[i].PartName) == name {
				return true
			} else if mbr.Particion[i].PartType == 'E' {
				extendida = i
			}
		}
		if extendida != -1 {
			File := getFile(path)
			File.Seek(int64(mbr.Particion[extendida].PartStart), 0)
			ebr := EBR{}
			ebrBytes := new(bytes.Buffer)
			json.NewEncoder(ebrBytes).Encode(ebr)
			num, _ := File.Read(ebrBytes.Bytes())
			offset, _ := File.Seek(0, io.SeekCurrent)
			for num != 0 && (int32(offset) < (mbr.Particion[extendida].PartSize + mbr.Particion[extendida].PartStart)) {
				if fmt.Sprint(ebr.PartName) == name {
					File.Close()
					return true
				}
				if ebr.PartNext == -1 {
					File.Close()
					return false
				}
			}
		}
	}
	return false
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
