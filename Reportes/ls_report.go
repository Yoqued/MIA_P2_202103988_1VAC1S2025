package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ReporteLs(id string, path string, filePath string) error {
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}
	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := GetFileNames(path)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
		node [shape=none margin=0]
		graph [nodesep=0.5]

		table [
			label=<
			<table border="1" cellborder="1" cellspacing="0" cellpadding="4">
				<tr>
					<td><b>Permisos</b></td>
					<td><b>Owner</b></td>
					<td><b>Grupo</b></td>
					<td><b>Size (en Bytes)</b></td>
					<td><b>Fecha</b></td>
					<td><b>Hora</b></td>
					<td><b>Tipo</b></td>
					<td><b>Name</b></td>
				</tr>
	`

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
	if len(fileParts) == 1 && fileParts[0] == "" {
		fileParts = []string{}
	}

	_, err = createLs(superblock, diskPath, 0, fileParts, id, &dotContent)
	if err != nil {
		return err
	}

	dotContent += `
				</table>
				>
			];
		}
		`
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

	fmt.Println("Imagen ls generada:", outputImage)
	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Reporte Ls %s creado con exito", filePath))
	return nil
}

func createLs(sb *estructuras.SUPERBLOCK, path string, inodeIndex int32, parentsDir []string, id string, dotContent *string) (bool, error) {
	// Crear un nuevo inodo
	inode := &estructuras.INODO{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.InodeStart+(inodeIndex*sb.InodeSize)))
	if err != nil {
		return false, err
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.Tipo[0] == '1' {
		return false, nil
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.Blocks {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &estructuras.BLOCK_CARPETA{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.BlockStart+(blockIndex*sb.BlockSize))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return false, err
		}
		inode.Print()
		fmt.Println("len", len(parentsDir), parentsDir)
		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.Content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.Content[indexContent]

			// Si el contenido está vacío, salir
			if content.Inodo == -1 {
				break
			}

			// Si las carpetas padre no están vacías, debemos buscar la carpeta padre más cercana
			if len(parentsDir) != 0 {
				// Obtenemos la carpeta padre más cercana
				parentDir, err := util.First(parentsDir)
				if err != nil {
					return false, err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.Name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					_, err := createLs(sb, path, content.Inodo, util.RemoveElement(parentsDir, 0), id, dotContent)
					if err != nil {
						return false, err
					}
					return true, nil
				}
			} else {
				contentName := strings.Trim(string(content.Name[:]), "\x00 ")
				fmt.Println("BLOQUE", contentName)
				inode := &estructuras.INODO{}
				// Deserializar el inodo
				err := inode.Deserialize(path, int64(sb.InodeStart+(content.Inodo*sb.InodeSize)))
				if err != nil {
					return false, err
				}

				var tipo string
				usersList := util.GruposId[id]
				fmt.Println(usersList.UNames)
				fmt.Println(usersList.GNames)
				owner := usersList.UNames[inode.Uid-1]
				group := usersList.GNames[inode.Gid-1]
				mtime := time.Unix(int64(1695050848), 0)
				fecha := mtime.Format("2006/01/02") // YYYY-MM-DD
				hora := mtime.Format("15:04:05")
				perm := string(inode.Perm[:])
				if string(inode.Tipo[:]) == "1" {
					tipo = "Archivo"
				} else {
					tipo = "Carpeta"
				}
				permisos, _ := ConvertToRWX(perm)

				*dotContent += fmt.Sprintf(`
				<tr>
					<td>%s</td> <!-- permisos -->
					<td>%s</td> <!-- owner -->
					<td>%s</td> <!-- grupo -->
					<td>%d</td> <!-- size -->
					<td>%s</td> <!-- fecha -->
					<td>%s</td> <!-- hora -->
					<td>%s</td> <!-- tipo -->
					<td>%s</td> <!-- nombre -->
				</tr>
				`, permisos, owner, group, inode.Size, fecha, hora, tipo, contentName)
			}
		}
	}
	return false, nil
}

func ConvertToRWX(perm string) (string, error) {
	if len(perm) != 3 {
		return "", fmt.Errorf("los permisos deben ser una cadena de tres dígitos")
	}

	var rwxPerms string

	for i := 0; i < 3; i++ {
		switch perm[i] {
		case '7':
			rwxPerms += "rwx"
		case '6':
			rwxPerms += "rw-"
		case '5':
			rwxPerms += "r-x"
		case '4':
			rwxPerms += "r--"
		case '3':
			rwxPerms += "-wx"
		case '2':
			rwxPerms += "-w-"
		case '1':
			rwxPerms += "--x"
		case '0':
			rwxPerms += "---"
		default:
			return "", fmt.Errorf("permiso no válido: %c", perm[i])
		}
		if i != 2 && perm[i] == '7' {
			rwxPerms += "-"
		}
	}

	return rwxPerms, nil
}
