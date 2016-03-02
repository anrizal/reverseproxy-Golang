package main

//import package
import (
    "fmt"
    "os"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strconv"
    "strings"
)

//data type declaration for config
type Config struct {
    Target      string
    ListenOn    string
    StatFolder  string
    CreateNewFileEveryXMinutes int32
}

//variable declaration
var conf Config
//start date point
var tiN time.Time
//current file name
var fname string

func getConfig(cFile string) {

    //open a file
    content, err := ioutil.ReadFile(cFile)
    if err!=nil{
        fmt.Print("Error:",err)
    }
    //decode json
    err = json.Unmarshal([]byte(content), &conf)
    if err!=nil{
        fmt.Print("Error:",err)
    }
}

//check error function
func check(e error) {
  if e != nil {
    log.Fatal(e)
    panic(e)
    //os.Exit(1)
  }
}

/*
Description: 
Utility function to copy header information

Input:
ResponseWriter
Request

Return:
return formated information
*/

//copy header function
func copyHeader(source http.Header, dest *http.Header){
  for n, v := range source {
      for _, vv := range v {
          dest.Add(n, vv)
      }
  }
}

/*
Description: 
request handler implementation

Input:
ResponseWriter
Request

Return:
return formated information
*/

func logStat(w http.ResponseWriter, r *http.Request){

//get target + the stuff after it
  uri := conf.Target+r.RequestURI

  //for logging- need to be REMOVED
  fmt.Println(r.Method + ": " + uri)

  //remove session
  cu := strings.Split(uri, "?")
//create new request to the target server with information from end user's browser
//The Information are method, header and body
  rr, err := http.NewRequest(r.Method, cu[0], r.Body)
  check(err)

//copy the header that we get from the browser
  copyHeader(r.Header, &rr.Header)

 // Create a client and query the target
 var transport http.Transport
 // no need to use http.Client.Do since we dont need to follow redirects nor handles cookies
 // and transport.RoundTrip is roughly 60% faster since it only makes single HTTP transaction
  resp, err := transport.RoundTrip(rr)
  check(err)

//Write stat
  writeStat(r.RequestURI)

  defer resp.Body.Close()
//read response body we got fromserver
  body, err := ioutil.ReadAll(resp.Body)
  check(err)

//create http.ResponseWriter Header
  dH := w.Header()
  copyHeader(resp.Header, &dH)

  w.Write(body)

}

/*
Description: 
Decide how new information will be written. 
Is it in the same file or new file

Input:
The formated result

Return:
return error if exist
*/

func writeStat(data string) {

 //declare time duration it needs to wait
 var gap time.Duration = time.Duration(conf.CreateNewFileEveryXMinutes) * time.Minute
 
 //calculate duration 
 duration := time.Since(tiN)

 if(duration > gap) {
      //if gap has been passed
      //reset time start and rename the file name that we are working on
      tiN = time.Now()
      fname = conf.StatFolder +"log-" + tiN.Format("2006.01.02-15:04:05")+ ".txt"
      fmt.Printf("new file created")      
  } 

  //open the file
  f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
  check(err)
  //split string, after ? is the user session
  s := strings.Split(data, "?")
  //write inside a file
  _, err = f.WriteString(strconv.FormatInt(time.Now().Unix(), 10)+","+s[1]+","+s[0]+"\n")
  check(err)

  defer f.Close()
  
}

/**
main function
*/
func main() {
  getConfig("config/conf.json") //get config
  tiN = time.Now()
  //reset time start and rename the file name that we are working on
  fname = conf.StatFolder +"log-" + tiN.Format("2006.01.02-15:04:05")+ ".txt"
  http.HandleFunc("/", logStat) //request handler for query stat
  log.Fatal(http.ListenAndServe(conf.ListenOn, nil)) //listen incoming request
}




