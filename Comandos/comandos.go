package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"

	Estructuras "../Estructuras"
	"github.com/fatih/color"
)

/*
 *   S T R U C T S
 */

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
type MBR struct {
	Size          int32
	FechaCreacion [19]byte
	DiskSignature int32
	DiskFit       byte
	Particion     [4]Particion
}

//EBR is...
type EBR struct {
	PartStatus byte
	PartFit    byte
	PartStart  int32
	PartSize   int32
	PartNext   int32
	PartName   [16]byte
}

/*
 *	S T R U C T   F A S E   2
 */

//SB is...
type SB struct {
	NombreHD                 [16]byte
	ArbolVirtualCount        int32
	DetalleDirectorioCount   int32
	InodosCount              int32
	BloquesCount             int32
	ArbolVirtualFree         int32
	DetalleDirectorioFree    int32
	InodosFree               int32
	BloquesFree              int32
	DateCreacion             [19]byte
	DateUltimoMontaje        [19]byte
	MontajesCount            int32
	StartBmArbolDirectorio   int32
	StartArbolDirectorio     int32
	StartBmDetalleDirectorio int32
	StartDetalleDirectorio   int32
	StartBmInodos            int32
	StartInodos              int32
	StartBmBloques           int32
	StartBloques             int32
	StartLog                 int32 //Bitacora.
	SizeStructAvd            int32 // = sizeof(arbolVirtual);
	SizeStructDd             int32 // sizeof(detalleDirectorio);
	SizeStructInodo          int32 // sizeof(InodoArchivo);
	SizeStructBloque         int32 // sizeof(bloqueDatos);
	FirstFreeAvd             int32
	FirstFreeDd              int32
	FirstFreeInodo           int32
	FirstFreeBloque          int32
	FirstFreeBitacora        int32
	MagicNum                 int32 //= 201807431;
}

//Bloque is...
type Bloque struct {
	Texto [25]byte
}

//Bitacora is...
type Bitacora struct {
	TipoOp    [20]byte
	Tipo      int32
	Nombre    [20]byte
	Contenido [20]byte
	Fecha     [10]byte
	Size      int32
}

//TablaInodo is...
type TablaInodo struct {
	ICountInodo            int32
	ISizeArchivo           int32
	ICountBloquesAsignados int32
	IArrayBloques          [4]int32
	IApIndirecto           int32
	IIDProper              int32
}

//DetalleDirectorio is...
type DetalleDirectorio struct {
	DDArrayFiles          [5]FileStruct // Los archivos que puede tener
	DDApDetalleDirectorio int32         // Apunta a la copia del DD por si ya no caben mas files
}

//FileStruct is...
type FileStruct struct {
	DDFileNombre           [16]byte
	DDFileApInodo          int32
	DDFileDateCreacion     [20]byte
	DDFileDateModificacion [20]byte
}

//Arbol is...
type Arbol struct {
	AVDFechaCreacion    [10]byte
	AVDNombreDirectorio [16]byte
	Subirectorios       [6]int32
	VirtualDirectorio   int32
	DetalleDirectorio   int32
	AVDProper           int32
}

//Usuario is...
type Usuario struct {
	IDUser   int32
	IDGrupo  int32
	UserName [12]byte
	PassWord [12]byte
	Group    [12]byte
}

/*
 *  L I S T A   P A R T I C I O N E S   M O N T A D A S
 */

var listaParticiones = Estructuras.Lista{
	Contador: 0,
	Primero:  nil,
}

/*
 *  V A R I A B L E S   P A R A   E L   M A N E J O   D E   LA  S E S I O N
 */

var userLoggeado Usuario
var isLogged bool = false

/*
 *  C O M A N D O S
 */

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
		os.MkdirAll(path, os.ModePerm)
		writeFile(path+name+".dsk", CalcularSize(size, unit), Disco)
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
	file.Seek(0, 0)
	var cero byte = '0'
	s1 := &cero
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())

	file.Seek(int64(size-1), 0)
	file.Write(binario2.Bytes())

	file.Seek(0, 0)
	//Meto el MBR
	s2 := &Disco
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s2)
	file.Write(binario.Bytes())

}

