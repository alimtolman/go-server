package main

import (
	"Server/DataBase"
	"Server/DataBase/Tables"
	"Server/Helpers/Convert"
	"Server/Helpers/File"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const http_port = "80"
const https_port = "443"
const cert_path = "_cert.pem"
const key_path = "_key.pem"
const root_msg = "It Works!"

/* AppData */

func AppDataAdd(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	q_id := request.URL.Query().Get("id")
	id := Convert.StrToUInt(q_id, 0)
	item := DataBase.NewAppData(id, string(body))

	Tables.AppData.AddOrSet(item)
	fmt.Fprint(writer, "OK")
}

func AppDataClear(writer http.ResponseWriter, request *http.Request) {
	Tables.AppData.RemoveAll()
	fmt.Fprint(writer, "OK")
}

func AppDataGet(writer http.ResponseWriter, request *http.Request) {
	q_id := request.URL.Query().Get("id")
	id := Convert.StrToUInt(q_id, 0)
	data, err := Tables.AppData.Get(func(item DataBase.AppData) bool { return item.Id == id })

	if err == nil {
		fmt.Fprint(writer, data.Data)
	} else {
		fmt.Fprint(writer, err.Error())
	}
}

func AppDataIds(writer http.ResponseWriter, request *http.Request) {
	dataCount := Tables.AppData.Count()
	ids := make([]uint32, 0)

	for i := uint32(0); i < dataCount; i++ {
		item, _ := Tables.AppData.GetAt(i)
		ids = append(ids, item.Id)
	}

	fmt.Fprint(writer, fmt.Sprint(ids))
}

/* Root */

func RootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		fmt.Fprint(writer, root_msg)
	} else {
		http.NotFound(writer, request)
	}
}

func main() {
	var waitGroup sync.WaitGroup

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/api/data/add", AppDataAdd)
	http.HandleFunc("/api/data/clear", AppDataClear)
	http.HandleFunc("/api/data/get", AppDataGet)
	http.HandleFunc("/api/data/ids", AppDataIds)

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		if err := http.ListenAndServe(":"+http_port, nil); err != nil {
			log.Println("Server listen error:", err)
		}
	}()

	if File.Exists(cert_path) && File.Exists(key_path) {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			if err := http.ListenAndServeTLS(":"+https_port, cert_path, key_path, nil); err != nil {
				log.Println("Server tls listen error:", err)
			}
		}()
	} else {
		log.Println("TLS cert/key not found, skipping HTTPS listener:", cert_path, key_path)
	}

	waitGroup.Wait()
}
