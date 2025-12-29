<template>
  <div class="col-12">
    <div class="form-section">
      <h3><i class="fas fa-download me-1"></i> ALL RESULT</h3>
      
      <div class="row">
        <div class="col-md-12">
          <div class="d-flex gap-2 mb-3">
            <button class="btn btn-success btn-lg" @click="downloadAllData" :disabled="downloading">
              <i class="fas fa-download"></i> 
              {{ downloading ? 'Downloading...' : 'Download All Data (Excel)' }}
            </button>
          </div>
                 <p class="text-muted">
                   Download data yang sudah masuk ke database dari section Product Rules, Addon Rule Details, dan Templates
                 </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AllResultSection',
  data() {
    return {
      API_URL: 'http://localhost:8080/api/property',
      downloading: false
    }
  },
  methods: {
    async downloadAllData() {
      try {
        this.downloading = true
        const response = await fetch(`${this.API_URL}/all-result/download`)
        
        if (!response.ok) {
          const errorText = await response.text()
          window.showCustomAlert('Failed to download all data: ' + errorText, 'error')
          return
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        
        // Get filename from Content-Disposition header or use default
        const contentDisposition = response.headers.get('Content-Disposition')
        let filename = `all_result_data_${new Date().toISOString().split('T')[0]}.xlsx`
        if (contentDisposition) {
          const filenameMatch = contentDisposition.match(/filename="?(.+)"?/i)
          if (filenameMatch) {
            filename = filenameMatch[1]
          }
        }
        
        a.download = filename
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
        window.showCustomAlert('All data downloaded successfully!', 'success')
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error: ' + error.message, 'error')
      } finally {
        this.downloading = false
      }
    }
  }
}
</script>

<style scoped>
.btn-lg {
  padding: 0.75rem 1.5rem;
  font-size: 1.1rem;
}

.text-muted {
  color: #6c757d;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}
</style>

