package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
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
	FechaCreacion [20]byte
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
	DateCreacion             [20]byte
	DateUltimoMontaje        [20]byte
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
	DDArrayFiles          [5]File
	DDApDetalleDirectorio int32
}

//File is...
type File struct {
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

//Sesion is...
type Sesion struct {
	IDUser        int32
	IDGrupo       int32
	InicioSupet   int32
	InicioJournal int32
	TipoSistema   int32
	Direccion     string
	Fit           byte
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

var sesionActual Sesion
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
		//Metodo para leer el struct MBR del Disco(archivo)
		//readMBR(path + name + ".dsk")
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
	/*fmt.Println("Part Start", m.PartStart)
	fmt.Println("Part Size", m.PartSize)
	fmt.Println("Part Fit", string(m.PartFit))
	fmt.Println("Part Name", string(m.PartName[:]))
	fmt.Println("Part Next", m.PartNext)
	fmt.Println("Part Status", string(m.PartStatus))
	fmt.Println("")*/
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

//BuscarParticion is...
func BuscarParticion(path string, nombre string) int {

	if VerificarRuta(path) {

		File := getFile(path)
		MBR := readMBR(File)

		for i := 0; i < 4; i++ {
			var nameByte [16]byte
			copy(nameByte[:], nombre)
			if MBR.Particion[i].PartStatus != '1' {
				if bytes.Compare(nameByte[:], MBR.Particion[i].PartName[:]) == 0 {
					return i
				}
			}
		}

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
		os.Create("graficaEBR.dot")
		graphDot := getFile("graficaEBR.dot")

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
				var status [3]byte
				if MB.Particion[i].PartStatus == '0' {
					copy(status[:], "0")
				} else if MB.Particion[i].PartStatus == '2' {
					copy(status[:], "2")
				}

				fmt.Fprintf(graphDot, "<tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>Particion%d</font></b></td></tr>\n", (i + 1))
				fmt.Fprintf(graphDot, "<tr>  <td>Status</td> <td>%s</td>  </tr>\n", string(status[:]))
				fmt.Fprintf(graphDot, "<tr>  <td>Type</td> <td>%c</td>  </tr>\n", MB.Particion[i].PartType)
				fmt.Fprintf(graphDot, "<tr>  <td>Fit</td> <td>%c</td>  </tr>\n", MB.Particion[i].PartFit)
				fmt.Fprintf(graphDot, "<tr>  <td>Start</td> <td>%d</td>  </tr>\n", MB.Particion[i].PartStart)
				fmt.Fprintf(graphDot, "<tr>  <td>Size</td> <td>%d</td>  </tr>\n", MB.Particion[i].PartSize)
				fmt.Fprintf(graphDot, "<tr>  <td>Name</td> <td>%s</td>  </tr>\n", string(MB.Particion[i].PartName[:]))
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
					var status [3]byte
					if extendedBoot.PartStatus == '0' {
						copy(status[:], "0")
					} else if extendedBoot.PartStatus == '2' {
						copy(status[:], "2")
					}

					fmt.Fprintf(graphDot, "<tr>  <td>Status</td> <td>%s</td>  </tr>\n", string(status[:]))
					fmt.Fprintf(graphDot, "<tr>  <td>Fit</td> <td>%c</td>  </tr>\n", extendedBoot.PartFit)
					fmt.Fprintf(graphDot, "<tr>  <td>Start</td> <td>%d</td>  </tr>\n", extendedBoot.PartStart)
					fmt.Fprintf(graphDot, "<tr>  <td>Size</td> <td>%d</td>  </tr>\n", extendedBoot.PartSize)
					fmt.Fprintf(graphDot, "<tr>  <td>Next</td> <td>%d</td>  </tr>\n", extendedBoot.PartNext)
					fmt.Fprintf(graphDot, "<tr>  <td>Name</td> <td>%s</td>  </tr>\n", string(extendedBoot.PartName[:]))
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
		exec.Command("dot", "-Tpng", "-o", "/home/jose/Escritorio/graficaEBR.png", "graficaEBR.dot").Output()
		SuccessMessage("[REP] -> Reporte del disco generado correctamente")

	}

}

//ReporteDisco is...
func ReporteDisco(direccion string, destino string, extension string) {

	var auxDir string = direccion

	if VerificarRuta(auxDir) {
		fp := getFile(auxDir)
		os.Create("graficaDisco.dot")
		graphDot := getFile("graficaDisco.dot")

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
		exec.Command("dot", "-Tpng", "-o", "/home/jose/Escritorio/grafica.png", "graficaDisco.dot").Output()

		SuccessMessage("[REP] -> Reporte del disco generado correctamente")
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

	os.Create("graficaSuperBloque.dot")
	graphDot := getFile("graficaSuperBloque.dot")

	//Empezamos a escribir en el archivo
	fmt.Fprintf(graphDot, "digraph G{ \n")
	fmt.Fprintf(graphDot, "node [shape=plaintext]\n")
	fmt.Fprintf(graphDot, "tbl[\nlabel=<\n")
	fmt.Fprintf(graphDot, "<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n")
	fmt.Fprintf(graphDot, " <tr ><td colspan='2' bgcolor= 'lightblue' ><b><font color='blue'>Super Bloque</font></b></td></tr>")
	fmt.Fprintf(graphDot, "<tr>  <td width='230'> <b>Atributo</b> </td> <td width='230'> <b>Valor</b> </td>  </tr>\n")

	fmt.Fprintf(graphDot, "<tr>  <td>NombreHD</td><td>%s</td>  </tr>\n", string(SB.NombreHD[:]))
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

/*
 *	C O M A N D O S   F A S E   2
 */

//MKFS is...
func MKFS(id string) {
	Formatear(id)
}

//Formatear is...
func Formatear(id string) {
	//TODO : Ver lo de la escritura , porque mueren las particiones logicas

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

		/*fmt.Println("SB", SBSize)
		fmt.Println("AVD", AVDSize)
		fmt.Println("DD", DDSize)
		fmt.Println("Inodo", InodoSize)
		fmt.Println("Bloque", BloqueSize)
		fmt.Println("Bitacora", BitacoraSize)*/

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
		//Magic Num
		SB.MagicNum = 201807431
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

		//CREAR RAIZ
		/*Carpeta folder;
		folder.makeDirectory("/", 0, id, montaje, false);*/
		/*
		 * SE CREA EL AVD QUE REPRESENTA AL ROOT
		 */
		CrearDirectorio("/", id, false, 0) //TODO : Verficar cuando vandar true de la bitacora
		fmt.Println("----------------------------------------------------")
		fmt.Println("-       Formateo LWH realizado correctamente       -")
		fmt.Println("----------------------------------------------------")

	} else {
		ErrorMessage("[MKFS] -> La particion no se encuentra montada")
	}
}

//CrearDirectorio is...
func CrearDirectorio(Ruta string, id string, bitacora bool, p int) {

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

	} else { //Si es una ruta diferente a root ('/')

		//Hacemos un split para obtener todos los nombres de las carpetas que se desean crear
		Rutas := strings.Split(Ruta, "/")
		Rutas = Rutas[1:]
		Rutas = Rutas[:len(Rutas)]
		// Creamos una variable llamada root que hace referencia al AVD del root que se crea en el formateo
		var Root Arbol
		// Se lee el root
		Root = readArbolVirtualDirectorio(File, int64(SB.StartArbolDirectorio))

		/* Mandamos esos valores para el metodo MKDIR que es el encargado de crear las carpetas
		le mandamos el root porque es la raiz y de ahi se desprenden todas las carpetas */
		MKDIR(Root, Rutas, Ruta, SB, 0)

		//Cerramos el archivo
		File.Close()
	}

	// Verificamos si es bitacora
	if !bitacora {
		dt := time.Now()
		fecha := dt.Format("01-02-2006 15:04:05")
		//Instanciamos la bitacora
		var log Bitacora
		copy(log.Fecha[:], fecha)
		log.Tipo = '0'
		copy(log.TipoOp[:], "mkdir")
		log.Size = int32(p)

		File2 := getFile(pathD)
		var BitacoraRaiz Bitacora
		BitacoraRaiz = readBitacora(File2, int64(SB.StartLog))
		sizeBitacora := int32(unsafe.Sizeof(log))
		//Lo escribimos en el first free de la bitacora.
		File2.Seek(int64(SB.StartLog+(BitacoraRaiz.Size*sizeBitacora)), 0) //TODO : verificar lo del size del bitacora raiz
		s := &log
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s)
		File2.Write(binario.Bytes())

		BitacoraRaiz.Size++
		File2.Seek(int64(SB.StartLog), 0)
		s1 := &log
		var binario1 bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s1)
		File2.Write(binario1.Bytes())

		File2.Close()
	}
}

