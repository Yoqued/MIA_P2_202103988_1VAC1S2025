package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Función para parsear el comando RMDISK
func ParseRmdisk(tokens []string) ([]string, error) {
	var driveLetter string
	var contador int

	for cont, token := range tokens {
		if token == "<newline>" {
			contador += 1
			break
		}

		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return tokens[cont:], fmt.Errorf("formato de parámetro inválido: %s", token)
		}
		key, value := strings.ToLower(parts[0]), parts[1]

		switch key {
		case "-driveletter":
			// Validar que sea una sola letra
			if len(value) != 1 || !isLetter(value[0]) {
				return tokens[cont:], fmt.Errorf("la letra del disco debe ser un carácter de la A-Z")
			}
			driveLetter = strings.ToUpper(value)
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	// Verificar si el parámetro -driveletter se proporcionó
	if driveLetter == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -driveletter")
	}

	// Construir la ruta completa del archivo
	path := filepath.Join("./MIA/P1/", fmt.Sprintf("%s.dsk", driveLetter))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("El disco con letra %s no existe en la ruta %s.\n", driveLetter, path)
		return tokens[contador:], fmt.Errorf("el disco con letra %s no existe", driveLetter)
	}

	// Eliminar el archivo directamente (sin confirmación)
	util.AgregarLetra(driveLetter) // Agregar la letra de disco al abecedario
	err := os.Remove(path)
	if err != nil {
		return tokens[contador:], fmt.Errorf("error al eliminar el disco: %v", err)
	}

	fmt.Printf("Disco %s eliminado con éxito desde %s.\n", driveLetter, path)
	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Disco %s removido con éxito desde %s!", driveLetter, path))

	return tokens[contador:], nil
}

// Función auxiliar para verificar si un byte es una letra
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
