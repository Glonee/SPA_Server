package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type Config struct {
	Port     string `json:"port"`
	Path     string `json:"path"`
	Homepage string `json:"homepage"`
}
type SinglePageFS struct {
	http.FileSystem
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	spfs := SinglePageFS{http.Dir(config.Path)}
	http.Handle(config.Homepage, http.StripPrefix(config.Homepage, http.FileServer(spfs)))
	fmt.Println("You can view it in the browser.")
	fmt.Println("Local: http://localhost:" + config.Port + config.Homepage)
	fmt.Println("On your network: http://" + getip() + ":" + config.Port + config.Homepage)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
func getip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
func (fs SinglePageFS) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)
	if err != nil {
		file, err = fs.FileSystem.Open("index.html")
	}
	return file, err
}
