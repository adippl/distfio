/*
 *  distfio - program designed to launch parallel fio tests on many clients
 *  Copyright (C) 2020 Adam Prycki (email: adam.prycki@gmail.com)
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */
package main

import "fmt"
import "log"
import "net/http"
import "os/exec"
import "os"
import "encoding/json"
import "io/ioutil"
import "strings"
import "time"
import "bytes"

var client_args Args
var job_args Args

type Fio_args struct {
	Direct		string `json:"direct"`
	Rw			string `json:"rw"`
	Bs			string `json:"bs"`
	Ioengine	string `json:"ioengine"`
	Iodepth		string `json:"iodepth"`
	Runtime		string `json:"runtime"`
	Numjobs		string `json:"numjobs"`
	
	// s3 specific options
	Http_verbose	string `json:"http_verbose"`
	Https			string `json:"https"`
	Http_s3_key		string `json:"http_s3_key"`
	Http_s3_keyid	string `json:"http_s3_keyid"`
	Http_host		string `json:"http_host"`
	Http_s3_region	string `json:"http_s3_region"`
	}

type Fio_job_options struct {
	Name	string `json:"name"`
	Size	string `json:"size"`
	}

type Fio_iodepth_level struct {
	L1		float64 `json:"1"`
	L2		float64 `json:"2"`
	L4		float64 `json:"4"`
	L8		float64 `json:"8"`
	L16		float64 `json:"16"`
	L32		float64 `json:"32"`
	L_more_eq_than_64	float64 `json:">=64"`
}

type Fio_latency struct {
	L2		float64 `json:"2"`
	L4		float64 `json:"4"`
	L10		float64 `json:"10"`
	L20		float64 `json:"20"`
	L50		float64 `json:"50"`
	L100	float64 `json:"100"`
	L250	float64 `json:"250"`
	L500	float64 `json:"500"`
	L750	float64 `json:"750"`
	L1000	float64 `json:"1000"`
}

type Fio_job struct {
	Jobname				string `json:"jobname"`
	Groupid				int `json:"groupid"`
	Job_error			int `json:"error"`
	Eta					int `json:"eta"`
	Elapsed				int `json:"elapsed"`
	Job_options			Fio_job_options `json:"job options"`
	Read				Fio_stats `json:"read"`
	Write				Fio_stats `json:"write"`
	Trim				Fio_stats `json:"trim"`
	Sync				Fio_stats `json:"sync"`
	Job_runtime			int64 `json:"job_runtime"`
	Usr_cpu				float64 `json:"usr_cpu"`
	Sys_cpu				float64 `json:"sys_cpu"`
	Ctx					int64 `json:"ctx"`
	Majf				int `json:"majf"`
	Minf				int64 `json:"minf"`
	Iodepth_level		Fio_iodepth_level `json:"iodepth_level"`
	Iodepth_submit		Fio_iodepth_level `json:"iodepth_submit"`
	Iodepth_complete	Fio_iodepth_level `json:"iodepth_complete"`
	Latency_ns			Fio_latency `json:"latency_ns"`
	Latency_us			Fio_latency `json:"latency_us"`
	Latency_ms			Fio_latency `json:"latency_ms"`
	Latency_depth		int64 `json:"latency_depth"`
	Latency_target		int64 `json:"latency_target"`
	Latency_percentile	float64 `json:"latency_percentile"`
	Latency_window		int64 `json:"latency_window"`
}

type Fio_slat_ns struct {
	Min		int64 `json:"min"`
	Max		int64 `json:"max"`
	Mean	float64 `json:"mean"`
	Stdev	float64 `json:"stdev"`
	N		int64 `json:"N		int64"`
}

type Fio_percentile struct {
	P_1_000000		int64 `json:"1.000000"`
	P_5_000000		int64 `json:"5.000000"`
	P_10_000000		int64 `json:"10.000000"`
	P_20_000000		int64 `json:"20.000000"`
	P_30_000000		int64 `json:"30.000000"`
	P_40_000000		int64 `json:"40.000000"`
	P_50_000000		int64 `json:"50.000000"`
	P_60_000000		int64 `json:"60.000000"`
	P_70_000000		int64 `json:"70.000000"`
	P_80_000000		int64 `json:"80.000000"`
	P_90_000000		int64 `json:"90.000000"`
	P_95_000000		int64 `json:"95.000000"`
	P_99_000000		int64 `json:"99.000000"`
	P_99_500000		int64 `json:"99.500000"`
	P_99_900000		int64 `json:"99.900000"`
	P_99_950000		int64 `json:"99.950000"`
	P_99_990000		int64 `json:"99.990000"`
}


