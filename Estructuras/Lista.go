package estructuras

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

//Constructor is...
func (Lista *Lista) Constructor() {
	Lista.Contador = 0
	Lista.Primero = nil
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

//getSize is...
func (Lista *Lista) getSize() int {
	return Lista.Contador
}

//EliminarNodo is...
func (Lista *Lista) EliminarNodo(ID string) int {
	Lista.Contador--
	aux := Lista.Primero
	tempID := "vd"
	tempID += string(aux.Letra)
	tempID += string(aux.Num)

	if ID == tempID {
		Lista.Primero = aux.Siguiente
		//free(aux)
		//return 1
	} else {
		var aux2 *Nodo = nil
		for aux != nil {
			tempID = "vd"
			tempID += string(aux.Letra)
			tempID += string(aux.Num)
			if ID == tempID {
				aux2.Siguiente = aux.Siguiente
				return 1
			}
			aux2 = aux
			aux = aux.Siguiente
		}
	}
	return 0
}
