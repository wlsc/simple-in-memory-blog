package main

/**
	Simple in-memory test blog playground
 */

import (
	"fmt"
	"net/http"
	"crypto/rand"
	"strconv"
	"./model/post"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

// *********************************
// ************ CONFIG *************
var assetFolder string = "asset"
var port int = 8080
var CHARSET string = "UTF-8"
// *********************************

var posts map[string]*model.Post

/*
 *	Generates random string of given length
 */
func GenerateId(length int) string {

	bytes := make([]byte, length)
	rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)
}

/**
 *	Index page and list of available in-memory blog entries
 */
func indexHandler(render render.Render) {
	render.HTML(200, "page/post-list", posts)
}

/**
 *	Loads create view for new post
 */
func newHandler(render render.Render) {
	render.HTML(200, "page/post-new", nil)
}

/**
 *	Loads update view for a given post
 */
func editHandler(render render.Render, params martini.Params) {

	id := params["id"]
	post, found := posts[id]

	if !found {
		render.Redirect("/")
		return
	}

	render.HTML(200, "page/post-new", post)
}

/**
 *	Creates new post within in-memory array or
 *	updates an old post
 */
func saveHandler(render render.Render, request *http.Request) {

	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *model.Post

	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id := GenerateId(8)
		posts[id] = model.Create(id, title, content)
	}

	render.Redirect("/")
}

/**
 *	Deletes post from in-memory array
 */
func deleteHandler(render render.Render, params martini.Params) {

	id := params["id"]

	if id == "" {
		render.Redirect("/")
		return
	}

	delete(posts, id)

	render.Redirect("/")
}

/**
 *	Entry point here
 */
func main() {

	posts = make(map[string]*model.Post, 0)

	m := martini.Classic()

	// RENDER SETUP
	m.Use(render.Renderer(render.Options{
		Directory:"template",
		Layout:"layout",
		Extensions:[]string{".htm", ".html", ".tmpl"},
		Charset:CHARSET,
		IndentJSON:true,
	}))

	// ROUTES SETUP
	// static
	staticOptions := martini.StaticOptions{
		Prefix : assetFolder}
	m.Use(martini.Static(assetFolder, staticOptions))
	// dynamic
	m.Get("/", indexHandler)
	m.Get("/new", newHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/save", saveHandler)

	m.RunOnAddr(":" + strconv.Itoa(port))
}