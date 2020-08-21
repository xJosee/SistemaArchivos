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

//TODO : Declarar lista de particiones
var listaParticiones = Estructuras.Lista{
	Contador: 0,
	Primero:  nil,
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
func reWriteEBR(file *os.File, Disco EBR, seek int64) {
	file.Seek(seek, 0)
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
func readMBR(file *os.File) MBR {
	m := MBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytesMBR(file, size)
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.BigEndian, &m)
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
func readEBR(file *os.File, seek int64) EBR {
	file.Seek(seek, 0)
	m := EBR{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytesEBR(file, size)
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

//FDISK is...
func FDISK(size int, unit byte, path string, Type byte, fit byte, delete string, name string, add int) bool {

	if delete != "" {

		EliminarParticion(path, name, delete)
	} else if Type == 'p' {
		if CrearParticionPrimaria(path, CalcularSize(size, unit), name, fit) {
			return true
		}
	} else if Type == 'e' {
		CrearParticionExtendida(path, CalcularSize(size, unit), name, fit)
	} else if Type == 'l' {
		CrearParticionLogica(path, name, CalcularSize(size, unit), fit)
	} else if delete != "" {

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

//ReporteDisco is...
func ReporteDisco(direccion string, destino string, extension string) {

	//TODO : Reporte del disco

	var auxDir string = direccion

	if VerificarRuta(auxDir) {
		fp := getFile(auxDir)
		os.Create("Reportes/grafica.dot")
		graphDot := getFile("Reportes/grafica.dot")

		fmt.Fprintf(graphDot, "digraph G{\n")
		fmt.Fprintf(graphDot, "  tbl [\n    shape=box\n    label=<\n")
		fmt.Fprintf(graphDot, "     <table border='0' cellborder='1' width='600' height='200' color='coral'>\n")
		fmt.Fprintf(graphDot, "     <tr>\n")
		fmt.Fprintf(graphDot, "     <td height='200' width='100'> MBR </td>\n")

		var masterboot MBR
		fp.Seek(0, 0)
		masterboot = readMBR(fp)

		var total int = int(masterboot.Size)
		var espacioUsado int = 0

		for i := 0; i < 4; i++ {
			var parcial int = int(masterboot.Particion[i].PartSize)
			if masterboot.Particion[i].PartStart != -1 {
				var porcentajeReal int = (parcial * 100) / total
				var porcentajeAux int = (porcentajeReal * 500) / 100
				espacioUsado += porcentajeReal

				if masterboot.Particion[i].PartStatus != '1' {

					if masterboot.Particion[i].PartType == 'P' { //Verificar Primaria

						fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>PRIMARIA <br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')

						if i != 3 {
							var p1 int = int(masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)
							var p2 int = int(masterboot.Particion[i+1].PartStart)

							if masterboot.Particion[i+1].PartStart != -1 {

								if (p2 - p1) != 0 { //Verficiar Si hay fragmentacion
									var fragmentacion int = p2 - p1
									var porcentajeReal int = (fragmentacion * 100) / total
									var porcentajeAux int = (porcentajeReal * 500) / 100

									fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>LIBRE<br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')
								}
							}

						} else {
							var p1 int = int(masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)
							var mbrTamano int = total + int(unsafe.Sizeof(masterboot))

							if (mbrTamano - p1) != 0 { // LIBRE
								var libre int = (mbrTamano - p1) + int(unsafe.Sizeof(masterboot))
								var porcentajeReal int = (libre * 100) / total
								var porcentajeAux int = (porcentajeReal * 500) / 100

								fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>LIBRE<br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')
							}
						}

					} else {
						//Particion Extendida
						extendedBoot := EBR{
							PartNext: -2,
						}
						fmt.Fprintf(graphDot, "     <td  height='200' width='%.1f'>\n     <table border='0'  height='200' WIDTH='%.1f' cellborder='1'>\n", float32(porcentajeReal), float32(porcentajeReal))
						fmt.Fprintf(graphDot, "     <tr>  <td height='60' colspan='15'>EXTENDIDA</td>  </tr>\n     <tr>\n")

						extendedBoot = readEBR(fp, int64(masterboot.Particion[i].PartStart))

						if extendedBoot.PartSize != 0 { //Si hay mas de alguna logica

							for extendedBoot.PartNext != -1 && (extendedBoot.PartNext < (masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)) {

								if extendedBoot.PartNext == -2 {
									extendedBoot = readEBR(fp, int64(masterboot.Particion[i].PartStart))
								} else {
									extendedBoot = readEBR(fp, int64(extendedBoot.PartNext))
								}
								parcial = int(extendedBoot.PartSize)
								porcentajeReal = (parcial * 100) / total

								if porcentajeReal != 0 {

									if extendedBoot.PartStatus != '1' {
										fmt.Fprintf(graphDot, "     <td height='140'>EBR</td>\n")
										fmt.Fprintf(graphDot, "     <td height='140'>LOGICA<br/> %.1f%c</td>\n", float32(porcentajeReal), '%')
									} else { //Espacio no asignado
										fmt.Fprintf(graphDot, "      <td height='150'>LIBRE<br/>  %.1f%c</td>\n", float32(porcentajeReal), '%')
									}
									if extendedBoot.PartNext == -1 {
										parcial = int((masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize) - (extendedBoot.PartStart + extendedBoot.PartSize))
										porcentajeReal = (parcial * 100) / total
										if porcentajeReal != 0 {
											fmt.Fprintf(graphDot, "     <td height='150'>LIBRE<br/>  %.1f%c </td>\n", float32(porcentajeReal), '%')
										}
										break
									}

								}
							}
						} else {
							fmt.Fprintf(graphDot, "     <td height='140'>  %.1f%c</td>", float32(porcentajeReal), '%')
						}

						fmt.Fprintf(graphDot, "     </tr>\n     </table>\n     </td>\n")

						if i != 3 {
							var p1 int = int(masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)
							var p2 int = int(masterboot.Particion[i+1].PartStart)

							if masterboot.Particion[i+1].PartStart != -1 {

								if (p2 - p1) != 0 { //Hay fragmentacion
									var fragmentacion int = p2 - p1
									var porcentajeReal int = (fragmentacion * 100) / total
									var porcentajeAux int = (porcentajeReal * 500) / 100
									fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>LIBRE<br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')
								}

							}
						} else {
							var p1 int = int(masterboot.Particion[i].PartStart + masterboot.Particion[i].PartSize)
							var mbrTamano int = total + int(unsafe.Sizeof(masterboot))

							if (mbrTamano - p1) != 0 { //Libre

								var libre int = (mbrTamano - p1) + int(unsafe.Sizeof(masterboot))
								var porcentajeReal int = (libre * 100) / total
								var porcentajeAux int = (porcentajeReal * 500) / 100

								fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>LIBRE<br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')
							}
						}

					}

				} else {
					fmt.Fprintf(graphDot, "     <td height='200' width='%.1f'>LIBRE <br/>  %.1f%c</td>\n", float32(porcentajeAux), float32(porcentajeReal), '%')
				}
			}
		}

		fmt.Fprintf(graphDot, "     <td height='200'> LIBRE %.1f%c\n     </td>", float32((100 - espacioUsado)), '%')

		fmt.Fprintf(graphDot, "     </tr> \n     </table>        \n>];\n\n}")
		graphDot.Close()
		fp.Close()

		var comando string = "dot -T" + "png" + " grafica.dot -o " + "/home/jose/Escritorio/"
		cmd := exec.Command(comando)
		cmd.Output()
		SuccessMessage("[REP] -> Reporte del disco generado correctamente")
	} else {
		ErrorMessage("[REP] -> No se encuentra el disco")
	}
}

//EliminarParticion is...
func EliminarParticion(path string, name string, delete string) {
	//TODO : Eliminar Particion

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
func AgregarQuitarEspacio() {
	//TODO : Agregar o Quitar espacio

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

			ebr := EBR{
				PartNext: -2,
			}
			for ebr.PartNext != -1 && (ebr.PartNext < masterboot.Particion[extendida].PartStart+masterboot.Particion[extendida].PartSize) {
				if ebr.PartNext == -2 {
					ebr = readEBR(File, int64(masterboot.Particion[extendida].PartStart))
				} else {
					ebr = readEBR(File, int64(ebr.PartNext))
				}
				var nameByte [16]byte
				copy(nameByte[:], name)
				if bytes.Compare(ebr.PartName[:], nameByte[:]) == 0 {
					return int((ebr.PartNext - int32(unsafe.Sizeof(ebr))))
				}
			}
		}
		File.Close()
	}
	return -1
}

//ReporteEBR is...
func ReporteEBR(path string) {

	if VerificarRuta(path) {

		File := getFile(path)
		graphDot := getFile("Reportes/grafica.dot")

		fmt.Fprintf(graphDot, "digraph G{ \n")
		fmt.Fprintf(graphDot, "node [shape=plaintext]\n")
		fmt.Fprintf(graphDot, "tbl[\nlabel=<\n")
		fmt.Fprintf(graphDot, "<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n")
		fmt.Fprintf(graphDot, " <tr ><td bgcolor= 'lightblue' ><b><font color='blue'>MBR</font></b></td></tr>")
		fmt.Fprintf(graphDot, "<tr>  <td width='150'> <b>Nombre</b> </td> <td width='150'> <b>Valor</b> </td>  </tr>\n")

		var MB MBR
		File.Seek(0, 0)
		MB = readMBR(File)

		var tamano int = int(MB.Size)

		fmt.Fprintf(graphDot, "<tr>  <td><b>Size</b></td><td>%d</td>  </tr>\n", tamano)

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

				fmt.Fprintf(graphDot, "<tr ><td bgcolor= 'lightblue' ><b><font color='blue'>Particion%d</font></b></td></tr>\n", (i + 1))
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
					fmt.Fprintf(graphDot, "<tr ><td bgcolor= 'lightblue' ><b><font color='blue'>EBR</font></b></td></tr>")
					fmt.Fprintf(graphDot, "<tr ><td width='150'><b>Nombre</b></td> <td width='150'><b>Valor</b></td>  </tr>\n")
					var status [3]byte
					if extendedBoot.PartStatus == '0' {
						copy(status[:], "0")
					} else if extendedBoot.PartStatus == '2' {
						copy(status[:], "2")
					}

					fmt.Fprintf(graphDot, "<tr>  <td><b>part_status_1</b></td> <td>%s</td>  </tr>\n", string(status[:]))
					fmt.Fprintf(graphDot, "<tr>  <td><b>part_fit_1</b></td> <td>%c</td>  </tr>\n", extendedBoot.PartFit)
					fmt.Fprintf(graphDot, "<tr>  <td><b>part_start_1</b></td> <td>%d</td>  </tr>\n", extendedBoot.PartStart)
					fmt.Fprintf(graphDot, "<tr>  <td><b>part_size_1</b></td> <td>%d</td>  </tr>\n", extendedBoot.PartSize)
					fmt.Fprintf(graphDot, "<tr>  <td><b>part_next_1</b></td> <td>%d</td>  </tr>\n", extendedBoot.PartNext)
					fmt.Fprintf(graphDot, "<tr>  <td><b>part_name_1</b></td> <td>%s</td>  </tr>\n", string(extendedBoot.PartName[:]))
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
		var comando string = "dot -T" + "png" + " grafica.dot -o " + "/home/jose/Escritorio/"
		cmd := exec.Command(comando)
		cmd.Output()
		SuccessMessage("[REP] -> Reporte del disco generado correctamente")

	}

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
