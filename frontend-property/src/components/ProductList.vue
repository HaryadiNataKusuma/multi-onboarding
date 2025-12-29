<template>
  <div class="product-list">
    <div class="table-section">
      <div class="table-header">
        <h4>Travel Products</h4>
      </div>
      <div v-if="loading" class="text-center p-4">
        <p>Loading products...</p>
      </div>
      <div v-else-if="products.length === 0" class="text-center p-4">
        <p>No products found. Create your first product!</p>
      </div>
      <div v-else class="table-responsive">
        <table class="table table-hover">
          <thead>
            <tr>
              <th>Logo</th>
              <th>Insurance Code</th>
              <th>Insurance Name</th>
              <th>Product Name</th>
              <th>Destination</th>
              <th>Category</th>
              <th>Price</th>
              <th>Duration</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in products" :key="product.id">
              <td>
                <img 
                  v-if="product.logo" 
                  :src="product.logo" 
                  :alt="product.insurance_name || 'Logo'"
                  class="table-logo"
                  @error="handleImageError"
                />
                <span v-else class="text-muted">No logo</span>
              </td>
              <td>
                <strong>{{ product.insurance_code || '-' }}</strong>
              </td>
              <td>{{ product.insurance_name || '-' }}</td>
              <td>
                <strong>{{ product.name }}</strong>
                <br>
                <small class="text-muted">{{ product.description }}</small>
              </td>
              <td>{{ product.destination }}</td>
              <td>
                <span class="badge bg-secondary">{{ product.category }}</span>
              </td>
              <td>${{ product.price.toFixed(2) }}</td>
              <td>{{ product.duration }} days</td>
              <td>
                <span 
                  :class="{
                    'badge bg-success': product.status === 'published',
                    'badge bg-warning': product.status === 'draft',
                    'badge bg-secondary': product.status === 'archived'
                  }"
                >
                  {{ product.status }}
                </span>
              </td>
              <td>
                <button
                  class="btn btn-sm btn-outline-primary me-1"
                  @click="handleEdit(product)"
                >
                  Edit
                </button>
                <button
                  class="btn btn-sm btn-outline-danger"
                  @click="handleDelete(product.id)"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { getProducts, deleteProduct } from '../services/api'

export default {
  name: 'ProductList',
  props: {
    products: {
      type: Array,
      default: () => []
    }
  },
  emits: ['refresh', 'edit'],
  setup(props, { emit }) {
    const loading = ref(false)

    const handleEdit = (product) => {
      emit('edit', product)
    }

    const handleDelete = async (id) => {
      if (!confirm('Are you sure you want to delete this product?')) {
        return
      }

      try {
        await deleteProduct(id)
        emit('refresh')
      } catch (error) {
        console.error('Error deleting product:', error)
        alert('Failed to delete product. Please try again.')
      }
    }

    const handleImageError = (event) => {
      event.target.style.display = 'none'
    }

    return {
      loading,
      handleEdit,
      handleDelete,
      handleImageError
    }
  }
}
</script>

<style scoped>
.product-list {
  margin-top: 2rem;
}

.text-center {
  text-align: center;
}

.p-4 {
  padding: 1.5rem;
}

.text-muted {
  color: #6c757d;
}

.me-1 {
  margin-right: 0.25rem;
}

.table-responsive {
  overflow-x: auto;
}

.table-logo {
  max-height: 50px;
  max-width: 100px;
  object-fit: contain;
  vertical-align: middle;
  border: 1px solid #dee2e6;
  padding: 4px;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.05);
}
</style>

