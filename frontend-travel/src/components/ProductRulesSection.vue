<template>
  <div class="col-12">
    <div class="form-section">
      <h3><i class="fas fa-file-excel me-1"></i> Product Rules Management</h3>
      
      <div class="row">
        <div class="col-md-12">
          <div class="d-flex gap-2 mb-3">
            <button class="btn btn-success" @click="downloadTemplate">
              <i class="fas fa-download"></i> Download Template CSV
            </button>
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
        </div>
      </div>
    </div>

    <!-- Draft Table Section -->
    <div v-if="showDraftSection && !showEditForm" class="table-section draft-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📊 Data Preview: Product Rules (Draft Save)</h4>
        <div>
          <button class="btn btn-success btn-sm me-2" @click="saveRules" :disabled="!canConfirm">
            <i class="fas fa-check"></i> Konfirmasi Data (DB)
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
                <th>Product Code</th>
                <th>Premium Type</th>
                <th>Start Days</th>
                <th>End Days</th>
                <th>Base Premium Value</th>
                <th>Rules</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="draftData.length === 0">
                <td colspan="7" class="text-center text-muted py-4">
                  <i class="fas fa-info-circle me-2"></i>No draft data available
                </td>
              </tr>
              <tr v-for="(item, index) in draftData" :key="index">
                <td>{{ item.product_code }}</td>
                <td>{{ item.premium_type }}</td>
                <td>{{ item.start_days }}</td>
                <td>{{ item.end_days }}</td>
                <td>{{ item.base_premium_value }}</td>
                <td>
                  <div style="max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.rules">
                    {{ item.rules }}
                  </div>
                </td>
                <td>
                  <button class="btn btn-sm btn-warning me-1" @click="editDraft(index)">
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

    <!-- Edit Draft Form Section -->
    <div v-if="showEditForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Draft Data</h3>
      
      <form @submit.prevent="updateDraft">
        <div class="form-grid form-grid-3">
          <div>
            <label class="form-label">Product Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editDraftForm.product_code" required>
          </div>
          <div>
            <label class="form-label">Premium Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editDraftForm.premium_type" required>
              <option value="">Select premium type...</option>
              <option value="DAILY">DAILY</option>
              <option value="ANNUAL">ANNUAL</option>
              <option value="PER_EXTRA_SEVEN_DAYS">PER_EXTRA_SEVEN_DAYS</option>
            </select>
          </div>
          <div>
            <label class="form-label">Start Days <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editDraftForm.start_days" required>
          </div>
          <div>
            <label class="form-label">End Days <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editDraftForm.end_days" required>
          </div>
          <div>
            <label class="form-label">Base Premium Value <span class="text-danger">*</span></label>
            <input type="number" step="0.01" class="form-control" v-model.number="editDraftForm.base_premium_value" required>
          </div>
          <div class="form-grid-3-span">
            <label class="form-label">Rules (JSON) <span class="text-danger">*</span></label>
            <textarea 
              class="form-control" 
              v-model="editDraftForm.rules" 
              rows="4"
              placeholder='{"min_child_age": 0, "max_child_age": 17}'
              required
            ></textarea>
            <small class="form-text text-muted">Enter valid JSON format</small>
          </div>
        </div>
        
        <div class="mt-3">
          <button type="submit" class="btn btn-success me-2">
            <i class="fas fa-check"></i> Update
          </button>
          <button type="button" class="btn btn-secondary" @click="cancelEditDraft">
            <i class="fas fa-times"></i> Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- Edit Result Form Section -->
    <div v-if="showEditResultForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Result Data</h3>
      
      <form @submit.prevent="updateResult">
        <div class="form-grid form-grid-3">
          <div>
            <label class="form-label">ID</label>
            <input type="text" class="form-control" v-model="editResultForm.id" readonly>
          </div>
          <div>
            <label class="form-label">Product Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editResultForm.product_code" required>
          </div>
          <div>
            <label class="form-label">Premium Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editResultForm.premium_type" required>
              <option value="">Select premium type...</option>
              <option value="DAILY">DAILY</option>
              <option value="ANNUAL">ANNUAL</option>
            </select>
          </div>
          <div>
            <label class="form-label">Start Days <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editResultForm.start_days" required>
          </div>
          <div>
            <label class="form-label">End Days <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editResultForm.end_days" required>
          </div>
          <div>
            <label class="form-label">Base Premium Value <span class="text-danger">*</span></label>
            <input type="number" step="0.01" class="form-control" v-model.number="editResultForm.base_premium_value" required>
          </div>
          <div class="form-grid-3-span">
            <label class="form-label">Rules (JSON) <span class="text-danger">*</span></label>
            <textarea 
              class="form-control" 
              v-model="editResultForm.rules" 
              rows="4"
              placeholder='{"min_child_age": 0, "max_child_age": 17}'
              required
            ></textarea>
            <small class="form-text text-muted">Enter valid JSON format</small>
          </div>
        </div>
        
        <div class="mt-3">
          <button type="submit" class="btn btn-success me-2">
            <i class="fas fa-check"></i> Update
          </button>
          <button type="button" class="btn btn-secondary" @click="cancelEditResult">
            <i class="fas fa-times"></i> Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- Result Table Section -->
    <div v-if="showResultSection && !showEditResultForm" class="table-section result-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📈 Data Result: Product Rules (From Database)</h4>
        <div>
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
                <th>ID</th>
                <th>Product Name</th>
                <th>Premium Type</th>
                <th>Start Days</th>
                <th>End Days</th>
                <th>Base Premium Value</th>
                <th>Rules</th>
                <th>Created By</th>
                <th>Created At</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="resultData.length === 0">
                <td colspan="10" class="text-center text-muted py-4">
                  <i class="fas fa-info-circle me-2"></i>No product rules found
                </td>
              </tr>
              <tr v-for="item in resultData" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.product_name || item.product_code }}</td>
                <td>{{ item.premium_type }}</td>
                <td>{{ item.start_days }}</td>
                <td>{{ item.end_days }}</td>
                <td>{{ item.base_premium_value }}</td>
                <td>
                  <div style="max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.rules">
                    {{ item.rules }}
                  </div>
                </td>
                <td>{{ item.created_by }}</td>
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
</template>

