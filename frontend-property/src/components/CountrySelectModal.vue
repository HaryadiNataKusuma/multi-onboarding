<template>
  <div>
    <!-- Backdrop -->
    <div 
      class="modal-backdrop" 
      :style="{ display: show ? 'block' : 'none' }"
      @click="close"
    ></div>
    
    <!-- Modal -->
    <div 
      class="country-select-modal" 
      :style="{ display: show ? 'block' : 'none' }"
    >
      <div class="modal-header">
        <h5 class="modal-title">
          <i class="fas fa-globe me-2"></i>Select Countries
        </h5>
        <button type="button" class="btn-close" @click="close"></button>
      </div>
      
      <div class="modal-body">
        <!-- Search and Select All -->
        <div class="mb-2 search-container">
          <div class="d-flex gap-2">
            <input 
              type="text" 
              class="form-control" 
              v-model="searchQuery"
              placeholder="Search countries..."
            />
            <button 
              type="button" 
              class="btn btn-primary"
              @click="selectAllCountries"
            >
              <i class="fas fa-check-double me-1"></i> Select All
            </button>
          </div>
        </div>
        
        <!-- Countries List -->
        <div class="countries-list">
          <div v-if="filteredCountries.length === 0" class="text-center p-3 text-muted">
            <p v-if="countries.length === 0">Loading countries...</p>
            <p v-else>No countries found matching "{{ searchQuery }}"</p>
          </div>
          <div 
            v-for="country in filteredCountries" 
            :key="country.id"
            class="country-item"
            :class="{ 'selected': isSelected(country.id) }"
            @click="toggleCountry(country.id)"
          >
            <span class="country-name">{{ country.name }}</span>
            <i v-if="isSelected(country.id)" class="fas fa-check text-success ms-auto"></i>
          </div>
        </div>
        
        <!-- Selected Country IDs Field -->
        <div class="country-ids-section">
          <label class="form-label">Country IDs:</label>
          <input 
            type="text" 
            class="form-control font-monospace country-ids-input" 
            :value="countryIdsJson"
            readonly
          />
        </div>
      </div>
      
      <!-- Action Buttons and Regions Section - Fixed at bottom -->
      <div class="modal-footer-custom">
        <!-- Regions Section - Exclude countries from existing regions -->
        <div v-if="regions.length > 0" class="regions-exclude-section">
          <label class="form-label mb-2">
            <i class="fas fa-map-marked-alt me-1"></i>Exclude Countries from Existing Regions:
          </label>
          <div class="regions-buttons">
            <button
              v-for="region in regions"
              :key="region.id"
              type="button"
              class="btn btn-sm btn-outline-warning region-exclude-btn"
              @click="excludeRegionCountries(region)"
              :title="getRegionTooltip(region)"
            >
              <i class="fas fa-minus-circle me-1"></i>
              {{ region.name }}
              <span class="badge bg-warning text-dark ms-1">{{ region.country_ids?.length || 0 }}</span>
            </button>
          </div>
        </div>
        
        <!-- Action Buttons -->
        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" @click="close">
            Cancel
          </button>
          <button type="button" class="btn btn-primary" @click="confirm">
            Confirm
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getCountries, getRegions } from '../services/api'

