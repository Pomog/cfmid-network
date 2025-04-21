package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form data", http.StatusBadRequest)
		return
	}
	smiles := r.FormValue("smiles")
	if smiles == "" {
		http.Error(w, "Missing 'smiles' parameter", http.StatusBadRequest)
		return
	}

	// Create temp files (unique per request)
	inFile, err := os.CreateTemp("", "cfm-in-*.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Create input file failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(inFile.Name())
	inFile.WriteString(smiles)
	inFile.Close()

	outFile, err := os.CreateTemp("", "cfm-out-*.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Create output file failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(outFile.Name())
	outFile.Close()

	// Timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cfm-predict",
		inFile.Name(),
		"0.001",
		"/trained_models_cfmid4.0/cfmid4/[M+H]+/param_output.log",
		"/trained_models_cfmid4.0/cfmid4/[M+H]+/param_config.txt",
		"0", // no fragment annotation
		outFile.Name(),
		"1", // postprocessing
		"0", // suppress exceptions
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Stdout failed: %v", err), http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Stderr failed: %v", err), http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Start failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Stream stderr first (warnings), then stdout
	go io.Copy(w, stderr)
	io.Copy(w, stdout)

	if err := cmd.Wait(); err != nil {
		log.Printf("cfm-predict error: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/predict", predictHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	srv := &http.Server{
		Addr:    ":5001",
		Handler: mux,
	}

	log.Println("🚀 CFM-ID wrapper starting on http://0.0.0.0:5001")
	log.Fatal(srv.ListenAndServe())
}
