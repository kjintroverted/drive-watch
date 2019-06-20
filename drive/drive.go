package drive

import (
	"strings"
	"os"
	"io"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func DocToHTML(input, outDir string) string {
	files, _ := ioutil.ReadDir(input)

	for _, file := range files {
		if strings.Index(file.Name(), "gdoc") < 0 {
			continue
		}
		content, _ := ioutil.ReadFile(input + "/" + file.Name())
		var doc map[string]string
		json.Unmarshal(content, &doc)
		response, _ := http.Get("https://docs.google.com/feeds/download/documents/export/Export?exportFormat=html&id=" + doc["doc_id"])
		newFile, _ := os.Create(outDir + "/" + strings.Replace(kebab(file.Name()), "gdoc", "html", -1))
		io.Copy(newFile, response.Body)
	}

	return outDir
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

// func HTMLtoMD(input, outDir string) {

// }
