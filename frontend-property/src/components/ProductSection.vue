<template>
  <div class="col-12">
    <div class="form-section">
      <h3><i class="fas fa-plus-circle me-1"></i> Tambah Data Product Baru</h3>
      
      <form @submit.prevent="saveProduct">
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Insurance Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="form.insurance_code" placeholder="Cth: DAMAI" readonly>
          </div>
          
          <div>
            <label class="form-label">Product Category <span class="text-danger">*</span></label>
            <select class="form-select" v-model="form.product_category" required>
              <option value="">Pilih Product Category</option>
              <option value="PR">PR</option>
              <option value="PR_HOME">PR_HOME</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Coverage Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="form.coverage_type" required>
              <option value="">Pilih Coverage Type</option>
              <option value="FLEXAS">FLEXAS</option>
              <option value="PACKAGE">PACKAGE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Product Name <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="form.name" placeholder="Cth: Damai Travel Domestic Individual" required>
          </div>
        </div>
        
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Hardcopy Admin Fee <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model="form.insurance_type" placeholder="Masukkan nominal" min="0" required>
          </div>
          
          <div>
            <label class="form-label">Is Active <span class="text-danger">*</span></label>
            <select class="form-select" v-model="form.is_active" required>
              <option value="1">1 (Active)</option>
              <option value="0">0 (Inactive)</option>
            </select>
          </div>
        </div>
        
        <div class="form-grid form-grid-1 mt-2">
          <div>
            <label class="form-label">Generated Product Code</label>
            <input type="text" class="form-control" v-model="form.code" readonly>
            <small class="text-muted">Product code akan otomatis di-generate berdasarkan field di atas</small>
          </div>
        </div>
        
        <div class="mt-3 d-flex align-items-center gap-2">
          <button type="button" class="btn btn-info" @click="loadDraft">
            <i class="fas fa-eye"></i> Tampilkan Draft Data
          </button>
          <button type="button" class="btn btn-warning" @click="clearDraft">
            <i class="fas fa-trash"></i> Clear Draft
          </button>
          <button type="submit" class="btn btn-primary">
            <i class="fas fa-save"></i> Save Product
          </button>
          <button type="button" class="btn btn-info" @click="loadResult">
            <i class="fas fa-chart-line"></i> RESULT
          </button>
        </div>
      </form>
    </div>

    <!-- Edit Form Section -->
    <div v-if="showEditForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Data Product</h3>
      
      <form @submit.prevent="updateProduct">
        <input type="hidden" v-model="editForm.id">
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Insurance Code <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.insurance_code" readonly>
          </div>
          
          <div>
            <label class="form-label">Product Category <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.product_category" required>
              <option value="">Pilih Product Category</option>
              <option value="PR">PR</option>
              <option value="PR_HOME">PR_HOME</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Product Name <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.name" placeholder="Cth: Damai Travel Domestic Individual" required>
          </div>
        </div>
        
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Hardcopy Admin Fee <span class="text-danger">*</span></label>
            <input type="number" class="form-control" v-model="editForm.insurance_type" placeholder="Masukkan nominal" min="0" required>
          </div>
          
          <div>
            <label class="form-label">Coverage Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.coverage_type" required>
              <option value="">Pilih Coverage Type</option>
              <option value="FLEXAS">FLEXAS</option>
              <option value="PACKAGE">PACKAGE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Is Active <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.is_active" required>
              <option value="1">1 (Active)</option>
              <option value="0">0 (Inactive)</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Logo</label>
            <input type="text" class="form-control" v-model="editForm.logo" placeholder="Logo URL">
          </div>
        </div>
        
        <div class="form-grid form-grid-1 mt-2">
          <div>
            <label class="form-label">Product Code</label>
            <input type="text" class="form-control" v-model="editForm.code" readonly>
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

    <!-- Draft Table -->
    <div v-if="showDraft" class="table-section draft-section show">
      <div class="table-header">
        <h4>📊 Data Preview: Products (Draft Save)</h4>
        <div class="table-header-actions">
        <button class="btn btn-success btn-confirm" :disabled="!canConfirm" @click="confirmProducts">
          <i class="fas fa-database me-2"></i> Konfirmasi Data (DB)
        </button>
          <button class="btn btn-secondary btn-close-draft" @click="closeDraft">
            <i class="fas fa-times me-2"></i> Close
        </button>
        </div>
      </div>
      <div class="table-content">
        <div class="table-responsive">
          <table class="table table-striped table-hover align-middle">
            <thead>
              <tr>
                <th>ID</th>
                <th>Code</th>
                <th>Insurance Code</th>
                <th>Name</th>
                <th>Logo</th>
                <th>Product Category</th>
                <th>Coverage Type</th>
                <th>Is Active</th>
                <th>Hardcopy Admin Fee</th>
                <th>Summary</th>
                <th>Insurance Detail</th>
                <th>Protection Detail</th>
                <th>How To Claim</th>
                <th>Deductible And Clause</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in draftData" :key="item.id">
                <td>{{ item.id || '-' }}</td>
                <td>{{ item.code || '-' }}</td>
                <td>{{ item.insurance_code || '-' }}</td>
                <td>{{ item.name || '-' }}</td>
                <td>
                  <img v-if="item.logo" :src="item.logo" class="table-logo" alt="Logo" @error="$event.target.style.display='none'">
                  <span v-else>-</span>
                </td>
                <td>{{ item.property_category || '-' }}</td>
                <td>{{ item.coverage_type || '-' }}</td>
                <td>{{ item.is_active || '-' }}</td>
                <td>{{ item.admin_fee != null && item.admin_fee !== undefined ? item.admin_fee : (item.insurance_type != null && item.insurance_type !== undefined ? item.insurance_type : '-') }}</td>
                <td>{{ item.summary || '-' }}</td>
                <td>{{ item.insurance_detail || '-' }}</td>
                <td>{{ item.protection_detail || '-' }}</td>
                <td>{{ item.how_to_claim || '-' }}</td>
                <td>{{ item.code ? 'DEDUCTION_AND_CLAUSE_' + item.code.replace(/-/g, '_') : '-' }}</td>
                <td>
                  <button class="btn btn-sm btn-warning me-1" @click="startEditDraft(item)">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button class="btn btn-sm btn-danger" @click="deleteDraft(item.id)">
                    <i class="fas fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Result Table -->
    <div v-if="showResult" class="table-section result-section show">
      <div class="table-header">
        <h4>📊 Data Result: Products (From Database)</h4>
        <div class="table-header-actions">
          <button class="btn btn-secondary" @click="closeResult">
            <i class="fas fa-times me-2"></i> Close
          </button>
        </div>
      </div>
      <div class="table-content">
        <div class="table-responsive">
          <table class="table table-striped table-hover align-middle">
            <thead>
              <tr>
                <th>ID</th>
                <th>Code</th>
                <th>Insurance Code</th>
                <th>Name</th>
                <th>Logo</th>
                <th>Product Category</th>
                <th>Coverage Type</th>
                <th>Is Active</th>
                <th>Hardcopy Admin Fee</th>
                <th>Summary</th>
                <th>Insurance Detail</th>
                <th>Protection Detail</th>
                <th>How To Claim</th>
                <th>Deductible And Clause</th>
                <th>Created At</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="resultData.length === 0">
                <td colspan="16" class="text-center text-muted py-4">
                  <i class="fas fa-info-circle me-2"></i>No data available
                </td>
              </tr>
              <tr v-for="item in resultData" :key="item.id">
                <td>{{ item.id || '-' }}</td>
                <td>{{ item.code || '-' }}</td>
                <td>{{ item.insurance_code || '-' }}</td>
                <td>{{ item.name || '-' }}</td>
                <td>
                  <img v-if="item.logo" :src="item.logo" class="table-logo" alt="Logo" @error="$event.target.style.display='none'">
                  <span v-else>-</span>
                </td>
                <td>{{ item.property_category || '-' }}</td>
                <td>{{ item.coverage_type || '-' }}</td>
                <td>{{ item.is_active || '-' }}</td>
                <td>{{ (item.admin_fee != null && item.admin_fee !== undefined) ? item.admin_fee : ((item.hardcopy_admin_fee != null && item.hardcopy_admin_fee !== undefined) ? item.hardcopy_admin_fee : '-') }}</td>
                <td>{{ item.summary || '-' }}</td>
                <td>{{ item.insurance_detail || '-' }}</td>
                <td>{{ item.protection_detail || '-' }}</td>
                <td>{{ item.how_to_claim || '-' }}</td>
                <td>{{ item.code ? 'DEDUCTION_AND_CLAUSE_' + item.code.replace(/-/g, '_') : '-' }}</td>
                <td>{{ item.created_at ? new Date(item.created_at).toLocaleString() : '-' }}</td>
                <td>
                  <button class="btn btn-sm btn-warning" @click="startEditResult(item)">
                    <i class="fas fa-edit"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Addons Modal -->
    <div v-if="showAddonsModal" class="modal-overlay" @click.self="closeAddonsModal">
      <div class="modal-content" style="max-width: 600px; max-height: 80vh; overflow-y: auto;">
        <div class="modal-header">
          <h5 class="modal-title">Pilih Addons</h5>
          <button type="button" class="btn-close" @click="closeAddonsModal"></button>
        </div>
        <div class="modal-body">
          <div class="mb-3 d-flex justify-content-between align-items-center">
            <div>
              <strong>Selected: {{ selectedAddons.length }} addon(s)</strong>
            </div>
            <div>
              <button type="button" class="btn btn-sm btn-primary me-2" @click="selectAllAddons">
                <i class="fas fa-check-square"></i> Select All
              </button>
              <button type="button" class="btn btn-sm btn-secondary" @click="deselectAllAddons">
                <i class="fas fa-square"></i> Deselect All
              </button>
            </div>
          </div>
          
          <div v-if="loadingAddons" class="text-center py-4">
            <div class="spinner-border" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>
          
          <div v-else-if="availableAddons.length === 0" class="text-center py-4 text-muted">
            No addons available
          </div>
          
          <div v-else class="addons-list">
            <div 
              v-for="addon in availableAddons" 
              :key="addon.id || addon.code" 
              class="addon-item"
            >
              <input 
                class="addon-checkbox" 
                type="checkbox" 
                :id="`addon-${addon.id || addon.code}`"
                :value="addon.code"
                v-model="selectedAddons"
              >
              <label class="addon-label" :for="`addon-${addon.id || addon.code}`">
                {{ addon.name }}
              </label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeAddonsModal">Cancel</button>
          <button type="button" class="btn btn-primary" @click="saveSelectedAddons">
            <i class="fas fa-check"></i> Save ({{ selectedAddons.length }})
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProductSection',
  props: {
    selectedInsuranceCode: String
  },
  data() {
    return {
      form: {
        code: '',
        insurance_code: '',
        product_category: '', // Product Category: PR, PR_HOME
        coverage_type: '', // FLEXAS, PACKAGE
        type: '',
        insurance_type: '', // Hardcopy Admin Fee (number)
        name: '',
        logo: '',
        is_active: '1'
      },
      editForm: {
        id: null,
        code: '',
        insurance_code: '',
        product_category: '', // Product Category: PR, PR_HOME
        coverage_type: '', // FLEXAS, PACKAGE
        type: '',
        insurance_type: '', // Hardcopy Admin Fee (number)
        name: '',
        logo: '',
        is_active: '1'
      },
      showEditForm: false,
      isEditingDraft: false,
      insurances: [],
      // regions removed for Property program
      draftData: [],
      showDraft: false,
      canConfirm: false,
      resultData: [],
      showResult: false,
      showAddonsModal: false,
      availableAddons: [],
      selectedAddons: [], // Temporary selection in modal
      loadingAddons: false,
      API_URL: 'http://localhost:8080/api/property', // Property program backend
      TEMPLATES_API_URL: 'http://localhost:8080/api/property', // Templates backend (now unified)
      generateCodeTimeout: null // For debouncing product code generation
    }
  },
  async mounted() {
    // Load insurances first, then auto-fill logo
    await this.loadInsurances()
    // loadRegions() removed for Property program
    // Set insurance code from prop if available
    if (this.selectedInsuranceCode) {
      this.form.insurance_code = this.selectedInsuranceCode
      // Auto-fill logo immediately since insurances are already loaded
      this.autoFillLogo()
      // Generate product code after a short delay to ensure regions are loaded
    this.$nextTick(() => {
        setTimeout(() => {
      this.generateProductCode()
        }, 100)
    })
    }
  },
  watch: {
    selectedInsuranceCode(newVal) {
      // Update insurance code in form
      this.form.insurance_code = newVal || ''
      // Auto-fill logo from insurance immediately
      // Use $nextTick to ensure insurances are loaded
      this.$nextTick(() => {
        this.autoFillLogo()
      })
      // Generate product code when insurance code changes
      this.$nextTick(() => {
      this.generateProductCode()
      })
      // Reload draft and result when insurance code changes
      if (this.showDraft) {
        this.loadDraft()
      }
      if (this.showResult) {
        this.loadResult()
      }
    },
    'form.insurance_code'() {
      this.debouncedGenerateProductCode()
      // Auto-fill logo from insurance when insurance_code changes
      // Use $nextTick to ensure insurances are loaded
      this.$nextTick(() => {
        this.autoFillLogo()
      })
    },
    // 'form.region'() removed for Property program
    'form.product_category'() {
      this.debouncedGenerateProductCode()
    },
    'form.coverage_type'() {
      this.debouncedGenerateProductCode()
    },
    'form.insurance_type'() {
      // No need to regenerate code for insurance_type (Hardcopy Admin Fee)
    }
  },
  methods: {
    debouncedGenerateProductCode() {
      // Clear existing timeout
      if (this.generateCodeTimeout) {
        clearTimeout(this.generateCodeTimeout)
      }
      // Set new timeout to debounce API calls
      this.generateCodeTimeout = setTimeout(() => {
        this.generateProductCode()
      }, 300) // Wait 300ms after user stops typing/selecting
    },
    generateProductCode() {
      // Format for Property based on product_category and coverage_type:
      // 1. PR + PACKAGE → PR-{insurance_code}-PAR-SME-01
      // 2. PR + FLEXAS → PR-{insurance_code}-PAR-SME-02
      // 3. PR_HOME + PACKAGE → PR-HOME-{insurance_code}-SR-01
      // 4. PR_HOME + FLEXAS → PR-HOME-{insurance_code}-SR-02
      
      // Check if all required fields are filled
      if (!this.form.insurance_code || !this.form.product_category || !this.form.coverage_type) {
        // Don't clear code if user is still typing, only clear if insurance_code is missing
        if (!this.form.insurance_code) {
          this.form.code = ''
        }
        return
      }
      
      const insuranceCode = this.form.insurance_code.toUpperCase()
      const productCategory = this.form.product_category // PR or PR_HOME
      const coverageType = this.form.coverage_type // PACKAGE or FLEXAS
      
      let pattern = ''
      
      if (productCategory === 'PR') {
        if (coverageType === 'PACKAGE') {
          pattern = `PR-${insuranceCode}-PAR-SME-01`
        } else if (coverageType === 'FLEXAS') {
          pattern = `PR-${insuranceCode}-PAR-SME-02`
        }
      } else if (productCategory === 'PR_HOME') {
        if (coverageType === 'PACKAGE') {
          pattern = `PR-HOME-${insuranceCode}-SR-01`
        } else if (coverageType === 'FLEXAS') {
          pattern = `PR-HOME-${insuranceCode}-SR-02`
        }
      }
      
      if (!pattern) {
        this.form.code = ''
        return
      }
      
      // Get next sequence number (replace last 01/02 with sequence)
      this.getNextSequenceNumber(pattern)
    },
    
    async getNextSequenceNumber(pattern) {
      try {
        const response = await fetch(`${this.API_URL}/products`)
        if (response.ok) {
          const products = await response.json()
          
          // Extract base pattern (without sequence number at the end)
          // Pattern format: PR-{insurance}-PAR-SME-01 or PR-HOME-{insurance}-SR-01
          // We need to match codes that have the same base (before the last -XX)
          const basePattern = pattern.replace(/-\d{2}$/, '')
          
          // Filter products that match the base pattern
          const matchingProducts = products.filter(p => {
            if (!p.code) return false
            // Check if code starts with base pattern (before sequence number)
            const codeWithoutSeq = p.code.replace(/-\d{2}$/, '')
            return codeWithoutSeq === basePattern
          })
          
          // Extract sequence numbers and find the max
          let maxSeq = 0
          matchingProducts.forEach(p => {
            const parts = p.code.split('-')
            const lastPart = parts[parts.length - 1]
            const seq = parseInt(lastPart) || 0
            if (seq > maxSeq) maxSeq = seq
          })
          
          // Generate new code with next sequence number
          // Keep the base pattern and replace the sequence number
          const nextSeq = (maxSeq + 1).toString().padStart(2, '0')
          this.form.code = `${basePattern}-${nextSeq}`
        } else {
          // If API fails, use pattern as is (already has 01 or 02)
          this.form.code = pattern
        }
      } catch (error) {
        console.error('Error generating product code:', error)
        // Fallback: use pattern as is (already has 01 or 02)
        this.form.code = pattern
      }
    },
    
    async generateAnnualProductCodes(pattern) {
      // For ANNUAL products, generate 2 consecutive product codes
      try {
        const response = await fetch(`${this.API_URL}/products`)
        if (response.ok) {
          const products = await response.json()
          
          // Filter products that match the pattern
          const matchingProducts = products.filter(p => {
            if (!p.code) return false
            const codeWithoutSeq = p.code.replace(/-\d{2}$/, '')
            return codeWithoutSeq === pattern
          })
          
          // Extract sequence numbers and find the max
          let maxSeq = 0
          matchingProducts.forEach(p => {
            const parts = p.code.split('-')
            const lastPart = parts[parts.length - 1]
            const seq = parseInt(lastPart) || 0
            if (seq > maxSeq) maxSeq = seq
          })
          
          // Generate consecutive sequence numbers
          const seq90 = (maxSeq + 1).toString().padStart(2, '0')
          const seq180 = (maxSeq + 2).toString().padStart(2, '0')
          
          // Update form with both codes
          this.form.code90Days = `${pattern}-${seq90}`
          this.form.code180Days = `${pattern}-${seq180}`
          // Also set form.code to the first one for backward compatibility
          this.form.code = `${pattern}-${seq90}`
        } else {
          // If API fails, use default sequence
          this.form.code90Days = `${pattern}-01`
          this.form.code180Days = `${pattern}-02`
          this.form.code = `${pattern}-01`
        }
      } catch (error) {
        console.error('Error generating annual product codes:', error)
        // If API fails, use default sequence
        this.form.code90Days = `${pattern}-01`
        this.form.code180Days = `${pattern}-02`
        this.form.code = `${pattern}-01`
      }
    },
    
    async loadInsurances() {
      try {
        const response = await fetch(`${this.API_URL}/insurances`)
        if (response.ok) {
          this.insurances = await response.json()
          console.log('Insurances loaded:', this.insurances.length)
          // Auto-fill logo if insurance_code is already set
          if (this.form.insurance_code) {
            this.autoFillLogo()
          }
        }
      } catch (error) {
        console.error('Error loading insurances:', error)
      }
    },
    
    autoFillLogo() {
      // Auto-fill logo from insurance table based on insurance_code
      if (!this.form.insurance_code) {
        console.log('Cannot auto-fill logo - insurance_code is empty')
        return
      }
      
      if (this.insurances.length === 0) {
        console.log('Cannot auto-fill logo - insurances not loaded yet')
        // Try to load insurances if not loaded
        this.loadInsurances().then(() => {
          this.autoFillLogo()
        })
        return
      }
      
      const insurance = this.insurances.find(ins => ins.code === this.form.insurance_code)
      if (insurance) {
        if (insurance.logo) {
          this.form.logo = insurance.logo
          console.log('✅ Logo auto-filled from insurance:', insurance.logo, 'for code:', this.form.insurance_code)
        } else {
          console.log('⚠️ Insurance found but no logo:', this.form.insurance_code)
          this.form.logo = ''
        }
      } else {
        console.log('❌ Insurance not found for code:', this.form.insurance_code, 'Available codes:', this.insurances.map(i => i.code))
        this.form.logo = ''
      }
    },
    
    // loadRegions() method removed for Property program - regions not needed
    async loadRegions_DISABLED() {
      try {
        const response = await fetch(`${this.API_URL}/regions`)
        if (response.ok) {
          const data = await response.json()
          this.regions = data || []
          console.log('Regions loaded:', this.regions.length)
        } else {
          const errorText = await response.text()
          console.error('Failed to load regions:', response.status, errorText)
          this.regions = []
          window.showCustomAlert('Failed to load regions', 'error')
        }
      } catch (error) {
        console.error('Error loading regions:', error)
        this.regions = []
        window.showCustomAlert('Connection error while loading regions', 'error')
      }
    },
    
    async saveProduct() {
      try {
        // Ensure logo is filled from insurance before saving
        if (!this.form.logo && this.form.insurance_code) {
          this.autoFillLogo()
        }
        
        // Region removed for Property program - no region_id needed
        // Type field removed - no min_adult, max_adult, min_child, max_child needed for Property

        // Convert code: replace "-" with "_" for SUMMARY, INSURANCE_DETAIL, PROTECTION_DETAIL, HOW_TO_CLAIM
        const codeWithUnderscore = this.form.code.replace(/-/g, '_')
        const summary = `SUMMARY_${codeWithUnderscore}`
        const insuranceDetail = `INSURANCE_DETAIL_${codeWithUnderscore}`
        const protectionDetail = `PROTECTION_DETAIL_${codeWithUnderscore}`
        const howToClaim = `HOW_TO_CLAIM_${codeWithUnderscore}`

        // Prepare draft data with all auto-filled fields
        // Ensure logo is included (should already be in form.logo from autoFillLogo)
        // If logo is still empty, try to get it from insurance one more time
        let logoToSave = this.form.logo
        if (!logoToSave && this.form.insurance_code && this.insurances.length > 0) {
          const insurance = this.insurances.find(ins => ins.code === this.form.insurance_code)
          if (insurance && insurance.logo) {
            logoToSave = insurance.logo
            this.form.logo = logoToSave
            console.log('Logo retrieved from insurance at save time:', logoToSave)
          }
        }
        
        // Save directly to database (no draft)
        // Parse hardcopy admin fee - ensure it's a valid integer
        const hardcopyAdminFee = this.form.insurance_type ? parseInt(this.form.insurance_type, 10) : 0
        if (isNaN(hardcopyAdminFee)) {
          window.showCustomAlert('Hardcopy Admin Fee harus berupa angka bulat yang valid', 'error')
          return
        }
        
        const productData = {
          code: this.form.code || '',
          insurance_code: this.form.insurance_code || '',
          name: this.form.name || '',
          logo: logoToSave || '',
          is_active: parseInt(this.form.is_active) || 1,
          admin_fee: hardcopyAdminFee, // Hardcopy Admin Fee - will be saved to hardcopy_admin_fee column
          property_category: this.form.product_category || '', // Product Category (PR, PR_HOME)
          coverage_type: this.form.coverage_type || '', // Coverage Type (FLEXAS, PACKAGE)
          min_coverage: 0, // Set default values
          max_coverage: 0, // Set default values
          summary: summary,
          insurance_detail: insuranceDetail,
          protection_detail: protectionDetail,
          how_to_claim: howToClaim,
          created_by: 1 // Default user ID
        }
        
        console.log('Saving product directly to database:', productData)
        console.log('Hardcopy Admin Fee value:', hardcopyAdminFee)

        const response = await fetch(`${this.API_URL}/products`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(productData)
        })
        
        if (!response.ok) {
          // Try to parse error response
          let errorMsg = 'Failed to save product'
          try {
            const errorData = await response.json()
            errorMsg = errorData.error || errorMsg
          } catch (e) {
            errorMsg = `Server error: ${response.status} ${response.statusText}`
          }
          window.showCustomAlert(errorMsg, 'error')
          return
        }
        
        const result = await response.json()
        window.showCustomAlert('Product saved successfully!', 'success')
        
        // No need for template drafts when saving directly
        
        this.resetForm()
        await this.loadResult()
        this.showResult = true
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      }
    },
    
    async createTemplateDrafts(summaryId, insuranceDetailId, productCode, insuranceCode) {
      try {
        // Create SUMMARY template draft
        const summaryPayload = {
          locale: 'id',
          id: summaryId,
          value: JSON.stringify(["SUMMARY1","SUMMARY2","SUMMARY3","SUMMARY4","SUMMARY5"]),
          insurance_code: insuranceCode,
          product_code: productCode,
          created_by: 1
        }
        
        // Create INSURANCE_DETAIL template draft
        const insuranceDetailPayload = {
          locale: 'id',
          id: insuranceDetailId,
          value: '', // Empty value, user will fill it later
          insurance_code: insuranceCode,
          product_code: productCode,
          created_by: 1
        }
        
        // Insert both templates
        const [summaryResponse, insuranceDetailResponse] = await Promise.all([
          fetch(`${this.TEMPLATES_API_URL}/templates/draft/add`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(summaryPayload)
          }),
          fetch(`${this.TEMPLATES_API_URL}/templates/draft/add`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(insuranceDetailPayload)
          })
        ])
        
        if (summaryResponse.ok && insuranceDetailResponse.ok) {
          console.log('✅ Template drafts created: SUMMARY and INSURANCE_DETAIL')
        } else {
          const summaryError = summaryResponse.ok ? null : await summaryResponse.json().catch(() => ({}))
          const insuranceDetailError = insuranceDetailResponse.ok ? null : await insuranceDetailResponse.json().catch(() => ({}))
          console.warn('⚠️ Some template drafts failed to create:', {
            summary: summaryError,
            insuranceDetail: insuranceDetailError
          })
        }
      } catch (error) {
        console.error('Error creating template drafts:', error)
        // Don't show error to user, just log it
      }
    },
    
    async loadDraft() {
      try {
        const response = await fetch(`${this.API_URL}/products/draft`)
        if (response.ok) {
          let data = await response.json()
          // Filter by selectedInsuranceCode if provided
          if (this.selectedInsuranceCode) {
            data = data.filter(item => item.insurance_code === this.selectedInsuranceCode)
          }
          this.draftData = data
          this.canConfirm = this.draftData.length > 0
          this.showDraft = this.canConfirm
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Clear all draft data?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.API_URL}/products/draft/clear`, { method: 'POST' })
        if (response.ok) {
          this.draftData = []
          this.canConfirm = false
          this.showDraft = false
          window.showCustomAlert('Draft cleared!', 'success')
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    async confirmProducts() {
      try {
        // Save draft data before confirming (so we can use it for addons)
        const draftsBeforeConfirm = [...this.draftData]
        
        const response = await fetch(`${this.API_URL}/products/draft/confirm`, { method: 'POST' })
        if (response.ok) {
          // Save addons for all confirmed products using saved draft data
          await this.saveProductAddons(draftsBeforeConfirm)
          
          window.showCustomAlert('Confirmed successfully!', 'success')
          this.draftData = []
          this.canConfirm = false
          this.showDraft = false
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Error confirming products: ' + error.message, 'error')
      }
    },
    
    getAddonNameMapping() {
      // Mapping addon code to name (all languages use the same name)
      return {
        'COVID-PROTECTION': 'Covid Protection',
        'TRAVEL-IBADAH': 'Travel Ibadah',
        'WINTER-SPORT': 'Winter Sports',
        'TRAVEL-PENDIDIKAN': 'Travel Pendidikan',
        'TRIP-CANCELLATION-EXTENSION': 'Trip Cancellation Extension',
        'GOLF-BENEFIT': 'Golf Benefit',
        'LOSS-OF-HOTEL-RESERVATION': 'Loss of Hotel Reservation',
        'PET-INSURANCE': 'Pet insurance',
        'MISS-EVENT': 'Missed event',
        'ADVENTURE': 'Adventures Activities Cover',
        'VISA-PROTECTION': 'Visa Protection',
        'AMATEUR-SPORT-EVENT': 'Amateur Sports Event',
        'BUSINESS-TRAVEL': 'Business Travel',
        'SNOW-SPORT': 'Snow Sport Cover'
      }
    },
    
    async saveProductAddons(drafts) {
      try {
        // Use provided drafts or fallback to this.draftData
        const draftsToUse = drafts || this.draftData
        const addonNameMap = this.getAddonNameMapping()
        
        console.log('Saving addons for drafts:', draftsToUse)
        console.log('Draft data details:', draftsToUse.map(d => ({ code: d.code, addons: d.addons })))
        
        for (const draft of draftsToUse) {
          // Check if draft has addons property and it's an array with items
          const draftAddons = draft.addons || []
          console.log(`Product ${draft.code} - addons:`, draftAddons, 'type:', typeof draftAddons, 'isArray:', Array.isArray(draftAddons))
          
          if (Array.isArray(draftAddons) && draftAddons.length > 0) {
            console.log(`Processing addons for product ${draft.code}:`, draftAddons)
            
            // Prepare addons with names from mapping
            const addonsWithNames = draftAddons.map(addonCode => {
              const addonName = addonNameMap[addonCode] || addonCode
              
              return {
                code: addonCode,
                name: addonName,
                name_my: addonName,
                name_en: addonName
              }
            })
            
            const payload = {
              product_code: draft.code,
              insurance_code: draft.insurance_code,
              addons: addonsWithNames
            }
            
            console.log('Sending payload:', payload)
            
            const response = await fetch(`${this.API_URL}/products/addons`, {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(payload)
            })
            
            if (!response.ok) {
              const errorText = await response.text()
              console.error(`Failed to save addons for product ${draft.code}:`, errorText)
              window.showCustomAlert(`Failed to save addons for product ${draft.code}`, 'error')
            } else {
              const result = await response.json()
              console.log(`Successfully saved addons for product ${draft.code}:`, result)
            }
          } else {
            console.log(`No addons to save for product ${draft.code}`)
          }
        }
      } catch (error) {
        console.error('Error saving product addons:', error)
        window.showCustomAlert('Error saving addons: ' + error.message, 'error')
      }
    },
    
    async loadResult() {
      try {
        let url = `${this.API_URL}/products`
        // Filter by selectedInsuranceCode if provided
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          let data = await response.json()
          // Additional client-side filter if needed
          if (this.selectedInsuranceCode) {
            data = data.filter(item => item.insurance_code === this.selectedInsuranceCode)
          }
          // Debug: log data to see admin_fee values
          console.log('Loaded products data:', data)
          console.log('Sample product admin_fee:', data.length > 0 ? data[0].admin_fee : 'no data')
          this.resultData = data
          this.showResult = true
        } else {
          window.showCustomAlert('Failed to load products', 'error')
        }
      } catch (error) {
        console.error('Error loading result:', error)
        window.showCustomAlert('Connection error while loading products', 'error')
      }
    },
    
    closeResult() {
      this.showResult = false
      this.resultData = []
    },
    
    closeDraft() {
      this.showDraft = false
    },
    
    async deleteDraft(id) {
      try {
        const response = await fetch(`${this.API_URL}/products/draft/${id}`, { method: 'DELETE' })
        if (response.ok) {
          window.showCustomAlert('Deleted!', 'success')
          await this.loadDraft()
        }
      } catch (error) {
        console.error('Error:', error)
      }
    },
    
    resetForm() {
      this.form = {
        code: '',
        code90Days: '',
        code180Days: '',
        insurance_code: this.selectedInsuranceCode || '',
        region: '',
        type: '',
        product_category: '', // Product Category
        coverage_type: '', // Coverage Type
        insurance_type: '', // Hardcopy Admin Fee
        protection_type: '',
        name: '',
        logo: '',
        is_active: '1',
        addons: []
      }
      this.selectedAddons = []
      this.$nextTick(() => {
        this.autoFillLogo()
        this.generateProductCode()
      })
    },
    
    async openAddonsModal() {
      this.showAddonsModal = true
      // Copy current selection from form or editForm depending on which is active
      if (this.showEditForm) {
        this.selectedAddons = [...(this.editForm.addons || [])]
      } else {
        this.selectedAddons = [...this.form.addons] // Copy current selection
      }
      await this.loadAddons()
    },
    
    closeAddonsModal() {
      this.showAddonsModal = false
    },
    
    async loadAddons() {
      this.loadingAddons = true
      try {
        // Define the addons that should be displayed (in order)
        const requiredAddons = [
          { code: 'COVID-PROTECTION', name: 'Covid Protection' },
          { code: 'TRAVEL-IBADAH', name: 'Travel Ibadah' },
          { code: 'WINTER-SPORT', name: 'Winter Sports' },
          { code: 'TRAVEL-PENDIDIKAN', name: 'Travel Pendidikan' },
          { code: 'TRIP-CANCELLATION-EXTENSION', name: 'Trip Cancellation Extension' },
          { code: 'GOLF-BENEFIT', name: 'Golf Benefit' },
          { code: 'LOSS-OF-HOTEL-RESERVATION', name: 'Loss of Hotel Reservation' },
          { code: 'PET-INSURANCE', name: 'Pet insurance' },
          { code: 'MISS-EVENT', name: 'Missed event' },
          { code: 'ADVENTURE', name: 'Adventures Activities Cover' },
          { code: 'VISA-PROTECTION', name: 'Visa Protection' },
          { code: 'AMATEUR-SPORT-EVENT', name: 'Amateur Sports Event' },
          { code: 'BUSINESS-TRAVEL', name: 'Business Travel' },
          { code: 'SNOW-SPORT', name: 'Snow Sport Cover' }
        ]
        
        // Try to load addons from API
        try {
          const response = await fetch(`${this.API_URL}/addons`)
          if (response.ok) {
            const allAddons = await response.json() || []
            const addonMap = new Map()
            
            // Create a map of addons from database
            allAddons.forEach(addon => {
              if (addon && addon.code) {
                addonMap.set(addon.code.toUpperCase(), addon)
              }
            })
            
            // Build available addons list, using database data if available, otherwise use default
            this.availableAddons = requiredAddons.map(reqAddon => {
              const dbAddon = addonMap.get(reqAddon.code.toUpperCase())
              if (dbAddon) {
                // Use data from database but ensure name matches
                return {
                  ...dbAddon,
                  name: reqAddon.name // Use the exact name from requirement
                }
              } else {
                // If not in database, create a default addon object
                return {
                  id: null,
                  code: reqAddon.code,
                  name: reqAddon.name,
                  is_active: 1
                }
              }
            })
          } else {
            // If API fails, use hardcoded list
            this.availableAddons = requiredAddons.map(reqAddon => ({
              id: null,
              code: reqAddon.code,
              name: reqAddon.name,
              is_active: 1
            }))
          }
        } catch (apiError) {
          console.error('Error loading addons from API:', apiError)
          // If API fails, use hardcoded list
          this.availableAddons = requiredAddons.map(reqAddon => ({
            id: null,
            code: reqAddon.code,
            name: reqAddon.name,
            is_active: 1
          }))
        }
      } catch (error) {
        console.error('Error in loadAddons:', error)
        window.showCustomAlert('Error loading addons', 'error')
        this.availableAddons = []
      } finally {
        this.loadingAddons = false
      }
    },
    
    selectAllAddons() {
      this.selectedAddons = this.availableAddons
        .filter(addon => addon.is_active === 1) // Only select active addons
        .map(addon => addon.code)
    },
    
    deselectAllAddons() {
      this.selectedAddons = []
    },
    
    saveSelectedAddons() {
      // Save to form or editForm depending on which is active
      if (this.showEditForm) {
        this.editForm.addons = [...this.selectedAddons]
        console.log('Addons saved to editForm:', this.editForm.addons)
      } else {
        this.form.addons = [...this.selectedAddons]
        console.log('Addons saved to form:', this.form.addons)
      }
      this.closeAddonsModal()
      window.showCustomAlert(`Selected ${this.selectedAddons.length} addon(s)`, 'success')
    },
    
    getAddonName(addonCode) {
      const addon = this.availableAddons.find(a => a.code === addonCode)
      return addon ? addon.name : addonCode
    },
    
    startEditDraft(item) {
      // Load draft item into edit form
      this.editForm = {
        id: item.id,
        code: item.code || '',
        insurance_code: item.insurance_code || '',
        product_category: item.property_category || '', // Product Category
        coverage_type: item.coverage_type || '', // Coverage Type
        region: item.region || '',
        type: item.type || '',
        insurance_type: item.insurance_type || '',
        protection_type: item.protection_type || '',
        name: item.name || '',
        logo: item.logo || '',
        is_active: String(item.is_active || '1'),
        addons: item.addons || []
      }
      this.isEditingDraft = true
      this.showEditForm = true
      this.showDraft = false
    },
    
    async startEditResult(item) {
      // Load result item into edit form for Property products
      // Get logo from insurance if not available
      let logoToEdit = item.logo || ''
      if (!logoToEdit && item.insurance_code && this.insurances.length > 0) {
        const insurance = this.insurances.find(ins => ins.code === item.insurance_code)
        if (insurance && insurance.logo) {
          logoToEdit = insurance.logo
        }
      }
      
      // Load hardcopy admin fee - check for null/undefined, not falsy (0 is valid)
      // Try admin_fee first, then hardcopy_admin_fee, then insurance_type as fallback
      const hardcopyAdminFee = (item.admin_fee != null && item.admin_fee !== undefined) 
        ? item.admin_fee 
        : ((item.hardcopy_admin_fee != null && item.hardcopy_admin_fee !== undefined) 
          ? item.hardcopy_admin_fee 
          : ((item.insurance_type != null && item.insurance_type !== undefined) ? item.insurance_type : ''))
      
      this.editForm = {
        id: item.id,
        code: item.code || '',
        insurance_code: item.insurance_code || '',
        product_category: item.property_category || '', // Product Category
        coverage_type: item.coverage_type || '', // Coverage Type
        insurance_type: hardcopyAdminFee, // Hardcopy Admin Fee
        name: item.name || '',
        logo: logoToEdit,
        is_active: String(item.is_active || '1')
      }
      
      console.log('Edit form loaded - Hardcopy Admin Fee:', hardcopyAdminFee, 'from item.admin_fee:', item.admin_fee, 'item.hardcopy_admin_fee:', item.hardcopy_admin_fee, 'full item:', item)
      console.log('Edit form loaded for result item:', this.editForm)
      this.isEditingDraft = false
      this.showEditForm = true
      this.showResult = false
    },
    
    async updateProduct() {
      try {
        // Region and Type removed for Property program - no region_id, min_adult, max_adult, min_child, max_child needed
        
        // Convert code: replace "-" with "_" for SUMMARY, INSURANCE_DETAIL, PROTECTION_DETAIL, HOW_TO_CLAIM
        const codeWithUnderscore = this.editForm.code.replace(/-/g, '_')
        const summary = `SUMMARY_${codeWithUnderscore}`
        const insuranceDetail = `INSURANCE_DETAIL_${codeWithUnderscore}`
        const protectionDetail = `PROTECTION_DETAIL_${codeWithUnderscore}`
        const howToClaim = `HOW_TO_CLAIM_${codeWithUnderscore}`
        
        // Get logo from insurance if not set
        let logoToUpdate = this.editForm.logo || ''
        if (!logoToUpdate && this.editForm.insurance_code && this.insurances.length > 0) {
          const insurance = this.insurances.find(ins => ins.code === this.editForm.insurance_code)
          if (insurance && insurance.logo) {
            logoToUpdate = insurance.logo
          }
        }
        
        // Parse hardcopy admin fee - ensure it's a valid integer
        const hardcopyAdminFee = this.editForm.insurance_type ? parseInt(this.editForm.insurance_type, 10) : 0
        if (isNaN(hardcopyAdminFee)) {
          window.showCustomAlert('Hardcopy Admin Fee harus berupa angka bulat yang valid', 'error')
          return
        }
        
        const payload = {
          code: this.editForm.code,
          insurance_code: this.editForm.insurance_code,
          name: this.editForm.name,
          logo: logoToUpdate,
          is_active: parseInt(this.editForm.is_active) || 1,
          admin_fee: hardcopyAdminFee, // Hardcopy Admin Fee - will be saved to hardcopy_admin_fee column
          property_category: this.editForm.product_category || '', // Product Category
          coverage_type: this.editForm.coverage_type || '', // Coverage Type
          min_coverage: 0,
          max_coverage: 0,
          summary: summary,
          insurance_detail: insuranceDetail,
          protection_detail: protectionDetail,
          how_to_claim: howToClaim
        }
        
        console.log('Updating product with payload:', payload)
        console.log('Hardcopy Admin Fee value:', hardcopyAdminFee)
        
        if (this.isEditingDraft) {
          // Update draft
          const response = await fetch(`${this.API_URL}/products/draft/${this.editForm.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          })
          
          if (response.ok) {
            window.showCustomAlert('Product draft updated successfully!', 'success')
            this.cancelEdit()
            await this.loadDraft()
            this.showDraft = true
          } else {
            const errorData = await response.json().catch(() => ({}))
            window.showCustomAlert(errorData.error || 'Failed to update draft', 'error')
          }
        } else {
          // Update database record
          const response = await fetch(`${this.API_URL}/products/${this.editForm.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          })
          
          if (response.ok) {
            window.showCustomAlert('Product updated successfully!', 'success')
            this.cancelEdit()
            await this.loadResult()
            this.showResult = true
          } else {
            const errorData = await response.json().catch(() => ({}))
            window.showCustomAlert(errorData.error || 'Failed to update product', 'error')
          }
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      }
    },
    
    async saveProductAddonsForEdit(productCode, insuranceCode, addonCodes) {
      try {
        const addonNameMap = this.getAddonNameMapping()
        const addonsWithNames = addonCodes.map(addonCode => {
          const addonName = addonNameMap[addonCode] || addonCode
          return {
            code: addonCode,
            name: addonName,
            name_my: addonName,
            name_en: addonName
          }
        })
        
        const payload = {
          product_code: productCode,
          insurance_code: insuranceCode,
          addons: addonsWithNames
        }
        
        const response = await fetch(`${this.API_URL}/products/addons`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
        
        if (!response.ok) {
          const errorText = await response.text()
          console.error(`Failed to save addons for product ${productCode}:`, errorText)
        } else {
          console.log(`Successfully saved addons for product ${productCode}`)
        }
      } catch (error) {
        console.error('Error saving addons for edit:', error)
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      this.isEditingDraft = false
      this.editForm = {
        id: null,
        code: '',
        code90Days: '',
        code180Days: '',
        insurance_code: '',
        product_category: '', // Product Category
        coverage_type: '', // Coverage Type
        region: '',
        type: '',
        insurance_type: '',
        protection_type: '',
        name: '',
        logo: '',
        is_active: '1',
        addons: []
      }
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1050;
}

.modal-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 90%;
  max-width: 600px;
}

.modal-header {
  padding: 1rem;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  opacity: 0.7;
}

.modal-body {
  padding: 1rem;
}

.modal-footer {
  padding: 1rem;
  border-top: 1px solid #dee2e6;
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.addons-list {
  max-height: 400px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
}

@media (max-width: 768px) {
  .addons-list {
    grid-template-columns: 1fr;
  }
}

.addon-item {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  margin-bottom: 0;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  transition: background-color 0.2s;
  gap: 0.75rem;
}

.addon-item:hover {
  background-color: #f8f9fa;
}

.addon-checkbox {
  width: 20px;
  height: 20px;
  min-width: 20px;
  min-height: 20px;
  margin: 0;
  cursor: pointer;
  flex-shrink: 0;
}

.addon-label {
  flex: 1;
  margin: 0;
  padding: 0;
  cursor: pointer;
  user-select: none;
  line-height: 1.5;
}

.addon-checkbox:checked + .addon-label {
  color: #0d6efd;
  font-weight: 500;
}
</style>

