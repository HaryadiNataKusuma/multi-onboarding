<template>
  <div id="wrapper">
    <Sidebar />
    <div id="page-content-wrapper">
      <Navbar />
      <InsuranceDropdown ref="insuranceDropdown" @insurance-selected="handleInsuranceSelected" />
      <div class="main-content-area">
        <div class="container-fluid">
          <div class="row">
            <InsuranceSection :selected-insurance-code="selectedInsuranceCode" @insurance-confirmed="handleInsuranceConfirmed" />
            <RegionsSection />
            <ProductSection :selected-insurance-code="selectedInsuranceCode" />
            <ProductRulesSection :selected-insurance-code="selectedInsuranceCode" />
            <TemplatesSection :selected-insurance-code="selectedInsuranceCode" />
            <CommissionSection :selected-insurance-code="selectedInsuranceCode" />
            <AddonsSection :selected-insurance-code="selectedInsuranceCode" />
            <HistoriesSection />
            <AllResultSection />
          </div>
        </div>
      </div>
    </div>
    <CustomAlertModal />
    <CustomConfirmModal />
  </div>
</template>

<script>
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
import InsuranceDropdown from './components/InsuranceDropdown.vue'
import InsuranceSection from './components/InsuranceSection.vue'
import ProductSection from './components/ProductSection.vue'
import RegionsSection from './components/RegionsSection.vue'
import TemplatesSection from './components/TemplatesSection.vue'
import ProductRulesSection from './components/ProductRulesSection.vue'
import CommissionSection from './components/CommissionSection.vue'
import AddonsSection from './components/AddonsSection.vue'
import HistoriesSection from './components/HistoriesSection.vue'
import AllResultSection from './components/AllResultSection.vue'
import CustomAlertModal from './components/CustomAlertModal.vue'
import CustomConfirmModal from './components/CustomConfirmModal.vue'

export default {
  name: 'App',
  components: {
    Sidebar,
    Navbar,
    InsuranceDropdown,
    InsuranceSection,
    ProductSection,
    RegionsSection,
    TemplatesSection,
    ProductRulesSection,
    CommissionSection,
    AddonsSection,
    HistoriesSection,
    AllResultSection,
    CustomAlertModal,
    CustomConfirmModal
  },
  data() {
    return {
      selectedInsuranceCode: ''
    }
  },
  methods: {
    handleInsuranceSelected(code) {
      this.selectedInsuranceCode = code
    },
    handleInsuranceConfirmed() {
      // Refresh insurance dropdown when new insurance is confirmed
      if (this.$refs.insuranceDropdown) {
        this.$refs.insuranceDropdown.loadInsurances()
      }
    }
  }
}
</script>


