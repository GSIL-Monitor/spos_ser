// User
package route

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
	"spos_auth/datas"
	"spos_auth/enforcer"
	"spos_auth/umisc"
	"strconv"
	"strings"
	"time"
)

const (
	SecretKey = "unitone_spos"
)

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("RegisterHandler: url=%s\n", r.URL.Path)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {
		if len(r.FormValue("name")) > 0 {
			fmt.Printf("RegisterHandler: name=%s;password=%s;phone=%s;mail=%s\n", r.FormValue("name"), r.FormValue("password"), r.FormValue("phone"), r.FormValue("mail"))

			if len(r.FormValue("name")) == 0 || len(r.FormValue("password")) < 6 || len(r.FormValue("password")) > 20 || len(r.FormValue("mail")) == 0 {
				fmt.Printf("RegisterHandler: Register info error\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("RegisterHandler: Register info error!"))
				return
			}

			if datas.IsExistUser(strings.Replace(r.FormValue("name"), " ", "", -1)) {
				fmt.Printf("RegisterHandler: ERROR,User is Exist!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("RegisterHandler: ERROR,User is Exist!"))
				return
			}

			fmt.Printf("RegisterHandler: new user codeing...\n")

			Encryptor := umisc.BcryptEncryptorNew(&umisc.Config{})
			EncryptedPassword, _ := Encryptor.Digest(r.FormValue("password"))

			timestamp := time.Now().Unix()
			fmt.Println(timestamp)

			basic_info := datas.Basic{UID: "U" + strconv.FormatInt(timestamp, 10), UserName: r.FormValue("name"), EncryptedPassword: EncryptedPassword, PhoneNum: r.FormValue("phone"), MailAddr: r.FormValue("mail")}
			user := &datas.UserInfo{Basic: basic_info}

			ret := user.AddUser()

			if ret != nil {
				fmt.Printf("RegisterHandler: Organization connext error, db record write error!\n")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("RegisterHandler: Organization connext error, db record write error!"))
				return
			}

			token := jwt.New(jwt.SigningMethodHS256)

			claims := basic_info.ToClaims()
			claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(24)).Unix()
			claims.IssuedAt = time.Now().Unix()

			token.Claims = claims

			//	tokenString, err := token.SignedString(signKey)
			tokenString, err := token.SignedString([]byte(SecretKey))

			fmt.Printf("RegisterHandler:tokenString=%s; err=%v\n", tokenString, err)

			if err != nil {
				fmt.Fprintln(w, "RegisterHandler:Error while signing the token")
				return
				//		fatal(err)
			}

			SendMailToActivate(r.FormValue("mail"), r.FormValue("name"), tokenString)

		} else {
			fmt.Printf("RegisterHandler: connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("RegisterHandler: connext error!"))
			return
		}
	} else {
		fmt.Printf("RegisterHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("RegisterHandler: please use get method!"))
		return
	}

	fmt.Printf("RegisterHandler: Account create success!\n")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, string("Account create success!"))

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//	var user UserCredentials

	fmt.Printf("LoginHandler: url=%s...\n", r.URL.Path)

	w.Header().Set("WWW-Authenticate", `Basic realm="Ambassador Realm"`)

	username, password, authOK := r.BasicAuth()

	fmt.Printf("LoginHandler: username=%s; password=%s,authOK=%v\n", username, password, authOK)

	if authOK == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "LoginHandler: Not authorized", 401)
		fmt.Printf("LoginHandler: Not authorized\n")
		return
	}

	var user datas.UserInfo
	user.GetUserByName(username)
	fmt.Printf("LoginHandler: password=%s\n", user.EncryptedPassword)

	l := &datas.SignLogs{}

	if len(user.EncryptedPassword) > 0 {
		Encryptor := umisc.BcryptEncryptorNew(&umisc.Config{})
		if err := Encryptor.Compare(user.EncryptedPassword, password); err == nil {
			fmt.Println("LoginHandler:Password ok")
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "LoginHandler:Password error", 401)
			fmt.Printf("LoginHandler:Password error\n")
			l.Log = "Password error"
			l.SignLogsUpdateToDB(user)
			return
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "account error, Not login authorized", 401)
		fmt.Printf("spos login account error\n")
		return
	}

	if !datas.IsActivateUser(user.UID) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "LoginHandler:Not mail confirmed! Please login your mail activate!", 401)
		fmt.Printf("LoginHandler:Not mail confirmed! Please login your mail activate!\n")
		l.Log = "Not mail confirmed!"
		l.SignLogsUpdateToDB(user)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	fmt.Printf("LoginHandler: user.UID=%s\n", user.UID)

	basic_info := datas.Basic{UID: user.UID}

	claims := basic_info.ToClaims()
	claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims.IssuedAt = time.Now().Unix()

	token.Claims = claims

	//	tokenString, err := token.SignedString(signKey)
	tokenString, err := token.SignedString([]byte(SecretKey))

	fmt.Printf("LoginHandler:tokenString=%s; err=%v\n", tokenString, err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "LoginHandler:Error while signing the token")
		l.Log = "Error while signing the token!"
		l.SignLogsUpdateToDB(user)
		return
		//		fatal(err)
	}

	response := Token{tokenString}

	fmt.Printf("LoginHandler:response=%v\n", response)
	JsonResponse(response, w)

	//cb_en.AddPolicy(user.UID, "/extauth/spos.Spos/SayHello", "POST") //tmp patli 20180907
	fmt.Printf("LoginHandler: User Login success!\n")
	l.Log = "User Login success!"
	l.SignLogsUpdateToDB(user)

}

func RegisteredHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("RegisteredHandler: url=%s\n", r.URL.Path)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.ParseForm()
	if r.Method == "GET" {

		if len(r.FormValue("Activate")) > 0 {
			fmt.Printf("RegisteredHandler: Activate=%s\n", r.FormValue("Activate"))
			j := umisc.NewJWT()
			// parseToken
			claims, err := j.ParseToken(r.FormValue("Activate"))
			fmt.Printf("RegisteredHandler: claims=%v; err=%v; uid=%s\n", claims, err, claims.UserID)
			if claims != nil {
				var user datas.UserInfo
				user.UID = claims.UserID
				if datas.IsActivateUser(user.UID) {
					fmt.Printf("RegisteredHandler: User Activated!")
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, string("RegisteredHandler: User Activated!"))
					return
				}

				ret := user.UpdateUserConfirmedAt()

				if ret != nil {
					fmt.Printf("RegisteredHandler: Organization connext error, db record update error!\n")
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, string("RegisteredHandler: Organization connext error, db record update error!"))
					return
				}
			} else {
				fmt.Printf("RegisteredHandler: jwt connext invaild!")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, string("RegisteredHandler: jwt connext invaild!"))
				return
			}

		} else {
			fmt.Printf("RegisteredHandler: Activate connext error\n")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, string("RegisteredHandler: Activate connext error!"))
			return
		}
	} else {
		fmt.Printf("RegisteredHandler: please use get method\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, string("RegisteredHandler: please use get method!"))
		return
	}

	fmt.Printf("RegisteredHandler: Account activate success!\n")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string("Account activate success!"))

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	fmt.Printf("JsonResponse: json create...\n")
	json, err := json.Marshal(response)

	fmt.Printf("JsonResponse: json=%v; err=%v\n", json, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("JsonResponse: send token json...\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, string(json))
}

func SendMailToActivate(to_addr string, username string, token string) {
	user := "spos@unitone.com.cn"
	password := "Spos@2017"
	host := "smtp.qiye.163.com:25"
	to := to_addr

	subject := "Spos account activate"

	body := `
 <html>
 <body>
 <h3>` + username + ", Please click the following URLs within 24 hours to activate, Thank you." + `</h3>
 <h3>` + "https://192.168.20.101:6443/api/v1/namespaces/default/services/ambassador/proxy/registered?Activate=" + token + `</h3>
 </body>
 </html>
 `
	fmt.Println("SendMailToActivate: Send email...")
	err := umisc.SendMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("SendMailToActivate: Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("SendMailToActivate: Send mail success!")
	}
}

func JwtTokenAuth(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("JwtTokenAuth: url=%s; method=%s...\n", r.URL.Path, r.Method)

	if strings.Contains(r.URL.Path, "spos.Lang") || strings.Contains(r.URL.Path, "/qotm/quote") {

		if strings.Contains(r.URL.Path, "spos.Lang/Login") {
			fmt.Printf("JwtTokenAuth: Login permit!\n")
			return
		}

		fmt.Printf("JwtTokenAuth: url=%s; need token!\n", r.URL.Path)

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(SecretKey), nil
			})

		fmt.Printf("JwtTokenAuth: token=%v\n", token)

		if err == nil {
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)

				uid := claims["userid"].(string)
				fmt.Printf("JwtTokenAuth: claims userid=%s; exp=%v; iat=%v\n", uid, claims["exp"], claims["iat"])
				fmt.Printf("JwtTokenAuth: Token Vaild\n")

				if enforcer.Enforce(uid, r.URL.Path, r.Method) == true {
					// permit alice to read data1
					fmt.Printf("JwtTokenAuth: permit the request\n")
				} else {
					fmt.Printf("JwtTokenAuth: deny the request\n")
					// deny the request, show an error

					w.WriteHeader(http.StatusUnauthorized)
					fmt.Printf("JwtTokenAuth: Unauthorized access to resources\n")
					fmt.Fprint(w, "Unauthorized access to resources")
				}

			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Printf("JwtTokenAuth: Token is not valid\n")
				fmt.Fprint(w, "Token is not valid")
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Printf("JwtTokenAuth: Unauthorized access to this resource\n")
			fmt.Fprint(w, "Unauthorized access to this resource")
		}

	} else {
		fmt.Printf("JwtTokenAuth: url=%s; not need token!\n", r.URL.Path)
	}

}

func UpdatePWHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("UpdatePWHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

}

func ForgetPWHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ForgetPWHandler: url=%s; method=%s...\n", r.URL.Path, r.Method)

}
