package nodes

import "bufio"

type Representation struct {
	Content string `bson:"c"`
}

////////////////////////////////////////////////////////////////////////////////
func (r *Representation) Render(w *bufio.Writer) error {
	w.WriteString(r.Content)
	return nil
}
