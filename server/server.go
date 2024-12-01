package main

import (
	//"flag"
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

	// PIPER SETTINGS
	D_MODEL_PATH = "/usr/local/share/piper/voices/en_US-hfc_female-medium.onnx"
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
			modelPath = D_MODEL_PATH
		} else {
			// Determine the model path based on the language
			modelPath = languageModels[req.Language]
		}

		if modelPath == "" {
			log.Printf("No model found for language: %s", req.Language)
			conn.WriteJSON(map[string]string{"error": "Unsupported language"})
			continue
		}

		// Create Piper command
		piperCmd := exec.Command("piper",
			"--model", modelPath,
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

	defaultValues := fmt.Sprintf("\nHelp Options:\n\t-h, --help'\tShow this help message\nDefault Values:\n\tport: %d\n\tlistening ip: %s\n\tmodel: %s\n", D_PORT, D_LISTENING_IP, D_MODEL_PATH)

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
