Feature: Component Hot Reload
  As a Gohan user
  I want components to automatically reload when themes change
  So that I see theme changes immediately without manual restarts

  Background:
    Given the theme system is initialized
    And Hyprland is running

  Scenario: Hyprland reloads after theme change
    Given the current theme is "mocha"
    When I apply the "latte" theme
    Then Hyprland should reload its configuration
    And the new theme colors should be visible

  Scenario: Waybar reloads after theme change
    Given the current theme is "mocha"
    And Waybar is running
    When I apply the "latte" theme
    Then Waybar should restart
    And the new theme should be applied to Waybar

  Scenario: Multiple components reload together
    Given the current theme is "mocha"
    And multiple components are running:
      | component |
      | Hyprland  |
      | Waybar    |
    When I apply the "frappe" theme
    Then all running components should reload
    And the reload should happen in the correct order

  Scenario: Reload continues despite component failures
    Given the current theme is "mocha"
    And Waybar is not running
    When I apply the "latte" theme
    Then Hyprland should still reload successfully
    And I should see a warning about Waybar
    And the theme should still be applied

  Scenario: Reload happens automatically
    Given the current theme is "mocha"
    When I apply the "latte" theme
    Then I should not need to manually reload any components
    And all changes should be visible immediately

  Scenario: Rollback triggers component reload
    Given I have changed from "mocha" to "latte"
    When I rollback the theme
    Then components should reload with the previous theme
    And the "mocha" theme should be visible

  Scenario: Dry-run mode skips reload
    Given I am running in dry-run mode
    When I apply the "latte" theme
    Then no components should be reloaded
    But the theme files should still be updated
