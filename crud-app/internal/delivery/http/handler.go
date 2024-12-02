package myhttp

import (
	"crud-app/internal/entity"
	"crud-app/internal/usecase"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	usecase   *usecase.RequestUsecase
	templates map[string]*template.Template
}

// Конструктор Handler с загрузкой шаблонов
func NewHandler(usecase *usecase.RequestUsecase) *Handler {
	return &Handler{
		usecase:   usecase,
		templates: loadTemplates(),
	}
}

// Загрузка HTML-шаблонов
func loadTemplates() map[string]*template.Template {
	return map[string]*template.Template{
		"index":    template.Must(template.ParseFiles("templates/index.html")),
		"create":   template.Must(template.ParseFiles("templates/create.html")),
		"requests": template.Must(template.ParseFiles("templates/requests.html")),
	}
}

// Главная страница
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.templates["index"].Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// Страница создания заявки
func (h *Handler) CreateRequestPage(w http.ResponseWriter, r *http.Request) {
	err := h.templates["create"].Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// Создание заявки (обработка формы)
func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			slog.Warn("Title or content missing in request")
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		request := entity.Request{
			Title:   title,
			Content: content,
			Status:  "Новая",
		}

		if err := h.usecase.CreateRequest(&request); err != nil {
			slog.Error("Failed to create request", slog.String("error", err.Error()))
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		slog.Info("Request created successfully", slog.String("title", title))
		http.Redirect(w, r, "/requests", http.StatusSeeOther)
		return
	}

	slog.Warn("Invalid request method")
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// Получение списка заявок с пагинацией
func (h *Handler) GetRequestsWithPagination(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 10
	}

	requests, err := h.usecase.GetRequestsWithPagination(offset, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve requests", http.StatusInternalServerError)
		return
	}

	err = h.templates["requests"].Execute(w, requests)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// Удаление заявки
func (h *Handler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteRequest(id); err != nil {
		http.Error(w, "Failed to delete request", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/requests", http.StatusSeeOther)
}
