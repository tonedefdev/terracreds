data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "rg" {
  count    = var.keyvault_only || var.test_az ? 1 : 0
  name     = var.rg_name
  location = var.location
}

resource "azurerm_public_ip" "pip" {
  count               = var.test_az ? 1 : 0
  name                = "terracreds-test-pip"
  location            = azurerm_resource_group.rg[0].location
  resource_group_name = azurerm_resource_group.rg[0].name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network" "vnet" {
  count               = var.test_az ? 1 : 0
  name                = "${var.vm_name}-vnet"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.rg[0].location
  resource_group_name = azurerm_resource_group.rg[0].name
}

resource "azurerm_subnet" "subnet" {
  count                = var.test_az ? 1 : 0
  name                 = "${var.vm_name}-subnet"
  resource_group_name  = azurerm_resource_group.rg[0].name
  virtual_network_name = azurerm_virtual_network.vnet[count.index].name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_network_interface" "nic" {
  count               = var.test_az ? 1 : 0
  name                = "${var.vm_name}-nic"
  location            = azurerm_resource_group.rg[0].location
  resource_group_name = azurerm_resource_group.rg[0].name

  ip_configuration {
    name                          = "external"
    public_ip_address_id          = azurerm_public_ip.pip[count.index].id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet[count.index].id
  }
}

resource "azurerm_windows_virtual_machine" "vm" {
  count               = var.test_az ? 1 : 0
  name                = var.vm_name
  location            = azurerm_resource_group.rg[0].location
  resource_group_name = azurerm_resource_group.rg[0].name
  size                = var.vm_size
  admin_username      = "terraadmin"
  admin_password      = "TerracredsR0cks!" // Yes, I know this is bad form, but it's only a test machine :P
  network_interface_ids = [
    azurerm_network_interface.nic[count.index].id,
  ]

  identity {
    type = "SystemAssigned"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }
}

resource "azurerm_key_vault" "keyvault" {
  count                       = var.keyvault_only || var.test_az ? 1 : 0
  name                        = var.keyvault_name
  location                    = azurerm_resource_group.rg[0].location
  resource_group_name         = azurerm_resource_group.rg[0].name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = var.keyvault_only ? var.object_id : azurerm_windows_virtual_machine.vm[count.index].identity[count.index].principal_id

    secret_permissions = [
      "Get",
      "List",
      "Set",
      "Delete"
    ]
  }
}