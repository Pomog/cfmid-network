package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
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

	// Read parameters
	prob := r.FormValue("prob_thresh")
	if prob == "" {
		prob = "0.001"
	}

	smiles := r.FormValue("smiles")
	if smiles == "" {
		http.Error(w, "Missing 'smiles' parameter", http.StatusBadRequest)
		return
	}

	// Create input file in file-mode: need an ID and SMILES per line
	inFile, err := os.CreateTemp("", "cfm-in-*.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Create input file failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(inFile.Name())
	// Write ID (M1) and SMILES
	inFile.WriteString(fmt.Sprintf("M1 %s\n", smiles))
	inFile.Close()

	// Create temp output file
	outFile, err := os.CreateTemp("", "cfm-out-*.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Create output file failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(outFile.Name())
	outFile.Close()

	// Run CFM-ID in file mode
	cmd := exec.Command(
		"cfm-predict",
		inFile.Name(), // input file with ID SMILES
		prob,          // prob_thresh
		"/trained_models_cfmid4.0/cfmid4/[M+H]+/param_output.log",
		"/trained_models_cfmid4.0/cfmid4/[M+H]+/param_config.txt",
		"1",            // annotate_fragments = YES
		outFile.Name(), // output file
		"1",            // apply_postproc
		"0",            // suppress_exceptions
	)

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("cfm-predict failed: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := os.ReadFile(outFile.Name())
	if err != nil {
		http.Error(w, fmt.Sprintf("Read result failed: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/predict", predictHandler)

	srv := &http.Server{
		Addr:    ":5001",
		Handler: mux,
	}

	log.Println("🚀 CFM-ID wrapper starting on http://0.0.0.0:5001")
	log.Fatal(srv.ListenAndServe())
}
