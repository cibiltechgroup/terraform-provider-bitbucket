package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccBitbucketBranchRestriction_basic(t *testing.T) {
	var branchRestrictionBefore, branchRestrictionAfter BranchRestriction

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketBranchRestrictionConfig := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-for-branch-restriction-test"
		}
		resource "bitbucket_branch_restriction" "test_repo_branch_restriction" {
			owner = "%s"
			repository = "${bitbucket_repository.test_repo.name}"
 			kind = "force"
 			pattern = "master"
		}
	`, testUser, testUser)

	testAccBitbucketBranchRestrictionConfig_removed := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-for-branch-restriction-test"
		}
	`, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketBranchRestrictionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketBranchRestrictionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketBranchRestrictionExists("bitbucket_branch_restriction.test_repo_branch_restriction", &branchRestrictionBefore),
				),
			},
			{
				Config: testAccBitbucketBranchRestrictionConfig_removed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketBranchRestrictionExists("bitbucket_branch_restriction.test_repo_branch_restriction", &branchRestrictionAfter),
				),
			},
		},
	})
}

func testAccCheckBitbucketBranchRestrictionDestroy(s *terraform.State) error {
	_, branch_restriction := s.RootModule().Resources["bitbucket_branch_restriction.test_repo_branch_restriction"]
	if branch_restriction {
		return fmt.Errorf("Found %s", "bitbucket_branch_restriction.test_repo_branch_restriction")
	}

	_, repository := s.RootModule().Resources["bitbucket_repository.test_repo"]
	if repository {
		return fmt.Errorf("Found %s", "bitbucket_repository.test_repo")
	}

	return nil
}

func testAccCheckBitbucketBranchRestrictionExists(n string, branchRestriction *BranchRestriction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No BranchRestriction ID is set")
		}
		return nil
	}
}
