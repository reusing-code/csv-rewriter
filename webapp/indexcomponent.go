package main

type IndexComponent struct {
}

func (c *IndexComponent) Render() string {
	html := `<div class="container-fluid">
    		   <div id="app"></div>
             </div>`

	return html

}
