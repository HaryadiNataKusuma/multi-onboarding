<template>
  <div>
    <div 
      class="custom-confirm-backdrop" 
      id="custom-confirm-backdrop"
      :style="{ display: showBackdrop ? 'block' : 'none' }"
      @click="close(false)"
    ></div>
    
    <div 
      class="custom-confirm-modal" 
      id="custom-confirm-modal"
      :style="{ display: showModal ? 'block' : 'none' }"
    >
      <div class="p-4">
        <div class="d-flex align-items-center mb-3">
          <i class="fas fa-question-circle text-warning me-2" style="font-size: 13px;"></i>
          <h5 class="mb-0 text-warning">Konfirmasi</h5>
        </div>
        <p class="mb-3">{{ message }}</p>
        <div class="text-end">
          <button type="button" class="btn btn-secondary me-2" @click="close(false)">Batal</button>
          <button type="button" class="btn btn-danger" @click="close(true)">Ya, Lanjutkan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CustomConfirmModal',
  data() {
    return {
      showModal: false,
      showBackdrop: false,
      message: '',
      callback: null
    }
  },
  mounted() {
    window.showCustomConfirm = this.show
    window.closeCustomConfirm = this.close
    window.customConfirmCallback = null
  },
  methods: {
    show(message, callback) {
      this.message = message
      this.showBackdrop = true
      this.showModal = true
      
      if (typeof callback === 'function') {
        window.customConfirmCallback = callback
      }
      
      return new Promise(resolve => {
        window.customConfirmCallback = (result) => resolve(result)
      })
    },
    close(result) {
      this.showBackdrop = false
      this.showModal = false
      
      if (window.customConfirmCallback) {
        window.customConfirmCallback(result)
        window.customConfirmCallback = null
      }
    }
  }
}
</script>

