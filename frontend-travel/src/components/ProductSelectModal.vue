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
      class="product-select-modal" 
      :style="{ display: show ? 'block' : 'none' }"
    >
      <div class="modal-header">
        <h5 class="modal-title">
          <i class="fas fa-box me-2"></i>Select Product Codes (Multiple)
        </h5>
        <button type="button" class="btn-close" @click="close"></button>
      </div>
      
      <div class="modal-body">
        <!-- Search -->
        <div class="mb-3 search-container">
          <input 
            type="text" 
            class="form-control" 
            v-model="searchQuery"
            placeholder="Search products..."
          />
        </div>
        
        <!-- Products List -->
        <div class="products-list">
          <div v-if="filteredProducts.length === 0" class="text-center p-3 text-muted">
            <p v-if="products.length === 0">Loading products...</p>
            <p v-else>No products found matching "{{ searchQuery }}"</p>
          </div>
          <div 
            v-for="product in filteredProducts" 
            :key="product.id"
            class="product-item"
            :class="{ 'selected': isSelected(product.code) }"
            @click="toggleProduct(product.code)"
          >
            <span class="product-name">{{ product.name }}</span>
            <i v-if="isSelected(product.code)" class="fas fa-check text-success ms-auto"></i>
          </div>
        </div>
        
        <!-- Selected Product Codes Field - Fixed at bottom -->
        <div class="product-code-section">
          <label class="form-label">Selected Product Codes ({{ localSelectedCodes.length }}):</label>
          <div class="input-group product-code-input-group">
            <input 
              type="text" 
              class="form-control font-monospace product-code-input" 
              :value="productCodesJson"
              readonly
            />
            <button 
              class="btn btn-outline-secondary" 
              type="button"
              @click="copyToClipboard"
              title="Copy to clipboard"
            >
              <i class="fas fa-copy"></i>
            </button>
          </div>
          <small class="form-text text-muted">
            Selected: {{ localSelectedCodes.length }} product code(s)
          </small>
          
          <!-- Action Buttons -->
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="close">
              Cancel
            </button>
            <button type="button" class="btn btn-primary" @click="confirm" :disabled="localSelectedCodes.length === 0">
              Confirm
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getProducts } from '../services/api'

export default {
  name: 'ProductSelectModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    selectedCodes: {
      type: Array,
      default: () => []
    }
  },
  emits: ['close', 'confirm'],
  data() {
    return {
      searchQuery: '',
      products: [],
      localSelectedCodes: []
    }
  },
  computed: {
    filteredProducts() {
      if (!this.searchQuery) {
        return this.products
      }
      const query = this.searchQuery.toLowerCase()
      return this.products.filter(product => 
        product.code.toLowerCase().includes(query) ||
        (product.name && product.name.toLowerCase().includes(query))
      )
    },
    productCodesJson() {
      if (this.localSelectedCodes.length === 0) {
        return '[]'
      }
      return JSON.stringify(this.localSelectedCodes.sort())
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        // Filter out invalid codes
        const validCodes = (this.selectedCodes || []).filter(code => code && code !== '')
        this.localSelectedCodes = validCodes.length > 0 ? [...validCodes] : []
        this.searchQuery = ''
        this.loadProducts()
      } else {
        this.localSelectedCodes = []
      }
    },
    selectedCodes(newVal) {
      if (!this.show) {
        const validCodes = (newVal || []).filter(code => code && code !== '')
        this.localSelectedCodes = [...validCodes]
      }
    }
  },
  methods: {
    async loadProducts() {
      try {
        this.products = await getProducts()
      } catch (error) {
        console.error('Error loading products:', error)
        window.showCustomAlert('Failed to load products', 'error')
      }
    },
    isSelected(productCode) {
      return this.localSelectedCodes.includes(productCode)
    },
    toggleProduct(productCode) {
      const index = this.localSelectedCodes.indexOf(productCode)
      if (index > -1) {
        // Remove if already selected
        this.localSelectedCodes.splice(index, 1)
      } else {
        // Add if not selected
        this.localSelectedCodes.push(productCode)
      }
    },
    copyToClipboard() {
      if (this.localSelectedCodes.length > 0) {
        navigator.clipboard.writeText(this.productCodesJson).then(() => {
          window.showCustomAlert('Product codes copied to clipboard!', 'success')
        }).catch(() => {
          window.showCustomAlert('Failed to copy', 'error')
        })
      }
    },
    close() {
      this.$emit('close')
    },
    confirm() {
      if (this.localSelectedCodes.length > 0) {
        this.$emit('confirm', this.localSelectedCodes)
        this.close()
      }
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

.product-select-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  z-index: 1050;
  width: 600px;
  max-width: 90%;
  height: 700px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  padding: 1.5rem;
  overflow: hidden;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  position: relative;
}

.search-container {
  flex-shrink: 0;
}

.products-list {
  flex: 1;
  min-height: 150px;
  max-height: 400px;
  overflow-y: auto;
  overflow-x: hidden;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 0.5rem;
  flex-shrink: 1;
}

.product-item {
  padding: 0.75rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  transition: background-color 0.2s;
  border: 1px solid transparent;
}

.product-item:hover {
  background-color: #f8f9fa;
  border-color: #dee2e6;
}

.product-item.selected {
  background-color: #e7f3ff;
  font-weight: 600;
  border-color: var(--primary-orange);
}

.product-name {
  flex: 1;
  color: #333;
  font-size: 1em;
  font-weight: 500;
}

.ms-auto {
  margin-left: auto;
}

.product-code-section {
  flex-shrink: 0;
  margin-top: auto;
  padding-bottom: 0;
  min-height: fit-content;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
  padding: 1rem 0 0.5rem 0;
  border-top: 1px solid #dee2e6;
  flex-shrink: 0;
  width: 100%;
  box-sizing: border-box;
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

.modal-actions .btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.modal-actions .btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.product-code-input-group {
  display: flex;
  width: 100%;
  overflow: hidden;
}

.product-code-input {
  font-size: 0.9em !important;
  word-break: break-all;
  overflow-wrap: break-word;
  white-space: normal;
  min-width: 0;
  flex: 1;
}

.product-code-input-group .btn {
  flex-shrink: 0;
  white-space: nowrap;
}
</style>
