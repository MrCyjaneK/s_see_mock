package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

const defaultMaxInputSize = 16 * 1024 * 1024 * 1024 // 16 GiB

const moduleVersion = "2.0.0"
const moduleExtension = "json"

var (
	loadErr string
	tempDir string
)

//export SPACESEE_mock_load
func SPACESEE_mock_load(workDir, assetsDir *C.char) C.int {
	loadErr = ""
	path := C.GoString(workDir)
	if path == "" {
		loadErr = "work dir is empty"
		return 1
	}
	if err := os.MkdirAll(path, 0o700); err != nil {
		loadErr = err.Error()
		return 1
	}
	tempDir = path
	_ = assetsDir
	return 0
}

//export SPACESEE_mock_unload
func SPACESEE_mock_unload() {}

//export SPACESEE_mock_isHeavy
func SPACESEE_mock_isHeavy() C.int { return 0 }

//export SPACESEE_mock_assetCount
func SPACESEE_mock_assetCount() C.int { return 0 }

//export SPACESEE_mock_asset
func SPACESEE_mock_asset(index C.int) *C.char {
	_ = index
	return nil
}

//export SPACESEE_mock_loadError
func SPACESEE_mock_loadError() *C.char {
	if loadErr == "" {
		return nil
	}
	return C.CString(loadErr)
}

//export SPACESEE_mock_version
func SPACESEE_mock_version() *C.char {
	return C.CString(moduleVersion)
}

//export SPACESEE_mock_extension
func SPACESEE_mock_extension() *C.char {
	return C.CString(moduleExtension)
}

//export SPACESEE_mock_filetypes
func SPACESEE_mock_filetypes() *C.char {
	return C.CString("*")
}

//export SPACESEE_mock_maxInputSize
func SPACESEE_mock_maxInputSize() C.longlong {
	return C.longlong(defaultMaxInputSize)
}

//export SPACESEE_mock_analyze
func SPACESEE_mock_analyze(inputPath *C.char, outputPath *C.char, dataDir *C.char) C.longlong {
	in := C.GoString(inputPath)
	out := C.GoString(outputPath)
	data := C.GoString(dataDir)
	jobID := newJob()
	go runAnalyze(jobID, in, out, data)
	return C.longlong(jobID)
}

func runAnalyze(jobID int64, inputPath, outputPath, dataDir string) {
	if dataDir == "" {
		setJobError(jobID, os.ErrInvalid)
		return
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		setJobError(jobID, err)
		return
	}
	touched := filepath.Join(dataDir, "touched")
	fTouched, err := os.OpenFile(touched, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err == nil {
		_, werr := fTouched.Write([]byte("mock\n"))
		cerr := fTouched.Close()
		if werr != nil {
			setJobError(jobID, werr)
			return
		}
		if cerr != nil {
			setJobError(jobID, cerr)
			return
		}
	} else if !os.IsExist(err) {
		setJobError(jobID, err)
		return
	}

	setJobProgress(jobID, 10, "reading input")
	f, err := os.Open(inputPath)
	if err != nil {
		setJobError(jobID, err)
		return
	}
	defer f.Close()

	setJobProgress(jobID, 40, "computing hashes")
	sha1H := sha1.New()
	sha224H := sha256.New224()
	sha256H := sha256.New()
	sha512H := sha512.New()
	md5H := md5.New()
	if _, err := io.Copy(io.MultiWriter(sha1H, sha224H, sha256H, sha512H, md5H), f); err != nil {
		setJobError(jobID, err)
		return
	}

	setJobProgress(jobID, 70, "writing output")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o700); err != nil {
		setJobError(jobID, err)
		return
	}
	payload := map[string]string{
		"sha1":   hex.EncodeToString(sha1H.Sum(nil)),
		"sha224": hex.EncodeToString(sha224H.Sum(nil)),
		"sha256": hex.EncodeToString(sha256H.Sum(nil)),
		"sha512": hex.EncodeToString(sha512H.Sum(nil)),
		"md5":    hex.EncodeToString(md5H.Sum(nil)),
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		setJobError(jobID, err)
		return
	}
	if err := os.WriteFile(outputPath, b, 0o600); err != nil {
		setJobError(jobID, err)
		return
	}
	setJobDone(jobID, "done")
}

//export SPACESEE_mock_progress
func SPACESEE_mock_progress(jobID C.longlong, scale C.int) C.int {
	return C.int(jobProgress(int64(jobID), int(scale)))
}

//export SPACESEE_mock_progressString
func SPACESEE_mock_progressString(jobID C.longlong) *C.char {
	return C.CString(jobProgressString(int64(jobID)))
}

//export SPACESEE_mock_free
func SPACESEE_mock_free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func main() {}
