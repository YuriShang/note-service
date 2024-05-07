package note

import (
	"encoding/json"
	"fmt"
	"net/http"
	"note_service/app/internal/apperror"
	"note_service/app/internal/client/user_client"
	"note_service/app/pkg/logging"
	"note_service/app/pkg/user"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	notesURL = "/notes"
	noteURL  = "/notes/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	NoteService Service
	UserClient  user_client.UserClient
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc( // GET /notes
		http.MethodGet,
		notesURL,
		user.Authentication(h.UserClient, apperror.Middleware(h.GetNotes)),
	)
	router.HandlerFunc( // GET /note/{uuid}
		http.MethodGet,
		noteURL,
		user.Authentication(h.UserClient, apperror.Middleware(h.GetNote)),
	)
	router.HandlerFunc( // POST /notes
		http.MethodPost,
		notesURL,
		user.Authentication(h.UserClient, user.Authorization(apperror.Middleware(h.CreateNote))),
	)
	router.HandlerFunc( // PATCH /note/{uuid}
		http.MethodPatch,
		noteURL,
		user.Authentication(h.UserClient, user.Authorization(apperror.Middleware(h.UpdateNote))),
	)
	router.HandlerFunc( // DELETE /note/{uuid}
		http.MethodDelete,
		noteURL,
		user.Authentication(h.UserClient, user.Authorization(apperror.Middleware(h.DeleteNote))),
	)
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET NOTES")
	w.Header().Set("Content-Type", "application/json")

	userUUID := r.Context().Value("userUUID").(uuid.UUID)

	note, err := h.NoteService.GetMany(r.Context(), userUUID)
	if err != nil {
		return err
	}
	noteBytes, err := json.Marshal(note)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(noteBytes)

	return nil
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	strNoteUUID := params.ByName("uuid")
	if strNoteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}
	userUUID := r.Context().Value("userUUID").(uuid.UUID)

	noteUUID, err := uuid.Parse(strNoteUUID)
	if err != nil {
		return apperror.BadRequestError("invalid uuid type")
	}
	note, err := h.NoteService.GetOne(r.Context(), noteUUID, userUUID)
	if err != nil {
		return err
	}
	noteBytes, err := json.Marshal(note)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(noteBytes)

	return nil
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE NOTE")

	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get userUUID from context")
	userUUID := r.Context().Value("userUUID").(uuid.UUID)

	h.Logger.Debug("decode create note dto")
	var dto CreateNoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	dto.UserUUID = &userUUID
	noteUUID, err := h.NoteService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", notesURL, noteUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	strNoteUUID := params.ByName("uuid")

	h.Logger.Debug("get userUUID from context")
	userUUID := r.Context().Value("userUUID").(uuid.UUID)

	if strNoteUUID == "" {
		return apperror.BadRequestError("id query parameter is required and must be a comma separated integers")
	}

	h.Logger.Debug("decode update note dto")
	var dto UpdateNoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	noteUUID, err := uuid.Parse(strNoteUUID)
	if err != nil {
		return apperror.BadRequestError("invalid uuid type")
	}

	dto.NoteUUID = &noteUUID

	err = h.NoteService.Update(r.Context(), dto, userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	strNoteUUID := params.ByName("uuid")
	if strNoteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	noteUUID, err := uuid.Parse(strNoteUUID)
	if err != nil {
		return apperror.BadRequestError("invalid uuid type")
	}
	userUUID := r.Context().Value("userUUID").(uuid.UUID)
	err = h.NoteService.Delete(r.Context(), noteUUID, userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
