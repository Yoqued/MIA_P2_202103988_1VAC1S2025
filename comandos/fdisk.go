package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseFdisk(tokens []string) (*objs.PARTICION, []string, error) {
	cmd := &objs.PARTICION{}
	var unit, driveLetter, tipo, fit, delete string
	var contador, Sizes, add int

	for cont, token := range tokens {
		if token == "<newline>" {
			contador += 1
			break
		}

		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, tokens[cont:], fmt.Errorf("formato de parámetro inválido: %s", token)
		}
		key, value := strings.ToLower(parts[0]), parts[1]

		switch key {
		case "-size":
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, tokens[cont:], errors.New("el tamaño debe ser un número entero positivo")
			}
			Sizes = size
		case "-unit":
			value = strings.ToUpper(value)
			if value != "B" && value != "K" && value != "M" {
				return nil, tokens[cont:], errors.New("la unidad debe ser K o M")
			}
			unit = value
		case "-driveletter":
			if len(value) != 1 || !isLetter(value[0]) {
				return nil, tokens[cont:], errors.New("la letra del disco debe ser un carácter de la A-Z")
			}
			driveLetter = strings.ToUpper(value)
		case "-fit":
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, tokens[cont:], errors.New("el ajuste debe ser BF, FF o WF")
			}
			fit = string(value[0])
		case "-type":
			value = strings.ToUpper(value)
			if value != "P" && value != "E" {
				return nil, tokens[cont:], errors.New("el tipo debe ser P, E o L")
			}
			tipo = value
			copy(cmd.Tipo[:], tipo)
		case "-name":
			fmt.Println(value)
			copy(cmd.Name[:], value)
		case "-delete":
			value = strings.ToUpper(value)
			if value != "FAST" && value != "FULL" {
				return nil, tokens[cont:], errors.New("el tipo de delete debe ser full o fast")
			}
			delete = value
		case "-add":
			size, err := strconv.Atoi(value)
			if err != nil {
				return nil, tokens[cont:], err
			}
			add = size
		default:
			return nil, tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if delete != "" && add != 0 {
		return nil, tokens[contador:], errors.New("no se puede hacer un delete y un add a la vez")
	}
	if driveLetter == "" {
		return nil, tokens[contador:], errors.New("faltan parametros requeridos: -driveletter")
	}
	if objs.IsEmptyByte(cmd.Name[:]) {
		return nil, tokens[contador:], errors.New("faltan parametros requeridos: -name")
	}

	// Construir la ruta del archivo basado en la letra del disco
	path := filepath.Join("/home/david/MIA/P1/", fmt.Sprintf("%s.dsk", driveLetter))

	if delete != "" {
		name := string(bytes.Trim(cmd.Name[:], "\x00"))
		err := deleteParticion(path, name, delete)
		if err != nil {
			return nil, tokens[contador:], err
		}
		return cmd, tokens[contador:], nil
	}

	if Sizes == 0 {
		return nil, tokens[contador:], errors.New("faltan parámetros requeridos: -size")
	}

	if unit == "" {
		unit = "K"
	}
	if fit == "" {
		fit = "W"
	}

	tamano, err := estructuras.ConvertToBytes(Sizes, unit)
	if err != nil {
		return nil, tokens[contador:], err
	}

	disk, err := objs.BuscarDisco(path)
	if err != nil {
		return nil, tokens[contador:], err
	}

	if disk.FreeSpace < tamano {
		return nil, tokens[contador:], errors.New("Espacio Insuficiente en el disco")
	} else {
		disk.FreeSpace -= tamano
	}

	binary.LittleEndian.PutUint32(cmd.Size[:], uint32(tamano))
	copy(cmd.Status[:], "0")
	copy(cmd.Fit[:], fit)

	mbr, err := objs.ReadMbr(path)
	if err != nil {
		fmt.Println("---------------ERROR2-------------", err)
		return nil, tokens[contador:], err
	}
	binary.LittleEndian.PutUint32(cmd.Start[:], uint32(CalcularStart(mbr)))

	if tipo == "P" {
		err = mbr.AgregarParticion(*cmd)
		if err != nil {
			fmt.Println("---------------ERROR-------------")
			fmt.Println(err)
			return nil, tokens[contador:], err
		}
		mbr.WriteToFile(path)
	} else if tipo == "E" {
		if !mbr.VerificarExtendida() {
			ebr := &objs.EBR{}
			copy(ebr.Fit[:], cmd.Fit[:])
			binary.LittleEndian.PutUint32(ebr.Next[:], uint32(0xFFFFFFFF))
			err = mbr.AgregarParticion(*cmd)
			if err != nil {
				fmt.Println("---------------ERROR-------------")
				fmt.Println(err)
				return nil, tokens[contador:], err
			}
			mbr.WriteToFile(path)
			ebr.WriteToFile(path, int(binary.LittleEndian.Uint32(cmd.Start[:])))
		} else {
			return nil, tokens[contador:], errors.New("ya existe una particion extendida en este disco")
		}
	} else if tipo == "L" {
		if mbr.VerificarExtendida() {
			startExt := mbr.StartExtendida()
			ebr, err := objs.ReadEBRsFromFile(path, startExt)
			if err != nil {
				return nil, tokens[contador:], err
			}
			startLogica, err := objs.StartLogica(path, startExt)
			if err != nil {
				return nil, tokens[contador:], err
			}
			startNextEbr := startLogica + tamano
			copy(ebr.Status[:], "0")
			copy(ebr.Name[:], cmd.Name[:])
			binary.LittleEndian.PutUint32(ebr.Size[:], uint32(tamano))
			binary.LittleEndian.PutUint32(ebr.Start[:], uint32(startLogica))
			binary.LittleEndian.PutUint32(ebr.Next[:], uint32(startNextEbr))
			ebr.WriteToFile(path, startLogica-30)
			nextebr := &objs.EBR{}
			copy(nextebr.Fit[:], ebr.Fit[:])
			binary.LittleEndian.PutUint32(nextebr.Next[:], uint32(0xFFFFFFFF))
			nextebr.WriteToFile(path, startNextEbr)
		} else {
			return nil, tokens[contador:], errors.New("no existe una particion extendida en este disco para la logica")
		}
	}
	name := string(bytes.Trim(cmd.Name[:], "\x00"))
	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Particion %s creada con exito", name))
	return cmd, tokens[contador:], nil
}