export default {
  name: 'CountrySelectModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    countries: {
      type: Array,
      default: () => []
    },
    selectedIds: {
      type: Array,
      default: () => []
    }
  },
  emits: ['close', 'confirm'],
  data() {
    return {
      searchQuery: '',
      localSelectedIds: [],
      regions: []
    }
  },
  computed: {
    filteredCountries() {
      if (!this.searchQuery) {
        return this.countries
      }
      const query = this.searchQuery.toLowerCase()
      return this.countries.filter(country => 
        country.name.toLowerCase().includes(query)
      )
    },
    selectedCountries() {
      return this.countries.filter(c => this.localSelectedIds.includes(c.id))
    },
    countryIdsJson() {
      // Filter out invalid IDs before displaying
      const validIds = this.localSelectedIds.filter(id => id && id !== 0 && !isNaN(id))
      if (validIds.length === 0) {
        return '[]'
      }
      return JSON.stringify(validIds.sort((a, b) => a - b))
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        // Filter out invalid IDs (0, null, undefined, NaN) from selectedIds
        const validIds = this.selectedIds.filter(id => id && id !== 0 && !isNaN(id) && typeof id === 'number')
        this.localSelectedIds = validIds.length > 0 ? [...validIds] : []
        this.searchQuery = ''
        // Load regions when modal opens
        this.loadRegions()
      } else {
        // Clear when modal closes
        this.localSelectedIds = []
        this.regions = []
      }
    },
    selectedIds(newVal) {
      if (!this.show) {
        // Filter out invalid IDs
        const validIds = newVal.filter(id => id && id !== 0)
        this.localSelectedIds = [...validIds]
      }
    }
  },
  methods: {
    isSelected(countryId) {
      // Filter out invalid IDs (0, null, undefined, NaN)
      if (!countryId || countryId === 0 || isNaN(countryId)) {
        return false
      }
      // Filter localSelectedIds to only valid IDs before checking
      const validIds = this.localSelectedIds.filter(id => id && id !== 0 && !isNaN(id))
      return validIds.includes(countryId)
    },
    toggleCountry(countryId) {
      // Filter out invalid IDs
      if (!countryId || countryId === 0 || isNaN(countryId)) {
        return
      }
      // Filter out invalid IDs from localSelectedIds first
      this.localSelectedIds = this.localSelectedIds.filter(id => id && id !== 0 && !isNaN(id))
      
      const index = this.localSelectedIds.indexOf(countryId)
      if (index > -1) {
        // Remove if already selected
        this.localSelectedIds.splice(index, 1)
      } else {
        // Add if not selected
        this.localSelectedIds.push(countryId)
      }
    },
    removeCountry(countryId) {
      const index = this.localSelectedIds.indexOf(countryId)
      if (index > -1) {
        this.localSelectedIds.splice(index, 1)
      }
    },
    selectAllCountries() {
      // Get all valid country IDs from filtered countries
      const validCountryIds = this.filteredCountries
        .map(c => c.id)
        .filter(id => id && id !== 0 && !isNaN(id))
      
      // Add all valid IDs that are not already selected
      validCountryIds.forEach(id => {
        if (!this.localSelectedIds.includes(id)) {
          this.localSelectedIds.push(id)
        }
      })
      
      // Filter out any invalid IDs
      this.localSelectedIds = this.localSelectedIds.filter(id => id && id !== 0 && !isNaN(id))
    },
    async loadRegions() {
      try {
        const regions = await getRegions()
        // Filter out invalid country_ids and ensure they're arrays
        this.regions = (regions || []).map(region => ({
          ...region,
          country_ids: Array.isArray(region.country_ids) 
            ? region.country_ids.filter(id => id && id !== 0 && !isNaN(id))
            : []
        })).filter(region => region.country_ids.length > 0) // Only show regions with countries
      } catch (error) {
        console.error('Error loading regions:', error)
        this.regions = []
      }
    },
    getRegionCountriesText(region) {
      if (!region || !region.country_ids || region.country_ids.length === 0) {
        return ''
      }
      
      // Get country names from country IDs
      const validCountryIds = region.country_ids.filter(id => id && id !== 0 && !isNaN(id))
      const countryNames = validCountryIds
        .map(id => {
          const country = this.countries.find(c => c.id === id)
          return country ? country.name : null
        })
        .filter(name => name !== null)
      
      if (countryNames.length === 0) {
        return `${validCountryIds.length} countries`
      }
      
      // Show first 5 countries, then "and X more" if there are more
      if (countryNames.length <= 5) {
        return countryNames.join(', ')
      } else {
        return countryNames.slice(0, 5).join(', ') + ` and ${countryNames.length - 5} more`
      }
    },
    getRegionTooltip(region) {
      if (!region || !region.country_ids || region.country_ids.length === 0) {
        return `Exclude countries from ${region.name}`
      }
      
      // Get country names from country IDs
      const validCountryIds = region.country_ids.filter(id => id && id !== 0 && !isNaN(id))
      const countryNames = validCountryIds
        .map(id => {
          const country = this.countries.find(c => c.id === id)
          return country ? country.name : `Country ID: ${id}`
        })
        .filter(name => name)
      
      const countriesText = countryNames.length > 0 
        ? countryNames.join(', ')
        : `${validCountryIds.length} countries`
      
      return `Exclude ${validCountryIds.length} countries from ${region.name}:\n${countriesText}`
    },
    excludeRegionCountries(region) {
      if (!region || !region.country_ids || region.country_ids.length === 0) {
        return
      }
      
      // Get valid country IDs from region
      const regionCountryIds = region.country_ids.filter(id => id && id !== 0 && !isNaN(id))
      
      // Get country names for display
      const countryNames = regionCountryIds
        .map(id => {
          const country = this.countries.find(c => c.id === id)
          return country ? country.name : null
        })
        .filter(name => name !== null)
      
      // Remove all country IDs that exist in this region from localSelectedIds
      this.localSelectedIds = this.localSelectedIds.filter(id => !regionCountryIds.includes(id))
      
      // Filter out any invalid IDs
      this.localSelectedIds = this.localSelectedIds.filter(id => id && id !== 0 && !isNaN(id))
      
      // Show feedback with country names
      const countriesList = countryNames.length > 0 
        ? countryNames.slice(0, 10).join(', ') + (countryNames.length > 10 ? ` and ${countryNames.length - 10} more` : '')
        : `${regionCountryIds.length} countries`
      
      window.showCustomAlert(
        `Excluded ${regionCountryIds.length} countries from "${region.name}" region:\n${countriesList}`,
        'success'
      )
    },
    close() {
      this.$emit('close')
    },
    confirm() {
      this.$emit('confirm', this.localSelectedIds)
      this.close()
    }
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1040;
}

