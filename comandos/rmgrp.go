package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type RMGRP struct {
	name string
}

func ParseRmgrp(tokens []string) ([]string, error) {
	cmd := &RMGRP{}
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
		case "-name":
			if len(value) > 10 {
				return tokens[cont:], errors.New("limite de caracteres, debe ser menos de 10")
			}
			cmd.name = value
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if cmd.name == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -user")
	}
	if reflect.DeepEqual(util.UserOnline, util.USER{}) {
		return tokens[contador:], errors.New("Se necesita iniciar sesion para ejecutar este comando")
	}
	if !groupExiste(cmd.name, util.GruposId[util.UserOnline.Id].GNames) {
		return tokens[contador:], errors.New("El grupo indicado no existe")
	}
	var UserString string

	// Extraer el superbloque de la particion
	startPart := int64(binary.LittleEndian.Uint32(util.UserOnline.Particion.Start[:]))
	superblock, _ := estructuras.ReadSuperBlock(util.UserOnline.Path, startPart)
	superblock.ReadFile(util.UserOnline.Path, []string{"users.txt"}, 0, &UserString)
	lineas := strings.Split(UserString, "\n")

	var newLines []string
	// Analizar el txt línea por línea
	for _, linea := range lineas {
		partes := strings.Split(linea, ",")
		if len(partes) >= 3 {
			if partes[1] == "G" {
				if partes[2] == cmd.name {
					partes[0] = "0"
				}
			}
		}
		newLines = append(newLines, strings.Join(partes, ","))
	}
	newUserString := strings.Join(newLines, "\n")
	fmt.Println("Nueva cadena a guardar:", newUserString)
	err := estructuras.ModifyUserGroup(superblock, util.UserOnline.Path, newUserString)
	if err != nil {
		return tokens[contador:], err
	}
	fmt.Println("Grupo removido con exito")
	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Grupo %s removido con exito!", cmd.name))

	//Guardar todos los cambios del superbloque
	superblock.WriteToFile(util.UserOnline.Path, int(startPart))

	return tokens[contador:], nil
}
