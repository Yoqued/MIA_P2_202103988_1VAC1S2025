package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

func ParseRename(tokens []string) ([]string, error) {
	var path, name string
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
			path = value
		case "-name":
			name = value
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if path == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -path")
	}
	if name == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -name")
	}

	startPart := int64(binary.LittleEndian.Uint32(util.UserOnline.Particion.Start[:]))
	superblock, _ := estructuras.ReadSuperBlock(util.UserOnline.Path, startPart)
	parentsDir, dest := util.GetParentDirectories(path)

	err := superblock.RenameFile(util.UserOnline.Path, parentsDir, 0, dest, name)
	if err != nil {
		return tokens[contador:], err
	}

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Se reemplazo el nombre de %s con % s exitosamente!", dest, name))

	return tokens[contador:], nil
}
