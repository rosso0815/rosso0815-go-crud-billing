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
	store "github.com/rosso0815/rosso0815-go-crud-billing/services"
	ui "github.com/rosso0815/rosso0815-go-crud-billing/web/ui"
)

func NewInvoice(store *store.Store, sessionManager *scs.SessionManager, cfg *config.Config, r *router.Router) *Web {
	m := Web{}
	m.cfg = cfg
	m.store = store
	m.Path = "invoice"
	m.SessionManager = sessionManager
	m.AddAction = ui.CrudActionAdd(fmt.Sprintf("%s/ui/%s/formask", m.cfg.WebPrefix, m.Path))
	m.cfg.Menus = append(m.cfg.Menus, config.Menu{Name: "Invoice", Path: fmt.Sprintf("ui/%s", m.Path)})
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s", m.cfg.WebPrefix, m.Path), m.invoiceListAll)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/formedit/{id}", m.cfg.WebPrefix, m.Path), m.invoiceFormEdit)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/formadd", m.cfg.WebPrefix, m.Path), m.invoiceFormAdd)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/formask", m.cfg.WebPrefix, m.Path), m.invoiceFormAsk)
	r.RegisterRoute(fmt.Sprintf("POST %s/ui/%s/formpost", m.cfg.WebPrefix, m.Path), m.invoiceFormPost)
	r.RegisterRoute(fmt.Sprintf("GET %s/ui/%s/print/{id}", m.cfg.WebPrefix, m.Path), m.invoicePrint)
	r.RegisterRoute(fmt.Sprintf("DELETE %s/ui/%s/{id}", m.cfg.WebPrefix, m.Path), m.invoiceDelete)
	r.RegisterRoute(fmt.Sprintf("POST %s/ui/%s/payed/{id}", m.cfg.WebPrefix, m.Path), m.invoicePayed)
	return &m
}

func (b *Web) invoiceListAll(w http.ResponseWriter, r *http.Request) {
	b.Base.Update(w, r)
	var header []templ.Component = []templ.Component{
		ui.CrudHeaderSort("InvoiceId", b.cfg),
		ui.CrudHeaderSort("Year", b.cfg),
		ui.CrudHeaderSort("Month", b.cfg),
		ui.CrudHeaderSort("Sum H", b.cfg),
		ui.CrudHeaderSort("Sum Sfr", b.cfg),
		ui.CrudHeaderSort("CustomerId", b.cfg),
		ui.CrudHeaderSort("FirstName", b.cfg),
		ui.CrudHeaderSort("LastName", b.cfg),
	}
	// log.Println("web search", r.URL.Query().Get("search"))
	var components []templ.Component
	invoices, err := b.store.InvoiceListBySearch(r.Context(), r.URL.Query().Get("search"), b.PageSize, b.PageCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	count := len(invoices)
	for _, item := range invoices {
		l := invoiceCrudLine(item, b.cfg)
		components = append(components, l)
	}
	if err := ui.CrudTable(
		ui.CrudList{
			Base:   b.Base,
			Header: header,
			Items:  components,
			Count:  count,
		}, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b *Web) invoicePrint(w http.ResponseWriter, r *http.Request) {
	b.Base.Update(w, r)
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = "Invalid invoice ID"
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	invoice, err := b.store.InvoiceGetById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if err := print(invoice, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) invoiceDelete(w http.ResponseWriter, r *http.Request) {
	b.Base.Update(w, r)
	idParam, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	err = b.store.InvoiceDeleteById(r.Context(), idParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	b.Base.MessageText = fmt.Sprintf("delete done of Invoice-Id %d", idParam)
	if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) invoiceFormAsk(w http.ResponseWriter, r *http.Request) {
	var customers []store.Customer
	s_customers, _ := b.store.CustomerList(r.Context())
	customers = append(customers, s_customers...)
	if err := invoiceFormAsk(b.cfg, b.Base, "choose", customers).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (m Web) invoiceFormEdit(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	idParam, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR strconv went wrong")
	}
	invoice, err := m.store.InvoiceGetById(r.Context(), idParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if invoice.BillPayed {
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = "bill payed"
		m.invoiceListAll(w, r)
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if err := invoiceForm(m.cfg, m.Base, &invoice).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (b Web) invoiceFormPost(w http.ResponseWriter, r *http.Request) {
	var invoice store.Invoice
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&invoice)
	if err != nil {
		log.Println("Error decode", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	err = b.store.InvoiceSave(r.Context(), invoice)
	if err != nil {
		log.Println("Error InvoiceSet", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		b.Base.MessageType = ui.Alert
		b.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(b.Base, b.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	b.Base.Update(w, r)
	b.invoiceListAll(w, r)
}

func (m Web) invoiceFormAdd(w http.ResponseWriter, r *http.Request) {
	var err error
	customerID, err := strconv.Atoi(r.URL.Query().Get("customer"))
	if err != nil {
		customerID = 0
	}
	workMonth, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		workMonth = 0
	}
	workYear, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		workYear = 0
	}
	invoice, err := m.store.InvoiceGetByCustomerAndMonth(r.Context(), customerID, workMonth, workYear)
	if err != nil {
		log.Println("Error InvoiceSet", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	customer, err := m.store.CustomerGetById(r.Context(), customerID)
	log.Println("form add customer alv", customer.Alv)
	if err != nil {
		log.Println("Error InvoiceSet", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if invoice.BillPayed {
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = "bill payed"
		m.invoiceListAll(w, r)
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	if err := invoiceForm(m.cfg, m.Base, &invoice).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

func (m Web) invoicePayed(w http.ResponseWriter, r *http.Request) {
	idParam, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	invoice, err := m.store.InvoiceGetById(r.Context(), idParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = err.Error()
		if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
			log.Println("render error:", err)
		}
		return
	}
	log.Println("invoice.payed", invoice.BillPayed)
	if !invoice.BillPayed {
		err = m.store.InvoicePayed(r.Context(), idParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			m.Base.MessageType = ui.Alert
			m.Base.MessageText = err.Error()
			if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
				log.Println("render error:", err)
			}
			return
		}
		invoice.BillPayed = true
	} else {
		err = m.store.InvoiceNotPayed(r.Context(), idParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			m.Base.MessageType = ui.Alert
			m.Base.MessageText = err.Error()
			if err := ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w); err != nil {
				log.Println("render error:", err)
			}
			return
		}
		invoice.BillPayed = false

	}
	if err := invoiceButtonPayed(invoice).Render(r.Context(), w); err != nil {
		log.Println("render error:", err)
	}
}

// // --- EOF
