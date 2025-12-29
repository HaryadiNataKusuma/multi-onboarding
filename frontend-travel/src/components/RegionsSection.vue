<template>
  <div class="col-12">
    <!-- Form Section -->
    <div class="form-section">
      <div class="d-flex justify-content-between align-items-center">
        <h3><i class="fas fa-globe me-1"></i> Tambah Data Region Baru</h3>
      </div>
      
      <form @submit.prevent="saveRegion">
        <div class="form-grid form-grid-3">
          <div>
            <label for="region-name" class="form-label">Region Name <span class="text-danger">*</span></label>
            <input 
              type="text" 
              class="form-control" 
              id="region-name" 
              v-model="form.name" 
              placeholder="Cth: Asia Pacific" 
              required
            >
          </div>
          
          <div>
            <label for="region-type" class="form-label">Type <span class="text-danger">*</span></label>
            <select class="form-select" id="region-type" v-model="form.type" required>
              <option value="">Select type...</option>
              <option value="WHITELIST">WHITELIST</option>
              <option value="BLACKLIST">BLACKLIST</option>
            </select>
          </div>
          
          <div>
            <label for="region-country" class="form-label">Countries <span class="text-danger">*</span></label>
            <button 
              type="button" 
              :class="['btn btn-countries-fixed w-100', form.country_ids.length > 0 ? 'btn-primary' : 'btn-outline-primary']"
              @click="openCountryModal"
            >
              <i class="fas fa-globe me-2"></i>
              <span class="btn-countries-text">
                {{ form.country_ids.length > 0 ? `${form.country_ids.length} countries selected` : 'Select Countries' }}
              </span>
            </button>
            <div v-if="form.country_ids.length > 0" class="mt-2">
              <small class="form-text text-muted">
                Country IDs: {{ JSON.stringify(form.country_ids.slice(0, 10)) }}{{ form.country_ids.length > 10 ? '...' : '' }}
              </small>
            </div>
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
      <h3><i class="fas fa-edit me-1"></i> Edit Data Region</h3>
      
      <form @submit.prevent="updateRegion">
        <input type="hidden" v-model="editForm.id">
        <div class="form-grid form-grid-3">
          <div>
            <label class="form-label">Region Name <span class="text-danger">*</span></label>
            <input type="text" class="form-control" v-model="editForm.name" required>
          </div>
          <div>
            <label class="form-label">Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.type" required>
              <option value="WHITELIST">WHITELIST</option>
              <option value="BLACKLIST">BLACKLIST</option>
            </select>
          </div>
          <div>
            <label class="form-label">Countries <span class="text-danger">*</span></label>
            <button 
              type="button" 
              :class="['btn btn-countries-fixed w-100', editForm.country_ids.length > 0 ? 'btn-primary' : 'btn-outline-primary']"
              @click="openCountryModalForEdit"
            >
              <i class="fas fa-globe me-2"></i>
              <span class="btn-countries-text">
                {{ editForm.country_ids.length > 0 ? `${editForm.country_ids.length} countries selected` : 'Select Countries' }}
              </span>
            </button>
            <div v-if="editForm.country_ids.length > 0" class="mt-2">
              <small class="form-text text-muted">
                Country IDs: {{ JSON.stringify(editForm.country_ids.slice(0, 10)) }}{{ editForm.country_ids.length > 10 ? '...' : '' }}
              </small>
            </div>
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
        <h4>📊 Data Preview: Regions (Draft Save)</h4>
        <div>
          <button 
            class="btn btn-success btn-confirm me-2" 
            :disabled="!canConfirm" 
            @click="confirmRegions"
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
                <th>Region Name</th>
                <th>Type</th>
                <th>Country</th>
                <th>Timestamp</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in draftData" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.name }}</td>
                <td>
                  <span :class="getTypeBadge(item.type).class">
                    {{ item.type }}
                  </span>
                </td>
                <td>
                  <span v-for="(countryId, idx) in item.country_ids" :key="countryId">
                    {{ getCountryName(countryId) }}<span v-if="idx < item.country_ids.length - 1">, </span>
                  </span>
                </td>
                <td>{{ formatDateOnly(item.timestamp || item.created_at) }}</td>
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
        <h4>📈 Data Result: Regions (From Database)</h4>
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
                <th>Region Name</th>
                <th>Type</th>
                <th>Country</th>
                <th>Created At</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in resultData" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ item.name }}</td>
                <td>
                  <span :class="getTypeBadge(item.type).class">
                    {{ item.type }}
                  </span>
                </td>
                <td>
                  <span v-for="(countryId, idx) in item.country_ids" :key="countryId">
                    {{ getCountryName(countryId) }}<span v-if="idx < item.country_ids.length - 1">, </span>
                  </span>
                </td>
                <td>{{ formatDateOnly(item.created_at) }}</td>
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
    
    <!-- Country Select Modal -->
    <CountrySelectModal 
      :show="showCountryModal"
      :countries="countries"
      :selected-ids="editingCountryIds.length > 0 ? editingCountryIds : form.country_ids"
      @close="closeCountryModal"
      @confirm="handleCountrySelect"
    />
  </div>
