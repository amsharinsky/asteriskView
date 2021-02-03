package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)


type AsteriskView struct{


	AsteriskAgent []string `json:"asteriskAgent"`
	Token string

}



func (AstView *AsteriskView) makeLog(errlog error){

	logfile,err:=os.OpenFile("/var/log/AsteriskView.log",os.O_RDWR|os.O_CREATE| os.O_APPEND, 0755)
	if err != nil{
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println(errlog.Error())

}
func (AstView *AsteriskView) openConfig(){

	f,err:=os.Open("./config/config.json")
	if err != nil{

		AstView.makeLog(err)

	}
	file,err:=ioutil.ReadAll(f)
	if err != nil{

		AstView.makeLog(err)
	}
	json.Unmarshal(file,AstView)


}

func (AstView *AsteriskView) handler(w http.ResponseWriter,r *http.Request){

	data:=make(map[string]map[string]string)

	transport:=&http.Transport{
		MaxIdleConns: 10,
		IdleConnTimeout: 5*time.Second ,
		}

	client:=&http.Client{Transport: transport}

	for _,value:=range AstView.AsteriskAgent{

		res,err:=client.Get(value+"?token="+AstView.Token)
		if err != nil{
			AstView.makeLog(err)
			continue
		}
		json.NewDecoder(res.Body).Decode(&data)

	}

	tpl:=template.Must(template.ParseFiles("./html/index.html"))
	tpl.Execute(w,data)

}

func main(){
	var astView AsteriskView
	astView.Token="b52c96bea30646abf8170f333bbd42b9"
	astView.openConfig()
	http.HandleFunc("/",astView.handler)
	err:=http.ListenAndServe(":8081",nil)
	if err != nil{
		astView.makeLog(err)
	}
}
