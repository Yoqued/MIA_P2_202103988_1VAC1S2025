package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type MOUNT struct {
	path string
	name string
}

var abecedario string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var currentedriveletter string = ""

func ParseMount(tokens []string) (*MOUNT, []string, error) {
	cmd := &MOUNT{}
	var contador int
	for cont, token := range tokens {

		if token == "<newline>" {
			contador += 1
			break
		}

		// Divide cada token en clave y valor usando "=" como delimitador
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, tokens[cont:], fmt.Errorf("formato de parámetro inválido: %s", token)
		}
		key, value := strings.ToLower(parts[0]), parts[1]

		// Key representa las palabras claves de cada atributo
		switch key {
		case "-driveletter":
			if len(value) != 1 || !isLetter(value[0]) {
				return nil, tokens[cont:], errors.New("la letra del disco debe ser un carácter de la A-Z")
			}
			driveLetter := strings.ToUpper(value)
			cmd.path = filepath.Join("/home/david/MIA/P1/", fmt.Sprintf("%s.dsk", driveLetter))
			currentedriveletter = driveLetter
		case "-name":
			cmd.name = value
		default:
			return nil, tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if cmd.path == "" {
		return nil, tokens[contador:], errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, tokens[contador:], errors.New("faltan parámetros requeridos: -name")
	}

	mbr, err := objs.ReadMbr(cmd.path)
	if err != nil {
		return nil, tokens[contador:], err
	}
	disco, err := objs.BuscarDisco(cmd.path)
	if err != nil {
		return nil, tokens[contador:], err
	}

	posparticion := mbr.PosicionParticion(cmd.name, cmd.path)
	if disco.Letter == "" {
		disco.Letter = IDLetra()
	}
	if posparticion == -1 {
		return nil, tokens[contador:], fmt.Errorf("La particion %s no existe", cmd.name)
	}
	var id string

	// Validar si es una particion primaria o extendida
	fmt.Println("POSPARTICION", posparticion)
	if posparticion < 4 {
		if string(mbr.Particiones[posparticion].Status[:]) == "1" {
			return nil, tokens[contador:], fmt.Errorf("La particion %s ya ha sido montada", cmd.name)
		}
		//Cambiar el estado de la particion a montada
		copy(mbr.Particiones[posparticion].Status[:], "1")

		//Crear id para la particion
		idPos := len(disco.Id)
		id = crearID(idPos, disco.Letter)
		fmt.Println("ID", id)
		copy(mbr.Particiones[posparticion].Id[:], id)
		disco.Id = append(disco.Id, id)

		mbr.WriteToFile(cmd.path)
	} else {
		return nil, tokens[contador:], errors.New("Solo se pueden montar particiones primarias")
	}

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Particion %s montada con el id %s", cmd.name, id))
	util.GruposId[id] = util.GRUPO{
		GNames:     []string{},
		UNames:     []string{},
		GroupCount: 0,
		UserCount:  0,
	}

	return cmd, tokens[contador:], nil // Devuelve el comando MKDISK creado
}

func crearID(posicion int, letter string) string {
	particionID := posicion + 1
	if letter == "" {
		letter = IDLetra()
	}
	return "88" + strconv.Itoa(particionID) + letter
}

func IDLetra() string {
	// Buscar y eliminar el currentDriveLetter del abecedario
	index := strings.Index(abecedario, strings.ToUpper(currentedriveletter))
	fmt.Println(index, "-----------------")
	if index == -1 {
		// Si la letra no está disponible, tomar la primera disponible
		if len(abecedario) == 0 {
			return ""
		}
		nextLetter := string(abecedario[0])
		abecedario = abecedario[1:]
		return nextLetter
	}

	nextLetter := string(abecedario[index])
	abecedario = abecedario[:index] + abecedario[index+1:]
	return nextLetter
}

func ParseUnmount(tokens []string) ([]string, error) {
	var id string
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
		case "-id":
			if len(value) != 4 {
				return nil, fmt.Errorf("El ID debe tener exactamente 4 caracteres")
			}

			// Extraemos las partes del ID original (ej: "A188")
			letra := string(value[0])  // "A"
			numero := string(value[1]) // "1"
			constante := value[2:]     // "88"

			// Reordenamos: constante + numero + letra → "88" + "1" + "A" = "881A"
			nuevoID := constante + numero + letra
			id = nuevoID
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if id == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -id")
	}
	if _, exists := util.GruposId[id]; !exists {
		return tokens[contador:], fmt.Errorf("El ID %s no existe", id)
	}

	path, err := objs.BuscarPathID(id)
	if err != nil {
		return tokens[contador:], err
	}

	mbr, err := objs.ReadMbr(path)
	if err != nil {
		return tokens[contador:], err
	}

	particion := mbr.BuscarIdParticion(id)
	if particion == nil {
		return tokens[contador:], fmt.Errorf("no se eoncontro particion con id: %s", id)
	}

	copy(particion.Correlative[:], make([]byte, len(particion.Correlative)))
	copy(particion.Status[:], "0")
	copy(particion.Id[:], make([]byte, len(particion.Id)))

	delete(util.GruposId, id)
	err = objs.RemoveIdDisk(path, id)
	if err != nil {
		return tokens[contador:], err
	}

	mbr.WriteToFile(path)

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Particion %s desmontada con exito", id))

	return tokens[contador:], nil
}