</template>

<script>
import { getRegions, createRegion, updateRegion, deleteRegion, getCountries } from '../services/api'
import CountrySelectModal from './CountrySelectModal.vue'

export default {
  name: 'RegionsSection',
  components: {
    CountrySelectModal
  },
  data() {
    return {
      form: {
        name: '',
        type: '',
        country_ids: []
      },
      editForm: {
        id: null,
        name: '',
        type: '',
        country_ids: []
      },
      showEditForm: false,
      showDraftSection: false,
      showResultSection: false,
      draftData: [],
      resultData: [],
      countries: [],
      canConfirm: false,
      showCountryModal: false,
      editingCountryIds: [],
      editingFromDraft: false,
      draftIdCounter: 1 // Counter for draft IDs
    }
  },
  async mounted() {
    await this.loadCountries()
  },
  methods: {
    async loadCountries() {
      try {
        this.countries = await getCountries()
        console.log('Countries loaded:', this.countries.length)
      } catch (error) {
        console.error('Error loading countries:', error)
        // Don't show alert on mount, only show if user tries to use the feature
        if (this.countries.length === 0) {
          console.warn('Countries not loaded yet, will retry when needed')
        }
      }
    },
    
    getCountryName(countryId) {
      const country = this.countries.find(c => c.id === countryId)
      return country ? country.name : `Country ID: ${countryId}`
    },
    
    openCountryModal() {
      // Ensure we start with clean state - filter out any invalid IDs
      const validIds = (this.form.country_ids || []).filter(id => id && id !== 0 && !isNaN(id))
      this.editingCountryIds = []
      this.form.country_ids = validIds
      this.showCountryModal = true
    },
    
    openCountryModalForEdit() {
      // Filter out invalid IDs
      const validIds = (this.editForm.country_ids || []).filter(id => id && id !== 0 && !isNaN(id))
      this.editingCountryIds = [...validIds]
      this.editForm.country_ids = validIds
      this.showCountryModal = true
    },
    
    closeCountryModal() {
      this.showCountryModal = false
      this.editingCountryIds = []
    },
    
    handleCountrySelect(selectedIds) {
      if (this.showEditForm) {
        this.editForm.country_ids = selectedIds
      } else {
        this.form.country_ids = selectedIds
      }
      this.closeCountryModal()
    },
    
    async saveRegion() {
      if (this.form.country_ids.length === 0) {
        window.showCustomAlert('Please select at least one country', 'error')
        return
      }
      
      try {
        // Save to draft (in-memory)
        const draftItem = {
          id: this.draftIdCounter++, // Simple incremental ID
          name: this.form.name,
          type: this.form.type,
          country_ids: [...this.form.country_ids],
          timestamp: new Date().toISOString()
        }
        
        this.draftData.push(draftItem)
        this.canConfirm = this.draftData.length > 0
        window.showCustomAlert('Region draft saved successfully!', 'success')
        this.resetForm()
        this.showDraftSection = true
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert(error.response?.data?.error || 'Failed to save region', 'error')
      }
    },
    
    async loadDraft() {
      // Draft is stored in-memory in draftData
      this.canConfirm = this.draftData.length > 0
    },
    
    async showDraft() {
      await this.loadDraft()
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data available', 'info')
        this.showDraftSection = false
        return
      }
      this.showDraftSection = true
      // Scroll to draft section
      this.$nextTick(() => {
        setTimeout(() => {
          const draftSection = document.querySelector('.draft-section')
          if (draftSection) {
            draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
          }
        }, 100)
      })
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Are you sure you want to clear all draft data?')
      if (!confirmed) return
      
      this.draftData = []
      this.canConfirm = false
      this.showDraftSection = false
      this.draftIdCounter = 1 // Reset counter
      window.showCustomAlert('Draft cleared successfully!', 'success')
    },
    
    hideDraft() {
      this.showDraftSection = false
    },
    
    async confirmRegions() {
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data to confirm', 'info')
        return
      }
      
      const confirmed = await window.showCustomConfirm('Are you sure you want to confirm all draft data to the database?')
      if (!confirmed) return
      
      try {
        // Save all draft items to database
        for (const item of this.draftData) {
          const regionData = {
            name: item.name,
            type: item.type,
            country_ids: item.country_ids
          }
          await createRegion(regionData)
        }
        
        window.showCustomAlert('Data confirmed successfully!', 'success')
        this.draftData = []
        this.canConfirm = false
        this.showDraftSection = false
        await this.loadResult()
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert(error.response?.data?.error || 'Failed to confirm regions', 'error')
      }
    },
    
    async showResult() {
      try {
        await this.loadResult()
        if (this.resultData.length === 0) {
          window.showCustomAlert('No regions found in database', 'info')
        }
        this.showResultSection = true
        // Scroll to result section
        this.$nextTick(() => {
          setTimeout(() => {
            const resultSection = document.querySelector('.result-section')
            if (resultSection) {
              resultSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
            }
          }, 100)
        })
      } catch (error) {
        console.error('Error in showResult:', error)
        window.showCustomAlert('Failed to load regions: ' + (error.message || 'Unknown error'), 'error')
        this.showResultSection = false
      }
    },
    
    async loadResult() {
      try {
        const response = await fetch('http://localhost:8080/api/travel/regions')
        if (response.ok) {
          const data = await response.json()
          this.resultData = data || []
          console.log('Regions loaded:', this.resultData.length)
        } else {
          const errorText = await response.text()
          console.error('Failed to load regions:', errorText)
          let errorMessage = 'Failed to load regions'
          try {
            const errorData = JSON.parse(errorText)
            errorMessage = errorData.error || errorMessage
          } catch (e) {
            errorMessage = errorText || errorMessage
          }
          window.showCustomAlert(errorMessage, 'error')
          throw new Error(errorMessage)
        }
      } catch (error) {
        console.error('Error loading result:', error)
        if (error.message) {
          throw error
        }
        throw new Error('Failed to load regions: ' + error.message)
      }
    },
    
    async refreshResult() {
      await this.loadResult()
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    async deleteDraft(id) {
      const confirmed = await window.showCustomConfirm('Are you sure you want to delete this draft item?')
      if (!confirmed) return
      
      try {
        // Remove from draft (in-memory)
        this.draftData = this.draftData.filter(d => d.id !== id)
        this.canConfirm = this.draftData.length > 0
        if (this.draftData.length === 0) {
          this.showDraftSection = false
        }
        window.showCustomAlert('Draft item deleted successfully!', 'success')
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to delete draft item', 'error')
      }
    },
    
    startEdit(item) {
      this.editForm = { ...item }
      this.editingCountryIds = [...(item.country_ids || [])]
      this.editForm.country_ids = [...(item.country_ids || [])]
      this.editingFromDraft = true // Flag to track if editing from draft
      this.showEditForm = true
      this.showDraftSection = false
    },
    
    startEditFromResult(item) {
      this.editForm = { ...item }
      this.editingCountryIds = [...(item.country_ids || [])]
      this.editForm.country_ids = [...(item.country_ids || [])]
      this.editingFromDraft = false // Flag to track if editing from result/database
      this.showEditForm = true
      this.showResultSection = false
    },
    
    async updateRegion() {
      if (this.editForm.country_ids.length === 0) {
        window.showCustomAlert('Please select at least one country', 'error')
        return
      }
      
      try {
        // Check if editing from draft or from database
        if (this.editingFromDraft) {
          // Update in draft (in-memory)
          const index = this.draftData.findIndex(d => d.id === this.editForm.id)
          if (index !== -1) {
            this.draftData[index] = {
              ...this.editForm,
              country_ids: [...this.editForm.country_ids],
              timestamp: this.draftData[index].timestamp || new Date().toISOString()
            }
            window.showCustomAlert('Region draft updated successfully!', 'success')
            this.cancelEdit()
            this.showDraftSection = true
            this.canConfirm = this.draftData.length > 0
          }
        } else {
          // Update in database
          const regionData = {
            name: this.editForm.name,
            type: this.editForm.type,
            country_ids: this.editForm.country_ids
          }
          
          await updateRegion(this.editForm.id, regionData)
          window.showCustomAlert('Region updated successfully!', 'success')
          this.cancelEdit()
          await this.loadResult()
          this.showResultSection = true
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert(error.response?.data?.error || 'Failed to update region', 'error')
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      this.editForm = { id: null, name: '', type: '', country_ids: [] }
      this.editingCountryIds = []
      this.editingFromDraft = false
    },
    
    resetForm() {
      this.form = { name: '', type: '', country_ids: [] }
    },
    
    getTypeBadge(type) {
      return {
        text: type,
        class: `badge ${type === 'WHITELIST' ? 'bg-success' : 'bg-danger'}`
      }
    },
    
    formatTimestamp(timestamp) {
      if (!timestamp) return '-'
      return new Date(timestamp).toLocaleString()
    },
    
    formatDateOnly(dateString) {
      if (!dateString) return '-'
      try {
        const date = new Date(dateString)
        if (isNaN(date.getTime())) return '-'
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, '0')
        const day = String(date.getDate()).padStart(2, '0')
        return `${year}-${month}-${day}`
      } catch (error) {
        console.error('Error formatting date:', error)
        return '-'
      }
    }
  }
}
</script>

<style scoped>
.btn-countries-fixed {
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

.btn-countries-fixed i {
  flex-shrink: 0;
}

.btn-countries-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
</style>

