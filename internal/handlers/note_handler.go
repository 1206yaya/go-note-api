package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1206yaya/go-note-api/internal/services"
	"github.com/labstack/echo/v4"
)

type NoteHandler struct {
	service services.NoteService
}

func NewNoteHandler(service services.NoteService) *NoteHandler {
	return &NoteHandler{
		service: service,
	}
}

func (h *NoteHandler) CreateNote(c echo.Context) error {
	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	note, err := h.service.CreateNote(c.Request().Context(), req.Title, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create note"})
	}

	return c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetNote(c echo.Context) error {
	id := c.Param("id")

	note, err := h.service.GetNote(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get note"})
	}

	if note == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Note not found"})
	}

	return c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c echo.Context) error {
	id := c.Param("id")

	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	note, err := h.service.UpdateNote(c.Request().Context(), id, req.Title, req.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update note"})
	}

	return c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteNote(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete note"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *NoteHandler) ListNotes(c echo.Context) error {
	
	notes, err := h.service.ListNotes(c.Request().Context())
	if err != nil {
		// エラーの詳細をログに出力
		log.Printf("Failed to list notes: %v", err)

		// クライアントに返すエラーメッセージ
		errorMessage := fmt.Sprintf("Failed to list notes: %v", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": errorMessage,
		})
	}
	return c.JSON(http.StatusOK, notes)
}
