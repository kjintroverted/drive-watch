package drive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	md "github.com/kjintroverted/html-to-markdown"
)

type fileInfo struct {
	Name    string
	ModTime time.Time
}

func AllDocToHTML(input, outDir string) {
	files, _ := ioutil.ReadDir(input)
	for _, file := range files {
		DocToHTML(input, file.Name(), outDir)
	}
}

func DocToHTML(path, name, outDir string) {
	if strings.Index(name, ".gdoc") < 0 {
		return
	}

	fileName := strings.ReplaceAll(Kebab(name), ".gdoc", "")

	// SAVE METADATA AS JSON
	info, _ := os.Stat(path + "/" + name)
	jsonbytes, _ := json.Marshal(fileInfo{Name: strings.ReplaceAll(info.Name(), ".gdoc", ""), ModTime: info.ModTime()})
	ioutil.WriteFile(outDir+"/"+fileName+".json", jsonbytes, 0644)

	// SAVE CONTENT AS HTML
	content, _ := ioutil.ReadFile(path + "/" + name)
	var doc map[string]string
	json.Unmarshal(content, &doc)
	response, _ := http.Get("https://docs.google.com/feeds/download/documents/export/Export?exportFormat=html&id=" + doc["doc_id"])
	newFile, _ := os.Create(outDir + "/" + fileName + ".html")
	io.Copy(newFile, response.Body)
	os.Chmod(outDir+"/"+fileName+".html", 0644)
}

func AllHTMLtoMD(input, outDir string) {
	files, _ := ioutil.ReadDir(input)
	for _, file := range files {
		HTMLtoMD(input, file.Name(), outDir)
	}
}

func HTMLtoMD(path, name, outDir string) {
	if strings.Index(name, ".html") < 0 {
		return
	}

	fileName := strings.ReplaceAll(Kebab(name), ".html", "")

	// CONVERT HTML
	converter := md.NewConverter("", true, nil)
	raw, _ := ioutil.ReadFile(path + "/" + name)
	mdContent, _ := converter.ConvertString(string(raw))

	mdContent = unescape(mdContent, "`", "\\")

	// ADD FRONT MATTER
	fileDataRaw, _ := ioutil.ReadFile(path + "/" + fileName + ".json")
	var fileData fileInfo
	json.Unmarshal(fileDataRaw, &fileData)
	mdContent = fmt.Sprintf(`---
title: %s
lastUpdated: %s
---
%s`, fileData.Name, fileData.ModTime, mdContent)

	err := ioutil.WriteFile(
		outDir+"/"+fileName+".md",
		[]byte(mdContent),
		0644)
	if err != nil {
		fmt.Println("ERROR", err)
	}
}

func Kebab(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	return strings.ReplaceAll(s, " ", "-")
}

func unescape(s string, characters ...string) string {
	for _, c := range characters {
		s = strings.ReplaceAll(s, "\\"+c, c)
	}
	return s
}
