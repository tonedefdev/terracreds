variable keyvault_name {
  type        = string
  default     = "terracreds-test-kv"
  description = "Azure Key Vault name"
}

variable location {
  type        = string
  default     = "West US2"
  description = "Resource location"
}

variable rg_name {
  type        = string
  default     = "terracreds-test-rg"
  description = "Resource group name"
}

variable test_az {
  type        = bool
  default     = true
  description = "A flag to enable testing of Azure resources"
}

variable test_aws {
  type        = bool
  default     = false
  description = "A flag to enable testing of AWS resources"
}

variable vm_name {
  type        = string
  default     = "terracreds-vm"
  description = "The name of the terracreds-test-vm"
}

variable vm_size {
  type        = string
  default     = "Standard_B2s"
  description = "Kubernetes node size"
}