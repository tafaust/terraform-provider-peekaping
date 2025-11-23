# Terraform Native Testing for Deferred Data Sources
# This file tests the fix for "inconsistent final plan" errors when list attributes
# reference data sources with depends_on or other deferred values

run "monitor_with_deferred_tag_ids" {
  command = plan

  # This test verifies that plan succeeds even when tag_ids contains
  # references to data sources that are unknown during plan time
  # The key test is that the plan completes without "inconsistent final plan" errors
  assert {
    condition     = length(peekaping_monitor.test) > 0
    error_message = "Plan should succeed without 'inconsistent final plan' errors when tag_ids references deferred data sources"
  }

  assert {
    condition     = peekaping_monitor.test["one"].name == "Test Monitor"
    error_message = "Monitor should be planned correctly"
  }
}

run "status_page_with_deferred_domains" {
  command = plan

  # This test verifies that plan succeeds even when domains contains
  # references to data sources that are unknown during plan time
  # The key test is that the plan completes without "inconsistent final plan" errors
  assert {
    condition     = peekaping_status_page.with_deferred_domains.title == "Status Page with Deferred Domains"
    error_message = "Plan should succeed without 'inconsistent final plan' errors when domains references deferred data sources"
  }

  assert {
    condition     = length(peekaping_status_page.with_deferred_domains.domains) >= 1
    error_message = "Status page should have domains planned"
  }
}

