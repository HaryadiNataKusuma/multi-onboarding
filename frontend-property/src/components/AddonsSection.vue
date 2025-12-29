<template>
  <div class="col-12">
    <div class="form-section">
      <h3><i class="fas fa-file-excel me-1"></i> Addon Rule Details Management</h3>
      
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

    <!-- Edit Form Section -->
    <div v-if="showEditForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Draft Data</h3>
      
      <form @submit.prevent="updateDraft">
        <div class="form-grid form-grid-3">
          <div>
            <label class="form-label">Product Name</label>
            <input type="text" class="form-control" v-model="editForm.product_name" readonly>
          </div>
          <div>
            <label class="form-label">Addon Name</label>
            <input type="text" class="form-control" v-model="editForm.addon_name" readonly>
          </div>
          <div>
            <label class="form-label">Product Code</label>
            <input type="text" class="form-control" v-model="editForm.product_code" readonly>
          </div>
          <div>
            <label class="form-label">Addon Code</label>
            <input type="text" class="form-control" v-model="editForm.addon_code" readonly>
          </div>
          <div>
            <label class="form-label">Addon Rule ID</label>
            <input type="text" class="form-control" v-model="editForm.addon_rule_id" readonly>
          </div>
          <div>
            <label class="form-label">Start Condition <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.start_condition" required>
          </div>
          <div>
            <label class="form-label">End Condition <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.end_condition" required>
          </div>
          <div>
            <label class="form-label">Value Type <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.value_type" required>
          </div>
          <div>
            <label class="form-label">Value <span class="text-danger">*</span></label>
            <input type="number" step="0.01" class="form-control" v-model.number="editForm.value" required>
          </div>
          <div>
            <label class="form-label">Duration Rule Type <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.duration_rule_type" required>
          </div>
          <div>
            <label class="form-label">Min Adult <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editForm.min_adult" required>
          </div>
          <div>
            <label class="form-label">Max Adult <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editForm.max_adult" required>
          </div>
          <div>
            <label class="form-label">Max Age <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model.number="editForm.max_age" required>
          </div>
        </div>
        
        <div class="mt-3">
          <button type="submit" class="btn btn-success me-2">
            <i class="fas fa-check"></i> Update
          </button>
          <button type="button" class="btn btn-secondary" @click="cancelEdit">
            <i class="fas fa-times"></i> Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- Edit Draft Form Section -->
    <div v-if="showEditDraftForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Draft Data</h3>
      
      <form @submit.prevent="updateDraftAddon">
        <div class="form-grid form-grid-2">
          <div>
            <label class="form-label">Addon Code</label>
            <input type="text" class="form-control" v-model="editDraftForm.addon_code" required>
          </div>
          <div>
            <label class="form-label">Product Code</label>
            <input type="text" class="form-control" v-model="editDraftForm.product_code" required>
          </div>
          <div>
            <label class="form-label">Construction Class ID</label>
            <input type="text" class="form-control" v-model="editDraftForm.construction_class_id" required>
          </div>
          <div>
            <label class="form-label">Occupation Code</label>
            <input type="text" class="form-control" v-model="editDraftForm.occupation_code" required>
          </div>
        </div>
        <div class="mt-3">
          <label class="form-label">Rules (JSON)</label>
          <textarea class="form-control" v-model="editDraftForm.rules" rows="5" required></textarea>
        </div>
        <div class="mt-3">
          <button type="submit" class="btn btn-primary">
            <i class="fas fa-save me-1"></i> Update
          </button>
          <button type="button" class="btn btn-secondary ms-2" @click="cancelEditDraft">
            <i class="fas fa-times me-1"></i> Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- Draft Table Section -->
    <div v-if="showDraftSection && !showEditForm && !showEditDraftForm" class="table-section draft-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📊 Data Preview: Addon Rule Details (Draft Save)</h4>
        <div>
          <button class="btn btn-success btn-sm me-2" @click="saveDetails" :disabled="!canConfirm">
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
                <th>Addon Code</th>
                <th>Product Code</th>
                <th>Construction Class ID</th>
                <th>Occupation Code</th>
                <th>Rules</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="draftData.length === 0">
                <td colspan="6" class="text-center text-muted py-4">
                  <i class="fas fa-info-circle me-2"></i>No draft data available
                </td>
              </tr>
              <tr v-for="(item, index) in draftData" :key="index">
                <td>{{ item.addon_code }}</td>
                <td>{{ item.product_code }}</td>
                <td>{{ item.construction_class_id }}</td>
                <td>{{ item.occupation_code }}</td>
                <td style="max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.rules">{{ item.rules || '-' }}</td>
                <td>
                  <button class="btn btn-sm btn-warning me-1" @click="editDraftAddon(index)">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button class="btn btn-sm btn-danger" @click="deleteDraft(index)">
                    <i class="fas fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Edit Result Form Section -->
    <div v-if="showEditResultForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Addon Rule</h3>
      
      <form @submit.prevent="updateResultAddon">
        <div class="form-grid form-grid-2">
          <div>
            <label class="form-label">ID</label>
            <input type="text" class="form-control" v-model="editResultForm.id" readonly>
          </div>
          <div>
            <label class="form-label">Addon Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editResultForm.addon_code" required>
          </div>
          <div>
            <label class="form-label">Product Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editResultForm.product_code" required>
          </div>
          <div>
            <label class="form-label">Construction Class ID <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editResultForm.construction_class_id" required>
          </div>
          <div>
            <label class="form-label">Occupation Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editResultForm.occupation_code" required>
          </div>
        </div>
        <div class="mt-3">
          <label class="form-label">Rules (JSON) <span class="text-danger">*</span></label>
          <textarea class="form-control" v-model="editResultForm.rules" rows="5" required></textarea>
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
        <h4>📈 Data Result: Addon Rule Details (From Database)</h4>
        <div>
          <button class="btn btn-info btn-sm me-2" @click="refreshResult">
            <i class="fas fa-sync-alt"></i> Refresh Data
          </button>
          <button type="button" class="btn btn-secondary btn-sm" @click="hideResult">
            <i class="fas fa-times"></i> Close
          </button>
        </div>
      </div>

      <!-- Filter Section -->
      <div class="p-3 border-bottom" style="background-color: #f8f9fa;">
        <div class="row g-3">
          <div class="col-md-4">
            <label class="form-label">Filter by Addon Code</label>
            <select class="form-select form-select-sm" v-model="resultFilter.addon_code" @change="applyResultFilter">
              <option value="">All Addon Codes</option>
              <option v-for="addon in uniqueAddonCodes" :key="addon" :value="addon">
                {{ addon }}
              </option>
            </select>
          </div>
          <div class="col-md-4">
            <label class="form-label">Filter by Product Code</label>
            <select class="form-select form-select-sm" v-model="resultFilter.product_code" @change="applyResultFilter">
              <option value="">All Product Codes</option>
              <option v-for="product in uniqueProductCodes" :key="product" :value="product">
                {{ product }}
              </option>
            </select>
          </div>
          <div class="col-md-4 d-flex align-items-end">
            <button class="btn btn-warning btn-sm" @click="clearResultFilter">
              <i class="fas fa-times"></i> Clear Filter
            </button>
          </div>
        </div>
      </div>

      <div class="table-content">
        <div class="table-responsive">
          <table class="table table-striped table-hover align-middle">
            <thead class="table-success">
              <tr>
                <th>Addon Code</th>
                <th>Product Code</th>
                <th>Construction Class ID</th>
                <th>Occupation Code</th>
                <th>Rules</th>
                <th>Created By</th>
                <th>Created At</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="filteredResultData.length === 0">
                <td colspan="8" class="text-center text-muted py-4">
                  <i class="fas fa-info-circle me-2"></i>No addon rules found
                </td>
              </tr>
              <tr v-for="item in filteredResultData" :key="item.id">
                <td>{{ item.addon_code || '-' }}</td>
                <td>{{ item.product_code || '-' }}</td>
                <td>{{ item.construction_class_id || '-' }}</td>
                <td>{{ item.occupation_code || '-' }}</td>
                <td style="max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="item.rules">{{ item.rules || '-' }}</td>
                <td>{{ item.created_by || '-' }}</td>
                <td>{{ formatDateOnly(item.created_at) }}</td>
                <td>
                  <button class="btn btn-sm btn-warning" @click="editResultAddon(item)">
                    <i class="fas fa-edit"></i>
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
  name: 'AddonsSection',
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
      canConfirm: false,
      editingIndex: -1,
      editForm: {
        product_name: '',
        addon_name: '',
        product_code: '',
        addon_code: '',
        addon_rule_id: '',
        start_condition: '-1',
        end_condition: '-1',
        value_type: 'FIXED',
        value: 0,
        duration_rule_type: 'DAILY',
        min_adult: 0,
        max_adult: 0,
        max_age: 0
      },
      API_URL: 'http://localhost:8080/api/property',
      ADDON_RULE_DETAILS_API_URL: 'http://localhost:8080/api/property/addon-rule-details', // Property-specific route
      addonRuleMap: new Map(), // Map addon_rule_id -> { addon_name, product_name }
      resultFilter: {
        addon_code: '',
        product_code: ''
      },
      showEditResultForm: false,
      editResultForm: {
        id: null,
        addon_code: '',
        product_code: '',
        construction_class_id: '',
        occupation_code: '',
        rules: ''
      },
      showEditDraftForm: false,
      editDraftIndex: -1,
      editDraftForm: {
        addon_code: '',
        product_code: '',
        construction_class_id: '',
        occupation_code: '',
        rules: ''
      }
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
  computed: {
    filteredResultData() {
      let filtered = this.resultData
      
      if (this.resultFilter.addon_code) {
        filtered = filtered.filter(item => item.addon_code === this.resultFilter.addon_code)
      }
      
      if (this.resultFilter.product_code) {
        filtered = filtered.filter(item => item.product_code === this.resultFilter.product_code)
      }
      
      return filtered
    },
    uniqueAddonCodes() {
      const codes = [...new Set(this.resultData.map(item => item.addon_code).filter(Boolean))]
      return codes.sort()
    },
    uniqueProductCodes() {
      const codes = [...new Set(this.resultData.map(item => item.product_code).filter(Boolean))]
      return codes.sort()
    },
    uniqueAddons() {
      const addons = new Set()
      this.resultData.forEach(item => {
        const addonName = this.getAddonName(item.addon_rule_id)
        if (addonName) {
          addons.add(addonName)
        }
      })
      return Array.from(addons).sort()
    },
    uniqueProducts() {
      const products = new Set()
      this.resultData.forEach(item => {
        const productName = this.getProductName(item.addon_rule_id)
        if (productName) {
          products.add(productName)
        }
      })
      return Array.from(products).sort()
    }
  },
  methods: {
    async downloadTemplate() {
      try {
        // Use property-specific endpoint for template download
        const response = await fetch(`${this.API_URL}/addon-rule-details/template`)
        if (!response.ok) {
          window.showCustomAlert('Failed to download template', 'error')
          return
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = 'property_addon_rule_details_template.csv'
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
        
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/upload`, {
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
      this.canConfirm = this.draftData && this.draftData.length > 0
      if (!this.draftData || this.draftData.length === 0) {
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
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/draft`)
        if (response.ok) {
          let data = await response.json() || []
          // Filter by selectedInsuranceCode if provided (filter directly from product_code in draft)
          if (this.selectedInsuranceCode && data.length > 0) {
            data = data.filter(item => {
              if (!item.product_code) return false
              // Match insurance code in product_code (e.g., PR-DAMAI-xxx or PR-HOME-DAMAI-xxx)
              const parts = item.product_code.split('-')
              if (parts.length >= 2) {
                // Check if second part is insurance code or third part (for PR-HOME-xxx pattern)
                return parts[1] === this.selectedInsuranceCode || 
                       (parts.length >= 3 && parts[2] === this.selectedInsuranceCode)
              }
              return false
            })
          }
          this.draftData = data
          this.canConfirm = this.draftData && this.draftData.length > 0
        } else {
          this.draftData = []
          this.canConfirm = false
        }
      } catch (error) {
        console.error('Error:', error)
        this.draftData = []
        this.canConfirm = false
      }
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Clear all draft data?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/draft/clear`, { method: 'POST' })
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
    
    async saveDetails() {
      if (!this.canConfirm || this.draftData.length === 0) {
        window.showCustomAlert('No data to save', 'warning')
        return
      }
      
      try {
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/draft/confirm`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ created_by: 1 })
        })
        
        if (response.ok) {
          window.showCustomAlert('Addon rule details saved successfully!', 'success')
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
        // Load data from property_service_development.addon_rules
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}`)
        if (response.ok) {
          const result = await response.json()
          // Handle both array and wrapped response formats
          let data = []
          if (Array.isArray(result)) {
            data = result
          } else if (result.data) {
            data = Array.isArray(result.data) ? result.data : []
          }
          
          // Filter by selectedInsuranceCode if provided (extract from product_code)
          if (this.selectedInsuranceCode) {
            data = data.filter(item => {
              if (!item.product_code) return false
              // Match insurance code in product_code (e.g., PR-DAMAI-xxx or PR-HOME-DAMAI-xxx)
              const parts = item.product_code.split('-')
              if (parts.length >= 2) {
                return parts[1] === this.selectedInsuranceCode || 
                       (parts.length >= 3 && parts[2] === this.selectedInsuranceCode)
              }
              return false
            })
          }
          
          this.resultData = data
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    async loadAddonAndProductNames() {
      try {
        // Get all unique addon_rule_ids
        const addonRuleIds = [...new Set(this.resultData.map(item => item.addon_rule_id).filter(id => id))]
        
        // Fetch all addon rules from shared route
        const addonRulesResponse = await fetch('http://localhost:8080/api/addon-rules')
        if (!addonRulesResponse.ok) {
          console.error('Failed to load addon rules:', addonRulesResponse.status)
          return
        }
        let addonRules = await addonRulesResponse.json() || []
        if (!Array.isArray(addonRules)) {
          addonRules = addonRules.data || []
        }
        
        // Filter by insurance_code if provided
        if (this.selectedInsuranceCode) {
          addonRules = addonRules.filter(ar => ar.insurance_code === this.selectedInsuranceCode)
        }
        
        // Fetch all addons from travel domain
        const addonsResponse = await fetch(`${this.API_URL}/addons`)
        if (!addonsResponse.ok) {
          console.error('Failed to load addons:', addonsResponse.status)
          return
        }
        let addons = await addonsResponse.json() || []
        if (!Array.isArray(addons)) {
          addons = addons.data || []
        }
        
        // Fetch all products (filter by insurance_code if provided)
        let productsUrl = `${this.API_URL}/products`
        if (this.selectedInsuranceCode) {
          productsUrl += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const productsResponse = await fetch(productsUrl)
        if (!productsResponse.ok) {
          console.error('Failed to load products:', productsResponse.status)
          return
        }
        let products = await productsResponse.json() || []
        if (!Array.isArray(products)) {
          products = products.data || []
        }
        
        // Create maps
        const addonMap = new Map()
        addons.forEach(addon => {
          if (addon && addon.code && addon.name) {
            addonMap.set(addon.code, addon.name)
          }
        })
        
        const productMap = new Map()
        products.forEach(product => {
          if (product && product.code && product.name) {
            productMap.set(product.code, product.name)
          }
        })
        
        // Build addonRuleMap: addon_rule_id -> { addon_name, product_name }
        this.addonRuleMap.clear()
        addonRules.forEach(rule => {
          if (rule && rule.id) {
            const addonName = rule.addon_code ? addonMap.get(rule.addon_code) : null
            const productName = rule.product_code ? productMap.get(rule.product_code) : null
            this.addonRuleMap.set(rule.id, {
              addon_name: addonName || '-',
              product_name: productName || '-'
            })
          }
        })
      } catch (error) {
        console.error('Error loading addon and product names:', error)
      }
    },
    
    getAddonName(addonRuleId) {
      if (!addonRuleId) return null
      const rule = this.addonRuleMap.get(addonRuleId)
      return rule ? rule.addon_name : null
    },
    
    getProductName(addonRuleId) {
      if (!addonRuleId) return null
      const rule = this.addonRuleMap.get(addonRuleId)
      return rule ? rule.product_name : null
    },
    
    applyResultFilter() {
      // Filter is applied via computed property filteredResultData
    },
    
    clearResultFilter() {
      this.resultFilter = {
        addon_code: '',
        product_code: ''
      }
    },
    
    async refreshResult() {
      await this.loadResult()
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    deleteDraft(index) {
      this.draftData.splice(index, 1)
      this.canConfirm = this.draftData.length > 0
    },
    
    editDraft(index) {
      const item = this.draftData[index]
      this.editForm = {
        product_name: item.product_name || '',
        addon_name: item.addon_name || '',
        product_code: item.product_code || '',
        addon_code: item.addon_code || '',
        addon_rule_id: item.addon_rule_id || '',
        start_condition: item.start_condition || '-1',
        end_condition: item.end_condition || '-1',
        value_type: item.value_type || 'FIXED',
        value: item.value || 0,
        duration_rule_type: item.duration_rule_type || 'DAILY',
        min_adult: item.min_adult || 0,
        max_adult: item.max_adult || 0,
        max_age: item.max_age || 0
      }
      this.editingIndex = index
      this.showEditForm = true
      this.showDraftSection = false
    },
    
    updateDraft() {
      if (this.editingIndex >= 0 && this.editingIndex < this.draftData.length) {
        this.draftData[this.editingIndex] = {
          ...this.draftData[this.editingIndex],
          start_condition: this.editForm.start_condition,
          end_condition: this.editForm.end_condition,
          value_type: this.editForm.value_type,
          value: this.editForm.value,
          duration_rule_type: this.editForm.duration_rule_type,
          min_adult: this.editForm.min_adult,
          max_adult: this.editForm.max_adult,
          max_age: this.editForm.max_age
        }
        window.showCustomAlert('Draft updated successfully!', 'success')
        this.cancelEdit()
        this.showDraftSection = true
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      this.editingIndex = -1
    },
    
    editDraftAddon(index) {
      const item = this.draftData[index]
      this.editDraftForm = {
        addon_code: item.addon_code || '',
        product_code: item.product_code || '',
        construction_class_id: item.construction_class_id || '',
        occupation_code: item.occupation_code || '',
        rules: item.rules || ''
      }
      this.editDraftIndex = index
      this.showEditDraftForm = true
      this.showDraftSection = false
    },
    
    updateDraftAddon() {
      if (this.editDraftIndex >= 0 && this.editDraftIndex < this.draftData.length) {
        this.draftData[this.editDraftIndex] = {
          addon_code: this.editDraftForm.addon_code,
          product_code: this.editDraftForm.product_code,
          construction_class_id: this.editDraftForm.construction_class_id,
          occupation_code: this.editDraftForm.occupation_code,
          rules: this.editDraftForm.rules
        }
        window.showCustomAlert('Draft updated successfully!', 'success')
        this.cancelEditDraft()
        this.showDraftSection = true
      }
    },
    
    cancelEditDraft() {
      this.showEditDraftForm = false
      this.editDraftIndex = -1
      this.editDraftForm = {
        addon_code: '',
        product_code: '',
        construction_class_id: '',
        occupation_code: '',
        rules: ''
      }
      this.editForm = {
        product_name: '',
        addon_name: '',
        product_code: '',
        addon_code: '',
        addon_rule_id: '',
        start_condition: '-1',
        end_condition: '-1',
        value_type: 'FIXED',
        value: 0,
        duration_rule_type: 'DAILY',
        min_adult: 0,
        max_adult: 0,
        max_age: 0
      }
    },
    
    editResultAddon(item) {
      this.editResultForm = {
        id: item.id,
        addon_code: item.addon_code || '',
        product_code: item.product_code || '',
        construction_class_id: item.construction_class_id || '',
        occupation_code: item.occupation_code || '',
        rules: item.rules || ''
      }
      this.showEditResultForm = true
      this.showResultSection = false
    },
    
    async deleteResultAddon(id) {
      const confirmed = await window.showCustomConfirm('Are you sure you want to delete this addon rule?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/${id}`, {
          method: 'DELETE'
        })
        
        if (response.ok) {
          window.showCustomAlert('Addon rule deleted successfully!', 'success')
          await this.loadResult()
        } else {
          const error = await response.json()
          window.showCustomAlert(error.error || 'Failed to delete', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    async updateResultAddon() {
      try {
        const response = await fetch(`${this.ADDON_RULE_DETAILS_API_URL}/${this.editResultForm.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            addon_code: this.editResultForm.addon_code,
            product_code: this.editResultForm.product_code,
            construction_class_id: this.editResultForm.construction_class_id,
            occupation_code: this.editResultForm.occupation_code,
            rules: this.editResultForm.rules
          })
        })
        
        if (response.ok) {
          window.showCustomAlert('Addon rule updated successfully!', 'success')
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
        addon_code: '',
        product_code: '',
        construction_class_id: '',
        occupation_code: '',
        rules: ''
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
.btn-product-code-fixed {
  width: 100%;
  min-height: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding: 0.375rem 0.75rem;
  white-space: nowrap;
  overflow: hidden;
  text-align: left;
  position: relative;
}

.btn-product-code-fixed i {
  flex-shrink: 0;
}

.btn-product-code-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
</style>
