<template>
  <div class="col-12">
    <div class="form-section templates-section">
      <div class="section-header">
        <div class="header-icon">
          <i class="fas fa-question"></i>
        </div>
        <h3>Templates Management</h3>
      </div>
      
      <div class="d-flex gap-3 flex-wrap mb-3">
        <button class="btn btn-primary" @click="triggerFileInput">
          <i class="fas fa-upload"></i> Upload CSV
        </button>
        <input 
          type="file" 
          ref="fileInput"
          accept=".xlsx,.xls,.csv"
          @change="handleFileChange"
          style="display: none;"
        >
        <button class="btn btn-info" @click="showDraft">
          <i class="fas fa-eye"></i> Tampilkan Draft Data
        </button>
        <button class="btn btn-warning" @click="clearDraft">
          <i class="fas fa-trash"></i> Clear Draft
        </button>
        <button class="btn btn-info" @click="showResult">
          <i class="fas fa-chart-line"></i> RESULT
        </button>
      </div>

      <!-- Add Template Form Section (Hidden) -->
      <div v-if="false" class="form-section add-form-section">
        <div class="section-header">
          <div class="header-icon" style="background-color: #17a2b8;">
            <i class="fas fa-plus"></i>
          </div>
          <h3>Tambah Template Baru</h3>
        </div>
        
        <form @submit.prevent="saveTemplate">
          <div class="form-grid form-grid-3">
            <div>
              <label class="form-label">Locale <span class="text-danger">*</span></label>
              <input 
                type="text" 
                class="form-control" 
                v-model="form.locale" 
                placeholder="Cth: id" 
                required>
            </div>
            <div>
              <label class="form-label">Template ID <span class="text-danger">*</span></label>
              <input 
                type="text" 
                class="form-control" 
                v-model="form.id" 
                placeholder="Template ID" 
                required>
            </div>
            <div>
              <label class="form-label">Product Code</label>
              <input 
                type="text" 
                class="form-control" 
                v-model="form.product_code" 
                placeholder="Product Code">
            </div>
          </div>
          
          <!-- HTML Converter Section (only for INSURANCE_DETAIL_TRAVEL_) -->
          <div v-if="isAddFormInsuranceDetailCar" class="html-converter-section mt-3 mb-3">
            <div class="html-converter-card">
              <div class="html-converter-header">
                <h6 class="mb-0">
                  <i class="fas fa-code me-2"></i>HTML Converter
                </h6>
              </div>
              <div class="html-converter-body">
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">DETAIL PRODUK</label>
                    <textarea 
                      class="form-control" 
                      v-model="addHtmlConverterForm.productInfo" 
                      rows="3"
                      placeholder="Masukkan detail produk asuransi...">
                    </textarea>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">PENGECUALIAN UMUM (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="addHtmlConverterForm.additionalProtection" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">SYARAT DAN KETENTUAN (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="addHtmlConverterForm.termsConditions" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="d-flex gap-2">
                  <button type="button" class="btn btn-primary btn-sm" @click="convertToHtmlForAdd">
                    <i class="fas fa-code"></i> Convert to HTML
                  </button>
                  <button type="button" class="btn btn-secondary btn-sm" @click="previewHtmlForAdd" :disabled="!addCurrentHtmlContent && !form.value">
                    <i class="fas fa-eye"></i> Preview
                  </button>
                </div>
              </div>
            </div>
            <small class="text-info mt-2 d-block">
              <i class="fas fa-info-circle"></i> Gunakan HTML Converter di atas untuk input data, atau edit langsung di textarea Value
            </small>
          </div>
          
          <div class="form-grid form-grid-1 mt-3">
            <div>
              <label class="form-label">Value <span class="text-danger">*</span></label>
              <textarea 
                class="form-control" 
                v-model="form.value" 
                rows="6"
                :placeholder="isAddFormInsuranceDetailCar ? 'HTML content akan muncul di sini setelah convert' : 'Template value/content'" 
                required></textarea>
            </div>
          </div>
          
          <div class="mt-3 d-flex align-items-center gap-3">
            <button type="submit" class="btn btn-success">
              <i class="fas fa-save"></i> Save Template
            </button>
            <button type="button" class="btn btn-secondary" @click="resetForm">
              <i class="fas fa-redo"></i> Reset
            </button>
          </div>
        </form>
      </div>

      <!-- Bulk Edit Form Section -->
      <div v-if="showBulkEditForm" class="form-section edit-form-section">
        <div class="section-header">
          <div class="header-icon" style="background-color: #ffc107;">
            <i class="fas fa-edit"></i>
          </div>
          <h3>Bulk Edit Templates ({{ isBulkEditingDraft ? 'Draft' : 'Result' }} - {{ isBulkEditingDraft ? selectedDraftItems.length : selectedResultItems.length }} items)</h3>
        </div>
        
        <form @submit.prevent="updateBulkTemplates">
          <div class="alert alert-info">
            <i class="fas fa-info-circle me-2"></i>
            You are editing <strong>{{ isBulkEditingDraft ? selectedDraftItems.length : selectedResultItems.length }} template(s)</strong>.
            The value below will be applied to all selected templates.
          </div>
          
          <!-- HTML Converter Section (only for INSURANCE_DETAIL_TV) -->
          <div v-if="isBulkEditInsuranceDetail" class="html-converter-section mt-3 mb-3">
            <div class="html-converter-card">
              <div class="html-converter-header">
                <h6 class="mb-0">
                  <i class="fas fa-code me-2"></i>HTML Converter
                </h6>
              </div>
              <div class="html-converter-body">
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">DETAIL PRODUK</label>
                    <textarea 
                      class="form-control" 
                      v-model="bulkEditHtmlConverterForm.productInfo" 
                      rows="3"
                      placeholder="Masukkan detail produk asuransi...">
                    </textarea>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">PENGECUALIAN UMUM (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="bulkEditHtmlConverterForm.additionalProtection" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">SYARAT DAN KETENTUAN (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="bulkEditHtmlConverterForm.termsConditions" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="d-flex gap-2">
                  <button type="button" class="btn btn-primary btn-sm" @click="convertToHtmlForBulkEdit">
                    <i class="fas fa-code"></i> Convert to HTML
                  </button>
                  <button type="button" class="btn btn-secondary btn-sm" @click="previewHtmlForBulkEdit" :disabled="!bulkEditCurrentHtmlContent && !bulkEditForm.value">
                    <i class="fas fa-eye"></i> Preview
                  </button>
                </div>
              </div>
            </div>
            <small class="text-info mt-2 d-block">
              <i class="fas fa-info-circle"></i> Gunakan HTML Converter di atas untuk input data, atau edit langsung di textarea Value
            </small>
          </div>
          
          <div class="form-grid form-grid-1 mt-3">
            <div>
              <label class="form-label">Value <span class="text-danger">*</span></label>
              <textarea 
                class="form-control" 
                v-model="bulkEditForm.value" 
                rows="6"
                :placeholder="isBulkEditInsuranceDetail ? 'HTML content akan muncul di sini setelah convert' : 'Template value/content (will be applied to all selected templates)'" 
                required></textarea>
            </div>
          </div>
          
          <div class="mt-3 d-flex align-items-center gap-3">
            <button type="submit" class="btn btn-success">
              <i class="fas fa-check"></i> Update All Selected Templates
            </button>
            <button type="button" class="btn btn-secondary" @click="cancelBulkEdit">
              <i class="fas fa-times"></i> Cancel
            </button>
          </div>
        </form>
      </div>

      <!-- Edit Form Section -->
      <div v-if="showEditForm && !showBulkEditForm" class="form-section edit-form-section">
        <div class="section-header">
          <div class="header-icon" style="background-color: #ffc107;">
            <i class="fas fa-edit"></i>
          </div>
          <h3>{{ isEditingResult ? 'Edit Data Template (From Database)' : 'Edit Data Template (From Draft)' }}</h3>
        </div>
        
        <form @submit.prevent="updateTemplate">
          <input type="hidden" v-model="editForm.id">
          <div class="form-grid form-grid-3">
            <div>
              <label class="form-label">Locale <span class="text-danger">*</span></label>
              <input 
                type="text" 
                class="form-control form-control-readonly" 
                v-model="editForm.locale" 
                readonly
                placeholder="Cth: id" 
                required>
            </div>
            <div>
              <label class="form-label">Template ID <span class="text-danger">*</span></label>
              <input 
                type="text" 
                class="form-control form-control-readonly" 
                v-model="editForm.id" 
                readonly
                placeholder="Template ID" 
                required>
            </div>
            <div>
              <label class="form-label">Product Code</label>
              <input 
                type="text" 
                class="form-control form-control-readonly" 
                v-model="editForm.product_code" 
                readonly
                placeholder="Product Code">
            </div>
          </div>
          
          <!-- HTML Converter Section (only for INSURANCE_DETAIL_TV) -->
          <div v-if="isInsuranceDetailCar" class="html-converter-section mt-3 mb-3">
            <div class="html-converter-card">
              <div class="html-converter-header">
                <h6 class="mb-0">
                  <i class="fas fa-code me-2"></i>HTML Converter
                </h6>
              </div>
              <div class="html-converter-body">
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">DETAIL PRODUK</label>
                    <textarea 
                      class="form-control" 
                      v-model="htmlConverterForm.productInfo" 
                      rows="3"
                      placeholder="Masukkan detail produk asuransi...">
                    </textarea>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">PENGECUALIAN UMUM (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="htmlConverterForm.additionalProtection" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="form-grid form-grid-1 mb-3">
                  <div>
                    <label class="form-label">SYARAT DAN KETENTUAN (pisahkan dengan enter, gunakan "-" untuk bullet)</label>
                    <textarea 
                      class="form-control" 
                      v-model="htmlConverterForm.termsConditions" 
                      rows="5"
                      placeholder="Item pertama&#10;Item kedua&#10;- Sub bullet 1&#10;- Sub bullet 2&#10;Item ketiga&#10;...">
                    </textarea>
                    <small class="text-muted">Baris yang dimulai dengan "-" akan menjadi bullet di dalam numbering item sebelumnya</small>
                  </div>
                </div>
                
                <div class="d-flex gap-2">
                  <button type="button" class="btn btn-primary btn-sm" @click="convertToHtml">
                    <i class="fas fa-code"></i> Convert to HTML
                  </button>
                  <button type="button" class="btn btn-secondary btn-sm" @click="previewHtml" :disabled="!currentHtmlContent && !editForm.value">
                    <i class="fas fa-eye"></i> Preview
                  </button>
                </div>
              </div>
            </div>
            <small class="text-info mt-2 d-block">
              <i class="fas fa-info-circle"></i> Gunakan HTML Converter di atas untuk input data, atau edit langsung di textarea Value
            </small>
          </div>
          
          <div class="form-grid form-grid-1 mt-3">
            <div>
              <label class="form-label">Value <span class="text-danger">*</span></label>
              <textarea 
                class="form-control" 
                v-model="editForm.value" 
                rows="6"
                :placeholder="isInsuranceDetailCar ? 'HTML content akan muncul di sini setelah convert' : 'Template value/content'" 
                required></textarea>
            </div>
          </div>
          
          <div class="mt-3 d-flex align-items-center gap-3">
            <button type="submit" class="btn btn-success">
              <i class="fas fa-check"></i> Update Template
            </button>
            <button type="button" class="btn btn-secondary" @click="cancelEdit">
              <i class="fas fa-times"></i> Cancel
            </button>
          </div>
        </form>
      </div>

      <!-- Draft Table -->
      <div v-if="showDraftSection && !showEditForm && !showBulkEditForm" class="table-section draft-section show">
        <div class="table-header d-flex justify-content-between align-items-center">
          <h4>📊 Data Preview: Templates (Draft Save)</h4>
          <div>
            <button 
              class="btn btn-info btn-sm me-2" 
              @click="openGenerateTemplatesModal"
            >
              <i class="fas fa-magic"></i> Generate from Products
            </button>
            <button 
              class="btn btn-warning btn-sm me-2" 
              @click="openBulkEdit(true)" 
              :disabled="!hasSelectedDraftItems"
              v-if="hasSelectedDraftItems"
            >
              <i class="fas fa-edit"></i> Bulk Edit ({{ selectedDraftItems.length }})
            </button>
            <button class="btn btn-success btn-sm me-2" @click="confirmTemplates" :disabled="!canConfirm">
              <i class="fas fa-check"></i> Confirm Templates Data
            </button>
            <button type="button" class="btn btn-secondary btn-sm" @click="hideDraft">
              <i class="fas fa-times"></i> Close
            </button>
          </div>
        </div>
        <div class="table-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead>
                <tr>
                  <th style="width: 50px;">
                    <input 
                      type="checkbox" 
                      :checked="isAllDraftSelected"
                      @change="toggleSelectAllDraft"
                      class="form-check-input"
                    >
                  </th>
                  <th>Locale</th>
                  <th>Template ID</th>
                  <th>Insurance</th>
                  <th>Product Name</th>
                  <th>Value</th>
                  <th>Created By</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="draftData.length === 0">
                  <td colspan="8" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No draft data available
                  </td>
                </tr>
                <tr v-for="(item, index) in draftData" :key="index">
                  <td>
                    <input 
                      type="checkbox" 
                      :checked="isDraftSelected(item)"
                      @change="toggleDraftSelection(item)"
                      class="form-check-input"
                    >
                  </td>
                  <td>{{ item.locale }}</td>
                  <td>{{ item.id }}</td>
                  <td>
                    <span v-if="item.insurance" class="badge bg-info">
                      {{ Array.isArray(item.insurance) ? item.insurance.join(', ') : item.insurance }}
                    </span>
                    <span v-else-if="item.insurance_code" class="badge bg-info">
                      {{ item.insurance_code }}
                    </span>
                    <span v-else class="text-muted">-</span>
                  </td>
                  <td>{{ getProductNameFromTemplateId(item.id) || getProductNameFromProductCode(item.product_code) || '-' }}</td>
                  <td>
                    <div style="max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.value">
                      {{ item.value || '-' }}
                    </div>
                  </td>
                  <td>{{ item.created_by }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="editDraft(item)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteDraft(index)">
                      <i class="fas fa-trash"></i> Hapus
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Result Table Section -->
      <div v-if="showResultSection && !showEditForm && !showBulkEditForm" class="table-section result-section show">
        <div class="table-header d-flex justify-content-between align-items-center">
          <h4>📈 Data Result: Templates (From Database - Summary & Insurance Detail)</h4>
          <div>
            <button 
              class="btn btn-warning btn-sm me-2" 
              @click="openBulkEdit(false)" 
              :disabled="!hasSelectedResultItems"
              v-if="hasSelectedResultItems"
            >
              <i class="fas fa-edit"></i> Bulk Edit ({{ selectedResultItems.length }})
            </button>
            <button class="btn btn-info btn-sm me-2" @click="showResult">
              <i class="fas fa-sync-alt"></i> Refresh Data
            </button>
            <button type="button" class="btn btn-secondary btn-sm" @click="hideResult">
              <i class="fas fa-times"></i> Close
            </button>
          </div>
        </div>

        <div class="table-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead class="table-success">
                <tr>
                  <th style="width: 50px;">
                    <input 
                      type="checkbox" 
                      :checked="isAllResultSelected"
                      @change="toggleSelectAllResult"
                      class="form-check-input"
                    >
                  </th>
                  <th>Locale</th>
                  <th>Template ID</th>
                  <th>Insurance Code</th>
                  <th>Product Name</th>
                  <th>Value</th>
                  <th>Created By</th>
                  <th>Created At</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="resultData.length === 0">
                  <td colspan="9" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No templates found (Summary & Insurance Detail)
                  </td>
                </tr>
                <tr v-for="item in resultData" :key="item.id">
                  <td>
                    <input 
                      type="checkbox" 
                      :checked="isResultSelected(item)"
                      @change="toggleResultSelection(item)"
                      class="form-check-input"
                    >
                  </td>
                  <td>{{ item.locale }}</td>
                  <td>
                    {{ item.id }}
                  </td>
                  <td>
                    <span v-if="getInsuranceCodeFromTemplateId(item.id)" class="badge bg-info">
                      {{ getInsuranceCodeFromTemplateId(item.id) }}
                    </span>
                    <span v-else-if="item.insurance_code" class="badge bg-info">
                      {{ item.insurance_code }}
                    </span>
                    <span v-else-if="item.insurance" class="badge bg-info">
                      {{ Array.isArray(item.insurance) ? item.insurance.join(', ') : item.insurance }}
                    </span>
                    <span v-else class="text-muted">-</span>
                  </td>
                  <td>{{ getProductNameFromTemplateId(item.id) || '-' }}</td>
                  <td>
                    <div style="max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.value">
                      {{ item.value || '-' }}
                    </div>
                  </td>
                  <td>{{ item.created_by || '-' }}</td>
                  <td>{{ formatDateOnly(item.created_at) }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning" @click="editResult(item)">
                      <i class="fas fa-edit"></i> Edit
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Generate Templates Modal -->
    <div v-if="showGenerateTemplatesModal" class="modal-backdrop" @click="closeGenerateTemplatesModal"></div>
    <div v-if="showGenerateTemplatesModal" class="modal-overlay">
      <div class="modal-content generate-templates-modal">
        <div class="modal-header">
          <div class="modal-header-content">
            <i class="fas fa-magic modal-header-icon"></i>
            <h5 class="modal-title">Generate Templates from Products</h5>
          </div>
          <button type="button" class="btn-close-modal" @click="closeGenerateTemplatesModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <!-- Header Section -->
          <div class="modal-section-header">
            <div class="d-flex justify-content-between align-items-center mb-3">
              <label class="form-label-modal">
                <i class="fas fa-box me-2"></i>Select Products
              </label>
              <button 
                type="button" 
                class="btn btn-sm btn-outline-primary btn-select-all"
                @click="toggleSelectAllProducts"
              >
                <i class="fas" :class="allProductsSelected ? 'fa-square' : 'fa-check-square'"></i>
                {{ allProductsSelected ? 'Deselect All' : 'Select All' }}
              </button>
            </div>
            
            <!-- Search -->
            <div class="search-wrapper">
              <i class="fas fa-search search-icon"></i>
              <input 
                type="text" 
                class="form-control search-input" 
                v-model="productSearchQuery"
                placeholder="Search products by code or name..."
              />
            </div>
          </div>
          
          <!-- Products List -->
          <div class="products-list-container">
            <div v-if="filteredProductsForGeneration.length === 0" class="empty-state">
              <i class="fas fa-inbox empty-icon"></i>
              <p v-if="productsForGeneration.length === 0" class="empty-text">Loading products...</p>
              <p v-else class="empty-text">No products found matching "{{ productSearchQuery }}"</p>
            </div>
            
            <div 
              v-for="product in filteredProductsForGeneration" 
              :key="product.id || product.code"
              class="product-item-list"
              :class="{ 'product-item-list-selected': selectedProductsForGeneration.includes(product.code) }"
              @click="toggleProductSelection(product.code)"
            >
              <input 
                class="product-checkbox-list" 
                type="checkbox" 
                :id="'gen-product-' + product.code"
                :value="product.code"
                v-model="selectedProductsForGeneration"
                @click.stop
              >
              <label class="product-label-list" :for="'gen-product-' + product.code">
                <span class="product-code-list">{{ product.code }}</span>
                <span class="product-name-list">{{ product.name }}</span>
                <span class="product-badge-list">
                  <i class="fas fa-shield-alt"></i> {{ product.insurance_code }}
                </span>
              </label>
            </div>
          </div>
          
          <!-- Selection Summary -->
          <div v-if="selectedProductsForGeneration.length > 0" class="selection-summary">
            <i class="fas fa-check-circle summary-icon"></i>
            <span class="summary-text">
              <strong>{{ selectedProductsForGeneration.length }}</strong> 
              {{ selectedProductsForGeneration.length === 1 ? 'product' : 'products' }} selected
            </span>
          </div>
        </div>
        
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary btn-modal" @click="closeGenerateTemplatesModal">
            <i class="fas fa-times me-2"></i>Cancel
          </button>
          <button 
            type="button" 
            class="btn btn-primary btn-modal" 
            @click="generateTemplatesFromSelectedProducts"
            :disabled="selectedProductsForGeneration.length === 0 || isGenerating"
          >
            <i class="fas fa-magic me-2"></i>
            {{ isGenerating ? 'Generating...' : `Generate Templates (${selectedProductsForGeneration.length})` }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// import { fetchWithAuth } from '../utils/apiHelper'

export default {
  name: 'TemplatesSection',
  props: {
    selectedInsuranceCode: String
  },
  data() {
    return {
      draftData: [],
      showDraftSection: false,
      canConfirm: false,
      resultData: [],
      showResultSection: false,
      showEditForm: false,
      isEditingResult: false,
      form: {
        locale: 'id',
        id: '',
        value: '',
        product_code: '',
        insurance_code: ''
      },
      editForm: {
        id: null,
        locale: '',
        id: '',
        value: '',
        product_code: '',
        insurance_code: ''
      },
      originalEditForm: null,
      API_URL: 'http://localhost:8080/api/travel',
      htmlConverterForm: {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      },
      addHtmlConverterForm: {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      },
      currentHtmlContent: '',
      addCurrentHtmlContent: '',
      showHtmlPreview: false,
      selectedDraftItems: [],
      selectedResultItems: [],
      showBulkEditForm: false,
      bulkEditForm: {
        value: ''
      },
      isBulkEditingDraft: false,
      products: [],
      bulkEditHtmlConverterForm: {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      },
      bulkEditCurrentHtmlContent: '',
      file: null,
      showGenerateTemplatesModal: false,
      productsForGeneration: [],
      selectedProductsForGeneration: [],
      productSearchQuery: '',
      isGenerating: false
    }
  },
  computed: {
    hasSelectedDraftItems() {
      return this.selectedDraftItems.length > 0
    },
    hasSelectedResultItems() {
      return this.selectedResultItems.length > 0
    },
    isAllDraftSelected() {
      return this.draftData.length > 0 && this.selectedDraftItems.length === this.draftData.length
    },
    isAllResultSelected() {
      return this.resultData.length > 0 && this.selectedResultItems.length === this.resultData.length
    },
    isInsuranceDetailCar() {
      return this.editForm.id && 
             this.editForm.id.includes('INSURANCE_DETAIL_TV')
    },
    isAddFormInsuranceDetailCar() {
      return this.form.id && 
             this.form.id.includes('INSURANCE_DETAIL_TV')
    },
    isBulkEditInsuranceDetail() {
      const itemsToCheck = this.isBulkEditingDraft ? this.selectedDraftItems : this.selectedResultItems
      if (itemsToCheck.length === 0) return false
      // Check if all selected items are INSURANCE_DETAIL_TV
      return itemsToCheck.every(item => {
        const templateId = item.item?.id || item.id
        return templateId && templateId.includes('INSURANCE_DETAIL_TV')
      })
    },
    filteredProductsForGeneration() {
      if (!this.productSearchQuery) {
        return this.productsForGeneration
      }
      const query = this.productSearchQuery.toLowerCase()
      return this.productsForGeneration.filter(product => 
        (product.code && product.code.toLowerCase().includes(query)) ||
        (product.name && product.name.toLowerCase().includes(query)) ||
        (product.insurance_code && product.insurance_code.toLowerCase().includes(query))
      )
    },
    allProductsSelected() {
      return this.productsForGeneration.length > 0 && 
             this.selectedProductsForGeneration.length === this.productsForGeneration.length
    }
  },
  watch: {
    selectedInsuranceCode(newVal) {
      // Update form insurance_code when selected insurance changes
      if (newVal) {
        this.form.insurance_code = newVal
        // Reload products when insurance code changes
        this.loadProducts()
        // Only auto-load draft if explicitly requested, not on every change
        // This prevents unwanted auto-loading when products are confirmed
        if (this.showDraftSection) {
          this.loadDraft()
        }
      } else {
        this.form.insurance_code = ''
        // Reload all products when no insurance selected
        this.loadProducts()
      }
    }
  },
  async mounted() {
    // Set insurance_code from prop when component is mounted
    if (this.selectedInsuranceCode) {
      this.form.insurance_code = this.selectedInsuranceCode
    }
    // Load products for product name lookup
    await this.loadProducts()
  },
  methods: {
    async saveTemplate() {
      try {
        // Auto-fill value for SUMMARY_TV templates
        let templateValue = (this.form.value || '').trim()
        if (this.form.id && this.form.id.startsWith('SUMMARY_TV') && !templateValue) {
          templateValue = JSON.stringify(["SUMMARY1","SUMMARY2","SUMMARY3","SUMMARY4","SUMMARY5"])
        }
        
        const payload = {
          locale: (this.form.locale || '').trim(),
          id: (this.form.id || '').trim(),
          value: templateValue,
          insurance_code: this.form.insurance_code || this.selectedInsuranceCode || '',
          product_code: (this.form.product_code || '').trim(),
          created_by: 1
        }

        if (!payload.locale || !payload.id || !payload.value) {
          window.showCustomAlert('Mohon lengkapi Locale, Template ID, dan Value', 'error')
          return
        }

        const response = await fetch(`${this.API_URL}/templates/draft/add`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
        
        const result = await response.json()
        if (response.ok) {
          window.showCustomAlert('Template draft saved successfully!', 'success')
          this.resetForm()
          await this.loadDraft()
          // Auto-show draft section after save
          this.showDraftSection = true
        } else {
          window.showCustomAlert(result.error || 'Failed to save template', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    resetForm() {
      this.form = {
        locale: 'id',
        id: '',
        value: '',
        product_code: '',
        insurance_code: this.selectedInsuranceCode || ''
      }
      this.addHtmlConverterForm = {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      }
      this.addCurrentHtmlContent = ''
    },
    
    async convertToHtmlForAdd() {
      try {
        const productInfo = (this.addHtmlConverterForm.productInfo || '').trim()
        
        // Validate required fields
        if (!productInfo) {
          window.showCustomAlert('DETAIL PRODUK harus diisi', 'error')
          return
        }
        
        // Convert textarea to arrays (split by newline, keep empty lines for structure)
        const additionalProtection = this.addHtmlConverterForm.additionalProtection
          ? this.addHtmlConverterForm.additionalProtection.split('\n')
          : []
        const termsConditions = this.addHtmlConverterForm.termsConditions
          ? this.addHtmlConverterForm.termsConditions.split('\n')
          : []
        
        // Prepare data for API
        const requestData = {
          template_type: 'INSURANCE_DETAIL_TV',
          content: {
            product_info: productInfo,
            additional_protection: additionalProtection,
            terms_conditions: termsConditions
          }
        }
        
        const response = await fetch(`${this.API_URL}/html/convert`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestData)
        })
        
        if (response.ok) {
          const result = await response.json()
          // Store HTML content and put it in the value field
          this.addCurrentHtmlContent = result.html_content
          this.form.value = result.html_content
          window.showCustomAlert('HTML berhasil dikonversi dan dimasukkan ke template!', 'success')
        } else {
          const error = await response.json()
          window.showCustomAlert(`Error: ${error.error}`, 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Terjadi kesalahan koneksi ke server API', 'error')
      }
    },
    
    previewHtmlForAdd() {
      // Use addCurrentHtmlContent if available, otherwise use form.value
      const htmlToPreview = this.addCurrentHtmlContent || this.form.value
      
      if (!htmlToPreview || !htmlToPreview.trim()) {
        window.showCustomAlert('Belum ada HTML yang dikonversi. Silakan convert terlebih dahulu.', 'warning')
        return
      }
      
      // Create preview window/modal
      this.showHtmlPreview = true
      
      // Create and show preview modal
      const previewWindow = window.open('', '_blank', 'width=800,height=600,scrollbars=yes')
      if (previewWindow) {
        previewWindow.document.write(`
          <!DOCTYPE html>
          <html>
          <head>
            <title>HTML Preview</title>
            <style>
              body {
                font-family: Arial, sans-serif;
                padding: 20px;
                max-width: 1200px;
                margin: 0 auto;
              }
              .insurance-detail {
                line-height: 1.6;
              }
              .insurance-detail h4 {
                color: #007bff;
                margin-top: 20px;
                margin-bottom: 10px;
                font-weight: bold;
              }
              .insurance-detail ol {
                margin-left: 20px;
                padding-left: 20px;
              }
              .insurance-detail ul {
                margin-left: 20px;
                padding-left: 20px;
                list-style-type: disc;
              }
              .insurance-detail ol li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ul li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ol > li > ul {
                margin-top: 5px;
                margin-bottom: 5px;
              }
              .insurance-detail p {
                margin-bottom: 10px;
              }
            </style>
          </head>
          <body>
            ${htmlToPreview}
          </body>
          </html>
        `)
        previewWindow.document.close()
      } else {
        window.showCustomAlert('Popup blocked. Please allow popups for this site.', 'warning')
      }
    },
    
    async loadDraft() {
      try {
        // Ensure products are loaded before loading draft
        await this.loadProducts()
        
        let url = `${this.API_URL}/templates/draft`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          const draftResponse = await response.json()
          // Backend returns {status: "Success", data: [...]} or direct array
          const draftItems = draftResponse.data || draftResponse || []
          
          // Filter out items with invalid insurance field (array with empty strings)
          // and ensure insurance field is properly formatted if it exists
          this.draftData = (draftItems || []).map(item => {
            // Fix insurance field if it's an array with empty strings ["","",""]
            if (item.insurance && Array.isArray(item.insurance)) {
              // Check if array contains only empty strings
              const isEmptyArray = item.insurance.every(val => !val || String(val).trim() === '')
              if (isEmptyArray) {
                // Replace empty array with proper insurance code string if available
                if (this.selectedInsuranceCode) {
                  item.insurance = this.selectedInsuranceCode
                } else {
                  // Remove insurance field if no selected insurance code
                  delete item.insurance
                }
              } else {
                // If array has valid values, join them
                const validValues = item.insurance.filter(val => val && String(val).trim() !== '')
                if (validValues.length > 0) {
                  item.insurance = validValues.join(', ')
                } else if (this.selectedInsuranceCode) {
                  item.insurance = this.selectedInsuranceCode
                } else {
                  delete item.insurance
                }
              }
            } else if (!item.insurance && item.insurance_code) {
              // Use insurance_code if insurance field doesn't exist
              item.insurance = item.insurance_code
            } else if (!item.insurance && !item.insurance_code && this.selectedInsuranceCode) {
              // Add insurance code if missing but we have selected insurance
              item.insurance = this.selectedInsuranceCode
            }
            return item
          })
          
          this.canConfirm = this.draftData.length > 0
          
          // Auto-show draft section if there's data (after save)
          if (this.draftData.length > 0) {
            this.showDraftSection = true
          }
        } else if (response.status === 404) {
          // Endpoint not found - backend might not implement templates draft yet
          this.draftData = []
          this.canConfirm = false
        } else {
          console.error('Failed to load templates draft:', response.statusText)
          this.draftData = []
          this.canConfirm = false
        }
      } catch (error) {
        console.error('Error loading templates draft:', error)
        this.draftData = []
        this.canConfirm = false
      }
    },
    
    triggerFileInput() {
      this.$refs.fileInput.click()
    },
    
    handleFileChange(event) {
      this.file = event.target.files[0]
      if (this.file) {
        this.uploadFile()
      }
    },
    
    async uploadFile() {
      if (!this.file) {
        window.showCustomAlert('Please select a file', 'warning')
        return
      }
      
      try {
        const formData = new FormData()
        formData.append('file', this.file)
        // Pass selectedInsuranceCode if available, so backend can use it if CSV doesn't have insurance_code
        if (this.selectedInsuranceCode) {
          formData.append('insurance_code', this.selectedInsuranceCode)
        }
        
        const response = await fetch(`${this.API_URL}/templates/upload`, {
          method: 'POST',
          body: formData
        })
        
        const result = await response.json()
        
        if (response.ok) {
          console.log('Upload response:', result)
          // Reload draft data after successful upload
          await this.loadDraft()
          console.log('Draft data after load:', this.draftData)
          this.showDraftSection = true
          this.canConfirm = this.draftData.length > 0
          if (this.draftData.length === 0) {
            window.showCustomAlert('File uploaded but no draft data found. Please check if data was saved correctly.', 'warning')
          } else {
            window.showCustomAlert(result.message || `File uploaded successfully! ${result.success_count || result.data?.length || 0} template(s) added.`, 'success')
          }
          this.$nextTick(() => {
            setTimeout(() => {
              const draftSection = document.querySelector('.templates-section .draft-section')
              if (draftSection) {
                draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
              }
            }, 100)
          })
          // Reset file input
          if (this.$refs.fileInput) {
            this.$refs.fileInput.value = ''
          }
          this.file = null
        } else {
          window.showCustomAlert(result.error || 'Upload failed', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      }
    },
    
    async showDraft() {
      await this.loadDraft()
      this.showDraftSection = true
      this.canConfirm = this.draftData.length > 0
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data available', 'info')
      } else {
        // Scroll to draft section
        this.$nextTick(() => {
          setTimeout(() => {
            const draftSection = document.querySelector('.templates-section .draft-section')
            if (draftSection) {
              draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
            }
          }, 100)
        })
      }
    },
    
    async loadAndShowDraft() {
      // Auto-load and show draft templates (called after products are confirmed)
      await this.loadDraft()
      if (this.draftData.length > 0) {
        this.showDraftSection = true
        this.canConfirm = true
      }
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Are you sure you want to clear all templates draft data?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.API_URL}/templates/draft/clear`, {
          method: 'POST'
        })
        
        if (response.ok) {
          this.draftData = []
          this.showDraftSection = false
          this.canConfirm = false
          window.showCustomAlert('Draft cleared successfully!', 'success')
        } else {
          window.showCustomAlert('Failed to clear draft', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    hideDraft() {
      this.showDraftSection = false
    },
    
    async confirmTemplates() {
      const confirmed = await window.showCustomConfirm('Apakah Anda yakin ingin mengkonfirmasi data draft templates ke database?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.API_URL}/templates/draft/confirm`, {
          method: 'POST'
        })
        
        const result = await response.json()
        
        if (response.ok) {
          window.showCustomAlert('Templates confirmed to database successfully!', 'success')
          this.draftData = []
          this.showDraftSection = false
          this.canConfirm = false
          await this.loadResult()
        } else {
          window.showCustomAlert(result.error || 'Failed to confirm templates', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    async showResult() {
      await this.loadResult()
      this.showResultSection = true
    },
    
    async loadResult() {
      try {
        let url = `${this.API_URL}/templates`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          const responseData = await response.json()
          // Backend returns {status: "Success", data: [...]} or direct array
          const data = responseData.data || responseData || []
          // Data from backend is already JOINed with products
          // Only shows templates where template_id matches summary or insurance_detail
          this.resultData = data
          if (data.length === 0) {
            window.showCustomAlert('Belum ada data templates yang dikonfirmasi (Summary & Insurance Detail)', 'info')
          }
        } else if (response.status === 404 || response.status === 501) {
          // Endpoint not implemented yet
          this.resultData = []
        } else {
          const result = await response.json()
          console.error('Failed to load templates result:', response.statusText)
          window.showCustomAlert(result.error || 'Failed to load templates result', 'error')
        }
      } catch (error) {
        console.error('Error loading templates result:', error)
        // Check if it's a connection error
        if (error.message && (error.message.includes('Failed to fetch') || error.message.includes('ERR_CONNECTION_REFUSED') || error.message.includes('ERR_NAME_NOT_RESOLVED'))) {
          window.showCustomAlert('Connection error. Please ensure backend is running on port 8080', 'error')
        } else {
          window.showCustomAlert('Error loading templates: ' + error.message, 'error')
        }
        this.resultData = []
      }
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    editResult(item) {
      console.log('Editing result item:', item)
      
      // Clear bulk selections when editing single item
      this.selectedDraftItems = []
      this.selectedResultItems = []
      
      // Load result item data into editForm
      this.editForm = {
        id: item.id,
        locale: item.locale || 'id',
        id: item.id || '',
        value: item.value || '',
        product_code: item.product_code || '',
        insurance_code: item.insurance_code || item.insurance || ''
      }
      
      // Store original for comparison (include created_by from item)
      this.originalEditForm = { 
        ...this.editForm,
        created_by: item.created_by || 1
      }
      this.isEditingResult = true
      
      // Reset HTML converter form first
      this.resetHtmlConverterForm()
      
      // Hide result section and show edit form
      this.showEditForm = true
      this.showResultSection = false
      this.showBulkEditForm = false
      
      // Parse HTML to form if it's INSURANCE_DETAIL_TV template
      // Use $nextTick to ensure DOM is ready for computed property
      this.$nextTick(() => {
        if (this.isInsuranceDetailCar && this.editForm.value) {
          this.parseHtmlToForm()
        }
      })
      
      console.log('Edit form loaded with result data:', this.editForm)
    },
    
    async updateTemplate() {
      try {
        const payload = {
          locale: (this.editForm.locale || '').trim(),
          id: (this.editForm.id || '').trim(),
          value: (this.editForm.value || '').trim(),
          created_by: this.originalEditForm.created_by || 1
        }

        if (!payload.locale || !payload.id || !payload.value) {
          window.showCustomAlert('Mohon lengkapi Locale, Template ID, dan Value', 'error')
          return
        }

        // Determine endpoint based on whether editing draft or result
        const url = this.isEditingResult
          ? `${this.API_URL}/templates/${this.editForm.id}`
          : `${this.API_URL}/templates/draft/${this.editForm.id}`

        // For draft updates, include insurance_code and product_code to preserve them
        if (!this.isEditingResult) {
          // Preserve insurance_code from original or current form, or use selected insurance
          // Priority: editForm (if not empty) > originalEditForm > selectedInsuranceCode
          let insuranceCode = ''
          if (this.editForm.insurance_code && this.editForm.insurance_code.trim() !== '') {
            insuranceCode = this.editForm.insurance_code
          } else if (this.originalEditForm && this.originalEditForm.insurance_code && this.originalEditForm.insurance_code.trim() !== '') {
            insuranceCode = this.originalEditForm.insurance_code
          } else if (this.selectedInsuranceCode && this.selectedInsuranceCode.trim() !== '') {
            insuranceCode = this.selectedInsuranceCode
          }
          
          // Preserve product_code from original or current form
          let productCode = ''
          if (this.editForm.product_code && this.editForm.product_code.trim() !== '') {
            productCode = this.editForm.product_code
          } else if (this.originalEditForm && this.originalEditForm.product_code && this.originalEditForm.product_code.trim() !== '') {
            productCode = this.originalEditForm.product_code
          }
          
          // Always include insurance_code and product_code in payload
          // Backend will preserve existing values if empty, but it's better to always send
          payload.insurance_code = insuranceCode || ''
          payload.product_code = productCode || ''
          
          console.log('Preserving insurance_code and product_code:', {
            insurance_code: insuranceCode,
            product_code: productCode,
            editForm_insurance_code: this.editForm.insurance_code,
            originalEditForm_insurance_code: this.originalEditForm?.insurance_code,
            selectedInsuranceCode: this.selectedInsuranceCode,
            payload_insurance_code: payload.insurance_code
          })
        }

        // Log payload for debugging
        console.log('Updating template:', {
          url,
          payload,
          isEditingResult: this.isEditingResult,
          editForm: this.editForm
        })

        if (this.isEditingResult) {
          // Editing RESULT - Save to backend first, then refresh data
          try {
            const response = await fetch(url, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(payload)
            })
            
            if (!response.ok) {
              const result = await response.json()
              console.error('Failed to save to backend:', result.error)
              window.showCustomAlert(result.error || 'Failed to save to backend', 'error')
              return
            }
            
            // Update successful - refresh data immediately
            window.showCustomAlert('Template updated successfully!', 'success')
            await this.loadResult()
            this.cancelEdit()
          } catch (error) {
            console.error('Error saving to backend:', error)
            window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
          }
        } else {
          // Editing DRAFT - Original behavior
          const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          })
          
          const result = await response.json()
          if (response.ok) {
            window.showCustomAlert('Template updated successfully!', 'success')
            this.cancelEdit()
            await this.loadDraft()
          } else {
            window.showCustomAlert(result.error || 'Failed to update', 'error')
          }
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      if (this.isEditingResult) {
        this.showResultSection = true
      } else {
        this.showDraftSection = true
      }
      this.isEditingResult = false
      this.editForm = {
        id: null,
        locale: '',
        id: '',
        value: '',
        product_code: '',
        insurance_code: ''
      }
      this.originalEditForm = null
      this.resetHtmlConverterForm()
      // Clear selections when canceling edit
      this.selectedDraftItems = []
      this.selectedResultItems = []
    },
    
    editDraft(item) {
      console.log('Editing draft item:', item)
      
      // Clear bulk selections when editing single item
      this.selectedDraftItems = []
      this.selectedResultItems = []
      
      // Load draft item data into editForm
      // Extract insurance_code from item (can be in insurance_code, insurance field, or as string/array)
      let insuranceCode = ''
      if (item.insurance_code && item.insurance_code.trim() !== '') {
        insuranceCode = item.insurance_code
      } else if (item.insurance) {
        if (Array.isArray(item.insurance)) {
          // If insurance is array, get first non-empty value
          const validValue = item.insurance.find(val => val && String(val).trim() !== '')
          insuranceCode = validValue || ''
        } else if (typeof item.insurance === 'string' && item.insurance.trim() !== '') {
          insuranceCode = item.insurance
        }
      }
      // Fallback to selectedInsuranceCode if still empty
      if (!insuranceCode && this.selectedInsuranceCode) {
        insuranceCode = this.selectedInsuranceCode
      }
      
      this.editForm = {
        id: item.id || '',
        locale: item.locale || 'id',
        value: item.value || '',
        product_code: item.product_code || '',
        insurance_code: insuranceCode || this.selectedInsuranceCode || ''
      }
      
      // Ensure insurance_code is always set
      if (!this.editForm.insurance_code && this.selectedInsuranceCode) {
        this.editForm.insurance_code = this.selectedInsuranceCode
      }
      
      // Store original for comparison (include created_by from item)
      this.originalEditForm = { 
        ...this.editForm,
        created_by: item.created_by || 1
      }
      
      console.log('editDraft - loaded insurance_code:', {
        insurance_code: insuranceCode,
        item_insurance_code: item.insurance_code,
        item_insurance: item.insurance,
        selectedInsuranceCode: this.selectedInsuranceCode
      })
      this.isEditingResult = false
      
      // Reset HTML converter form first
      this.resetHtmlConverterForm()
      
      // Hide draft section and show edit form
      this.showEditForm = true
      this.showDraftSection = false
      this.showBulkEditForm = false
      
      // Parse HTML to form if it's INSURANCE_DETAIL_TV template
      // Use $nextTick to ensure DOM is ready for computed property
      this.$nextTick(() => {
        if (this.isInsuranceDetailCar && this.editForm.value) {
          this.parseHtmlToForm()
        }
      })
      
      console.log('Edit form loaded with draft data:', this.editForm)
    },
    
    async deleteDraft(index) {
      const item = this.draftData[index]
      if (!item || !item.id) {
        window.showCustomAlert('Invalid item to delete', 'error')
        return
      }
      
      const confirmed = await window.showCustomConfirm('Apakah Anda yakin ingin menghapus data ini?')
      if (!confirmed) return
      
      try {
        // Get locale from item, default to 'id'
        const locale = item.locale || 'id'
        const url = `${this.API_URL}/templates/draft/${item.id}?locale=${locale}`
        
        const response = await fetch(url, {
          method: 'DELETE'
        })
        
        if (response.ok) {
          window.showCustomAlert('Data berhasil dihapus!', 'success')
          await this.loadDraft()
          if (this.draftData.length === 0) {
            this.showDraftSection = false
            this.canConfirm = false
          }
        } else {
          const errorData = await response.json().catch(() => ({ error: 'Failed to delete draft' }))
          window.showCustomAlert(errorData.error || 'Gagal menghapus data', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    formatTimestamp(timestamp) {
      if (!timestamp) return '-'
      try {
        return new Date(timestamp).toLocaleString()
      } catch (e) {
        return timestamp
      }
    },
    
    async openGenerateTemplatesModal() {
      this.showGenerateTemplatesModal = true
      this.selectedProductsForGeneration = []
      this.productSearchQuery = ''
      await this.loadProductsForGeneration()
    },
    
    closeGenerateTemplatesModal() {
      this.showGenerateTemplatesModal = false
      this.selectedProductsForGeneration = []
      this.productSearchQuery = ''
    },
    
    async loadProductsForGeneration() {
      try {
        let url = `${this.API_URL}/products`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          this.productsForGeneration = await response.json()
        } else {
          window.showCustomAlert('Failed to load products', 'error')
        }
      } catch (error) {
        console.error('Error loading products:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      }
    },
    
    toggleSelectAllProducts() {
      if (this.allProductsSelected) {
        this.selectedProductsForGeneration = []
      } else {
        this.selectedProductsForGeneration = this.productsForGeneration.map(p => p.code)
      }
    },
    
    toggleProductSelection(productCode) {
      const index = this.selectedProductsForGeneration.indexOf(productCode)
      if (index > -1) {
        this.selectedProductsForGeneration.splice(index, 1)
      } else {
        this.selectedProductsForGeneration.push(productCode)
      }
    },
    
    async generateTemplatesFromSelectedProducts() {
      if (this.selectedProductsForGeneration.length === 0) {
        window.showCustomAlert('Please select at least one product', 'warning')
        return
      }
      
      const confirmed = await window.showCustomConfirm(
        `Generate template drafts (SUMMARY and INSURANCE_DETAIL) for ${this.selectedProductsForGeneration.length} selected product(s)?`
      )
      if (!confirmed) return
      
      this.isGenerating = true
      try {
        const response = await fetch(`${this.API_URL}/products/generate-templates/selected`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            product_codes: this.selectedProductsForGeneration,
            created_by: 1
          })
        })
        
        if (response.ok) {
          const result = await response.json()
          window.showCustomAlert(
            `Success! Generated templates for ${result.products_processed} product(s).`, 
            'success'
          )
          // Close modal and reload draft data
          this.closeGenerateTemplatesModal()
          await this.loadDraft()
          this.showDraftSection = true
        } else {
          const error = await response.json()
          window.showCustomAlert(error.error || 'Failed to generate templates', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      } finally {
        this.isGenerating = false
      }
    },
    
    resetHtmlConverterForm() {
      this.htmlConverterForm = {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      }
      this.currentHtmlContent = ''
      this.showHtmlPreview = false
    },
    
    parseNestedList(olElement) {
      // Parse ordered list with nested unordered lists
      const items = []
      const liElements = olElement.querySelectorAll('> li')
      
      liElements.forEach(li => {
        // Clone to avoid modifying original
        const clonedLi = li.cloneNode(true)
        const ul = clonedLi.querySelector('ul')
        
        // Remove nested ul to get main text
        if (ul) {
          ul.remove()
        }
        
        const mainText = clonedLi.textContent.trim()
        if (mainText) {
          items.push(mainText)
        }
        
        // Get nested bullet items
        if (ul) {
          const bulletItems = ul.querySelectorAll('li')
          bulletItems.forEach(bulletLi => {
            const bulletText = bulletLi.textContent.trim()
            if (bulletText) {
              items.push('  - ' + bulletText)
            }
          })
        }
      })
      
      return items
    },
    
    parseHtmlToForm() {
      // Parse HTML from editForm.value and populate htmlConverterForm
      // Only works for INSURANCE_DETAIL_TV template
      if (!this.isInsuranceDetailCar || !this.editForm.value) {
        return
      }
      
      try {
        const htmlContent = this.editForm.value.trim()
        if (!htmlContent) {
          return
        }
        
        // Create a temporary DOM element to parse HTML
        const tempDiv = document.createElement('div')
        tempDiv.innerHTML = htmlContent
        
        // Extract Product Info
        const productInfoSection = tempDiv.querySelector('.product-info-section')
        if (productInfoSection) {
          const pTag = productInfoSection.querySelector('p')
          if (pTag) {
            this.htmlConverterForm.productInfo = pTag.textContent.trim()
          }
        }
        
        // Extract Additional Protection (PENGECUALIAN UMUM) - with nested bullets
        const additionalProtectionSection = tempDiv.querySelector('.additional-protection-section')
        if (additionalProtectionSection) {
          const olTag = additionalProtectionSection.querySelector('ol')
          if (olTag) {
            const items = this.parseNestedList(olTag)
            this.htmlConverterForm.additionalProtection = items.join('\n')
          }
        }
        
        // Extract Terms and Conditions (SYARAT DAN KETENTUAN) - with nested bullets
        const termsConditionsSection = tempDiv.querySelector('.terms-conditions-section')
        if (termsConditionsSection) {
          const olTag = termsConditionsSection.querySelector('ol')
          if (olTag) {
            const items = this.parseNestedList(olTag)
            this.htmlConverterForm.termsConditions = items.join('\n')
          }
        }
        
        // Store current HTML content for preview
        this.currentHtmlContent = htmlContent
        
        console.log('HTML parsed and form populated:', this.htmlConverterForm)
      } catch (error) {
        console.error('Error parsing HTML:', error)
        // If parsing fails, just reset the form
        this.resetHtmlConverterForm()
      }
    },
    
    parseHtmlToBulkEditForm(htmlContent) {
      // Parse HTML and populate bulkEditHtmlConverterForm
      // Only works for INSURANCE_DETAIL_TV template
      if (!htmlContent) {
        return
      }
      
      // Check if all selected items are INSURANCE_DETAIL_TV
      const itemsToCheck = this.isBulkEditingDraft ? this.selectedDraftItems : this.selectedResultItems
      if (itemsToCheck.length === 0) return false
      const allInsuranceDetail = itemsToCheck.every(item => {
        const templateId = item.item?.id || item.id
        return templateId && templateId.includes('INSURANCE_DETAIL_TV')
      })
      
      if (!allInsuranceDetail) {
        return
      }
      
      try {
        const html = htmlContent.trim()
        if (!html) {
          return
        }
        
        // Create a temporary DOM element to parse HTML
        const tempDiv = document.createElement('div')
        tempDiv.innerHTML = html
        
        // Extract Product Info
        const productInfoSection = tempDiv.querySelector('.product-info-section')
        if (productInfoSection) {
          const pTag = productInfoSection.querySelector('p')
          if (pTag) {
            this.bulkEditHtmlConverterForm.productInfo = pTag.textContent.trim()
          }
        }
        
        // Extract Additional Protection (PENGECUALIAN UMUM) - with nested bullets
        const additionalProtectionSection = tempDiv.querySelector('.additional-protection-section')
        if (additionalProtectionSection) {
          const olTag = additionalProtectionSection.querySelector('ol')
          if (olTag) {
            const items = this.parseNestedList(olTag)
            this.bulkEditHtmlConverterForm.additionalProtection = items.join('\n')
          }
        }
        
        // Extract Terms and Conditions (SYARAT DAN KETENTUAN) - with nested bullets
        const termsConditionsSection = tempDiv.querySelector('.terms-conditions-section')
        if (termsConditionsSection) {
          const olTag = termsConditionsSection.querySelector('ol')
          if (olTag) {
            const items = this.parseNestedList(olTag)
            this.bulkEditHtmlConverterForm.termsConditions = items.join('\n')
          }
        }
        
        // Store current HTML content for preview
        this.bulkEditCurrentHtmlContent = html
        
        console.log('HTML parsed and bulk edit form populated:', this.bulkEditHtmlConverterForm)
      } catch (error) {
        console.error('Error parsing HTML for bulk edit:', error)
        // If parsing fails, just reset the form
        this.bulkEditHtmlConverterForm = {
          productInfo: '',
          additionalProtection: '',
          termsConditions: ''
        }
        this.bulkEditCurrentHtmlContent = ''
      }
    },
    
    async convertToHtmlForBulkEdit() {
      try {
        const productInfo = (this.bulkEditHtmlConverterForm.productInfo || '').trim()
        
        // Validate required fields
        if (!productInfo) {
          window.showCustomAlert('DETAIL PRODUK harus diisi', 'error')
          return
        }
        
        // Convert textarea to arrays (split by newline, keep structure for nested bullets)
        const additionalProtection = this.bulkEditHtmlConverterForm.additionalProtection
          ? this.bulkEditHtmlConverterForm.additionalProtection.split('\n')
          : []
        const termsConditions = this.bulkEditHtmlConverterForm.termsConditions
          ? this.bulkEditHtmlConverterForm.termsConditions.split('\n')
          : []
        
        // Prepare data for API
        const requestData = {
          template_type: 'INSURANCE_DETAIL_TV',
          content: {
            product_info: productInfo,
            additional_protection: additionalProtection,
            terms_conditions: termsConditions
          }
        }
        
        const response = await fetch(`${this.API_URL}/html/convert`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestData)
        })
        
        if (response.ok) {
          const result = await response.json()
          // Store HTML content and put it in the value field
          this.bulkEditCurrentHtmlContent = result.html_content
          this.bulkEditForm.value = result.html_content
          window.showCustomAlert('HTML berhasil dikonversi dan dimasukkan ke template!', 'success')
        } else {
          const error = await response.json()
          window.showCustomAlert(`Error: ${error.error}`, 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Terjadi kesalahan koneksi ke server API', 'error')
      }
    },
    
    previewHtmlForBulkEdit() {
      // Use bulkEditCurrentHtmlContent if available, otherwise use bulkEditForm.value
      const htmlToPreview = this.bulkEditCurrentHtmlContent || this.bulkEditForm.value
      
      if (!htmlToPreview || !htmlToPreview.trim()) {
        window.showCustomAlert('Belum ada HTML yang dikonversi. Silakan convert terlebih dahulu.', 'warning')
        return
      }
      
      // Create preview window/modal
      this.showHtmlPreview = true
      
      // Create and show preview modal
      const previewWindow = window.open('', '_blank', 'width=800,height=600,scrollbars=yes')
      if (previewWindow) {
        previewWindow.document.write(`
          <!DOCTYPE html>
          <html>
          <head>
            <title>HTML Preview</title>
            <style>
              body {
                font-family: Arial, sans-serif;
                padding: 20px;
                max-width: 1200px;
                margin: 0 auto;
              }
              .insurance-detail {
                line-height: 1.6;
              }
              .insurance-detail h4 {
                color: #007bff;
                margin-top: 20px;
                margin-bottom: 10px;
                font-weight: bold;
              }
              .insurance-detail ol {
                margin-left: 20px;
                padding-left: 20px;
              }
              .insurance-detail ul {
                margin-left: 20px;
                padding-left: 20px;
                list-style-type: disc;
              }
              .insurance-detail ol li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ul li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ol > li > ul {
                margin-top: 5px;
                margin-bottom: 5px;
              }
              .insurance-detail p {
                margin-bottom: 10px;
              }
            </style>
          </head>
          <body>
            ${htmlToPreview}
          </body>
          </html>
        `)
        previewWindow.document.close()
      } else {
        window.showCustomAlert('Popup blocked. Please allow popups for this site.', 'warning')
      }
    },
    
    async convertToHtml() {
      try {
        const productInfo = (this.htmlConverterForm.productInfo || '').trim()
        
        // Validate required fields
        if (!productInfo) {
          window.showCustomAlert('DETAIL PRODUK harus diisi', 'error')
          return
        }
        
        // Convert textarea to arrays (split by newline, keep structure for nested bullets)
        const additionalProtection = this.htmlConverterForm.additionalProtection
          ? this.htmlConverterForm.additionalProtection.split('\n')
          : []
        const termsConditions = this.htmlConverterForm.termsConditions
          ? this.htmlConverterForm.termsConditions.split('\n')
          : []
        
        // Prepare data for API
        const requestData = {
          template_type: 'INSURANCE_DETAIL_TV',
          content: {
            product_info: productInfo,
            additional_protection: additionalProtection,
            terms_conditions: termsConditions
          }
        }
        
        const response = await fetch(`${this.API_URL}/html/convert`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestData)
        })
        
        if (response.ok) {
          const result = await response.json()
          // Store HTML content and put it in the value field
          this.currentHtmlContent = result.html_content
          this.editForm.value = result.html_content
          window.showCustomAlert('HTML berhasil dikonversi dan dimasukkan ke template!', 'success')
        } else {
          const error = await response.json()
          window.showCustomAlert(`Error: ${error.error}`, 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Terjadi kesalahan koneksi ke server API', 'error')
      }
    },
    
    previewHtml() {
      // Use currentHtmlContent if available, otherwise use editForm.value
      const htmlToPreview = this.currentHtmlContent || this.editForm.value
      
      if (!htmlToPreview || !htmlToPreview.trim()) {
        window.showCustomAlert('Belum ada HTML yang dikonversi. Silakan convert terlebih dahulu.', 'warning')
        return
      }
      
      // Create preview window/modal
      this.showHtmlPreview = true
      
      // Create and show preview modal
      const previewWindow = window.open('', '_blank', 'width=800,height=600,scrollbars=yes')
      if (previewWindow) {
        previewWindow.document.write(`
          <!DOCTYPE html>
          <html>
          <head>
            <title>HTML Preview</title>
            <style>
              body {
                font-family: Arial, sans-serif;
                padding: 20px;
                max-width: 1200px;
                margin: 0 auto;
              }
              .insurance-detail {
                line-height: 1.6;
              }
              .insurance-detail h4 {
                color: #007bff;
                margin-top: 20px;
                margin-bottom: 10px;
                font-weight: bold;
              }
              .insurance-detail ol {
                margin-left: 20px;
                padding-left: 20px;
              }
              .insurance-detail ul {
                margin-left: 20px;
                padding-left: 20px;
                list-style-type: disc;
              }
              .insurance-detail ol li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ul li {
                margin-bottom: 8px;
                padding-left: 5px;
              }
              .insurance-detail ol > li > ul {
                margin-top: 5px;
                margin-bottom: 5px;
              }
              .insurance-detail p {
                margin-bottom: 10px;
              }
            </style>
          </head>
          <body>
            ${htmlToPreview}
          </body>
          </html>
        `)
        previewWindow.document.close()
      } else {
        window.showCustomAlert('Popup blocked. Please allow popups for this site.', 'warning')
      }
    },
    
    formatDateOnly(dateString) {
      if (!dateString) return '-'
      const date = new Date(dateString)
      if (isNaN(date.getTime())) {
        if (typeof dateString === 'string') {
          const datePart = dateString.split(' ')[0]
          return datePart || dateString
        }
        return dateString
      }
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    },
    
    async loadProducts() {
      try {
        let url = 'http://localhost:8080/api/products'
        // Filter by selectedInsuranceCode if available
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const response = await fetch(url)
        if (response.ok) {
          this.products = await response.json()
          console.log('Products loaded:', this.products.length)
        }
      } catch (error) {
        console.error('Error loading products:', error)
      }
    },
    
    extractProductCodeFromTemplateId(templateId) {
      if (!templateId) return null
      // Template ID format: 
      // - SUMMARY_TV_DAMAI_WW_IND_01
      // - INSURANCE_DETAIL_TV_DAMAI_WW_IND_01
      // - SUMMARY_TV_DAMAI_DOM_IND_ANN_03 (ANNUAL with suffix)
      // Product code format: TV-DAMAI-WW-IND-01 (with dashes, not underscores)
      // So we need to extract and convert underscores to dashes
      const parts = templateId.split('_')
      if (parts.length >= 2) {
        // Find index of TV
        const tvIndex = parts.findIndex(p => p === 'TV')
        if (tvIndex !== -1 && tvIndex < parts.length - 1) {
          // Get everything after TV and convert underscores to dashes
          const codeParts = parts.slice(tvIndex + 1)
          // Convert to product code format: TV-DAMAI-WW-IND-01
          let productCode = 'TV-' + codeParts.join('-')
          
          // Try to find exact match first
          let product = this.products.find(p => p.code === productCode)
          if (product) {
            return productCode
          }
          
          // If not found, try without the last part (suffix like _03, _04)
          // This handles cases like INSURANCE_DETAIL_TV_DAMAI_DOM_IND_ANN_03
          // where product code might be TV-DAMAI-DOM-IND-ANN-01 or TV-DAMAI-DOM-IND-ANN-02
          if (codeParts.length > 0) {
            // Remove last part and try again
            const codePartsWithoutSuffix = codeParts.slice(0, -1)
            if (codePartsWithoutSuffix.length > 0) {
              const baseProductCode = 'TV-' + codePartsWithoutSuffix.join('-')
              // Try to find product code that starts with baseProductCode
              product = this.products.find(p => p.code && p.code.startsWith(baseProductCode))
              if (product) {
                return product.code
              }
            }
            
            // If still not found, try to match by removing numeric suffix
            // For example: INSURANCE_DETAIL_TV_DAMAI_DOM_IND_ANN_03 -> TV-DAMAI-DOM-IND-ANN-01
            // Check if last part is a number
            const lastPart = codeParts[codeParts.length - 1]
            if (lastPart && /^\d+$/.test(lastPart)) {
              // Remove numeric suffix and try to find matching product
              const codePartsWithoutNumeric = codeParts.slice(0, -1)
              if (codePartsWithoutNumeric.length > 0) {
                const baseCode = 'TV-' + codePartsWithoutNumeric.join('-')
                // Find all products that start with baseCode
                const matchingProducts = this.products.filter(p => p.code && p.code.startsWith(baseCode))
                if (matchingProducts.length > 0) {
                  // Return the first matching product code
                  return matchingProducts[0].code
                }
              }
            }
          }
          
          // Return the extracted code anyway, even if not found
          return productCode
        }
      }
      return null
    },
    
    getProductNameFromTemplateId(templateId) {
      const productCode = this.extractProductCodeFromTemplateId(templateId)
      if (!productCode) return null
      
      // Try exact match first
      let product = this.products.find(p => p.code === productCode)
      if (product) {
        return product.name
      }
      
      // If not found, try to find product code that starts with the extracted code
      // This handles cases where template ID has suffix but product code doesn't
      product = this.products.find(p => p.code && p.code.startsWith(productCode))
      if (product) {
        return product.name
      }
      
      return null
    },
    
    getProductNameFromProductCode(productCode) {
      if (!productCode) return null
      const product = this.products.find(p => p.code === productCode)
      return product ? product.name : null
    },
    
    getInsuranceCodeFromTemplateId(templateId) {
      const productCode = this.extractProductCodeFromTemplateId(templateId)
      if (!productCode) return null
      
      const product = this.products.find(p => p.code === productCode)
      return product ? product.insurance_code : null
    },
    
    // Bulk edit methods
    isDraftSelected(item) {
      return this.selectedDraftItems.some(selected => 
        selected.id === item.id && selected.locale === item.locale
      )
    },
    
    isResultSelected(item) {
      return this.selectedResultItems.some(selected => 
        selected.id === item.id && selected.locale === item.locale
      )
    },
    
    toggleDraftSelection(item) {
      const index = this.selectedDraftItems.findIndex(selected => 
        selected.id === item.id && selected.locale === item.locale
      )
      if (index > -1) {
        this.selectedDraftItems.splice(index, 1)
      } else {
        this.selectedDraftItems.push({ id: item.id, locale: item.locale, item: item })
      }
    },
    
    toggleResultSelection(item) {
      const index = this.selectedResultItems.findIndex(selected => 
        selected.id === item.id && selected.locale === item.locale
      )
      if (index > -1) {
        this.selectedResultItems.splice(index, 1)
      } else {
        this.selectedResultItems.push({ id: item.id, locale: item.locale, item: item })
      }
    },
    
    toggleSelectAllDraft(event) {
      if (event.target.checked) {
        this.selectedDraftItems = this.draftData.map(item => ({
          id: item.id,
          locale: item.locale,
          item: item
        }))
      } else {
        this.selectedDraftItems = []
      }
    },
    
    toggleSelectAllResult(event) {
      if (event.target.checked) {
        this.selectedResultItems = this.resultData.map(item => ({
          id: item.id,
          locale: item.locale,
          item: item
        }))
      } else {
        this.selectedResultItems = []
      }
    },
    
    openBulkEdit(isDraft) {
      this.isBulkEditingDraft = isDraft
      this.showBulkEditForm = true
      this.bulkEditForm.value = ''
      this.bulkEditHtmlConverterForm = {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      }
      this.bulkEditCurrentHtmlContent = ''
      
      // Hide draft/result section
      if (isDraft) {
        this.showDraftSection = false
      } else {
        this.showResultSection = false
      }
      
      // Parse HTML from first selected item if it's INSURANCE_DETAIL_TV
      this.$nextTick(() => {
        const itemsToCheck = isDraft ? this.selectedDraftItems : this.selectedResultItems
        if (itemsToCheck.length > 0) {
          // Check if all selected items are INSURANCE_DETAIL_TV
          const allInsuranceDetail = itemsToCheck.every(item => {
            const templateId = (item.item?.id || item.id)
            return templateId && templateId.includes('INSURANCE_DETAIL_TV')
          })
          
          if (allInsuranceDetail) {
            const firstItem = itemsToCheck[0].item || itemsToCheck[0]
            if (firstItem.value) {
              this.parseHtmlToBulkEditForm(firstItem.value)
            }
          }
        }
      })
    },
    
    cancelBulkEdit() {
      this.showBulkEditForm = false
      this.bulkEditForm.value = ''
      this.bulkEditHtmlConverterForm = {
        productInfo: '',
        additionalProtection: '',
        termsConditions: ''
      }
      this.bulkEditCurrentHtmlContent = ''
      this.selectedDraftItems = []
      this.selectedResultItems = []
      
      // Show draft/result section back
      if (this.isBulkEditingDraft) {
        this.showDraftSection = true
      } else {
        this.showResultSection = true
      }
    },
    
    async updateBulkTemplates() {
      if (!this.bulkEditForm.value || !this.bulkEditForm.value.trim()) {
        window.showCustomAlert('Value is required', 'error')
        return
      }
      
      const itemsToUpdate = this.isBulkEditingDraft ? this.selectedDraftItems : this.selectedResultItems
      
      if (itemsToUpdate.length === 0) {
        window.showCustomAlert('No items selected', 'error')
        return
      }
      
      const confirmed = await window.showCustomConfirm(
        `Are you sure you want to update ${itemsToUpdate.length} template(s) with the new value?`
      )
      if (!confirmed) return
      
      try {
        let successCount = 0
        let errorCount = 0
        
        // Update each template
        for (const selectedItem of itemsToUpdate) {
          const item = selectedItem.item || selectedItem
          const payload = {
            locale: (item.locale || 'id').trim(),
            id: (item.id || '').trim(),
            value: this.bulkEditForm.value.trim(),
            created_by: item.created_by || 1
          }
          
          if (!payload.locale || !payload.id) {
            errorCount++
            continue
          }
          
          try {
            const url = this.isBulkEditingDraft
              ? `${this.API_URL}/templates/draft/${item.id}`
              : `${this.API_URL}/templates/${item.id}`
            
            const method = 'PUT'
            
            const response = await fetch(url, {
              method: method,
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(payload)
            })
            
            if (response.ok) {
              successCount++
            } else {
              const result = await response.json()
              console.error(`Failed to update ${item.id}:`, result.error)
              errorCount++
            }
          } catch (error) {
            console.error(`Error updating ${item.id}:`, error)
            errorCount++
          }
        }
        
        if (successCount > 0) {
          window.showCustomAlert(
            `Successfully updated ${successCount} template(s)${errorCount > 0 ? `. ${errorCount} failed.` : ''}`,
            successCount === itemsToUpdate.length ? 'success' : 'warning'
          )
          
          // Refresh data
          if (this.isBulkEditingDraft) {
            await this.loadDraft()
            this.showDraftSection = true
          } else {
            await this.loadResult()
            this.showResultSection = true
          }
          
          // Clear selections and close form
          this.cancelBulkEdit()
        } else {
          window.showCustomAlert('Failed to update templates', 'error')
        }
      } catch (error) {
        console.error('Error in bulk update:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    }
  }
}
</script>

<style scoped>
.templates-section {
  position: relative;
  border-left: 4px solid #28a745;
  padding-left: 1.5rem;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 1.5rem;
}

.header-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #28a745;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.header-icon i {
  color: white;
  font-size: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.add-form-section {
  border-left: 4px solid #17a2b8;
  margin-top: 2rem;
}

.edit-form-section {
  border-left: 4px solid #ffc107;
  margin-top: 2rem;
}

.form-grid {
  display: grid;
  gap: 1rem;
  margin-bottom: 1rem;
  align-items: start;
}

.form-grid > div {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.form-grid-1 {
  grid-template-columns: 1fr;
}

.form-grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.form-label {
  font-weight: bold;
  color: #333;
  margin-bottom: 0.5rem;
  font-size: 14px;
  min-height: 20px;
  line-height: 1.4;
  display: block;
}

.form-control,
.form-select,
textarea.form-control {
  width: 100%;
  border-radius: 0.4rem;
  font-size: 0.85rem;
  box-shadow: none;
  border: 1px solid #ced4da;
  padding: 0.375rem 0.75rem;
}

.form-control {
  height: 40px;
}

textarea.form-control {
  height: auto;
  resize: vertical;
}

.form-control:focus,
.form-select:focus,
textarea.form-control:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.form-control-readonly {
  background-color: #e9ecef;
  cursor: not-allowed;
  opacity: 0.7;
}

.form-control[readonly] {
  background-color: #e9ecef;
  cursor: not-allowed;
}

@media (max-width: 1200px) {
  .form-grid-3 {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .form-grid-1,
  .form-grid-3 {
    grid-template-columns: 1fr;
  }
}

.table-section {
  margin-top: 2rem;
  border: 1px solid #dee2e6;
  border-radius: 0.5rem;
  background-color: #fff;
  overflow: hidden;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.table-header h4 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.table-content {
  padding: 1rem;
}

.result-section {
  border-left: 4px solid #007bff;
}

.result-section .table-header {
  background-color: #e7f3ff;
}

.draft-section {
  border-left: 4px solid #28a745;
}

.draft-section .table-header {
  background-color: #e8f5e9;
}

.badge {
  padding: 0.35em 0.65em;
  font-size: 0.875em;
  font-weight: 600;
}

.bg-info {
  background-color: #17a2b8 !important;
}

.bg-primary {
  background-color: #007bff !important;
}

.html-converter-section {
  margin: 1rem 0;
}

.html-converter-card {
  border: 1px solid #007bff;
  border-radius: 0.5rem;
  overflow: hidden;
  background-color: #fff;
}

.html-converter-header {
  background-color: #007bff;
  color: white;
  padding: 0.75rem 1rem;
  font-weight: 600;
}

.html-converter-header h6 {
  margin: 0;
  color: white;
  display: flex;
  align-items: center;
}

.html-converter-body {
  padding: 1rem;
}

.html-converter-body .btn {
  font-size: 0.875rem;
}

.text-info {
  color: #17a2b8 !important;
}

/* Generate Templates Modal Styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(2px);
  z-index: 1040;
  animation: fadeIn 0.2s ease;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1050;
  padding: 1rem;
  overflow-y: auto;
}

.generate-templates-modal {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 650px;
  max-width: 95%;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.generate-templates-modal .modal-header {
  padding: 1rem 1.5rem;
  border-bottom: 2px solid #17a2b8;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  background: #17a2b8;
  color: white;
}

.modal-header-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.modal-header-icon {
  font-size: 1.25rem;
  color: white;
}

.generate-templates-modal .modal-title {
  margin: 0;
  font-weight: 600;
  font-size: 1.1rem;
  color: white;
}

.btn-close-modal {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  color: white;
  padding: 0.4rem;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.btn-close-modal:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.generate-templates-modal .modal-body {
  padding: 1.25rem;
  overflow-y: auto;
  flex: 1;
  background: #f8f9fa;
}

.modal-section-header {
  margin-bottom: 1rem;
}

.form-label-modal {
  font-weight: 600;
  font-size: 0.9rem;
  color: #333;
  margin: 0;
  display: flex;
  align-items: center;
}

.btn-select-all {
  border-radius: 6px;
  padding: 0.4rem 0.75rem;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.btn-select-all:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(23, 162, 184, 0.2);
}

.search-wrapper {
  position: relative;
  margin-bottom: 1rem;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: #6c757d;
  z-index: 1;
  font-size: 0.85rem;
}

.search-input {
  padding-left: 2.25rem;
  border-radius: 6px;
  border: 1px solid #ced4da;
  height: 38px;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.search-input:focus {
  border-color: #17a2b8;
  box-shadow: 0 0 0 0.15rem rgba(23, 162, 184, 0.15);
}

.products-list-container {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 0.5rem;
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 1rem;
}

.products-list-container::-webkit-scrollbar {
  width: 6px;
}

.products-list-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

.products-list-container::-webkit-scrollbar-thumb {
  background: #17a2b8;
  border-radius: 10px;
}

.products-list-container::-webkit-scrollbar-thumb:hover {
  background: #138496;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #6c757d;
}

.empty-icon {
  font-size: 3rem;
  color: #dee2e6;
  margin-bottom: 1rem;
}

.empty-text {
  margin: 0;
  font-size: 0.95rem;
}

/* Compact List Style */
.product-item-list {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  margin-bottom: 0.5rem;
  cursor: pointer;
  transition: all 0.15s ease;
  background: white;
}

.product-item-list:hover {
  border-color: #17a2b8;
  background: #f0f9fa;
}

.product-item-list-selected {
  border-color: #17a2b8;
  background: #e6f7f9;
}

.product-checkbox-list {
  width: 18px;
  height: 18px;
  cursor: pointer;
  margin: 0;
  flex-shrink: 0;
  accent-color: #17a2b8;
}

.product-label-list {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.4;
}

.product-code-list {
  font-weight: 600;
  color: #333;
  min-width: 180px;
  font-size: 0.85rem;
}

.product-name-list {
  color: #555;
  flex: 1;
  font-size: 0.85rem;
}

.product-badge-list {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.2rem 0.5rem;
  background: #17a2b8;
  color: white;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
}

.product-badge-list i {
  font-size: 0.7rem;
}

.selection-summary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(23, 162, 184, 0.1);
  border: 1px solid #17a2b8;
  border-radius: 6px;
  margin-top: 0.75rem;
}

.summary-icon {
  font-size: 1.1rem;
  color: #28a745;
}

.summary-text {
  font-size: 0.875rem;
  color: #333;
}

.generate-templates-modal .modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #dee2e6;
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  flex-shrink: 0;
  background: white;
}

.btn-modal {
  padding: 0.5rem 1.25rem;
  border-radius: 6px;
  font-weight: 500;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  min-width: 100px;
}

.btn-modal:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-modal:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .generate-templates-modal {
    width: 100%;
    max-width: 100%;
    max-height: 100vh;
    border-radius: 0;
  }
  
  .generate-templates-modal .modal-header {
    padding: 1.25rem 1.5rem;
  }
  
  .generate-templates-modal .modal-body {
    padding: 1.5rem;
  }
  
  .product-item {
    padding: 0.75rem;
  }
}
</style>

