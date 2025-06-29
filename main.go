package main

import (
	analizar "MIA_P2_202103988_1VAC1S2025/Analizador"
	util "MIA_P2_202103988_1VAC1S2025/Util"
	"MIA_P2_202103988_1VAC1S2025/comandos"
	"MIA_P2_202103988_1VAC1S2025/estructuras"
	"MIA_P2_202103988_1VAC1S2025/objs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Define the CommandRequest and CommandResponse structs
type CommandRequest struct {
	Commands string `json:"commands"`
}

type CommandResponse struct {
	Output string `json:"output"`
}

// Estructura para la solicitud de login
type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	IdParticion string `json:"partitionId"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"` // Incluye el token si el login es exitoso
}

type DiscosResponse struct {
	Discos []objs.DISK `json:"discos"`
}

// Funcion que analizara los comandos
func executeCommands(commands string) string {
	var output string
	answer, errs := analizar.AnalizarTexto(commands)
	if len(errs) > 0 {
		// Si hay errores, convertimos cada uno a string
		var errMessages []string
		for _, err := range errs {
			errMessages = append(errMessages, err.Error())
		}
		// Unimos todos los errores en una cadena
		errorsText := strings.Join(errMessages, "\n")
		output = errorsText
	} else {
		output = "Comandos ejecutados con éxito!"
	}

	// Imprimir el contenido de myInterface si no es nil
	if answer != nil {
		output += "\nResultados:\n"

		switch v := answer.(type) {
		case []string:
			// Si myInterface es un slice de strings
			output += strings.Join(v, "\n")
		case string:
			// Si myInterface es un string
			output += v
		case []int:
			// Si myInterface es un slice de enteros
			for _, val := range v {
				output += fmt.Sprintf("%d\n", val)
			}
		default:
			output += fmt.Sprintf("%v", answer)
		}
	}

	return output
}