//MKDIR is...
func MKDIR(AVD Arbol, paths []string, RutaDisco string, SuperBloque SB, Apuntador int) {
	// TODO : Hacer el mkdir

	// Se leer el archivo que contiene al disco
	File := getFile(RutaDisco)
	// Se declara un char para escribir en los bitmaps
	var uno byte = '1'

	// Esta variable hace referencia al apuntador actual en los subdirectorios
	var apuntador int32 = 0
	var apuntadorAVD int32 = 0

	// Se recorren las 6 posiciones de los subdirectorios del AVD
	for i := 0; i < 6; i++ {

		// Se obtiene el apuntador
		apuntador = AVD.Subirectorios[i]

		if apuntador == -1 { // Si el apuntador es -1 significa que esta vacia la posicion

			// Se crear una variable tipo AVD que hace referencia a la carpeta que se creara
			var Carpeta Arbol
			// Se inicializan los valores
			dt := time.Now()
			fecha := dt.Format("01-02-2006 15:04:05")
			copy(Carpeta.AVDNombreDirectorio[:], paths[0])
			copy(Carpeta.AVDFechaCreacion[:], fecha)

			//El detalle directorio de la carpeta es la primera posicion libre del detalle direcotiro del superbloque
			Carpeta.DetalleDirectorio = SuperBloque.FirstFreeDd
			//El apuntador al subdirectorio es la primera posicion libre del AVd del superbloque
			AVD.Subirectorios[i] = SuperBloque.FirstFreeAvd

			//El apuntador AVD es la primera posicion libre del AVd del superbloque
			apuntadorAVD = SuperBloque.FirstFreeAvd

			//Formula para saber en que posicion escribir
			Posicion := SuperBloque.StartArbolDirectorio + (SuperBloque.FirstFreeAvd * int32(unsafe.Sizeof(Carpeta)))

			/*
			 * Escribimos el AVD de la carpeta que deseamos crear
			 */
			File.Seek(int64(Posicion), 0)
			s1 := &Carpeta
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, s1)
			File.Write(binario.Bytes())

			//Escribimos el 1 en el bitmap para representar que esta lleno
			File.Seek(int64(SuperBloque.StartBmArbolDirectorio+SuperBloque.FirstFreeAvd), 0)
			s2 := &uno
			var binario1 bytes.Buffer
			binary.Write(&binario1, binary.BigEndian, s2)
			File.Write(binario1.Bytes())

			SuperBloque.ArbolVirtualFree--
			SuperBloque.FirstFreeAvd++

			/*
			 * Reescribimos el AVD que recibimos como parametro en el metodo
			 */
			File.Seek(int64(SuperBloque.StartArbolDirectorio+(apuntador*int32(unsafe.Sizeof(AVD)))), 0)
			s3 := &AVD
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s3)
			File.Write(binario2.Bytes())

			//Eliminamos del vector de que contiene los nombres de las carpetas , el nombre de la carpeta que acabamos de crear
			paths = paths[1:]
			//Cerramos el archivo
			File.Close()
			/*
			 * llamamos recursivamente a este metodo para crear todos las carpetas y se detiene hasta que
			 * el vector de los nombres de las carpetas este vacio
			 */
			if len(paths) == 0 {
				return
			}
			MKDIR(Carpeta, paths, RutaDisco, SuperBloque, int(apuntadorAVD))
			return
		}
		/*
		 * Esta parte significa que la posicion no esta vacia y que ya existe una creada
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
			MKDIR(CarpetaHija, paths, RutaDisco, SuperBloque, int(apuntador))
			return
		}

		/*
		 *  APUNTADORES INDIRECTOS
		 *  Hace referencia al AVD que se crea si un AVD ya esta lleno , entonces se crea una copia
		 */
		// El valor del apuntador es el siguiente o sea la copia
		apuntador = AVD.VirtualDirectorio

		if apuntador == -1 { //Significa que no hay una copia aun

			//Creamos el AVD que hace referencia a la copia
			var CopiaAVD Arbol
			//Le seteamos el nombre del AVD anterior
			CopiaAVD.AVDNombreDirectorio = AVD.AVDNombreDirectorio
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
			//Escribimos el La copia del AVD
			s := &CopiaAVD
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, s)
			File.Write(binario.Bytes())
			File.Seek(int64(SuperBloque.StartBmArbolDirectorio+SuperBloque.FirstFreeAvd), 0)
			//Escribimos el 1 en el bitmap para representar el espacio ocupado
			s1 := &uno
			var binario1 bytes.Buffer
			binary.Write(&binario1, binary.BigEndian, s1)
			File.Write(binario1.Bytes())

			SuperBloque.ArbolVirtualFree--
			SuperBloque.FirstFreeAvd++

			//Reescribo avd actual
			File.Seek(int64(SuperBloque.StartArbolDirectorio+(apuntador*int32(unsafe.Sizeof(AVD)))), 0)
			s2 := &AVD
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s2)
			File.Write(binario2.Bytes())

			//Cerramos el archivo
			File.Close()
			// Llamamos al metodo recursivamente para que cree la carpeta en la copia del AVD y le mandamos la copia
			MKDIR(CopiaAVD, paths, RutaDisco, SuperBloque, int(apuntadorAVD))
			return

		}
		//Significa que ya existe una copia de ese AVD
		var CarpetaCopia Arbol
		CarpetaCopia = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio+(apuntador*int32(unsafe.Sizeof(CarpetaCopia)))))
		//Cerramos el archivo
		File.Close()
		//Llamamos recursivamente al metodo MKDIR para que repita todo el proceso en la copia del AVD
		MKDIR(CarpetaCopia, paths, RutaDisco, SuperBloque, int(apuntador))
		return

	}

}