type Fio_clat_ns struct {
	Min			int64 `json:"min"`
	Max			int64 `json:"max"`
	Mean		float64 `json:"mean"`
	Stdev		float64 `json:"stdev"`
	N			int64 `json:"N			int64"`
	Percentile	Fio_percentile `json:"percentile"`
}


type Fio_lat_ns struct {
	Min		int64 `json:"min"`
	Max		int64 `json:"max"`
	Mean	float64 `json:"mean"`
	Stdev	float64 `json:"stdev"`
	N		int64 `json:"N"`
}

type Fio_stats struct {
	Io_bytes		int64 `json:"io_bytes"`
	Io_kbytes		int64 `json:"io_kbytes"`
	Bw_bytes		int64 `json:"bw_bytes"`
	Bw				int64 `json:"bw"`
	Iops			float64 `json:"iops"`
	Runtime			int64 `json:"runtime"`
	Total_iop		int64 `json:"total_iop"`
	Short_ios		int64 `json:"short_ios"`
	Drop_ios		int64 `json:"drop_ios"`
	Slat_ns			Fio_slat_ns `json:"slat_ns"`
	Clat_ns			Fio_clat_ns `json:"clat_ns"`
	Lat_ns			Fio_lat_ns `json:"lat_ns"`
	Bw_min			int64 `json:"bw_min"`
	Bw_max			int64 `json:"bw_max"`
	Bw_agg			float64 `json:"bw_agg"`
	Bw_mean			float64 `json:"bw_mean"`
	Bw_dev			float64 `json:"bw_dev"`
	Bw_samples		int64 `json:"bw_samples"`
	Iops_min		int64 `json:"iops_min"`
	Iops_max		int64 `json:"iops_max"`
	Iops_mean		float64 `json:"iops_mean"`
	Iops_stddev		float64 `json:"iops_stddev"`
	Iops_samples	int64 `json:"iops_samples"`
}

type Fio_results struct {
	Fio_version		string `json:"fio version"`
	Timestamp		int64 `json:"timestamp"`
	Timestamp_ms	int64 `json:"timestamp_ms"`
// commented out. fio date format causes unmarshall-er error
//	Time			time.Time `json:"time"`
	Global_options	Fio_args `json:"global options"`
	Jobs			[]Fio_job `json:"jobs"`
}

type Node struct {
	Hostname	string
	Port		int
	}

type Args struct {
	Nodes		[]Node
	Dir			string
	Device		string
	Bs			string
	Iodepth		string
	Direct		string
	Ioengine	string
	Runtime		string
	Numjobs		string
	Name		string
	Rwmode		string
	Group_reporting		bool
	Time_based			bool
	}

func write_example_args(){
	var err error
	example_args := Args{
		Nodes:		[]Node{
			Node{
				Hostname:	"localhost",
				Port:		10000,
				},
			},
		Dir:		"/home/adip/tmp/fio/",
		Device:		"",
		Bs:			"4k",
		Iodepth:	"32",
		Direct:		"1",
		Ioengine:	"direct",
		Runtime:	"",
		Numjobs:	"8",
		Name:		"raw_randread",
		Rwmode:		"randread",
		}
	file,err := json.MarshalIndent(example_args, "", "\t")
	if err != nil {
		fmt.Println(err)
		panic(err)}
	
	fmt.Println("writing example config to ./arguments_example.json")
	err = ioutil.WriteFile("./arguments_example.json", file, 0644)
	if err != nil {
		panic(err)}}


func client_load_config(){
	var err error
	var raw []byte
//	var args Args
	
	if len(os.Args) < 2 {
		fmt.Println("config file not specified \nexiting...")
		os.Exit(2)}
	
	raw,err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Argv[1] should point to file with client config")
		panic(err)}
	err = json.Unmarshal(raw, &client_args)
	if err != nil {
		fmt.Println("failed to unmarshal client config")
		fmt.Println(err)
		os.Exit(3)}
	fmt.Printf("loaded config %#v\n", client_args)}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")}

func handle_fio_job_submission(w http.ResponseWriter, r *http.Request){
	var err  error
	var dec  *json.Decoder
	var args Args

	log.Println("\n\nEndpoint Hit: handle_fio_job_submission")
	//debug	log.Println(r.Header.Get("Content-Type"))
	if r.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		fmt.Fprintf(w, "fio job submission page")
		return}
	//limit request to 1MiB
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	//	//debug
	//	buf := new(bytes.Buffer)
	//	buf.ReadFrom(r.Body)
	//	log.Println(buf.String())
	//	//debug
	dec = json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&args)
	if err!=nil {
		log.Println(err)}
	log.Println("received job submission")
	log.Printf("%#v\n", args)
	//TODO start job and send result in response
	//delay simulatin job work
	time.Sleep(time.Second * time.Duration(5))
	}

