/*
* @author Colton J. McCurdy
*	GitHub: mccurdyc
* Email:  mccurdyc22@gmail.com
* Date: 2018-01-22
 */

package handlers

import (
	"io"
	"net/http"
)

// Hello is a handler function that will print "hello" followed by a new line.
func Hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello\n")
}
