/*
 * Copyright 2015 Xuyuan Pang
 * Author: Xuyuan Pang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net/http"

	"github.com/Xuyuanp/hador"
	"github.com/Xuyuanp/hador/swagger"
	"github.com/hador-contrib/cors"
)

var nextID = 10001

// User struct
type User struct {
	ID       int    `json:"id"`
	NickName string `json:"nick_name"`
	Password string `json:"password"`
}

// UserList struct
type UserList struct {
	UserCount int    `json:"user_count"`
	Users     []User `json:"users"`
}

var fakeStore = map[int]*User{
	10000: &User{ID: 10000, NickName: "Jack", Password: "foobar"},
}

func main() {
	h := hador.Default()

	h.AddFilters(
		// cors support
		cors.Allow(&cors.CORSOptions{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "DELETE", "PUT"},
		}),
	)

	h.Group(`/v1`, func(v1 hador.Router) {
		v1.Group(`/users`, func(root hador.Router) {

			// GET /v1/users
			root.Get(`/`, hador.HandlerFunc(getUserList)).
				DocOperation().
				DocSumDesc("get user list", "").
				DocResponseModel("200", "user list", UserList{})

			// POST /v1/users
			root.Post(`/`, hador.HandlerFunc(newUser)).
				DocOperation().
				DocSumDesc("new user", "").
				DocParameterBody("user", "user info", User{}, true).
				DocResponseModel("200", "user info", User{})

			root.Group(`/{user-id:\d+}`, func(userRouter hador.Router) {

				// GET /v1/users/{user-id}
				userRouter.Get(`/`, hador.HandlerFunc(getUser)).
					DocOperation().
					DocSumDesc("get user info", "").
					DocParameterPath("user-id", "integer", "user id", true).
					DocResponseModel("200", "user info", User{})

				// DELETE /v1/users/{user-id}
				userRouter.Delete(`/`, hador.HandlerFunc(delUser)).
					DocOperation().
					DocSumDesc("delete user info", "").
					DocParameterPath("user-id", "integer", "user id", true).
					DocResponseModel("200", "user info", User{}).
					DocResponseSimple("404", "not found")

				// PUT /v1/users/{user-id}
				userRouter.Put(`/`, hador.HandlerFunc(updateUser)).
					DocOperation().
					DocSumDesc("update user info", "").
					DocParameterPath("user-id", "user id", "integer", true).
					DocParameterBody("user", "user info", User{}, true).
					DocResponseModel("200", "user info", User{}).
					DocResponseSimple("404", "not found")

			}, UIDFilter())
		}, errorFilter())
	})

	// swagger support
	// open http://127.0.0.1:9090/apidocs in your broswer
	// and enter http://127.0.0.1:9090/apidocs.json in the api input field
	h.SwaggerDocument().
		DocInfo("User Manager", "user CRUD", "v1", "http://your.term.of.service.addr").
		DocHost("127.0.0.1:9090")

	h.Swagger(swagger.Config{
		// your swagger-ui file path
		UIFilePath: "/path/to/your/swagger-ui/dist",

		// swagger json api
		APIPath: "/apidocs.json",

		// swagger-ui web location
		UIPrefix: "/apidocs",
	})

	h.Run(":9090")
}

// UIDFilter resolved user-id param
func UIDFilter() hador.FilterFunc {
	return func(ctx *hador.Context, next hador.Handler) {
		uid, err := ctx.Params().GetInt("user-id")
		if err != nil {
			ctx.OnError(http.StatusBadRequest, err)
			return
		}
		ctx.Set("user-id", uid)
		defer ctx.Delete("user-id")

		next.Serve(ctx)
	}
}

// errorFilter handle error
func errorFilter() hador.FilterFunc {
	return func(ctx *hador.Context, next hador.Handler) {
		// set error message as json format
		ctx.Err4XXHandler = func(status int, args ...interface{}) {
			text := http.StatusText(status)
			ctx.RenderJSON(text, status)
		}
		next.Serve(ctx)
	}
}

func getUserList(ctx *hador.Context) {
	if len(fakeStore) == 0 {
		ctx.OnError(http.StatusNotFound)
		return
	}
	users := make([]User, len(fakeStore))
	i := 0
	for _, u := range fakeStore {
		users[i] = *u
		i++
	}
	result := UserList{
		UserCount: len(users),
		Users:     users,
	}
	ctx.RenderJSON(result, http.StatusOK)
}

func getUser(ctx *hador.Context) {
	uid, _ := ctx.Get("user-id").(int)
	user, ok := fakeStore[uid]
	if !ok {
		ctx.OnError(http.StatusNotFound)
		return
	}
	ctx.RenderJSON(user)
}

func delUser(ctx *hador.Context) {
	uid, _ := ctx.Get("user-id").(int)
	user, ok := fakeStore[uid]
	if !ok {
		ctx.OnError(http.StatusNotFound)
		return
	}
	delete(fakeStore, uid)
	ctx.RenderJSON(user)
}

func newUser(ctx *hador.Context) {
	user := User{}
	if err := ctx.ResolveJSON(&user); err != nil {
		ctx.OnError(http.StatusBadRequest, err)
		return
	}
	user.ID = nextID
	fakeStore[nextID] = &user
	nextID++

	ctx.RenderJSON(user, http.StatusCreated)
}

func updateUser(ctx *hador.Context) {
	uid, _ := ctx.Get("user-id").(int)
	user, ok := fakeStore[uid]
	if !ok {
		ctx.OnError(http.StatusNotFound)
		return
	}
	newUser := User{}
	if err := ctx.ResolveJSON(&newUser); err != nil {
		ctx.OnError(http.StatusBadRequest, err)
		return
	}
	user.Password = newUser.Password
	user.NickName = newUser.NickName

	ctx.RenderJSON(user)
}
