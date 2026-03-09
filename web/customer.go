// Package user tbd
package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
	router "github.com/rosso0815/rosso0815-go-crud-billing/router"
	services "github.com/rosso0815/rosso0815-go-crud-billing/services"
	ui "github.com/rosso0815/rosso0815-go-crud-billing/web/ui"
)

// New
func NewCustomer(store *services.Store, sessionManager *scs.SessionManager, cfg *config.Config, r *router.Router) *Web {
	b := Web{}
	b.store = store
	b.Path = "customer"
	b.SessionManager = sessionManager
	b.cfg = cfg
	b.AddAction = ui.CrudActionAdd(fmt.Sprintf("%s/ui/%s/form_add", b.cfg.WebPrefix, b.Path))
	cfg.Menus = append(cfg.Menus, config.Menu{Name: "Customer", Path: fmt.Sprintf("ui/%s", b.Path)})
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s", b.cfg.WebPrefix, b.Path), b.customerListAll)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/edit/{id}", b.cfg.WebPrefix, b.Path), b.customer_FormEdit)
	r.RegisterRoute(fmt.Sprintf("PUT %s/ui/%s/{id}", b.cfg.WebPrefix, b.Path), b.customer_Edit)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/form_add", b.cfg.WebPrefix, b.Path), b.customerFormAdd)
	r.RegisterRoute(fmt.Sprintf("POST %s/ui/%s", b.cfg.WebPrefix, b.Path), b.customerAdd)
	r.RegisterRoute(fmt.Sprintf("DELETE %s/ui/%s/{id}", b.cfg.WebPrefix, b.Path), b.customerDelete)
	return &b
}

// ListAll
func (b *Web) customerListAll(w http.ResponseWriter, r *http.Request) {
	b.Update(w, r)
	var header []templ.Component = []templ.Component{
		ui.CrudHeaderSort("CustomerId", b.cfg),
		ui.CrudHeaderSort("FirstName", b.cfg),
		ui.CrudHeaderSort("LastName", b.cfg),
		ui.CrudHeaderSort("Street", b.cfg),
		ui.CrudHeaderSort("Town", b.cfg),
		ui.CrudHeaderSort("Remark", b.cfg),
		ui.CrudHeaderSort("Action", b.cfg),
	}
	rCount, err := b.store.CustomerListBySearchCount(r.Context(), b.Search)
	if err != nil {
		log.Println("ERROR getting customer count:", err)
		rCount = 0
	}
	items, err := b.store.CustomerListBySearch(r.Context(), b.Search, b.PageSize, b.PageCount)
	if err != nil {
		log.Println("ERROR getting customers:", err)
		items = []services.Customer{}
	}
	var components []templ.Component
	for _, item := range items {
		l := customerLine(item, b.cfg)
		components = append(components, l)
	}
	if err := ui.CrudTable(
		ui.CrudList{
			Base:   b.Base,
			Header: header,
			Items:  components,
			Count:  rCount,
		}, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) customerDelete(w http.ResponseWriter, r *http.Request) {
	b.Update(w, r)
	idParam, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := ui.CrudMessageAlert(b.Base, err.Error(), b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	err = b.store.CustomerDelete(r.Context(), idParam)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := ui.CrudMessageAlert(b.Base, err.Error(), b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	b.Base.MessageText = fmt.Sprintf("delete done of Customer-Id %d", idParam)
	if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) customer_FormEdit(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	idParam, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	b.Base.Update(w, r)
	customer, err := b.store.CustomerGetById(r.Context(), idParam)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := ui.CrudMessageAlert(b.Base, err.Error(), b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if err := customer_Form(b.Base, ui.Edit, "edit", customer, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) customer_Edit(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var c services.Customer
	err := decoder.Decode(&c)
	if err != nil {
		log.Println("ERROR decoding JSON:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		b.MessageText = err.Error()
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	err = b.store.CustomerUpdate(r.Context(), c)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	b.Base.Update(w, r)
	b.Base.MessageText = fmt.Sprintf("Customer %d updated", c.CustomerID)
	b.customerListAll(w, r)
}

func (b Web) customerFormAdd(w http.ResponseWriter, r *http.Request) {
	b.Base.Update(w, r)
	if err := customer_Form(b.Base, ui.Add, "add", services.Customer{}, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

// Add Customer
func (b Web) customerAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("@@@ customer_add")
	b.Base.Update(w, r)
	decoder := json.NewDecoder(r.Body)
	var c services.Customer
	err := decoder.Decode(&c)
	if err != nil {
		log.Println("ERROR", err)
		w.WriteHeader(http.StatusBadRequest)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	_, err = b.store.CustomerCreate(r.Context(), c)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	b.Base.Update(w, r)
	b.Base.MessageText = "customer added"
	b.customerListAll(w, r)
}

// --- EOF
