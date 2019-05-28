package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/swiftstack/ProxyFS/conf"
	"github.com/swiftstack/ProxyFS/jrpcfs"
)

type jrpcRequestStruct struct {
	JSONrpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	ID      uint64         `json:"id"`
	Params  [1]interface{} `json:"params"`
}

type jrpcRequestEmptyParamStruct struct{}

type jrpcResponseIDAndErrorStruct struct {
	ID    uint64 `json:"id"`
	Error string `json:"error"`
}

type jrpcResponseNoErrorStruct struct {
	ID     uint64      `json:"id"`
	Result interface{} `json:"result"`
}

type envStruct struct {
	AuthToken     string
	StorageURL    string
	LastMessageID uint64
	MountID       jrpcfs.MountIDAsString
}

var (
	authKey  string
	authURL  string
	authUser string
	envPath  string
	verbose  bool

	env *envStruct
)

func main() {
	var (
		args         []string
		cmd          string
		confFilePath string
		confMap      conf.ConfMap
		err          error
	)

	// Parse arguments

	args = os.Args[1:]

	if 2 > len(args) {
		log.Fatalf("Must specify a .conf and a command")
	}

	confFilePath = args[0]

	confMap, err = conf.MakeConfMapFromFile(confFilePath)
	if nil != err {
		log.Fatalf("Failed to load %s: %v", confFilePath, err)
	}

	authURL, err = confMap.FetchOptionValueString("JRPCTool", "AuthURL")
	if nil != err {
		log.Fatalf("Failed to parse %s value for JRPCTool.AuthURL: %v", confFilePath, err)
	}
	authUser, err = confMap.FetchOptionValueString("JRPCTool", "AuthUser")
	if nil != err {
		log.Fatalf("Failed to parse %s value for JRPCTool.AuthUser: %v", confFilePath, err)
	}
	authKey, err = confMap.FetchOptionValueString("JRPCTool", "AuthKey")
	if nil != err {
		log.Fatalf("Failed to parse %s value for JRPCTool.AuthKey: %v", confFilePath, err)
	}
	envPath, err = confMap.FetchOptionValueString("JRPCTool", "EnvPath")
	if nil != err {
		log.Fatalf("Failed to parse %s value for JRPCTool.EnvPath: %v", confFilePath, err)
	}
	verbose, err = confMap.FetchOptionValueBool("JRPCTool", "Verbose")
	if nil != err {
		log.Fatalf("Failed to parse %s value for JRPCTool.Verbose: %v", confFilePath, err)
	}

	cmd = args[1]

	switch cmd {
	case "a":
		doAuth()
	case "m":
		doMount()
	case "r":
		doJRPC(args[2:]...)
	case "c":
		doClean()
	default:
		log.Fatalf("Could not understand command \"%s\"", cmd)
	}
}

func doAuth() {
	var (
		err          error
		httpClient   *http.Client
		httpRequest  *http.Request
		httpResponse *http.Response
	)

	httpRequest, err = http.NewRequest("GET", authURL, nil)
	if nil != err {
		log.Fatalf("Failed to create GET of authURL==\"%s\": %v", authURL, err)
	}

	httpRequest.Header["X-Auth-User"] = []string{authUser}
	httpRequest.Header["X-Auth-Key"] = []string{authKey}

	if verbose {
		fmt.Printf("Auth Request:\n")
		fmt.Println("  httpRequest.URL:   ", httpRequest.URL)
		fmt.Println("  httpRequest.Header:", httpRequest.Header)
	}

	httpClient = &http.Client{}

	httpResponse, err = httpClient.Do(httpRequest)
	if nil != err {
		log.Fatalf("Failed to issue GET of authURL==\"%s\": %v", authURL, err)
	}
	if http.StatusOK != httpResponse.StatusCode {
		log.Fatalf("Received unexpected HTTP Status for Get of authURL==\"%s\": %s", authURL, httpResponse.Status)
	}

	env = &envStruct{
		AuthToken:     httpResponse.Header.Get("X-Auth-Token"),
		StorageURL:    httpResponse.Header.Get("X-Storage-Url"),
		LastMessageID: 0,
		MountID:       "",
	}

	if verbose {
		fmt.Printf("Auth Response:\n")
		fmt.Printf("  env.AuthToken:  %s\n", env.AuthToken)
		fmt.Printf("  env.StorageURL: %s\n", env.StorageURL)
	}

	writeEnv()
}

