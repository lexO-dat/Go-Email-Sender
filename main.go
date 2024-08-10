package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"email-api/mail"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}
}

// Estructura para recibir los datos del correo
type EmailRequest struct {
	Mail    string `json:"mail"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Habilitar CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Handler para enviar el correo
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Habilitar CORS para todas las solicitudes

	// Manejar las solicitudes OPTIONS para CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Asegurarse de que sea un POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el JSON de la solicitud
	var emailReq EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		http.Error(w, "Error al procesar el JSON", http.StatusBadRequest)
		return
	}

	// Construir el contenido del correo
	content := fmt.Sprintf(`
		<h1>%s</h1>
		<h2>%s</h2>
		<p>%s</p>
	`, emailReq.Mail, emailReq.Subject, emailReq.Body)

	// Crear el remitente usando el paquete mail
	sender := mail.NewGmailSender(
		os.Getenv("EMAIL_SENDER_NAME"),
		os.Getenv("EMAIL_SENDER_ADDRESS"),
		os.Getenv("EMAIL_SENDER_PASSWORD"),
	)

	// Enviar el correo
	to := []string{os.Getenv("DESTINATION_EMAIL")}
	attachFiles := []string{} // Puedes agregar archivos si es necesario

	err = sender.SendEmail(emailReq.Subject, content, to, nil, nil, attachFiles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al enviar el correo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Correo enviado exitosamente"))
}

func main() {
	http.HandleFunc("/send-email", sendEmailHandler)

	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
