// UserRole
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
	"spos_auth/enforcer"
)

func AddUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddUserRoleHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("uid")) > 0 {

			rid := r.FormValue("rid")
			uid := r.FormValue("uid")

			fmt.Printf("AddUserRoleHandler: uid=%s;rid=%s\n", uid, rid)

			if len(rid) == 0 || len(uid) == 0 {
				fmt.Printf("AddUserRoleHandler: UserRole info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserRoleHandler: UserRole info error!"))
				return
			}

			if datas.IsExistUserRole(rid, uid) {
				fmt.Printf("AddUserRoleHandler: ERROR,UserRole is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserRoleHandler: ERROR,UserRole is Exist!"))
				return
			}

			fmt.Printf("AddUserRoleHandler: new UserRole codeing...\n")
			ret := datas.AddUserRole(uid, rid)
			if ret != nil {
				fmt.Printf("AddUserRoleHandler: UserRole connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddUserRoleHandler: UserRole connext error, db record write error!"))
				return
			} else {
				if enforcer.AddRoleForUser(uid, rid) {
					fmt.Printf("AddUserRoleHandler: Add role of user success!\n")

				} else {
					fmt.Printf("AddUserRoleHandler: Add role of user error!\n")

				}
			}
		} else {
			fmt.Printf("AddUserRoleHandler: UserRole connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddUserRoleHandler: UserRole connext error!"))
			return
		}
	} else {
		fmt.Printf("AddUserRoleHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddUserRoleHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("UserRole create success!"))
	fmt.Printf("AddUserRoleHandler: UserRole create success!\n")
}