//reWriteMBR is...
func reWriteMBR(file *os.File, Disco MBR) {
	file.Seek(0, 0)
	s1 := &Disco
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//reWriteEBR is...
func reWriteEBR(file *os.File, Disco EBR, seek int64) {
	file.Seek(seek, 0)
	s1 := &Disco
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

func reWriteSuperBloque(file *os.File, SuperB SB, seek int64) {
	file.Seek(seek, 0)
	s1 := &SuperB
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

/*
 *  READS DE LOS STRUCTS EN LOS FILES
 */

//readMBR is...
func readMBR(file *os.File) MBR {
	m := MBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m

}

//readEBR is...
func readEBR(file *os.File, seek int64) EBR {
	file.Seek(seek, 0)
	m := EBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readByte is...
func readByte(file *os.File, seek int64) byte {
	file.Seek(seek, 0)
	var m byte
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readSuperBloque is...
func readSuperBloque(file *os.File, seek int64) SB {
	file.Seek(seek, 0)
	m := SB{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readArbolVirtualDirectorio is...
func readArbolVirtualDirectorio(file *os.File, seek int64) Arbol {
	file.Seek(seek, 0)
	m := Arbol{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readInodo is...
func readInodo(file *os.File, seek int64) TablaInodo {
	file.Seek(seek, 0)
	m := TablaInodo{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readBloque is...
func readBloque(file *os.File, seek int64) Bloque {
	file.Seek(seek, 0)
	m := Bloque{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readDetalleDirectorio
func readDetalleDirectorio(file *os.File, seek int64) DetalleDirectorio {
	file.Seek(seek, 0)
	m := DetalleDirectorio{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readBitacora is...
func readBitacora(file *os.File, seek int64) Bitacora {
	file.Seek(seek, 0)
	m := Bitacora{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

//readNextBytesEBR is...
func readNextBytes(file *os.File, number int) []byte {
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
	if unit == 'M' || unit == 'm' {
		return size * 1024 * 1024
	} else if unit == 'K' || unit == 'k' {
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

//ParticionExist is...
func ParticionExist(path string, name string) bool {
	//extendida := -1
	if VerificarRuta(path) {
		File := getFile(path)
		mbr := readMBR(File)
		extendida := -1
		for i := 0; i < 4; i++ {

			var nameByte [16]byte
			copy(nameByte[:], name)
			//fmt.Println("ParticionExisit", string(nameByte[:]), string(mbr.Particion[i].PartName[:]), bytes.Compare(nameByte[:], mbr.Particion[i].PartName[:]))
			if bytes.Compare(nameByte[:], mbr.Particion[i].PartName[:]) == 0 {
				fmt.Println("Si son iguales")
				File.Close()
				return true
			} else if mbr.Particion[i].PartType == 'E' {
				extendida = i
			}
		}

		if extendida != -1 {

			ebr := EBR{
				PartNext: -2,
			}

			for ebr.PartNext != -1 && (ebr.PartNext < (mbr.Particion[extendida].PartSize + mbr.Particion[extendida].PartStart)) {
				if ebr.PartNext == -2 {
					ebr = readEBR(File, int64(mbr.Particion[extendida].PartStart))
				} else {
					ebr = readEBR(File, int64(ebr.PartNext))
				}
				var nameByte [16]byte
				copy(nameByte[:], name)

				if bytes.Compare(nameByte[:], ebr.PartName[:]) == 0 {
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

//EliminarParticion is...
func EliminarParticion(path string, name string, delete string) {
	//TODO : Eliminar Particion Logica

	if VerificarRuta(path) {
		File := getFile(path)

		var mount bool = listaParticiones.BuscarNodo(path, name)

		if !mount {

			var masterboot MBR
			File.Seek(0, 0)
			masterboot = readMBR(File)
			var index int = -1
			var flagExtendida bool = false
			//int index_Extendida = -1

			for i := 0; i < 4; i++ {
				var nameByte [16]byte
				copy(nameByte[:], name)

				if bytes.Compare(nameByte[:], masterboot.Particion[i].PartName[:]) == 0 {
					index = i
					if masterboot.Particion[i].PartType == 'E' {
						flagExtendida = true
					}
					break
				} else if masterboot.Particion[i].PartType == 'E' {
					//index_Extendida = i
				}
			}

			/*fmt.Println("[FDISK] -> Seguro que desea eliminar la particion? (S/N)")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()*/

			if index != -1 {
				if !flagExtendida {
					if strings.ToLower(delete) == "fast" {

						masterboot.Particion[index].PartStatus = '1'
						masterboot.Particion[index].PartType = '0'
						masterboot.Particion[index].PartFit = '0'
						masterboot.Particion[index].PartSize = 0
						masterboot.Particion[index].PartStart = -1
						copy(masterboot.Particion[index].PartName[:], "")

						reWriteMBR(File, masterboot)
						SuccessMessage("[FDISK] -> Particion eliminada correctamente")
					} else {
						masterboot.Particion[index].PartStatus = '1'
						masterboot.Particion[index].PartType = '0'
						masterboot.Particion[index].PartFit = '0'
						masterboot.Particion[index].PartSize = 0
						masterboot.Particion[index].PartStart = -1
						copy(masterboot.Particion[index].PartName[:], "")
						reWriteMBR(File, masterboot)
						File.Seek(int64(masterboot.Particion[index].PartStart), 0)
						for i := 0; i < int(masterboot.Particion[index].PartSize); i++ {
							File.Write([]byte{0})
						}
						SuccessMessage("[FDISK] -> Particion eliminada correctamente")
					}
				} else {
					if strings.ToLower(delete) == "fast" {

						masterboot.Particion[index].PartStatus = '1'
						masterboot.Particion[index].PartType = '0'
						masterboot.Particion[index].PartFit = '0'
						masterboot.Particion[index].PartSize = 0
						masterboot.Particion[index].PartStart = -1
						copy(masterboot.Particion[index].PartName[:], "")
						reWriteMBR(File, masterboot)
						SuccessMessage("[FDISK] -> Particion eliminada correctamente")

					} else if strings.ToLower(delete) == "full" {
						masterboot.Particion[index].PartStatus = '1'
						masterboot.Particion[index].PartType = '0'
						masterboot.Particion[index].PartFit = '0'
						masterboot.Particion[index].PartSize = 0
						masterboot.Particion[index].PartStart = -1
						copy(masterboot.Particion[index].PartName[:], "")
						reWriteMBR(File, masterboot)
						File.Seek(int64(masterboot.Particion[index].PartStart), 0)
						for i := 0; i < int(masterboot.Particion[index].PartSize); i++ {
							File.Write([]byte{0})
						}
						SuccessMessage("[FDISK] -> Particion eliminada correctamente")
					}

				}
			}

		} else {
			ErrorMessage("[FDISK] -> No se puede eliminar una particion montada")
		}

	}
}

//AgregarQuitarEspacio is...
func AgregarQuitarEspacio(path string, name string, add int, unit byte) {
	var tipo string = ""

	if add > 0 {
		tipo = "add"
	}

	if tipo != "add" {
		add = add * -1
	}

	var size int = CalcularSize(add, unit)

	if VerificarRuta(path) {
		File := getFile(path)
		var mount bool = listaParticiones.BuscarNodo(path, name)

		if !mount {
			var masterboot MBR
			File.Seek(0, 0)
			masterboot = readMBR(File)
			var index int = -1
			var indexExtendida int = -1
			var flagExtendida bool = false

			for i := 0; i < 4; i++ {
				var nameByte [16]byte
				copy(nameByte[:], name)

				if bytes.Compare(masterboot.Particion[i].PartName[:], nameByte[:]) == 0 {
					index = i
					if masterboot.Particion[i].PartType == 'E' {
						flagExtendida = true
					}
					break
				} else if masterboot.Particion[i].PartType == 'E' {
					indexExtendida = i
				}
			}

			if index != -1 {
				if !flagExtendida {
					if tipo == "add" {
						if index != 3 {
							var p1 int = int(masterboot.Particion[index].PartStart + masterboot.Particion[index].PartSize)
							var p2 int = int(masterboot.Particion[index+1].PartStart)
							if (p2 - p1) != 0 {
								var fragmentacion int = p2 - p1
								if fragmentacion >= size {
									masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
									File.Seek(0, 0)
									reWriteMBR(File, masterboot)
									SuccessMessage("[FDISK] -> Espacio agregado correctamente")

								} else {
									ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
								}
							} else {
								if masterboot.Particion[index+1].PartStatus != '1' {
									if masterboot.Particion[index+1].PartSize >= int32(size) {
										masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
										masterboot.Particion[index+1].PartSize = (masterboot.Particion[index+1].PartSize - int32(size))
										masterboot.Particion[index+1].PartSize = masterboot.Particion[index+1].PartStart + int32(size)
										File.Seek(0, 0)
										reWriteMBR(File, masterboot)
										SuccessMessage("[FDISK] -> Espacio agregado correctamente")
									} else {
										ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
									}
								}
							}
						} else {
							var p int = int(masterboot.Particion[index].PartStart + masterboot.Particion[index].PartSize)
							var total int = int(masterboot.Size + int32(unsafe.Sizeof(masterboot)))
							if (total - p) != 0 {
								var fragmentacion int = total - p
								if fragmentacion >= size {
									masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
									File.Seek(0, 0)
									reWriteMBR(File, masterboot)
									SuccessMessage("[FDISK] -> Espacio agregado correctamente")
								} else {
									ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
								}
							} else {
								ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
							}
						}
					} else {
						if int32(size) >= masterboot.Particion[index].PartSize {
							ErrorMessage("[FDISK] -> No se puede disminuir esa cantidad de espacio")
						} else {
							masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize - int32(size)
							File.Seek(0, 0)
							reWriteMBR(File, masterboot)
							SuccessMessage("[FDISK] -> Espacio reducido correctamente")
						}
					}

				} else {
					if tipo == "add" {
						if index != 3 {
							var p1 int = int(masterboot.Particion[index].PartStart + masterboot.Particion[index].PartSize)
							var p2 int = int(masterboot.Particion[index+1].PartStart)
							if (p2 - p1) != 0 {
								var fragmentacion int = p2 - p1
								if fragmentacion >= size {
									masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
									File.Seek(0, 0)
									reWriteMBR(File, masterboot)
									SuccessMessage("[FDISK] -> Espacio agregado correctamente")

								} else {
									ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
								}
							} else {
								if masterboot.Particion[index+1].PartStatus != '1' {
									if masterboot.Particion[index+1].PartSize >= int32(size) {
										masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
										masterboot.Particion[index+1].PartSize = (masterboot.Particion[index+1].PartSize - int32(size))
										masterboot.Particion[index+1].PartStart = masterboot.Particion[index+1].PartStart + int32(size)
										File.Seek(0, 0)
										reWriteMBR(File, masterboot)
										SuccessMessage("[FDISK] -> Espacio agregado correctamente")
									} else {
										ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
									}
								}
							}
						} else {
							var p int = int(masterboot.Particion[index].PartStart + masterboot.Particion[index].PartSize)
							var total int = int(masterboot.Size + int32(unsafe.Sizeof(masterboot)))
							if (total - p) != 0 {
								var fragmentacion int = total - p

								if fragmentacion >= size {
									masterboot.Particion[index].PartSize = masterboot.Particion[index].PartSize + int32(size)
									File.Seek(0, 0)
									reWriteMBR(File, masterboot)
									SuccessMessage("[FDISK] -> Espacio agregado correctamente")
								} else {
									ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
								}
							} else {
								ErrorMessage("[FDISK] -> No cuenta con suficiente espacio")
							}
						}
					} else {
						if int32(size) >= masterboot.Particion[indexExtendida].PartSize {
							ErrorMessage("[FDISK] -> No es posible reducir esa cantidad de espacio")
						} else {
							var extendedBoot EBR
							extendedBoot = readEBR(File, int64(masterboot.Particion[indexExtendida].PartStart))

							for (extendedBoot.PartNext != -1) && (extendedBoot.PartNext < (masterboot.Particion[indexExtendida].PartSize + masterboot.Particion[indexExtendida].PartStart)) {
								extendedBoot = readEBR(File, int64(extendedBoot.PartNext))
							}
							var ultimaLogica int = int(extendedBoot.PartStart + extendedBoot.PartSize)
							var aux int = int((masterboot.Particion[indexExtendida].PartStart + masterboot.Particion[indexExtendida].PartSize) - int32(size))
							if aux > ultimaLogica { //No toca ninguna logica
								masterboot.Particion[indexExtendida].PartSize = masterboot.Particion[indexExtendida].PartSize - int32(size)
								File.Seek(0, 0)
								reWriteMBR(File, masterboot)
								SuccessMessage("[FDISK] -> Espacio reducido correctamente")
							} else {
								ErrorMessage("[FDISk] -> No se puede reducir esa cantidad de espacio")
							}
						}
					}
				}
			} else {
				if indexExtendida != -1 {
					var logica int = ParticionLogicaExist(path, name)
					if logica != -1 {
						if tipo == "add" {

							var extendedBoot EBR
							extendedBoot = readEBR(File, int64(logica))
							_ = extendedBoot

						} else {

							var extendedBoot EBR
							extendedBoot = readEBR(File, int64(logica))

							if int32(size) >= extendedBoot.PartSize {
								ErrorMessage("[FDISk] -> No se puede reducir esa cantidad de espacio")
							} else {
								extendedBoot.PartSize = extendedBoot.PartSize - int32(size)
								reWriteEBR(File, extendedBoot, int64(logica))
								SuccessMessage("[FDISK] -> Espacio reducido correctamente")
							}
						}
					} else {
						ErrorMessage("[FDISK] -> No se encuentra la particion")
					}
				} else {
					ErrorMessage("[FDISK] -> No se encuentra la particion")
				}
			}

		} else {
			ErrorMessage("[FDISK] -> No se puede editar una particion montada")
		}
	}
}

//MOUNT is...
func MOUNT(path string, name string) {

	indexP := ParticionExtendidaExist(path, name)

	if indexP != -1 {
		File := getFile(path)

		if VerificarRuta(path) {
			var masterboot MBR
			File.Seek(0, 0)

			masterboot = readMBR(File)

			masterboot.Particion[indexP].PartStatus = '2'

			reWriteMBR(File, masterboot)
			File.Close()

			letra := listaParticiones.BuscarLetra(path, name)

			if letra == -1 {
				ErrorMessage("[MOUNT] -> La particion ya se encuentra montada")
			} else {
				num := listaParticiones.BuscarNumero(path, name)
				auxLetra := byte(letra)
				id := "vd"
				id += string(auxLetra) + string(num)

				n := Estructuras.Nodo{
					Direccion: path,
					Nombre:    name,
					Letra:     auxLetra,
					Num:       num,
					Siguiente: nil,
					PartStart: int(masterboot.Particion[indexP].PartStart),
					PartSize:  int(masterboot.Particion[indexP].PartSize),
				}

				listaParticiones.Insertar(&n)
				listaParticiones.Listar()

				SuccessMessage("[MOUNT] -> Particion montada correctamente")

			}
		} else {
			ErrorMessage("[MOUNT] -> El disco no existe")
		}
	} else {
		indexP := ParticionLogicaExist(path, name)
		if indexP != -1 {

			File := getFile(path)

			if VerificarRuta(path) {

				var extendedBoot EBR
				extendedBoot = readEBR(File, int64(indexP))
				extendedBoot.PartStatus = '2'
				reWriteEBR(File, extendedBoot, int64(indexP))
				File.Close()

				letra := listaParticiones.BuscarLetra(path, name)

				if letra == -1 {
					ErrorMessage("[MOUNT:ERROR] : La particion ya se encuentra montada")
				} else {
					num := listaParticiones.BuscarNumero(path, name)
					auxLetra := byte(letra)
					id := "vd"
					id += string(auxLetra) + string(num)

					n := Estructuras.Nodo{
						Direccion: path,
						Nombre:    name,
						Letra:     auxLetra,
						Num:       num,
						Siguiente: nil,
					}

					listaParticiones.Insertar(&n)
					listaParticiones.Listar()
					SuccessMessage("[MOUNT] -> Particion montada correctamente")

				}
			} else {
				ErrorMessage("[MOUNT] -> El disco no existe")
			}
		} else {
			ErrorMessage("[MOUNT] -> La particion no se encuentra")
		}
	}

}

//UNMOUNT is...
func UNMOUNT(id string) bool {
	var eliminado int = listaParticiones.EliminarNodo(id)
	if eliminado == 1 {
		SuccessMessage("[UNMOUNT] -> Particion desmontada correctamente")
		listaParticiones.Listar()
		return true
	}
	ErrorMessage("[UNMOUNT] ->  La particion a desmontar no se encuentra")
	return false
}

/*
 *  C R E A R   P A R T I C I O N E S
 */

//CrearParticionPrimaria is...
func CrearParticionPrimaria(path string, size int, name string, fit byte) bool {
	File := getFile(path)
	var mbr MBR
	if VerificarRuta(path) {
		Bandera := false
		numParticion := 0
		File.Seek(0, 0)
		mbr = readMBR(File)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartStart == -1 || (mbr.Particion[i].PartStatus == '1' && mbr.Particion[i].PartSize >= int32(size)) {
				Bandera = true
				numParticion = i
				break
			}
		}
		//Bandera -> Indica si tiene espacio para crear la particion
		if Bandera {
			//Verificar el LIBRE del disco
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
							var x byte = 1
							var start bytes.Buffer
							binary.Write(&start, binary.BigEndian, x)
							File.Write(start.Bytes())
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

//CrearParticionLogica is...
func CrearParticionLogica(path string, name string, size int, fit byte) {
	var mbr MBR
	if VerificarRuta(path) {
		File := getFile(path)
		var numExtendida int = -1
		File.Seek(0, 0)
		mbr = readMBR(File)
		for i := 0; i < 4; i++ {
			if mbr.Particion[i].PartType == 'E' || mbr.Particion[i].PartType == 'e' {
				numExtendida = i
				break
			}
		}

		if !ParticionExist(path, name) {
			if numExtendida != -1 {
				var EB EBR
				cont := mbr.Particion[numExtendida].PartStart
				EB = readEBR(File, int64(cont))

				if mbr.Particion[numExtendida].PartSize < int32(size) {
					ErrorMessage("[FDISK] -> La particion logica que desea crear excede en size a la extendida")
				} else {
					for (EB.PartNext != -1) && (EB.PartNext < (mbr.Particion[numExtendida].PartSize + mbr.Particion[numExtendida].PartStart)) {
						EB = readEBR(File, int64(EB.PartNext))
					}

					espacioNecesario := EB.PartStart + EB.PartSize + int32(size)

					if espacioNecesario <= (mbr.Particion[numExtendida].PartSize + mbr.Particion[numExtendida].PartStart) {

						// Escribimos el EBR anterior
						EB.PartNext = EB.PartStart + int32(size) + int32(unsafe.Sizeof(EB))
						EB.PartSize = int32(size)
						copy(EB.PartName[:], name)
						reWriteEBR(File, EB, int64(EB.PartStart))

						//Escribimos el nuevo EBR
						EB.PartStatus = 0
						EB.PartFit = fit
						EB.PartStart = EB.PartNext
						EB.PartSize = 0
						EB.PartNext = -1

						reWriteEBR(File, EB, int64(EB.PartStart))
						SuccessMessage("[FDISK] -> Particion logica creada correctamente")

					} else {
						ErrorMessage("[FDISK] -> La particion logica es mas grande que la extendida")
					}
				}

			} else {
				ErrorMessage("[FDISK] -> Para crear una particion logica debe existir una extendida")
			}
		} else {
			ErrorMessage("[FDISK] -> Ya existe una particion con ese nombre")
		}

		File.Close()
	} else {
		ErrorMessage("[FDISK] -> No existe el disco")
	}
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
		mbr = readMBR(File)
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
				//Verificar el LIBRE del disco
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

							var EB EBR
							EB.PartFit = fit
							EB.PartStatus = '0'
							EB.PartStart = mbr.Particion[numParticion].PartStart
							EB.PartSize = 0
							EB.PartNext = -1
							copy(EB.PartName[:], "")

							reWriteEBR(File, EB, int64(mbr.Particion[numParticion].PartStart))

							for i := 0; i < size-int(unsafe.Sizeof(EB)); i++ {
								var x byte = 1
								var start bytes.Buffer
								binary.Write(&start, binary.BigEndian, x)
								File.Write(start.Bytes())
							}
							SuccessMessage("[FDISK] -> Particion extendida creada correctamente")
						}
					} else {
						ErrorMessage("[FDISK] -> Ya existe una particion con ese nombre")
					}
				} else {
					ErrorMessage("[FDISK] -> La particion a crear es mayor al LIBRE del disco")
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

/*
 *  V E R I F I C A C I O N E S
 */

//ParticionExtendidaExist is...
func ParticionExtendidaExist(path string, name string) int {
	if VerificarRuta(path) {
		File := getFile(path)
		var masterboot MBR
		File.Seek(0, 0)
		masterboot = readMBR(File)
		for i := 0; i < 4; i++ {
			if masterboot.Particion[i].PartStatus != '1' {
				var nameByte [16]byte
				copy(nameByte[:], name)
				if bytes.Compare(nameByte[:], masterboot.Particion[i].PartName[:]) == 0 {
					return i
				}
			}
		}

	}
	return -1
}

//ParticionLogicaExist is...
func ParticionLogicaExist(path string, name string) int {
	if VerificarRuta(path) {
		File := getFile(path)
		var extendida int = -1
		var masterboot MBR
		File.Seek(0, 0)
		masterboot = readMBR(File)

		for i := 0; i < 4; i++ {
			if masterboot.Particion[i].PartType == 'E' {
				extendida = i
				break
			}
		}
		if extendida != -1 {

			ebr := EBR{}
			ebr = readEBR(File, int64(masterboot.Particion[extendida].PartStart))

			for ebr.PartNext != -1 && (ebr.PartNext < masterboot.Particion[extendida].PartStart+masterboot.Particion[extendida].PartSize) {
				var nameByte [16]byte
				copy(nameByte[:], name)
				if bytes.Compare(ebr.PartName[:], nameByte[:]) == 0 {
					return int((ebr.PartNext - int32(unsafe.Sizeof(ebr))))
				}
				if ebr.PartNext == -1 {

				} else {
					ebr = readEBR(File, int64(ebr.PartNext))
				}
			}
		}
		File.Close()
	}
	return -1
}

/*
 *	R E P O R T E S
 */

//**** FASE 1 ****/

//ReporteEBR is...
func ReporteEBR(path string) {

	if VerificarRuta(path) {

		File := getFile(path)
		os.Create("Reportes/graficaEBR.dot")
		graphDot := getFile("Reportes/graficaEBR.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")
		fmt.Fprintf(graphDot, "tbl[\nlabel=<\n")
		fmt.Fprintf(graphDot, "<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n")
		fmt.Fprintf(graphDot, " <tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>MBR</font></b></td></tr>")
		fmt.Fprintf(graphDot, "<tr>  <td width='150'> <b>Nombre</b> </td> <td width='150'> <b>Valor</b> </td>  </tr>\n")

		var MB MBR
		File.Seek(0, 0)
		MB = readMBR(File)

		var tamano int = int(MB.Size)

		fmt.Fprintf(graphDot, "<tr>  <td>Size</td><td>%d</td>  </tr>\n", tamano)

		//Obteniendo la fehca
		dt := time.Now()
		fecha := dt.Format("01-02-2006 15:04:05")

		fmt.Fprintf(graphDot, "<tr>  <td>Fecha</td> <td>%s</td>  </tr>\n", string(fecha))
		fmt.Fprintf(graphDot, "<tr>  <td>Signature</td> <td>%d</td>  </tr>\n", MB.DiskSignature)
		fmt.Fprintf(graphDot, "<tr>  <td>Fit</td> <td>%c</td>  </tr>\n", MB.DiskFit)

		var posExtendida int = -1

		for i := 0; i < 4; i++ {

			if MB.Particion[i].PartStart != -1 && MB.Particion[i].PartStatus != '1' {
				if MB.Particion[i].PartType == 'E' {
					posExtendida = i
				}
				var status string
				if MB.Particion[i].PartStatus == '0' {
					status = "0"
				} else if MB.Particion[i].PartStatus == '2' {
					status = "2"
				} else if MB.Particion[i].PartStatus == '1' {
					status = "1"
				}

				fmt.Fprintf(graphDot, "<tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>Particion%d</font></b></td></tr>\n", (i + 1))
				fmt.Fprintf(graphDot, "<tr>  <td>Status</td> <td>%s</td>  </tr>\n", status)
				fmt.Fprintf(graphDot, "<tr>  <td>Type</td> <td>%c</td>  </tr>\n", MB.Particion[i].PartType)
				fmt.Fprintf(graphDot, "<tr>  <td>Fit</td> <td>%c</td>  </tr>\n", MB.Particion[i].PartFit)
				fmt.Fprintf(graphDot, "<tr>  <td>Start</td> <td>%d</td>  </tr>\n", MB.Particion[i].PartStart)
				fmt.Fprintf(graphDot, "<tr>  <td>Size</td> <td>%d</td>  </tr>\n", MB.Particion[i].PartSize)
				PartName := string(MB.Particion[i].PartName[:])
				PartName = strings.Replace(PartName, "\x00", "", -1)
				fmt.Fprintf(graphDot, "<tr>  <td>Name</td> <td>%s</td>  </tr>\n", PartName)
			}
		}

		fmt.Fprintf(graphDot, "</table>\n")
		fmt.Fprintf(graphDot, ">];\n")

		if posExtendida != -1 {

			var posEBR int = 1
			var extendedBoot EBR
			File.Seek(0, 0)
			extendedBoot = readEBR(File, int64(MB.Particion[posExtendida].PartStart))

			for extendedBoot.PartNext != -1 && (extendedBoot.PartNext < MB.Particion[posExtendida].PartStart+MB.Particion[posExtendida].PartSize) {
				if extendedBoot.PartStatus != '1' {

					fmt.Fprintf(graphDot, "\ntbl_%d[\nlabel=<\n ", posEBR)
					fmt.Fprintf(graphDot, "<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n ")
					fmt.Fprintf(graphDot, "<tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>EBR</font></b></td></tr>")
					fmt.Fprintf(graphDot, "<tr ><td width='150'><b>Nombre</b></td> <td width='150'><b>Valor</b></td>  </tr>\n")
					var status string
					if extendedBoot.PartStatus == '0' {
						status = "0"
					} else if extendedBoot.PartStatus == '2' {
						status = "2"
					} else if extendedBoot.PartStatus == '1' {
						status = "1"
					}

					fmt.Fprintf(graphDot, "<tr>  <td>Status</td> <td>%s</td>  </tr>\n", string(status[:]))
					fmt.Fprintf(graphDot, "<tr>  <td>Fit</td> <td>%c</td>  </tr>\n", extendedBoot.PartFit)
					fmt.Fprintf(graphDot, "<tr>  <td>Start</td> <td>%d</td>  </tr>\n", extendedBoot.PartStart)
					fmt.Fprintf(graphDot, "<tr>  <td>Size</td> <td>%d</td>  </tr>\n", extendedBoot.PartSize)
					fmt.Fprintf(graphDot, "<tr>  <td>Next</td> <td>%d</td>  </tr>\n", extendedBoot.PartNext)
					PartNameExt := string(extendedBoot.PartName[:])
					PartNameExt = strings.Replace(PartNameExt, "\x00", "", -1)
					fmt.Fprintf(graphDot, "<tr>  <td>Name</td> <td>%s</td>  </tr>\n", PartNameExt)
					fmt.Fprintf(graphDot, "</table>\n")
					fmt.Fprintf(graphDot, ">];\n")

					posEBR++

				}
				if extendedBoot.PartNext == -1 {
				} else {
					extendedBoot = readEBR(File, int64(extendedBoot.PartNext))
				}

			}

		}

		fmt.Fprintf(graphDot, "}\n")
		graphDot.Close()
		File.Close()
		exec.Command("dot", "-Tpng", "-o", "/home/jose/Escritorio/graficaEBR.png", "Reportes/graficaEBR.dot").Output()
	}

}

//ReporteDisco is...
func ReporteDisco(direccion string) {

	var auxDir string = direccion

	if VerificarRuta(auxDir) {
		fp := getFile(auxDir)
		os.Create("Reportes/graficaDisco.dot")
		graphDot := getFile("Reportes/graficaDisco.dot")

		fmt.Fprintf(graphDot, "digraph G{\n")
		fmt.Fprintf(graphDot, "  tbl [\n    shape=box\n    label=<\n")
		fmt.Fprintf(graphDot, "     <table border='0' cellborder='1' width='600' height='200' color='lightblue'>\n")
		fmt.Fprintf(graphDot, "     <tr>\n")
		fmt.Fprintf(graphDot, "     <td  cellspacing= '0' height='200' width='100'> MBR </td>\n")

		var masterboot MBR
		fp.Seek(0, 0)
		masterboot = readMBR(fp)

		for i := 0; i < 4; i++ {

			if masterboot.Particion[i].PartStart != -1 {

				if masterboot.Particion[i].PartStatus != '1' {

					if masterboot.Particion[i].PartType == 'P' { //Verificar Primaria

						fmt.Fprintf(graphDot, "     <td cellspacing= '0' height='200' width='200'>PRIMARIA </td>\n")

					} else {
						//Particion Extendida

						fmt.Fprintf(graphDot, "     <td cellspacing= '0' height='200' width='200'>\n     <table border='0'  height='200' WIDTH='200' cellborder='1'>\n")
						fmt.Fprintf(graphDot, "     <tr>  <td height='60' colspan='15'>EXTENDIDA</td>  </tr>\n     <tr>\n")

						var extendedBoot EBR
						fp.Seek(0, 0)
						extendedBoot = readEBR(fp, int64(masterboot.Particion[i].PartStart))

						if extendedBoot.PartSize != 0 { //Si hay mas de alguna logica

							for extendedBoot.PartNext != -1 && (extendedBoot.PartNext < (masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)) {

								fmt.Fprintf(graphDot, "     <td cellspacing= '0' height='140'>EBR</td>\n")
								fmt.Fprintf(graphDot, "     <td cellspacing= '0' height='140'>LOGICA</td>\n")

								if extendedBoot.PartNext == -1 {

								} else {
									extendedBoot = readEBR(fp, int64(extendedBoot.PartNext))
								}
							}
						} else {
							fmt.Fprintf(graphDot, "     <td cellspacing= '0' height='140'></td>")
						}

						fmt.Fprintf(graphDot, "     </tr>\n     </table>\n     </td>\n")

					}
				}
			}
		}

		fmt.Fprintf(graphDot, "     <td height='200'> LIBRE </td>")

		fmt.Fprintf(graphDot, "     </tr> \n     </table>        \n>];\n\n}")
		graphDot.Close()
		fp.Close()
		exec.Command("dot", "-Tpng", "-o", "/home/jose/Escritorio/grafica.png", "Reportes/graficaDisco.dot").Output()
	} else {
		ErrorMessage("[REP] -> No se encuentra el disco")
	}
}

/*** FASE 2 ****/

//ReporteSuperBloque is...
func ReporteSuperBloque(ID string) {
	//TODO : Reporte SuperBloque
	PartStart := listaParticiones.GetPartStart(ID)
	PathDisco := listaParticiones.GetDireccion(ID)
	File := getFile(PathDisco)
	SB := readSuperBloque(File, int64(PartStart))

	os.Create("Reportes/graficaSuperBloque.dot")
	graphDot := getFile("Reportes/graficaSuperBloque.dot")

	//Empezamos a escribir en el archivo
	fmt.Fprintf(graphDot, "digraph G{ \n")
	fmt.Fprintf(graphDot, "node [shape=plaintext]\n")
	fmt.Fprintf(graphDot, "tbl[\nlabel=<\n")
	fmt.Fprintf(graphDot, "<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n")
	fmt.Fprintf(graphDot, " <tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>Super Bloque</font></b></td></tr>")
	fmt.Fprintf(graphDot, "<tr>  <td width='230'> <b>Atributo</b> </td> <td width='230'> <b>Valor</b> </td>  </tr>\n")
	SBName := string(SB.NombreHD[:])
	SBName = strings.Replace(SBName, "\x00", "", -1)
	fmt.Fprintf(graphDot, "<tr>  <td>NombreHD</td><td>%s</td>  </tr>\n", SBName)
	fmt.Fprintf(graphDot, "<tr>  <td>ArbolVirtualCount</td><td>%d</td>  </tr>\n", SB.ArbolVirtualCount)
	fmt.Fprintf(graphDot, "<tr>  <td>DetalleDirectorioCount</td><td>%d</td>  </tr>\n", SB.DetalleDirectorioCount)
	fmt.Fprintf(graphDot, "<tr>  <td>InodosCount</td><td>%d</td>  </tr>\n", SB.InodosCount)
	fmt.Fprintf(graphDot, "<tr>  <td>BloquesCount</td><td>%d</td>  </tr>\n", SB.BloquesCount)
	fmt.Fprintf(graphDot, "<tr>  <td>ArbolVirtualFree</td><td>%d</td>  </tr>\n", SB.ArbolVirtualFree)
	fmt.Fprintf(graphDot, "<tr>  <td>DetalleDirectorioFree</td><td>%d</td>  </tr>\n", SB.DetalleDirectorioFree)
	fmt.Fprintf(graphDot, "<tr>  <td>InodosFree</td><td>%d</td>  </tr>\n", SB.InodosFree)
	fmt.Fprintf(graphDot, "<tr>  <td>BloquesFree</td><td>%d</td>  </tr>\n", SB.BloquesFree)
	fmt.Fprintf(graphDot, "<tr>  <td>FechaCreacion</td><td>%s</td>  </tr>\n", string(SB.DateCreacion[:]))
	fmt.Fprintf(graphDot, "<tr>  <td>FechaUltimoMontaje</td><td>%s</td>  </tr>\n", string(SB.DateUltimoMontaje[:]))
	fmt.Fprintf(graphDot, "<tr>  <td>MontajesCount</td><td>%d</td>  </tr>\n", SB.MontajesCount)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBmArbolDirectorio</td><td>%d</td>  </tr>\n", SB.StartBmArbolDirectorio)
	fmt.Fprintf(graphDot, "<tr>  <td>StartArbolDirectorio</td><td>%d</td>  </tr>\n", SB.StartArbolDirectorio)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBmArbolDirectorio</td><td>%d</td>  </tr>\n", SB.StartBmDetalleDirectorio)
	fmt.Fprintf(graphDot, "<tr>  <td>StartDetalleDirectorio</td><td>%d</td>  </tr>\n", SB.StartDetalleDirectorio)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBmInodos</td><td>%d</td>  </tr>\n", SB.StartBmInodos)
	fmt.Fprintf(graphDot, "<tr>  <td>StartInodos</td><td>%d</td>  </tr>\n", SB.StartInodos)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBmBloques</td><td>%d</td>  </tr>\n", SB.StartBmBloques)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBloques</td><td>%d</td>  </tr>\n", SB.StartBloques)
	fmt.Fprintf(graphDot, "<tr>  <td>StartBitacora</td><td>%d</td>  </tr>\n", SB.StartLog)
	fmt.Fprintf(graphDot, "<tr>  <td>SizeAVD</td><td>%d</td>  </tr>\n", SB.SizeStructAvd)
	fmt.Fprintf(graphDot, "<tr>  <td>SizeDD</td><td>%d</td>  </tr>\n", SB.SizeStructDd)
	fmt.Fprintf(graphDot, "<tr>  <td>SizeInodo</td><td>%d</td>  </tr>\n", SB.SizeStructInodo)
	fmt.Fprintf(graphDot, "<tr>  <td>SizeBloque</td><td>%d</td>  </tr>\n", SB.SizeStructBloque)
	fmt.Fprintf(graphDot, "<tr>  <td>FirstFreeAVD</td><td>%d</td>  </tr>\n", SB.FirstFreeAvd)
	fmt.Fprintf(graphDot, "<tr>  <td>FirstFreeDD</td><td>%d</td>  </tr>\n", SB.FirstFreeDd)
	fmt.Fprintf(graphDot, "<tr>  <td>FirstFreeInodos</td><td>%d</td>  </tr>\n", SB.FirstFreeInodo)
	fmt.Fprintf(graphDot, "<tr>  <td>FirstFreeBloque</td><td>%d</td>  </tr>\n", SB.FirstFreeInodo)
	fmt.Fprintf(graphDot, "<tr>  <td>MagicNum</td><td>%d</td>  </tr>\n", SB.MagicNum)

	fmt.Fprintf(graphDot, "</table>\n")
	fmt.Fprintf(graphDot, ">];\n}")
}

//ReporteTreeComplete is...
func ReporteTreeComplete(path string, id string) {

	RutaDisco := listaParticiones.GetDireccion(id)
	var Grafica string

	if RutaDisco != "null" {

		PartStart := listaParticiones.GetPartStart(id)
		File := getFile(RutaDisco)

		SuperBloque := readSuperBloque(File, int64(PartStart))
		Root := readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

		os.Create("Reportes/graficaTreeComplete.dot")
		graphDot := getFile("Reportes/graficaTreeComplete.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")

		RecorrerArbolReporte(Root, SuperBloque, File, &Grafica, false, 0, false, true)

		fmt.Fprintf(graphDot, Grafica)

		fmt.Fprintf(graphDot, "}")

	} else {
		ErrorMessage("[REP] -> Particion no montada")
	}

}

//RecorrerArbolReporte is...
func RecorrerArbolReporte(arbol Arbol, Superbloque SB, file *os.File, Grafica *string, avd bool, ptr int, onlyAVD bool, showInodes bool) {

	//fmt.Println("Carpeta", string(arbol.AVDNombreDirectorio[:]))
	var Graph string
	texto := string(arbol.AVDNombreDirectorio[:])
	texto = strings.Replace(texto, "\x00", "", -1)

	*Grafica += "tbl"
	var num string
	num = fmt.Sprint(num, ptr)
	*Grafica += num
	*Grafica += "[label=<\n"
	*Grafica += "<table border='0' cellborder='1' cellspacing='0'>\n"
	*Grafica += "<tr>\n"
	if !avd {
		*Grafica += "<td colspan='8' bgcolor= 'lightblue' >" + texto + "</td>\n"
	} else {
		*Grafica += "<td colspan='8' bgcolor= 'forestgreen' >" + texto + "</td>\n"
	}
	*Grafica += "</tr>\n"
	*Grafica += "<tr>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >1</td>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >2</td>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >3</td>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >4</td>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >5</td>\n"
	*Grafica += "<td bgcolor='lightblue' width='20' >6</td>\n"
	*Grafica += "<td bgcolor='deepskyblue4' width='30' >DD</td>\n"
	*Grafica += "<td bgcolor='forestgreen' width='20' >AVD</td>\n"
	*Grafica += "</tr>\n"
	*Grafica += "<tr>\n"
	for i := 0; i < 6; i++ {

		*Grafica += "<td width='20'>"
		var aux string
		aux = fmt.Sprint(aux, arbol.Subirectorios[i])
		*Grafica += aux
		*Grafica += "</td>\n"

	}

	*Grafica += "<td width='20'>"
	var aux string
	aux = fmt.Sprint(aux, arbol.DetalleDirectorio)
	*Grafica += aux
	*Grafica += "</td>\n"
	*Grafica += "<td width='20'>"

	var aux1 string
	aux1 = fmt.Sprint(aux1, arbol.VirtualDirectorio)
	*Grafica += aux1
	*Grafica += "</td>\n"
	*Grafica += "</tr>\n"
	*Grafica += "</table>\n"
	*Grafica += ">];\n"

	/*
	 * LEEMOS EL DETALLE DIRECTORIO
	 */

	if !onlyAVD {
		//Apuntador al DD
		*Grafica += "tbl" + num
		*Grafica += "->"
		*Grafica += "tbl" + num + "DD\n"
		GenerarDetalleDirectorio(num, texto+"DD", int(arbol.DetalleDirectorio), file, Superbloque, &Graph, showInodes)
		*Grafica += Graph
	}

	if arbol.VirtualDirectorio != -1 {
		ApuntadorAVDCopia := arbol.VirtualDirectorio
		var CopiaCarpeta Arbol
		CopiaCarpeta = readArbolVirtualDirectorio(file, int64(Superbloque.StartArbolDirectorio+(ApuntadorAVDCopia*int32(unsafe.Sizeof(arbol)))))
		/*
		 *	PUNTERO PARA EL AVD VIRTUAL
		 */
		*Grafica += "tbl" + num
		*Grafica += "->"
		var num2 string
		num2 = fmt.Sprint(num2, ApuntadorAVDCopia)
		*Grafica += "tbl" + num2 + "\n"

		RecorrerArbolReporte(CopiaCarpeta, Superbloque, file, Grafica, true, int(ApuntadorAVDCopia), onlyAVD, showInodes)
	}

	/*
	 * RECORREMOS RECURSIVAMENTE
	 */
	for j := 0; j < 6; j++ {
		Apuntador := arbol.Subirectorios[j]
		if Apuntador != -1 {
			var CarpetaHija Arbol
			CarpetaHija = readArbolVirtualDirectorio(file, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CarpetaHija)))))
			/*
			 * APUNTADOR EN EL GRAPHVIZ
			 */
			*Grafica += "tbl" + num
			*Grafica += "->"
			var num3 string
			num3 = fmt.Sprint(num3, Apuntador)
			*Grafica += "tbl" + num3 + "\n"
			/*
			 * LLAMAMOS RECURSIVAMENTE AL METODO
			 */
			RecorrerArbolReporte(CarpetaHija, Superbloque, file, Grafica, false, int(Apuntador), onlyAVD, showInodes)
		}
	}

}

//GenerarDetalleDirectorio is...
func GenerarDetalleDirectorio(Nombre string, name string, puntero int, file *os.File, super SB, GraficaDD *string, showInodes bool) {
	var Detalle DetalleDirectorio
	PosicionDD := super.StartDetalleDirectorio + int32(puntero*int(unsafe.Sizeof(Detalle)))
	Detalle = readDetalleDirectorio(file, int64(PosicionDD))

	Nombre += "DD"

	*GraficaDD += "tbl" + Nombre + "[label=<\n"
	*GraficaDD += "<table border='0' cellborder='1' cellspacing='0'>\n"
	*GraficaDD += "<tr>"
	*GraficaDD += "<td bgcolor='deepskyblue4' width='100' colspan='2'>" + name + "</td>\n"
	*GraficaDD += "</tr>\n"
	*GraficaDD += "<tr>\n"
	*GraficaDD += "<td>Virtual</td>\n"
	var ptr string
	ptr = fmt.Sprint(ptr, Detalle.DDApDetalleDirectorio)
	*GraficaDD += "<td>" + ptr + "</td>\n"
	*GraficaDD += "</tr>\n"
	for i := 0; i < 5; i++ {
		var apuntador string
		nombreFile := string(Detalle.DDArrayFiles[i].DDFileNombre[:])
		nombreFile = strings.Replace(nombreFile, "\x00", "", -1)
		apuntador = fmt.Sprint(apuntador, Detalle.DDArrayFiles[i].DDFileApInodo)

		*GraficaDD += "<tr>\n"
		if nombreFile == "" {
			*GraficaDD += "<td>" + "-1" + "</td>\n"
		} else {
			*GraficaDD += "<td>" + nombreFile + "</td>\n"
		}
		*GraficaDD += "<td>" + apuntador + "</td>\n"
		*GraficaDD += "</tr>\n"

	}
	*GraficaDD += "</table>\n>];\n"

	if showInodes {
		for i := 0; i < 5; i++ {

			nombreFile := string(Detalle.DDArrayFiles[i].DDFileNombre[:])
			nombreFile = strings.Replace(nombreFile, "\x00", "", -1)

			var Inodo TablaInodo
			ApInodo := Detalle.DDArrayFiles[i].DDFileApInodo

			if ApInodo != -1 {
				*GraficaDD += "tbl" + Nombre + "->"
				var NameInodo string
				NameInodo = fmt.Sprint(NameInodo, ApInodo)
				*GraficaDD += "tblInodo" + NameInodo + "\n"

				Inodo = readInodo(file, int64(super.StartInodos+(ApInodo*int32(unsafe.Sizeof(Inodo)))))
				GenerarReporteInodo(Inodo, file, super, GraficaDD, NameInodo, nombreFile)

			}
		}
	}
	if Detalle.DDApDetalleDirectorio != -1 {
		var aux string
		aux = fmt.Sprint(aux, Detalle.DDApDetalleDirectorio)
		fmt.Println("Test", "tbl", aux, "DD")
		*GraficaDD += "tbl" + Nombre + "-> " + "tbl" + aux + "DDDD" + ";\n"
		aux += "DD"
		GenerarDetalleDirectorio(aux, name, int(Detalle.DDApDetalleDirectorio), file, super, GraficaDD, showInodes)
	}
}

//GenerarReporteInodo is...
func GenerarReporteInodo(Inodo TablaInodo, File *os.File, SuperBloque SB, GraficaDD *string, NameInodo string, nombreFile string) {
	/*
	 * TABLA INODO
	 */
	*GraficaDD += "tblInodo" + NameInodo
	*GraficaDD += "[label=<\n"
	*GraficaDD += "<table border='0' cellborder='1' cellspacing='0'>\n"
	*GraficaDD += "<tr><td bgcolor='darkorange' width='100' colspan='5'>" + nombreFile + "</td>\n"
	*GraficaDD += "</tr>\n"
	*GraficaDD += "<tr>\n"
	*GraficaDD += "<td>1</td>\n"
	*GraficaDD += "<td>2</td>\n"
	*GraficaDD += "<td>3</td>\n"
	*GraficaDD += "<td>4</td>\n"
	*GraficaDD += "<td bgcolor='darkorange' width='25'>Inodo</td>\n"
	*GraficaDD += "</tr>\n"
	*GraficaDD += "<tr>\n"
	var ptr1 string
	ptr1 = fmt.Sprint(ptr1, Inodo.IApIndirecto)
	*GraficaDD += "<td>" + ptr1 + "</td>\n"
	for i := 0; i < 4; i++ {
		var aux string
		aux = fmt.Sprint(aux, Inodo.IArrayBloques[i])
		*GraficaDD += "<td>" + aux + "</td>"
	}
	*GraficaDD += "</tr>\n"
	*GraficaDD += "</table>\n>];"

	for i := 0; i < 4; i++ {
		ApBloque := Inodo.IArrayBloques[i]

		if ApBloque != -1 {
			var aux1 string
			aux1 = fmt.Sprint(aux1, ApBloque)
			var Block Bloque
			Block = readBloque(File, int64(SuperBloque.StartBloques+(ApBloque*int32(unsafe.Sizeof(Block)))))

			texto := string(Block.Texto[:])
			texto = strings.Replace(texto, "\x00", "", -1)

			*GraficaDD += "tblInodo" + NameInodo + "->" + "tblBloque" + aux1 + "\n"
			*GraficaDD += "tblBloque" + aux1 + "[label=<\n"
			*GraficaDD += "<table border='0' cellborder='1' cellspacing='0'>\n"
			*GraficaDD += "<tr>\n"
			*GraficaDD += "<td width='200' bgcolor= 'lightblue' >'" + texto + "'</td>\n"
			*GraficaDD += "</tr>\n"
			*GraficaDD += "</table>\n"
			*GraficaDD += ">];\n"
		}
	}
	if Inodo.IApIndirecto != -1 {
		var apunt string
		apunt = fmt.Sprint(apunt, Inodo.IApIndirecto)
		apunt += "Inode"
		*GraficaDD += "tblInodo" + NameInodo + "->tblInodo" + apunt + "\n"
		var InodoIndirecto TablaInodo
		InodoIndirecto = readInodo(File, int64(SuperBloque.StartInodos+(Inodo.IApIndirecto*int32(unsafe.Sizeof(InodoIndirecto)))))
		GenerarReporteInodo(InodoIndirecto, File, SuperBloque, GraficaDD, apunt, nombreFile)
	}
}

//ReporteTreeFile is...
func ReporteTreeFile(carpeta string, id string, path string) {
	PathDisco := listaParticiones.GetDireccion(id)

	if PathDisco != "null" {

		var Grafica string

		PartStart := listaParticiones.GetPartStart(id)
		File := getFile(PathDisco)

		SuperBloque := readSuperBloque(File, int64(PartStart))
		Root := readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

		os.Create("Reportes/graficaTreeFile.dot")
		graphDot := getFile("Reportes/graficaTreeFile.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")

		Rutas := strings.Split(carpeta, "/")
		Rutas = Rutas[1:]
		fmt.Println(Rutas)
		BuscarCarpeta(Root, Rutas, PathDisco, SuperBloque, &Grafica, true)

		fmt.Fprintf(graphDot, Grafica)

		fmt.Fprintf(graphDot, "}")

	} else {
		ErrorMessage("[REP] -> Particion no montada")
	}
}

//BuscarCarpeta is...
func BuscarCarpeta(root Arbol, Rutas []string, PathDisco string, Superbloque SB, Grafica *string, mostrarInodes bool) {
	if len(Rutas) == 0 {
		return
	}

	if VerificarRuta(PathDisco) {

		File := getFile(PathDisco) //Leemos el disco
		var Apuntador int32 = 0

		//recorremos los 6 subdirectorios del AVD
		for i := 0; i < 6; i++ {
			Apuntador = root.Subirectorios[i]
			var CarpetaHija Arbol
			CarpetaHija = readArbolVirtualDirectorio(File, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CarpetaHija)))))

			//Verificamos el nombre de la carpeta
			var nameCarpetaByte [16]byte
			copy(nameCarpetaByte[:], Rutas[0])

			if bytes.Compare(CarpetaHija.AVDNombreDirectorio[:], nameCarpetaByte[:]) == 0 {
				Rutas = Rutas[1:]
				if len(Rutas) == 0 {
					Apuntador = CarpetaHija.DetalleDirectorio
					var Archivos DetalleDirectorio
					Archivos = readDetalleDirectorio(File, int64(Superbloque.StartDetalleDirectorio+(Apuntador*int32(unsafe.Sizeof(Archivos)))))
					//Cerramos el archivo
					if mostrarInodes {
						RecorrerArbolReporte(CarpetaHija, Superbloque, File, Grafica, false, int(Apuntador), false, true)
					} else {
						RecorrerArbolReporte(CarpetaHija, Superbloque, File, Grafica, false, int(Apuntador), false, false)
					}
					return
				}
				File.Close()
				BuscarCarpeta(CarpetaHija, Rutas, PathDisco, Superbloque, Grafica, mostrarInodes)
				return
			}

		}
		// Buscamos en la copia de la carpeta
		Apuntador = root.VirtualDirectorio
		var CopiaCarpeta Arbol
		CopiaCarpeta = readArbolVirtualDirectorio(File, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CopiaCarpeta)))))
		File.Close()
		/*
		 * LLAMAMOS EL METODO RECURSIVAMENTE
		 */
		BuscarCarpeta(CopiaCarpeta, Rutas, PathDisco, Superbloque, Grafica, mostrarInodes)
		return
	}
	ErrorMessage("[MKFILE] -> No hay ningun disco en la ruta indicada")
}

/*
 *	REPORTE DIRECTORIO
 */

//ReporteDirectorio is...
func ReporteDirectorio(path string, id string) {
	RutaDisco := listaParticiones.GetDireccion(id)
	var Grafica string

	if RutaDisco != "null" {

		PartStart := listaParticiones.GetPartStart(id)
		File := getFile(RutaDisco)

		SuperBloque := readSuperBloque(File, int64(PartStart))
		Root := readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

		os.Create("Reportes/graficaDirectorio.dot")
		graphDot := getFile("Reportes/graficaDirectorio.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")

		RecorrerArbolReporte(Root, SuperBloque, File, &Grafica, false, 0, true, false)

		fmt.Fprintf(graphDot, Grafica)

		fmt.Fprintf(graphDot, "}")

	} else {
		ErrorMessage("[REP] -> Particion no montada")
	}
}

//ReporteTreeDirectorio is...
func ReporteTreeDirectorio(carpeta string, path string, id string) {
	PathDisco := listaParticiones.GetDireccion(id)

	if PathDisco != "null" {

		var Grafica string

		PartStart := listaParticiones.GetPartStart(id)
		File := getFile(PathDisco)

		SuperBloque := readSuperBloque(File, int64(PartStart))
		Root := readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

		os.Create("Reportes/graficaTreeDirectorio.dot")
		graphDot := getFile("Reportes/graficaTreeDirectorio.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")

		Rutas := strings.Split(carpeta, "/")
		Rutas = Rutas[1:]
		fmt.Println(Rutas)
		BuscarCarpeta(Root, Rutas, PathDisco, SuperBloque, &Grafica, false)

		fmt.Fprintf(graphDot, Grafica)

		fmt.Fprintf(graphDot, "}")

	} else {
		ErrorMessage("[REP] -> Particion no montada")
	}

}

/*
 * METODOS PARA ESCRIBIR
 */

//WriteOne is...
func WriteOne(file *os.File, uno byte) {
	s1 := &uno
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteDD is...
func WriteDD(file *os.File, DD DetalleDirectorio) {
	s1 := &DD
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteSB is...
func WriteSB(file *os.File, Super SB) {
	s1 := &Super
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteAVD is...
func WriteAVD(file *os.File, AVD Arbol) {
	s1 := &AVD
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteInode is...
func WriteInode(file *os.File, Inode TablaInodo) {
	s1 := &Inode
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteBloque is...
func WriteBloque(file *os.File, Block Bloque) {
	s1 := &Block
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

//WriteBitacora is...
func WriteBitacora(file *os.File, log Bitacora) {
	s1 := &log
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	file.Write(binario2.Bytes())
}

/*
 * REPORTES BITMAPS
 */

//ReporteBMarbdir is...
func ReporteBMarbdir(path string, id string) {

	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		if VerificarRuta(RutaDisco) {
			PartStart := listaParticiones.GetPartStart(id)
			File := getFile(RutaDisco)
			SuperBloque := readSuperBloque(File, int64(PartStart))
			CantidadEstructuras := int(SuperBloque.MagicNum)
			StartBitMap := int(SuperBloque.StartBmArbolDirectorio)

			os.Create("Reportes/BitMapArbolDirectorio.txt")
			BitMapTxt := getFile("Reportes/BitMapArbolDirectorio.txt")
			var contador int
			fmt.Fprintf(BitMapTxt, "\n\nBitMap Arbol Virtual de Directorio\n\n")
			for i := 0; i < CantidadEstructuras; i++ {
				contador++
				i, _ := strconv.Atoi(string(readByte(File, int64(StartBitMap+i))))
				var Byte string
				Byte = fmt.Sprint(Byte, i)
				if contador == 51 {
					fmt.Fprintf(BitMapTxt, "|\n")
					contador = 0
				} else {
					fmt.Fprintf(BitMapTxt, "| "+Byte+" ")
				}
			}

		} else {
			ErrorMessage("[REP] -> No se encuentra el disco")
		}

	} else {
		ErrorMessage("[REP] -> No hay ninguna particion montada con ese id")
	}

}

//ReporteBMdetdir is...
func ReporteBMdetdir(path string, id string) {
	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		if VerificarRuta(RutaDisco) {
			PartStart := listaParticiones.GetPartStart(id)
			File := getFile(RutaDisco)
			SuperBloque := readSuperBloque(File, int64(PartStart))
			CantidadEstructuras := int(SuperBloque.MagicNum)
			StartBitMap := int(SuperBloque.StartBmDetalleDirectorio)

			os.Create("Reportes/BitMapDetalleDirectorio.txt")
			BitMapTxt := getFile("Reportes/BitMapDetalleDirectorio.txt")
			var contador int
			fmt.Fprintf(BitMapTxt, "\n\nBitMap Arbol Detalle Directorio\n\n")
			for i := 0; i < CantidadEstructuras; i++ {
				contador++
				i, _ := strconv.Atoi(string(readByte(File, int64(StartBitMap+i))))
				var Byte string
				Byte = fmt.Sprint(Byte, i)
				if contador == 51 {
					fmt.Fprintf(BitMapTxt, "\n")
					contador = 0
				} else {
					fmt.Fprintf(BitMapTxt, "| "+Byte+" ")
				}
			}

		} else {
			ErrorMessage("[REP] -> No se encuentra el disco")
		}

	} else {
		ErrorMessage("[REP] -> No hay ninguna particion montada con ese id")
	}
}

//ReporteBMinode is...
func ReporteBMinode(path string, id string) {
	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		if VerificarRuta(RutaDisco) {
			PartStart := listaParticiones.GetPartStart(id)
			File := getFile(RutaDisco)
			SuperBloque := readSuperBloque(File, int64(PartStart))
			CantidadEstructuras := int(SuperBloque.MagicNum)
			StartBitMap := int(SuperBloque.StartBmInodos)

			os.Create("Reportes/BitMapInode.txt")
			BitMapTxt := getFile("Reportes/BitMapInode.txt")
			var contador int
			fmt.Fprintf(BitMapTxt, "\n\nBitMap Inode\n\n")
			for i := 0; i < CantidadEstructuras; i++ {
				contador++
				i, _ := strconv.Atoi(string(readByte(File, int64(StartBitMap+i))))
				var Byte string
				Byte = fmt.Sprint(Byte, i)
				if contador == 51 {
					fmt.Fprintf(BitMapTxt, "\n")
					contador = 0
				} else {
					fmt.Fprintf(BitMapTxt, "| "+Byte+" ")
				}
			}

		} else {
			ErrorMessage("[REP] -> No se encuentra el disco")
		}

	} else {
		ErrorMessage("[REP] -> No hay ninguna particion montada con ese id")
	}
}

//ReporteBMblock is...
func ReporteBMblock(path string, id string) {
	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		if VerificarRuta(RutaDisco) {
			PartStart := listaParticiones.GetPartStart(id)
			File := getFile(RutaDisco)
			SuperBloque := readSuperBloque(File, int64(PartStart))
			CantidadEstructuras := int(SuperBloque.MagicNum)
			StartBitMap := int(SuperBloque.StartBmBloques)

			os.Create("Reportes/BitMapBloque.txt")
			BitMapTxt := getFile("Reportes/BitMapBloque.txt")
			var contador int
			fmt.Fprintf(BitMapTxt, "\n\nBitMap Bloque")
			for i := 0; i < CantidadEstructuras; i++ {
				contador++
				i, _ := strconv.Atoi(string(readByte(File, int64(StartBitMap+i))))
				var Byte string
				Byte = fmt.Sprint(Byte, i)
				if contador == 51 {
					fmt.Fprintf(BitMapTxt, "\n")
					contador = 0
				} else {
					fmt.Fprintf(BitMapTxt, "| "+Byte+" ")
				}
			}

		} else {
			ErrorMessage("[REP] -> No se encuentra el disco")
		}

	} else {
		ErrorMessage("[REP] -> No hay ninguna particion montada con ese id")
	}
}

/*
 *	REPORTE BITACORA
 */

//ReporteBitacora is...
func ReporteBitacora(path string, id string) {

	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		if VerificarRuta(RutaDisco) {

			os.Create("Reportes/Bitacora.dot")
			logDot := getFile("Reportes/Bitacora.dot")

			PartStart := listaParticiones.GetPartStart(id)
			File := getFile(RutaDisco)
			SuperBloque := readSuperBloque(File, int64(PartStart))
			var log Bitacora
			fmt.Fprintf(logDot, "digraph G{\nrankdir=\"LR\"\nnode [shape=plaintext]\n")
			for i := 0; i < int(SuperBloque.FirstFreeBitacora); i++ {
				log = readBitacora(File, int64(SuperBloque.StartLog+int32(i*int(unsafe.Sizeof(log)))))
				var Index string
				Index = fmt.Sprint(Index, i)
				fmt.Fprintf(logDot, "tbl"+Index+"[label=<\n")
				fmt.Fprintf(logDot, "<table border='0' cellborder='1' cellspacing='0'>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				var namelog string = "log "
				namelog = fmt.Sprint(namelog, i+1)
				fmt.Fprintf(logDot, "<td colspan='2' bgcolor= 'lightblue' >%s</td>\n", namelog)
				fmt.Fprintf(logDot, "</tr>")
				fmt.Fprintf(logDot, "<tr>\n")
				//Nombre
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Nombre</td>\n")
				nombrelog := string(log.Nombre[:])
				nombrelog = strings.Replace(nombrelog, "\x00", "", -1)
				fmt.Fprintf(logDot, "<td width='300' >%s</td>\n", nombrelog)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				//Contenido
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Contenido</td>\n")
				contenido := string(log.Contenido[:])
				contenido = strings.Replace(contenido, "\x00", "", -1)
				fmt.Fprintf(logDot, "<td width='300' >%s</td>\n", contenido)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				//Fecha
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Fecha</td>\n")
				date := string(log.Fecha[:])
				date = strings.Replace(date, "\x00", "", -1)
				fmt.Fprintf(logDot, "<td width='300' >%s</td>\n", date)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				//TipoOperacion
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Tipo Operacion</td>\n")
				tipoOp := string(log.TipoOp[:])
				tipoOp = strings.Replace(tipoOp, "\x00", "", -1)
				fmt.Fprintf(logDot, "<td width='300' >%s</td>\n", tipoOp)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				//Tipo
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Tipo</td>\n")
				fmt.Fprintf(logDot, "<td width='300' >%d</td>\n", log.Tipo)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "<tr>\n")
				//Size
				fmt.Fprintf(logDot, "<td bgcolor='lightblue' width='200' >Size</td>\n")
				fmt.Fprintf(logDot, "<td width='300' >%d</td>\n", log.Size)
				fmt.Fprintf(logDot, "</tr>\n")
				fmt.Fprintf(logDot, "</table>\n>];")
			}
			fmt.Fprintf(logDot, "\n}")
		} else {
			ErrorMessage("[REP] -> No se encuentra el disco")
		}

	} else {
		ErrorMessage("[REP] -> No hay ninguna particion montada con ese id")
	}
}

/*
 *	REPORTE LS
 */

//ReporteLS is...
func ReporteLS(path string, id string) {

}

/*
 *	C O M A N D O S   F A S E   2
 */

//MKFS is...
func MKFS(id string) {
	Formatear(id)
}

//Formatear is...
func Formatear(id string) {

	pathD := listaParticiones.GetDireccion(id) //Obtenemos la direccion del disco

	if pathD != "null" {

		/*
		 *   INSTANCIANDO LOS STRUCTS
		 */
		SB := SB{}
		AVD := InicializarAVD(Arbol{})
		DD := InicializarDD(DetalleDirectorio{})
		Inodo := InicializarInodo(TablaInodo{})
		Bloque := InicializarBloque(Bloque{})
		Bitacora := InicializarBitacora(Bitacora{})
		/*
		 *  OBTENIENDO EL SIZE
		 */
		SBSize := int(unsafe.Sizeof(SB))
		AVDSize := int(unsafe.Sizeof(AVD))
		DDSize := int(unsafe.Sizeof(DD))
		InodoSize := int(unsafe.Sizeof(Inodo))
		BloqueSize := int(unsafe.Sizeof(Bloque))
		BitacoraSize := int(unsafe.Sizeof(Bitacora))

		File := getFile(pathD)                         //Obtenemos el disco
		PartName := listaParticiones.GetPartName(id)   // Obtenemos el PartName
		PartSize := listaParticiones.GetPartSize(id)   // Obtenemos el PartSize
		PartStart := listaParticiones.GetPartStart(id) // Obtenemos el partStart

		//fmt.Println("PartStart", PartStart)

		//Formula de Cantidad de estructuras
		var CantidadEstructuras int = (PartSize - (2 * SBSize)) / (27 + AVDSize + DDSize + (5*InodoSize + (20 * BloqueSize) + BitacoraSize))

		//Cantidad de elementos de cada Struct
		var CantidadAVD int = CantidadEstructuras
		var CantidadDD int = CantidadEstructuras
		var CantidadInodos int = 5 * CantidadEstructuras
		var CantidadBloques int = 20 * CantidadEstructuras
		//var CantidadBitacora int = CantidadEstructuras
		//fmt.Println(CantidadAVD, CantidadDD, CantidadInodos, CantidadBloques, CantidadBitacora)
		/*
		 *  INICIALIZANDO EL SUPER BLOQUE
		 */
		copy(SB.NombreHD[:], PartName)
		//Cantidad de elementos
		SB.ArbolVirtualCount = int32(CantidadAVD)
		SB.DetalleDirectorioCount = int32(CantidadDD)
		SB.InodosCount = int32(CantidadInodos)
		SB.BloquesCount = int32(CantidadBloques)
		SB.ArbolVirtualFree = int32(CantidadAVD)     //Se crea carpeta raiz.
		SB.DetalleDirectorioFree = int32(CantidadDD) //Detalle de directorio de la carpeta raiz
		SB.InodosFree = int32(CantidadInodos)
		SB.BloquesFree = int32(CantidadBloques)
		//Obtenemos fecha actual
		dt := time.Now()
		fecha := dt.Format("01-02-2006 15:04:05")
		copy(SB.DateCreacion[:], fecha)
		copy(SB.DateUltimoMontaje[:], fecha)
		SB.MontajesCount = 0
		//Start Cada Struct
		SB.StartBmArbolDirectorio = int32(PartStart + int(unsafe.Sizeof(SB)))
		SB.StartArbolDirectorio = SB.StartBmArbolDirectorio + int32(CantidadEstructuras)
		SB.StartBmDetalleDirectorio = SB.StartArbolDirectorio + int32((CantidadEstructuras * int(unsafe.Sizeof(AVD))))
		SB.StartDetalleDirectorio = SB.StartBmDetalleDirectorio + int32(CantidadEstructuras)
		SB.StartBmInodos = SB.StartDetalleDirectorio + int32((CantidadEstructuras * int(unsafe.Sizeof(DD))))
		SB.StartInodos = SB.StartBmInodos + int32((5 * CantidadEstructuras))
		SB.StartBmBloques = SB.StartInodos + int32((5 * CantidadEstructuras * int(unsafe.Sizeof(Inodo))))
		SB.StartBloques = SB.StartBmBloques + int32((20 * CantidadEstructuras))
		SB.StartLog = SB.StartBloques + int32((20 * CantidadEstructuras * int(unsafe.Sizeof(Bloque)))) //Bitacora.
		//Espacios Libres
		SB.FirstFreeAvd = 0
		SB.FirstFreeDd = 0
		SB.FirstFreeInodo = 0
		SB.FirstFreeBloque = 0
		SB.FirstFreeBitacora = 0
		//Magic Num
		SB.MagicNum = int32(CantidadEstructuras)
		//Struct Size
		SB.SizeStructAvd = int32(unsafe.Sizeof(AVD))
		SB.SizeStructDd = int32(unsafe.Sizeof(DD))
		SB.SizeStructInodo = int32(unsafe.Sizeof(Inodo))
		SB.SizeStructBloque = int32(unsafe.Sizeof(Bloque))

		//fmt.Println("Inicio Bitacora", SB.StartLog+int32(CantidadEstructuras))

		//Escribo el super bloque al inicio de la particion
		reWriteSuperBloque(File, SB, int64(PartStart))

		//BitMap AVD
		for i := 0; i < CantidadEstructuras; i++ {
			File.Seek(int64((SB.StartBmArbolDirectorio + int32(i))), 0)
			var cero = '0'
			s1 := &cero
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//AVD = Arbol
		for i := 0; i < CantidadEstructuras; i++ {
			File.Seek(int64((SB.StartArbolDirectorio + int32((i * int(unsafe.Sizeof(AVD)))))), 0)
			s1 := &AVD
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//BitMap DD
		for i := 0; i < CantidadEstructuras; i++ {
			File.Seek(int64(SB.StartBmDetalleDirectorio+int32(i)), 0)
			var cero = '0'
			s1 := &cero
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//DD
		for i := 0; i < CantidadEstructuras; i++ {
			File.Seek(int64(SB.StartDetalleDirectorio+int32((i*int(unsafe.Sizeof(DD))))), 0)
			s1 := &DD
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//BitMap Inodo
		for i := 0; i < (5 * CantidadEstructuras); i++ {
			File.Seek(int64(SB.StartBmInodos+int32(i)), 0)
			var cero = '0'
			s1 := &cero
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//Inodo
		for i := 0; i < (5 * CantidadEstructuras); i++ {
			File.Seek(int64(SB.StartInodos+int32((i*int(unsafe.Sizeof(Inodo))))), 0)
			s1 := &Inodo
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//BitMap Bloque
		for i := 0; i < (20 * CantidadEstructuras); i++ {
			File.Seek(int64(SB.StartBmBloques+int32(i)), 0)
			var cero = '0'
			s1 := &cero
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//Bloque
		for i := 0; i < (20 * CantidadEstructuras); i++ {
			File.Seek(int64(SB.StartBloques+int32((i*int(unsafe.Sizeof(Bloque))))), 0)
			s1 := &Bloque
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		//Bitacora
		for i := 0; i < CantidadEstructuras; i++ {
			File.Seek(int64(SB.StartLog+int32((i*int(unsafe.Sizeof(Bitacora))))), 0)
			s1 := &Bitacora
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s1)
			File.Write(binario2.Bytes())
		}

		// Guardar Bitacora
		File.Seek(int64(SB.StartLog), 0)
		s1 := &Bitacora
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s1)
		File.Write(binario2.Bytes())

		File.Close()

		/*
		 * SE CREA EL AVD QUE REPRESENTA AL ROOT
		 */
		CrearRoot("/", id, 0) //TODO : Verficar cuando vandar true de la bitacora

		File = getFile(pathD)
		SuperBlock := readSuperBloque(File, int64(PartStart))
		DDroot := readDetalleDirectorio(File, int64(SuperBlock.StartDetalleDirectorio))
		Ruta := strings.Split("user.txt", "/")
		/*
		 *	MANDANDO A CREAR EL USER.TXT
		 */
		CrearArchivo(DDroot, Ruta, 0, pathD, SuperBlock, 0, "1,G,root\n1,U,root,123\n2,G,Usuarios")

		fmt.Println("----------------------------------------------------")
		fmt.Println("-       Formateo LWH realizado correctamente       -")
		fmt.Println("----------------------------------------------------")

	} else {
		ErrorMessage("[MKFS] -> La particion no se encuentra montada")
	}
}

//ComandoMKDIR is...
func ComandoMKDIR(id string, path string, p bool) {

	RutaDisco := listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {

		File := getFile(RutaDisco)
		PartStart := listaParticiones.GetPartStart(id)

		SuperBloque := readSuperBloque(File, int64(PartStart))

		Root := readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

		Rutas := strings.Split(path, "/")

		File.Close()

		MKDIR(Root, Rutas, RutaDisco, SuperBloque, 0, p)
		SuccessMessage("[MKDIR] -> Carpeta creada correctamente")

		/*
		 *	CREAR LA BITACORA
		 */
		File = getFile(RutaDisco)
		Bitacora := InicializarBitacora(Bitacora{})
		dt := time.Now()
		fecha := dt.Format("01-02-2020 15:04:05")
		copy(Bitacora.Fecha[:], fecha)
		copy(Bitacora.Nombre[:], path)
		copy(Bitacora.TipoOp[:], "mkdir")
		Bitacora.Tipo = 0
		File.Seek(int64(SuperBloque.StartLog+(SuperBloque.FirstFreeBitacora*int32(unsafe.Sizeof(Bitacora)))), 0)
		WriteBitacora(File, Bitacora)
		SuperBloque.FirstFreeBitacora++
		File.Seek(int64(PartStart), 0)
		WriteSB(File, SuperBloque)

	} else {
		ErrorMessage("[MKDIR] -> No existe ningun path asociado a ese id")
	}

}

//CrearRoot is...
func CrearRoot(Ruta string, id string, p int) {
	pathD := listaParticiones.GetDireccion(id)     //Obtenemos la direccion del disco
	SB := SB{}                                     // Instanciamos un SuperBloque
	File := getFile(pathD)                         //Obtenemos el file que contiene el disco
	PartStart := listaParticiones.GetPartStart(id) // Obtenemos el partStart

	SB = readSuperBloque(File, int64(PartStart)) //Leemos el superbloque

	if Ruta == "/" { //Aca creamos el AVD que hace referencia al root
		// y es el primero que se crea al hacer el formateo
		//Se crea una variable carpeta de tipo AVD que hace referencia al root
		var CarpetaRoot Arbol = InicializarAVD(Arbol{})
		//Se inicializan los valores del AVD
		copy(CarpetaRoot.AVDNombreDirectorio[:], "root")
		dt := time.Now()
		fecha := dt.Format("01-02-2006 15:04:05")
		copy(CarpetaRoot.AVDFechaCreacion[:], fecha)
		CarpetaRoot.DetalleDirectorio = SB.FirstFreeDd

		//Nos posicionamos en el start del AVD del superbloque
		File.Seek(int64(SB.StartArbolDirectorio), 0)
		//Escribimos la carpeta root
		s1 := &CarpetaRoot
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s1)
		File.Write(binario2.Bytes())

		//Escribimos 1 en el bitmap para indicar que esta ocupado
		File.Seek(int64(SB.StartBmArbolDirectorio), 0)
		var uno byte = '1'
		s2 := &uno
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s2)
		File.Write(binario3.Bytes())

		SB.ArbolVirtualFree--
		SB.FirstFreeAvd++

		/*
		 * INICIALIZAMOS EL DETALLE DIRECTORIO EN EL AVD
		 */

		//Escribiendo el detalle directorio
		Detalle := InicializarDD(DetalleDirectorio{})
		File.Seek(int64(SB.StartDetalleDirectorio), 0)
		s3 := &Detalle
		var binario4 bytes.Buffer
		binary.Write(&binario4, binary.BigEndian, s3)
		File.Write(binario4.Bytes())

		//Escribimos 1 en el bitmap para indicar que esta ocupado
		File.Seek(int64(SB.StartBmDetalleDirectorio), 0)
		var uno1 byte = '1'
		s4 := &uno1
		var binario5 bytes.Buffer
		binary.Write(&binario5, binary.BigEndian, s4)
		File.Write(binario5.Bytes())

		SB.DetalleDirectorioFree--
		SB.FirstFreeDd++

		/*
		 * REESCRIBIMOS EL SUPERBLOQUE EN EL PARTSART
		 */
		File.Seek(int64(PartStart), 0)
		s5 := &SB
		var binario6 bytes.Buffer
		binary.Write(&binario6, binary.BigEndian, s5)
		File.Write(binario6.Bytes())

		//Cerramos el archivo
		File.Close()

	}
}

//MKDIR is...
func MKDIR(AVD Arbol, paths []string, RutaDisco string, SuperBloque SB, puntero int, p bool) {
	// Se leer el archivo que contiene al disco
	if !VerificarRuta(RutaDisco) {
		ErrorMessage("[MKDIR] -> No existe ningun disco en la ruta del id")
		return
	}
	File := getFile(RutaDisco)

	if len(paths) == 0 {
		return
	}

	var apuntadorAVD int32 = 0

	// Se recorren las 6 posiciones de los subdirectorios del AVD
	for i := 0; i < 6; i++ {

		// Se obtiene el apuntador
		apuntador := AVD.Subirectorios[i]

		if apuntador == -1 { // Si el apuntador es -1 significa que esta vacia la posicion

			// Se crear una variable tipo AVD que hace referencia a la carpeta que se creara

			if p || len(paths) == 1 {
				/*Deja crear las carpetas*/
			} else {
				return
			}
			/*
			 * CARPETA A CREAR
			 */
			Carpeta := InicializarAVD(Arbol{})
			dt := time.Now()
			fecha := dt.Format("01-02-2006 15:04:05")
			copy(Carpeta.AVDNombreDirectorio[:], paths[0]) //Le seteamos el nombre a la carpeta
			copy(Carpeta.AVDFechaCreacion[:], fecha)       // Le seteamos la fecha

			//El detalle directorio de la carpeta es la primera posicion libre del detalle directorio del superbloque
			Carpeta.DetalleDirectorio = SuperBloque.FirstFreeDd
			// Guardamos el puntero del DD
			PosicionDetalleDirectorio := SuperBloque.FirstFreeDd
			//El apuntador al subdirectorio es la primera posicion libre del AVd del superbloque
			AVD.Subirectorios[i] = SuperBloque.FirstFreeAvd

			//El apuntador AVD es la primera posicion libre del AVd del superbloque
			apuntadorAVD = SuperBloque.FirstFreeAvd

			//Formula para saber en que posicion escribir el AVD
			Posicion := SuperBloque.StartArbolDirectorio + (apuntadorAVD * int32(unsafe.Sizeof(Carpeta)))

			/*
			 * Escribimos el AVD de la carpeta que deseamos crear
			 */
			File.Seek(int64(Posicion), 0)
			WriteAVD(File, Carpeta)

			//Escribimos el 1 en el bitmap para representar que esta lleno
			File.Seek(int64(SuperBloque.StartBmArbolDirectorio+apuntadorAVD), 0)
			WriteOne(File, '1')

			SuperBloque.ArbolVirtualFree--
			SuperBloque.FirstFreeAvd++

			/*
			 * Reescribimos el AVD que recibimos como parametro en el metodo
			 */
			File.Seek(int64(SuperBloque.StartArbolDirectorio+int32((puntero*int(unsafe.Sizeof(AVD))))), 0)
			WriteAVD(File, AVD)

			File.Seek(int64(SuperBloque.StartBmDetalleDirectorio+SuperBloque.FirstFreeDd), 0)
			WriteOne(File, '1')

			/*
			 * INICIALIZAMOS EL DETALLE DIRECTORIO EN EL AVD
			 */

			Detalle := InicializarDD(DetalleDirectorio{})
			File.Seek(int64(SuperBloque.StartDetalleDirectorio+(PosicionDetalleDirectorio*int32(unsafe.Sizeof(Detalle)))), 0)
			WriteDD(File, Detalle)

			//Escribimos 1 en el bitmap para indicar que esta ocupado
			File.Seek(int64(SuperBloque.StartBmDetalleDirectorio), 0)
			WriteOne(File, '1')

			SuperBloque.DetalleDirectorioFree--
			SuperBloque.FirstFreeDd++

			/*
			 *	Reescribimos el superbloque
			 */
			var PosicionSB int32 = SuperBloque.StartBmArbolDirectorio - int32(unsafe.Sizeof(SuperBloque))
			File.Seek(int64(PosicionSB), 0)
			WriteSB(File, SuperBloque)

			//Eliminamos del vector de que contiene los nombres de las carpetas , el nombre de la carpeta que acabamos de crear
			paths = paths[1:]
			//Cerramos el archivo
			File.Close()
			/*
			 * llamamos recursivamente a este metodo para crear todos las carpetas y se detiene hasta que
			 * el vector de los nombres de las carpetas este vacio
			 */
			MKDIR(Carpeta, paths, RutaDisco, SuperBloque, int(apuntadorAVD), p)
			return

		}
		/*
		 * Esta parte significa que la posicion no esta vacia y que ya existe una carpeta creada
		 */
		var CarpetaHija Arbol
		//Leemos la carpeta que esta creada en esa posicion
		CarpetaHija = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio+(apuntador*int32(unsafe.Sizeof(CarpetaHija)))))
		var nameByte [16]byte
		copy(nameByte[:], paths[0])
		//Comparamos si el nombre es igual
		if bytes.Compare(nameByte[:], CarpetaHija.AVDNombreDirectorio[:]) == 0 { // Si es igual significa que la carpeta ya estaba creada
			//Se saca del vector
			paths = paths[1:]
			//Se cierra el archivo
			File.Close()
			//Se llama recursivamente para crear las demas carpetas del vector
			MKDIR(CarpetaHija, paths, RutaDisco, SuperBloque, int(apuntador), p)
			return
		}
	}

	/*
	 *  APUNTADORES INDIRECTOS
	 *  Hace referencia al AVD que se crea si un AVD ya esta lleno , entonces se crea una copia
	 */
	// El valor del apuntador es el siguiente o sea la copia
	apuntador := AVD.VirtualDirectorio

	if apuntador == -1 { //Significa que no hay una copia aun

		//Creamos el AVD que hace referencia a la copia
		CopiaAVD := InicializarAVD(Arbol{})
		//Le seteamos el nombre del AVD anterior
		Nombre := string(AVD.AVDNombreDirectorio[:])
		Nombre = strings.Replace(Nombre, "\x00", "", -1)
		copy(CopiaAVD.AVDNombreDirectorio[:], Nombre)
		//Obtenemos la fecha actual
		dt := time.Now()
		fecha := dt.Format("01-02-2006 15:04:05")
		//Le setemos la fecha
		copy(CopiaAVD.AVDFechaCreacion[:], fecha)

		CopiaAVD.DetalleDirectorio = SuperBloque.FirstFreeDd
		AVD.VirtualDirectorio = SuperBloque.FirstFreeAvd

		apuntadorAVD = SuperBloque.FirstFreeAvd

		var Posicion int = int(SuperBloque.StartArbolDirectorio + (SuperBloque.FirstFreeAvd * int32(unsafe.Sizeof(CopiaAVD))))

		File.Seek(int64(Posicion), 0)
		WriteAVD(File, CopiaAVD)

		File.Seek(int64(SuperBloque.StartBmArbolDirectorio+SuperBloque.FirstFreeAvd), 0)
		WriteOne(File, '1')

		SuperBloque.ArbolVirtualFree--
		SuperBloque.FirstFreeAvd++

		//Reescribo avd que recibo como parametro
		File.Seek(int64(SuperBloque.StartArbolDirectorio+int32((puntero*int(unsafe.Sizeof(AVD))))), 0)
		WriteAVD(File, AVD)

		File.Seek(int64(SuperBloque.StartBmDetalleDirectorio+SuperBloque.FirstFreeDd), 0)
		WriteOne(File, '1')

		/*
		 * INICIALIZAMOS EL DETALLE DIRECTORIO EN EL AVD
		 */

		Detalle := InicializarDD(DetalleDirectorio{})
		File.Seek(int64(SuperBloque.StartDetalleDirectorio+(SuperBloque.FirstFreeDd*int32(unsafe.Sizeof(Detalle)))), 0)
		WriteDD(File, Detalle)

		//Escribimos 1 en el bitmap para indicar que esta ocupado
		File.Seek(int64(SuperBloque.StartBmDetalleDirectorio), 0)
		WriteOne(File, '1')

		SuperBloque.DetalleDirectorioFree--
		SuperBloque.FirstFreeDd++

		/*
		 *	Reescribimos el superbloque
		 */
		var PosicionSB int32 = SuperBloque.StartBmArbolDirectorio - int32(unsafe.Sizeof(SuperBloque))
		File.Seek(int64(PosicionSB), 0)
		WriteSB(File, SuperBloque)

		//Cerramos el archivo
		File.Close()

		// Llamamos al metodo recursivamente para que cree la carpeta en la copia del AVD y le mandamos la copia
		MKDIR(CopiaAVD, paths, RutaDisco, SuperBloque, int(apuntadorAVD), p)
		return

	}
	//Significa que ya existe una copia de ese AVD
	var CarpetaCopia Arbol
	CarpetaCopia = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio+(apuntador*int32(unsafe.Sizeof(CarpetaCopia)))))
	//Cerramos el archivo
	File.Close()

	//Llamamos recursivamente al metodo MKDIR para que repita todo el proceso en la copia del AVD
	MKDIR(CarpetaCopia, paths, RutaDisco, SuperBloque, int(apuntador), p)
	return

}

//MKFILE is...
func MKFILE(id string, path string, p bool, size int, count string) {

	if isLogged {

		if path == "" {
			return
		}

		PathDisco := listaParticiones.GetDireccion(id)
		if PathDisco != "null" {

			if VerificarRuta(PathDisco) {
				File := getFile(PathDisco)
				//PartSize := listaParticiones.GetPartSize(id)
				PartStart := listaParticiones.GetPartStart(id)
				//PartName := listaParticiones.GetPartName(id)

				//Leemos el SuperBloque
				var SuperBloque SB
				SuperBloque = readSuperBloque(File, int64(PartStart))

				Rutas := strings.Split(path, "/") // TODO : Eliminar si vienen / al inicio y final
				RutasMKFILE := strings.Split(path, "/")

				var Root Arbol
				Root = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

				//Eliminamos el espacio vacio y el nombre del archivo home/jose
				Rutas = Rutas[1:]
				Rutas = Rutas[:len(Rutas)-1]
				RutasMKFILE = RutasMKFILE[1:]

				File.Close()

				//Mandamos a crear las carpetas si no estan creadas
				MKDIR(Root, Rutas, PathDisco, SuperBloque, 0, p)

				File1 := getFile(PathDisco)
				SuperB := readSuperBloque(File1, int64(PartStart))

				Raiz := readArbolVirtualDirectorio(File1, int64(SuperB.StartArbolDirectorio))
				DetalleRoot := readDetalleDirectorio(File1, int64(SuperB.StartDetalleDirectorio))

				File1.Close()

				if len(RutasMKFILE) == 1 {
					CrearArchivo(DetalleRoot, RutasMKFILE, 0, PathDisco, SuperB, size, count)
					return
				}

				RecorrerArbol(Raiz, RutasMKFILE, PathDisco, SuperB, size, count)

				SuccessMessage("[MKFILE] -> Archivo creado correctamente")

			} else {
				ErrorMessage("[MKFILE] -> No se encuentra ningun disco en esa ruta")
			}

		} else {
			ErrorMessage("[MKFILE] -> No hay ninguna particion montada con ese id")
		}
	} else {
		ErrorMessage("[MKFILE] -> Para ejecutar este comando tienes que estar loggeado")
	}
}

//RecorrerArbol is...
func RecorrerArbol(root Arbol, Rutas []string, PathDisco string, Superbloque SB, size int, count string) {

	if len(Rutas) == 0 {
		return
	}

	if VerificarRuta(PathDisco) {

		File := getFile(PathDisco) //Leemos el disco
		var Apuntador int32 = 0

		//recorremos los 6 subdirectorios del AVD
		for i := 0; i < 6; i++ {
			Apuntador = root.Subirectorios[i]
			var CarpetaHija Arbol
			CarpetaHija = readArbolVirtualDirectorio(File, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CarpetaHija)))))

			//Verificamos el nombre de la carpeta
			var nameCarpetaByte [16]byte
			copy(nameCarpetaByte[:], Rutas[0])

			if bytes.Compare(CarpetaHija.AVDNombreDirectorio[:], nameCarpetaByte[:]) == 0 {
				Rutas = Rutas[1:]
				if len(Rutas) == 1 {
					Apuntador = CarpetaHija.DetalleDirectorio
					var Archivos DetalleDirectorio
					Archivos = readDetalleDirectorio(File, int64(Superbloque.StartDetalleDirectorio+(Apuntador*int32(unsafe.Sizeof(Archivos)))))
					//Cerramos el archivo
					File.Close()
					CrearArchivo(Archivos, Rutas, int(Apuntador), PathDisco, Superbloque, size, count)
					return
				}
				File.Close()
				RecorrerArbol(CarpetaHija, Rutas, PathDisco, Superbloque, size, count)
				return
			}

		}
		// Buscamos en la copia de la carpeta
		Apuntador = root.VirtualDirectorio
		var CopiaCarpeta Arbol
		CopiaCarpeta = readArbolVirtualDirectorio(File, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CopiaCarpeta)))))
		File.Close()
		/*
		 * LLAMAMOS EL METODO RECURSIVAMENTE
		 */
		RecorrerArbol(CopiaCarpeta, Rutas, PathDisco, Superbloque, size, count)
		return
	}
	ErrorMessage("[MKFILE] -> No hay ningun disco en la ruta indicada")

}

//CrearArchivo is...
func CrearArchivo(Archivo DetalleDirectorio, Rutas []string, apuntador int, RutaDisco string, SuperB SB, size int, count string) {

	if VerificarRuta(RutaDisco) {
		//Recorreremos el array de Files del detalle directorio
		File := getFile(RutaDisco)
		for i := 0; i < 5; i++ {
			Puntero := Archivo.DDArrayFiles[i].DDFileApInodo
			if Puntero == -1 { //Quiere decir que esta disponible
				/*
				 * CREAMOS UN STRUCT DE FILE QUE REPRESENTA EL ARCHIVO QUE DESEAMOS CREAR
				 */
				var FileCrear FileStruct
				dt := time.Now()
				fecha := dt.Format("01-02-2006 15:04:05")
				// Inicializamos los atributos
				copy(FileCrear.DDFileDateCreacion[:], fecha)
				copy(FileCrear.DDFileDateModificacion[:], fecha)
				copy(FileCrear.DDFileNombre[:], Rutas[0])
				FileCrear.DDFileApInodo = SuperB.FirstFreeInodo
				/*
				 * ASIGNAMOS EL FILE QUE ACABAMOS DE CREAR AL DETALLE DIRECTORIO
				 */
				Archivo.DDArrayFiles[i] = FileCrear

				/*
				 * REESCRIBIMOS EL DETALLE DIRECTORIO CON LAS MODIFICACIONES REALIZADAS
				 */
				File.Seek(int64(SuperB.StartDetalleDirectorio+(int32(apuntador)*int32(unsafe.Sizeof(Archivo)))), 0)
				WriteDD(File, Archivo)

				//Obtenemos el valor del apuntador al inodo
				ApuntadorInodo := SuperB.FirstFreeInodo
				/*
				 *  CREAMOS EL INODO Y SETEAMOS SUS VALORES
				 */
				Inodo := InicializarInodo(TablaInodo{})
				Inodo.ICountInodo = ApuntadorInodo
				Inodo.ISizeArchivo = int32(size)
				/*
				 *	CALCULAMOS LA CANTIDAD
				 */
				var nBloques float64
				var num float64 = 25

				if count != "" && size == 0 { // EL CONTENIDO DEFINE LA CANTIDAD DE BLOQUES
					nBloques = float64(len(count)) / num
					if nBloques-math.Floor(nBloques) != 0 {
						nBloques = nBloques + 1
					}
				} else if count == "" && size != 0 { // EL SIZE DEFINE LA CANTIDAD DE BLOQUES
					nBloques = float64(size) / num
					if nBloques-math.Floor(nBloques) != 0 {
						nBloques = nBloques + 1
					}
				} else if count != "" && size != 0 { // SI LOS 2 VIENEN
					if len(count) > size { // EL CONTENIDO DEFINE LA CANTIDAD DE BLOQUES
						nBloques = float64(len(count)) / num
						if nBloques-math.Floor(nBloques) != 0 {
							nBloques = nBloques + 1
						}
					} else { // EL SIZE DEFINE LA CANTIDAD DE BLOQUES
						nBloques = float64(size) / num
						if nBloques-math.Floor(nBloques) != 0 {
							nBloques = nBloques + 1
						}
					}
				}
				if nBloques <= 4 {
					Inodo.ICountBloquesAsignados = int32(nBloques)
				} else {
					Inodo.ICountBloquesAsignados = 4
				}
				SuperB.InodosFree--
				SuperB.FirstFreeInodo++
				/*
				 *  ESCRIBIMOS EL 1 EN EL BITMAP DE INODO
				 */
				File.Seek(int64(SuperB.StartBmInodos+Inodo.ICountInodo), 0)
				var uno byte = '1'
				s := &uno
				var binario1 bytes.Buffer
				binary.Write(&binario1, binary.BigEndian, s)
				File.Write(binario1.Bytes())
				/*
				 *	REESCRIBIMOS LE SUPERBLOQUE
				 */
				File.Seek(int64(SuperB.StartBmArbolDirectorio-int32(unsafe.Sizeof(SuperB))), 0)
				s2 := &SuperB
				var binario2 bytes.Buffer
				binary.Write(&binario2, binary.BigEndian, s2)
				File.Write(binario2.Bytes())
				//Cerramos el archivo
				File.Close()

				CrearInodo(Inodo, Rutas[0], RutaDisco, SuperB, count, int(nBloques))
				return
			}
			/*
			 *	LEEMOS EL INODO PARA SUSTITUIR LA DATA
			 */
			var nameDirec [16]byte
			copy(nameDirec[:], Rutas[0])

			if bytes.Compare(nameDirec[:], Archivo.DDArrayFiles[i].DDFileNombre[:]) == 0 {
				/*var InodoB TablaInodo
				InodoB = readInodo(File, int64(SuperB.StartInodos+(Puntero*int32(unsafe.Sizeof(InodoB)))))
				SustituirData(InodoB, SuperB, File, count)*/
				ErrorMessage("[MKFILE] -> Ya existe un file con el mismo nombre")
				return
			}

		}
		/*
		 * NOS MOVEMOS A BUSCAR EN LA COPIA
		 */

		var ApuntadorCopia int32 = Archivo.DDApDetalleDirectorio
		if ApuntadorCopia == -1 {
			/*
			 *	TENEMOS QUE CREAR UNA COPIA DEL DETALLE DIRECTORIO
			 */
			NuevoDetalleDirectorio := InicializarDD(DetalleDirectorio{})
			PosNuevoDD := SuperB.StartDetalleDirectorio + (SuperB.FirstFreeDd * int32(unsafe.Sizeof(NuevoDetalleDirectorio)))
			ApuntadorCopia = SuperB.FirstFreeDd
			Archivo.DDApDetalleDirectorio = SuperB.FirstFreeDd

			/*
			 * ESCRIBIMOS EL NUEVO DETALLE DIRECTORIO (COPIA)
			 */
			File.Seek(int64(PosNuevoDD), 0)
			s3 := &NuevoDetalleDirectorio
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s3)
			File.Write(binario3.Bytes())

			/*
			 *  ESCRIBIMOS EL 1 EN EL BITMAP DEL DETALLE DIRECTORIO
			 */
			File.Seek(int64(SuperB.StartBmDetalleDirectorio+SuperB.FirstFreeDd), 0)
			var unoo byte = '1'
			s4 := &unoo
			var binario4 bytes.Buffer
			binary.Write(&binario4, binary.BigEndian, s4)
			File.Write(binario4.Bytes())

			SuperB.DetalleDirectorioFree--
			SuperB.FirstFreeDd++

			/*
			 *  REESCRIBIMOS EL DETALLE DIRECTORIO PADRE PARA APLICARLE LOS CAMBIOS
			 */
			File.Seek(int64(SuperB.StartDetalleDirectorio+int32(apuntador*int(unsafe.Sizeof(Archivo)))), 0)
			s5 := &Archivo
			var binario5 bytes.Buffer
			binary.Write(&binario5, binary.BigEndian, s5)
			File.Write(binario5.Bytes())

			/*
			 * REESCRIBIMOS EL SUPERBLOQUE PARA APLICARLE LOS CAMBIOS
			 */
			File.Seek(int64(SuperB.StartBmArbolDirectorio-int32(unsafe.Sizeof(SuperB))), 0)
			s6 := &SuperB
			var binario6 bytes.Buffer
			binary.Write(&binario6, binary.BigEndian, s6)
			File.Write(binario6.Bytes())

			//Cerramos el archivo
			File.Close()

			/*
			 * YA QUE CREAMOS EL DD COPIA MANDAMOS A LLAMAR AL METODO CREARARCHIVO RECURSIVAMENTE
			 * PERO LE MANDAMOS COMO DD LA COPIA
			 */
			CrearArchivo(NuevoDetalleDirectorio, Rutas, int(ApuntadorCopia), RutaDisco, SuperB, size, count)
			return
		}
		var DetalleDirectorioCopia DetalleDirectorio
		DetalleDirectorioCopia = readDetalleDirectorio(File, int64(SuperB.StartDetalleDirectorio+(ApuntadorCopia*int32(unsafe.Sizeof(DetalleDirectorioCopia)))))
		//Cerramos el archivo
		File.Close()
		/*
		 * LLAMAMOS RECURSIVAMENTE AL METODO CREAR ARCHIVO PERO LE MANDAMOS COMO DD LA COPIA
		 */
		CrearArchivo(DetalleDirectorioCopia, Rutas, int(ApuntadorCopia), RutaDisco, SuperB, size, count)
		return

	}
	ErrorMessage("[MKFILE] -> No existe ningun disco en esa ruta")

}

//SustituirData is...
func SustituirData(Inodo TablaInodo, SuperBloque SB, file *os.File, cont string) {

	for i := 0; i < 4; i++ {
		ApuntadorBloque := Inodo.IArrayBloques[i]

		if ApuntadorBloque != -1 {
			var Block Bloque
			Block = readBloque(file, int64(SuperBloque.StartBloques+(ApuntadorBloque*int32(unsafe.Sizeof(Block)))))
			copy(Block.Texto[:], cont)
			file.Seek(int64(SuperBloque.StartBloques+(ApuntadorBloque*int32(unsafe.Sizeof(Block)))), 0)
			WriteBloque(file, Block)
		}
	}

}

//CrearInodo is...
func CrearInodo(Inodo TablaInodo, NombreArchivo string, RutaDisco string, SuperBloque SB, contenido string, Bloques int) {

	if VerificarRuta(RutaDisco) {
		File := getFile(RutaDisco)

		for i := 0; i < int(Inodo.ICountBloquesAsignados); i++ {
			//Verificamos si el contenido viene o no
			if contenido == "" {
				contenido = "abcdefghijklmnopqrstuvwxy"
			}
			/*
			 * CREAMOS EL BLOQUE
			 */
			var BloqueDatos Bloque
			Rune := []rune(contenido)
			SubString := string(Rune[0:25])
			copy(BloqueDatos.Texto[:], SubString)
			//Verificamos el valor del contenido a ingresar
			if len(contenido) > 25 {
				SubStringCont := string(Rune[25:len(contenido)])
				contenido = SubStringCont
			} else {
				contenido = ""
			}
			Inodo.IArrayBloques[i] = SuperBloque.FirstFreeBloque
			/*
			 *  ESCRIBIMOS EL BLOQUE EN EL FILE
			 */
			File.Seek(int64(SuperBloque.StartBloques+(SuperBloque.FirstFreeBloque)*int32(unsafe.Sizeof(BloqueDatos))), 0)
			s := &BloqueDatos
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, s)
			File.Write(binario.Bytes())
			/*
			 * PONEMOS UN 1 EN EL BITMAP DE BLOQUES
			 */
			File.Seek(int64(SuperBloque.StartBmBloques+SuperBloque.FirstFreeBloque), 0)
			var uno byte = '1'
			s1 := &uno
			var binario1 bytes.Buffer
			binary.Write(&binario1, binary.BigEndian, s1)
			File.Write(binario1.Bytes())
			/*
			 * REESCRIBIMOS EL INODO
			 */
			File.Seek(int64(SuperBloque.StartInodos+(Inodo.ICountInodo*int32(unsafe.Sizeof(Inodo)))), 0)
			s2 := &Inodo
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s2)
			File.Write(binario2.Bytes())

			SuperBloque.BloquesFree--
			SuperBloque.FirstFreeBloque++
			/*
			 * REESCRIBIMOS EL SUPERBLOQUE CON LOS CAMBIOS
			 */
			File.Seek(int64(SuperBloque.StartBmArbolDirectorio-int32(unsafe.Sizeof(SuperBloque))), 0)
			s3 := &SuperBloque
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s3)
			File.Write(binario3.Bytes())
		}
		/*
		 *  VERIFICAMOS SI LOS BLOQUES SON MAYOR A 4
		 */
		if Bloques > 4 {
			Bloques = Bloques - 4
			InodoCopia := InicializarInodo(TablaInodo{})
			InodoCopia.ICountInodo = SuperBloque.FirstFreeInodo
			Inodo.IApIndirecto = InodoCopia.ICountInodo
			InodoCopia.ISizeArchivo = Inodo.ISizeArchivo

			SuperBloque.InodosFree--
			SuperBloque.FirstFreeInodo++
			/*
			 *	REESCRIBIMOS EL INODOS CON LOS CAMBIOS
			 */
			File.Seek(int64(SuperBloque.StartInodos+(Inodo.ICountInodo*int32(unsafe.Sizeof(InodoCopia)))), 0)
			s4 := &Inodo
			var binario4 bytes.Buffer
			binary.Write(&binario4, binary.BigEndian, s4)
			File.Write(binario4.Bytes())
			/*
			 *	ESCRIBIMOS UN 1 EN EL BITMAP DEL INODO
			 */
			File.Seek(int64(SuperBloque.StartBmInodos+InodoCopia.ICountInodo), 0)
			var uno byte = '1'
			s5 := &uno
			var binario5 bytes.Buffer
			binary.Write(&binario5, binary.BigEndian, s5)
			File.Write(binario5.Bytes())
			/*
			 *	REESCRIBIMOS EL SUPERBLOQUE CON LOS CAMBIOS
			 */
			File.Seek(int64(SuperBloque.StartBmArbolDirectorio-int32(unsafe.Sizeof(SuperBloque))), 0)
			s6 := &SuperBloque
			var binario6 bytes.Buffer
			binary.Write(&binario6, binary.BigEndian, s6)
			File.Write(binario6.Bytes())
			if Bloques <= 4 {
				InodoCopia.ICountBloquesAsignados = int32(Bloques)
			} else {
				InodoCopia.ICountBloquesAsignados = 4
			}
			File.Close()
			CrearInodo(InodoCopia, NombreArchivo, RutaDisco, SuperBloque, contenido, Bloques)
			return
		}
		File.Close()
	} else {
		ErrorMessage("[MKFILE] -> No existe ningun disco en esta ruta")
	}

}

//ComandoEditFile is...
func ComandoEditFile() {

}

/*
 *  METODOS PARA INICIALIZAR LOS STRUCTS DE LA FASE 2
 */

//InicializarAVD is...
func InicializarAVD(Arbol Arbol) Arbol {
	for i := 0; i < 6; i++ {
		Arbol.Subirectorios[i] = -1
	}
	Arbol.AVDProper = -1
	Arbol.VirtualDirectorio = -1
	Arbol.DetalleDirectorio = 1
	return Arbol
}

//InicializarDD is...
func InicializarDD(DetalleD DetalleDirectorio) DetalleDirectorio {
	DetalleD.DDApDetalleDirectorio = -1
	for i := 0; i < 5; i++ {
		DetalleD.DDArrayFiles[i].DDFileApInodo = -1
	}
	return DetalleD
}

//InicializarInodo is...
func InicializarInodo(Inodo TablaInodo) TablaInodo {
	for i := 0; i < 4; i++ {
		Inodo.IArrayBloques[i] = -1
	}
	Inodo.IApIndirecto = -1
	Inodo.IIDProper = -1
	return Inodo
}

//InicializarBloque is...
func InicializarBloque(Bloque Bloque) Bloque {
	return Bloque
}

//InicializarBitacora is...
func InicializarBitacora(Bitacora Bitacora) Bitacora {
	Bitacora.Tipo = -1
	return Bitacora
}

/*
 *  Manejo de la sesion
 */

//Login is ...
func Login(user string, password string, id string) {
	var RutaDisco string = listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {
		File := getFile(RutaDisco)
		PartStart := listaParticiones.GetPartStart(id)
		SuperBloque := readSuperBloque(File, int64(PartStart))
		DetallDirectorioRoot := readDetalleDirectorio(File, int64(SuperBloque.StartDetalleDirectorio))
		ObtenerDataUserTXT(DetallDirectorioRoot, File, SuperBloque, user, password)

	} else {
		ErrorMessage("[LOGIN] -> No se encuentra ninguna particion montada con ese id")
	}
}

//ObtenerDataUserTXT is...
func ObtenerDataUserTXT(DDroot DetalleDirectorio, File *os.File, SuperB SB, user string, pass string) {
	//Obtenemos el apuntador al Inodo
	ApuntadorInodo := DDroot.DDArrayFiles[0].DDFileApInodo
	var InodoB TablaInodo
	//Leemos el Inodo
	InodoB = readInodo(File, int64(SuperB.StartInodos+(ApuntadorInodo*int32(unsafe.Sizeof(InodoB)))))
	//Obtenemos el apuntador al bloque
	var Split []string
	var ContenidoUserTxt string
	for i := 0; i < 4; i++ {
		ApuntadorBloque := InodoB.IArrayBloques[i]
		if ApuntadorBloque != -1 {
			var Block Bloque
			//Leemos el bloque
			Block = readBloque(File, int64(SuperB.StartBloques+(ApuntadorBloque*int32(unsafe.Sizeof(Block)))))
			ContenidoUserTxt += string(Block.Texto[:])
		}
	}
	ContenidoUserTxt = strings.Replace(ContenidoUserTxt, "\x00", "", -1)
	Split = strings.Split(ContenidoUserTxt, "\n")

	var Datos []string
	var log bool
	for i := 0; i < len(Split); i++ {
		Datos = strings.Split(Split[i], ",")
		if Datos[1] == "U" {
			log = VerificarDatos(Datos[2], Datos[3], user, pass, Datos[0])
			if log {
				SuccessMessage("[LOGIN] -> Loggeado Correctamente")
				break
			}
		}
	}
	if !log {
		ErrorMessage("[LOGIN] -> Datos Incorrectos")
	}
}

//VerificarDatos is...
func VerificarDatos(u string, p string, user string, pass string, UserID string) bool {

	if u == user && p == pass {
		isLogged = true
		copy(userLoggeado.UserName[:], user)
		copy(userLoggeado.PassWord[:], pass)
		i, _ := strconv.Atoi(UserID)
		userLoggeado.IDUser = int32(i)
		return true
	}

	return false
}

//Logout is...
func Logout() {
	// TODO : Hacer logout
}

/*
 *  M E N S A J E S
 */

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
