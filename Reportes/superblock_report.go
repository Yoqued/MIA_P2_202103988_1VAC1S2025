package reportes

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func ReporteSuperblock(superblock *estructuras.SUPERBLOCK, path string) error {
	// Crear las carpetas padre si no existen
	err := CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := GetFileNames(path)

	// Transformando tiempos
	mtime := time.Unix(int64(superblock.MontTime), 0).Format(time.RFC3339)
	umtime := time.Unix(int64(superblock.UMontTime), 0).Format(time.RFC3339)

	// Definir el contenido DOT con una tabla estilizada para el SUPERBLOCK
	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext, fontname="Helvetica", fontsize=12]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="5" cellpadding="10">
                <tr><td bgcolor="lightblue" colspan="2"> REPORTE SUPERBLOCK </td></tr>
                <tr><td bgcolor="lightgray">Filesystem</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">CountInodos</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">CountBlocks</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">FreeBlocks</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">FreeInodos</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">MontTime</td><td>%s</td></tr>
                <tr><td bgcolor="lightgray">UMontTime</td><td>%s</td></tr>
                <tr><td bgcolor="lightgray">CountMount</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">Magic</td><td>%d</td></tr>
                <tr><td bgcolor="lightgray">InodeSize</td><td>%d bytes</td></tr>
                <tr><td bgcolor="lightgray">BlockSize</td><td>%d bytes</td></tr>
                <tr><td bgcolor="lightgray">FirstFreeInode</td><td>%d (byte)</td></tr>
                <tr><td bgcolor="lightgray">FirstFreeBlock</td><td>%d (byte)</td></tr>
                <tr><td bgcolor="lightgray">InodeBitmapStart</td><td>%d (byte)</td></tr>
                <tr><td bgcolor="lightgray">BlockBitmapStart</td><td>%d (byte)</td></tr>
                <tr><td bgcolor="lightgray">InodeStart</td><td>%d (byte)</td></tr>
                <tr><td bgcolor="lightgray">BlockStart</td><td>%d (byte)</td></tr>
            </table>>] }`,
		superblock.Filesystem,
		superblock.CountInodos,
		superblock.CountBlocks,
		superblock.FreeBlocks,
		superblock.FreeInodos,
		mtime,
		umtime,
		superblock.CountMount,
		superblock.Magic,
		superblock.InodeSize,
		superblock.BlockSize,
		superblock.FirstFreeInode,
		superblock.FirstFreeBlock,
		superblock.InodeBitmapStart,
		superblock.BlockBitmapStart,
		superblock.InodeStart,
		superblock.BlockStart,
	)

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

	fmt.Println("Imagen de la tabla de SUPERBLOCK generada:", outputImage)
	util.Respuestas = append(util.Respuestas, "Reporte de SuperBloque creado con exito")
	return nil
}
