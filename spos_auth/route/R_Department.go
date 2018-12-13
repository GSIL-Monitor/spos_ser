// Department
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
)

func AddDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddDepartmentHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("name")) > 0 {

			orid := r.FormValue("orid")
			name := r.FormValue("name")
			describe := r.FormValue("describe")

			fmt.Printf("AddDepartmentHandler: name=%s;orid=%s;describe=%s\n", name, orid, describe)

			if len(name) == 0 || len(orid) == 0 {
				fmt.Printf("AddDepartmentHandler: Department info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentHandler: Department info error!"))
				return
			}

			if datas.IsExistDepartment(name) {
				fmt.Printf("AddDepartmentHandler: ERROR,Department is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentHandler: ERROR,Department is Exist!"))
				return
			}

			fmt.Printf("AddDepartmentHandler: new Department codeing...\n")

			ret := datas.AddDepartment(orid, name, describe)

			if ret != nil {
				fmt.Printf("AddDepartmentHandler: Department connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddDepartmentHandler: Department connext error, db record write error!"))
				return
			}
		} else {
			fmt.Printf("AddDepartmentHandler: Department connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddDepartmentHandler: Department connext error!"))
			return
		}
	} else {
		fmt.Printf("AddDepartmentHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddDepartmentHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("Department create success!"))
	fmt.Printf("AddDepartmentRoleHandler: Department create success!\n")

}
