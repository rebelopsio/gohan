Feature: Theme Switching
  As a Gohan user
  I want to switch between themes
  So that I can change my desktop's appearance quickly

  Background:
    Given the theme system is initialized
    And the "mocha" theme is active
    And my desktop environment is configured with:
      | component       |
      | window manager  |
      | status bar      |
      | terminal        |

  Scenario: Change to a different theme safely
    When I switch to the "latte" theme
    Then my desktop should display the "latte" theme
    And my previous settings should be saved for later
    And I should see confirmation of the change

  Scenario: Theme switch shows progress
    When I switch to the "frappe" theme
    Then I should see the theme being applied
    And I should be notified when it's complete

  Scenario: Theme switch saves previous settings safely
    When I switch to the "macchiato" theme
    Then my previous settings should be saved first
    And all affected settings should be included in the save
    And I should be able to restore them later if needed

  Scenario: Theme switch updates all desktop components
    When I switch to the "gohan" theme
    Then my window manager should display the "gohan" theme
    And my status bar should display the "gohan" theme
    And my terminal should display the "gohan" theme

  Scenario: Failed theme switch does not leave partial state
    Given the status bar cannot be updated
    When I attempt to switch to the "latte" theme
    Then the theme switch should fail
    And the active theme should still be "mocha"
    And my original settings should be restored
    And I should be notified of the problem

  Scenario: Switch to already active theme
    Given the "mocha" theme is active
    When I switch to the "mocha" theme
    Then I should be informed the theme is already active
    And no configuration changes should be made
    And no backup should be created

  Scenario: Switch to non-existent theme
    When I attempt to switch to the "nonexistent" theme
    Then I should receive an error
    And the error should indicate the theme does not exist
    And the active theme should remain unchanged

  Scenario: Theme changes appear instantly
    When I switch to the "frappe" theme
    Then I should see the new theme immediately
    And my applications should not need to restart

  Scenario: Theme switch preserves user customizations
    Given I have custom keybindings in my window manager
    When I switch to the "latte" theme
    Then my custom keybindings should be preserved
    And only appearance settings should change

  Scenario: Multiple rapid theme switches
    When I switch to the "latte" theme
    And I immediately switch to the "frappe" theme
    Then the final active theme should be "frappe"
    And both theme changes should have backups
    And all configurations should reflect the final theme
