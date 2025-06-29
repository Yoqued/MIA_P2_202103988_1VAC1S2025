package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

func ParseMove(tokens []string) ([]string, error) {
	var Copypath, Destpath string
	var contador int
	for cont, token := range tokens {

		if token == "<newline>" {
			contador += 1
			break
		}

		// Divide cada token en clave y valor usando "=" como delimitador
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return tokens[cont:], fmt.Errorf("formato de parámetro inválido: %s", token)
		}
		key, value := strings.ToLower(parts[0]), parts[1]

		// Key representa las palabras claves de cada atributo
		switch key {
		case "-path":
			Copypath = value
		case "-destino":
			Destpath = value
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if Destpath == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -destino")
	}
	if Copypath == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -path")
	}

	startPart := int64(binary.LittleEndian.Uint32(util.UserOnline.Particion.Start[:]))
	superblock, _ := estructuras.ReadSuperBlock(util.UserOnline.Path, startPart)
	parentsDir, Copyname := util.GetParentDirectories(Copypath)

	// Obtemos el inodo al que apunta el directorio que deseamos copiar
	InodoId, err := superblock.SearchFileInodo(util.UserOnline.Path, 0, parentsDir, Copyname)
	if err != nil {
		return tokens[contador:], err
	}

	// Removemos el elemento para cambiarlo de lugar
	err = superblock.RemoveFile(util.UserOnline.Path, parentsDir, 0, Copyname)
	if err != nil {
		return tokens[contador:], err
	}

	// Manejo de copiar el archivo/carpeta en la ruta indicada
	fileParts := strings.Split(Destpath, "/")
	fileParts = fileParts[1:]
	if len(fileParts) > 0 && fileParts[len(fileParts)-1] == "" {
		fileParts = fileParts[:len(fileParts)-1]
	}
	if len(fileParts) == 0 || (len(fileParts) == 1 && fileParts[0] == "") {
		fileParts = []string{}
	}

	_, err = superblock.CopyFile(util.UserOnline.Path, 0, fileParts, Copyname, InodoId)
	if err != nil {
		return tokens[contador:], err
	}

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Carpeta %s movida exitosamente en %s", Copyname, Destpath))

	return tokens[contador:], nil
}
