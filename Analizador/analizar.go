package analizar

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/comandos"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func AnalizarTexto(entrada string) (interface{}, []error) {
	util.Respuestas = []string{}
	lineas := strings.Split(entrada, "\n")
	var tokens []string
	var errores []error
	var myInterface interface{}

	// Iterar sobre cada línea
	fmt.Println(lineas)
	for _, linea := range lineas {
		fmt.Println("PASO1", linea)
		// Verificar si la línea es un comentario
		if strings.HasPrefix(linea, "#") {
			util.Respuestas = append(util.Respuestas, linea)
			continue // No procesar más la línea si es un comentario
		}

		// Separar cada línea respetando las comillas
		palabras, err := parseLine(linea)
		if err != nil {
			util.Respuestas = append(util.Respuestas, err.Error())
			return nil, errores
		}

		tokens = append(tokens, palabras...)
		tokens = append(tokens, "<newline>")

		for len(tokens) > 0 {
			token := tokens[0]

			// Ignorar salto de linea
			if token == "<newline>" {
				util.Respuestas = append(util.Respuestas, "\n")
				tokens = tokens[1:] // Remover <newline> y continuar
				continue
			}

			token = strings.ToUpper(token)

			if token == "MKDISK" {
				fmt.Println("ENTRO A MKDIR")
				objeto, newtokens, err := comandos.ParseMkdisk(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break // Termina el ciclo actual y continúa con el próximo conjunto de tokens
				}
				fmt.Println("objeto:", objeto)
				tokens = newtokens
			} else if token == "RMDISK" {
				newtokens, err := comandos.ParseRmdisk(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
				util.Respuestas = append(util.Respuestas, "Disco eliminado con exito!")
			} else if token == "FDISK" {
				objeto, newtokens, err := comandos.ParseFdisk(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				fmt.Println("objeto:", objeto)
				tokens = newtokens
			} else if token == "MOUNT" {
				objeto, newtokens, err := comandos.ParseMount(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				fmt.Println("objeto:", objeto)
				tokens = newtokens
			} else if token == "UNMOUNT" {
				newtokens, err := comandos.ParseUnmount(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "MKFS" {
				objeto, newtokens, err := comandos.ParseMkfs(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				fmt.Println("objeto:", objeto)
				tokens = newtokens
			} else if token == "CAT" {
				newtokens, err, exitCat := comandos.ParseCat(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
				util.Respuestas = append(util.Respuestas, exitCat)
			} else if token == "LOGIN" {
				newtokens, err := comandos.ParseLogin(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "LOGOUT" {
				if reflect.DeepEqual(util.UserOnline, util.USER{}) {
					util.Respuestas = append(util.Respuestas, "No hay ninguna sesion activa para hacer Logout")
					tokens = recortarTokens(tokens)
					break
				} else {
					util.UserOnline = util.USER{}
					util.Respuestas = append(util.Respuestas, "Se ha cerrado la sesion")
				}
				tokens = tokens[2:]
			} else if token == "MKGRP" {
				newtokens, err := comandos.ParseMkgrp(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "MKUSR" {
				newtokens, err := comandos.ParseMkusr(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "RMGRP" {
				newtokens, err := comandos.ParseRmgrp(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "RMUSR" {
				newtokens, err := comandos.ParseRmusr(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "CHGRP" {
				newtokens, err := comandos.ParseChgrp(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "REP" {
				newtokens, err := comandos.ParseRep(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "MKFILE" {
				newtokens, err := comandos.ParseMkfile(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "MKDIR" {
				newtokens, err := comandos.ParseMkdir(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "RENAME" {
				newtokens, err := comandos.ParseRename(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "EDIT" {
				newtokens, err := comandos.ParseEdit(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "REMOVE" {
				newtokens, err := comandos.ParseRemove(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "COPY" {
				newtokens, err := comandos.ParseCopy(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "MOVE" {
				newtokens, err := comandos.ParseMove(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else if token == "FIND" {
				newtokens, err := comandos.ParseFind(tokens[1:])
				if err != nil {
					util.Respuestas = append(util.Respuestas, err.Error())
					tokens = recortarTokens(tokens)
					break
				}
				tokens = newtokens
			} else {
				// Si el comando no es reconocido, devuelve un error
				errores = append(errores, fmt.Errorf("comando desconocido: %s", token))
				return nil, errores
			}
		}
	}

	myInterface = util.Respuestas
	return myInterface, errores
}

// parseLine analiza una línea para separar las palabras considerando las comillas
func parseLine(linea string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder
	inComillas := false

	for _, char := range linea {
		switch {
		case char == '"' && !inComillas:
			inComillas = true
		case char == '"' && inComillas:
			//Indica la comilla de cierre entonces la cadena armada se almacena, luego se resetea
			inComillas = false
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		case unicode.IsSpace(char) && !inComillas:
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		default:
			currentToken.WriteRune(char)
		}
	}

	// Añadir el último token si existe
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	// Comprobar si hay un número impar de comillas
	if inComillas {
		return nil, errors.New("error de sintaxis: comillas no cerradas")
	}

	return tokens, nil
}

func recortarTokens(tokens []string) []string {
	fmt.Println("recortando tokens", tokens)
	for i, token := range tokens {
		if token == "<newline>" {
			return tokens[i+1:] // Devuelve los tokens después de <newline>
		}
	}
	return []string{} // Si no encuentra <newline>, retorna un array vacío
}
