package controller

import "net/http"

// Analytics handles a web POST "/analytics" request. It responds with a view (HTML partial) containing all
// relevant analytics related to the current bank.
func (c *ComplianceCheckController) Analytics(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./app/web/static/views/analytics.html")
}
