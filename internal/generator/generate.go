package generator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const fileMode = 0o600

type StaticSite struct {
	OutDir    string
	StaticDir string
}

// GenerateStaticSite generates a static site by iterating over the routes in a Gin router
// and creating corresponding HTML files in the specified output directory.
func GenerateStaticSite(r *gin.Engine, props StaticSite) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(props.OutDir, fileMode); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Iterate over the routes in the Gin router
	for _, route := range r.Routes() {
		if route.Method != http.MethodGet {
			continue
		}

		// Handle static file routes
		if strings.HasPrefix(route.Path, "/static") {
			staticFilePath := strings.TrimPrefix(route.Path, "/static")
			srcPath := filepath.Join(props.StaticDir, staticFilePath)
			destPath := filepath.Join(props.OutDir, staticFilePath)

			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy static file %s: %w", srcPath, err)
			}
			continue
		}

		b, err := GeneratePage(r, route)
		if err != nil {
			return fmt.Errorf("failed to generate page for route %s: %w", route.Path, err)
		}

		// Determine the output file path
		outputPath := filepath.Join(props.OutDir, "index.html")
		if route.Path != "/" {
			outputPath = filepath.Join(props.OutDir, route.Path, "index.html")
		}

		// Ensure the directory for the output file exists
		if err := os.MkdirAll(filepath.Dir(outputPath), fileMode); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
		}

		// Write the response to the output file
		if err := os.WriteFile(outputPath, b, fileMode); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Create the directory structure for the destination file
	if err := os.MkdirAll(filepath.Dir(dst), fileMode); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", dst, err)
	}

	// Create the destination file for writing
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents from %s to %s: %w", src, dst, err)
	}

	return nil
}

// GeneratePage simulates a request to the given route
// Returns the page body and an error
func GeneratePage(r *gin.Engine, route gin.RouteInfo) ([]byte, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route.Path, http.NoBody)
	if err != nil {
		return nil, err
	}
	r.ServeHTTP(w, req)
	return w.Body.Bytes(), nil
}
