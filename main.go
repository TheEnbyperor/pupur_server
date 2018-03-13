package main

import (
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qor/admin"
	"github.com/qor/media"
	"github.com/qor/media/media_library"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/rs/cors"
	"gopkg.in/hlandau/passlib.v1"
	"log"
	"net/http"
	"os"
)

const addr = ":8080"

var db *gorm.DB

type User struct {
	gorm.Model
	Login string
	Pass  string
}

func (u *User) DisplayName() string {
	return u.Login
}

type Recipe struct {
	gorm.Model
	Name        string                 `json:"name"`
	Author      string                 `json:"author"`
	Description string                 `json:"description"`
	Image       media_library.MediaBox `json:"image"`
	Icon        media_library.MediaBox `json:"icon"`
	Ingredients []*Ingredient          `json:"ingredients"`
	Steps       []*Step                `json:"steps"`
}

type Ingredient struct {
	gorm.Model
	RecipeId    uint
	Description string `json:"description"`
	Amount      string `json:"amount"`
}

type Step struct {
	gorm.Model
	RecipeId    uint
	Image       media_library.MediaBox `json:"image"`
	Description string                 `json:"description"`
}

var sessionManager = scs.NewCookieManager("i2D3qgZyClwuNHRgEzc4uiUKR1mFR9fo")

var ingredientType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Ingredient",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				ingredient, isOk := params.Source.(*Ingredient)
				if isOk {
					return ingredient.Model.ID, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"amount": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var stepType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Step",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				step, isOk := params.Source.(*Step)
				if isOk {
					return step.Model.ID, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"image": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				step, isOk := params.Source.(*Step)
				if isOk {
					return step.Image.URL(), nil
				}
				return nil, nil
			},
		},
	},
})

var recipeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Recipe",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipe, isOk := params.Source.(*Recipe)
				if isOk {
					return recipe.Model.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"author": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"image": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipe, isOk := params.Source.(*Recipe)
				if isOk {
					return recipe.Image.URL(), nil
				}
				return nil, nil
			},
		},
		"icon": &graphql.Field{
			Type: graphql.String,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipe, isOk := params.Source.(*Recipe)
				if isOk {
					return recipe.Icon.URL(), nil
				}
				return nil, nil
			},
		},
		"ingredients": &graphql.Field{
			Type: graphql.NewList(ingredientType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipe, isOk := params.Source.(*Recipe)
				if isOk {
					ingredients := make([]*Ingredient, 0)

					err := db.Model(&recipe).Related(&ingredients).Error
					if err != nil {
						return nil, err
					}

					return ingredients, nil
				}
				return nil, nil
			},
		},
		"steps": &graphql.Field{
			Type: graphql.NewList(stepType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipe, isOk := params.Source.(*Recipe)
				if isOk {
					steps := make([]*Step, 0)

					err := db.Model(&recipe).Related(&steps).Error
					if err != nil {
						return nil, err
					}

					return steps, nil
				}
				return nil, nil
			},
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"recipes": &graphql.Field{
			Type: graphql.NewList(recipeType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				recipes := make([]*Recipe, 0)

				err := db.Find(&recipes).Error
				if err != nil {
					return nil, err
				}

				return recipes, nil
			},
		},
	},
})

type AdminAuth struct{}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	session := sessionManager.Load(c.Request)
	userId, err := session.GetInt("userID")
	if err == nil {
		user := &User{}
		err := db.First(user, userId).Error
		if err == nil {
			return user
		}
	}
	return nil
}

const loginPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="login">
    <label for="password">Password</label>
    <input type="password" id="password" name="pass">
    <button type="submit">Login</button>
</form>
`

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session := sessionManager.Load(r)
		r.ParseForm()
		name := r.FormValue("login")
		pass := r.FormValue("pass")

		redirectTarget := "/login"

		if name != "" && pass != "" {
			user := &User{
				Login: name,
			}
			err := db.First(user, user).Error
			if err == nil {
				newHash, err := passlib.Verify(pass, user.Pass)
				if err == nil {
					if newHash != "" {
						user.Pass = newHash
						db.Save(user)
					}
					redirectTarget = "/"
					session.PutInt(w, "userID", int(user.Model.ID))
				} else {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		}
		http.Redirect(w, r, redirectTarget, 302)
	} else {
		fmt.Fprint(w, loginPage)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.Load(r)
	session.Destroy(w)

	http.Redirect(w, r, "/", 302)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.Load(r)
	userId, err := session.GetInt("userID")
	if err == nil {
		user := &User{}
		err := db.First(user, userId).Error
		if err == nil {
			http.Redirect(w, r, "/admin", 302)
			return
		}
	}
	http.Redirect(w, r, "/login", 302)
}

func main() {
	ingredientType.AddFieldConfig("recipe", &graphql.Field{
		Type: recipeType,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ingredient, isOk := params.Source.(*Ingredient)
			if isOk {
				recipe := &Recipe{}

				err := db.First(recipe, ingredient.RecipeId).Error
				if err != nil {
					return nil, err
				}

				return recipe, nil
			}
			return nil, nil
		},
	})
	stepType.AddFieldConfig("recipe", &graphql.Field{
		Type: recipeType,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			step, isOk := params.Source.(*Step)
			if isOk {
				recipe := &Recipe{}

				err := db.First(recipe, step.RecipeId).Error
				if err != nil {
					return nil, err
				}

				return recipe, nil
			}
			return nil, nil
		},
	})

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbName := os.Getenv("DB_NAME")

	log.Println("Connecting to database...")

	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", dbUser, dbPass, dbHost, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	media.RegisterCallbacks(db)
	err = db.AutoMigrate(&media_library.MediaLibrary{}, &Recipe{}, &Ingredient{}, &Step{}, &User{}).Error
	if err != nil {
		panic(err)
	}

	Admin := admin.New(&admin.AdminConfig{
		SiteName: "Pupur",
		DB:       db,
		Auth:     AdminAuth{},
	})

	Admin.AddResource(&media_library.MediaLibrary{})
	recipe := Admin.AddResource(&Recipe{})
	user := Admin.AddResource(&User{})
	user.IndexAttrs("ID", "Login")
	user.EditAttrs("Login", "Password")
	user.Meta(&admin.Meta{
		Name: "Password",
		Type:   "password",
		Valuer: func(interface{}, *qor.Context) interface{} { return "" },
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			if newPassword := utils.ToString(metaValue.Value); newPassword != "" {
				hash, _ := passlib.Hash(newPassword)
				user := record.(*User)
				user.Pass = hash
				db.Save(user)
			}
		},
	})

	recipe.GetMeta("Steps").Resource.Meta(&admin.Meta{
		Name: "Description",
		Type: "text",
	})

	adminMux := http.NewServeMux()
	Admin.MountTo("/admin", adminMux)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	if err != nil {
		panic(err)
	}

		h := handler.New(&handler.Config{
			Schema:   &schema,
			Pretty:   true,
			GraphiQL: true,
		})
		corsH := cors.Default().Handler(h)
		r := mux.NewRouter()

		r.Path("/").HandlerFunc(indexHandler)

		r.Path("/login").Methods("GET", "POST").HandlerFunc(loginHandler)
		r.Path("/logout").HandlerFunc(logoutHandler)
		r.PathPrefix("/admin").Handler(adminMux)
		r.Handle("/graphql", corsH)
		r.PathPrefix("/media").Handler(http.StripPrefix("/media", utils.FileServer(http.Dir("media"))))
		r.PathPrefix("/").Handler(utils.FileServer(http.Dir("public")))

		log.Printf("Listening on %s\n", addr)
		http.ListenAndServe(addr, r)
}
