package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"fmt"
	"os"
	"strings"
)

func ReporteBMInodo(superblock *estructuras.SUPERBLOCK, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Abrir el archivo de disco
	file, err := os.Open(diskPath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo de disco: %v", err)
	}
	defer file.Close()

	// Calcular el número total de inodos
	totalInodes := superblock.CountInodos + superblock.FreeInodos

	// Obtener el contenido del bitmap de inodos
	var bitmapContent strings.Builder

	for i := int32(0); i < totalInodes; i++ {
		// Establecer el puntero
		_, err := file.Seek(int64(superblock.InodeBitmapStart+i), 0)
		if err != nil {
			return fmt.Errorf("error al establecer el puntero en el archivo: %v", err)
		}

		// Leer un byte (carácter '0' o '1')
		char := make([]byte, 1)
		_, err = file.Read(char)
		if err != nil {
			return fmt.Errorf("error al leer el byte del archivo: %v", err)
		}

		// Agregar el carácter al contenido del bitmap
		bitmapContent.WriteByte(char[0])
		bitmapContent.WriteString("\t")

		// Agregar un carácter de nueva línea cada 20 caracteres (20 inodos)
		if (i+1)%20 == 0 {
			bitmapContent.WriteString("\n")
		}
	}

	// Crear el archivo TXT
	txtFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error al crear el archivo TXT: %v", err)
	}
	defer txtFile.Close()

	// Escribir el contenido del bitmap en el archivo TXT
	_, err = txtFile.WriteString(bitmapContent.String())
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo TXT: %v", err)
	}

	fmt.Println("Archivo del bitmap de inodos generado:", path)
	util.Respuestas = append(util.Respuestas, "Reporte de Bitmap de Inodos creado con exito")
	return nil
}

func ReporteBMBLock(superblock *estructuras.SUPERBLOCK, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Abrir el archivo de disco
	file, err := os.Open(diskPath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo de disco: %v", err)
	}
	defer file.Close()

	// Calcular el número total de bloques
	totalBlocks := superblock.CountBlocks + superblock.FreeBlocks

	// Obtener el contenido del bitmap de bloques
	var bitmapContent strings.Builder

	for i := int32(0); i < totalBlocks; i++ {
		// Establecer el puntero
		_, err := file.Seek(int64(superblock.BlockBitmapStart+i), 0)
		if err != nil {
			return fmt.Errorf("error al establecer el puntero en el archivo: %v", err)
		}

		// Leer un byte (carácter '0' o '1')
		char := make([]byte, 1)
		_, err = file.Read(char)
		if err != nil {
			return fmt.Errorf("error al leer el byte del archivo: %v", err)
		}

		// Agregar el carácter al contenido del bitmap
		bitmapContent.WriteByte(char[0])
		bitmapContent.WriteString("\t")

		// Agregar un carácter de nueva línea cada 20 caracteres (20 inodos)
		if (i+1)%20 == 0 {
			bitmapContent.WriteString("\n")
		}
	}

	// Crear el archivo TXT
	txtFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error al crear el archivo TXT: %v", err)
	}
	defer txtFile.Close()

	// Escribir el contenido del bitmap en el archivo TXT
	_, err = txtFile.WriteString(bitmapContent.String())
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo TXT: %v", err)
	}

	fmt.Println("Archivo del bitmap de bloques generado:", path)
	util.Respuestas = append(util.Respuestas, "Reporte de Bitmap de Bloques creado con exito")
	return nil
}