//MKFILE is...
func MKFILE(id string, path string, p bool, size int, count string) {

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

			var Root Arbol
			Root = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

			//Eliminamos el espacio vacio y el nombre del archivo home/jose
			NombreArchivo := Rutas[len(Rutas)-1]
			Rutas = Rutas[1:]
			Rutas = Rutas[:len(Rutas)-1]

			fmt.Println(NombreArchivo)
			File.Close()
			//Mandamos a crear las carpetas si no estan creadas
			MKDIR(Root, Rutas, PathDisco, SuperBloque, 0)

			File = getFile(PathDisco)
			SuperBloque = readSuperBloque(File, int64(PartStart))

			Root = readArbolVirtualDirectorio(File, int64(SuperBloque.StartArbolDirectorio))

			File.Close()

			RecorrerArbol(Root, Rutas, PathDisco, SuperBloque, size, count)

		} else {
			ErrorMessage("[MKFILE] -> No se encuentra ningun disco en esa ruta")
		}

	} else {
		ErrorMessage("[MKFILE] -> No hay ninguna particion montada con ese id")
	}
}

//RecorrerArbol is...
func RecorrerArbol(root Arbol, Rutas []string, PathDisco string, Superbloque SB, size int, count string) {

	if VerificarRuta(PathDisco) {

		File := getFile(PathDisco)
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
					File.Close()
					CrearArchivo(Archivos, int(Apuntador), Rutas, PathDisco, Superbloque, size, count)
					return
				}
				File.Close()
				RecorrerArbol(CarpetaHija, Rutas, PathDisco, Superbloque, size, count)
				return

			}

		}
		// Buscamos en la copida de la carpeta
		Apuntador = root.VirtualDirectorio
		var CopiaCarpeta Arbol
		CopiaCarpeta = readArbolVirtualDirectorio(File, int64(Superbloque.StartArbolDirectorio+(Apuntador*int32(unsafe.Sizeof(CopiaCarpeta)))))
		File.Close()
		RecorrerArbol(CopiaCarpeta, Rutas, PathDisco, Superbloque, size, count)
		return

	}
	ErrorMessage("[MKFILE] -> No hay ningun disco en la ruta indicada")

}

