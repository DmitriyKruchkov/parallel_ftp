package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)


func main() {
	port := flag.String("port", "8000", "port to listen")
	ftp_root_dir := flag.String("dir", ".", "Root directory for FTP server")

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // Например, обрыв соединения
			continue
		}
		go handleConn(conn, ftp_root_dir) // Обработка подключения
	}
}
func handleConn(c net.Conn, root_dir *string) {
	defer c.Close()
	current_dir := *root_dir
	for {
		reader := bufio.NewReader(c)
		message, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	cmd_and_params := strings.Split(message, " ")
	switch strings.Replace(cmd_and_params[0], "\n", "", -1) {
		case "cd":
			current_dir = get_way(cmd_and_params[1], *root_dir)
			os.MkdirAll(current_dir, 0666)
		case "ls":
			files := func(entries []os.DirEntry, err error) string {
				if err != nil {
					log.Println("Dir read error", err)
				}
				list := []string{}
				for _, entry := range entries {
					list = append(list, entry.Name())
				}
				return strings.Join(list, "\n")
			}(os.ReadDir(current_dir))
			io.WriteString(c, files + "\n")
		case "get":
			filePath := filepath.Join(current_dir, strings.TrimSpace(cmd_and_params[1]))
    
			data, err := os.ReadFile(filePath)
			if err != nil {
				io.WriteString(c, "File not found\n")
			}
			
			io.WriteString(c, string(data) + "\n")
			
		
		case "close":
			return
	}
	}
}

func get_way(path string, root_path string) string {
	var abs_path string
	if path[0] == '/' {
		path = path[1:]
	}
	abs_path = root_path + "/" + path
	return abs_path[:len(abs_path)-1]

}

