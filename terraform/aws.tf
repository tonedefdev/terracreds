resource "aws_vpc" "vpc" {
  count      = var.test_aws ? 1 : 0
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "vm_subnet" {
  count             = var.test_aws ? 1 : 0
  availability_zone = "us-west-2a"
  vpc_id            = aws_vpc.vpc[count.index].id
  cidr_block        = "10.0.1.0/24"

  tags = {
    Name = "vm_subnet"
  }
}

resource "aws_security_group" "allow_ssh" {
  count       = var.test_aws ? 1 : 0
  name        = "allow_ssh"
  description = "Allow SSH inbound traffic"
  vpc_id      = aws_vpc.vpc[count.index].id

  ingress = [
    {
      description      = "SSH from Internet"
      cidr_blocks      = ["0.0.0.0/0"]
      from_port        = 22
      to_port          = 22
      ipv6_cidr_blocks = ["::/0"]
      prefix_list_ids  = []
      protocol         = "tcp"
      security_groups  = []
      self             = false
    }
  ]

  tags = {
    Name = "allow_ssh"
  }
}

resource "aws_key_pair" "test_vm_keys" {
  count      = var.test_aws ? 1 : 0
  key_name   = "test_vm_keys"
  public_key = file("C:\\Temp\\aws_rsa.pub")
}

resource "aws_instance" "test_vm" {
  count                       = var.test_aws ? 1 : 0
  ami                         = "ami-005e54dee72cc1d00"
  associate_public_ip_address = true
  availability_zone           = "us-west-2a"
  instance_type               = "t2.micro"
  key_name                    = aws_key_pair.test_vm_keys[count.index].key_name
  subnet_id                   = aws_subnet.vm_subnet[0].id

  vpc_security_group_ids = [
    aws_security_group.allow_ssh[0].id
  ]
}