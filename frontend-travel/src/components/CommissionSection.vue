<template>
  <div class="col-12">
    <!-- Form Section -->
    <div class="form-section">
      <div class="d-flex justify-content-between align-items-center">
        <h3><i class="fas fa-chart-bar me-1"></i> Tambah Data Commission Baru</h3>
      </div>
      
      <form @submit.prevent="saveCommission">
        <!-- Product Code Selection -->
        <div class="form-section mb-4" style="border: 1px solid #dee2e6; border-radius: 0.375rem; padding: 1rem;">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <label class="form-label mb-0"><strong>Pilih Product Code <span class="text-danger">*</span></strong></label>
            <div>
              <button type="button" class="btn btn-sm btn-outline-primary me-2" @click="selectAllProducts">
                <i class="fas fa-check-square"></i> Select All
              </button>
              <button type="button" class="btn btn-sm btn-outline-secondary me-2" @click="deselectAllProducts">
                <i class="fas fa-square"></i> Deselect All
              </button>
              <button type="button" class="btn btn-sm btn-outline-info" @click="loadProducts" title="Refresh product list">
                <i class="fas fa-sync-alt"></i> Refresh
              </button>
            </div>
          </div>
          <div v-if="filteredProductsForSelection.length === 0" class="text-muted">
            <i class="fas fa-info-circle"></i> No products available. Please select an insurance code first.
          </div>
          <div v-else class="product-code-selection" style="max-height: 200px; overflow-y: auto; border: 1px solid #dee2e6; border-radius: 0.375rem; padding: 0.75rem;" @click="refreshProductsOnClick">
            <div class="form-check mb-2" v-for="product in filteredProductsForSelection" :key="product.code">
              <input 
                class="form-check-input" 
                type="checkbox" 
                :id="'product-' + product.code"
                :value="product.code"
                v-model="selectedProductCodes"
              >
              <label class="form-check-label" :for="'product-' + product.code">
                <strong>{{ product.code }}</strong> - {{ product.name }}
              </label>
            </div>
          </div>
          <div v-if="selectedProductCodes.length > 0" class="mt-2 text-muted small">
            <i class="fas fa-check-circle text-success"></i> {{ selectedProductCodes.length }} product(s) selected
          </div>
        </div>
        
        <div class="form-grid form-grid-4">
          <div>
            <label class="form-label">Commission Percentage <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="form.commission_percentage" 
              placeholder="0.00" 
              required
            >
          </div>
          
          <div>
            <label class="form-label">Commission VAT Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="form.commission_vat_type" required>
              <option value="">Select type...</option>
              <option value="EXCLUSIVE">EXCLUSIVE</option>
              <option value="INCLUSIVE">INCLUSIVE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">AF Percentage <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="form.af_percentage" 
              placeholder="0.00" 
              required
            >
          </div>
          
          <div>
            <label class="form-label">AF VAT Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="form.af_vat_type" required>
              <option value="">Select type...</option>
              <option value="EXCLUSIVE">EXCLUSIVE</option>
              <option value="INCLUSIVE">INCLUSIVE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Admin Fee <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="form.admin_fee" 
              placeholder="0.00" 
              required
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
      <h3><i class="fas fa-edit me-1"></i> Edit Data Commission</h3>
      
      <form @submit.prevent="updateCommission">
        <input type="hidden" v-model="editForm.id">
        <div class="form-grid form-grid-4">
          <div style="display: none;">
            <label class="form-label">Product Code</label>
            <select class="form-select" v-model="editForm.product_code" @focus="loadProducts">
              <option value="">Auto-generate from products</option>
              <option v-for="product in products" :key="product.code" :value="product.code">
                {{ product.code }} - {{ product.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Commission Percentage <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="editForm.commission_percentage" 
              required
            >
          </div>
          
          <div>
            <label class="form-label">Commission VAT Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.commission_vat_type" required>
              <option value="">Select type...</option>
              <option value="EXCLUSIVE">EXCLUSIVE</option>
              <option value="INCLUSIVE">INCLUSIVE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">AF Percentage <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="editForm.af_percentage" 
              required
            >
          </div>
          
          <div>
            <label class="form-label">AF VAT Type <span class="text-danger">*</span></label>
            <select class="form-select" v-model="editForm.af_vat_type" required>
              <option value="">Select type...</option>
              <option value="EXCLUSIVE">EXCLUSIVE</option>
              <option value="INCLUSIVE">INCLUSIVE</option>
            </select>
          </div>
          
          <div>
            <label class="form-label">Admin Fee <span class="text-danger">*</span></label>
            <input 
              type="number" 
              step="0.01"
              class="form-control" 
              v-model="editForm.admin_fee" 
              required
            >
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
    <div v-if="showDraftSection && !showEditForm" class="table-section draft-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📊 Data Preview: Commissions (Draft Save)</h4>
        <div>
          <button 
            class="btn btn-success btn-confirm me-2" 
            :disabled="!canConfirm" 
            @click="confirmCommissions"
          >
            <i class="fas fa-database me-2"></i> Konfirmasi Data (DB)
          </button>
          <button type="button" class="btn btn-secondary btn-sm" @click="hideDraft">
            <i class="fas fa-times"></i> Close
          </button>
        </div>
      </div>

      <!-- Tabs for each table -->
      <ul class="nav nav-tabs mt-3" role="tablist">
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: draftActiveTab === 'commissions' }"
            @click="draftActiveTab = 'commissions'"
            type="button"
          >
            commissions.commissions
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: draftActiveTab === 'plan_commissions' }"
            @click="draftActiveTab = 'plan_commissions'"
            type="button"
          >
            plan_commissions
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: draftActiveTab === 'default_config_products' }"
            @click="draftActiveTab = 'default_config_products'"
            type="button"
          >
            default_config_products
          </button>
        </li>
      </ul>

      <div class="table-content">
        <!-- Tab: commissions.commissions -->
        <div v-show="draftActiveTab === 'commissions'" class="tab-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Product Code</th>
                  <th>Insurance Code</th>
                  <th>Agent Level</th>
                  <th>Corporate ID</th>
                  <th>Basic Commission</th>
                  <th>Note</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="draftCommissions.length === 0">
                  <td colspan="10" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No draft data available for commissions.commissions
                  </td>
                </tr>
                <tr v-for="item in sortedDraftCommissions" :key="item.id">
                  <td>{{ item.id }}</td>
                  <td>{{ item.product_code }}</td>
                  <td>{{ item.insurance_code || '-' }}</td>
                  <td><span class="badge bg-secondary">{{ item.agent_level }}</span></td>
                  <td>{{ item.corporate_id || '-' }}</td>
                  <td>{{ item.basic_commission || '-' }}</td>
                  <td>{{ item.note || '-' }}</td>
                  <td>{{ formatDateOnly(item.start_date) }}</td>
                  <td>{{ formatDateOnly(item.end_date) }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEdit(item)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteDraft(item.id, 'commissions')">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Tab: plan_commissions -->
        <div v-show="draftActiveTab === 'plan_commissions'" class="tab-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Plan Code</th>
                  <th>Product ID</th>
                  <th>Insurer ID</th>
                  <th>Commission %</th>
                  <th>Comm VAT Type</th>
                  <th>AF %</th>
                  <th>AF VAT Type</th>
                  <th>Admin Fee</th>
                  <th>Hardcopy Fee</th>
                  <th>Is Active</th>
                  <th>Version</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="draftPlanCommissions.length === 0">
                  <td colspan="13" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No draft data available for plan_commissions
                  </td>
                </tr>
                <tr v-for="item in draftPlanCommissions" :key="item.id">
                  <td>{{ item.id }}</td>
                  <td>{{ item.plan_code || '-' }}</td>
                  <td>{{ item.product_id || '-' }}</td>
                  <td>{{ item.insurer_id || '-' }}</td>
                  <td>{{ item.commission_percentage || '-' }}</td>
                  <td>
                    <span :class="getVATTypeBadge(item.commission_vat_type).class">
                      {{ item.commission_vat_type || '-' }}
                    </span>
                  </td>
                  <td>{{ item.af_percentage || '-' }}</td>
                  <td>
                    <span :class="getVATTypeBadge(item.af_vat_type).class">
                      {{ item.af_vat_type || '-' }}
                    </span>
                  </td>
                  <td>{{ item.admin_fee || '-' }}</td>
                  <td>{{ item.hardcopy_fee || '-' }}</td>
                  <td>
                    <span :class="item.is_active ? 'badge bg-success' : 'badge bg-secondary'">
                      {{ item.is_active ? 'Yes' : 'No' }}
                    </span>
                  </td>
                  <td>{{ item.version || '-' }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEdit(item)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteDraft(item.id, 'plan_commissions')">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Tab: default_config_products -->
        <div v-show="draftActiveTab === 'default_config_products'" class="tab-content">
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead>
                <tr>
                  <th>Level</th>
                  <th>Sequence</th>
                  <th>Kind</th>
                  <th>Category</th>
                  <th>Insurer</th>
                  <th>Product</th>
                  <th>Selected</th>
                  <th>Renewal Count</th>
                  <th>Commission</th>
                  <th>Max Discount</th>
                  <th>Order</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="draftDefaultConfig.length === 0">
                  <td colspan="12" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No draft data available for default_config_products
                  </td>
                </tr>
                <tr v-for="item in sortedDraftDefaultConfig" :key="item.id">
                  <td><span class="badge bg-info">{{ item.level }}</span></td>
                  <td>{{ item.sequence || '-' }}</td>
                  <td>{{ item.kind || '-' }}</td>
                  <td>{{ item.category || '-' }}</td>
                  <td>{{ item.insurer || '-' }}</td>
                  <td>{{ item.product || '-' }}</td>
                  <td>
                    <span :class="item.selected ? 'badge bg-success' : 'badge bg-secondary'">
                      {{ item.selected ? 'Yes' : 'No' }}
                    </span>
                  </td>
                  <td>{{ item.renewal_count || 0 }}</td>
                  <td>{{ item.commission || '-' }}</td>
                  <td>{{ item.max_discount || '-' }}</td>
                  <td>{{ item.order || '-' }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEdit(item)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteDraft(item.id, 'default_config')">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Result Table Section -->
    <div v-if="showResultSection && !showEditForm" class="table-section result-section show">
      <div class="table-header d-flex justify-content-between align-items-center">
        <h4>📈 Data Result: Commissions (From Database)</h4>
        <div>
          <button class="btn btn-info btn-sm me-2" @click="refreshResult">
            <i class="fas fa-sync-alt"></i> Refresh Data
          </button>
          <button type="button" class="btn btn-secondary btn-sm" @click="hideResult">
            <i class="fas fa-times"></i> Close
          </button>
        </div>
      </div>

      <!-- Tabs for each table -->
      <ul class="nav nav-tabs mt-3" role="tablist">
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: resultActiveTab === 'commissions' }"
            @click="resultActiveTab = 'commissions'; loadResultCommissions()"
            type="button"
          >
            commissions.commissions
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: resultActiveTab === 'plan_commissions' }"
            @click="resultActiveTab = 'plan_commissions'; loadResultPlanCommissions()"
            type="button"
          >
            plan_commissions
          </button>
        </li>
        <li class="nav-item" role="presentation">
          <button 
            class="nav-link" 
            :class="{ active: resultActiveTab === 'default_config_products' }"
            @click="resultActiveTab = 'default_config_products'; loadResultDefaultConfig()"
            type="button"
          >
            default_config_products
          </button>
        </li>
      </ul>

      <div class="table-content">
        <!-- Tab: commissions.commissions -->
        <div v-show="resultActiveTab === 'commissions'" class="tab-content">
          <!-- Filter Section -->
          <div class="p-3 border-bottom" style="background-color: #f8f9fa;">
            <div class="row g-3">
              <div class="col-md-4">
                <label class="form-label">Filter by Product Code</label>
                <input 
                  type="text" 
                  class="form-control form-control-sm" 
                  v-model="filterCommissions.product_code" 
                  placeholder="Enter product code..."
                >
              </div>
              <div class="col-md-4">
                <label class="form-label">Filter by Agent Level</label>
                <input 
                  type="text" 
                  class="form-control form-control-sm" 
                  v-model="filterCommissions.agent_level" 
                  placeholder="Enter agent level..."
                >
              </div>
              <div class="col-md-4 d-flex align-items-end">
                <button class="btn btn-warning btn-sm" @click="clearFilterCommissions">
                  <i class="fas fa-times"></i> Clear Filter
                </button>
              </div>
            </div>
          </div>
          
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead class="table-success">
                <tr>
                  <th>ID</th>
                  <th>Product Code</th>
                  <th>Insurance Code</th>
                  <th>Agent Level</th>
                  <th>Corporate ID</th>
                  <th>Basic Commission</th>
                  <th>Note</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th>Created At</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="filteredResultCommissions.length === 0">
                  <td colspan="11" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No data found in commissions.commissions
                  </td>
                </tr>
                <tr v-for="item in filteredResultCommissions" :key="item.id">
                  <td>{{ item.id }}</td>
                  <td>{{ item.product_code || '-' }}</td>
                  <td>{{ item.insurance_code || '-' }}</td>
                  <td><span class="badge bg-secondary">{{ item.agent_level || '-' }}</span></td>
                  <td>{{ item.corporate_id || '-' }}</td>
                  <td>{{ item.basic_commission || '-' }}</td>
                  <td>{{ item.note || '-' }}</td>
                  <td>{{ formatDateOnly(item.start_date) }}</td>
                  <td>{{ formatDateOnly(item.end_date) }}</td>
                  <td>{{ formatDateOnly(item.created_at) }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEditFromResult(item, 'commissions')">
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Tab: plan_commissions -->
        <div v-show="resultActiveTab === 'plan_commissions'" class="tab-content">
          <!-- Filter Section -->
          <div class="p-3 border-bottom" style="background-color: #f8f9fa;">
            <div class="row g-3">
              <div class="col-md-8">
                <label class="form-label">Filter by Plan Code</label>
                <input 
                  type="text" 
                  class="form-control form-control-sm" 
                  v-model="filterPlanCommissions.plan_code" 
                  placeholder="Enter plan code..."
                >
              </div>
              <div class="col-md-4 d-flex align-items-end">
                <button class="btn btn-warning btn-sm" @click="clearFilterPlanCommissions">
                  <i class="fas fa-times"></i> Clear Filter
                </button>
              </div>
            </div>
          </div>
          
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead class="table-success">
                <tr>
                  <th>ID</th>
                  <th>Plan Code</th>
                  <th>Product ID</th>
                  <th>Insurer ID</th>
                  <th>Commission %</th>
                  <th>Comm VAT Type</th>
                  <th>AF %</th>
                  <th>AF VAT Type</th>
                  <th>Admin Fee</th>
                  <th>Hardcopy Fee</th>
                  <th>Is Active</th>
                  <th>Version</th>
                  <th>Created At</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="filteredResultPlanCommissions.length === 0">
                  <td colspan="14" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No data found in plan_commissions
                  </td>
                </tr>
                <tr v-for="item in filteredResultPlanCommissions" :key="item.id">
                  <td>{{ item.id }}</td>
                  <td>{{ item.plan_code || '-' }}</td>
                  <td>{{ item.product_id || '-' }}</td>
                  <td>{{ item.insurer_id || '-' }}</td>
                  <td>{{ item.commission_percentage || '-' }}</td>
                  <td>
                    <span :class="getVATTypeBadge(item.commission_vat_type).class">
                      {{ item.commission_vat_type || '-' }}
                    </span>
                  </td>
                  <td>{{ item.af_percentage || '-' }}</td>
                  <td>
                    <span :class="getVATTypeBadge(item.af_vat_type).class">
                      {{ item.af_vat_type || '-' }}
                    </span>
                  </td>
                  <td>{{ item.admin_fee || '-' }}</td>
                  <td>{{ item.hardcopy_fee || '-' }}</td>
                  <td>
                    <span :class="item.is_active ? 'badge bg-success' : 'badge bg-secondary'">
                      {{ item.is_active ? 'Yes' : 'No' }}
                    </span>
                  </td>
                  <td>{{ item.version || '-' }}</td>
                  <td>{{ formatDateOnly(item.created_at) }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEditFromResult(item, 'plan_commissions')">
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Tab: default_config_products -->
        <div v-show="resultActiveTab === 'default_config_products'" class="tab-content">
          <!-- Filter Section -->
          <div class="p-3 border-bottom" style="background-color: #f8f9fa;">
            <div class="row g-3">
              <div class="col-md-4">
                <label class="form-label">Filter by Product</label>
                <input 
                  type="text" 
                  class="form-control form-control-sm" 
                  v-model="filterDefaultConfig.product" 
                  placeholder="Enter product..."
                >
              </div>
              <div class="col-md-4">
                <label class="form-label">Filter by Level</label>
                <input 
                  type="text" 
                  class="form-control form-control-sm" 
                  v-model="filterDefaultConfig.level" 
                  placeholder="Enter level..."
                >
              </div>
              <div class="col-md-4 d-flex align-items-end">
                <button class="btn btn-warning btn-sm" @click="clearFilterDefaultConfig">
                  <i class="fas fa-times"></i> Clear Filter
                </button>
              </div>
            </div>
          </div>
          
          <div class="table-responsive">
            <table class="table table-striped table-hover align-middle">
              <thead class="table-success">
                <tr>
                  <th>Level</th>
                  <th>Sequence</th>
                  <th>Kind</th>
                  <th>Category</th>
                  <th>Insurer</th>
                  <th>Product</th>
                  <th>Selected</th>
                  <th>Renewal Count</th>
                  <th>Commission</th>
                  <th>Max Discount</th>
                  <th>Order</th>
                  <th>Created At</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="filteredResultDefaultConfig.length === 0">
                  <td colspan="13" class="text-center text-muted py-4">
                    <i class="fas fa-info-circle me-2"></i>No data found in default_config_products
                  </td>
                </tr>
                <tr v-for="item in filteredResultDefaultConfig" :key="item.id">
                  <td><span class="badge bg-info">{{ item.level || '-' }}</span></td>
                  <td>{{ item.sequence || '-' }}</td>
                  <td>{{ item.kind || '-' }}</td>
                  <td>{{ item.category || '-' }}</td>
                  <td>{{ item.insurer || '-' }}</td>
                  <td>{{ item.product || '-' }}</td>
                  <td>
                    <span :class="item.selected ? 'badge bg-success' : 'badge bg-secondary'">
                      {{ item.selected ? 'Yes' : 'No' }}
                    </span>
                  </td>
                  <td>{{ item.renewal_count || 0 }}</td>
                  <td>{{ item.commission || '-' }}</td>
                  <td>{{ item.max_discount || '-' }}</td>
                  <td>{{ item.order || '-' }}</td>
                  <td>{{ formatDateOnly(item.created_at) }}</td>
                  <td>
                    <button class="btn btn-sm btn-warning me-1" @click="startEditFromResult(item, 'default_config')">
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
  </div>
</template>

<script>
export default {
  name: 'CommissionSection',
  props: {
    selectedInsuranceCode: String
  },
  data() {
    return {
      form: {
        product_code: '',
        commission_percentage: 0,
        commission_vat_type: '',
        af_percentage: 0,
        af_vat_type: '',
        admin_fee: 0
      },
      editForm: {
        id: null,
        product_code: '',
        commission_percentage: 0,
        commission_vat_type: '',
        af_percentage: 0,
        af_vat_type: '',
        admin_fee: 0,
        tableType: 'commissions',
        plan_code: '',
        product_id: null,
        insurer_id: null,
        hardcopy_fee: 0,
        version: 1,
        is_active: 1
      },
      showEditForm: false,
      showDraftSection: false,
      showResultSection: false,
      draftData: [],
      draftCommissions: [],
      draftPlanCommissions: [],
      draftDefaultConfig: [],
      resultData: [],
      resultCommissions: [],
      resultPlanCommissions: [],
      resultDefaultConfig: [],
      products: [],
      selectedProductCodes: [], // Array of selected product codes
      canConfirm: false,
      editingFromDraft: false,
      draftIdCounter: 1,
      draftActiveTab: 'commissions',
      resultActiveTab: 'commissions',
      API_URL: 'http://localhost:8080/api/travel',
      COMMISSIONS_API_URL: 'http://localhost:8080/api/commissions', // Shared legacy route for commission tables
      // Filter states
      filterCommissions: {
        product_code: '',
        agent_level: ''
      },
      filterDefaultConfig: {
        product: '',
        level: ''
      },
      filterPlanCommissions: {
        plan_code: ''
      }
    }
  },
  computed: {
    sortedDraftCommissions() {
      return [...this.draftCommissions].sort((a, b) => {
        const productCompare = (a.product_code || '').localeCompare(b.product_code || '')
        if (productCompare !== 0) {
          return productCompare
        }
        return (a.agent_level || '').localeCompare(b.agent_level || '')
      })
    },
    sortedDraftDefaultConfig() {
      return [...this.draftDefaultConfig].sort((a, b) => {
        const productCompare = (a.product || '').localeCompare(b.product || '')
        if (productCompare !== 0) {
          return productCompare
        }
        return (a.level || '').localeCompare(b.level || '')
      })
    },
    sortedResultCommissions() {
      return [...this.resultCommissions].sort((a, b) => {
        const productCompare = (a.product_code || '').localeCompare(b.product_code || '')
        if (productCompare !== 0) {
          return productCompare
        }
        return (a.agent_level || '').localeCompare(b.agent_level || '')
      })
    },
    sortedResultDefaultConfig() {
      return [...this.resultDefaultConfig].sort((a, b) => {
        const productCompare = (a.product || '').localeCompare(b.product || '')
        if (productCompare !== 0) {
          return productCompare
        }
        return (a.level || '').localeCompare(b.level || '')
      })
    },
    // Filtered data
    filteredResultCommissions() {
      let filtered = this.sortedResultCommissions
      
      if (this.filterCommissions.product_code) {
        filtered = filtered.filter(item => 
          (item.product_code || '').toLowerCase().includes(this.filterCommissions.product_code.toLowerCase())
        )
      }
      
      if (this.filterCommissions.agent_level) {
        filtered = filtered.filter(item => 
          (item.agent_level || '').toLowerCase().includes(this.filterCommissions.agent_level.toLowerCase())
        )
      }
      
      return filtered
    },
    filteredResultDefaultConfig() {
      let filtered = this.sortedResultDefaultConfig
      
      if (this.filterDefaultConfig.product) {
        filtered = filtered.filter(item => 
          (item.product || '').toLowerCase().includes(this.filterDefaultConfig.product.toLowerCase())
        )
      }
      
      if (this.filterDefaultConfig.level) {
        filtered = filtered.filter(item => 
          (item.level || '').toLowerCase().includes(this.filterDefaultConfig.level.toLowerCase())
        )
      }
      
      return filtered
    },
    filteredResultPlanCommissions() {
      let filtered = [...this.resultPlanCommissions]
      
      if (this.filterPlanCommissions.plan_code) {
        filtered = filtered.filter(item => 
          (item.plan_code || '').toLowerCase().includes(this.filterPlanCommissions.plan_code.toLowerCase())
        )
      }
      
      return filtered
    },
    // Unique values for filter dropdowns
    uniqueProductCodes() {
      return [...new Set(this.resultCommissions.map(item => item.product_code).filter(Boolean))].sort()
    },
    uniqueAgentLevels() {
      return [...new Set(this.resultCommissions.map(item => item.agent_level).filter(Boolean))].sort()
    },
    uniqueProducts() {
      return [...new Set(this.resultDefaultConfig.map(item => item.product).filter(Boolean))].sort()
    },
    uniqueLevels() {
      return [...new Set(this.resultDefaultConfig.map(item => item.level).filter(Boolean))].sort()
    },
    uniquePlanCodes() {
      return [...new Set(this.resultPlanCommissions.map(item => item.plan_code).filter(Boolean))].sort()
    },
    // Filtered products for selection (based on selected insurance code)
    filteredProductsForSelection() {
      if (!this.selectedInsuranceCode) {
        return []
      }
      return this.products.filter(p => p.insurance_code === this.selectedInsuranceCode)
    }
  },
  async mounted() {
    await this.loadProducts()
  },
  watch: {
    selectedInsuranceCode(newCode, oldCode) {
      // Reload products when insurance code changes
      if (newCode !== oldCode) {
        this.loadProducts()
      }
      // Reload result data when insurance code changes
      if (this.showResultSection && newCode !== oldCode) {
        this.loadResult()
      }
      // Clear selected product codes when insurance changes
      if (newCode !== oldCode) {
        this.selectedProductCodes = []
      }
    }
  },
  methods: {
    async loadProducts() {
      try {
        // Use travel domain endpoint to get products filtered by insurance_code
        let url = `${this.API_URL}/products`
        // Filter by selectedInsuranceCode if provided
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const response = await fetch(url)
        if (response.ok) {
          let data = await response.json()
          // Handle both array and wrapped response formats
          if (Array.isArray(data)) {
            this.products = data
          } else if (data.data && Array.isArray(data.data)) {
            this.products = data.data
          } else {
            this.products = []
          }
          // Additional client-side filter if needed
          if (this.selectedInsuranceCode) {
            this.products = this.products.filter(item => item.insurance_code === this.selectedInsuranceCode)
          }
          console.log('Products loaded for commission section:', this.products.length, 'products for insurance:', this.selectedInsuranceCode)
        } else {
          console.error('Failed to load products:', response.status, response.statusText)
          this.products = []
        }
      } catch (error) {
        console.error('Error loading products:', error)
        this.products = []
      }
    },
    
    async saveCommission() {
      // Validate required fields
      if (!this.form.commission_percentage && this.form.commission_percentage !== 0) {
        window.showCustomAlert('Commission Percentage is required', 'error')
        return
      }
      if (!this.form.commission_vat_type) {
        window.showCustomAlert('Commission VAT Type is required', 'error')
        return
      }
      if (!this.form.af_percentage && this.form.af_percentage !== 0) {
        window.showCustomAlert('AF Percentage is required', 'error')
        return
      }
      if (!this.form.af_vat_type) {
        window.showCustomAlert('AF VAT Type is required', 'error')
        return
      }
      if (!this.form.admin_fee && this.form.admin_fee !== 0) {
        window.showCustomAlert('Admin Fee is required', 'error')
        return
      }
      
      // Check if insurance code is selected
      if (!this.selectedInsuranceCode) {
        window.showCustomAlert('Please select an insurance code first', 'error')
        return
      }
      
      // Check if at least one product code is selected
      if (!this.selectedProductCodes || this.selectedProductCodes.length === 0) {
        window.showCustomAlert('Please select at least one product code', 'error')
        return
      }
      
      try {
        // Get all products for the selected insurance
        const productsResponse = await fetch(`${this.API_URL}/products`)
        if (!productsResponse.ok) {
          throw new Error('Failed to load products')
        }
        const allProducts = await productsResponse.json()
        
        // Filter products by insurance code and selected product codes
        const filteredProducts = allProducts.filter(p => 
          p.insurance_code === this.selectedInsuranceCode && 
          this.selectedProductCodes.includes(p.code)
        )
        
        if (filteredProducts.length === 0) {
          window.showCustomAlert('No products found for the selected insurance and product codes', 'error')
          return
        }
        
        // Agent levels to generate
        const agentLevels = [
          'GREEN', 'SILVER', 'GOLD', 'DIAMOND', 'PLATINUM', 'CORPORATE', 'CORPORATE2',
          'CORPORATE3', 'CORPORATE4', 'TIED', 'SKB', 'SHOPDRIVE', 'MVPARTNERSHIP1', 'DIRECTPROPERTY'
        ]
        
        // Generate draft items for each product and each agent level
        // Split into three tables based on structure
        for (const product of filteredProducts) {
          // Table 2: quotation_service_development.plan_commissions (only once per product, no tiering)
          const planCommissionId = this.draftIdCounter++
          const planTimestamp = new Date().toISOString()
          const planCommissionItem = {
            id: planCommissionId,
            plan_code: product.code,
            product_id: 5,
            insurer_id: null,
            commission_percentage: parseFloat(this.form.commission_percentage) || 0,
            commission_vat_type: this.form.commission_vat_type,
            af_percentage: parseFloat(this.form.af_percentage) || 0,
            af_vat_type: this.form.af_vat_type,
            admin_fee: parseFloat(this.form.admin_fee) || 0,
            hardcopy_fee: 0,
            is_active: 1,
            version: 1,
            timestamp: planTimestamp
          }
          this.draftPlanCommissions.push(planCommissionItem)
          
          // Generate for each agent level (for commissions and default_config_products)
          for (const agentLevel of agentLevels) {
            const baseId = this.draftIdCounter++
            const timestamp = new Date().toISOString()
            
            // Table 1: commission_service_development.commissions
            // Set corporate_id: 1 for CORPORATE, CORPORATE2, CORPORATE3, CORPORATE4, otherwise 0
            const corporateLevels = ['CORPORATE', 'CORPORATE2', 'CORPORATE3', 'CORPORATE4']
            const corporateId = corporateLevels.includes(agentLevel) ? 1 : 0
            
            const commissionItem = {
              id: baseId,
              product_code: product.code,
              insurance_code: this.selectedInsuranceCode,
              agent_level: agentLevel,
              corporate_id: corporateId,
              basic_commission: parseFloat(this.form.commission_percentage) || 0,
              bonus_point: 0,
              note: `Commission for ${product.code} - ${agentLevel}`,
              start_date: timestamp,
              end_date: timestamp,
              timestamp: timestamp
            }
            this.draftCommissions.push(commissionItem)
            
            // Table 3: agent_service_development.default_config_products
            const commissionValue = parseFloat(this.form.commission_percentage) || 0
            const defaultConfigItem = {
              id: baseId,
              level: agentLevel,
              sequence: 9,
              kind: 'PL',
              category: 'TV',
              insurer: this.selectedInsuranceCode,
              product: product.code,
              payment_frequency: null,
              policy_year: null,
              selected: 1,
              renewal_count: 0,
              commission: commissionValue,
              point: 0,
              upline_bonus: 0,
              bonus: 0,
              max_discount: commissionValue,
              order: 1000000,
              timestamp: timestamp
            }
            this.draftDefaultConfig.push(defaultConfigItem)
            
            // Keep old draftData for backward compatibility
            const draftItem = {
              id: baseId,
              product_code: product.code,
              product_name: product.name,
              agent_level: agentLevel,
              insurance_code: this.selectedInsuranceCode,
              commission_percentage: parseFloat(this.form.commission_percentage) || 0,
              commission_vat_type: this.form.commission_vat_type,
              af_percentage: parseFloat(this.form.af_percentage) || 0,
              af_vat_type: this.form.af_vat_type,
              admin_fee: parseFloat(this.form.admin_fee) || 0,
              timestamp: timestamp
            }
            this.draftData.push(draftItem)
          }
        }
        
        this.canConfirm = this.draftData.length > 0
        const totalRows = filteredProducts.length * agentLevels.length
        window.showCustomAlert(`Commission draft saved successfully! Generated ${totalRows} rows (${filteredProducts.length} products × ${agentLevels.length} agent levels)`, 'success')
        this.resetForm()
        this.showDraftSection = true
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to save commission: ' + (error.message || 'Unknown error'), 'error')
      }
    },
    
    async showDraft() {
      await this.loadDraft()
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data available', 'info')
        this.showDraftSection = false
        return
      }
      this.showDraftSection = true
      this.$nextTick(() => {
        setTimeout(() => {
          const draftSection = document.querySelector('.draft-section')
          if (draftSection) {
            draftSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
          }
        }, 100)
      })
    },
    
    async loadDraft() {
      try {
        const response = await fetch(`${this.API_URL}/commissions/draft`)
        if (response.ok) {
          const data = await response.json()
          this.draftData = data || []
          this.canConfirm = this.draftData.length > 0
        }
      } catch (error) {
        console.error('Error loading draft:', error)
      }
    },
    
    async clearDraft() {
      const confirmed = await window.showCustomConfirm('Are you sure you want to clear all draft data?')
      if (!confirmed) return
      
      try {
        // Clear in-memory draft
        this.draftData = []
        this.draftCommissions = []
        this.draftPlanCommissions = []
        this.draftDefaultConfig = []
        this.canConfirm = false
        this.showDraftSection = false
        this.draftIdCounter = 1
        window.showCustomAlert('Draft cleared successfully!', 'success')
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Failed to clear draft', 'error')
      }
    },
    
    hideDraft() {
      this.showDraftSection = false
    },
    
    async confirmCommissions() {
      if (this.draftData.length === 0) {
        window.showCustomAlert('No draft data to confirm', 'info')
        return
      }
      
      // Check if insurance code is selected
      if (!this.selectedInsuranceCode) {
        window.showCustomAlert('Please select an insurance code first', 'error')
        return
      }
      
      const confirmed = await window.showCustomConfirm('Are you sure you want to confirm all draft data? This will generate commission data to 3 tables (commission_service_development.commissions, agent_service_development.default_config_products, quotation_service_development.plan_commissions)')
      if (!confirmed) return
      
      try {
        // Check backend connection first
        try {
          const healthCheck = await fetch(`${this.API_URL}/products`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
          })
          if (!healthCheck.ok && healthCheck.status === 404) {
            throw new Error('Backend server is not responding. Please ensure the backend server is running on http://localhost:8080')
          }
        } catch (healthError) {
          if (healthError.message.includes('Failed to fetch') || healthError.message.includes('NetworkError')) {
            throw new Error('Cannot connect to backend server. Please ensure the backend server is running on http://localhost:8080')
          }
          throw healthError
        }
        
        // Get unique product codes and commission data from draft
        const uniqueProducts = [...new Set(this.draftData.map(item => item.product_code).filter(code => code))]
        const firstDraftItem = this.draftData[0]
        
        if (uniqueProducts.length === 0) {
          window.showCustomAlert('No product codes found in draft data', 'error')
          return
        }
        
        // Check for duplicates before generating
        const checkDuplicateResponse = await fetch(`${this.COMMISSIONS_API_URL}/check-duplicates`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ product_codes: uniqueProducts })
        })
        
        if (!checkDuplicateResponse.ok) {
          const errorData = await checkDuplicateResponse.json().catch(() => ({ error: 'Failed to check duplicates' }))
          window.showCustomAlert(errorData.error || 'Failed to check for duplicates', 'error')
          return
        }
        
        const duplicateData = await checkDuplicateResponse.json()
        
        // Check if there are any duplicates
        // duplicateData is an object with table names as keys and arrays of duplicate codes as values
        // If empty object {}, it means no duplicates
        const duplicateTables = Object.keys(duplicateData).filter(table => {
          const codes = duplicateData[table]
          return Array.isArray(codes) && codes.length > 0
        })
        
        if (duplicateTables.length > 0) {
          // Build error message with duplicate information
          let errorMessage = 'The following product codes already exist in the database:\n\n'
          for (const table of duplicateTables) {
            const codes = duplicateData[table]
            if (Array.isArray(codes) && codes.length > 0) {
              errorMessage += `${table}:\n  - ${codes.join('\n  - ')}\n\n`
            }
          }
          errorMessage += 'Please remove these product codes from your selection or update the existing records.'
          window.showCustomAlert(errorMessage, 'error')
          return
        }
        
        // No duplicates found, continue with generation
        console.log('No duplicates found, proceeding with commission generation')
        
        // Call generate endpoint
        const generateData = {
          insurance_code: this.selectedInsuranceCode,
          product_codes: uniqueProducts,
          commission_percentage: firstDraftItem.commission_percentage,
          commission_vat_type: firstDraftItem.commission_vat_type,
          af_percentage: firstDraftItem.af_percentage,
          af_vat_type: firstDraftItem.af_vat_type,
          admin_fee: firstDraftItem.admin_fee
        }
        
        console.log('Sending generate request:', generateData)
        
        const response = await fetch(`${this.COMMISSIONS_API_URL}/generate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(generateData)
        })
        
        console.log('Response status:', response.status, response.statusText)
        
        if (!response.ok) {
          let errorMessage = 'Failed to generate commissions'
          if (response.status === 404) {
            errorMessage = `Endpoint not found (404). Please check:\n1. Backend server is running on http://localhost:8080\n2. Endpoint /api/commissions/generate is registered\n3. Check backend logs for errors`
          } else if (response.status === 500) {
            errorMessage = 'Server error. Please check backend logs for details.'
          } else {
            try {
              const errorData = await response.json()
              errorMessage = errorData.error || errorData.message || errorMessage
            } catch (e) {
              const errorText = await response.text().catch(() => '')
              errorMessage = errorText || errorMessage
            }
          }
          throw new Error(errorMessage)
        }
        
        const result = await response.json().catch(() => ({ message: 'Success' }))
        console.log('Generate result:', result)
        
        // Agent levels count (should match backend)
        const agentLevelsCount = 14 // Updated: 14 agent levels (GREEN, SILVER, GOLD, DIAMOND, PLATINUM, CORPORATE, CORPORATE2, CORPORATE3, CORPORATE4, TIED, SKB, SHOPDRIVE, MVPARTNERSHIP1, DIRECTPROPERTY)
        
        // Calculate total rows: 
        // - commissions: agentLevelsCount rows per product
        // - default_config_products: agentLevelsCount rows per product  
        // - plan_commissions: 1 row per product
        // Total: (agentLevelsCount × 2) + 1 = 29 rows per product
        const rowsPerProduct = (agentLevelsCount * 2) + 1
        const totalRows = uniqueProducts.length * rowsPerProduct
        window.showCustomAlert(`Data confirmed successfully! Generated ${totalRows} rows to 3 tables (${uniqueProducts.length} products × ${agentLevelsCount} agent levels).`, 'success')
        this.draftData = []
        this.draftCommissions = []
        this.draftPlanCommissions = []
        this.draftDefaultConfig = []
        this.canConfirm = false
        this.showDraftSection = false
        this.draftIdCounter = 1
        await this.loadResult()
      } catch (error) {
        console.error('Error in confirmCommissions:', error)
        window.showCustomAlert(error.message || 'Failed to confirm commissions', 'error')
      }
    },
    
    async showResult() {
      try {
        await this.loadResult()
        this.showResultSection = true
        const totalCount = this.resultCommissions.length + this.resultPlanCommissions.length + this.resultDefaultConfig.length
        if (totalCount === 0) {
          window.showCustomAlert('No commissions found in database', 'info')
        }
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
        window.showCustomAlert('Failed to load commissions: ' + (error.message || 'Unknown error'), 'error')
        // Still show the section even if there's an error, so user can see the error message
        this.showResultSection = true
      }
    },
    
    async loadResult() {
      try {
        // Load all three tables
        await Promise.all([
          this.loadResultCommissions(),
          this.loadResultPlanCommissions(),
          this.loadResultDefaultConfig()
        ])
      } catch (error) {
        console.error('Error loading result:', error)
      }
    },
    
    async loadResultCommissions() {
      try {
        // Use shared legacy route /api/commissions/table/commissions
        let url = `${this.COMMISSIONS_API_URL}/table/commissions`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const response = await fetch(url)
        if (response.ok) {
          const data = await response.json()
          // Handle both array and wrapped response formats
          if (Array.isArray(data)) {
            this.resultCommissions = data
          } else if (data.data && Array.isArray(data.data)) {
            this.resultCommissions = data.data
          } else {
            this.resultCommissions = []
          }
          console.log('Commissions loaded:', this.resultCommissions.length, 'for insurance:', this.selectedInsuranceCode)
        } else {
          console.error('Failed to load commissions:', response.status, response.statusText)
          this.resultCommissions = []
        }
      } catch (error) {
        console.error('Error loading commissions table:', error)
        this.resultCommissions = []
      }
    },
    
    async loadResultPlanCommissions() {
      try {
        // Use shared legacy route /api/commissions/table/plan_commissions
        let url = `${this.COMMISSIONS_API_URL}/table/plan_commissions`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const response = await fetch(url)
        if (response.ok) {
          const data = await response.json()
          // Handle both array and wrapped response formats
          if (Array.isArray(data)) {
            this.resultPlanCommissions = data
          } else if (data.data && Array.isArray(data.data)) {
            this.resultPlanCommissions = data.data
          } else {
            this.resultPlanCommissions = []
          }
        } else {
          console.error('Failed to load plan_commissions:', response.status, response.statusText)
          this.resultPlanCommissions = []
        }
      } catch (error) {
        console.error('Error loading plan_commissions table:', error)
        this.resultPlanCommissions = []
      }
    },
    
    async loadResultDefaultConfig() {
      try {
        // Use shared legacy route /api/commissions/table/default_config_products
        let url = `${this.COMMISSIONS_API_URL}/table/default_config_products`
        if (this.selectedInsuranceCode) {
          url += `?insurance_code=${encodeURIComponent(this.selectedInsuranceCode)}`
        }
        const response = await fetch(url)
        if (response.ok) {
          const data = await response.json()
          // Handle both array and wrapped response formats
          if (Array.isArray(data)) {
            this.resultDefaultConfig = data
          } else if (data.data && Array.isArray(data.data)) {
            this.resultDefaultConfig = data.data
          } else {
            this.resultDefaultConfig = []
          }
        } else {
          console.error('Failed to load default_config_products:', response.status, response.statusText)
          this.resultDefaultConfig = []
        }
      } catch (error) {
        console.error('Error loading default_config_products table:', error)
        this.resultDefaultConfig = []
      }
    },
    
    async refreshResult() {
      await this.loadResult()
    },
    
    hideResult() {
      this.showResultSection = false
    },
    
    clearFilterCommissions() {
      this.filterCommissions = {
        product_code: '',
        agent_level: ''
      }
    },
    
    clearFilterDefaultConfig() {
      this.filterDefaultConfig = {
        product: '',
        level: ''
      }
    },
    
    clearFilterPlanCommissions() {
      this.filterPlanCommissions = {
        plan_code: ''
      }
    },
    
    async deleteDraft(id, tableType) {
      const confirmed = await window.showCustomConfirm('Are you sure you want to delete this draft item?')
      if (!confirmed) return
      
      try {
        // Remove from all draft arrays
        this.draftData = this.draftData.filter(d => d.id !== id)
        if (tableType === 'commissions') {
          this.draftCommissions = this.draftCommissions.filter(d => d.id !== id)
        } else if (tableType === 'plan_commissions') {
          this.draftPlanCommissions = this.draftPlanCommissions.filter(d => d.id !== id)
        } else if (tableType === 'default_config') {
          this.draftDefaultConfig = this.draftDefaultConfig.filter(d => d.id !== id)
        } else {
          // Remove from all if no specific table
          this.draftCommissions = this.draftCommissions.filter(d => d.id !== id)
          this.draftPlanCommissions = this.draftPlanCommissions.filter(d => d.id !== id)
          this.draftDefaultConfig = this.draftDefaultConfig.filter(d => d.id !== id)
        }
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
      this.editForm = {
        id: item.id,
        product_code: item.product_code,
        commission_percentage: item.commission_percentage || 0,
        commission_vat_type: item.commission_vat_type || '',
        af_percentage: item.af_percentage || 0,
        af_vat_type: item.af_vat_type || '',
        admin_fee: item.admin_fee || 0
      }
      this.editingFromDraft = true
      this.showEditForm = true
      this.showDraftSection = false
    },
    
    startEditFromResult(item, tableType) {
      this.editForm = {
        id: item.id,
        product_code: item.product_code || item.product || '',
        commission_percentage: item.commission_percentage || item.commission || item.basic_commission || 0,
        commission_vat_type: item.commission_vat_type || '',
        af_percentage: item.af_percentage || item.point || item.bonus_point || 0,
        af_vat_type: item.af_vat_type || '',
        admin_fee: item.admin_fee || 0,
        tableType: tableType || 'commissions',
        // Additional fields for plan_commissions versioning
        plan_code: item.plan_code || '',
        product_id: item.product_id || null,
        insurer_id: item.insurer_id || null,
        hardcopy_fee: item.hardcopy_fee || 0,
        version: item.version || 1,
        is_active: item.is_active !== undefined ? item.is_active : 1
      }
      this.editingFromDraft = false
      this.showEditForm = true
      this.showResultSection = false
    },
    
    async updateCommission() {
      try {
        if (this.editingFromDraft) {
          // Update in draft (in-memory)
          const index = this.draftData.findIndex(d => d.id === this.editForm.id)
          if (index !== -1) {
            this.draftData[index] = {
              ...this.draftData[index],
              ...this.editForm,
              commission_percentage: parseFloat(this.editForm.commission_percentage) || 0,
              af_percentage: parseFloat(this.editForm.af_percentage) || 0,
              admin_fee: parseFloat(this.editForm.admin_fee) || 0
            }
            window.showCustomAlert('Commission draft updated successfully!', 'success')
            this.cancelEdit()
            this.showDraftSection = true
            this.canConfirm = this.draftData.length > 0
          }
        } else {
          // Update in database
          const tableType = this.editForm.tableType || 'commissions'
          
          // Special handling for plan_commissions with versioning
          if (tableType === 'plan_commissions') {
            // Step 1: Deactivate existing record (set is_active = 0)
            const deactivateData = {
              is_active: 0,
              updated_by: 1
            }
            
            const deactivateResponse = await fetch(`${this.API_URL}/commissions/plan_commissions/${this.editForm.id}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(deactivateData)
            })
            
            if (!deactivateResponse.ok) {
              const errorData = await deactivateResponse.json().catch(() => ({ error: 'Failed to deactivate existing record' }))
              window.showCustomAlert(errorData.error || 'Failed to deactivate existing record', 'error')
              return
            }
            
            // Step 2: Create new record with version + 1
            const newVersion = (this.editForm.version || 1) + 1
            const newPlanCommissionData = {
              plan_code: this.editForm.plan_code || this.editForm.product_code,
              product_id: this.editForm.product_id || 5,
              insurer_id: this.editForm.insurer_id || null,
              commission_percentage: parseFloat(this.editForm.commission_percentage) || 0,
              commission_vat_type: this.editForm.commission_vat_type,
              af_percentage: parseFloat(this.editForm.af_percentage) || 0,
              af_vat_type: this.editForm.af_vat_type,
              admin_fee: parseFloat(this.editForm.admin_fee) || 0,
              hardcopy_fee: parseFloat(this.editForm.hardcopy_fee) || 0,
              is_active: 1,
              version: newVersion,
              created_by: 1
            }
            
            const createResponse = await fetch(`${this.API_URL}/commissions/plan_commissions`, {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(newPlanCommissionData)
            })
            
            if (createResponse.ok) {
              window.showCustomAlert(`Plan commission updated successfully! New version ${newVersion} created.`, 'success')
              this.cancelEdit()
              await this.loadResultPlanCommissions()
              this.showResultSection = true
              this.resultActiveTab = 'plan_commissions'
            } else {
              const errorData = await createResponse.json().catch(() => ({ error: 'Failed to create new version' }))
              window.showCustomAlert(errorData.error || 'Failed to create new version', 'error')
            }
          } else {
            // Regular update for other tables
            const commissionData = {
              product_code: this.editForm.product_code,
              commission_percentage: parseFloat(this.editForm.commission_percentage) || 0,
              commission_vat_type: this.editForm.commission_vat_type,
              af_percentage: parseFloat(this.editForm.af_percentage) || 0,
              af_vat_type: this.editForm.af_vat_type,
              admin_fee: parseFloat(this.editForm.admin_fee) || 0,
              updated_by: 1
            }
            
            const response = await fetch(`${this.API_URL}/commissions/${this.editForm.id}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(commissionData)
            })
            
            if (response.ok) {
              window.showCustomAlert('Commission updated successfully!', 'success')
              this.cancelEdit()
              await this.loadResult()
              this.showResultSection = true
            } else {
              const errorData = await response.json().catch(() => ({ error: 'Failed to update' }))
              window.showCustomAlert(errorData.error || 'Failed to update', 'error')
            }
          }
        }
      } catch (error) {
        console.error('Error:', error)
        window.showCustomAlert('Connection error', 'error')
      }
    },
    
    cancelEdit() {
      this.showEditForm = false
      this.editForm = {
        id: null,
        product_code: '',
        commission_percentage: 0,
        commission_vat_type: '',
        af_percentage: 0,
        af_vat_type: '',
        admin_fee: 0,
        tableType: 'commissions',
        plan_code: '',
        product_id: null,
        insurer_id: null,
        hardcopy_fee: 0,
        version: 1,
        is_active: 1
      }
      this.editingFromDraft = false
    },
    
    resetForm() {
      this.form = {
        product_code: '',
        commission_percentage: 0,
        commission_vat_type: '',
        af_percentage: 0,
        af_vat_type: '',
        admin_fee: 0
      }
      // Don't clear selectedProductCodes - user might want to keep selection
    },
    
    selectAllProducts() {
      this.selectedProductCodes = this.filteredProductsForSelection.map(p => p.code)
    },
    
    deselectAllProducts() {
      this.selectedProductCodes = []
    },
    
    async refreshProductsOnClick(event) {
      // Only refresh if clicking on the container itself, not on checkboxes or labels
      if (event.target === event.currentTarget || event.target.classList.contains('product-code-selection')) {
        await this.loadProducts()
      }
    },
    
    getVATTypeBadge(type) {
      return {
        text: type,
        class: `badge ${type === 'EXCLUSIVE' ? 'bg-warning' : 'bg-success'}`
      }
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
