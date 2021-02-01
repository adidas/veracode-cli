package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}

// newFormRequest Creates a new form http request with optional extra params
func newFormRequest(uri string, params map[string]string) (*http.Request, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request, nil
}

func makeRequest(uri string, user string, pass string, formParams map[string]string, requestName string) (string, error) {
	var request *http.Request
	var method string
	var err error
	if uri == "" {
		return "", errors.New("URI cannot be empty")
	}
	if formParams == nil {
		method = "GET"
		request, err = http.NewRequest("GET", uri, nil)
		if err != nil {
			return "", err
		}
	} else {
		method = "POST"
		request, err = newFormRequest(uri, formParams)
		if err != nil {
			return "", err
		}
	}

	authHeader, err := generateHeader(uri, method, user, pass)

	if err != nil {
		return "Error on generateHeader", err
	}

	request.Header.Set("Authorization", authHeader)
	log.Println(requestName)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		// handle err
		return "", errors.New(API_CON_ERROR)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	//log.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return "", errors.New(API_NOT_FOUND)
		}
		if resp.StatusCode == 403 {
			return "", errors.New(FORBIDEN_ACCESS)
		}
		if resp.StatusCode == 401 {
			return "", errors.New(INVALID_CREDS)
		}
		s := strconv.Itoa(resp.StatusCode)+"\n"+bodyString[len(bodyString)-1300 : len(bodyString)-900]
		return "", errors.New(API_UNEXPECTED_ERROR + "\n\t" + s)
	}
	
	if strings.Contains(bodyString,"Access denied"){
		return "", errors.New(FORBIDEN_ACCESS)
	}
	
	if strings.Contains(bodyString,"<error>"){
		//s := strconv.Itoa(resp.StatusCode)+"\n"+bodyString[len(bodyString)-1300 : len(bodyString)-900]
		s := strings.Split(bodyString, "<error>")[1]
		return "", errors.New(API_UNEXPECTED_ERROR + "\n" +"Error message: "+ s)
	}
	//	log.Println("bodyString:", bodyString)

	return bodyString, nil

}

// DownloadFile download a file form URL with basic authentication
func DownloadFile(filepath string, uri string, user string, pass string) error {
	var request *http.Request
	var err error
	if uri == "" {
		return errors.New("URI cannot be empty")
	}

	request, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	authHeader, err := generateHeader(uri, "GET", user, pass)

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", authHeader)

	log.Println("Starting request")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		// handle err
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