// HTTP handler for executing commands
func commandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var cmdRequest CommandRequest
		fmt.Println("ENTRO A ANALIZAR, Texto ingresado:")
		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decodificar el JSON en la estructura CommandRequest
		err = json.Unmarshal(body, &cmdRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Ejecutar los comandos y obtener el resultado
		output := executeCommands(cmdRequest.Commands)

		// Crear la respuesta
		response := CommandResponse{Output: output}

		// Convertir la respuesta a JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Establecer el tipo de contenido a JSON y devolver la respuesta
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Lógica de autenticación
func authenticate(username, password, partitionId string) (bool, string) {
	success, err := comandos.Login(util.USER{
		User: username,
		Pass: password,
		Id:   partitionId,
	})

	// Si hay un error, devolvemos el error en el mensaje
	if err != nil {
		fmt.Println("Error en el login:", err)
		return false, "Autenticación fallida: " + err.Error()
	}

	// Si el login es exitoso
	if success {
		fmt.Println("Autenticación exitosa")
		return true, "Login exitoso"
	}

	// Si no fue exitoso, devolver el mensaje de error
	return false, "Autenticación fallida: usuario o contraseña incorrectos"
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var loginRequest LoginRequest

		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decodificar el JSON en la estructura LoginRequest
		err = json.Unmarshal(body, &loginRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validar las credenci
		// Extraemos las partes del ID original (ej: "A188")
		letra := string(loginRequest.IdParticion[0])  // "A"
		numero := string(loginRequest.IdParticion[1]) // "1"
		constante := loginRequest.IdParticion[2:]     // "88"

		// Reordenamos: constante + numero + letra → "88" + "1" + "A" = "881A"
		nuevoID := constante + numero + letra
		success, token := authenticate(loginRequest.Username, loginRequest.Password, nuevoID)
		var response LoginResponse

		if success {
			response = LoginResponse{
				Message: "Login exitoso",
				Success: true,
				Token:   token,
			}
		} else {
			response = LoginResponse{
				Message: token,
				Success: false,
			}
		}

		// Convertir la respuesta a JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Establecer el tipo de contenido a JSON y devolver la respuesta
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// Reseteando el usuario logeado
		util.UserOnline = util.USER{}
		var response LoginResponse

		response = LoginResponse{
			Message: "Logout Exitoso",
			Success: true,
		}

		// Convertir la respuesta a JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Logout Exitoso")

		// Establecer el tipo de contenido a JSON y devolver la respuesta
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Handler para obtener la lista de discos
func discosHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Mandando discos")
	if r.Method == http.MethodGet {
		response := DiscosResponse{
			Discos: objs.Discos,
		}

		// Convertir la respuesta a JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Establecer el tipo de contenido a JSON y devolver la respuesta
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Se obtiene la ruta del disco seleccionado
func selectDiscoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := r.URL.Query().Get("path") // Obtener el parámetro 'path' de la URL
		if path == "" {
			http.Error(w, "El parámetro 'path' es requerido", http.StatusBadRequest)
			return
		}

		fmt.Println("Este es el path", path)
		_, particiones, err := util.GetMBR_Partitions(path)
		if err != nil {
			http.Error(w, "Problemas con extraer contenido del disco", http.StatusNotFound)
		}

		// Enviar el listado de particiones
		responseJSON, err := json.Marshal(particiones)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func selectPartitionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PARTITIONHANDLER")
	files := []estructuras.SHOWFILES{}
	if r.Method == http.MethodGet {
		startpart := r.URL.Query().Get("start") // Obtener el parámetro 'path' de la URL
		if startpart == "" {
			http.Error(w, "El parámetro 'start' es requerido", http.StatusBadRequest)
			return
		}
		startPart, err := strconv.Atoi(strings.TrimSpace(startpart))
		if err != nil {
			fmt.Println("Error de conversión:", err)
		}

		path := r.URL.Query().Get("path") // Obtener el parámetro 'path' de la URL
		if path == "" {
			http.Error(w, "El parámetro 'path' es requerido", http.StatusBadRequest)
			return
		}

		file := r.URL.Query().Get("fileparts") // Obtener el parámetro 'path' de la URL
		if file == "" {
			http.Error(w, "El parámetro 'path' es requerido", http.StatusBadRequest)
			return
		}

		fileParts := strings.Split(file, "/")
		fileParts = fileParts[1:]
		if len(fileParts) > 0 && fileParts[len(fileParts)-1] == "" {
			fileParts = fileParts[:len(fileParts)-1]
		}
		if len(fileParts) == 0 || (len(fileParts) == 1 && fileParts[0] == "") {
			fileParts = []string{}
		}

		fmt.Println("FILE", file)
		fmt.Println("FILEPARTS", fileParts)
		fmt.Println("Este es el start", startPart)
		fmt.Println("Este es el path", path)
		superblock, _ := estructuras.ReadSuperBlock(path, int64(startPart))

		// Validar si se extrajo un superbloque que si existe con mkfiles
		err = superblock.ValidateSuperBlock()
		if err != nil {
			fmt.Println("ERROR", err)
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		}

		err = superblock.ShowFiles(path, fileParts, 0, &files)
		if err != nil {
			fmt.Println("ERROR", err)
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		}

		// Enviar el listado de archivos
		responseJSON, err := json.Marshal(files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func selectContentHandler(w http.ResponseWriter, r *http.Request) {
	var cadena string
	if r.Method == http.MethodGet {
		startpart := r.URL.Query().Get("start") // Obtener el parámetro 'path' de la URL
		if startpart == "" {
			http.Error(w, "El parámetro 'start' es requerido", http.StatusBadRequest)
			return
		}
		startPart, err := strconv.Atoi(strings.TrimSpace(startpart))
		if err != nil {
			fmt.Println("Error de conversión:", err)
		}

		path := r.URL.Query().Get("path") // Obtener el parámetro 'path' de la URL
		if path == "" {
			http.Error(w, "El parámetro 'path' es requerido", http.StatusBadRequest)
			return
		}

		file := r.URL.Query().Get("fileparts") // Obtener el parámetro 'path' de la URL
		if file == "" {
			http.Error(w, "El parámetro 'path' es requerido", http.StatusBadRequest)
			return
		}

		fileParts := strings.Split(file, "/")
		fileParts = fileParts[1:]
		if len(fileParts) > 0 && fileParts[len(fileParts)-1] == "" {
			fileParts = fileParts[:len(fileParts)-1]
		}
		if len(fileParts) == 0 || (len(fileParts) == 1 && fileParts[0] == "") {
			fileParts = []string{}
		}

		fmt.Println("FILE", file)
		fmt.Println("FILEPARTS", fileParts)
		fmt.Println("Este es el start", startPart)
		fmt.Println("Este es el path", path)
		superblock, _ := estructuras.ReadSuperBlock(util.UserOnline.Path, int64(startPart))
		err = superblock.ReadFile(path, fileParts, 0, &cadena)
		if err != nil {
			fmt.Println("ERROR", err)
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		}

		// Enviar el listado de archivos
		responseJSON, err := json.Marshal(cadena)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Wrapper to handle CORS
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Permitir orígenes
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			// Responder a la solicitud OPTIONS
			w.WriteHeader(http.StatusOK)
			return
		}

		// Llamar al handler original
		h(w, r)
	}
}

func InicializarDisks() {
	filePath := "./PR1_AC1.smia"

	content, err := readSmiaFile(filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	analizar.AnalizarTexto(string(content))
}

func readSmiaFile(path string) ([]byte, error) {
	// Abrir el archivo
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %v", err)
	}
	defer file.Close()

	// Leer todo el contenido del archivo
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo: %v", err)
	}

	return content, nil
}

// Main function to start the server
func main() {
	http.HandleFunc("/execute", withCORS(commandHandler))
	http.HandleFunc("/login", withCORS(loginHandler))
	http.HandleFunc("/logout", withCORS(logoutHandler))
	http.HandleFunc("/discos", withCORS(discosHandler))
	http.HandleFunc("/disco-select", withCORS(selectDiscoHandler))
	http.HandleFunc("/partition-select", withCORS(selectPartitionHandler))
	http.HandleFunc("/file-select", withCORS(selectPartitionHandler))
	http.HandleFunc("/content-select", withCORS(selectContentHandler))
	fmt.Println("Servidor iniciado en http://localhost:8080")
	InicializarDisks()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
