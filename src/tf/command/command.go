package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Request URL constant
const TRANSFER string = "https://transfer.sh"

// Storage file path
var STORAGE_FILE string = os.Getenv("HOME") + "/.tfstorage.json"

/**
Donwload File
@param []string upload file list
@param bool     secure flag
@return void
*/
func Download(args []string, secure bool) {
	var (
		file     string
		fileName string
		buffer   bytes.Buffer
	)

	file = strings.Trim(args[0], "\n")
	fileName = path.Base(file)
	fmt.Printf("Getting file: %s\n", file)

	cmd := exec.Command("curl", file)
	cmd.Stdout = &buffer
	fmt.Printf("Downloading %s...", fileName)

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v", err)
		return
	}
	cwd, _ := os.Getwd()
	if writeErr := ioutil.WriteFile(cwd+"/"+fileName, buffer.Bytes(), 0644); writeErr != nil {
		log.Fatalf("%v", writeErr)
		return
	}
	fmt.Println("Success!")
}

/**
Upload File
@param []string upload file list
@param string   file tag
@param bool     secure flag
@return void
*/
func Upload(args []string, tag string, secure bool) {
	if _, statErr := os.Stat(args[0]); statErr != nil {
		log.Fatalf("%v", statErr)
		return
	}

	var buffer bytes.Buffer
	file := path.Base(args[0])
	cmd := exec.Command("curl", "--upload-file", args[0], TRANSFER+"/"+file)
	cmd.Stdout = &buffer
	fmt.Printf("Uploading %s...\n", file)

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v", err)
		return
	}

	fmt.Printf("Upload Success: file URL is %s\n", buffer.String())

	// need to save tag?
	if len(tag) > 0 {
		fmt.Printf("Tag saved: %s\n", tag)
		makeTag(tag, buffer.String())
	}
}

/**
Find File from marked tag
@param string   file tag
@param bool     secure flag
@return void
*/
func Find(tag string) {
	var dat = make(map[string][]string)
	var (
		fp       *os.File
		fileList []string
		input    string
	)

	if stat, statErr := os.Stat(STORAGE_FILE); statErr != nil {
		fp, _ = os.Create(STORAGE_FILE)
		defer fp.Close()
	} else {
		fp, _ := os.Open(STORAGE_FILE)
		buffer := make([]byte, stat.Size())
		defer fp.Close()
		if _, readErr := fp.Read(buffer); readErr != nil {
			log.Fatalf("FileRead Error: %v", readErr)
		}
		json.Unmarshal(buffer, &dat)
	}

	if len(tag) == 0 {
		for k, v := range dat {
			fmt.Println("Tag: " + k)
			for i, f := range v {
				fmt.Printf("[%d] %s", i, f)
				fileList = append(fileList, f)
			}
		}
	} else if v, ok := dat[tag]; ok {
		fmt.Println("Tag: " + tag)
		for i, f := range v {
			fmt.Printf("[%d] %s", i, f)
			fileList = append(fileList, f)
		}
	} else {
		fmt.Println("Tag " + tag + " not found")
	}

	// Getting file input
	fmt.Printf("\nGetting file? [0-%d]:", len(fileList)-1)
	fmt.Scanln(&input)

	num, _ := strconv.Atoi(input)
	if num >= 0 && num < len(fileList) {
		Download([]string{strings.Trim(fileList[num], "\n")}, false)
	} else {
		fmt.Println("Index not found.")
	}
}
func Help() {
	fmt.Println("Transfer.sh (https://transfer.sh/) command tool")
	fmt.Println("===============================================")
	fmt.Println("Usage:")
	fmt.Println("    tf [-sth] operation file")
	fmt.Println("Options:")
	fmt.Println("   -s : Secure upload/download")
	fmt.Println("   -t : Tag mark")
	fmt.Println("   -h : Show this help")
	fmt.Println("Operation::")
	fmt.Println("   upload   : Upload file")
	fmt.Println("   download : Download file")
	fmt.Println("   find     : Find file from Tag")
	fmt.Println("   help     : Show this help")
}

func makeTag(tag, file string) {
	var dat = make(map[string][]string)
	var fp *os.File

	STORAGE_FILE := os.Getenv("HOME") + "/.tfstorage.json"
	_, statErr := os.Stat(STORAGE_FILE)
	if statErr != nil {
		fp, _ = os.Create(STORAGE_FILE)
		defer fp.Close()
	} else {
		var buffer []byte
		fp, _ := os.Open(STORAGE_FILE)
		defer fp.Close()
		fp.Read(buffer)
		json.Unmarshal(buffer, &dat)
	}

	if v, ok := dat[tag]; ok {
		v = append(v, file)
		dat[tag] = v
	} else {
		dat[tag] = []string{file}
	}

	writeData, _ := json.Marshal(dat)
	fp.Write(writeData)
}
