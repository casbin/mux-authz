package authz

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, handler http.HandlerFunc, c *CasbinAuthorizer, user string, path string, method string, code int) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user,"111")

	r := mux.NewRouter()
	r.HandleFunc("/{url:[A-Za-z0-9\\/]+}", handler)
	r.Use(c.Middleware)


	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr,req)


	fmt.Println(req.URL.String())
	// Check the status code is what we expect.
	if rr.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, rr.Code, code)
	}

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func TestHealthCheckHandler(t *testing.T) {

	e, _ := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")

	c:= &CasbinAuthorizer{Enforcer: e}

	handler:= http.HandlerFunc(HealthCheckHandler)

	testRequest(t, handler,	c, "alice", "/dataset1/resource1", "GET", 200)
	testRequest(t, handler,	c, "alice", "/dataset1/resource1", "POST", 200)
	testRequest(t, handler,	c, "alice", "/dataset1/resource2", "GET", 200)
	testRequest(t, handler,	c, "alice", "/dataset1/resource2", "POST", 403)

	testRequest(t, handler, c, "bob", "/dataset2/resource1", "GET", 200)
	testRequest(t, handler, c, "bob", "/dataset2/resource1", "POST", 200)
	testRequest(t, handler, c, "bob", "/dataset2/resource1", "DELETE", 200)
	testRequest(t, handler, c, "bob", "/dataset2/resource2", "GET", 200)
	testRequest(t, handler, c, "bob", "/dataset2/resource2", "POST", 403)
	testRequest(t, handler, c, "bob", "/dataset2/resource2", "DELETE", 403)

	testRequest(t, handler, c, "bob", "/dataset2/folder1/item1", "GET", 403)
	testRequest(t, handler, c, "bob", "/dataset2/folder1/item1", "POST", 200)
	testRequest(t, handler, c, "bob", "/dataset2/folder1/item1", "DELETE", 403)
	testRequest(t, handler, c, "bob", "/dataset2/folder1/item2", "GET", 403)
	testRequest(t, handler, c, "bob", "/dataset2/folder1/item2", "POST", 200)
	testRequest(t, handler, c, "bob", "/dataset2/folder1/item2", "DELETE", 403)

	// cathy can access all /dataset1/* resources via all methods because it has the dataset1_admin role.
	testRequest(t, handler, c, "cathy", "/dataset1/item", "GET", 200)
	testRequest(t, handler, c, "cathy", "/dataset1/item", "POST", 200)
	testRequest(t, handler, c, "cathy", "/dataset1/item", "DELETE", 200)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "GET", 403)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "POST", 403)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "DELETE", 403)

	// delete all roles on user cathy, so cathy cannot access any resources now.
	c.Enforcer.DeleteRolesForUser("cathy")

	testRequest(t, handler, c, "cathy", "/dataset1/item", "GET", 403)
	testRequest(t, handler, c, "cathy", "/dataset1/item", "POST", 403)
	testRequest(t, handler, c, "cathy", "/dataset1/item", "DELETE", 403)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "GET", 403)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "POST", 403)
	testRequest(t, handler, c, "cathy", "/dataset2/item", "DELETE", 403)

}