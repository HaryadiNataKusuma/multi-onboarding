package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HTMLHandler handles HTML conversion requests
type HTMLHandler struct{}

// NewHTMLHandler creates a new HTML handler
func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

// ConvertHTML handles POST /api/html/convert
func (h *HTMLHandler) ConvertHTML(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		TemplateType string `json:"template_type"`
		Content      struct {
			ProductInfo          string   `json:"product_info"`
			AdditionalProtection []string `json:"additional_protection"`
			TermsConditions      []string `json:"terms_conditions"`
		} `json:"content"`
		// Alternative format: direct fields (for backward compatibility)
		ProductInfo          string   `json:"productInfo"`
		AdditionalProtection []string `json:"additionalProtection"`
		TermsConditions      []string `json:"termsConditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Support both formats: nested Content or direct fields
	var productInfo string
	var additionalProtection, termsConditions []string

	if requestData.Content.ProductInfo != "" || len(requestData.Content.AdditionalProtection) > 0 || len(requestData.Content.TermsConditions) > 0 {
		// Use nested Content format
		productInfo = requestData.Content.ProductInfo
		additionalProtection = requestData.Content.AdditionalProtection
		termsConditions = requestData.Content.TermsConditions
	} else {
		// Use direct fields format
		productInfo = requestData.ProductInfo
		additionalProtection = requestData.AdditionalProtection
		termsConditions = requestData.TermsConditions
	}

	// Generate HTML content
	htmlContent := generateInsuranceDetailHTML(
		productInfo,
		additionalProtection,
		termsConditions,
	)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"html_content": htmlContent,
	})
}

// generateInsuranceDetailHTML generates HTML for INSURANCE_DETAIL_TV template
// Supports nested bullet points within numbered lists
func generateInsuranceDetailHTML(productInfo string, additionalProtection, termsConditions []string) string {
	var html strings.Builder

	html.WriteString(`<div class="insurance-detail">`)

	// Product Info Section (DETAIL PRODUK)
	if productInfo != "" {
		html.WriteString(`<div class="product-info-section">`)
		html.WriteString(`<h4>DETAIL PRODUK</h4>`)
		html.WriteString(`<p>`)
		html.WriteString(escapeHTML(productInfo))
		html.WriteString(`</p>`)
		html.WriteString(`</div>`)
	}

	// Additional Protection Section (PENGECUALIAN UMUM) - using ordered list with nested bullets
	if len(additionalProtection) > 0 {
		html.WriteString(`<div class="additional-protection-section">`)
		html.WriteString(`<h4>PENGECUALIAN UMUM</h4>`)
		html.WriteString(generateNumberedListWithNestedBullets(additionalProtection))
		html.WriteString(`</div>`)
	}

	// Terms and Conditions Section (SYARAT DAN KETENTUAN) - using ordered list with nested bullets
	if len(termsConditions) > 0 {
		html.WriteString(`<div class="terms-conditions-section">`)
		html.WriteString(`<h4>SYARAT DAN KETENTUAN</h4>`)
		html.WriteString(generateNumberedListWithNestedBullets(termsConditions))
		html.WriteString(`</div>`)
	}

	html.WriteString(`</div>`)

	return html.String()
}

// generateNumberedListWithNestedBullets generates <ol> with nested <ul> for bullet points
// Format: Lines starting with "-" = bullet items (can be indented or not), other lines = numbered items
func generateNumberedListWithNestedBullets(items []string) string {
	var html strings.Builder
	html.WriteString(`<ol>`)

	var currentBullets []string
	var lastNumberedItem string

	for i, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}

		// Check if it starts with bullet marker "-" or "•" (with or without indentation)
		hasBulletMarker := strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "•")
		
		// Remove leading dashes/bullets and trim
		content := trimmed
		if hasBulletMarker {
			// Remove "-" or "•" and any leading spaces after it
			content = strings.TrimPrefix(trimmed, "-")
			content = strings.TrimPrefix(content, "•")
			content = strings.TrimSpace(content)
		}

		// If has bullet marker, it's a bullet item (regardless of indentation)
		if hasBulletMarker {
			// Add bullet to current list (will be attached to previous numbered item)
			currentBullets = append(currentBullets, content)
		} else {
			// It's a numbered item - close previous numbered item if it exists
			if lastNumberedItem != "" {
				html.WriteString(`<li>`)
				html.WriteString(escapeHTML(lastNumberedItem))
				if len(currentBullets) > 0 {
					html.WriteString(`<ul>`)
					for _, bullet := range currentBullets {
						html.WriteString(`<li>`)
						html.WriteString(escapeHTML(bullet))
						html.WriteString(`</li>`)
					}
					html.WriteString(`</ul>`)
					currentBullets = []string{}
				}
				html.WriteString(`</li>`)
			}
			lastNumberedItem = content
		}

		// If this is the last item, close it
		if i == len(items)-1 {
			if hasBulletMarker {
				// Last item is a bullet - need to attach to previous numbered item or create new one
				if lastNumberedItem != "" {
					html.WriteString(`<li>`)
					html.WriteString(escapeHTML(lastNumberedItem))
					if len(currentBullets) > 0 {
						html.WriteString(`<ul>`)
						for _, bullet := range currentBullets {
							html.WriteString(`<li>`)
							html.WriteString(escapeHTML(bullet))
							html.WriteString(`</li>`)
						}
						html.WriteString(`</ul>`)
					}
					html.WriteString(`</li>`)
				} else if len(currentBullets) > 0 {
					// Only bullets, no numbered item - create empty numbered item with bullets
					html.WriteString(`<li>`)
					html.WriteString(`<ul>`)
					for _, bullet := range currentBullets {
						html.WriteString(`<li>`)
						html.WriteString(escapeHTML(bullet))
						html.WriteString(`</li>`)
					}
					html.WriteString(`</ul>`)
					html.WriteString(`</li>`)
				}
			} else {
				// Last item is a numbered item
				html.WriteString(`<li>`)
				html.WriteString(escapeHTML(lastNumberedItem))
				if len(currentBullets) > 0 {
					html.WriteString(`<ul>`)
					for _, bullet := range currentBullets {
						html.WriteString(`<li>`)
						html.WriteString(escapeHTML(bullet))
						html.WriteString(`</li>`)
					}
					html.WriteString(`</ul>`)
				}
				html.WriteString(`</li>`)
			}
		}
	}

	html.WriteString(`</ol>`)
	return html.String()
}

// escapeHTML escapes HTML special characters
// Note: Must escape & first to avoid double-escaping
func escapeHTML(s string) string {
	// Replace & first to avoid double-escaping
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

