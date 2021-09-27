package route53resolver_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfroute53resolver "github.com/hashicorp/terraform-provider-aws/internal/service/route53resolver"
	tfroute53resolver "github.com/hashicorp/terraform-provider-aws/internal/service/route53resolver"
)

func init() {
	resource.AddTestSweepers("aws_route53_resolver_firewall_rule", &resource.Sweeper{
		Name: "aws_route53_resolver_firewall_rule",
		F:    testSweepRoute53ResolverFirewallRules,
		Dependencies: []string{
			"aws_route53_resolver_firewall_rule_association",
		},
	})
}

func testSweepRoute53ResolverFirewallRules(region string) error {
	client, err := acctest.SharedRegionalSweeperClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).Route53ResolverConn
	var sweeperErrs *multierror.Error

	err = conn.ListFirewallRulesPages(&route53resolver.ListFirewallRulesInput{}, func(page *route53resolver.ListFirewallRulesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, firewallRule := range page.FirewallRules {
			id := tfroute53resolver.FirewallRuleCreateID(*firewallRule.FirewallRuleGroupId, *firewallRule.FirewallDomainListId)

			log.Printf("[INFO] Deleting Route53 Resolver DNS Firewall rule: %s", id)
			r := ResourceFirewallRule()
			d := r.Data(nil)
			d.SetId(id)
			err := r.Delete(d, client)

			if err != nil {
				log.Printf("[ERROR] %s", err)
				sweeperErrs = multierror.Append(sweeperErrs, err)
				continue
			}
		}

		return !lastPage
	})
	if acctest.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Route53 Resolver DNS Firewall rules sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Route53 Resolver DNS Firewall rules: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSRoute53ResolverFirewallRule_basic(t *testing.T) {
	var v route53resolver.FirewallRule
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53_resolver_firewall_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallRuleExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttrPair(resourceName, "firewall_rule_group_id", "aws_route53_resolver_firewall_rule_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "firewall_domain_list_id", "aws_route53_resolver_firewall_domain_list.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "priority", "100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute53ResolverFirewallRule_block(t *testing.T) {
	var v route53resolver.FirewallRule
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53_resolver_firewall_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallRuleConfig_block(rName, "NODATA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallRuleExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "block_response", "NODATA"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute53ResolverFirewallRule_blockOverride(t *testing.T) {
	var v route53resolver.FirewallRule
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53_resolver_firewall_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallRuleConfig_blockOverride(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallRuleExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "block_override_dns_type", "CNAME"),
					resource.TestCheckResourceAttr(resourceName, "block_override_domain", "example.com."),
					resource.TestCheckResourceAttr(resourceName, "block_override_ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "block_response", "OVERRIDE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute53ResolverFirewallRule_disappears(t *testing.T) {
	var v route53resolver.FirewallRule
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53_resolver_firewall_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallRuleExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceFirewallRule(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRoute53ResolverFirewallRuleDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53_resolver_firewall_rule" {
			continue
		}

		// Try to find the resource
		_, err := tfroute53resolver.FindFirewallRuleByID(conn, rs.Primary.ID)
		// Verify the error is what we want
		if tfawserr.ErrMessageContains(err, route53resolver.ErrCodeResourceNotFoundException, "") {
			continue
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("Route 53 Resolver DNS Firewall rule still exists: %s", rs.Primary.ID)
	}

	return nil
}

func testAccCheckRoute53ResolverFirewallRuleExists(n string, v *route53resolver.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route 53 Resolver DNS Firewall rule ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn
		out, err := tfroute53resolver.FindFirewallRuleByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		*v = *out

		return nil
	}
}

func testAccRoute53ResolverFirewallRuleConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_rule_group" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_rule" "test" {
  name                    = %[1]q
  action                  = "ALLOW"
  firewall_rule_group_id  = aws_route53_resolver_firewall_rule_group.test.id
  firewall_domain_list_id = aws_route53_resolver_firewall_domain_list.test.id
  priority                = 100
}
`, rName)
}

func testAccRoute53ResolverFirewallRuleConfig_block(rName, blockResponse string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_rule_group" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_rule" "test" {
  name                    = %[1]q
  action                  = "BLOCK"
  block_response          = %[2]q
  firewall_rule_group_id  = aws_route53_resolver_firewall_rule_group.test.id
  firewall_domain_list_id = aws_route53_resolver_firewall_domain_list.test.id
  priority                = 100
}
`, rName, blockResponse)
}

func testAccRoute53ResolverFirewallRuleConfig_blockOverride(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_rule_group" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
}

resource "aws_route53_resolver_firewall_rule" "test" {
  name                    = %[1]q
  action                  = "BLOCK"
  block_override_dns_type = "CNAME"
  block_override_domain   = "example.com."
  block_override_ttl      = 60
  block_response          = "OVERRIDE"
  firewall_rule_group_id  = aws_route53_resolver_firewall_rule_group.test.id
  firewall_domain_list_id = aws_route53_resolver_firewall_domain_list.test.id
  priority                = 100
}
`, rName)
}