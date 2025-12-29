import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/travel', // Using travel program routes
  headers: {
    'Content-Type': 'application/json'
  }
})

export const getProducts = async () => {
  const response = await api.get('/products')
  return response.data
}

export const getProduct = async (id) => {
  const response = await api.get(`/products/${id}`)
  return response.data
}

export const createProduct = async (product) => {
  const response = await api.post('/products', product)
  return response.data
}

export const updateProduct = async (id, product) => {
  const response = await api.put(`/products/${id}`, product)
  return response.data
}

export const deleteProduct = async (id) => {
  await api.delete(`/products/${id}`)
}

// Region API
export const getRegions = async () => {
  const response = await api.get('/regions')
  return response.data
}

export const getRegion = async (id) => {
  const response = await api.get(`/regions/${id}`)
  return response.data
}

export const createRegion = async (region) => {
  const response = await api.post('/regions', region)
  return response.data
}

export const updateRegion = async (id, region) => {
  const response = await api.put(`/regions/${id}`, region)
  return response.data
}

export const deleteRegion = async (id) => {
  await api.delete(`/regions/${id}`)
}

// Country API
export const getCountries = async () => {
  const response = await api.get('/countries')
  return response.data
}

// Addon API
export const getAddons = async (insuranceCode = null) => {
  const url = insuranceCode 
    ? `/addons?insurance_code=${insuranceCode}`
    : '/addons'
  const response = await api.get(url)
  return response.data
}

export const getAddon = async (id) => {
  const response = await api.get(`/addons/${id}`)
  return response.data
}

export const createAddon = async (addon) => {
  const response = await api.post('/addons', addon)
  return response.data
}

export const updateAddon = async (id, addon) => {
  const response = await api.put(`/addons/${id}`, addon)
  return response.data
}

export const deleteAddon = async (id) => {
  await api.delete(`/addons/${id}`)
}

// Addon Rule API
export const getAddonRules = async (productCode = null, addonCode = null) => {
  let url = '/addon-rules'
  const params = []
  if (productCode) params.push(`product_code=${productCode}`)
  if (addonCode) params.push(`addon_code=${addonCode}`)
  if (params.length > 0) url += '?' + params.join('&')
  const response = await api.get(url)
  return response.data
}

export const getAddonRule = async (id) => {
  const response = await api.get(`/addon-rules/${id}`)
  return response.data
}

export const createAddonRule = async (addonRule) => {
  const response = await api.post('/addon-rules', addonRule)
  return response.data
}

export const createAddonRulesBatch = async (addonRules) => {
  const response = await api.post('/addon-rules/batch', addonRules)
  return response.data
}

export const deleteAddonRule = async (id) => {
  await api.delete(`/addon-rules/${id}`)
}

// Insurance Product Addon Mapping API
export const createInsuranceProductAddonMappingsBatch = async (mappings) => {
  const response = await api.post('/insurance-product-addon-mappings/batch', mappings)
  return response.data
}

