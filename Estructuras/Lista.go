package estructuras

import "fmt"

//Nodo is...
type Nodo struct {
	Direccion string
	Nombre    string
	Letra     byte
	Num       int
	Siguiente *Nodo
}

//Lista is...
type Lista struct {
	Contador int
	Primero  *Nodo
}

//Insertar is...
func (Lista *Lista) Insertar(Nodo *Nodo) {
	Lista.Contador++
	aux := Lista.Primero
	if Lista.Primero == nil {
		Lista.Primero = Nodo
	} else {
		for aux.Siguiente != nil {
			aux = aux.Siguiente
		}
		aux.Siguiente = Nodo
	}
}

//GetSize is...
func (Lista *Lista) GetSize() int {
	return Lista.Contador
}

//EliminarNodo is...
func (Lista *Lista) EliminarNodo(ID string) int {
	Lista.Contador--
	aux := Lista.Primero
	var tempID string = "vd"
	tempID += string(aux.Letra)
	tempID = fmt.Sprint(tempID, aux.Num)

	if ID == tempID {
		Lista.Primero = aux.Siguiente
		//unsafe.Pointer(aux)
		return 1
	}
	var aux2 *Nodo = nil
	for aux != nil {
		tempID = "vd"
		tempID += string(aux.Letra)
		tempID = fmt.Sprint(tempID, aux.Num)
		if ID == tempID {
			aux2.Siguiente = aux.Siguiente
			return 1
		}
		aux2 = aux
		aux = aux.Siguiente
	}

	return 0
}

//BuscarLetra is...
func (Lista *Lista) BuscarLetra(direccion string, nombre string) int {
	aux := Lista.Primero
	var retorno int = 'a'

	for aux != nil {
		if (direccion == aux.Direccion) && (nombre == aux.Nombre) {
			return -1
		}

		if direccion == aux.Direccion {
			return int(aux.Letra)
		} else if retorno == int(aux.Letra) {
			retorno++
		}

		aux = aux.Siguiente
	}

	return retorno
}

//BuscarNumero is...
func (Lista *Lista) BuscarNumero(direccion string, nombres string) int {
	var retorno int = 1
	aux := Lista.Primero
	for aux != nil {
		if (direccion == aux.Direccion) && (retorno == aux.Num) {
			retorno++
		}
		aux = aux.Siguiente
	}
	return retorno
}

//GetDireccion is...
func (Lista *Lista) GetDireccion(id string) string {
	aux := Lista.Primero
	for aux != nil {
		tempID := "vd"
		tempID += string(aux.Letra)
		tempID = fmt.Sprint(tempID, aux.Num)
		if id == tempID {
			return aux.Direccion
		}
		aux = aux.Siguiente
	}
	return "null"
}

//GetPartName is...
func (Lista *Lista) GetPartName(id string) string {
	aux := Lista.Primero
	for aux != nil {
		tempID := "vd"
		tempID += string(aux.Letra)
		tempID = fmt.Sprint(tempID, aux.Num)
		if id == tempID {
			return aux.Nombre
		}
		aux = aux.Siguiente
	}
	return "null"
}

//BuscarNodo is...
func (Lista *Lista) BuscarNodo(direccion string, nombre string) bool {
	aux := Lista.Primero
	for aux != nil {
		if (aux.Direccion == direccion) && (aux.Nombre == nombre) {
			return true
		}
		aux = aux.Siguiente
	}
	return false
}

//Listar is...
func (Lista *Lista) Listar() {
	aux := Lista.Primero
	for aux != nil {
		fmt.Println("NOMBRE : ", aux.Nombre, " <----> ", "ID : vd", string(aux.Letra), aux.Num)
		aux = aux.Siguiente
	}
}
