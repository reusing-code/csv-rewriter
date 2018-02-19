package main

type IndexComponent struct {
}

func (c *IndexComponent) Render() string {
	html :=
		`
<div class="container">
<h1 class="display-1">csvrewrite</h1>
<p class="lead">Rewrite .comdirect csv to YNAB4 compatible csv</p>
<form>
  <div class="form-group">
    <label for="csvinput">Input .csv file</label>
    <input type="file" class="form-control-file" id="csvinput">
  </div>
</form>
<div id="rewrite-output-container" hidden>
<h3>Warnings</h3>
<div class="alert alert-warning" role="alert">
<pre>
<samp id="rewrite-output-content">
</samp>
</pre>
</div>
</div>
</div>`

	return html

}
