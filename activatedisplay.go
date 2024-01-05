package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/activate-display", activateDisplay)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func activateDisplay(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("xset", "dpms", "force", "on")
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	err := cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Display activated"))
}
