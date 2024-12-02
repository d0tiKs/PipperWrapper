package main

import (
	//"flag"

	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
)

const (
	// CONNECTION PROPERTIES
	D_PORT         = 18080
	D_LISTENING_IP = "127.0.0.1"

	//WEBSOCKET SETTINGS
	D_VOICE_PACKET_BUFFER = 4096

	// PIPER SETTINGS
	D_MODEL             = "en"
	D_TEMP_FILE_DIR     = "/var/tmp/tts/"
	D_AUDIO_FILE_FORMAT = "wav"
)

var languageModels = map[string]string{
	"en": "/usr/local/share/piper/voices/en_US-hfc_female-medium.onnx",
	"fr": "/usr/local/share/piper/voices/fr_FR-upmc-medium.onnx",
	// Add more language models as needed
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TTSRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	Voice    string `json:"voice,omitempty"`
}

func ttsToHost(conn *websocket.Conn, modelPath string, text string) {
	ctx, _ := context.WithCancel(context.Background())

	// Create Piper command
	piperCmd := exec.CommandContext(ctx, "piper", "--model", modelPath, "--output-raw")

	// Create aplay command
	aplayCmd := exec.CommandContext(ctx, "aplay", "-r", "22050", "-f", "S16_LE", "-t", "raw")

	// Create pipes
	piperStdin, err := piperCmd.StdinPipe()
	if err != nil {
		log.Println("Piper stdin pipe error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not create stdin pipe to piper"})
		return
	}

	piperStdout, err := piperCmd.StdoutPipe()
	if err != nil {
		log.Println("Piper stdout pipe error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not create stdout pipe from piper"})
		return
	}

	aplayStdin, err := aplayCmd.StdinPipe()
	if err != nil {
		log.Println("aplay stdin pipe error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not create stdin pipe to aplay"})
		return
	}

	// Start commands
	if err := piperCmd.Start(); err != nil {
		log.Println("Piper start error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not start piper"})
		return
	}

	if err := aplayCmd.Start(); err != nil {
		log.Println("aplay start error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not start aplay"})
		return
	}
	// Write text to piper
	go func() {
		defer piperStdin.Close()
		if _, err := piperStdin.Write([]byte(text)); err != nil {
			log.Println("Error writing to piper stdin:", err)
		}
	}()

	// Pipe piper output to aplay
	go func() {
		defer aplayStdin.Close()
		if _, err := io.Copy(aplayStdin, piperStdout); err != nil {
			log.Println("Error copying from piper stdout to aplay stdin:", err)
		}
	}()

	// Wait for commands to complete
	go func() {
		defer conn.WriteJSON(map[string]string{"status": "success"})
		if err := aplayCmd.Wait(); err != nil {
			log.Println("aplay wait error:", err)
		}
		if err := piperCmd.Wait(); err != nil {
			log.Println("piper wait error:", err)
		}
		defer piperStdin.Close()
		defer aplayStdin.Close()
		defer piperStdout.Close()
	}()
}