.country-select-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  z-index: 1050;
  width: 700px;
  max-width: 95vw;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* Prevent content from overflowing */
}

.modal-header {
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.modal-title {
  margin: 0;
  font-weight: 600;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  opacity: 0.5;
}

.btn-close:hover {
  opacity: 1;
}

.modal-body {
  padding: 0.75rem 1.25rem;
  overflow: hidden;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  position: relative;
}

.search-container {
  flex-shrink: 0;
}

.countries-list {
  flex: 1;
  min-height: 120px;
  max-height: calc(85vh - 420px);
  overflow-y: auto;
  overflow-x: hidden;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 0.4rem;
  flex-shrink: 1;
}

.country-item {
  padding: 0.75rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background-color 0.2s;
  border: 1px solid transparent;
}

.country-item:hover {
  background-color: #f8f9fa;
  border-color: #dee2e6;
}

.country-item.selected {
  background-color: #e7f3ff;
  font-weight: 600;
  border-color: var(--primary-orange);
}

.country-name {
  flex: 1;
}

.ms-auto {
  margin-left: auto;
}

.country-ids-section {
  flex-shrink: 0;
  padding-top: 0.4rem;
  padding-bottom: 0;
  min-height: fit-content;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.country-ids-section .form-label {
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}

.modal-footer-custom {
  padding: 0.6rem 1.25rem;
  border-top: 1px solid #dee2e6;
  background-color: #fff;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}


.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  flex-shrink: 0;
  width: 100%;
}

.modal-actions .btn {
  min-width: 100px;
  padding: 0.625rem 1.5rem;
  font-weight: 500;
  font-size: 0.95rem;
  line-height: 1.5;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.modal-actions .btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.selected-countries-container {
  flex-shrink: 0;
  max-height: 80px;
  overflow-y: auto;
  overflow-x: hidden;
  order: -1; /* Move above country-ids-section */
}

.selected-countries {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.country-ids-input {
  font-size: 0.9em !important;
  word-break: break-all;
  overflow-wrap: break-word;
  white-space: normal;
  width: 100%;
}

.regions-exclude-section {
  flex: 1;
  min-width: 0;
  flex-shrink: 0;
}

.regions-exclude-section .form-label {
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.35rem;
  display: block;
  font-size: 0.9rem;
}

.regions-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  max-height: 80px;
  overflow-y: auto;
  padding: 0.4rem;
  background-color: #f8f9fa;
  border-radius: 6px;
  width: 100%;
}

.region-exclude-btn {
  white-space: nowrap;
  font-size: 0.85rem;
  padding: 0.4rem 0.75rem;
  border-color: #ffc107;
  color: #856404;
  transition: all 0.2s ease;
}

.region-exclude-btn:hover {
  background-color: #ffc107;
  color: #000;
  border-color: #ffc107;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(255, 193, 7, 0.3);
}

.region-exclude-btn .badge {
  font-size: 0.7rem;
  padding: 0.2rem 0.4rem;
}
</style>

