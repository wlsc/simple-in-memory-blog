package model

/**
 *	Holds one particular blog entry
 */
type Post struct {
	Id      string
	Title   string
	Content string
}

/**
 *	Creates new instance of the post structure
 */
func Create(id, title, content string) *Post {
	return &Post{id, title, content}
}