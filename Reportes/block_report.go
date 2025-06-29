package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReporteBlocks(superblock *estructuras.SUPERBLOCK, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := GetFileNames(path)

	// Iniciar el contenido DOT con estilos
	dotContent := `digraph G {
		node [shape=plaintext]
		graph [bgcolor="#f0f0f0"];
		edge [color="#aaaaaa"];
	`

	// Iterar sobre cada bloque
	for i := int32(0); i < superblock.CountBlocks; i++ {
		blockStart := int64(superblock.BlockStart + (i * superblock.BlockSize))
		buffer := make([]byte, 64)
		file, err := os.Open(diskPath)
		if err != nil {
			return err
		}
		_, err = file.ReadAt(buffer, blockStart)
		if err != nil {
			return err
		}
		file.Close()

		switch estructuras.BlockTypes[int(i)] {
		case "CARPETA":
			bcarpeta := estructuras.BLOCK_CARPETA{}
			err = bcarpeta.Deserialize(diskPath, blockStart)
			if err != nil {
				return err
			}

			dotContent += fmt.Sprintf(
				`Block%d [label=<
					<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" BGCOLOR="#d9edf7">
						<TR><TD COLSPAN="2" BGCOLOR="#bce8f1">Bloque %d - Carpeta</TD></TR>
						<TR><TD><B>Nombre</B></TD><TD><B>Inodo</B></TD></TR>`,
				i, i)
			for _, content := range bcarpeta.Content {
				name := string(bytes.Trim(content.Name[:], "\x00"))
				dotContent += fmt.Sprintf("<TR><TD>%s</TD><TD>%d</TD></TR>", name, content.Inodo)
			}
			dotContent += "</TABLE>>];"

		case "APUNTADOR":
			bapuntador := &estructuras.BLOCK_APUNTADOR{}
			err = bapuntador.Deserialize(diskPath, blockStart)
			if err != nil {
				return err
			}

			// Crear una tabla para los apuntadores
			dotContent += fmt.Sprintf(
				`Block%d [label=<
					<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" BGCOLOR="#f9f9f9">
						<TR><TD COLSPAN="%d" BGCOLOR="#f2f2f2">Bloque Apuntadores %d</TD></TR>
						<TR>`, i, len(bapuntador.Pointers), i)

			// Agregar los apuntadores en filas con un máximo de 4 columnas
			for j, pointer := range bapuntador.Pointers {
				if j%4 == 0 && j != 0 {
					dotContent += "</TR><TR>" // Cierra la fila y comienza una nueva
				}
				dotContent += fmt.Sprintf("<TD>%d</TD>", pointer)
			}
			dotContent += "</TR></TABLE>>];"

		case "CONTARCHIVO":
			bcontarchivo, err := estructuras.ReadCont_Archivo(diskPath, blockStart)
			if err != nil {
				return err
			}

			content := string(bytes.Trim(bcontarchivo.Content[:], "\x00"))
			content = strings.ReplaceAll(content, "\n", "<BR/>")

			dotContent += fmt.Sprintf(
				`Block%d [label=<
					<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" BGCOLOR="#f2f2f2">
						<TR><TD>Bloque Archivo %d</TD></TR>
						<TR><TD>%s</TD></TR>
					</TABLE>>];`, i, i, content)
		}
	}

	// Cerrar el contenido DOT
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

	util.Respuestas = append(util.Respuestas, "Reporte de Bloques creado con éxito")
	return nil
}
