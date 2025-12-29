package http

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type handler struct{}

// NewHTMLConverterHandler creates a new HTML converter handler
func NewHTMLConverterHandler() *handler {
	return &handler{}
}

// ConvertHTMLRequest represents the request structure for HTML conversion
type ConvertHTMLRequest struct {
	TemplateType string                 `json:"template_type" binding:"required"`
	Content      map[string]interface{} `json:"content" binding:"required"`
}

// ConvertHTMLResponse represents the response structure
type ConvertHTMLResponse struct {
	HTMLContent string `json:"html_content"`
}

// ConvertHTML handles HTML conversion based on template type
func (h *handler) ConvertHTML(c echo.Context) error {
	var req ConvertHTMLRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "Failed",
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	// Only support INSURANCE_DETAIL_MOTOR_VEHICLE_CAR for now
	if req.TemplateType != "INSURANCE_DETAIL_MOTOR_VEHICLE_CAR" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "Failed",
			"message": fmt.Sprintf("Template type '%s' is not supported", req.TemplateType),
		})
	}

	// Extract content fields
	productInfo := getStringFromContent(req.Content, "product_info")
	additionalProtection := getArrayFromContent(req.Content, "additional_protection")
	termsConditions := getArrayFromContent(req.Content, "terms_conditions")
	ownRisk := getArrayFromContent(req.Content, "own_risk")

	// Generate HTML
	htmlContent := generateInsuranceDetailHTML(productInfo, additionalProtection, termsConditions, ownRisk)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"html_content": htmlContent,
	})
}

// Helper function to get string from content map
func getStringFromContent(content map[string]interface{}, key string) string {
	if val, ok := content[key]; ok {
		if str, ok := val.(string); ok {
			return strings.TrimSpace(str)
		}
	}
	return ""
}

// Helper function to get array from content map
func getArrayFromContent(content map[string]interface{}, key string) []string {
	if val, ok := content[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			var result []string
			for _, item := range arr {
				if str, ok := item.(string); ok {
					trimmed := strings.TrimSpace(str)
					if trimmed != "" {
						result = append(result, trimmed)
					}
				}
			}
			return result
		}
		// If it's already a []string, convert it
		if arr, ok := val.([]string); ok {
			var result []string
			for _, item := range arr {
				trimmed := strings.TrimSpace(item)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
			return result
		}
	}
	return []string{}
}

// generateInsuranceDetailHTML generates HTML content for INSURANCE_DETAIL_MOTOR_VEHICLE_CAR template
// Returns compact HTML without unnecessary whitespace
func generateInsuranceDetailHTML(productInfo string, additionalProtection, termsConditions, ownRisk []string) string {
	var sb strings.Builder

	// Start HTML structure (compact - no newlines or extra spaces)
	sb.WriteString("<div class=\"insurance-detail\">")

	// Product Info Section
	if productInfo != "" {
		sb.WriteString("<div class=\"product-info-section\"><p>")
		sb.WriteString(html.EscapeString(productInfo))
		sb.WriteString("</p></div>")
	}

	// Additional Protection Section
	if len(additionalProtection) > 0 {
		sb.WriteString("<div class=\"additional-protection-section\"><h4>Perlindungan Tambahan</h4><ul>")
		for _, item := range additionalProtection {
			if item != "" {
				sb.WriteString("<li>")
				sb.WriteString(html.EscapeString(item))
				sb.WriteString("</li>")
			}
		}
		sb.WriteString("</ul></div>")
	}

	// Terms and Conditions Section
	if len(termsConditions) > 0 {
		sb.WriteString("<div class=\"terms-conditions-section\"><h4>Syarat dan Ketentuan</h4><ol>")
		for _, item := range termsConditions {
			if item != "" {
				sb.WriteString("<li>")
				sb.WriteString(html.EscapeString(item))
				sb.WriteString("</li>")
			}
		}
		sb.WriteString("</ol></div>")
	}

	// Own Risk Section
	if len(ownRisk) > 0 {
		sb.WriteString("<div class=\"own-risk-section\"><h4>Risiko Sendiri</h4><ul>")
		for _, item := range ownRisk {
			if item != "" {
				sb.WriteString("<li>")
				sb.WriteString(html.EscapeString(item))
				sb.WriteString("</li>")
			}
		}
		sb.WriteString("</ul></div>")
	}

	// End HTML structure
	sb.WriteString("</div>")

	return sb.String()
}