func doMount() {
	var (
		accountName                 string
		err                         error
		httpClient                  *http.Client
		httpRequest                 *http.Request
		httpResponse                *http.Response
		httpResponseBody            []byte
		httpResponseBodyBytesBuffer bytes.Buffer
		jrpcRequest                 *jrpcRequestStruct
		jrpcRequestBuf              []byte
		jrpcRequestBytesBuffer      bytes.Buffer
		jrpcResponseIDAndError      *jrpcResponseIDAndErrorStruct
		jrpcResponseNoError         *jrpcResponseNoErrorStruct
		mountReply                  *jrpcfs.MountByAccountNameReply
		mountRequest                *jrpcfs.MountByAccountNameRequest
		storageURLSplit             []string
		thisMessageID               uint64
	)

	readEnv()

	thisMessageID = env.LastMessageID + 1
	env.LastMessageID = thisMessageID

	writeEnv()

	storageURLSplit = strings.Split(env.StorageURL, "/")
	if 0 == len(storageURLSplit) {
		log.Fatalf("Attempt to compute accountName from strings.Split(env.StorageURL, \"/\") failed")
	}

	accountName = storageURLSplit[len(storageURLSplit)-1]

	mountRequest = &jrpcfs.MountByAccountNameRequest{
		AccountName:  accountName,
		MountOptions: 0,
		AuthUserID:   0,
		AuthGroupID:  0,
	}

	jrpcRequest = &jrpcRequestStruct{
		JSONrpc: "2.0",
		Method:  "Server.RpcMountByAccountName",
		ID:      thisMessageID,
		Params:  [1]interface{}{mountRequest},
	}

	jrpcRequestBuf, err = json.Marshal(jrpcRequest)
	if nil != err {
		log.Fatalf("Attempt to marshal jrpcRequest(mount) failed: %v", err)
	}

	if verbose {
		_ = json.Indent(&jrpcRequestBytesBuffer, jrpcRequestBuf, "", "    ")
		fmt.Printf("jrpcRequestBuf:\n%s\n", jrpcRequestBytesBuffer.Bytes())
	}

	httpRequest, err = http.NewRequest("PROXYFS", env.StorageURL, bytes.NewReader(jrpcRequestBuf))
	if nil != err {
		log.Fatalf("Failed to create httpRequest for mount: %v", err)
	}

	httpRequest.Header["X-Auth-Token"] = []string{env.AuthToken}
	httpRequest.Header["Content-Type"] = []string{"application/json"}

	httpClient = &http.Client{}

	httpResponse, err = httpClient.Do(httpRequest)
	if nil != err {
		log.Fatalf("Failed to issue PROXYFS Mount: %v", err)
	}
	if http.StatusOK != httpResponse.StatusCode {
		log.Fatalf("Received unexpected HTTP Status for PROXYFS Mount: %s", httpResponse.Status)
	}

	httpResponseBody, err = ioutil.ReadAll(httpResponse.Body)
	if nil != err {
		log.Fatalf("Failed to read httpResponse.Body: %v", err)
	}
	err = httpResponse.Body.Close()
	if nil != err {
		log.Fatalf("Failed to close httpResponse.Body: %v", err)
	}

	if verbose {
		_ = json.Indent(&httpResponseBodyBytesBuffer, httpResponseBody, "", "    ")
		fmt.Printf("httpResponseBody:\n%s\n", httpResponseBodyBytesBuffer.Bytes())
	}

	jrpcResponseIDAndError = &jrpcResponseIDAndErrorStruct{}

	err = json.Unmarshal(httpResponseBody, jrpcResponseIDAndError)
	if nil != err {
		log.Fatalf("Failed to json.Unmarshal(httpResponseBody) [Case 1]: %v", err)
	}
	if thisMessageID != jrpcResponseIDAndError.ID {
		log.Fatalf("Got unexpected MessageID in httpResponseBody [Case 1]")
	}
	if "" != jrpcResponseIDAndError.Error {
		log.Fatalf("Got JRPC Failure on PROXYFS Mount: %v", err)
	}

	mountReply = &jrpcfs.MountByAccountNameReply{}
	jrpcResponseNoError = &jrpcResponseNoErrorStruct{Result: mountReply}

	err = json.Unmarshal(httpResponseBody, jrpcResponseNoError)
	if nil != err {
		log.Fatalf("Failed to json.Unmarshal(httpResponseBody) [Case 2]: %v", err)
	}
	if thisMessageID != jrpcResponseIDAndError.ID {
		log.Fatalf("Got unexpected MessageID in httpResponseBody [Case 2]")
	}

	env.MountID = mountReply.MountID

	writeEnv()
}

