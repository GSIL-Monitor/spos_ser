// IndustryCategory
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
)

func AddIndustryCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddIndustryCategoryHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("name")) > 0 {

			name := r.FormValue("name")
			describe := r.FormValue("describe")

			fmt.Printf("AddIndustryCategoryHandler: name=%s;describe=%s\n", name, describe)

			if len(name) == 0 {
				fmt.Printf("AddIndustryCategoryHandler: IndustryCategory info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddIndustryCategoryHandler: IndustryCategory info error!"))
				return
			}

			if datas.IsExistIndustryCategory(name) {
				fmt.Printf("AddIndustryCategoryHandler: ERROR,IndustryCategory is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddIndustryCategoryHandler: ERROR,IndustryCategory is Exist!"))
				return
			}

			fmt.Printf("AddIndustryCategoryHandler: new IndustryCategory codeing...\n")

			ret := datas.AddIndustryCategory(name, describe)

			if ret != nil {
				fmt.Printf("AddIndustryCategoryHandler: IndustryCategory connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddIndustryCategoryHandler: IndustryCategory connext error, db record write error!"))
				return
			}
		} else {
			fmt.Printf("AddIndustryCategoryHandler: IndustryCategory connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddIndustryCategoryHandler: IndustryCategory connext error!"))
			return
		}
	} else {
		fmt.Printf("AddIndustryCategoryHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddIndustryCategoryHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("IndustryCategory create success!"))
	fmt.Printf("AddIndustryCategoryHandler: IndustryCategory create success!\n")

}
