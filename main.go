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
    log.Fatal(http.ListenAndServe(":10000", nil))
}



func main() {
    handleRequests()
}
