package comandos

import (
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseMkdisk(tokens []string) (*objs.DISK, []string, error) {
	cmd := &objs.DISK{}
	var contador, Size int
	var unit string
	basePath := "/home/david/MIA/P1/" // Ruta base donde se crearán los discos

	fmt.Println(tokens)
	for _, token := range tokens {
		if token == "<newline>" {
			contador += 1
			break
		}

		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, tokens[contador:], fmt.Errorf("formato de parámetro inválido: %s", token)
		}
		key, value := strings.ToLower(parts[0]), parts[1]

		key = strings.ToLower(key)
		switch key {
		case "-size":
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, tokens[contador:], errors.New("el tamaño debe ser un número entero positivo")
			}
			Size = size
		case "-unit":
			value = strings.ToUpper(value)
			if value != "K" && value != "M" {
				return nil, tokens[contador:], errors.New("la unidad debe ser K o M")
			}
			unit = value
		case "-fit":
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, tokens[contador:], errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.Fit = value
		default:
			return nil, tokens[contador:], fmt.Errorf("parámetro desconocido: %s", key)
		}
		contador += 1
	}

	if Size == 0 {
		return nil, tokens[contador:], errors.New("faltan parámetros requeridos: -size")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if unit == "" {
		unit = "M"
	}
	if cmd.Fit == "" {
		cmd.Fit = "FF"
	}

	// Generar nombre de archivo secuencial (A.dsk, B.dsk, etc.)
	diskName := generateNextDiskName()
	cmd.Path = filepath.Join(basePath, diskName)

	// Crear directorio si no existe
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, tokens[contador:], fmt.Errorf("error al crear directorio: %v", err)
	}

	// Llama a la función CreateBinaryFile del paquete disk para crear el archivo binario
	cmd.Size, _ = estructuras.ConvertToBytes(Size, unit)
	cmd.FreeSpace, _ = estructuras.ConvertToBytes(Size, unit)
	err := estructuras.CreateBinaryFile(cmd.Path, cmd.Size)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, tokens[contador:], err
	}
	cmd.Name = diskName
	objs.Discos = append(objs.Discos, *cmd)

	CrearMbr(cmd)
	fmt.Println("Se creo el disco")
	util.Respuestas = append(util.Respuestas, fmt.Sprintf("Disco %s creado con exito!", cmd.Path))
	return cmd, tokens[contador:], nil
}

// generateNextDiskName genera el siguiente nombre de disco en secuencia alfabética
func generateNextDiskName() string {
	letra := util.EliminarPrimeraLetra()
	letraG := strings.ToUpper(letra)
	return letraG + ".dsk"
}

// Resto del código permanece igual...
func CrearMbr(cmd *objs.DISK) error {
	mbr := &objs.MBR{}

	binary.LittleEndian.PutUint32(mbr.Size[:], uint32(cmd.Size))

	//Ingresa la fecha y hora actual
	now := time.Now()
	yearBase := 2000
	daysSinceBase := (now.Year()-yearBase)*365 + now.YearDay()
	secondsInDay := now.Hour()*3600 + now.Minute()*60 + now.Second()
	dayFraction := float32(secondsInDay) / 86400.0
	fechaHora := float32(daysSinceBase) + dayFraction
	// Convertir a bytes
	binary.LittleEndian.PutUint32(mbr.Fecha[:], math.Float32bits(fechaHora))

	// Generar y asignar un signature aleatorio
	binary.LittleEndian.PutUint32(mbr.Signature[:], GenerateDiskID())

	//Primer caracter de fit
	copy(mbr.Fit[:], string(cmd.Fit[0]))

	if err := mbr.WriteToFile(cmd.Path); err != nil {
		fmt.Println("Error al escribir el archivo:", err)
		return err
	}

	fmt.Println("Mbr creado con exito")

	return nil
}

func GenerateDiskID() uint32 {
	rand.Seed(time.Now().UnixNano())
	id := rand.Uint32()
	fmt.Println("ID del disco generado:", id, "---------------------------------")
	return id
}
