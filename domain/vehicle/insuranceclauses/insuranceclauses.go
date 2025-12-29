package insuranceclauses

import (
	"database/sql"
	"io"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// UpdateClauseRequest represents request to update a clause
type UpdateClauseRequest struct {
	UpdateData map[string]interface{}
}

// UpdateClauseDraftRequest represents request to update a draft clause
type UpdateClauseDraftRequest struct {
	ID            int
	InsuranceCode string
	ProductType   string
	Code          string
	Title         string
	Content       string
	Description   *string
	IsMandatory   *bool
	IsActive      *bool
}

// ConfirmClausesResult represents result of confirming clauses
type ConfirmClausesResult struct {
	InsertedCount int
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Template & Upload
	DownloadClausesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error)
	UploadClausesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error)

	// Draft operations
	GetClausesDraft(c echo.Context, insuranceCode string) ([]utils.ClausesDraft, error)
	UpdateClauseDraft(c echo.Context, req UpdateClauseDraftRequest) error
	DeleteClauseDraft(c echo.Context, id int) error
	ClearClausesDraft(c echo.Context) error
	ConfirmClauses(c echo.Context) (*ConfirmClausesResult, error)

	// Confirmed operations
	GetClauses(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	UpdateClause(c echo.Context, id int, req UpdateClauseRequest, userName string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Template & Upload
	GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]string, error)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Draft operations
	GetAllClausesDraft(c echo.Context, insuranceCode string) ([]utils.ClausesDraft, error)
	UpdateClauseDraftInMemory(c echo.Context, draftData *utils.DraftData, req UpdateClauseDraftRequest) error
	DeleteClauseDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error
	ClearClausesDraftInMemory(c echo.Context, draftData *utils.DraftData)

	// Confirmed operations
	GetAllClauses(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetClauseBeforeUpdate(c echo.Context, id int) (*ClauseBeforeUpdate, error)
	UpdateClause(c echo.Context, id int, updateData map[string]interface{}) (int64, error)
	InsertClause(c echo.Context, clause utils.ClausesDraft) error
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

// ClauseBeforeUpdate represents clause data before update (for history)
type ClauseBeforeUpdate struct {
	InsuranceCode sql.NullString
	ProductType   sql.NullString
	Code          sql.NullString
	Title         sql.NullString
	Content       sql.NullString
	Description   sql.NullString
	IsMandatory   sql.NullBool
	IsActive      sql.NullBool
}


