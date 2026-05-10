package products

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
	"github.com/CassRamos/go-ecom-api.git/internal/json"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id format", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		log.Println(err)

		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) FilterProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	nameParam := query.Get("name")
	categoryParam := query.Get("category")

	var params repository.FilterProductsParams

	if nameParam != "" {
		params.Name = pgtype.Text{String: nameParam, Valid: true}
	}

	if categoryParam != "" {
		params.Category = pgtype.Text{String: categoryParam, Valid: true}
	}

	params.MinPrice = parseIntToPgType(query.Get("min_price"))
	params.MaxPrice = parseIntToPgType(query.Get("max_price"))
	params.MinQuantity = parseIntToPgType(query.Get("min_quantity"))
	params.MaxQuantity = parseIntToPgType(query.Get("max_quantity"))

	products, err := h.service.FilterProducts(r.Context(), params)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to filter products", http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []repository.Product{}
	}

	json.Write(w, http.StatusOK, products)

}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var params repository.CreateProductParams

	if err := json.Read(r, &params); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), params)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, product)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id format", http.StatusBadRequest)
		return
	}

	var params repository.UpdateProductParams
	if err := json.Read(r, &params); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params.ID = id

	product, err := h.service.UpdateProduct(r.Context(), params)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func parseIntToPgType(s string) pgtype.Int4 {
	if s == "" {
		return pgtype.Int4{Valid: false}
	}

	parsed, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return pgtype.Int4{Valid: false}
	}

	return pgtype.Int4{Int32: int32(parsed), Valid: true}
}