func ttsToWebSocket(conn *websocket.Conn, modelPath string, text string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a temporary file to store the audio data
	tempFile, err := os.CreateTemp(D_TEMP_FILE_DIR, "audio-*."+D_AUDIO_FILE_FORMAT)
	if err != nil {
		log.Println("Error creating temporary file:", err)
		conn.WriteJSON(map[string]string{"error": "Could not create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Create Piper command to output to the temporary file
	piperCmd := exec.CommandContext(ctx, "piper", "-t", D_AUDIO_FILE_FORMAT, "--model", modelPath, "--output_file", tempFile.Name())

	// Create stdin pipe
	piperStdin, err := piperCmd.StdinPipe()
	if err != nil {
		log.Println("Piper stdin pipe error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not create stdin pipe to piper"})
		return
	}
	defer piperStdin.Close()

	// Start the piper command
	if err := piperCmd.Start(); err != nil {
		log.Println("Piper start error:", err)
		conn.WriteJSON(map[string]string{"error": "Could not start piper"})
		return
	}

	// Write text to piper
	go func() {
		defer piperStdin.Close()
		if _, err := piperStdin.Write([]byte(text)); err != nil {
			log.Println("Error writing to piper stdin:", err)
		}
	}()

	// Wait for the piper command to complete
	if err := piperCmd.Wait(); err != nil {
		log.Println("Piper wait error:", err)
	}

	// Read the audio data from the temporary file and send it over WebSocket
	audioData, err := io.ReadAll(tempFile)
	if err != nil {
		log.Println("Error reading audio data from temporary file:", err)
		conn.WriteJSON(map[string]string{"error": "Could not read audio data from temporary file"})
		return
	}

	// Encode the audio data using Base64
	encodedAudioData := base64.StdEncoding.EncodeToString(audioData)

	// Send the encoded audio data over WebSocket
	if err := conn.WriteMessage(websocket.TextMessage, []byte(encodedAudioData)); err != nil {
		log.Println("Error writing to WebSocket:", err)
		return
	}

	if err := conn.WriteJSON(map[string]string{"status": "success"}); err != nil {
		log.Println("Error sending success status:", err)
	}
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

		var req_body = fmt.Sprintf("body :\n%s", req)
		log.Println(req_body)

		modelPath := ""
		if req.Language == "" {
			// Default model
			modelPath = languageModels[D_MODEL]
		} else {
			// Determine the model path based on the language
			modelPath = languageModels[req.Language]
		}

		if modelPath == "" {
			log.Printf("No model found for language: %s", req.Language)
			conn.WriteJSON(map[string]string{"error": "Unsupported language"})
			continue
		}

		ttsToHost(conn, modelPath, req.Text)
		//ttsToWebSocket(conn, modelPath, req.Text)
	}
}

/* type ServerConfig struct {
	port       int
	listen_ip  string
	model_path string
}

var config *ServerConfig

func initCommandArguments() {
	// Define configuration struct to store our values
	config = &ServerConfig{}

	// Create a new FlagSet for our server command
	arguments := flag.NewFlagSet("server", flag.ExitOnError)

	// Port flags
	arguments.IntVar(&config.port, "p", D_PORT, "Port to run the server on")
	arguments.IntVar(&config.port, "port", D_PORT, "Port to run the server on")

	// Listen IP flags
	arguments.StringVar(&config.listen_ip, "l", D_LISTENING_IP, "IP address to listen on")
	arguments.StringVar(&config.listen_ip, "i", D_LISTENING_IP, "IP address to listen on")
	arguments.StringVar(&config.listen_ip, "listen_ip", D_LISTENING_IP, "IP address to listen on")

	// Model flags
	arguments.StringVar(&config.model_path, "m", D_MODEL_PATH, "Path to the onnx model to use")
	arguments.StringVar(&config.model_path, "model", D_MODEL_PATH, "Path to the onnx model to useModel to use")

	const usage = `Usage:
piper-websocket-tts-server [--port PORT] [--listen_ip IP] [--model PATH]

Options:
-p, --port PORT             Port to run the server on. Default is 8080.
-l, --listen_ip IP          IP address to listen on. Default is 0.0.0.0.
-m, --model PATH            Path to the onnx model to use. Default is ./model.onnx.
`

	arguments.Usage = func() { fmt.Print(usage) }

	// Parse the flags
	if err := arguments.Parse(os.Args[1:]); err != nil {
		arguments.Usage()
		os.Exit(1)
	}
} */

type Options struct {
	Port       int    `short:"p" long:"port" description:"Port to listen on"`
	ListenIP   string `short:"l" long:"listen" description:"IP address to listen on"`
	ConfigFile string `short:"c" long:"config" description:"Path to config file"`
}

func initOptions() {
	// Create an Options instance to hold our flags
	var opts Options

	// Create a new parser
	parser := flags.NewParser(&opts, flags.Default)

	defaultValues := fmt.Sprintf("\nHelp Options:\n\t-h, --help'\tShow this help message\nDefault Values:\n\tport: %d\n\tlistening ip: %s\n\tmodel: %s\n", D_PORT, D_LISTENING_IP, D_MODEL)

	parser.Usage = defaultValues
	// Parse the command line arguments
	_, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:

			//fmt.Print(defaultValues)
			os.Exit(1)
		}
	}
}

func main() {
	//initCommandArguments()
	//initOptions()
	http.HandleFunc("/tts", handleTTS)

	fmt.Printf("Server starting on port %d\n", D_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", D_PORT), nil))
}
