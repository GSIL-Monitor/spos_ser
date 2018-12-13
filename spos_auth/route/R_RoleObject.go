// RoleObject
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
	"spos_auth/enforcer"
)

func AddRoleObjectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddRoleObjectHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("rid")) > 0 {

			rid := r.FormValue("rid")
			oid := r.FormValue("oid")

			fmt.Printf("AddRoleObjectHandler: rid=%s; oid=%s\n", rid, oid)

			if len(rid) == 0 || len(oid) == 0 {
				fmt.Printf("AddRoleObjectHandler: RoleObject info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleObjectHandler: RoleObject info error!"))
				return
			}

			if datas.IsExistRoleObject(rid, oid) {
				fmt.Printf("AddRoleObjectHandler: ERROR,RoleObject is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleObjectHandler: ERROR,RoleObject is Exist!"))

				RoleAddPolicy(rid, oid)
				return
			}

			fmt.Printf("AddRoleObjectHandler: new RoleObject codeing...\n")
			ret := datas.AddRoleObject(rid, oid)

			if ret != nil {
				fmt.Printf("AddRoleObjectHandler: RoleObject connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddRoleObjectHandler: RoleObject connext error, db record write error!"))
				return
			} else {
				RoleAddPolicy(rid, oid)
			}

		} else {
			fmt.Printf("AddRoleObjectHandler: RoleObject connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddRoleObjectHandler: RoleObject connext error!"))
			return
		}
	} else {
		fmt.Printf("AddRoleObjectHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddRoleObjectHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("RoleObject create success!"))
	fmt.Printf("AddRoleObjectHandler: RoleObject create success!\n")
}

func RoleAddPolicy(rid, oid string) {
	or := datas.GetObjectResource(oid)
	if enforcer.AddPolicy(rid, or.Content, or.Action) {
		fmt.Printf("AddRoleObjectHandler: Add Policy of Role Object success!\n")
	} else {
		fmt.Printf("AddRoleObjectHandler: Add Policy of Role Object error!\n")
	}
}
