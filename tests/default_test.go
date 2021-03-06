package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"test/db"
	"test/models"
	_ "test/routers"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func SeedDB() {
	db.InitDatabase()
	var jsonStr = []byte(`
		{
		  "ProductID": 1,
		  "Quatity": 10
		}`)
	var jsonStr2 = []byte(`
		{
		  "ProductID": 2,
		  "Quatity": 5
		}`)
	var jsonStr3 = []byte(`
		{
		  "ProductID": 3,
		  "Quatity": 10
		}`)
	r, _ := http.NewRequest("POST", "/v1/product/AddProduct", bytes.NewBuffer(jsonStr))
	r2, _ := http.NewRequest("POST", "/v1/product/AddProduct", bytes.NewBuffer(jsonStr2))
	r3, _ := http.NewRequest("POST", "/v1/product/AddProduct", bytes.NewBuffer(jsonStr3))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	beego.BeeApp.Handlers.ServeHTTP(w, r2)
	beego.BeeApp.Handlers.ServeHTTP(w, r3)
}
func TestPostPurchasesCase1(t *testing.T) {
	SeedDB()
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 1
		}
	      ]`)
	r, _ := http.NewRequest("POST", "/v1/product/purchases", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			product1, _ := models.FindProduct(1)
			product2, _ := models.FindProduct(2)
			So(w.Code, ShouldEqual, 200)
			So(product1.Quatity, ShouldEqual, 8)
			So(product2.Quatity, ShouldEqual, 4)
		})
	})
}

func TestPostPurchasesCase2(t *testing.T) {
	SeedDB()
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 6
		}
	      ]`)
	r, _ := http.NewRequest("POST", "/v1/product/purchases", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 422", func() {
			So(w.Code, ShouldEqual, 422)
			product1, _ := models.FindProduct(1)
			product2, _ := models.FindProduct(2)
			So(product1.Quatity, ShouldEqual, 10)
			So(product2.Quatity, ShouldEqual, 5)
		})
	})
}
func CreatRequest(jsonStr []byte, requests chan httptest.ResponseRecorder) {
	r, _ := http.NewRequest("POST", "/v1/product/purchases", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	requests <- *w
}
func TestPostPurchasesCase3(t *testing.T) {
	SeedDB()
	requests := make(chan httptest.ResponseRecorder, 2)
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 1
		}
	      ]`)

	var jsonStr2 = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 1
		},{
		  "ProductID": 2,
		  "Quatity": 2
		}
	      ]`)
	go CreatRequest(jsonStr, requests)
	go CreatRequest(jsonStr2, requests)
	w1 := <-requests
	w2 := <-requests
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w1.Code, ShouldEqual, 200)
			So(w2.Code, ShouldEqual, 200)
			product1, _ := models.FindProduct(1)
			product2, _ := models.FindProduct(2)
			So(product1.Quatity, ShouldEqual, 7)
			So(product2.Quatity, ShouldEqual, 2)
		})
	})
}

func TestPostPurchasesCase4(t *testing.T) {
	SeedDB()
	requests := make(chan httptest.ResponseRecorder, 2)
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 1
		}
	      ]`)

	var jsonStr2 = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 1
		},{
		  "ProductID": 2,
		  "Quatity": 5
		}
	      ]`)
	go CreatRequest(jsonStr, requests)
	go CreatRequest(jsonStr2, requests)
	w1 := <-requests
	w2 := <-requests
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w1.Code, ShouldEqual, 200)
			So(w2.Code, ShouldEqual, 200)
		})
	})
}

func TestPostPurchasesCase5(t *testing.T) {
	SeedDB()
	requests := make(chan httptest.ResponseRecorder, 2)
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 1
		}
	      ]`)

	var jsonStr2 = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 1
		}
	      ]`)
	go CreatRequest(jsonStr, requests)
	go CreatRequest(jsonStr2, requests)
	w1 := <-requests
	w2 := <-requests
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w1.Code, ShouldEqual, 200)
			So(w2.Code, ShouldEqual, 200)
		})
	})
}

func TestPostPurchasesCase6(t *testing.T) {
	SeedDB()
	requests := make(chan httptest.ResponseRecorder, 2)
	var jsonStr = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 2
		},{
		  "ProductID": 2,
		  "Quatity": 1
		}
	      ]`)

	var jsonStr2 = []byte(`[
		{
		  "ProductID": 1,
		  "Quatity": 1
		},
		{
		  "ProductID": 3,
		  "Quatity": 2
		}
	      ]`)
	go CreatRequest(jsonStr, requests)
	go CreatRequest(jsonStr2, requests)
	w1 := <-requests
	w2 := <-requests
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w1.Code, ShouldEqual, 200)
			So(w2.Code, ShouldEqual, 200)
		})
	})
}
