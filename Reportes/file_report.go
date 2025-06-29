package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReporteFile(id string, path string, filePath string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	var salida string

	diskPath, err := objs.BuscarPathID(id)
	if err != nil {
		return err
	}

	mbr, _ := objs.ReadMbr(diskPath)
	particion := mbr.BuscarIdParticion(id)
	startPart := int64(binary.LittleEndian.Uint32(particion.Start[:]))
	superblock, _ := estructuras.ReadSuperBlock(diskPath, startPart)

	fileParts := strings.Split(filePath, "/")
	fileParts = fileParts[1:]

	// Buscamos el archivo segun la ruta especificada
	err = superblock.ReadFile(diskPath, fileParts, 0, &salida)
	if err != nil {
		return err
	}

	fileName := fileParts[len(fileParts)-1]

	if filepath.Ext(path) != "" {
		// Si ya tiene un archivo, recortamos esa parte del path
		path = filepath.Dir(path)
	}

	fileFullPath := filepath.Join(path, fileName)

	fmt.Println(fileFullPath)
	err = os.WriteFile(fileFullPath, []byte(salida), 0644)
	if err != nil {
		return err
	}

	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Reporte file %s creado con exito", filePath))

	return nil
}
