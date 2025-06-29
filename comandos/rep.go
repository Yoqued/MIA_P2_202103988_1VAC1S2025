package comandos

import (
	reportes "MIA_P2_202103988_1VAC1S2025/Reportes"
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type REP struct {
	name string
	path string
	id   string
	ruta string
}

func ParseRep(tokens []string) ([]string, error) {
	cmd := &REP{}
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
			cmd.name = value
		case "-path":
			cmd.path = value
		case "-id":
			if len(value) != 4 {
				return tokens[cont:], fmt.Errorf("El ID debe tener exactamente 4 caracteres")
			}

			// Extraemos las partes del ID original (ej: "A188")
			letra := string(value[0])  // "A"
			numero := string(value[1]) // "1"
			constante := value[2:]     // "88"

			// Reordenamos: constante + numero + letra → "88" + "1" + "A" = "881A"
			nuevoID := constante + numero + letra
			cmd.id = nuevoID
		case "-ruta":
			cmd.ruta = value
		default:
			return tokens[cont:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if cmd.name == "" {
		return tokens[contador:], errors.New("faltan parámetros requeridos: -name")
	}
	if cmd.path == "" {
		return tokens[contador:], errors.New("faltan parametros requeridos: -path")
	}
	if cmd.id == "" {
		return tokens[contador:], errors.New("faltan parametros requeridos: -id")
	}
	fmt.Println("REEEP----", util.GruposId)
	if _, exists := util.GruposId[cmd.id]; !exists {
		return tokens[contador:], fmt.Errorf("El ID %s no existe", cmd.id)
	}

	path, err := objs.BuscarPathID(cmd.id)
	if err != nil {
		return nil, err
	}
	fmt.Println(path, cmd.id)
	mbr, _ := objs.ReadMbr(path)
	particion := mbr.BuscarIdParticion(cmd.id)
	fmt.Println(particion)
	startPart := int64(binary.LittleEndian.Uint32(particion.Start[:]))
	superblock, _ := estructuras.ReadSuperBlock(path, startPart)

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Reporte del ID: %s", cmd.id))

	switch cmd.name {
	case "mbr":
		fmt.Println("ENTRANDO A REPORTE MBR")
		err := reportes.ReporteMbr(mbr, cmd.path, path)
		if err != nil {
			return tokens[contador:], err
		}
	case "disk":
		fmt.Println("ENTRANDO A REPORTE DISK")
		err := reportes.ReporteDisk(cmd.id, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "inode":
		fmt.Println("ENTRANDO A REPORTE INODO")
		err := reportes.ReporteInodo(superblock, path, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "block":
		fmt.Println("ENTRANDO A REPORTE BLOCK")
		err := reportes.ReporteBlocks(superblock, path, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "bm_inode":
		fmt.Println("ENTRANDO A REPORTE BMINODO")
		err := reportes.ReporteBMInodo(superblock, path, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "bm_block":
		fmt.Println("ENTRANDO A REPORTE BMBLOCK")
		err := reportes.ReporteBMBLock(superblock, path, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "sb":
		fmt.Println("ENTRANDO A REPORTE SUPERBLOCK")
		err := reportes.ReporteSuperblock(superblock, cmd.path)
		if err != nil {
			return tokens[contador:], err
		}
	case "file":
		fmt.Println("ENTRANDO A REPORTE FILE")
		err := reportes.ReporteFile(cmd.id, cmd.path, cmd.ruta)
		if err != nil {
			return tokens[contador:], err
		}
	case "ls":
		fmt.Println("ENTRANDO A REPORTE LS")
		err := reportes.ReporteLs(cmd.id, cmd.path, cmd.ruta)
		if err != nil {
			return tokens[contador:], err
		}
	}

	return tokens[contador:], nil
}
