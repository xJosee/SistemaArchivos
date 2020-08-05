package comandos

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
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
	err := ioutil.WriteFile(path+name+".disk", []byte(""), 0755)
	if err == nil {
		writeFile(path+name+".disk", CalcularSize(size, unit))
		CrearRaid(size, fit, unit, path, name)
	}
}

//CrearRaid is ...
func CrearRaid(size int, fit byte, unit byte, path string, name string) {
	err := ioutil.WriteFile(path+name+"Raid.disk", []byte(""), 0755)
	if err == nil {
		writeFile(path+name+"Raid.disk", CalcularSize(size, unit))
	}
}

//WriteFile is ...
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
		var binBuf bytes.Buffer
		binary.Write(&binBuf, binary.BigEndian, s)
		//b :=binbuf.Bytes()
		//l := len(b)
		//fmt.Println(l)
		writeNextBytes(file, binBuf.Bytes())

	}
}
func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}

//CalcularSize is ...
func CalcularSize(size int, unit byte) int {
	if unit == 'M' {
		return 63 * size * 1000
	} else if unit == 'K' {
		return 63 * size
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