//CrearArchivo is...
func CrearArchivo(Archivo DetalleDirectorio, Apuntador int, Rutas []string, RutaDisco string, SuperB SB, size int, count string) {

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
	return DetalleD
}

//InicializarInodo is...
func InicializarInodo(Inodo TablaInodo) TablaInodo {
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
	fmt.Println(user, password, id)
	// TODO : Hacer login
	var RutaDisco string = listaParticiones.GetDireccion(id)

	if RutaDisco != "null" {
		var PartName string = listaParticiones.GetPartName(id)
		PosParticion := BuscarParticion(RutaDisco, PartName)

		if PosParticion != -1 {

			/*var masterboot MBR
			var SuperBloque SB
			var Inodo TablaInodo
			File := getFile(RutaDisco)
			masterboot = readMBR(File)
			SuperBloque = readSuperBloque(File, int64(masterboot.Particion[PosParticion].PartStart))
			Inodo = readInodo(File, int64(SuperBloque.StartInodos+int32(unsafe.Sizeof(Inodo))))

			File.Seek(int64(SuperBloque.StartInodos+int32(unsafe.Sizeof(Inodo))), 0)

			inodo.i_atime = time(nullptr)
			fwrite(&inodo, sizeof(InodoTable), 1, fp)
			fclose(fp)
			currentSession.inicioSuper = masterboot.mbr_partition[index].part_start
			currentSession.fit = masterboot.mbr_partition[index].part_fit
			currentSession.inicioJournal = masterboot.mbr_partition[index].part_start+static_cast < int > (sizeof(SuperBloque))
			currentSession.tipo_sistema = super.s_filesystem_type
			return verificarDatos(user, password, direccion)*/

		} else {

		}

	} else {
		ErrorMessage("[LOGIN] -> No se encuentra ninguna particion montada con ese id")
	}
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
