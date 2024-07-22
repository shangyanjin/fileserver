package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileInfo holds information about a file
type FileInfo struct {
	Name  string
	Path  string
	IsDir bool
}

const (
	defaultPort = 3389         // 默认端口
	defaultDir  = "./download" // 默认文件目录
)

func main() {
	port := findAvailablePort(defaultPort)
	fmt.Printf("Serving on http://localhost:%d\n", port)

	r := gin.Default()
	r.GET("/*path", fileHandler)

	r.Run(fmt.Sprintf(":%d", port))
}

func findAvailablePort(startPort int) int {
	port := startPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port
		}
		port++
	}
}

func fileHandler(c *gin.Context) {
	dir := defaultDir // 使用默认文件目录
	requestedPath := c.Param("path")[1:]
	fullPath := filepath.Join(dir, requestedPath)
	fileInfo, err := os.Stat(fullPath)

	if err != nil {
		c.String(http.StatusNotFound, "File or directory not found")
		return
	}

	if fileInfo.IsDir() {
		serveDirectory(c, fullPath, requestedPath)
	} else {
		serveFile(c, fullPath)
	}
}

func serveDirectory(c *gin.Context, dirPath, requestedPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to read directory")
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:  file.Name(),
			Path:  filepath.Join(requestedPath, file.Name()),
			IsDir: file.IsDir(),
		})
	}

	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File Server</title>
        <style>
            body { font-family: Arial, sans-serif; }
            table { width: 100%; border-collapse: collapse; }
            th, td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
            tr:hover { background-color: #f5f5f5; }
            a { text-decoration: none; color: #00f; }
            a:hover { text-decoration: underline; }
        </style>
    </head>
    <body>
        <h1>File Server</h1>
        <table>
            <tr>
                <th>Name</th>
                <th>Type</th>
            </tr>
            {{range .}}
            <tr>
                <td><a href="/{{.Path}}">{{.Name}}</a></td>
                <td>{{if .IsDir}}Directory{{else}}File{{end}}</td>
            </tr>
            {{end}}
        </table>
    </body>
    </html>
    `
	t, err := template.New("filelist").Parse(tmpl)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to load template")
		return
	}

	// Replace backslashes with forward slashes for cross-platform compatibility
	for i := range fileInfos {
		fileInfos[i].Path = strings.Replace(fileInfos[i].Path, "\\", "/", -1)
	}

	t.Execute(c.Writer, fileInfos)
}

func serveFile(c *gin.Context, filePath string) {
	http.ServeFile(c.Writer, c.Request, filePath)
}
