// main.go
package main

import (
	"fmt"
	"net/http"
	"spos_auth/datas"
	"spos_auth/route"
	//	"github.com/dgrijalva/jwt-go"
	//	"spos_auth/umisc"
	//	"strconv"
	//	"strings"
	//	"time"
)

func main() {
	fmt.Printf("Spos auth running ...\n")

	if datas.OpenUnDbAndInitCasbin() != nil {
		fmt.Printf("main: db open error!\n")
		panic("db connection failure")
	}

	//  http.HandleFunc("/extauth/qotm/quote", use(myHandler, basicAuth))
	//	http.HandleFunc("/", use(myHandler, basicAuth))

	//	http.HandleFunc("/", basicAuth)

	http.HandleFunc("/", route.JwtTokenAuth)
	http.HandleFunc("/extauth/login", route.LoginHandler)
	http.HandleFunc("/extauth/register", route.RegisterHandler)
	http.HandleFunc("/extauth/registered", route.RegisteredHandler)
	http.HandleFunc("/extauth/updatepassword", route.UpdatePWHandler)
	http.HandleFunc("/extauth/forgetpassword", route.ForgetPWHandler)

	http.HandleFunc("/extauth/addIndustryCategory", route.AddIndustryCategoryHandler)
	http.HandleFunc("/extauth/addOrganization", route.AddOrganizationHandler)
	http.HandleFunc("/extauth/addRole", route.AddRoleHandler)
	http.HandleFunc("/extauth/addObjectResource", route.AddObjectResourceHandler)
	http.HandleFunc("/extauth/addRoleObject", route.AddRoleObjectHandler)
	http.HandleFunc("/extauth/addUserRole", route.AddUserRoleHandler)
	http.HandleFunc("/extauth/addDepartment", route.AddDepartmentHandler)
	http.HandleFunc("/extauth/addDepartmentOR", route.AddDepartmentORHandler)
	http.HandleFunc("/extauth/addUserOrg", route.AddUserOrgHandler)

	http.ListenAndServe(":3000", nil)

	fmt.Printf("main end\n")

}

//func RegisterHandler(w http.ResponseWriter, r *http.Request) {

//	fmt.Printf("RegisterHandler: url=%s\n", r.URL.Path)
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

//	r.ParseForm()
//	if r.Method == "GET" {
//		if len(r.FormValue("name")) > 0 {
//			fmt.Printf("RegisterHandler: name=%s;password=%s;phone=%s;mail=%s\n", r.FormValue("name"), r.FormValue("password"), r.FormValue("phone"), r.FormValue("mail"))

//			if len(r.FormValue("name")) == 0 || len(r.FormValue("password")) < 6 || len(r.FormValue("password")) > 20 || len(r.FormValue("mail")) == 0 {
//				fmt.Printf("RegisterHandler: Register info error\n")
//				w.WriteHeader(http.StatusBadRequest)
//				fmt.Fprintln(w, string("RegisterHandler: Register info error!"))
//				return
//			}

//			if datas.IsExistUser(strings.Replace(r.FormValue("name"), " ", "", -1)) {
//				fmt.Printf("RegisterHandler: ERROR,User is Exist!\n")
//				w.WriteHeader(http.StatusBadRequest)
//				fmt.Fprintln(w, string("RegisterHandler: ERROR,User is Exist!"))
//				return
//			}

//			fmt.Printf("RegisterHandler: new user codeing...\n")

//			Encryptor := umisc.BcryptEncryptorNew(&umisc.Config{})
//			EncryptedPassword, _ := Encryptor.Digest(r.FormValue("password"))

//			timestamp := time.Now().Unix()
//			fmt.Println(timestamp)

//			basic_info := datas.Basic{UID: strconv.FormatInt(timestamp, 10), UserName: r.FormValue("name"), EncryptedPassword: EncryptedPassword, PhoneNum: r.FormValue("phone"), MailAddr: r.FormValue("mail")}
//			user := &datas.UserInfo{Basic: basic_info}

//			ret := user.AddUser()

//			if ret != nil {
//				fmt.Printf("RegisterHandler: Organization connext error, db record write error!\n")
//				w.WriteHeader(http.StatusBadRequest)
//				fmt.Fprintln(w, string("RegisterHandler: Organization connext error, db record write error!"))
//				return
//			}

//			token := jwt.New(jwt.SigningMethodHS256)

//			claims := basic_info.ToClaims()
//			claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(24)).Unix()
//			claims.IssuedAt = time.Now().Unix()

//			token.Claims = claims

//			//	tokenString, err := token.SignedString(signKey)
//			tokenString, err := token.SignedString([]byte(route.SecretKey))

//			fmt.Printf("RegisterHandler:tokenString=%s; err=%v\n", tokenString, err)

//			if err != nil {
//				fmt.Fprintln(w, "RegisterHandler:Error while signing the token")
//				return
//				//		fatal(err)
//			}

//			route.SendMailToActivate(r.FormValue("mail"), r.FormValue("name"), tokenString)

//		} else {
//			fmt.Printf("RegisterHandler: connext error\n")
//			w.WriteHeader(http.StatusBadRequest)
//			fmt.Fprintln(w, string("RegisterHandler: connext error!"))
//			return
//		}
//	} else {
//		fmt.Printf("RegisterHandler: please use get method\n")
//		w.WriteHeader(http.StatusBadRequest)
//		fmt.Fprintln(w, string("RegisterHandler: please use get method!"))
//		return
//	}

//	w.WriteHeader(http.StatusCreated)
//	fmt.Fprintln(w, string("Account create success!"))
//	fmt.Printf("RegisterHandler: Account create success!\n")
//}
