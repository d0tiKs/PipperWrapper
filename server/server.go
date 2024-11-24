package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TTSRequest struct {
	Text  string `json:"text"`
	Voice string `json:"voice,omitempty"`
}

func handleTTS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var req TTSRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var req_body = fmt.Sprintf("body :\n%s", req.Text)
		log.Println(req_body)

		// Create Piper command
		piperCmd := exec.Command("piper",
			"--model", "/usr/local/share/piper/voices/en_US-hfc_female-medium.onnx",
			"--output-raw")

		// Create aplay command
		aplayCmd := exec.Command("aplay", "-r", "22050", "-f", "S16_LE", "-t", "raw")

		// Create pipes
		piperStdin, err := piperCmd.StdinPipe()
		if err != nil {
			log.Println("Piper stdin pipe error:", err)
			conn.WriteJSON(map[string]string{"error": "Could not create stdin pipe to piper"})
			continue
		}

		piperStdout, err := piperCmd.StdoutPipe()
		if err != nil {
			log.Println("Piper stdout pipe error:", err)
			conn.WriteJSON(map[string]string{"error": "Could not create stdout pipe from piper"})
			continue
		}

		aplayStdin, err := aplayCmd.StdinPipe()
		if err != nil {
			log.Println("aplay stdin pipe error:", err)
			conn.WriteJSON(map[string]string{"error": "Could not create stdin pipe to aplay"})
			continue
		}

		// Start commands
		if err := piperCmd.Start(); err != nil {
			log.Println("Piper start error:", err)
			conn.WriteJSON(map[string]string{"error": "Could not start piper"})
			continue
		}

		if err := aplayCmd.Start(); err != nil {
			log.Println("aplay start error:", err)
			conn.WriteJSON(map[string]string{"error": "Could not start aplay"})
			continue
		}

		// Write text to piper
		go func() {
			defer piperStdin.Close()
			piperStdin.Write([]byte(req.Text))
		}()

		// Pipe piper output to aplay
		go func() {
			defer aplayStdin.Close()
			io.Copy(aplayStdin, piperStdout)
		}()

		// Wait for commands to complete
		go func() {
			aplayCmd.Wait()
			piperCmd.Wait()
			defer conn.WriteJSON(map[string]string{"status": "success"})
		}()
	}
}

func main() {
	http.HandleFunc("/tts", handleTTS)

	port := 18080
	fmt.Printf("Server starting on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
