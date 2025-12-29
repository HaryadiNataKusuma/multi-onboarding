<template>
  <div class="bg-light border-bottom p-3">
    <div class="container-fluid">
      <div class="row align-items-center">
        <div class="col-md-3">
          <label for="insurance-dropdown" class="form-label mb-0">
            <i class="fas fa-shield-alt me-2"></i><strong>List Asuransi:</strong>
          </label>
        </div>
        <div class="col-md-6">
          <select 
            class="form-select" 
            id="insurance-dropdown" 
            v-model="selectedCode"
            @change="onChange"
          >
            <option value="">Pilih Kode Asuransi...</option>
            <option v-for="insurance in insurances" :key="insurance.code" :value="insurance.code">
              {{ insurance.code }} - {{ insurance.name }}
            </option>
          </select>
        </div>
        <div class="col-md-3">
          <small class="text-muted">
            <i class="fas fa-info-circle me-1"></i>
            Memudahkan edit data yang sudah masuk ke database
          </small>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'InsuranceDropdown',
  data() {
    return {
      insurances: [],
      selectedCode: ''
    }
  },
  mounted() {
    this.loadInsurances()
  },
  methods: {
    async loadInsurances() {
      try {
        console.log('🔄 Loading insurances from http://localhost:8080/api/travel/insurances...')
        // Using Travel program backend
        const response = await fetch('http://localhost:8080/api/travel/insurances')
        console.log('📡 Response status:', response.status, response.statusText)
        
        if (response.ok) {
          const data = await response.json()
          console.log('📦 Received data:', data)
          console.log('📊 Data type:', Array.isArray(data) ? 'Array' : typeof data)
          console.log('📈 Data length:', Array.isArray(data) ? data.length : 'N/A')
          
          // Ensure data is an array
          if (Array.isArray(data)) {
            this.insurances = data
            console.log('✅ Insurances loaded successfully:', data.length, 'items')
            if (data.length > 0) {
              console.log('📋 First insurance:', data[0])
            } else {
              console.warn('⚠️  No insurance data found (empty array)')
            }
          } else {
            console.error('❌ Response is not an array:', data)
            this.insurances = []
          }
        } else {
          const errorText = await response.text()
          console.error('❌ Failed to load insurances:', response.status, response.statusText)
          console.error('📄 Error response:', errorText)
          this.insurances = []
        }
      } catch (error) {
        console.error('❌ Error loading insurances:', error)
        console.error('💡 Make sure the backend server is running on http://localhost:8080')
        this.insurances = []
      }
    },
    onChange() {
      this.$emit('insurance-selected', this.selectedCode)
    }
  }
}
</script>

