// Package settings regroups some API methods to facilitate the usage of the
// io.cozy settings documents.
package settings

import (
	"encoding/json"
	"net/http"

	"github.com/cozy/cozy-stack/pkg/consts"
	"github.com/cozy/cozy-stack/pkg/couchdb"
	"github.com/cozy/cozy-stack/pkg/sessions"
	"github.com/cozy/cozy-stack/web/jsonapi"
	"github.com/cozy/cozy-stack/web/middlewares"
	"github.com/cozy/cozy-stack/web/permissions"
	"github.com/labstack/echo"
)

type apiSession struct {
	s *sessions.Session
}

func (s *apiSession) ID() string                             { return s.s.ID() }
func (s *apiSession) Rev() string                            { return s.s.Rev() }
func (s *apiSession) DocType() string                        { return consts.Sessions }
func (s *apiSession) Clone() couchdb.Doc                     { return s }
func (s *apiSession) SetID(_ string)                         {}
func (s *apiSession) SetRev(_ string)                        {}
func (s *apiSession) Relationships() jsonapi.RelationshipMap { return nil }
func (s *apiSession) Included() []jsonapi.Object             { return nil }
func (s *apiSession) Links() *jsonapi.LinksList              { return nil }
func (s *apiSession) MarshalJSON() ([]byte, error)           { return json.Marshal(s.s) }

func getSessions(c echo.Context) error {
	inst := middlewares.GetInstance(c)

	if err := permissions.AllowWholeType(c, permissions.GET, consts.Sessions); err != nil {
		return err
	}

	sessions, err := sessions.GetAll(inst)
	if err != nil {
		return err
	}

	objs := make([]jsonapi.Object, len(sessions))
	for i, s := range sessions {
		objs[i] = &apiSession{s}
	}

	return jsonapi.DataList(c, http.StatusOK, objs, nil)
}

// Routes sets the routing for the settings service
func Routes(router *echo.Group) {
	router.GET("/disk-usage", diskUsage)

	router.POST("/passphrase", registerPassphrase)
	router.PUT("/passphrase", updatePassphrase)

	router.GET("/instance", getInstance)
	router.PUT("/instance", updateInstance)
	router.GET("/sessions", getSessions)

	router.GET("/clients", listClients)
	router.DELETE("/clients/:id", revokeClient)
	router.POST("/synchronized", synchronized)

	router.GET("/onboarded", onboarded)
	router.GET("/context", context)
}
