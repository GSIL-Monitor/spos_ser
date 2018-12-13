// Role
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
)

func AddRoleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddRoleHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("name")) > 0 {

			name := r.FormValue("name")
			describe := r.FormValue("describe")
			fmt.Printf("AddRoleHandler: name=%s;describe=%s\n", name, describe)

			if len(name) == 0 {
				fmt.Printf("AddRoleHandler: Role info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleHandler: Role info error!"))
				return
			}

			if datas.IsExistOrganization(name) {
				fmt.Printf("AddRoleHandler: ERROR,Role is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleHandler: ERROR,Role is Exist!"))
				return
			}

			fmt.Printf("AddRoleHandler: new Role codeing...\n")

			ret := datas.AddRole(name, describe)

			if ret != nil {
				fmt.Printf("AddRoleHandler: Role connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleHandler: Role connext error, db record write error!"))
				return
			}
		} else {
			fmt.Printf("AddRoleHandler: Role connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddRoleHandler: Role connext error!"))
			return
		}
	} else {
		fmt.Printf("AddRoleHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddRoleHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("Role create success!"))
	fmt.Printf("AddRoleHandler: Role create success!\n")

}
