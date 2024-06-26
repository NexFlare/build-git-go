package object

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/git-starter-go/internal/util"
)

func GetFileData(filesha string) []byte {
	fileContent, err := os.ReadFile(fmt.Sprintf(".git/objects/%s/%s", filesha[:2], filesha[2:]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	bytesReader := bytes.NewReader(fileContent)
	bytesWriter := &bytes.Buffer{}
	reader, err := zlib.NewReader(bytesReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong: %s\n", err.Error())
		os.Exit(1)
	}
	defer reader.Close()
	_, err = io.Copy(bytesWriter, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	return bytesWriter.Bytes()
}

func GetHeader(data []byte) ObjectHeader {
	header := ObjectHeader{}
	_, err := fmt.Sscanf(string(data), "%s %d", &header.Type, &header.Size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading header: %s\n", err.Error())
		os.Exit(1)
	}
	return header
}

func CreateBlog(fileContent []byte, writeToObject bool) []byte {
	blobContent := fmt.Sprintf("blob %d\x00%s", len(string(fileContent)), string(fileContent))
	hashedValue := util.GetHash([]byte(blobContent))
	if writeToObject {
		util.CreateObjectFile([]byte(blobContent), hex.EncodeToString(hashedValue))
	}
	return hashedValue
}
