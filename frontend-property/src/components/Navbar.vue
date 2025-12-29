<template>
  <nav class="navbar navbar-expand-lg navbar-light bg-light border-bottom border-warning border-5">
    <div class="container-fluid">
      <h2 class="h4 mb-0">Product Onboarding Dashboard</h2>
      <div class="d-flex align-items-center">
        <span class="badge bg-success me-2" :class="dbStatusClass">
          {{ dbStatusText }}
        </span>
        <i class="fas fa-user-circle fa-2x text-secondary"></i>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      dbStatus: 'connected' // 'connected', 'disconnected'
    }
  },
  computed: {
    dbStatusClass() {
      return this.dbStatus === 'connected' ? 'bg-success' : 'bg-danger'
    },
    dbStatusText() {
      return this.dbStatus === 'connected' ? '✅ MySQL Connected' : '❌ MySQL Disconnected'
    }
  },
  mounted() {
    this.checkDbStatus()
  },
  methods: {
    async checkDbStatus() {
      try {
        // You can add API call here to check DB status
        this.dbStatus = 'connected'
      } catch (error) {
        this.dbStatus = 'disconnected'
      }
    }
  }
}
</script>

