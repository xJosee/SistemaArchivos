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
