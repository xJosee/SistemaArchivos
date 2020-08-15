package comandos

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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
	FechaCreacion [20]byte
	DiskSignature int32
	DiskFit       byte
	Particion     [4]Particion
}

//EBR is...
type EBR struct { //22
	PartStatus byte
	PartFit    byte
	PartStart  int32
	PartSize   int32
	PartNext   int32
	PartName   [16]byte
}

//MKDISK is...
func MKDISK(size int, fit byte, unit byte, path string, name string) bool {
	//Creando una instancia del struct MBR que representa al disco
	dt := time.Now()
	fecha := dt.Format("01-02-2006 15:04:05")
	Disco := MBR{}
	Disco.Size = int32(CalcularSize(size, unit))
	copy(Disco.FechaCreacion[:], fecha)
	Disco.DiskSignature = int32(rand.Int())
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
	if !VerificarRuta(path + name + ".dsk") {
		//Metodo que escribe el disco(archivo)
		writeFile(path+name+".dsk", CalcularSize(size, unit), Disco)
		//Metodo para leer el struct MBR del Disco(archivo)
		readMBR(path + name + ".dsk")
		//Crea una copia del disco (RAID)
		writeFile(path+name+"_raid.dsk", CalcularSize(size, unit), Disco)

		return true

	}

	return false
}

