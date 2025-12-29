<template>
  <div class="product-form">
    <div class="form-section">
      <h3>{{ editingProduct ? 'Edit Product' : 'Add New Product' }}</h3>
      <form @submit.prevent="handleSubmit">
        <div class="form-grid form-grid-2">
          <div class="mb-3">
            <label for="insurance_code" class="form-label">Insurance Code</label>
            <input
              type="text"
              id="insurance_code"
              v-model="form.insurance_code"
              class="form-control"
              placeholder="Enter insurance code"
              required
            />
          </div>

          <div class="mb-3">
            <label for="insurance_name" class="form-label">Insurance Name</label>
            <input
              type="text"
              id="insurance_name"
              v-model="form.insurance_name"
              class="form-control"
              placeholder="Enter insurance name"
              required
            />
          </div>

          <div class="mb-3">
            <label for="name" class="form-label">Product Name</label>
            <input
              type="text"
              id="name"
              v-model="form.name"
              class="form-control"
              placeholder="Enter product name"
              required
            />
          </div>

          <div class="mb-3">
            <label for="destination" class="form-label">Destination</label>
            <input
              type="text"
              id="destination"
              v-model="form.destination"
              class="form-control"
              placeholder="Enter destination"
              required
            />
          </div>

          <div class="mb-3">
            <label for="category" class="form-label">Category</label>
            <select
              id="category"
              v-model="form.category"
              class="form-select"
              required
            >
              <option value="">Select category</option>
              <option value="adventure">Adventure</option>
              <option value="beach">Beach</option>
              <option value="cultural">Cultural</option>
              <option value="family">Family</option>
              <option value="luxury">Luxury</option>
            </select>
          </div>

          <div class="mb-3">
            <label for="price" class="form-label">Price (USD)</label>
            <input
              type="number"
              id="price"
              v-model.number="form.price"
              class="form-control"
              placeholder="0.00"
              min="0"
              step="0.01"
              required
            />
          </div>

          <div class="mb-3">
            <label for="duration" class="form-label">Duration (days)</label>
            <input
              type="number"
              id="duration"
              v-model.number="form.duration"
              class="form-control"
              placeholder="Enter duration in days"
              min="1"
              required
            />
          </div>

          <div class="mb-3">
            <label for="status" class="form-label">Status</label>
            <select
              id="status"
              v-model="form.status"
              class="form-select"
              required
            >
              <option value="draft">Draft</option>
              <option value="published">Published</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </div>

        <div class="mb-3">
          <label for="logo" class="form-label">Logo URL</label>
          <input
            type="url"
            id="logo"
            v-model="form.logo"
            class="form-control"
            placeholder="https://example.com/logo.png"
          />
          <small class="form-text text-muted">Enter URL to the logo image</small>
          <div v-if="form.logo" class="mt-2">
            <img :src="form.logo" alt="Logo preview" class="logo-preview" @error="handleImageError" />
          </div>
        </div>

        <div class="mb-3">
          <label for="description" class="form-label">Description</label>
          <textarea
            id="description"
            v-model="form.description"
            class="form-control"
            rows="4"
            placeholder="Enter product description"
            required
          ></textarea>
        </div>

        <div class="d-flex gap-2">
          <button type="submit" class="btn btn-primary">
            {{ editingProduct ? 'Update Product' : 'Create Product' }}
          </button>
          <button type="button" class="btn btn-secondary" @click="handleCancel">
            Cancel
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import { createProduct, updateProduct } from '../services/api'

export default {
  name: 'ProductForm',
  props: {
    editingProduct: {
      type: Object,
      default: null
    }
  },
  emits: ['product-saved', 'cancel'],
  setup(props, { emit }) {
    const form = ref({
      insurance_code: '',
      insurance_name: '',
      name: '',
      description: '',
      destination: '',
      price: 0,
      duration: 1,
      category: '',
      status: 'draft',
      logo: ''
    })

    // Watch for editing product changes
    watch(() => props.editingProduct, (newProduct) => {
      if (newProduct) {
        form.value = {
          insurance_code: newProduct.insurance_code || '',
          insurance_name: newProduct.insurance_name || '',
          name: newProduct.name || '',
          description: newProduct.description || '',
          destination: newProduct.destination || '',
          price: newProduct.price || 0,
          duration: newProduct.duration || 1,
          category: newProduct.category || '',
          status: newProduct.status || 'draft',
          logo: newProduct.logo || ''
        }
      } else {
        resetForm()
      }
    }, { immediate: true })

    const resetForm = () => {
      form.value = {
        insurance_code: '',
        insurance_name: '',
        name: '',
        description: '',
        destination: '',
        price: 0,
        duration: 1,
        category: '',
        status: 'draft',
        logo: ''
      }
    }

    const handleImageError = (event) => {
      event.target.style.display = 'none'
    }

    const handleSubmit = async () => {
      try {
        if (props.editingProduct) {
          await updateProduct(props.editingProduct.id, form.value)
        } else {
          await createProduct(form.value)
        }
        emit('product-saved')
        resetForm()
      } catch (error) {
        console.error('Error saving product:', error)
        alert('Failed to save product. Please try again.')
      }
    }

    const handleCancel = () => {
      resetForm()
      emit('cancel')
    }

    return {
      form,
      handleSubmit,
      handleCancel,
      handleImageError
    }
  }
}
</script>

<style scoped>
.product-form {
  margin-bottom: 2rem;
}

.form-grid {
  display: grid;
  gap: 1rem;
}

.form-grid-2 {
  grid-template-columns: repeat(2, 1fr);
}

@media (max-width: 768px) {
  .form-grid-2 {
    grid-template-columns: 1fr;
  }
}

.gap-2 {
  gap: 0.5rem;
}

.d-flex {
  display: flex;
}

.logo-preview {
  max-width: 150px;
  max-height: 150px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 8px;
  background: #f8f9fa;
  object-fit: contain;
}

.form-text {
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.text-muted {
  color: #6c757d;
}

.mt-2 {
  margin-top: 0.5rem;
}
</style>

