// Organization
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
)

func AddOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddOrganizationHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("name")) > 0 {
			fmt.Printf("AddOrganizationHandler: name=%s\n", r.FormValue("name"))

			name := r.FormValue("name")
			icid := r.FormValue("icid")
			agentname := r.FormValue("agentname")
			agentphone := r.FormValue("agentphone")
			agentmail := r.FormValue("agentmail")
			describe := r.FormValue("describe")

			fmt.Printf("AddOrganizationHandler: name=%s,icid=%s,agentname=%s,agentphone=%s,agentmail=%s,describe=%s\n", name, icid, agentname, agentphone, agentmail, describe)

			if len(name) == 0 || len(icid) == 0 ||
				len(agentname) == 0 || len(agentphone) == 0 || len(agentmail) == 0 {
				fmt.Printf("AddOrganizationHandler: Organization info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddOrganizationHandler: Organization info error!"))
				return
			}

			if datas.IsExistOrganization(name) {
				fmt.Printf("AddOrganizationHandler: ERROR,Organization is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddOrganizationHandler: ERROR,Organization is Exist!"))
				return
			}

			fmt.Printf("AddOrganizationHandler: new Organization codeing...\n")

			ret := datas.AddOrganization(icid, name, agentname, agentphone, agentmail, describe)

			if ret != nil {
				fmt.Printf("AddOrganizationHandler: Organization connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddOrganizationHandler: Organization connext error, db record write error!"))
				return
			}

		} else {
			fmt.Printf("AddOrganizationHandler: Organization connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddOrganizationHandler: Organization connext error!"))
			return
		}
	} else {
		fmt.Printf("AddOrganizationHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddOrganizationHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("Organization create success!"))
	fmt.Printf("AddOrganizationHandler: Organization create success!\n")

}
