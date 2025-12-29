<template>
  <div>
    <div 
      class="custom-alert-backdrop" 
      id="custom-alert-backdrop"
      :style="{ display: showBackdrop ? 'block' : 'none' }"
      @click="close"
    ></div>
    
    <div 
      class="custom-alert-modal" 
      id="custom-alert-modal"
      :style="{ display: showModal ? 'block' : 'none' }"
    >
      <div class="p-4">
        <div class="d-flex align-items-center mb-3">
          <i :class="alertIcon" :style="{ fontSize: '13px' }"></i>
          <h5 class="mb-0" :class="titleClass">{{ title }}</h5>
        </div>
        <p class="mb-3">{{ message }}</p>
        <div class="text-end">
          <button type="button" class="btn btn-primary" @click="close">OK</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CustomAlertModal',
  data() {
    return {
      showModal: false,
      showBackdrop: false,
      title: '',
      message: '',
      type: 'success'
    }
  },
  computed: {
    alertIcon() {
      const icons = {
        'success': 'fas fa-check-circle text-success me-2',
        'error': 'fas fa-exclamation-circle text-danger me-2',
        'warning': 'fas fa-exclamation-triangle text-warning me-2',
        'info': 'fas fa-info-circle text-info me-2'
      }
      return icons[this.type] || icons['info']
    },
    titleClass() {
      const classes = {
        'success': 'text-success',
        'error': 'text-danger',
        'warning': 'text-warning',
        'info': 'text-info'
      }
      return classes[this.type] || classes['info']
    }
  },
  mounted() {
    window.showCustomAlert = this.show
    window.closeCustomAlert = this.close
  },
  methods: {
    show(message, type = 'success') {
      this.message = message
      this.type = type
      this.showBackdrop = true
      this.showModal = true
    },
    close() {
      this.showBackdrop = false
      this.showModal = false
    }
  }
}
</script>

