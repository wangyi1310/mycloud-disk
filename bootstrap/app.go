package bootstrap

import "fmt"

func InitApplication() {
	fmt.Print(`
   ___ _                 _                    
  / __\ | ___  _   _  __| |_ __ _____   _____ 
 / /  | |/ _ \| | | |/ _  | '__/ _ \ \ / / _ \	
/ /___| | (_) | |_| | (_| | | |  __/\ V /  __/
\____/|_|\___/ \__,_|\__,_|_|  \___| \_/ \___|

   V` + conf.BackendVersion + `  Commit #` + conf.LastCommit + `  Pro=` + conf.IsPro + `
================================================

`)
	go CheckUpdate()
}