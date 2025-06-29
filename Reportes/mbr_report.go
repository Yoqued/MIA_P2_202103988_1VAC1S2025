package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func ReporteMbr(mbr *objs.MBR, outputPath string, diskPath string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(outputPath)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := GetFileNames(outputPath)

	// Convertir mbr.Tamano a entero
	mbrTamano := binary.LittleEndian.Uint32(mbr.Size[:])

	// Convertir mbr.Fecha a time.Time
	mbrFecha := time.Unix(int64(binary.LittleEndian.Uint32(mbr.Fecha[:])), 0)

	// Definir el contenido DOT con una tabla
	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext, fontname="Helvetica", fontsize=12]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="5" cellpadding="10">
                <tr><td bgcolor="lightblue" colspan="2"> REPORTE MBR </td></tr>
                <tr><td bgcolor="#b3cde3">mbr_tamano</td><td>%d bytes</td></tr>
                <tr><td bgcolor="#b3cde3">mrb_fecha_creacion</td><td>%s</td></tr>
                <tr><td bgcolor="#b3cde3">mbr_disk_signature</td><td>%d</td></tr>
            `, mbrTamano, mbrFecha.Format("2006-01-02 15:04:05"), mbr.Signature[:])

	// Agregar las particiones primarias y extendidas
	for i, part := range mbr.Particiones {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partName := strings.TrimRight(string(part.Name[:]), "\x00")
		if partName == "" {
			continue
		}
		// Convertir Status, Tipo y Fit a char
		partStatus := string(part.Status[:])
		partType := string(part.Tipo[:])
		partFit := string(part.Fit[:])

		// Convertir Start, Size, Correlative y Id a enteros
		partStart := binary.LittleEndian.Uint32(part.Start[:])
		partSize := binary.LittleEndian.Uint32(part.Size[:])
		partCorrelative := binary.LittleEndian.Uint32(part.Correlative[:])
		partId := binary.LittleEndian.Uint32(part.Id[:])

		// Determinar el color de fondo según el tipo de partición
		bgColor := "#fff496" // Soft yellow for primary partitions
		attrColor := "#fff2cc"
		if partType == "E" {
			bgColor = "#b8eea0" // Soft green for extended partitions
			attrColor = "#d9f2d9"
		}

		// Agregar la partición a la tabla
		dotContent += fmt.Sprintf(`
				<tr><td bgcolor="%s" colspan="2"> PARTICIÓN %d </td></tr>
				<tr><td bgcolor="%s">part_status</td><td>%s</td></tr>
				<tr><td bgcolor="%s">part_type</td><td>%s</td></tr>
				<tr><td bgcolor="%s">part_fit</td><td>%s</td></tr>
				<tr><td bgcolor="%s">part_start</td><td>%d (byte)</td></tr>
				<tr><td bgcolor="%s">part_size</td><td>%d bytes</td></tr>
				<tr><td bgcolor="%s">part_correlative</td><td>%d</td></tr>
				<tr><td bgcolor="%s">part_id</td><td>%d</td></tr>
				<tr><td bgcolor="%s">part_name</td><td>%s</td></tr>
			`, bgColor, i+1, attrColor, partStatus, attrColor, partType, attrColor, partFit, attrColor, partStart, attrColor, partSize, attrColor, partCorrelative, attrColor, partId, attrColor, partName)

		// Si es una partición extendida, agregar las particiones lógicas
		if partType == "E" {
			dotContent, err = LogicalPartitionsDot(diskPath, int(partStart), dotContent)
			if err != nil {
				return err
			}
		}
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>] }"

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	fmt.Println("Imagen de la tabla generada:", outputImage)
	util.Respuestas = append(util.Respuestas, "Reporte de MBR creado con exito")
	return nil
}

func LogicalPartitionsDot(path string, start int, dotContent string) (string, error) {
	ebr := &objs.EBR{}
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	currentOffset := int64(start)

	for {
		// Mueve el puntero del archivo a la posición del EBR actual
		_, err := file.Seek(currentOffset, 0)
		if err != nil {
			return "", fmt.Errorf("error al mover el puntero del archivo: %v", err)
		}

		// Lee el EBR desde el archivo
		if err := binary.Read(file, binary.LittleEndian, ebr); err != nil {
			return "", fmt.Errorf("error al leer el EBR: %v", err)
		}

		// Convertir los campos del EBR a sus respectivos tipos
		partName := strings.TrimRight(string(ebr.Name[:]), "\x00")
		partStatus := string(ebr.Status[:])
		partFit := string(ebr.Fit[:])
		partStart := binary.LittleEndian.Uint32(ebr.Start[:])
		partSize := binary.LittleEndian.Uint32(ebr.Size[:])
		nextOffset := int64(binary.LittleEndian.Uint32(ebr.Next[:]))
		if partName == "" {
			break
		}

		// Agregar la información de la partición lógica con un color suave coral
		dotContent += fmt.Sprintf(`
			<tr><td bgcolor="#fe9c91" colspan="2"> PARTICIÓN LÓGICA </td></tr>
			<tr><td bgcolor="#fddcdc">part_name</td><td>%s</td></tr>
			<tr><td bgcolor="#fddcdc">part_status</td><td>%s</td></tr>
			<tr><td bgcolor="#fddcdc">part_next</td><td>%d</td></tr>
			<tr><td bgcolor="#fddcdc">part_fit</td><td>%s</td></tr>
			<tr><td bgcolor="#fddcdc">part_start</td><td>%d (byte)</td></tr>
			<tr><td bgcolor="#fddcdc">part_size</td><td>%d bytes</td></tr>
		`, partName, partStatus, nextOffset, partFit, partStart, partSize)

		// Verifica si hay otro EBR
		if nextOffset == int64(0xFFFFFFFF) {
			break // No hay más EBRs, sal del bucle
		}

		currentOffset = nextOffset
	}

	return dotContent, nil
}

func CreateParentDirs(path string) error {
	dir := filepath.Dir(path)
	// os.MkdirAll no sobrescribe las carpetas existentes, solo crea las que no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear las carpetas padre: %v", err)
	}
	return nil
}

func GetFileNames(path string) (string, string) {
	dir := filepath.Dir(path)
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dotFileName := filepath.Join(dir, baseName+".dot")
	outputImage := path
	return dotFileName, outputImage
}
