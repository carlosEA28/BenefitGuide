package server

import "net/http"

type sourceHandler struct {
	service officialSourceService
}

func newSourceHandler(service officialSourceService) *sourceHandler {
	return &sourceHandler{service: service}
}

func (h *sourceHandler) createSource(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	source, err := decodeSource(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	source.PublicServiceID = serviceID
	created, err := h.service.CreateOfficialSource(r.Context(), source)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sourceResponseFor(created))
}

func (h *sourceHandler) listSources(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	sources, err := h.service.ListOfficialSources(r.Context(), serviceID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	items := make([]sourceResponse, 0, len(sources))
	for i := range sources {
		items = append(items, sourceResponseFor(&sources[i]))
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *sourceHandler) getSource(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	source, err := h.service.GetOfficialSource(r.Context(), serviceID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sourceResponseFor(source))
}

func (h *sourceHandler) updateSource(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	changes, err := decodeSource(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	changes.ID = id
	updated, err := h.service.UpdateOfficialSource(r.Context(), serviceID, changes)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sourceResponseFor(updated))
}

func (h *sourceHandler) deleteSource(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.service.DeleteOfficialSource(r.Context(), serviceID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
