// DepartmentOR
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
	"spos_auth/enforcer"
)

func AddDepartmentORHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddDepartmentORHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("deid")) > 0 {

			oid := r.FormValue("oid")
			deid := r.FormValue("deid")

			fmt.Printf("AddDepartmentORHandler: deid=%s;oid=%s\n", deid, oid)

			if len(oid) == 0 || len(deid) == 0 {
				fmt.Printf("AddDepartmentORHandler: DepartmentOR info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentORHandler: DepartmentOR info error!"))
				return
			}

			if datas.IsExistDepartmentOR(deid, oid) {
				fmt.Printf("AddDepartmentORHandler: ERROR,DepartmentOR is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentORHandler: ERROR,DepartmentOR is Exist!"))

				DepartmentAddPolicy(deid, oid)
				return
			}

			fmt.Printf("AddDepartmentORHandler: new DepartmentOR codeing...\n")
			ret := datas.AddDepartmentOR(deid, oid)
			if ret != nil {
				fmt.Printf("AddDepartmentORHandler: DepartmentOR connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentORHandler: DepartmentOR connext error, db record write error!"))
				return
			} else {
				DepartmentAddPolicy(deid, oid)
			}
		} else {
			fmt.Printf("AddDepartmentORHandler: DepartmentOR connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddDepartmentORHandler: DepartmentOR connext error!"))
			return
		}
	} else {
		fmt.Printf("AddDepartmentORHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddDepartmentORHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("DepartmentOR create success!"))
	fmt.Printf("AddDepartmentORHandler: DepartmentOR create success!\n")
}

func DepartmentAddPolicy(deid, oid string) {
	or := datas.GetObjectResource(oid)
	if enforcer.AddPolicy(deid, or.Content, or.Action) {
		fmt.Printf("AddRoleObjectHandler: Add Policy of Role Object success!\n")
	} else {
		fmt.Printf("AddRoleObjectHandler: Add Policy of Role Object error!\n")
	}
}
