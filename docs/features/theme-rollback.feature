Feature: Theme Rollback
  As a Gohan user
  I want to rollback theme changes
  So that I can recover if I don't like a new theme

  Background:
    Given the theme system is initialized
    And the "mocha" theme is active

  Scenario: Restore previous theme
    Given I recently changed from "mocha" to "latte"
    When I undo the theme change
    Then my desktop should display the "mocha" theme again
    And my previous settings should be restored
    And I should see confirmation of the restoration

  Scenario: Restore to specific earlier theme
    Given I have changed themes multiple times:
      | from      | to        | when                |
      | mocha     | latte     | earlier today       |
      | latte     | frappe    | a few hours ago     |
      | frappe    | macchiato | recently            |
    When I restore my appearance from a few hours ago
    Then my desktop should display the "latte" theme
    And my desktop should look as it did then

  Scenario: Restoration shows progress
    Given I recently changed from "mocha" to "gohan"
    When I undo the theme change
    Then I should see the restoration in progress
    And I should be notified when complete

  Scenario: View my theme change history
    Given I changed themes twice today:
      | from  | to     | when            |
      | mocha | latte  | 2 hours ago     |
      | latte | frappe | 30 minutes ago  |
    When I view my theme history
    Then I should see 2 previous themes I can restore
    And they should be sorted newest first
    And each should show the theme transition

  Scenario: Attempt restore with no history
    Given no theme changes have been made
    When I attempt to undo changes
    Then I should receive an error
    And the error should indicate no previous themes are available
    And the current theme should remain active

  Scenario: Restore to invalid point in history
    When I attempt to restore to a non-existent point in history
    Then I should receive an error
    And the error should indicate that point was not found
    And the current theme should remain unchanged

  Scenario: Failed restoration preserves current state
    Given I switched from "mocha" to "latte"
    And the saved settings are corrupted
    When I attempt to undo the change
    Then I should receive an error
    And the error should describe the problem
    And the "latte" theme should remain active
    And my current settings should be unchanged

  Scenario: Undo multiple theme changes sequentially
    Given I switched from "mocha" to "latte"
    And I switched from "latte" to "frappe"
    And I switched from "frappe" to "macchiato"
    When I undo the most recent change
    Then the active theme should be "frappe"
    When I undo again
    Then the active theme should be "latte"
    When I undo again
    Then the active theme should be "mocha"
