Feature: Theme Preview
  As a Gohan user
  I want to preview themes before applying them
  So that I can see how they look without committing to the change

  Background:
    Given the theme system is initialized
    And the "mocha" theme is active

  Scenario: Preview a theme's colors
    When I preview the "latte" theme
    Then I should see a preview of its colors
    And the preview should show:
      | element      |
      | background   |
      | text         |
      | accents      |
      | highlights   |
      | success      |
      | errors       |

  Scenario: Preview theme with visual representation
    When I preview the "frappe" theme
    Then I should see a visual representation
    And it should show sample colors
    And it should indicate it is a "dark" theme

  Scenario: Preview shows theme information
    When I preview the "macchiato" theme
    Then I should see the display name "Catppuccin Macchiato"
    And I should see the author "Catppuccin"
    And I should see it is suitable for nighttime use
    And I should see a description

  Scenario: Preview without applying
    When I preview the "latte" theme
    Then my active theme should still be "mocha"
    And my configuration files should be unchanged
    And no backup should be created

  Scenario: Preview non-existent theme
    When I attempt to preview the "nonexistent" theme
    Then I should receive an error
    And the error should indicate the theme was not found

  Scenario: Compare multiple themes
    When I preview the "latte" theme
    And I preview the "mocha" theme
    Then I should be able to see differences in their color schemes
    And I should see "latte" is suitable for daytime use
    And I should see "mocha" is suitable for nighttime use

  Scenario: Preview shows affected components
    When I preview the "gohan" theme
    Then I should see which desktop components will be themed:
      | component        |
      | window manager   |
      | status bar       |
      | terminal         |
      | application menu |

  Scenario: Preview with detailed color information
    When I preview the "macchiato" theme with detailed output
    Then I should see hex color codes for all colors
    And I should see RGB values
    And I should see color names and purposes

  Scenario: Preview active theme
    Given the "mocha" theme is active
    When I preview the "mocha" theme
    Then I should see its colors
    And it should be marked as currently active
    And I should be informed this is the active theme
