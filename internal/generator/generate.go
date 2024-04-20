package generator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const fileMode fs.FileMode = 0o777 // TODO: why is it only working with 777 and not even 666?

// GenerateStaticSite generates a static site by iterating over the routes in a Gin router
// and creating corresponding HTML files in the specified output directory.
func GenerateStaticSite(r *gin.Engine, outDir string, fsys fs.FS) error {
	log := slog.Default()
	// Ensure the output directory exists
	if err := os.MkdirAll(outDir, fileMode); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Iterate over the routes in the Gin router
	for _, route := range r.Routes() {
		if route.Method != http.MethodGet {
			continue
		}

		// Handle static file routes
		if strings.HasPrefix(route.Path, "/static") {
			err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.HasPrefix(path, "static") {
					srcPath := path
					destPath := filepath.Join(outDir, srcPath)
					if err := copyFile(fsys, srcPath, destPath); err != nil {
						return fmt.Errorf("failed to copy static file %s: %w", srcPath, err)
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to copy static files: %w", err)
			}
			continue
		}

		b, err := generatePage(r, route)
		if err != nil {
			return fmt.Errorf("failed to generate page for route %s: %w", route.Path, err)
		}

		// Determine the output file path
		outputPath := filepath.Join(outDir, "index.html")
		if route.Path != "/" {
			outputPath = filepath.Join(outDir, route.Path, "index.html")
		}
		log.Info("Generating page", "route", route.Path, "outputDir", outputPath)

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

// copyFile copies a file from the source to the destination
func copyFile(fsys fs.FS, src, dst string) error {
	srcFile, err := fsys.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	if err = os.MkdirAll(filepath.Dir(dst), fileMode); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", dst, err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents from %s to %s: %w", src, dst, err)
	}

	return nil
}

// generatePage simulates a request to the given route
// Returns the page body and an error if one occurred
func generatePage(r *gin.Engine, route gin.RouteInfo) (b []byte, err error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route.Path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer func(b io.Closer) {
		if cErr := b.Close(); cErr != nil {
			err = errors.Join(err, cErr)
		}
	}(req.Body)

	r.ServeHTTP(w, req)
	return w.Body.Bytes(), nil
}
