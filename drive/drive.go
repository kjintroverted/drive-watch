package drive

import (
	"fmt"
	"strings"
	"os"
	"io"
	"encoding/json"
	"io/ioutil"
	"net/http"
	md "github.com/kjintroverted/html-to-markdown"
)

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
	content, _ := ioutil.ReadFile(path + "/" + name)
	var doc map[string]string
	json.Unmarshal(content, &doc)
	response, _ := http.Get("https://docs.google.com/feeds/download/documents/export/Export?exportFormat=html&id=" + doc["doc_id"])
	fileName := strings.ReplaceAll(kebab(name), "gdoc", "html")
	newFile, _ := os.Create(outDir + "/" + fileName)
	io.Copy(newFile, response.Body)
	os.Chmod(outDir+"/"+fileName, 0644)
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

	raw, _ := ioutil.ReadFile(path+"/"+name)
	
	converter := md.NewConverter("", true, nil)
	mdContent, _ := converter.ConvertString(string(raw))
	err := ioutil.WriteFile(
		outDir+"/"+strings.ReplaceAll(kebab(name), "html", "md"), 
		[]byte(mdContent), 
		0644) 
	if err != nil {
		fmt.Println("ERROR", err)
	}
}

func kebab(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	return strings.ReplaceAll(s, " ", "-")
}
