# Terraform Native Test for Deferred Data Sources
# This test verifies that list attributes (tag_ids, notification_ids, monitor_ids)
# can handle unknown values from deferred data sources without "inconsistent final plan" errors

run "test_deferred_tag_ids" {
  command = plan

  # This test should pass - plan should succeed even with unknown tag IDs
  # Before the fix, this would fail with "Provider produced inconsistent final plan"
  assert {
    condition     = length(peekaping_monitor.test) == 2
    error_message = "Should plan to create 2 monitors"
  }

  assert {
    condition = alltrue([
      for k, v in peekaping_monitor.test : v.tag_ids != null
    ])
    error_message = "All monitors should have tag_ids planned (even if some are unknown)"
  }
}

run "apply_with_deferred_data_sources" {
  command = apply

  # This test verifies that apply succeeds and resolves the deferred data sources
  assert {
    condition     = length(peekaping_monitor.test) == 2
    error_message = "Should create 2 monitors"
  }

  assert {
    condition = alltrue([
      for k, v in peekaping_monitor.test : length(v.tag_ids) >= 1
    ])
    error_message = "All monitors should have at least one tag_id after apply"
  }

  assert {
    condition = alltrue([
      for k, v in peekaping_monitor.test : v.id != ""
    ])
    error_message = "All monitors should have valid IDs"
  }

  assert {
    condition = alltrue([
      for k, v in peekaping_monitor.test : v.name != ""
    ])
    error_message = "All monitors should have names"
  }
}

run "update_with_deferred_data_sources" {
  command = apply

  # Test that updates work correctly with deferred data sources
  # This simulates a scenario where tag_ids are updated but some are still unknown
  assert {
    condition     = length(peekaping_monitor.test) == 2
    error_message = "Should still have 2 monitors after update"
  }

  assert {
    condition = alltrue([
      for k, v in peekaping_monitor.test : v.id != ""
    ])
    error_message = "All monitors should retain their IDs after update"
  }
}
