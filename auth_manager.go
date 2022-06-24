package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type User struct {
	id                  int
	user, pass, project string
	created             time.Time
}

type Response struct {
	User, Project, msg, Err string
	Msg                     []string
}

var projectVal = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func validateUser(username string, password string, project string) error {
	if len(username) == 0 {
		return errors.New("User name must not be empty!")
	}

	if !projectVal.MatchString(project) {
		return errors.New("Project name must not be empty! Only Letters, numbers and '_' are allowd!")
	}

	if len(password) < 5 {
		return errors.New("The password must contain at least 5 characters!")
	}

	return nil
}

func AddUserProject(w http.ResponseWriter, r *http.Request, user User) {
	fmt.Println("method:", r.Method, user.user) //get request method
	if user.user != "admin" {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var res Response
	if r.Method == "POST" {
		if checkErr(r.ParseForm()) {
			// logic part of log in
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			project := r.Form.Get("project")
			res.User = username
			res.Project = project
			if err := validateUser(username, password, project); err == nil {
				sha256_pass := sha256.Sum256([]byte(password))

				if addUser(username, sha256_pass, project) {
					checkErr(os.MkdirAll(path.Join(cfg.Root_dir, "projects", project), 0761))
					InfoLogger.Println("New User Created: ", username)
					res.msg = fmt.Sprintf("Added new user! ELN file watcher CMD command:\n"+
						"efw.exe -duration <integer> -src <folder> -dst %s%s/%s -user %s -pass %s [-zip]", cfg.Host, cfg.Prefix_url, project, username, password)
				} else {
					res.Err = fmt.Sprint("ERROR no new User: ", username)
				}
			} else {
				res.Err = fmt.Sprint("ERROR: ", err)
			}
		}

	} // write data to response
	if len(res.msg) > 0 {
		res.Msg = strings.Split(res.msg, "\n")
	}
	t, _ := template.ParseFiles("views/new_user.gtpl")
	err := t.Execute(w, res)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

func basicAuth(next func(w http.ResponseWriter, r *http.Request, user User)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {
			user, ok := getUser(username)
			if ok {
				// Calculate SHA-256 hashes for the provided and expected
				// usernames and passwords.
				usernameHash := sha256.Sum256([]byte(username))
				passwordHash := sha256.Sum256([]byte(password))
				expectedUsernameHash := sha256.Sum256([]byte(user.user))
				expectedPasswordHash := []byte(user.pass)

				// Use the subtle.ConstantTimeCompare() function to check if
				// the provided username and password hashes equal the
				// expected username and password hashes. ConstantTimeCompare
				// will return 1 if the values are equal, or 0 otherwise.
				// Importantly, we should to do the work to evaluate both the
				// username and password before checking the return values to
				// avoid leaking information.
				usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
				passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

				// If the username and password are correct, then call
				// the next handler in the chain. Make sure to return
				// afterwards, so that none of the code below is run.
				if usernameMatch && passwordMatch {
					next(w, r, user)
					return
				}
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
