scenario "e2e_vault" {
  terraform_cli = terraform_cli.default
  terraform     = terraform.default
  providers = [
    provider.aws.default,
    provider.enos.default
  ]

  matrix {
    builder = ["local", "crt"]
  }

  locals {
    aws_ssh_private_key_path = abspath(var.aws_ssh_private_key_path)
    boundary_install_dir     = abspath(var.boundary_install_dir)
    local_boundary_dir       = abspath(var.local_boundary_dir)
    build_path = {
      "local" = "/tmp",
      "crt"   = var.crt_bundle_path == null ? null : abspath(var.crt_bundle_path)
    }
  }

  step "check_jq_is_installed" {
    module = module.binary_finder

    variables {
      name = "jq"
    }
  }

  step "find_azs" {
    module = module.az_finder

    variables {
      instance_type = [
        var.worker_instance_type,
        var.controller_instance_type
      ]
    }
  }

  step "create_db_password" {
    module = module.random_stringifier
  }

  step "build_boundary" {
    module = matrix.builder == "crt" ? module.build_crt : module.build_local

    variables {
      path = local.build_path[matrix.builder]
    }
  }

  step "create_base_infra" {
    module = module.infra

    variables {
      availability_zones = step.find_azs.availability_zones
    }
  }

  step "create_boundary_cluster" {
    module = module.boundary
    depends_on = [
      step.create_base_infra,
      step.build_boundary
    ]

    variables {
      boundary_install_dir     = local.boundary_install_dir
      controller_instance_type = var.controller_instance_type
      controller_count         = var.controller_count
      db_pass                  = step.create_db_password.string
      kms_key_arn              = step.create_base_infra.kms_key_arn
      local_artifact_path      = step.build_boundary.artifact_path
      ubuntu_ami_id            = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      vpc_id                   = step.create_base_infra.vpc_id
      worker_count             = var.worker_count
      worker_instance_type     = var.worker_instance_type
    }
  }

  step "create_vault_cluster" {
    module = module.vault
    depends_on = [
      step.create_base_infra,
    ]

    variables {
      ami_id            = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      instance_type     = var.vault_instance_type
      instance_count    = 1
      kms_key_arn       = step.create_base_infra.kms_key_arn
      storage_backend   = "raft"
      sg_additional_ips = step.create_boundary_cluster.controller_ips
      unseal_method     = "awskms"
      vault_release = {
        version = "1.11.0"
        edition = "oss"
      }
      vpc_id = step.create_base_infra.vpc_id
    }
  }

  step "create_target" {
    module     = module.target
    depends_on = [step.create_base_infra]

    variables {
      ami_id               = step.create_base_infra.ami_ids["ubuntu"]["amd64"]
      aws_ssh_keypair_name = var.aws_ssh_keypair_name
      enos_user            = var.enos_user
      instance_type        = var.target_instance_type
      vpc_id               = step.create_base_infra.vpc_id
      target_count         = 1
    }
  }

  step "run_e2e_test" {
    module = module.test_e2e
    depends_on = [
      step.create_boundary_cluster,
      step.create_target,
      step.create_vault_cluster
    ]

    variables {
      test_package             = "github.com/hashicorp/boundary/testing/internal/e2e/credential/vault"
      alb_boundary_api_addr    = step.create_boundary_cluster.alb_boundary_api_addr
      auth_method_id           = step.create_boundary_cluster.auth_method_id
      auth_login_name          = step.create_boundary_cluster.auth_login_name
      auth_password            = step.create_boundary_cluster.auth_password
      local_boundary_dir       = local.local_boundary_dir
      aws_ssh_private_key_path = local.aws_ssh_private_key_path
      target_ips               = step.create_target.target_ips
      target_user              = "ubuntu"
      vault_addr               = step.create_vault_cluster.instance_public_ips[0]
      vault_root_token         = step.create_vault_cluster.vault_root_token
    }
  }

  output "test_results" {
    value = step.run_e2e_test.test_results
  }
}
