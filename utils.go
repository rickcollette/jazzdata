package jazzdata

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

var ImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var DefaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (Linux; Android 7.0; BLN-L22) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
}

func DownloadCover(url string, targetDir string, coverName string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range DefaultHeaders {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to download cover: " + resp.Status)
	}

	outPath := filepath.Join(targetDir, coverName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	return err
}

func CompareCovers(path1, path2 string) (int, error) {
	// Open and decode the images.
	file1, err := os.Open(path1)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	img1, _, err := image.Decode(file1)
	if err != nil {
		return 0, err
	}
	file2, err := os.Open(path2)
	if err != nil {
		return 0, err
	}
	defer file2.Close()
	img2, _, err := image.Decode(file2)
	if err != nil {
		return 0, err
	}

	hash1, err := goimagehash.AverageHash(img1)
	if err != nil {
		return 0, err
	}
	hash2, err := goimagehash.AverageHash(img2)
	if err != nil {
		return 0, err
	}

	return hash1.Distance(hash2)
}

func GetAlbumPaths(root string, mustHaveCover bool) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() || path == root {
			return nil
		}
		if IsAlbumDir(path) || (!mustHaveCover && HasSong(path)) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func GetBasePath(p string) string {
	info, err := os.Stat(p)
	if err == nil && !info.IsDir() {
		return filepath.Dir(p)
	}
	return filepath.Clean(p)
}

func GetExtensionFromURL(urlStr string) string {
	// Optionally send a HEAD request; here we use the MIME type if available.
	resp, err := http.Head(urlStr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ctype := resp.Header.Get("Content-Type")
	exts, err := mime.ExtensionsByType(ctype)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ""
}

func GetCover(p string) string {
	info, err := os.Stat(p)
	if err != nil {
		return ""
	}
	dir := p
	if !info.IsDir() {
		dir = filepath.Dir(p)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.ToLower(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ImageExtensions[ext] && (name == "folder" || name == "poster" || name == "cover" || name == "default") {
			return filepath.Join(dir, file.Name())
		}
	}

	return ""
}

func HasSong(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if SupportedSongExtensions[ext] {
			return true
		}
	}
	return false
}

func HasCover(p string) bool {
	return GetCover(p) != ""
}

func IsAlbumDir(dir string) bool {
	return HasCover(dir) && HasSong(dir)
}
