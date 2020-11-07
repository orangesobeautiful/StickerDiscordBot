package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"ourbot/DatabaseOperation-Go/imageproxy/googledriver"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadFile(folderPath string, url string) (string, int, error) {

	// Get the data
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var statusCode = resp.StatusCode
	if statusCode != 200 {
		return "", statusCode, errors.New(strconv.Itoa(resp.StatusCode))
	}

	_, params, _ := mime.ParseMediaType(resp.Header["Content-Disposition"][0])
	fileName := params["filename"]

	// Create the file
	out, err := os.Create(path.Join(folderPath, fileName))
	if err != nil {
		return fileName, statusCode, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return fileName, statusCode, err
}

var imgFolderPath string
var dbURL string = os.Getenv("DATABASE_URL")
var dbOp *googledriver.DataBaseOperation

func googleDriverImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	//fmt.Println(id)
	var err error
	var statusCode int
	localFileName := dbOp.GetLocalFileName(id)
	var filePath string
	filePath = path.Join(imgFolderPath, localFileName)
	if localFileName == "" || !fileExists(filePath) {
		dbOp.RemoveImgByID(id)
		localFileName, statusCode, err = downloadFile(imgFolderPath, "https://drive.google.com/uc?export=view&id="+id)
		if statusCode != 200 {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, "QQ "+strconv.Itoa(statusCode))
		}else{
			if err != nil {
				panic(err)
			}
		}
		
		filePath = path.Join(imgFolderPath, localFileName)
		dbOp.AddNewImg(id, localFileName)
	} else {
		dbOp.UpdateUseTime(id)
	}

	w.Header().Set("Content-Disposition", "inline;filename=\""+localFileName+"\";filename*=UTF-8''"+localFileName)
	http.ServeFile(w, r, filePath)
}

func main() {
	var err error
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	defaultImgFolder := path.Join(exeDir, "ImageProxyFolder")

	var cliHost = flag.String("host", "127.0.0.1", "Listen Host")
	var cliPort = flag.String("port", "5001", "Listen Port")
	var cliLogFile = flag.String("log", "", "Log file path")
	var cliProxyHeader = flag.Bool("proxy", false, "Has Proxy Header?")
	var cliImgFolderPath = flag.String("imagefolder", defaultImgFolder, "Image Folder Path")

	flag.Parse()

	host := *cliHost
	port := *cliPort
	output := *cliLogFile
	isProxyHeader := *cliProxyHeader
	imgFolderPath = *cliImgFolderPath

	err = os.MkdirAll(imgFolderPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	if dbURL == "" {
		dbURL = "postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable"
	}

	dbOp = googledriver.DbOpen(dbURL)
	//fmt.Println(dbOp.GetLocalFileName("123"))
	defer dbOp.DbClose()

	r := httprouter.New()
	r.GET("/img-proxy/google-driver/:id", googleDriverImage)

	var f *os.File
	var serverHandler http.Handler

	if output != "" {
		f, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		serverHandler = handlers.LoggingHandler(f, r)
		serverHandler = handlers.LoggingHandler(os.Stdout, serverHandler)

		defer f.Close()
	} else {
		serverHandler = handlers.LoggingHandler(os.Stdout, r)
	}

	if isProxyHeader {
		serverHandler = handlers.ProxyHeaders(serverHandler)
	}

	fmt.Println("Image Proxy Start Listen!")
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), serverHandler)
}
