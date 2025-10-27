Feature: Theme Management
  As a Gohan user
  I want to manage visual themes
  So that I can customize my Hyprland environment's appearance

  Background:
    Given the theme system is initialized
    And the following themes are available:
      | name       | variant | author      |
      | mocha      | dark    | Catppuccin  |
      | latte      | light   | Catppuccin  |
      | frappe     | dark    | Catppuccin  |
      | macchiato  | dark    | Catppuccin  |
      | gohan      | dark    | Gohan Team  |

  Scenario: List available themes
    When I view available themes
    Then I should see 5 themes
    And each theme should have a name
    And each theme should indicate if it's suitable for day or night use
    And each theme should show its creator

  Scenario: Identify active theme
    Given the "mocha" theme is active
    When I request a list of available themes
    Then the "mocha" theme should be marked as active
    And all other themes should not be marked as active

  Scenario: View theme information
    When I view the "latte" theme
    Then I should see it is a light theme
    And I should see it was created by "Catppuccin"
    And I should see a preview of its colors

  Scenario: Attempt to get non-existent theme
    When I request details for the "nonexistent" theme
    Then I should receive an error
    And the error should indicate the theme was not found

  Scenario: Find themes suitable for nighttime use
    When I look for dark themes
    Then I should see 4 themes
    And all themes should be suitable for low-light environments

  Scenario: Find themes suitable for daytime use
    When I look for light themes
    Then I should see 1 theme
    And the theme should be "latte"

  Scenario: Get active theme information
    Given the "frappe" theme is active
    When I request the active theme
    Then I should receive the "frappe" theme
    And it should be marked as active

  Scenario: System has default theme when none set
    Given no theme has been explicitly set
    When I request the active theme
    Then I should receive the "mocha" theme
    And it should be marked as active
