package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReporteDisk(id string, path string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := GetFileNames(path)

	diskPath, err := objs.BuscarPathID(id)
	if err != nil {
		return err
	}

	mbr, _ := objs.ReadMbr(diskPath)

	// Calcular el tamaño total del disco
	totalSize := int64(binary.LittleEndian.Uint32(mbr.Size[:]))

	// Iniciar la generación del contenido del archivo DOT
	dotContent := `digraph DiskStructure {`
	dotContent += `  node [shape=none, fontname="Helvetica,Arial,sans-serif"];`
	dotContent += `  rankdir=TB;`

	// Crear el nodo del disco
	dotContent += `  disk [label=<
					<table border='0' cellborder='1' cellspacing='0' cellpadding='10' style='rounded' bgcolor='#F5F5F5'>
					<tr>
					<td colspan='2' bgcolor='#0044CC'><font color='white' face='Verdana' point-size='20'><b>Reporte Disk</b></font></td>
					</tr>
					<tr>`

	// Agregar MBR con nuevo diseño
	mbrPercentage := float64(153) / float64(totalSize) * 100
	dotContent += fmt.Sprintf(`<td width='%d' bgcolor='#FF4500' align='center' style='border:2px solid black;'>
		<font color='black'><b>MBR</b><br/>%.5f%%</font></td>\n`, int(mbrPercentage), mbrPercentage)

	// Agregar particiones y espacio libre
	freeSpace := totalSize
	for i, particion := range mbr.Particiones {
		partitionSize := int(binary.LittleEndian.Uint32(particion.Size[:]))
		if partitionSize > 0 {
			partitionPercentage := float64(partitionSize) / float64(totalSize) * 100
			partitionType := string(particion.Tipo[0])
			partitionName := strings.TrimRight(string(particion.Name[:]), "\x00")
			color := getPartitionColor(partitionType)

			dotContent += fmt.Sprintf(`<td width='%d' bgcolor='%s' align='center' style='border:2px solid black;'>
			<font color='black'><b>%s %d</b><br/>%s<br/>%.2f%%</font></td>\n`,
				int(partitionPercentage), color, getPartitionTypeLabel(partitionType), i+1, partitionName, partitionPercentage)

			freeSpace -= int64(partitionSize)

			if partitionType == "E" {
				// Agregar EBRs y particiones lógicas con nuevo diseño
				var cont int
				startExt := mbr.StartExtendida()
				currentOffset := int64(startExt)

				for {
					cont++
					// Lee el EBR desde el archivo
					ebr := &objs.EBR{}
					err := ebr.Deserialize(diskPath, currentOffset)
					if err != nil {
						return err
					}
					dotContent += `<td width='50' bgcolor='#FFD700' align='center' style='border: 2px solid blue;'>
					<font color='black'><b>EBR</b></font></td>\n`

					nextOffset := int64(binary.LittleEndian.Uint32(ebr.Next[:]))
					if nextOffset == int64(0xFFFFFFFF) {
						break // No hay más EBRs, sal del bucle
					}

					partitionName := strings.TrimRight(string(ebr.Name[:]), "\x00")
					logicalSize := int(binary.LittleEndian.Uint32(ebr.Size[:]))
					logicalPercentage := (float64(logicalSize) / float64(totalSize)) * 100

					dotContent += fmt.Sprintf(`<td width='%d' bgcolor='#FFA07A' align='center' style='border: 4px solid green;'>
					<font color='black'><b>Logica %d</b><br/>%s<br/>%.3f%%</font></td>\n`, int(logicalPercentage), cont, partitionName, logicalPercentage)

					currentOffset = nextOffset
				}
			}
		}
	}

	// Agregar espacio libre restante con nuevo diseño
	if freeSpace > 0 {
		freePercentage := float64(freeSpace) / float64(totalSize) * 100
		dotContent += fmt.Sprintf(`<td width='%d' bgcolor='#FFFFFF' align='center' style='border:2px solid red;'>
		<font color='black'><b>Espacio Libre</b><br/>%.2f%%</font></td>\n`, int(freePercentage), freePercentage)
	}

	dotContent += `      </tr>`
	dotContent += `    </table>`
	dotContent += `  >];`
	dotContent += "}"

	// Crear el archivo DOT
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer dotFile.Close()

	// Escribir el contenido DOT en el archivo
	_, err = dotFile.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar la imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Imagen del disco generada:", outputImage)
	util.Respuestas = append(util.Respuestas, "Reporte de Disco creado con exito")
	return nil
}

// Función para asignar colores a las particiones
func getPartitionColor(partitionType string) string {
	switch partitionType {
	case "P":
		return "#87CEEB" // Azul claro para particiones primarias
	case "E":
		return "#32CD32" // Verde para particiones extendidas
	default:
		return "#D3D3D3" // Gris para cualquier otra
	}
}

// Función para obtener las etiquetas de los tipos de partición
func getPartitionTypeLabel(partitionType string) string {
	switch partitionType {
	case "P":
		return "Primaria"
	case "E":
		return "Extendida"
	default:
		return "Desconocida"
	}
}