func client_fio_job_submission(){
	var err			error
	var jsonData	[]byte
	var req			*http.Request
	var postUrl		string
	var resp		*http.Response
	var resp_data	[]byte
	var client		*http.Client
	
	jsonData,err = json.Marshal(client_args)
	if err!=nil {
		log.Fatal(err)}
	
	
	for _,v:=range client_args.Nodes{
		postUrl = fmt.Sprintf("http://%s:%d/job_submission",
			v.Hostname,
			v.Port)
		req,err = http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client = &http.Client{}
		resp,err = client.Do(req)
		if err != nil {
			log.Fatal(err)
			os.Exit(105)}
		resp_data,err = ioutil.ReadAll(resp.Body)
		if err!=nil {
			log.Fatal(err)}
		log.Println(string(resp_data))}
	}

func fio_processes_dump(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Fio process dump\n")
	fmt.Println("Endpoint Hit: fio_processes_dump")
	cmd := exec.Command("bash", "-c", "ps aux |grep fio")
	stdout,err := cmd.Output()
	if err!=nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err.Error())}
	fmt.Fprintf(w, string(stdout))}

func delayed_exit(ms int){
	time.Sleep(time.Millisecond * time.Duration(ms))
	// TODO kill fio benchmarks
	os.Exit(0)}
	
func handle_exit(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "killing distfio server")
	fmt.Println("Endpoint Hit: /exit")
	go delayed_exit(100)}

func handleRequests(){
	http.HandleFunc("/", homePage)
	http.HandleFunc("/fio_dump", fio_processes_dump)
	http.HandleFunc("/exit", handle_exit)
	http.HandleFunc("/job_submission", handle_fio_job_submission)
	log.Fatal(http.ListenAndServe(":10000", nil))}

func client_req_kill(){
	var err error
	var resp *http.Response
	var resp_data []byte
	
	if len(client_args.Nodes) == 0 {
		log.Println("Node list is empty, couldn't kill any")
		return}
	for _,v:=range client_args.Nodes{
		resp,err = http.Get(fmt.Sprintf("http://%s:%d/exit",
			v.Hostname,
			v.Port))
		if err != nil {
			log.Println(err)
			os.Exit(105)}
		resp_data,err = ioutil.ReadAll(resp.Body)
		if err!=nil {
			log.Fatal(err)}
		log.Println(string(resp_data))}
	log.Println("kill commands sent do all job servers.\nexiting...")
	os.Exit(0)}


func client_arguments_handle(){
	if len(os.Args) <3 {
		fmt.Println("clientd didn't receive any arguments")
		fmt.Println("exiting...")
		os.Exit(4)}
	switch os.Args[2] {
		case "kill":
			fmt.Println("killing jobs servers")
			client_req_kill()
			//TODO kill job servers
		// TODO kill all JOBS on servers
		case "start":
			client_fio_job_submission()
		default:
			fmt.Printf("'%s' command doesn't exist", os.Args[2])}}


func test_fio_json_ingest(){
	var results Fio_results
	var file_path string = "fio_results_example.json"
	
	fmt.Printf("\n\nfio json output file test\n")
	jsonFile, err := os.Open(file_path)
	if err != nil {
	    fmt.Println(err)}
	defer jsonFile.Close()
	
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("ERR reading cluster.json")
		os.Exit(10);}
		
	err = nil
	err = json.Unmarshal(raw,&results)
	if err!=nil {
		fmt.Println("ERR while trying to unmarshal test file", err)
		os.Exit(11);}
	fmt.Printf("test_fio_json_ingest loaded results from file %s \n %+v\n",
		file_path,
		&results)
	fmt.Printf("\nend of fio json output file test\n\n")}


func server() {
	test_fio_json_ingest()
	handleRequests()
}

func client() {
	fmt.Println("client not implemented yet")
	write_example_args()
	client_load_config()
	client_arguments_handle()
}

func main() {
	var args0split []string
	var progname string
	
	// get name of the program
	fmt.Println("arg0", os.Args[0])
	args0split = strings.Split(os.Args[0], "/")
	progname = args0split[len(args0split)-1]
	
	if progname == "distfio" {
		fmt.Println("starting in daemon mode")
		server()
		os.Exit(0)
	}else if progname == "distfio-client" {
		fmt.Println("started in client mode")
		client()
		os.Exit(0)
	}else{
		fmt.Fprintln(os.Stderr, "incorrect filename! ", progname)
		os.Exit(1)}}