//writeFile is...
func writeFile(path string, size int, Disco MBR) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < size; i++ {
		var numParticion int8 = int8(0)
		err := binary.Write(file, binary.LittleEndian, numParticion)
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

//reWriteMBR is...
func reWriteMBR(file *os.File, Disco MBR) {
	file.Seek(0, 0)
	s1 := &Disco
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	writeNextBytes(file, binario2.Bytes())
}

//reWriteEBR is...
func reWriteEBR(file *os.File, Disco EBR) {
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
	/*for i := 0; i < 4; i++ {
		fmt.Println("Tipo", m.Particion[i].PartType)
		fmt.Println("Size", m.Particion[i].PartSize)
	}*/
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
	file, err := os.OpenFile(path, os.O_RDWR, 0755)
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
func FDISK(size int, unit byte, path string, Type byte, fit byte, delete string, name string, add int) bool {
	if Type == 'p' {
		if CrearParticionPrimaria(path, CalcularSize(size, unit), name, fit) {
			readMBR(path)
			return true
		}
	} else if Type == 'e' {
		CrearParticionExtendida(path, CalcularSize(size, unit), name, fit)
		readMBR(path)
	} else if Type == 'l' {
		CrearParticionLogica()
	} else if delete != "" {
		EliminarParticion()
	} else if add != 0 {
		AgregarQuitarEspacio()
	}

	return false
}

//CrearParticionPrimaria is...
func CrearParticionPrimaria(path string, size int, name string, fit byte) bool {
	File := getFile(path)
	var mbr MBR
	if VerificarRuta(path) {
		Bandera := false
		numParticion := 0
		File.Seek(0, 0)
		mbr = readMBR(path)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartStart == -1 || (mbr.Particion[i].PartStatus == '1' && mbr.Particion[i].PartSize >= int32(size)) {
				Bandera = true
				numParticion = i
				break
			}
		}
		//Bandera -> Indica si tiene espacio para crear la particion
		if Bandera {
			//Verificar el espacio libre del disco
			espacioUsado := 0
			for i := 0; i < 4; i++ {
				if mbr.Particion[i].PartStatus != '1' {
					espacioUsado += int(mbr.Particion[i].PartSize)
				}
			}
			EspacioLibre := (int(mbr.Size) - espacioUsado)
			//fmt.Println("EspacioDisponible : ", EspacioLibre)
			//fmt.Println("EspacioRequerido : ", size)
			if EspacioLibre >= size {
				if !ParticionExist(path, name) {
					if mbr.DiskFit == 'F' || mbr.DiskFit == 'f' { //FIRST FIT
						mbr.Particion[numParticion].PartType = 'P'
						mbr.Particion[numParticion].PartFit = fit
						//start
						if numParticion == 0 {
							mbr.Particion[numParticion].PartStart = int32(unsafe.Sizeof(mbr))
						} else {
							mbr.Particion[numParticion].PartStart = mbr.Particion[numParticion-1].PartStart + mbr.Particion[numParticion-1].PartSize
						}
						mbr.Particion[numParticion].PartSize = int32(size)
						mbr.Particion[numParticion].PartStatus = '0'
						copy(mbr.Particion[numParticion].PartName[:], name)
						//Se guarda de nuevo el MBR
						reWriteMBR(File, mbr)
						//Se guardan los bytes de la particion
						File.Seek(int64(mbr.Particion[numParticion].PartStart), 0)
						for i := 0; i < size; i++ {
							File.Write([]byte{1})
						}
						SuccessMessage("[FDISK] -> Particion Primaria creada correctamente")
						return true
					}
				} else {
					ErrorMessage("[FDISK] -> Ya existe una particion con el mismo nombre")
				}
			} else {
				ErrorMessage("[FDISK] -> La particion a crear excede el size del disco")
			}
		} else {
			ErrorMessage("[FDISK] -> El disco ya cuenta con 4 particiones")
		}
		File.Close()
	} else {
		ErrorMessage("[FDISK] -> El disco no existe")
	}

	return false

}

//ParticionExist is...
func ParticionExist(path string, name string) bool {
	extendida := -1
	if VerificarRuta(path) {
		mbr := readMBR(path)
		for i := 0; i < 4; i++ {

			nameParticionString := string(mbr.Particion[i].PartName[:])
			//TODO : Ver bien la comparacion de nombres
			if strings.Compare(nameParticionString, name) == 1 {
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
			numParticion, _ := File.Read(ebrBytes.Bytes())
			offset, _ := File.Seek(0, io.SeekCurrent)

			for numParticion != 0 && (int32(offset) < (mbr.Particion[extendida].PartSize + mbr.Particion[extendida].PartStart)) {
				numParticion, _ = File.Read(ebrBytes.Bytes())
				offset, _ = File.Seek(0, io.SeekCurrent)
				nameParticionString := string(ebr.PartName[:])
				if nameParticionString == name {
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
	//TODO : Crear particion logica
}

//CrearParticionExtendida is...
func CrearParticionExtendida(path string, size int, name string, fit byte) {

	File := getFile(path)
	var mbr MBR
	if VerificarRuta(path) {
		var flagParticion bool = false
		var flagExtendida bool = false
		var numParticion int = 0
		File.Seek(0, 0)
		mbr = readMBR(path)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartType == 'E' || mbr.Particion[i].PartType == 'e' {
				flagExtendida = true
				break
			}
		}
		if !flagExtendida {
			//Verificar si existe una particion disponible
			for i := 0; i < 4; i++ {
				if mbr.Particion[i].PartStart == -1 || (mbr.Particion[i].PartStatus == '1' && mbr.Particion[i].PartSize >= int32(size)) {
					flagParticion = true
					numParticion = i
					break
				}
			}
			if flagParticion {
				//Verificar el espacio libre del disco
				var espacioUsado int = 0
				for i := 0; i < 4; i++ {
					if mbr.Particion[i].PartStatus != '1' {
						espacioUsado += int(mbr.Particion[i].PartSize)
					}
				}
				EspacioDisponible := mbr.Size - int32(espacioUsado)
				if EspacioDisponible >= int32(size) {
					if !ParticionExist(path, name) {
						if mbr.DiskFit == 'F' || mbr.DiskFit == 'f' {
							mbr.Particion[numParticion].PartType = 'E'
							mbr.Particion[numParticion].PartFit = fit
							//start
							if numParticion == 0 {
								mbr.Particion[numParticion].PartStart = int32(unsafe.Sizeof(mbr))
							} else {
								mbr.Particion[numParticion].PartStart = mbr.Particion[numParticion-1].PartStart + mbr.Particion[numParticion-1].PartSize
							}
							mbr.Particion[numParticion].PartSize = int32(size)
							mbr.Particion[numParticion].PartStatus = '0'
							copy(mbr.Particion[numParticion].PartName[:], name)
							//Se guarda de nuevo el MBR
							reWriteMBR(File, mbr)
							//Se guardan los bytes de la particion
							File.Seek(int64(mbr.Particion[numParticion].PartStart), 0)

							var EB EBR
							EB.PartFit = fit
							EB.PartStatus = '0'
							EB.PartStart = mbr.Particion[numParticion].PartStart
							EB.PartSize = 0
							EB.PartNext = -1
							copy(EB.PartName[:], "")

							reWriteEBR(File, EB)

							for i := 0; i < size-int(unsafe.Sizeof(EB)); i++ {
								File.Write([]byte{1})
							}
							SuccessMessage("[FDISK] -> Particion extendida creada correctamente")
						}
					} else {
						ErrorMessage("[FDISK] -> Ya existe una particion con ese nombre")
					}
				} else {
					ErrorMessage("[FDISK] -> La particion a crear es mayor al espacio libre del disco")
				}
			} else {
				ErrorMessage("[FDISK] -> El disco ya cuenta con 4 particiones")
			}
		} else {
			ErrorMessage("[FDISK] -> El disco ya cuenta con una particion extendida")
		}
		File.Close()
	} else {
		ErrorMessage("[FDISK] -> No existe el disco")
	}
}

//EliminarParticion is...
func EliminarParticion() {
	//TODO : Eliminar Particion
}

//AgregarQuitarEspacio is...
func AgregarQuitarEspacio() {
	//TODO : Agregar o Quitar espacio
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
