package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func predictHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "POST only", http.StatusMethodNotAllowed)
    return
  }

  // Parse form parameter "smiles"
  if err := r.ParseForm(); err != nil {
    http.Error(w, "Bad form", http.StatusBadRequest)
    return
  }
  smiles := r.FormValue("smiles")
  if smiles == "" {
    http.Error(w, "Missing 'smiles'", http.StatusBadRequest)
    return
  }

  // Write input file for cfm-predict
  if err := os.WriteFile("input.txt", []byte(smiles), 0644); err != nil {
    http.Error(w, "Cannot write input.txt", http.StatusInternalServerError)
    return
  }

  // Run cfm-predict CLI
  cmd := exec.Command("cfm-predict",
    "input.txt",
    "0.001",
    "/trained_models_cfmid4.0/[M+H]+/param_output.log",
    "/trained_models_cfmid4.0/[M+H]+/param_config.txt",
    "0",        // no fragment annotation
    "output.txt",
    "1",        // apply postprocessing
    "0",        // do not suppress exceptions
  )
  output, err := cmd.CombinedOutput()
  if err != nil {
    http.Error(w, fmt.Sprintf("cfm-predict error: %s\n%s", err, output),
      http.StatusInternalServerError)
    return
  }

  // Read and return the predicted spectrum
  result, err := os.ReadFile("output.txt")
  if err != nil {
    http.Error(w, "Cannot read output.txt", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "text/plain")
  w.Write(result)
}

func main() {
  http.HandleFunc("/predict", predictHandler)
  fmt.Println("CFM-ID HTTP wrapper listening on :5001")
  http.ListenAndServe(":5001", nil)
}