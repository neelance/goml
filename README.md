# goml
experiment on syntactic sugar for putting markup into Go code

## Example

```go
package templates

import (
	"fmt"
)

func RenderPage(e element, data *Data) {
	<div (.container, .large)> {
		if data.ShowHelloWorld {
			<h1> { % "Hello World!" }
		}

		for _, item := range data.Items {
			<div> {
				<img (src: item.ImgSrc)>
				% item.Text
			}
		}
	}
}
```

turns into

```go
package templates

import (
	"fmt"
)

func RenderPage(e element, data *Data) {
	e.AppendElement("div", attributes{".container": true, ".large": true}, func(e element) {
		if data.ShowHelloWorld {
			e.AppendElement("h1", nil, func(e element) { e.AppendTextNode("Hello World!") })
		}

		for _, item := range data.Items {
			e.AppendElement("div", nil, func(e element) {
				e.AppendElement("img", attributes{"src": item.ImgSrc}, nil)
				e.AppendTextNode(item.Text)
			})
		}
	})
}
```