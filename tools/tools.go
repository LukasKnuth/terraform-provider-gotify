// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

import (
	// Documentation generation
	// We do this here because otherwise `go mod tidy` will get rid of it and `go generate` will fail
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