func CalcularStart(mbr *objs.MBR) int {
	var start = 153
	particiones := mbr.Particiones
	for _, particion := range particiones {
		size := int(binary.LittleEndian.Uint32(particion.Size[:]))
		startpart := int(binary.LittleEndian.Uint32(particion.Start[:]))
		if objs.IsEmptyByte(particion.Name[:]) {
			break
		}
		start = startpart + size
	}
	return start
}

func deleteParticion(path string, name string, tipo string) error {
	mbr, err := objs.ReadMbr(path)
	if err != nil {
		return err
	}

	particion := mbr.BuscarParticion(name)
	if particion == nil {
		return fmt.Errorf("la partición '%s' no fue encontrada", name)
	}

	if tipo == "FULL" {
		file, err := os.OpenFile(path, os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("no se pudo abrir el archivo: %v", err)
		}
		defer file.Close()

		start := int64(binary.LittleEndian.Uint32(particion.Start[:]))
		_, err = file.Seek(start, 0)
		if err != nil {
			return fmt.Errorf("no se pudo mover a la posición de inicio de la partición: %v", err)
		}

		size := int(binary.LittleEndian.Uint32(particion.Size[:]))
		zeroes := make([]byte, size)

		_, err = file.Write(zeroes)
		if err != nil {
			return fmt.Errorf("no se pudo sobrescribir la partición: %v", err)
		}
	}

	particion.Clear()
	err = mbr.WriteToFile(path)
	if err != nil {
		return err
	}

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Particion %s removida con exito, con el protocolo %s", name, tipo))

	return nil
}