func doJRPC(s ...string) {
	var (
		arbitraryRequestMethod      string
		arbitraryRequestParam       string
		err                         error
		httpClient                  *http.Client
		httpRequest                 *http.Request
		httpResponse                *http.Response
		httpResponseBody            []byte
		httpResponseBodyBytesBuffer bytes.Buffer
		jrpcRequest                 *jrpcRequestStruct
		jrpcRequestBuf              []byte
		jrpcRequestBytesBuffer      bytes.Buffer
		jrpcResponseIDAndError      *jrpcResponseIDAndErrorStruct
		thisMessageID               uint64
	)

	readEnv()

	if "" == env.MountID {
		log.Fatalf("Attempt to issue JSON RPC to unmounted volume")
	}

	thisMessageID = env.LastMessageID + 1
	env.LastMessageID = thisMessageID

	writeEnv()

	arbitraryRequestMethod = s[0]

	switch len(s) {
	case 1:
		arbitraryRequestParam = "\"MountID\":\"" + string(env.MountID) + "\""
	case 2:
		arbitraryRequestParam = "\"MountID\":\"" + string(env.MountID) + "\"," + s[1]
	default:
		log.Fatalf("JSON RPC must be either [Method] or [Method Param]... not %v", s)
	}

	jrpcRequest = &jrpcRequestStruct{
		JSONrpc: "2.0",
		Method:  arbitraryRequestMethod,
		ID:      thisMessageID,
		Params:  [1]interface{}{&jrpcRequestEmptyParamStruct{}},
	}

	jrpcRequestBuf, err = json.Marshal(jrpcRequest)
	if nil != err {
		log.Fatalf("Attempt to marshal jrpcRequest failed: %v", err)
	}

	jrpcRequestBuf = append(jrpcRequestBuf[:len(jrpcRequestBuf)-3], arbitraryRequestParam...)
	jrpcRequestBuf = append(jrpcRequestBuf, "}]}"...)

	if verbose {
		_ = json.Indent(&jrpcRequestBytesBuffer, jrpcRequestBuf, "", "    ")
		fmt.Printf("jrpcRequestBuf:\n%s\n", jrpcRequestBytesBuffer.Bytes())
	}

	httpRequest, err = http.NewRequest("PROXYFS", env.StorageURL, bytes.NewReader(jrpcRequestBuf))
	if nil != err {
		log.Fatalf("Failed to create httpRequest: %v", err)
	}

	httpRequest.Header["X-Auth-Token"] = []string{env.AuthToken}
	httpRequest.Header["Content-Type"] = []string{"application/json"}

	httpClient = &http.Client{}

	httpResponse, err = httpClient.Do(httpRequest)
	if nil != err {
		log.Fatalf("Failed to issue PROXYFS: %v", err)
	}
	if http.StatusOK != httpResponse.StatusCode {
		log.Fatalf("Received unexpected HTTP Status: %s", httpResponse.Status)
	}

	httpResponseBody, err = ioutil.ReadAll(httpResponse.Body)
	if nil != err {
		log.Fatalf("Failed to read httpResponse.Body: %v", err)
	}
	err = httpResponse.Body.Close()
	if nil != err {
		log.Fatalf("Failed to close httpResponse.Body: %v", err)
	}

	if verbose {
		_ = json.Indent(&httpResponseBodyBytesBuffer, httpResponseBody, "", "    ")
		fmt.Printf("httpResponseBody:\n%s\n", httpResponseBodyBytesBuffer.Bytes())
	}

	jrpcResponseIDAndError = &jrpcResponseIDAndErrorStruct{}

	err = json.Unmarshal(httpResponseBody, jrpcResponseIDAndError)
	if nil != err {
		log.Fatalf("Failed to json.Unmarshal(httpResponseBody) [Case 2]: %v", err)
	}
	if thisMessageID != jrpcResponseIDAndError.ID {
		log.Fatalf("Got unexpected MessageID in httpResponseBody [Case 2]")
	}
	if "" != jrpcResponseIDAndError.Error {
		log.Fatalf("Got JRPC Failure: %v", err)
	}
}

func doClean() {
	_ = os.RemoveAll(envPath)
}

func writeEnv() {
	var (
		envBuf []byte
		err    error
	)

	envBuf, err = json.Marshal(env)
	if nil != err {
		log.Fatalf("Failed to json.Marshal(env): %v", err)
	}

	err = ioutil.WriteFile(envPath, envBuf, 0644)
	if nil != err {
		log.Fatalf("Failed to persist ENV to envPath==\"%s\": %v", envPath, err)
	}
}

func readEnv() {
	var (
		envBuf []byte
		err    error
	)

	env = &envStruct{}
	envBuf, err = ioutil.ReadFile(envPath)
	if nil != err {
		log.Fatalf("Failed to recover ENV from envPath==\"%s\": %v", envPath, err)
	}

	err = json.Unmarshal(envBuf, env)
	if nil != err {
		log.Fatalf("Failed to json.Unmarshal(envBuf): %v", err)
	}
}
