// UserOrg
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
	"spos_auth/enforcer"
)

func AddUserOrgHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddUserOrgHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("uid")) > 0 {

			uid := r.FormValue("uid")
			orid := r.FormValue("orid")
			deid := r.FormValue("deid")

			fmt.Printf("AddUserOrgHandler: uid=%s;orid=%s;deid=%s\n", uid, orid, deid)

			if len(uid) == 0 || len(orid) == 0 || len(deid) == 0 {
				fmt.Printf("AddUserOrgHandler: UserOrg info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserOrgHandler: UserOrg info error!"))
				return
			}

			if datas.IsExistUserOrg(uid, orid, deid) {
				fmt.Printf("AddUserOrgHandler: ERROR,UserOrg is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserOrgHandler: ERROR,UserOrg is Exist!"))
				return
			}

			fmt.Printf("AddUserOrgHandler: new UserOrg codeing...\n")
			ret := datas.AddUserOrg(uid, orid, deid)

			if ret != nil {
				fmt.Printf("AddUserOrgHandler: UserOrg connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserOrgHandler: UserOrg connext error, db record write error!"))
				return
			} else {
				if enforcer.AddRoleForUser(uid, deid) {
					fmt.Printf("AddUserOrgHandler: Add department of user success!\n")

				} else {
					fmt.Printf("AddUserOrgHandler: Add department of user error!\n")

				}
			}
		} else {
			fmt.Printf("AddUserOrgHandler: UserOrg connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddUserOrgHandler: UserOrg connext error!"))
			return
		}
	} else {
		fmt.Printf("AddUserOrgHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddUserOrgHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("UserOrg create success!"))
	fmt.Printf("AddUserOrgHandler: UserOrg create success!\n")
}
