<template>
  <div class="col-12">
    <!-- Form Section -->
    <div class="form-section">
      <div class="d-flex justify-content-between align-items-center">
        <h3><i class="fas fa-plus-circle me-1"></i> Tambah Data Asuransi Baru</h3>
      </div>
      
      <form @submit.prevent="saveInsurance">
        <div class="form-grid form-grid-4">
          <div>
            <label for="ins-code" class="form-label">Code Asuransi <span class="text-danger">*</span></label>
            <input 
              type="text" 
              class="form-control" 
              id="ins-code" 
              v-model="form.code" 
              placeholder="Cth: TUGU" 
              required
            >
          </div>
          <div>
            <label for="ins-name" class="form-label">Nama Asuransi <span class="text-danger">*</span></label>
            <input 
              type="text" 
              class="form-control" 
              id="ins-name" 
              v-model="form.name" 
              placeholder="Cth: Tugu Insurance" 
              required
            >
          </div>
          <div>
            <label for="ins-status" class="form-label">Status <span class="text-danger">*</span></label>
            <select class="form-select" id="ins-status" v-model="form.status" required>
              <option value="1">1 (Active)</option>
              <option value="0">0 (Inactive)</option>
            </select>
          </div>
          <div>
            <label for="ins-logo" class="form-label">Logo URL (Optional)</label>
            <input 
              type="url" 
              class="form-control" 
              id="ins-logo" 
              v-model="form.logo" 
              placeholder="Cth: https://cdn.com/tugu.png"
            >
          </div>
        </div>
        
        <div class="mt-3 d-flex align-items-center gap-2">
          <button type="button" class="btn btn-info" @click="showDraft">
            <i class="fas fa-eye"></i> Tampilkan Draft Data
          </button>
          <button type="button" class="btn btn-warning" @click="clearDraft">
            <i class="fas fa-trash"></i> Clear Draft
          </button>
          <button type="submit" class="btn btn-primary">
            <i class="fas fa-save"></i> Save
          </button>
          <button type="button" class="btn btn-info" @click="showResult">
            <i class="fas fa-chart-line"></i> RESULT
          </button>
        </div>
      </form>
    </div>

    <!-- Edit Form Section -->
    <div v-if="showEditForm" class="form-section">
      <h3><i class="fas fa-edit me-1"></i> Edit Data Asuransi</h3>
      
      <form @submit.prevent="updateInsurance">
        <input type="hidden" v-model="editForm.id">
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Code Asuransi <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.code" required>
          </div>
          <div>
            <label class="form-label">Nama Asuransi <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.name" required>
          </div>
          <div>
            <label class="form-label">Status <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.status" required>
              <option value="1">1 (Active)</option>
              <option value="0">0 (Inactive)</option>
            </select>
          </div>
          <div>
            <label class="form-label">Logo URL (Optional)</label>
            <input type="url" class="form-control" v-model="editForm.logo">
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

    <!-- Draft Table Section -->
    <div v-if="showDraftSection" class="table-section draft-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📊 Data Preview: Insurances (Draft Save)</h4>
        <div>
        <button 
            class="btn btn-success btn-confirm me-2" 
          :disabled="!canConfirm" 
          @click="confirmInsurances"
        >
          <i class="fas fa-database me-2"></i> Konfirmasi Data (DB)
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
                <th>ID</th>
                <th>Code</th>
                <th>Name</th>
                <th>Logo</th>
                <th>Status</th>
                <th>Timestamp</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in draftData" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.code }}</td>
                <td>{{ item.name }}</td>
                <td>
                  <img v-if="item.logo" :src="item.logo" class="table-logo" alt="Logo">
                  <span v-else>-</span>
                </td>
                <td>
                  <span :class="getStatusBadge(item.status).class">
                    {{ getStatusBadge(item.status).text }}
                  </span>
                </td>
                <td>{{ formatTimestamp(item.timestamp) }}</td>
                <td>
                  <button class="btn btn-sm btn-warning me-1" @click="startEdit(item)">
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

    <!-- Result Table Section -->
    <div v-if="showResultSection" class="table-section result-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📈 Data Result: Insurances (From Database)</h4>
        <div>
          <button class="btn btn-info btn-sm me-2" @click="refreshResult">
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
                <th>Code</th>
                <th>Name</th>
                <th>Logo</th>
                <th>Status</th>
                <th>Created At</th>
                <th>Updated At</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in resultData" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.code }}</td>
                <td>{{ item.name }}</td>
                <td>
                  <img v-if="item.logo" :src="item.logo" class="table-logo" alt="Logo">
                  <span v-else>-</span>
                </td>
                <td>
                  <span :class="getStatusBadge(item.status).class">
                    {{ getStatusBadge(item.status).text }}
                  </span>
                </td>
                <td>{{ formatDateOnly(item.created_at) }}</td>
                <td>{{ item.updated_at }}</td>
                <td>
                  <button class="btn btn-sm btn-warning me-1" @click="startEditFromResult(item)">
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
  name: 'InsuranceSection',
  props: {
    selectedInsuranceCode: String
  },
  watch: {
    selectedInsuranceCode(newVal) {
      // Reload result when insurance code changes
      if (this.showResultSection) {
        this.loadResult()
      }
    }
  },
  data() {
    return {
      form: {
        code: '',
        name: '',
        status: '1',
        logo: ''
      },
      editForm: {
        id: null,
        code: '',
        name: '',
        status: '1',
        logo: ''
      },
      showEditForm: false,
      showDraftSection: false,
      showResultSection: false,
      draftData: [],
      resultData: [],
      canConfirm: false,
      API_URL: 'http://localhost:8080/api/travel' // Travel program backend
    }
  },
  methods: {
    async saveInsurance() {
      try {
        // Convert status to integer
        const payload = {
          ...this.form,
          status: parseInt(this.form.status) || 1
        }
        
        const response = await fetch(`${this.API_URL}/insurances/draft/add`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({ error: 'Failed to save' }))
          window.showCustomAlert(errorData.error || 'Failed to save', 'error')
          return
        }
        
        const result = await response.json()
        
        window.showCustomAlert('Insurance draft saved successfully!', 'success')
        this.resetForm()
        // Load draft and show section
        await this.loadDraft()
        this.showDraftSection = true
        this.canConfirm = this.draftData.length > 0
        this.$nextTick(() => {
          // Scroll to draft section
          const draftSection = document.querySelector('.draft-section')
          if (draftSection) {
            draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
          }
        })
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      }
    },
    
    async loadDraft() {
      try {
        const response = await fetch(`${this.API_URL}/insurances/draft`)
        if (response.ok) {
          this.draftData = await response.json()
          this.canConfirm = this.draftData.length > 0
        } else {
          console.error('Failed to load draft:', response.status)
          this.draftData = []
          this.canConfirm = false
        }
      } catch (error) {
        console.error('Error loading draft:', error)
        this.draftData = []
        this.canConfirm = false
      }
    },
    
    async showDraft() {
      try {
        await this.loadDraft()
        this.showDraftSection = true
        this.canConfirm = this.draftData.length > 0
        if (this.draftData.length === 0) {
          window.showCustomAlert('No draft data available', 'info')
        } else {
          // Scroll to draft section
          this.$nextTick(() => {
            setTimeout(() => {
              const draftSection = document.querySelector('.draft-section')
              if (draftSection) {
                draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
              }
            }, 100)
          })
        }
      } catch (error) {
        console.error('Error showing draft:', error)
        window.showCustomAlert('Failed to load draft data: ' + error.message, 'error')
      }
    },
    
    hideDraft() {
      this.showDraftSection = false
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Are you sure you want to clear all draft data?')
      if (!confirmed) return
      
      try {
        const response = await fetch(`${this.API_URL}/insurances/draft/clear`, {
          method: 'POST'
        })
        
        if (response.ok) {
          this.draftData = []
          this.canConfirm = false
          this.showDraftSection = false
          window.showCustomAlert('Draft cleared successfully!', 'success')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to clear draft', 'error')
      }
    },
    
    async confirmInsurances() {
      try {
        const response = await fetch(`${this.API_URL}/insurances/draft/confirm`, {
          method: 'POST'
        })
        
        const result = await response.json()
        
        if (response.ok) {
          window.showCustomAlert('Data confirmed successfully!', 'success')
          this.draftData = []
          this.canConfirm = false
          this.showDraftSection = false
          // Emit event to parent to refresh insurance dropdown
          this.$emit('insurance-confirmed')
        } else {
          window.showCustomAlert(result.error || 'Failed to confirm', 'error')
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    async showResult() {
      await this.loadResult()
      this.showResultSection = true
    },
    
    async loadResult() {
      try {
        const response = await fetch(`${this.API_URL}/insurances`)
        if (response.ok) {
          let data = await response.json()
          // Filter by selectedInsuranceCode if provided
          if (this.selectedInsuranceCode) {
            data = data.filter(item => item.code === this.selectedInsuranceCode)
          }
          this.resultData = data
        } else {
          console.error('Failed to load insurances:', response.status)
          this.resultData = []
        }
      } catch (error) {
        console.error('Error loading result:', error)
        window.showCustomAlert('Connection error while loading insurances', 'error')
        this.resultData = []
      }
    },
    
    async refreshResult() {
      await this.loadResult()
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    async deleteDraft(id) {
      try {
        const response = await fetch(`${this.API_URL}/insurances/draft/${id}`, {
          method: 'DELETE'
        })
        
        if (response.ok) {
          window.showCustomAlert('Item deleted successfully!', 'success')
          await this.loadDraft()
          if (this.draftData.length === 0) {
            this.showDraftSection = false
          }
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to delete', 'error')
      }
    },
    
    startEdit(item) {
      // Convert status from int to string for form
      this.editForm = {
        id: item.id,
        code: item.code,
        name: item.name,
        status: String(item.status || 1),
        logo: item.logo || ''
      }
      this.showEditForm = true
      this.showDraftSection = false
    },
    
    startEditFromResult(item) {
      // Convert status from int to string for form
      this.editForm = {
        id: item.id,
        code: item.code,
        name: item.name,
        status: String(item.status || 1),
        logo: item.logo || ''
      }
      this.showEditForm = true
      this.showResultSection = false
    },
    
    async updateInsurance() {
      try {
        // Convert status to integer
        const payload = {
          code: this.editForm.code,
          name: this.editForm.name,
          status: parseInt(this.editForm.status) || 1,
          logo: this.editForm.logo || ''
        }
        
        // Check if editing draft or database record
        // If editing from draft, the item will have a timestamp field
        // If editing from result, it will have created_at and updated_at
        const isDraft = this.draftData.some(d => d.id === this.editForm.id)
        
        if (isDraft) {
          // Update draft
          const response = await fetch(`${this.API_URL}/insurances/draft/${this.editForm.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          })
          
          if (response.ok) {
            window.showCustomAlert('Insurance draft updated successfully!', 'success')
            this.cancelEdit()
            await this.loadDraft()
            this.showDraftSection = true
          } else {
            const result = await response.json()
            window.showCustomAlert(result.error || 'Failed to update', 'error')
          }
        } else {
          // Update database record
          const response = await fetch(`${this.API_URL}/insurances/${this.editForm.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
          })
          
          if (response.ok) {
            window.showCustomAlert('Insurance updated successfully!', 'success')
            this.cancelEdit()
            await this.loadResult()
            this.showResultSection = true
            // Emit event to parent to refresh insurance dropdown
            this.$emit('insurance-confirmed')
          } else {
            const result = await response.json()
            window.showCustomAlert(result.error || 'Failed to update', 'error')
          }
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      this.editForm = { id: null, code: '', name: '', status: '1', logo: '' }
    },
    
    resetForm() {
      this.form = { code: '', name: '', status: '1', logo: '' }
    },
    
    getStatusBadge(status) {
      return {
        text: status == 1 ? 'Active' : 'Inactive',
        class: `badge ${status == 1 ? 'bg-success' : 'bg-secondary'}`
      }
    },
    
    formatTimestamp(timestamp) {
      if (!timestamp) return '-'
      return new Date(timestamp).toLocaleString()
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

