variable "keyvault_name" {
  type        = string
  default     = "terracreds-test-kv"
  description = "Azure Key Vault name"
}

variable "location" {
  type        = string
  default     = "West US2"
  description = "Resource location"
}

variable "keyvault_only" {
  type        = bool
  default     = true
  description = "Create only the Azure Key Vault resources and not any VMs"
}

variable "object_id" {
  type        = string
  default     = "5277cb3b-c86d-4bdc-8034-6220bdd5b4d3"
  description = "The Object ID of the user to add to the Azure Key Vault Access Policy"
}

variable "rg_name" {
  type        = string
  default     = "terracreds-test-rg"
  description = "Resource group name"
}

variable "test_az" {
  type        = bool
  default     = false
  description = "A flag to enable testing of Azure resources"
}

variable "test_aws" {
  type        = bool
  default     = false
  description = "A flag to enable testing of AWS resources"
}

variable "vm_name" {
  type        = string
  default     = "terracreds-vm"
  description = "The name of the terracreds-test-vm"
}

variable "vm_size" {
  type        = string
  default     = "Standard_B2s"
  description = "Kubernetes node size"
}