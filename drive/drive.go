package drive

import (
	"fmt"
	"strings"
	"os"
	"io"
	"encoding/json"
	"io/ioutil"
	"net/http"
	md "github.com/lunny/html2md"
)

func DocToHTML(input, outDir string) {
	files, _ := ioutil.ReadDir(input)

	for _, file := range files {
		if strings.Index(file.Name(), ".gdoc") < 0 {
			continue
		}
		content, _ := ioutil.ReadFile(input + "/" + file.Name())
		var doc map[string]string
		json.Unmarshal(content, &doc)
		response, _ := http.Get("https://docs.google.com/feeds/download/documents/export/Export?exportFormat=html&id=" + doc["doc_id"])
		fileName := strings.ReplaceAll(kebab(file.Name()), "gdoc", "html")
		newFile, _ := os.Create(outDir + "/" + fileName)
		io.Copy(newFile, response.Body)
		os.Chmod(outDir+"/"+fileName, 0644)
	}
}

func HTMLtoMD(input, outDir string) {
	files, _ := ioutil.ReadDir(input)
	for _, file := range files {
		if strings.Index(file.Name(), ".html") < 0 {
			continue
		}
		raw, _ := ioutil.ReadFile(input+"/"+file.Name())
		err := ioutil.WriteFile(
			outDir+"/"+strings.ReplaceAll(kebab(file.Name()), "html", "md"), 
			[]byte(md.Convert(string(raw))), 
			0644) 
		if err != nil {
			fmt.Println("ERROR", err)
		}
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
