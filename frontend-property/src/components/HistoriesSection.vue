<template>
  <div class="col-12">
    <div class="form-section">
      <h3><i class="fas fa-history me-1"></i> Histories Management</h3>
      
      <div class="row mb-3">
        <div class="col-md-12">
          <div class="d-flex gap-2 align-items-center flex-wrap">
            <div class="filter-section">
              <label class="form-label me-2 mb-0">Filter by Section:</label>
              <select class="form-select d-inline-block" style="width: auto; min-width: 200px;" v-model="selectedSection" @change="loadHistories">
                <option value="">All Sections</option>
                <option v-for="section in availableSections" :key="section" :value="section">
                  {{ section }}
                </option>
              </select>
            </div>
            <button class="btn btn-info" @click="loadHistories">
              <i class="fas fa-sync-alt"></i> Load Histories
            </button>
            <button class="btn btn-secondary" @click="hideHistories">
              <i class="fas fa-times"></i> Close
            </button>
          </div>
        </div>
      </div>

      <!-- Histories Table Section -->
      <div v-if="showHistories" class="table-section result-section show">
        <div class="table-header d-flex justify-content-between align-items-center">
          <h4>📜 History Records</h4>
          <div>
            <button class="btn btn-info btn-sm me-2" @click="refreshHistories">
              <i class="fas fa-sync-alt"></i> Refresh Data
            </button>
          </div>
        </div>

        <div class="table-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead class="table-success">
                <tr>
                  <th>ID</th>
                  <th>User Name</th>
                  <th>Section</th>
                  <th>Action</th>
                  <th>Record ID</th>
                  <th>Record Type</th>
                  <th>Data Before</th>
                  <th>Data After</th>
                  <th>Created At</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="histories.length === 0">
                  <td colspan="9" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No history records found
                  </td>
                </tr>
                <tr v-for="item in histories" :key="item.id">
                  <td>{{ item.id }}</td>
                  <td>{{ item.user_name || '-' }}</td>
                  <td>{{ item.section || '-' }}</td>
                  <td>
                    <span class="badge bg-info">{{ item.action || '-' }}</span>
                  </td>
                  <td>{{ item.record_id || '-' }}</td>
                  <td>{{ item.record_type || '-' }}</td>
                  <td>
                    <div v-if="item.data_before" class="data-preview">
                      <pre class="mb-0">{{ formatJSON(item.data_before) }}</pre>
                    </div>
                    <span v-else class="text-muted">-</span>
                  </td>
                  <td>
                    <div v-if="item.data_after" class="data-preview">
                      <pre class="mb-0">{{ formatJSON(item.data_after) }}</pre>
                    </div>
                    <span v-else class="text-muted">-</span>
                  </td>
                  <td>{{ formatDateOnly(item.created_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HistoriesSection',
  data() {
    return {
      histories: [],
      showHistories: false,
      selectedSection: '',
      availableSections: [],
      API_URL: 'http://localhost:8080/api/property'
    }
  },
  methods: {
    async loadHistories() {
      try {
        // Load all histories first to populate sections dropdown if empty
        if (this.availableSections.length === 0) {
          await this.loadAllHistoriesForSections()
        }
        
        await this.fetchHistories()
        this.showHistories = true
        // Removed scroll behavior - keep scroll position at current location
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to load histories: ' + (error.message || 'Unknown error'), 'error')
      }
    },
    
    async loadAllHistoriesForSections() {
      // Load all histories first to get available sections
      try {
        const response = await fetch(`${this.API_URL}/histories`)
        if (response.ok) {
          const data = await response.json()
          const allData = Array.isArray(data) ? data : []
          
          // Extract unique sections
          const sections = new Set()
          allData.forEach(h => {
            if (h.section) {
              sections.add(h.section)
            }
          })
          this.availableSections = Array.from(sections).sort()
        }
      } catch (error) {
        console.error('Error loading sections:', error)
      }
    },
    
    async fetchHistories() {
      try {
        let url = `${this.API_URL}/histories`
        if (this.selectedSection) {
          url = `${this.API_URL}/histories/table/${encodeURIComponent(this.selectedSection)}`
        }
        
        const response = await fetch(url)
        if (response.ok) {
          const data = await response.json()
          this.histories = Array.isArray(data) ? data : []
          
          // If loading all histories, extract unique sections for dropdown
          if (!this.selectedSection) {
            const sections = new Set()
            this.histories.forEach(h => {
              if (h.section) {
                sections.add(h.section)
              }
            })
            this.availableSections = Array.from(sections).sort()
          }
        } else {
          const error = await response.json()
          throw new Error(error.error || 'Failed to load histories')
        }
      } catch (error) {
        console.error('Error:', error)
        this.histories = []
        throw error
      }
    },
    
    async refreshHistories() {
      await this.fetchHistories()
      window.showCustomAlert('Histories refreshed successfully!', 'success')
    },
    
    hideHistories() {
      this.showHistories = false
    },
    
    formatJSON(jsonString) {
      if (!jsonString) return '-'
      try {
        const parsed = JSON.parse(jsonString)
        return JSON.stringify(parsed, null, 2)
      } catch (e) {
        return jsonString
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
.data-preview {
  max-width: 300px;
  max-height: 150px;
  overflow: auto;
  font-size: 0.85rem;
}

.data-preview pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', monospace;
  background: #f8f9fa;
  padding: 8px;
  border-radius: 4px;
  margin: 0;
}
</style>
