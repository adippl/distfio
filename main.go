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
//import "encoding/json"
import "time"

type fio_args struct {
	direct		int64
	rw			string
	bs			string
	ioengine	string
	iodepth		int64
	runtime		int64
	numjobs		int64
	}

type fio_job_options struct {
	name	string
	size	string
	}

type fio_job struct {
	jobname			string
	groupid			int
	job_error		int `json:"error"`
	eta				int
	elapsed			int
	job_options		fio_job_options `json:"job options"`
}

type fio_slat_ns struct {
	min		int64
	max		int64
	mean	float64
	stdev	float64
	N		int64
}

type fio_percentile struct {
	p_1_000000		int64 `json:"1.000000"`
	p_5_000000		int64 `json:"5.000000"`
	p_10_000000		int64 `json:"10.000000"`
	p_20_000000		int64 `json:"20.000000"`
	p_30_000000		int64 `json:"30.000000"`
	p_40_000000		int64 `json:"40.000000"`
	p_50_000000		int64 `json:"50.000000"`
	p_60_000000		int64 `json:"60.000000"`
	p_70_000000		int64 `json:"70.000000"`
	p_80_000000		int64 `json:"80.000000"`
	p_90_000000		int64 `json:"90.000000"`
	p_95_000000		int64 `json:"95.000000"`
	p_99_000000		int64 `json:"99.000000"`
	p_99_500000		int64 `json:"99.500000"`
	p_99_900000		int64 `json:"99.900000"`
	p_99_950000		int64 `json:"99.950000"`
	p_99_990000		int64 `json:"99.990000"`
}


type fio_clat_ns struct {
	min			int64
	max			int64
	mean		float64
	stdev		float64
	N			int64
	percentile	fio_percentile }


type fio_lat_ns struct {
	min		int64
	max		int64
	mean	float64
	stdev	float64
	N		int64 }

type fio_read struct {
	io_bytes		int64
	io_kbytes		int64
	bw_bytes		int64
	bw				int64
	iops			float64
	runtime			int64
	total_iop		int64
	short_ios		int64
	drop_ios		int64
	slat_ns			fio_slat_ns
	clat_ns			fio_clat_ns
	lat_ns			fio_lat_ns
	bw_min			int64
	bw_max			int64
	bw_agg			float64
	bw_mean			float64
	bw_dev			float64
	bw_samples		int64
	iops_min		int64
	iops_max		int64
	iops_mean		float64
	iops_stddev		float64
	iops_samples	int64
	}



type fio_results struct {
	fio_version		string `json:"fio version"`
	timestmap		int64
	timestamp_ms	int64
	time			time.Time
	fio_args		fio_args `json:"global options"`
	jobs			[]fio_job
	}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")}


func fio_processes_dump(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Fio process dump\n")
	fmt.Println("Endpoint Hit: fio_processes_dump")
	cmd := exec.Command("bash", "-c", "ps aux |grep fio")
	stdout,err := cmd.Output()
	if err!=nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err.Error())}
	fmt.Fprintf(w, string(stdout))}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/fio_dump", fio_processes_dump)
	log.Fatal(http.ListenAndServe(":10000", nil))}

func test_fio_json_ingest(){
	var file string = "fio_results_example.json"
	
	jsonFile, err := os.Open("users.json")
	if err != nil {
	    fmt.Println(err)}
	defer jsonFile.Close()
	}



func main() {
    handleRequests()
}