<script>
export default {
  name: 'ProductRulesSection',
  props: {
    selectedInsuranceCode: String
  },
  data() {
    return {
      file: null,
      draftData: [],
      resultData: [],
      showDraftSection: false,
      showResultSection: false,
      showEditForm: false,
      showEditResultForm: false,
      canConfirm: false,
      editingDraftIndex: null,
      editDraftForm: {
        product_code: '',
        premium_type: '',
        start_days: 0,
        end_days: 0,
        base_premium_value: 0,
        rules: ''
      },
      editResultForm: {
        id: null,
        product_code: '',
        premium_type: '',
        start_days: 0,
        end_days: 0,
        base_premium_value: 0,
        rules: '',
        updated_by: 1
      },
      API_URL: 'http://localhost:8080/api/travel',
      PRODUCT_RULES_API_URL: 'http://localhost:8080/api/product-rules' // Shared legacy route
    }
  },
  watch: {
    selectedInsuranceCode(newVal) {
      // Reload draft and result when insurance code changes
      if (this.showDraftSection) {
        this.loadDraft()
      }
      if (this.showResultSection) {
        this.loadResult()
      }
    }
  },
  methods: {
    async downloadTemplate() {
      try {
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/template`)
        if (!response.ok) {
          window.showCustomAlert('Failed to download template', 'error')
          return
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = 'product_rules_template.csv'
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
        window.showCustomAlert('Template downloaded successfully!', 'success')
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
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
        
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/upload`, {
          method: 'POST',
          body: formData
        })
        
        const result = await response.json()
        
        if (response.ok) {
          this.draftData = result.data || []
          this.canConfirm = this.draftData.length > 0
          this.showDraftSection = true
          window.showCustomAlert('File uploaded successfully!', 'success')
          this.$nextTick(() => {
            setTimeout(() => {
              const draftSection = document.querySelector('.draft-section')
              if (draftSection) {
                draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
              }
            }, 100)
          })
        } else {
          window.showCustomAlert(result.error || 'Upload failed', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    async showDraft() {
      await this.loadDraft()
      this.showDraftSection = true
      this.canConfirm = this.draftData.length > 0
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data available', 'info')
      } else {
        this.$nextTick(() => {
          setTimeout(() => {
            const draftSection = document.querySelector('.draft-section')
            if (draftSection) {
              draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
            }
          }, 100)
        })
      }
    },
    
    async loadDraft() {
      try {
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/draft`)
        if (response.ok) {
          let data = await response.json()
          // Filter by selectedInsuranceCode if provided
          // Extract insurance code from product_code (format: TV-{INSURANCE_CODE}-...)
          if (this.selectedInsuranceCode) {
            data = data.filter(item => {
              if (!item.product_code) return false
              const parts = item.product_code.split('-')
              return parts.length >= 2 && parts[1] === this.selectedInsuranceCode
            })
          }
          this.draftData = data
          this.canConfirm = this.draftData.length > 0
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Clear all draft data?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/draft/clear`, { method: 'POST' })
        if (response.ok) {
          this.draftData = []
          this.canConfirm = false
          this.showDraftSection = false
          window.showCustomAlert('Draft cleared!', 'success')
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    hideDraft() {
      this.showDraftSection = false
    },
    
    async saveRules() {
      if (!this.canConfirm || this.draftData.length === 0) {
        window.showCustomAlert('No data to save', 'warning')
        return
      }
      
      try {
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/draft/confirm`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ created_by: 1 })
        })
        
        if (response.ok) {
          window.showCustomAlert('Product rules saved successfully!', 'success')
          this.draftData = []
          this.canConfirm = false
          this.showDraftSection = false
          await this.loadResult()
        } else {
          const error = await response.json()
          window.showCustomAlert(error.error || 'Failed to save', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    async showResult() {
      await this.loadResult()
      this.showResultSection = true
      this.$nextTick(() => {
        setTimeout(() => {
          const resultSection = document.querySelector('.result-section')
          if (resultSection) {
            resultSection.scrollIntoView({ behavior: 'smooth', block: 'center' })
          }
        }, 200)
      })
    },
    
    async loadResult() {
      try {
        // Use shared legacy route /api/product-rules instead of /api/travel/product-rules
        let url = `${this.PRODUCT_RULES_API_URL}`
        // Add insurance_code filter if provided
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          const result = await response.json()
          // Handle both response formats: {data: [...]} or {status: "Success", data: [...]}
          if (Array.isArray(result)) {
            this.resultData = result
          } else if (result.data) {
            this.resultData = Array.isArray(result.data) ? result.data : []
          } else {
            this.resultData = []
          }
        } else {
          console.error('Failed to load product rules:', response.status, response.statusText)
          this.resultData = []
        }
      } catch (error) {
        console.error('Error loading product rules:', error)
        this.resultData = []
      }
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    editDraft(index) {
      const item = this.draftData[index]
      if (!item) {
        window.showCustomAlert('Invalid item to edit', 'error')
        return
      }
      
      // Handle rules field - convert to JSON string if it's an object
      let rulesString = ''
      if (item.rules) {
        if (typeof item.rules === 'string') {
          // Try to parse and format if it's already a JSON string
          try {
            const parsed = JSON.parse(item.rules)
            rulesString = JSON.stringify(parsed, null, 2)
          } catch (e) {
            rulesString = item.rules
          }
        } else {
          // If it's an object, stringify it
          rulesString = JSON.stringify(item.rules, null, 2)
        }
      } else {
        rulesString = '{}'
      }
      
      this.editDraftForm = {
        product_code: item.product_code || '',
        premium_type: item.premium_type || '',
        start_days: item.start_days || 0,
        end_days: item.end_days || 0,
        base_premium_value: item.base_premium_value || 0,
        rules: rulesString
      }
      this.editingDraftIndex = index
      this.showEditForm = true
      this.showDraftSection = false
    },
    
    async updateDraft() {
      if (this.editingDraftIndex === null) {
        window.showCustomAlert('Invalid draft index', 'error')
        return
      }
      
      try {
        // Validate JSON format for rules
        let rulesValue = this.editDraftForm.rules
        try {
          JSON.parse(rulesValue)
        } catch (e) {
          window.showCustomAlert('Invalid JSON format in Rules field', 'error')
          return
        }
        
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/draft/${this.editingDraftIndex}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.editDraftForm)
        })
        
        if (response.ok) {
          window.showCustomAlert('Draft updated successfully!', 'success')
          await this.loadDraft()
          this.cancelEditDraft()
        } else {
          const errorData = await response.json().catch(() => ({ error: 'Failed to update draft' }))
          window.showCustomAlert(errorData.error || 'Failed to update draft', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error. Please ensure backend is running', 'error')
      }
    },
    
    cancelEditDraft() {
      this.showEditForm = false
      this.editingDraftIndex = null
      this.editDraftForm = {
        product_code: '',
        premium_type: '',
        start_days: 0,
        end_days: 0,
        base_premium_value: 0,
        rules: ''
      }
      if (this.draftData.length > 0) {
        this.showDraftSection = true
      }
    },
    
    async deleteDraft(index) {
      const item = this.draftData[index]
      if (!item) {
        window.showCustomAlert('Invalid item to delete', 'error')
        return
      }
      
      const confirmed = await window.showCustomConfirm('Apakah Anda yakin ingin menghapus data ini?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/draft/${index}`, {
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
    
    editResult(item) {
      console.log('Editing result item:', item)
      
      // Handle rules field - convert to JSON string if it's an object
      let rulesString = ''
      if (item.rules) {
        if (typeof item.rules === 'string') {
          // Try to parse and format if it's already a JSON string
          try {
            const parsed = JSON.parse(item.rules)
            rulesString = JSON.stringify(parsed, null, 2)
          } catch (e) {
            rulesString = item.rules
          }
        } else {
          // If it's an object, stringify it
          rulesString = JSON.stringify(item.rules, null, 2)
        }
      } else {
        rulesString = '{}'
      }
      
      this.editResultForm = {
        id: item.id || null,
        product_code: item.product_code || '',
        premium_type: item.premium_type || '',
        start_days: item.start_days || 0,
        end_days: item.end_days || 0,
        base_premium_value: item.base_premium_value || 0,
        rules: rulesString,
        updated_by: 1
      }
      console.log('Edit form loaded with result data:', this.editResultForm)
      this.showEditResultForm = true
      this.showResultSection = false
    },
    
    async updateResult() {
      try {
        // Validate JSON format for rules
        try {
          JSON.parse(this.editResultForm.rules)
        } catch (e) {
          window.showCustomAlert('Invalid JSON format in Rules field', 'error')
          return
        }
        
        const response = await fetch(`${this.PRODUCT_RULES_API_URL}/${this.editResultForm.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            product_code: this.editResultForm.product_code,
            premium_type: this.editResultForm.premium_type,
            start_days: this.editResultForm.start_days,
            end_days: this.editResultForm.end_days,
            base_premium_value: this.editResultForm.base_premium_value,
            rules: this.editResultForm.rules,
            updated_by: this.editResultForm.updated_by
          })
        })
        
        if (response.ok) {
          window.showCustomAlert('Product rule updated successfully!', 'success')
          this.cancelEditResult()
          await this.loadResult()
          this.showResultSection = true
        } else {
          const error = await response.json()
          window.showCustomAlert(error.error || 'Failed to update', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    cancelEditResult() {
      this.showEditResultForm = false
      this.showResultSection = true
      this.editResultForm = {
        id: null,
        product_code: '',
        premium_type: '',
        start_days: 0,
        end_days: 0,
        base_premium_value: 0,
        rules: '',
        updated_by: 1
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
    }
  }
}
</script>

<style scoped>
.form-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.table-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.table-header {
  margin-bottom: 15px;
}

.table-content {
  overflow-x: auto;
}
</style>
