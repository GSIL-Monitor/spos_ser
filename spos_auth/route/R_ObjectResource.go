// ObjectResource
package route

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
)

func AddObjectResourceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AddObjectResourceHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("name")) > 0 {

			name := r.FormValue("name")
			content := r.FormValue("content")
			action := r.FormValue("action")
			describe := r.FormValue("describe")

			fmt.Printf("AddObjectResourceHandler: name=%s;content=%s;action=%s;describe=%s\n", name, content, action, describe)

			if len(name) == 0 || len(content) == 0 {
				fmt.Printf("AddObjectResourceHandler: ObjectResource info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddObjectResourceHandler: ObjectResource info error!"))
				return
			}

			if datas.IsExistObjectResource(name) {
				fmt.Printf("AddObjectResourceHandler: ERROR,ObjectResource is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddObjectResourceHandler: ERROR,ObjectResource is Exist!"))
				return
			}

			fmt.Printf("AddObjectResourceHandler: new ObjectResource codeing...\n")

			ret := datas.AddObjectResource(name, content, action, describe)

			if ret != nil {
				fmt.Printf("AddObjectResourceHandler: ObjectResource connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("AddObjectResourceHandler: ObjectResource connext error, db record write error!"))
				return
			}
		} else {
			fmt.Printf("AddObjectResourceHandler: ObjectResource connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("AddObjectResourceHandler: ObjectResource connext error!"))
			return
		}
	} else {
		fmt.Printf("AddObjectResourceHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("AddObjectResourceHandler: please use get method!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("ObjectResource create success!"))
	fmt.Printf("AddObjectResourceHandler: ObjectResource create success!\n")

}
