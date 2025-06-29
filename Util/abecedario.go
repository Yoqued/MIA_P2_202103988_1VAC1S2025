package util

import (
	"strings"
)

var AbecedarioMayuscula = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "Ñ", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// Elimina el primer elemento y devuelve la letra eliminada
func EliminarPrimeraLetra() string {
	letra := AbecedarioMayuscula[0]
	AbecedarioMayuscula = AbecedarioMayuscula[1:]
	return letra
}

func AgregarLetra(letra string) {
	// Verifica si la letra ya está en el slice
	letra = strings.ToUpper(letra)
	for _, l := range AbecedarioMayuscula {
		if l == letra {
			return // Ya existe, no la agrega
		}
	}
	// Insertar letra al inicio
	AbecedarioMayuscula = append([]string{letra}, AbecedarioMayuscula...)
}
